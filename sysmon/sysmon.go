// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sysmon

import (
	"runtime"
	"sync/atomic"
	"time"

	"changkun.de/x/pkg/rp"
)

// Monitor lifecycle states.
const (
	running int32 = iota
	stopping
	stopped
)

type sysmon struct {
	p                       rp.RandomProcess
	observeFunc, syscapFunc func() int64
	actionFunc              func(int64) any
	interval                time.Duration
	state                   atomic.Int32
}

var sysmon0 sysmon

// Init prepares the monitor. observeFunc reports the events seen since the
// last check, syscapFunc reports the capacity available now, and actionFunc
// applies the scaling suggestion.
func Init(
	windowSize int, confidence float64,
	observeFunc, syscapFunc func() int64,
	actionFunc func(int64) any,
	interval time.Duration,
) {
	sysmon0 = sysmon{
		p:           rp.NewCountProcess(float64(windowSize), confidence),
		observeFunc: observeFunc,
		syscapFunc:  syscapFunc,
		actionFunc:  actionFunc,
		interval:    interval,
	}
}

// Run starts an initialized system monitor. It returns immediately; the
// monitor runs until Stop.
func Run() {
	go func() {
		for {
			if sysmon0.state.CompareAndSwap(stopping, stopped) {
				return
			}

			time.Sleep(sysmon0.interval)

			nevents := sysmon0.observeFunc()
			capacity := sysmon0.syscapFunc()

			// Store number of events during sleep.
			if nevents > 0 {
				sysmon0.p.Store(float64(nevents))
			}

			suggestion, ok := sysmon0.p.Acceptable(float64(capacity))
			if capacity > suggestion {
				if !sysmon0.p.Significant() || ok {
					continue
				}
			}

			sysmon0.actionFunc(suggestion)
		}
	}()
}

// Stop stops system monitoring. If wait is true it returns only once the
// monitor goroutine has exited, so the caller may then read whatever the
// callbacks were writing.
func Stop(wait bool) {
	sysmon0.state.CompareAndSwap(running, stopping)
	if !wait {
		return
	}
	for sysmon0.state.Load() != stopped {
		runtime.Gosched()
	}
}
