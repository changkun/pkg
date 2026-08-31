// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package errhandling contrasts two ways of writing the same file copy: one
// that checks the error of every step at the call site, and one that hides the
// checks behind a type which keeps the first error and turns every later step
// into a no-op.
//
// Inspired by https://www.youtube.com/watch?v=1B71SL6Y0kA
package errhandling

import (
	"fmt"
	"io"
	"os"
)

// CopyFile copies src to dst and checks each step at the call site. A failed
// copy leaves no partial dst behind.
func CopyFile(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("copy %s %s: %w", src, dst, err)
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("copy %s %s: %w", src, dst, err)
	}

	if _, err := io.Copy(w, r); err != nil {
		w.Close()
		_ = os.Remove(dst)
		return fmt.Errorf("copy %s %s: %w", src, dst, err)
	}

	if err := w.Close(); err != nil {
		_ = os.Remove(dst)
		return fmt.Errorf("copy %s %s: %w", src, dst, err)
	}
	return nil
}

// SafeCopyFile copies src to dst through SafeCopy. The steps read as a
// straight line because SafeCopy carries the error, and only the end of the
// function checks it.
func SafeCopyFile(src, dst string) error {
	c := safeOpen(src)
	c.Create(dst)
	c.Copy()
	c.Close()

	if c.err != nil {
		_ = os.Remove(dst)
		return fmt.Errorf("copy %s %s: %w", src, dst, c.err)
	}
	return nil
}

// SafeCopy holds a copy in progress together with the first error it hit.
// Every step after that error does nothing, so the caller checks once instead
// of after each call.
type SafeCopy struct {
	r, w     *os.File
	src, dst string
	err      error
}

// safeOpen starts a copy by opening src for reading.
func safeOpen(src string) SafeCopy {
	r, err := os.Open(src)
	return SafeCopy{r: r, src: src, err: err}
}

// Create opens dst for writing.
func (c *SafeCopy) Create(dst string) {
	c.dst = dst
	if c.err != nil {
		return
	}
	c.w, c.err = os.Create(dst)
}

// Copy streams the source into the destination.
func (c *SafeCopy) Copy() {
	if c.err != nil {
		return
	}
	_, c.err = io.Copy(c.w, c.r)
}

// Close releases both files. It always runs, so an aborted copy does not leak
// a descriptor, and it keeps the earliest error rather than the last one.
func (c *SafeCopy) Close() {
	if c.r != nil {
		if err := c.r.Close(); err != nil && c.err == nil {
			c.err = err
		}
	}
	if c.w != nil {
		if err := c.w.Close(); err != nil && c.err == nil {
			c.err = err
		}
	}
}
