// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gen_test

import (
	"math/rand"
	"testing"

	"changkun.de/x/pkg/gen"
)

func TestRandomString(t *testing.T) {
	s := gen.RandomString(12)
	if len(s) != 12 {
		t.Fatalf("want 12 chars, got: %v", s)
	}
}

func TestFastrand(t *testing.T) {
	if gen.Fastrand(1, 2) != 132101 {
		t.Fatalf("fastrand: wrong impl")
	}
}

func BenchmarkFastrand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gen.Fastrand(0, 0)
	}
}

func BenchmarkRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Uint32()
	}
}
