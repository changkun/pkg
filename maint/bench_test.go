// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a GPL-3.0
// license that can be found in the LICENSE file.

package maint_test

import (
	"fmt"
	"testing"

	"changkun.de/x/pkg/maint"
)

func BenchmarkCall(b *testing.B) {
	f1 := func() {}
	f2 := func() {}

	maint.Init(func() {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if i%2 == 0 {
				maint.Call(f1)
			} else {
				maint.Call(f2)
			}
		}
	})
}

func ExampleInit() {
	maint.Init(func() {
		maint.Call(func() {
			fmt.Println("from main thread")
		})
	})
	// Output: from main thread
}
