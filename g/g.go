// Package g is a goroutine life cycle management package.
package g

import (
	"fmt"
	"sync/atomic"
)

// Go runs a given f concurrently and returns
// a coresponding goroutine instance.
func Go(f func() error) *G {
	g := &G{
		f:      f,
		cancel: make(chan struct{}),
		wait:   make(chan struct{}),
	}
	go g.worker()
	return g
}

// G is a goroutine
type G struct {
	f   func() error
	err error

	undo   atomic.Value
	cancel chan struct{}
	wait   chan struct{}
	done   uint32
}

// WithUndo sets up a function that can be executed
// if a goroutine is canceled.
func (g *G) WithUndo(undo func()) *G {
	g.undo.Store(undo)
	return g
}

// Cancel cancels the goroutine. Note that the goroutine
// does not get stopped and a good practice is to use
// Cancel with WithUndo function:
//
//   g := g.Go(func() error { ... }).WithUndo(func() { ... })
//   ...
//   err := g.Cancel().Wait().Err()
//   if err != nil { ... }
//
// Then the settled undo function will be executed after
// the cancellation.
func (g *G) Cancel() *G {
	g.cancel <- struct{}{}
	return g
}

// Wait wait until the goroutine terminates.
func (g *G) Wait() *G {
	if atomic.LoadUint32(&g.done) != 1 {
		<-g.wait
	}
	atomic.StoreUint32(&g.done, 1)
	return g
}

// Err reports if the goroutine returns error,
// This call will wait until the goroutine terminates.
func (g *G) Err() error {
	g.Wait()
	return g.err
}

func (g *G) worker() {
	done := make(chan error, 1)
	defer func() {
		g.wait <- struct{}{}
	}()
	go func() {
		var err error
		defer func() {
			if r := recover(); r != nil {
				done <- fmt.Errorf("%v", r)
				return
			}
			done <- err
		}()
		err = g.f()
	}()
	select {
	case err := <-done:
		g.err = err
	case <-g.cancel:
		close(g.cancel)
		g.err = <-done // wait until execution is done
		undo, ok := g.undo.Load().(func())
		if ok {
			undo()
		}
	}
}
