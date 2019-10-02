// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package rp implements a random process where identifies peak load.
//
// 1. Identify peak load
// The current approach:  Mann–Whitney U test
//
//   Consider a count process: N1= (n_1, ..., n_{p/2}), N2 = (n_{p/2}, ..., n_{p})
//   and our goal is to check if N2 is significant higher than N1.
//   Let hypothesis: H0: \mu_{N_1} \leq \mu_{N_2}, and H1: \mu_{N_1} > \mu_{N_2}
//   Then rejection region is:
//
//     z = \left|\frac{\mu_{N_2} - \mu_{N_1}}{\sigma_{N_1}}\right| \geq z(c)
//
//   We reject H0 under given confidence level.
//   Since this approach assume the count process obey normal dist, it is
//   recommended to have a window size higher than 15.
//
//
// 2. Check if system resource can handle next peak load: Poisson + Moving average
// The current approach: Markov Poisson process + Moving average
//
//   Consider p windows, and each number of events obey Poisson process,
//   then \lambda = \frac{1}{p} \sum_{i=1}^{p} n_i, and in the next window,
//   the probability of having k number of events is:
//
//     P(N=k) = \frac{\lambda^k \exp{-\lambda}}{k!}, k = 0, 1, ...
//
//   Assume a system can handle Q requests, then we have:
//
//     P(N\leq Q) = \sum_{i=0}^{Q} P(N=i) = \sum_{i=0}^{Q} \frac{\lambda^k \exp{\left(-\lambda\right)}}{k!}
//
package rp

import (
	"math"
	"sync"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat/distuv"
)

// RandomProcess defines two type of tests for given observations.
type RandomProcess interface {
	Store(nevent float64)
	Significant() bool
	Acceptable(k float64) (int64, bool)
}

// CountProcess caches historical number of events,
// which provides significant test and poission process evaluation.
type countProcess struct {
	maxsize    float64
	confidence float64 // for z test and poission process
	nevents    []float64

	mu sync.RWMutex
}

// NewCountProcess creates a new counting random process.
// maxsize represents the observation window size of the count process.
// confidenceLevel gives the confidence level of the significant test
// and acceptance test, it is recommended to set below 0.05.
func NewCountProcess(maxsize, confidenceLevel float64) RandomProcess {
	return &countProcess{
		maxsize:    maxsize,
		nevents:    []float64{},
		confidence: confidenceLevel,
	}
}

// Store stores given number of events
func (w *countProcess) Store(nevent float64) {
	w.mu.Lock()
	if float64(len(w.nevents)) >= w.maxsize {
		w.nevents = append(w.nevents[1:], nevent)
	} else {
		w.nevents = append(w.nevents, nevent)
	}
	w.mu.Unlock()
}

// Significant performs a significant greater test,
// here we implement by z-test.
//
// a historical window of #events [n_1, n_2, ..., n_p],
// we compare: [n_1, ..., n_{p/2}] and [n_{p/2}, ..., n_{p}]
// through ztest.
//
// confidence coefficient is configurable
func (w *countProcess) Significant() bool {
	return w.ztest()
}

// Acceptable calculates whether the system can accept k events
// on next interval, here we implement by poission process.
//
// Since we have tested half of the current window is significant larger
// than previous window, therefore the probability distribution between
// two windows are different.
//
// Now, we assume the recent half of window obey poission process.
// For a historical window of #requests [n_1, n_2, ..., n_p],
// We evaluates: [n_{p/2}, ..., n_{p}].
func (w *countProcess) Acceptable(k float64) (int64, bool) {
	// make a copy
	w.mu.RLock()
	events := w.nevents[len(w.nevents)/2:]
	w.mu.RUnlock()

	// moving average
	future := mean(events)
	if future > k {
		return int64(future), false
	}
	if 4*future <= k {
		return int64(future), true
	}

	// magnify avg arr accor. acceptance
	lambda := floats.Sum(events) / (float64(len(events)))
	acceptprob := (distuv.Poisson{Lambda: lambda}).CDF(k)

	return int64(future), acceptprob > 1-w.confidence
}

// ztest implements a one-tailed (right) z test (assumption: Guassian).
//
// H0: w.events[len(w.events)/2:] <= w.events[0:len(w.events)/2]
// H1: w.events[len(w.events)/2:] >  w.events[0:len(w.events)/2]
// we should have w.confidence to reject H0
func (w *countProcess) ztest() bool {
	// copy
	w.mu.RLock()
	events := make([]float64, len(w.nevents))
	copy(events, w.nevents)
	w.mu.RUnlock()

	if len(events) < 2 {
		return false
	}

	x1 := events[len(events)/2:]    // sample
	x2 := events[0 : len(events)/2] // population
	return (mean(x1)-mean(x2))/stddiv(x2) > -distuv.UnitNormal.Quantile(w.confidence)
}

func mean(x []float64) float64 {
	return floats.Sum(x) / float64(len(x))
}

func stddiv(x []float64) float64 {
	m := mean(x)
	var sum float64
	for _, i := range x {
		sum += math.Pow(float64(i)-m, 2)
	}
	return math.Sqrt(sum / float64(len(x)))
}
