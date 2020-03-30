// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/changkun/gobase/mkill"
)

func main() {
	mkill.GOMAXTHREADS(10)
	mkill.SetDebug(true)
	for {
		time.Sleep(time.Second)
		go func() {
			time.Sleep(time.Second * 10)
		}()
	}
}
