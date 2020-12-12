// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package g_test

import (
	"errors"
	"testing"
	"time"

	"changkun.de/x/pkg/g"
)

func TestGo(t *testing.T) {
	g := g.Go(func() error {
		println("hello")
		return nil
	})
	g.Wait()
}

func TestGoCancel(t *testing.T) {
	g := g.Go(func() error {
		time.Sleep(time.Millisecond)
		println("hello")
		return nil
	})
	g.Cancel().Wait()
}

func TestGoErr(t *testing.T) {
	g := g.Go(func() error {
		return errors.New("throw")
	})
	if err := g.Err(); err == nil {
		t.Fatalf("cannot get returned error")
	}
	t.Logf("get returned error: %v", g.Err())
}
func TestGoPanicErr(t *testing.T) {
	g := g.Go(func() error {
		panic("throw")
	})
	if err := g.Err(); err == nil {
		t.Fatalf("cannot get returned error")
	}
	t.Logf("get returned error: %v", g.Err())
}

func TestGoUndo(t *testing.T) {
	g := g.Go(func() error {
		time.Sleep(time.Millisecond)
		println("wtf")
		return nil
	}).WithUndo(func() {
		println("undo")
	})
	err := g.Cancel().Wait().Err()
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
}
