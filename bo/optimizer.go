// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package bo

import (
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"gonum.org/v1/gonum/optimize"
)

const (
	// DefaultRounds is the default number of rounds to run.
	DefaultRounds = 20
	// DefaultRandomRounds is the default number of random rounds to run.
	DefaultRandomRounds = 5
	// DefaultMinimize is the default value of minimize.
	DefaultMinimize = true

	// NumRandPoints is the maximum allowed number of evaluations.
	NumRandPoints = 100000
	// NumGradPoints is the number of random points of gradient descent
	NumGradPoints = 256
)

var (
	// DefaultExploration uses UCB with 95 confidence interval.
	DefaultExploration = UCB{Kappa: 1.96}
	// DefaultBarrierFunc sets the default barrier function to use.
	DefaultBarrierFunc = LogBarrier{}
)

// Optimizer is a blackbox gaussian process optimizer.
type Optimizer struct {
	mu struct {
		sync.Mutex
		gp                          *GP
		params                      []Param
		round, randomRounds, rounds int
		exploration                 Exploration
		minimize                    bool
		barrierFunc                 BarrierFunc

		explorationErr error
	}
	running uint32 // atomic, 0 false 1 true
}

// NewOptimizer creates a new optimizer with the specified optimizable parameters and
// options.
func NewOptimizer(params []Param, opts ...OptimizerOption) *Optimizer {
	o := &Optimizer{}
	o.mu.gp = NewGP(MaternCov{}, 0)
	o.mu.params = params

	// Set default values.
	o.mu.randomRounds = DefaultRandomRounds
	o.mu.rounds = DefaultRounds
	o.mu.exploration = DefaultExploration
	o.mu.minimize = DefaultMinimize
	o.mu.barrierFunc = DefaultBarrierFunc

	o.updateNames("")

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// updateNames sets the gaussian process names.
func (o *Optimizer) updateNames(outputName string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	var inputNames []string
	for _, p := range o.mu.params {
		inputNames = append(inputNames, p.GetName())
	}
	o.mu.gp.SetNames(inputNames, outputName)
}

// GP returns the underlying gaussian process. Primary for use with plotting
// behavior.
func (o *Optimizer) GP() *GP {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.mu.gp
}

func sampleParams(params []Param) []float64 {
	x := make([]float64, len(params))
	for i, p := range params {
		x[i] = p.Sample()
	}
	return x
}

func sampleParamsMap(params []Param) map[Param]float64 {
	x := map[Param]float64{}
	for i, v := range sampleParams(params) {
		x[params[i]] = v
	}
	return x
}

type randerFunc func([]float64) []float64

func (f randerFunc) Rand(x []float64) []float64 {
	return f(x)
}

func isFatalErr(err error) bool {
	if err == nil {
		return false
	}

	// Only recurse 100 times before breaking.
	for i := 0; i < 100; i++ {
		parent := errors.Cause(err)
		if parent == err {
			break
		}
		err = parent
	}

	if _, ok := err.(optimize.ErrFunc); ok {
		return false
	}
	switch err {
	case optimize.ErrLinesearcherFailure, optimize.ErrNoProgress:
		return false
	default:
		return true
	}
}

// Next returns the next best x values to explore. If more than rounds have
// elapsed, nil is returned. If parallel is true, that round can happen in
// parallel to other rounds.
func (o *Optimizer) Next() (x map[Param]float64, parallel bool, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Return if we've exceeded max # of rounds, or if there was an error while
	// doing exploration which is likely caused by numerical precision errors.
	if o.mu.round >= o.mu.rounds || o.mu.explorationErr != nil {
		return nil, false, nil
	}

	// If we don't have enough random rounds, run more.
	if o.mu.round < o.mu.randomRounds {
		x = sampleParamsMap(o.mu.params)
		o.mu.round++
		// Don't return parallel on the last random round.
		return x, o.mu.round != o.mu.randomRounds, nil
	}
	return o.eval()
}

func (o *Optimizer) eval() (x map[Param]float64, parallel bool, err error) {

	var fErr error
	f := func(x []float64) float64 {
		v, err := o.mu.exploration.Estimate(o.mu.gp, o.mu.minimize, x)
		if err != nil {
			fErr = errors.Wrap(err, "exploration error")
		}

		if o.mu.minimize {
			return v
		}
		return -v
	}
	problem := optimize.Problem{
		Func: f,
		Grad: func(grad, x []float64) {
			g, err := o.mu.gp.Gradient(x)
			if err != nil {
				fErr = errors.Wrap(err, "gradient error")
			}
			copy(grad, g)
		},
	}

	// Randomly query a bunch of points to get a good estimate of maximum.
	result, err := optimize.Minimize(problem, make([]float64, len(o.mu.params)), &optimize.Settings{
		FuncEvaluations: NumRandPoints,
	}, &optimize.GuessAndCheck{
		Rander: randerFunc(func(x []float64) []float64 {
			return sampleParams(o.mu.params)
		}),
	})
	if err != nil {
		return nil, false, errors.Wrapf(err, "random sample failed")
	}
	if fErr != nil {
		o.mu.explorationErr = fErr
	}
	min := result.F
	minX := result.X

	// Run gradient descent on the best point.
	method := optimize.LBFGS{}
	grad := BoundsMethod{
		Method: &method,
		Bounds: o.mu.params,
	}
	// TODO: Bounded line searcher.
	{
		result, err := optimize.Minimize(problem, minX, nil, grad)
		if isFatalErr(err) {
			o.mu.explorationErr = errors.Wrapf(err, "random sample optimize failed")
		}
		if fErr != nil {
			o.mu.explorationErr = fErr
		}
		if result != nil && result.F < min {
			min = result.F
			minX = result.X
		}
	}

	// Attempt to use gradient descent on random points.
	for i := 0; i < NumGradPoints; i++ {
		x := sampleParams(o.mu.params)
		result, err := optimize.Minimize(problem, x, nil, grad)
		if isFatalErr(err) {
			o.mu.explorationErr = errors.Wrapf(err, "gradient descent failed: i %d, x %+v, result%+v", i, x, result)
		}
		if fErr != nil {
			o.mu.explorationErr = fErr
		}
		if result != nil && result.F < min {
			min = result.F
			minX = result.X
		}
	}

	if o.mu.explorationErr != nil {
		return nil, false, nil
	}

	m := map[Param]float64{}
	for i, x := range minX {
		m[o.mu.params[i]] = x
	}

	o.mu.round++
	return m, false, nil
}

// ExplorationErr returns the error of exploration
func (o *Optimizer) ExplorationErr() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.mu.explorationErr
}

