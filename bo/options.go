// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package bo

// OptimizerOption sets an option on the optimizer.
type OptimizerOption func(*Optimizer)

// WithOutputName sets the outputs name. Only really matters if you're planning
// on using gp/plot.
func WithOutputName(name string) OptimizerOption {
	return func(o *Optimizer) {
		o.updateNames(name)
	}
}

// WithRandomRounds sets the number of random rounds to run.
func WithRandomRounds(rounds int) OptimizerOption {
	return func(o *Optimizer) {
		o.mu.randomRounds = rounds
	}
}

// WithRounds sets the total number of rounds to run.
func WithRounds(rounds int) OptimizerOption {
	return func(o *Optimizer) {
		o.mu.rounds = rounds
	}
}

// WithExploration sets the exploration function to use.
func WithExploration(exploration Exploration) OptimizerOption {
	return func(o *Optimizer) {
		o.mu.exploration = exploration
	}
}

// WithMinimize sets whether or not to minimize. Passing false, maximizes
// instead.
func WithMinimize(minimize bool) OptimizerOption {
	return func(o *Optimizer) {
		o.mu.minimize = minimize
	}
}

// WithBarrierFunc sets the barrier function to use.
func WithBarrierFunc(bf BarrierFunc) OptimizerOption {
	return func(o *Optimizer) {
		o.mu.barrierFunc = bf
	}
}
