// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package promise

import (
	"fmt"
	"sync"
	"time"

	"github.com/changkun/gobase/rt"
)

// All a list of function to be executed
func All(ff ...func()) (err error) {
	var wg sync.WaitGroup
	wg.Add(len(ff))
	for _, f := range ff {
		go func(g func()) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("service panic: %s, errors: %v", rt.GetFuncName(g), r)
				}
			}()
			g()
		}(f)
	}
	wg.Wait()
	return
}

// Run check and run ff return error if one return error
func Run(ff ...func() error) (err error) {
	for _, f := range ff {
		err = f()
	}
	return
}

// WaitGroup wrapper
type WaitGroup struct {
	waitgroup sync.WaitGroup
}

// NewWaitGroup creates a new WaitGroup of sync.WaitGroup
func NewWaitGroup() *WaitGroup {
	return &WaitGroup{
		// waitgroup: new(sync.WaitGroup),
	}
}

// Add given functions
func (wg *WaitGroup) Add(fs ...func()) {
	wg.waitgroup.Add(len(fs))
	for _, f := range fs {
		go func(f func()) {
			defer wg.waitgroup.Done()
			f()
		}(f)
	}
}

// Wait until group is done
func (wg *WaitGroup) Wait() {
	wg.waitgroup.Wait()
}

// WaitUntilTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out. Note that function added to waitgroup
// is still running without cancellation.
func (wg *WaitGroup) WaitUntilTimeout(timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		wg.waitgroup.Wait()
		close(c)
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timeout
	}
}