// Log adds given x and y to the gaussian process
func (o *Optimizer) Log(x map[Param]float64, y float64) {
	o.mu.Lock()
	defer o.mu.Unlock()

	var xa []float64
	for _, p := range o.mu.params {
		xa = append(xa, x[p])
	}
	o.mu.gp.Add(xa, y)
}

// Predict consumes all historical X and Y and predict the next x
func (o *Optimizer) Predict(X []map[Param]float64, Y []float64) (x map[Param]float64, err error) {
	for i := range X {
		o.Log(X[i], Y[i])
	}
	x, _, err = o.eval()
	if x == nil || err != nil {
		return nil, errors.Wrapf(err, "failed to get next point")
	}
	return
}

// RunSerial will call f sequentially without parallelism.
// It blocks until all rounds have elapsed, or Stop is called.
func (o *Optimizer) RunSerial(f func(map[Param]float64) float64) (x map[Param]float64, y float64, err error) {
	for {
		status := atomic.LoadUint32(&o.running)
		if status == 1 {
			return nil, 0, errors.New("optimizer is already running")
		}
		if atomic.CompareAndSwapUint32(&o.running, status, 1) {
			break
		}
	}

	for {
		if !o.Running() {
			return nil, 0, errors.New("optimizer got stop signal")
		}

		x, _, err := o.Next()
		if err != nil {
			return nil, 0, errors.Wrapf(err, "failed to get next point")
		}
		if x == nil {
			break
		}
		o.Log(x, f(x))
	}

	atomic.StoreUint32(&o.running, 0)

	var xa []float64
	if o.mu.minimize {
		xa, y = o.mu.gp.Minimum()
	} else {
		xa, y = o.mu.gp.Maximum()
	}
	x = map[Param]float64{}
	for i, v := range xa {
		x[o.mu.params[i]] = v
	}

	return x, y, nil
}

// Run will call f the fewest times as possible while trying to maximize
// the output value. It blocks until all rounds have elapsed, or Stop is called.
func (o *Optimizer) Run(f func(map[Param]float64) float64) (x map[Param]float64, y float64, err error) {
	for {
		status := atomic.LoadUint32(&o.running)
		if status == 1 {
			return nil, 0, errors.New("optimizer is already running")
		}
		if atomic.CompareAndSwapUint32(&o.running, status, 1) {
			break
		}
	}

	var wg sync.WaitGroup
	for {
		if !o.Running() {
			return nil, 0, errors.New("optimizer got stop signal")
		}

		x, parallel, err := o.Next()
		if err != nil {
			return nil, 0, errors.Wrapf(err, "failed to get next point")
		}
		if x == nil {
			break
		}
		if parallel {
			wg.Add(1)
			go func() {
				defer wg.Done()

				o.Log(x, f(x))
			}()
		} else {
			wg.Wait()
			o.Log(x, f(x))
		}
	}

	atomic.StoreUint32(&o.running, 0)

	var xa []float64
	if o.mu.minimize {
		xa, y = o.mu.gp.Minimum()
	} else {
		xa, y = o.mu.gp.Maximum()
	}
	x = map[Param]float64{}
	for i, v := range xa {
		x[o.mu.params[i]] = v
	}

	return x, y, nil
}

// Stop stops Optimize.
func (o *Optimizer) Stop() {
	atomic.StoreUint32(&o.running, 0)
}

// Running returns whether or not the optimizer is running.
func (o *Optimizer) Running() bool {
	return atomic.LoadUint32(&o.running) == 1
}

// Rounds is the number of rounds that have been run.
func (o *Optimizer) Rounds() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.mu.round
}
