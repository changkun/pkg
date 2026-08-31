// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gls_test

import (
	"sync"
	"testing"

	"changkun.de/x/pkg/gls"
)

func TestGLS(t *testing.T) {
	gls.Store("hello", "world")

	// The point of goroutine local storage is that the value stays on the
	// goroutine that stored it. Report from the test goroutine: t.Fatalf calls
	// runtime.Goexit, which does not end the test when called anywhere else.
	var leaked bool
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, leaked = gls.Get("hello")
	}()
	wg.Wait()
	if leaked {
		t.Fatal("stored value is visible from another goroutine")
	}

	v, ok := gls.Get("hello")
	if !ok {
		t.Fatal("stored value is not visible from the goroutine that stored it")
	}
	if v != "world" {
		t.Errorf("stored value = %v, want world", v)
	}

	gls.Clear()
	if _, ok := gls.Get("hello"); ok {
		t.Error("value is still visible after Clear")
	}
}
