// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sysmon

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/changkun/gobase/rp"
)

type sysmon struct {
	p                       rp.RandomProcess
	observeFunc, syscapFunc func() int64
	actionFunc              func(int64) interface{}
	interval                time.Duration
	stopped                 int32
}

var sysmon0 sysmon

// Init ...
func Init(
	windowSize int, confidence float64,
	observeFunc, syscapFunc func() int64,
	actionFunc func(int64) interface{},
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

// Run runs an initialized system momitor
func Run() {
	go func() {
		for {
			if atomic.CompareAndSwapInt32(&sysmon0.stopped, 1, 2) {
				break
			}

			time.Sleep(sysmon0.interval)

			ob := sysmon0.observeFunc()
			cap := sysmon0.syscapFunc()
			fmt.Printf("sysmon: check, ob: %v, cap: %v\n", ob, cap)

			// store number of events during sleep
			if ob > 0 {
				sysmon0.p.Store(float64(ob))
			}

			suggestion, ok := sysmon0.p.Acceptable(float64(sysmon0.syscapFunc()))
			if cap > suggestion {
				if !sysmon0.p.Significant() || ok {
					continue
				}
			}

			sysmon0.actionFunc(suggestion)
		}
	}()
}

// Stop stops system monitoring
func Stop(wait bool) {
	atomic.CompareAndSwapInt32(&sysmon0.stopped, 0, 1)
	if !wait {
		return
	}
	for {
		if atomic.LoadInt32(&sysmon0.stopped) == 2 {
			break
		}
		runtime.Gosched()
	}
}
