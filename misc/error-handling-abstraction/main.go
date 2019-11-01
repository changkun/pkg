// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// inspiring from https://www.youtube.com/watch?v=1B71SL6Y0kA
package main

import (
	"fmt"
	"io"
	"os"
)

// CopyFile : General Approach
func CopyFile(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("copy %s %s: %v", src, dst, err)
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("copy %s %s: %v", src, dst, err)
	}

	if _, err := io.Copy(w, r); err != nil {
		w.Close()
		os.Remove(dst)
		return fmt.Errorf("copy %s %s: %v", src, dst, err)
	}

	if err := w.Close(); err != nil {
		os.Remove(dst)
		return fmt.Errorf("copy %s %s: %v", src, dst, err)
	}

	return nil
}

// SafeCopyFile : Error handling by abstraction
func SafeCopyFile(src, dst string) error {
	c := safeOpen(src)
	c.Create(dst)
	c.Copy()
	c.Close()

	if c.err != nil {
		os.Remove(dst)
		return fmt.Errorf("copy %s %s: %v", src, dst, c.err)
	}
	return nil
}

type SafeCopy struct {
	r, w     *os.File
	src, dst string
	err      error
}

func safeOpen(src string) SafeCopy {
	r, err := os.Open(src)
	return SafeCopy{r: r, src: src, err: err}
}

func (c *SafeCopy) Create(dst string) {
	c.dst = dst
	if c.err != nil {
		c.err = fmt.Errorf("copy %s %s: %v", c.src, c.dst, c.err)
		return
	}
	c.w, c.err = os.Create(c.dst)
}

func (c *SafeCopy) Copy() {
	if c.err != nil {
		c.r.Close()
		c.err = fmt.Errorf("copy %s %s: %v", c.src, c.dst, c.err)
		return
	}
	_, c.err = io.Copy(c.r, c.w)
}

func (c *SafeCopy) Close() {
	if c.err != nil {
		if c.w != nil {
			c.w.Close()
			os.Remove(c.dst)
		}
		c.err = fmt.Errorf("copy %s %s: %v", c.src, c.dst, c.err)
		return
	}

	c.err = c.r.Close()
	c.err = c.w.Close()
}
