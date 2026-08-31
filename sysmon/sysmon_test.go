// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sysmon_test

import (
	"math/rand/v2"
	"testing"
	"time"

	"changkun.de/x/pkg/sysmon"
)

const (
	// interval is short enough to keep the test under a second and long
	// enough that the monitor's own bookkeeping is not the dominant cost.
	interval = 5 * time.Millisecond
	// rounds is the number of intervals to simulate.
	rounds = 200
	// windowSize is the observation window. The package documents that the
	// z test wants more than 15 samples to justify its normality assumption.
	windowSize = 16
	// confidence is the level for both the significance and acceptance tests.
	confidence = 0.05
	// peakDemand is the largest number of events one interval can bring.
	peakDemand = 20
	// startCapacity is what the system can serve before any scaling.
	startCapacity = 30
	// warmup is the number of intervals to ignore when judging the outcome.
	// The monitor cannot predict a distribution it has not observed yet.
	warmup = 2 * windowSize
)

// system is a simulated resource pool driven by the monitor. Every field is
// touched only from the monitor goroutine, and Stop(true) waits for that
// goroutine to exit before the test reads them.
type system struct {
	rng      *rand.Rand
	capacity int64

	rounds   int
	starved  int   // intervals where demand exceeded capacity
	scaled   int   // intervals where the monitor scaled up
	minSpare int64 // lowest capacity observed
}

// observe consumes one interval of demand. Consumption is capped by what is
// actually available: a pool cannot serve requests it has no capacity for, so
// capacity floors at zero rather than going negative.
func (s *system) observe() int64 {
	s.rounds++
	demand := s.rng.Int64N(peakDemand)
	served := min(demand, s.capacity)
	if demand > s.capacity {
		s.starved++
	}
	s.capacity -= served
	if s.capacity < s.minSpare {
		s.minSpare = s.capacity
	}
	return served
}

// syscap reports the capacity left for the next interval.
func (s *system) syscap() int64 { return s.capacity }

// scale applies the monitor's suggestion.
func (s *system) scale(suggestion int64) any {
	s.scaled++
	s.capacity += 4 * suggestion
	return nil
}

// TestSysmonKeepsCapacityAhead drives the monitor against a stationary demand
// and checks its contract: capacity stays ahead of the load often enough that
// the pool rarely runs dry.
//
// The random source is explicit. The global source in math/rand has been
// seeded randomly since go 1.20, which would make this a different test on
// every run.
func TestSysmonKeepsCapacityAhead(t *testing.T) {
	s := &system{
		rng:      rand.New(rand.NewPCG(1, 2)),
		capacity: startCapacity,
		minSpare: startCapacity,
	}

	sysmon.Init(
		windowSize, confidence,
		s.observe, s.syscap, s.scale,
		interval,
	)
	sysmon.Run()
	time.Sleep(rounds * interval)
	sysmon.Stop(true)

	if s.rounds < warmup {
		t.Fatalf("only %d intervals ran, need more than the %d warmup intervals to judge", s.rounds, warmup)
	}
	if s.minSpare < 0 {
		t.Errorf("capacity went to %d; the simulation must never serve more than it has", s.minSpare)
	}
	if s.scaled == 0 {
		t.Error("monitor never scaled up, so the test proved nothing about it")
	}

	// Demand averages peakDemand/2 per interval against a starting capacity of
	// startCapacity, so an unscaled pool starves within a handful of
	// intervals. A monitor that tracks the load keeps that rate low.
	const maxStarveRatio = 0.25
	if ratio := float64(s.starved) / float64(s.rounds); ratio > maxStarveRatio {
		t.Errorf("starved in %d of %d intervals (%.0f%%), want at most %.0f%%",
			s.starved, s.rounds, ratio*100, maxStarveRatio*100)
	}
	t.Logf("intervals=%d scaled=%d starved=%d minSpare=%d", s.rounds, s.scaled, s.starved, s.minSpare)
}
