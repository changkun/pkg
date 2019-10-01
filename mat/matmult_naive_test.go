// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

import (
	"fmt"
	"testing"
)

func TestMatrix_DotNaive(t *testing.T) {
	A, err := NewDense(2, 3)(
		1, 2, 3,
		2, 3, 1,
	)
	if err != nil {
		t.Errorf("New(2, 3) error, expect nil")
	}
	B, err := NewDense(3, 1)(
		3,
		2,
		1,
	)
	if err != nil {
		t.Errorf("New(3, 1) error, expect nil")
	}
	T, err := NewDense(2, 1)(
		10,
		13,
	)
	if err != nil {
		t.Errorf("New(2, 1) error, expect nil")
	}

	fs := map[string]func(Matrix, Matrix) error{
		"DotNaive":     A.DotNaive,
		"DotNaiveP":    A.DotNaiveP,
		"DotNaiveIJK":  A.DotNaiveIJK,
		"DotNaiveIKJ":  A.DotNaiveIKJ,
		"DotNaiveJIK":  A.DotNaiveJIK,
		"DotNaiveJKI":  A.DotNaiveJKI,
		"DotNaiveKIJ":  A.DotNaiveKIJ,
		"DotNaiveKJI":  A.DotNaiveKJI,
		"DotNaiveIJKP": A.DotNaiveIJKP,
		"DotNaiveIKJP": A.DotNaiveIKJP,
		"DotNaiveJIKP": A.DotNaiveJIKP,
		"DotNaiveJKIP": A.DotNaiveJKIP,
		"DotNaiveKIJP": A.DotNaiveKIJP,
		"DotNaiveKJIP": A.DotNaiveKJIP,
	}

	for name, f := range fs {
		t.Run(name, func(t *testing.T) {
			C := Zero(2, 1)
			err := f(B, C)
			if err != nil {
				t.Errorf("DotBlockIJK(B, C) error, expect nil")
			}
			if !T.Equal(C) {
				t.Errorf("DotBlockIJK() not euqal, expect euqal, got:")
				C.Print()
			}
		})
	}
}

// ----------------------- benchmarks ----------------------------

func BenchmarkDotNaive(b *testing.B) {
	for n := 80; n < 100; n++ {
		A := Rand(n, n)
		B := Rand(n, n)
		C := Zero(n, n)

		fs := map[string]func(Matrix, Matrix) error{
			fmt.Sprintf("DotNaiveIJK()-%dx%d", n, n):  A.DotNaiveIJK,
			fmt.Sprintf("DotNaiveIKJ()-%dx%d", n, n):  A.DotNaiveIKJ,
			fmt.Sprintf("DotNaiveJIK()-%dx%d", n, n):  A.DotNaiveJIK,
			fmt.Sprintf("DotNaiveJKI()-%dx%d", n, n):  A.DotNaiveJKI,
			fmt.Sprintf("DotNaiveKIJ()-%dx%d", n, n):  A.DotNaiveKIJ,
			fmt.Sprintf("DotNaiveKJI()-%dx%d", n, n):  A.DotNaiveKJI,
			fmt.Sprintf("DotNaiveIJKP()-%dx%d", n, n): A.DotNaiveIJKP,
			fmt.Sprintf("DotNaiveIKJP()-%dx%d", n, n): A.DotNaiveIKJP,
			fmt.Sprintf("DotNaiveJIKP()-%dx%d", n, n): A.DotNaiveJIKP,
			fmt.Sprintf("DotNaiveJKIP()-%dx%d", n, n): A.DotNaiveJKIP,
			fmt.Sprintf("DotNaiveKIJP()-%dx%d", n, n): A.DotNaiveKIJP,
			fmt.Sprintf("DotNaiveKJIP()-%dx%d", n, n): A.DotNaiveKJIP,
		}

		for name, f := range fs {
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					f(B, C)
				}
			})
		}
	}
}
