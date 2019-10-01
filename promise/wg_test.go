// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package promise_test

import (
	"sync"
	"testing"
	"time"

	"github.com/changkun/gobase/promise"
)

func TestWaitGroupOnTimeout(t *testing.T) {
	var mutex sync.Mutex
	wg := promise.NewWaitGroup()
	a, b, c := 0, 0, 0

	sleepOneMillisecond := func() {
		time.Sleep(time.Millisecond)
		mutex.Lock()
		a = 1
		mutex.Unlock()
	}
	sleepTwoMillisecond := func() {
		time.Sleep(time.Millisecond * 2)
		mutex.Lock()
		b = 2
		mutex.Unlock()
	}
	sleepThreeMillisecond := func() {
		time.Sleep(time.Millisecond * 3)
		mutex.Lock()
		c = 3
		mutex.Unlock()
	}
	wg.Add(
		sleepOneMillisecond,
		sleepTwoMillisecond,
		sleepThreeMillisecond,
	)
	wg.WaitUntilTimeout(time.Nanosecond)

	mutex.Lock()
	if a != 0 {
		t.Fatalf("want 0, got: %v", a)
	}
	if b != 0 {
		t.Fatalf("want 0, got: %v", b)
	}
	if c != 0 {
		t.Fatalf("want 0, got: %v", c)
	}
	mutex.Unlock()

	time.Sleep(time.Second)
	mutex.Lock()
	if a != 1 {
		t.Fatalf("want 0, got: %v", a)
	}
	if b != 2 {
		t.Fatalf("want 0, got: %v", b)
	}
	if c != 3 {
		t.Fatalf("want 0, got: %v", c)
	}
	mutex.Unlock()
}
