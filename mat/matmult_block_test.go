// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

import (
	"fmt"
	"math"
	"testing"
)

func TestMatrix_MultBlock(t *testing.T) {
	A, err := NewDense(5, 8)(
		1, 2, 3, 4, 5, 6, 7, 8,
		9, 8, 7, 6, 5, 4, 3, 2,
		1, 2, 3, 4, 5, 6, 7, 8,
		9, 8, 7, 6, 5, 4, 3, 2,
		1, 2, 3, 4, 5, 6, 7, 8,
	)
	if err != nil {
		t.Errorf("New(5, 8) error, expect nil")
	}
	B, err := NewDense(8, 3)(
		9, 8, 7,
		6, 5, 4,
		3, 2, 1,
		2, 3, 4,
		5, 6, 7,
		8, 9, 8,
		7, 6, 5,
		4, 3, 2,
	)
	if err != nil {
		t.Errorf("New(8, 3) error, expect nil")
	}
	T, err := NewDense(5, 3)(
		192, 186, 168,
		248, 234, 212,
		192, 186, 168,
		248, 234, 212,
		192, 186, 168,
	)
	if err != nil {
		t.Errorf("New(5, 3) error, expect nil")
	}

	t.Run("DotBlock()", func(t *testing.T) {
		C := Zero(5, 3)
		err := A.DotBlock(B, C)
		if err != nil {
			t.Errorf("DotBlock(B, C) error, expect nil")
		}
		if !T.Equal(C) {
			t.Errorf("DotBlock() not euqal, expect euqal, got:")
			C.Print()
		}
	})
	t.Run("DotBlockP()", func(t *testing.T) {
		C := Zero(5, 3)
		err := A.DotBlockP(B, C)
		if err != nil {
			t.Errorf("DotBlockP(B, C) error, expect nil")
		}
		if !T.Equal(C) {
			t.Errorf("DotBlockP() not euqal, expect euqal, got:")
			C.Print()
		}
	})

	fs := map[string]func(int, Matrix, Matrix) error{
		"DotBlockIJK":  A.DotBlockIJK,
		"DotBlockIKJ":  A.DotBlockIKJ,
		"DotBlockJIK":  A.DotBlockJIK,
		"DotBlockJKI":  A.DotBlockJKI,
		"DotBlockKIJ":  A.DotBlockKIJ,
		"DotBlockKJI":  A.DotBlockKJI,
		"DotBlockIJKP": A.DotBlockIJKP,
		"DotBlockIKJP": A.DotBlockIKJP,
		"DotBlockJIKP": A.DotBlockJIKP,
		"DotBlockJKIP": A.DotBlockJKIP,
		"DotBlockKIJP": A.DotBlockKIJP,
		"DotBlockKJIP": A.DotBlockKJIP,
	}

	for name, f := range fs {
		t.Run(name, func(t *testing.T) {
			C := Zero(5, 3)
			err := f(2, B, C)
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

func BenchmarkDotBlock(b *testing.B) {
	for n := 80; n < 100; n++ {
		A := Rand(n, n)
		B := Rand(n, n)
		C := Zero(n, n)

		// L2 cache line: 256K
		// In core i7 the line sizes in L1 , L2 and L3 are the same: that is 64 Bytes.
		// see: https://stackoverflow.com/questions/14707803/line-size-of-l1-and-l2-caches
		// blockSize = sqrt(#CacheLines/3)
		nCacheLines := 256 * 1024 / 64
		blockSize := int(math.Sqrt(float64(nCacheLines) / 3))

		fs := map[string]func(int, Matrix, Matrix) error{
			fmt.Sprintf("DotBlockIJK()-block-size-%d-%dx%d", blockSize, n, n):  A.DotBlockIJK,
			fmt.Sprintf("DotBlockIKJ()-block-size-%d-%dx%d", blockSize, n, n):  A.DotBlockIKJ,
			fmt.Sprintf("DotBlockJIK()-block-size-%d-%dx%d", blockSize, n, n):  A.DotBlockJIK,
			fmt.Sprintf("DotBlockJKI()-block-size-%d-%dx%d", blockSize, n, n):  A.DotBlockJKI,
			fmt.Sprintf("DotBlockKIJ()-block-size-%d-%dx%d", blockSize, n, n):  A.DotBlockKIJ,
			fmt.Sprintf("DotBlockKJI()-block-size-%d-%dx%d", blockSize, n, n):  A.DotBlockKJI,
			fmt.Sprintf("DotBlockIJKP()-block-size-%d-%dx%d", blockSize, n, n): A.DotBlockIJKP,
			fmt.Sprintf("DotBlockIKJP()-block-size-%d-%dx%d", blockSize, n, n): A.DotBlockIKJP,
			fmt.Sprintf("DotBlockJIKP()-block-size-%d-%dx%d", blockSize, n, n): A.DotBlockJIKP,
			fmt.Sprintf("DotBlockJKIP()-block-size-%d-%dx%d", blockSize, n, n): A.DotBlockJKIP,
			fmt.Sprintf("DotBlockKIJP()-block-size-%d-%dx%d", blockSize, n, n): A.DotBlockKIJP,
			fmt.Sprintf("DotBlockKJIP()-block-size-%d-%dx%d", blockSize, n, n): A.DotBlockKJIP,
		}

		for name, f := range fs {
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					f(blockSize, B, C)
				}
			})
		}
	}
}
