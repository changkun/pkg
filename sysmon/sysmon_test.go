// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sysmon_test

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/changkun/gobase/sysmon"
)

func TestSysmon(t *testing.T) {
	total := int64(30)
	factor := int64(20)

	sysmon.Init(
		15, 0.05,
		func() int64 {
			ob := rand.Int63() % factor
			println("consume: ", ob)
			atomic.AddInt64(&total, -ob)
			return ob
		}, func() int64 {
			cap := atomic.LoadInt64(&total)
			if cap < 0 {
				panic(fmt.Sprintf("sysmon fail to scale: %v", cap))
			}
			return cap
		}, func(suggestion int64) interface{} {
			// double
			atomic.AddInt64(&total, 4*suggestion)
			println("scale! ", 4*suggestion)
			return nil
		}, time.Second,
	)

	go sysmon.Run()

	time.Sleep(20 * time.Second)

	sysmon.Stop(true)
}
