// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

import (
	"fmt"
	"testing"
)

var (
	m = 5
	n = 10
)

func TestZero(t *testing.T) {
	A := Zero(m, n)
	for i := 0; i < A.Row(); i++ {
		for j := 0; j < A.Col(); j++ {
			if A.At(i, j) != 0.0 {
				t.Errorf("Zero() expect zero element, got: %.2f", A.At(i, j))
			}
		}
	}
}

func TestRand(t *testing.T) {
	t.Run("Rand()", func(t *testing.T) {
		if mat := Rand(m, n); mat == nil {
			t.Errorf("Rand() expect non nil matrix, got: nil")
		}
	})
	t.Run("RandP()", func(t *testing.T) {
		if mat := RandP(m, n); mat == nil {
			t.Errorf("RandP() expect non nil matrix, got: nil")
		}
	})
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		m, n int
		arr  []float64
	}{
		{
			name: "3x3 matrix",
			m:    3, n: 3,
			arr: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "2x4 matrix",
			m:    2, n: 4,
			arr: []float64{1, 2, 3, 4, 5, 6, 7, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDense(tt.m, tt.n)(tt.arr...)
			if err != nil {
				t.Errorf("New() error = %v, want nil", err)
			}
			for i := 0; i < tt.m; i++ {
				for j := 0; j < tt.n; j++ {
					if got.At(i, j) != tt.arr[i*tt.n+j] {
						t.Errorf("New(i, j) = %.2f, want %.2f", tt.arr[i*tt.n+j], got.At(i, j))
					}
				}
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDenseP(tt.m, tt.n)(tt.arr...)
			if err != nil {
				t.Errorf("NewP() error = %v, want nil", err)
			}
			for i := 0; i < tt.m; i++ {
				for j := 0; j < tt.n; j++ {
					if got.At(i, j) != tt.arr[i*tt.n+j] {
						t.Errorf("NewP(i, j) = %.2f, want %.2f", tt.arr[i*tt.n+j], got.At(i, j))
					}
				}
			}
		})
	}
}

func TestNewFail(t *testing.T) {
	tests := []struct {
		name string
		m, n int
		arr  []float64
	}{
		{
			name: "3x3 matrix",
			m:    3, n: 3,
			arr: []float64{1},
		},
		{
			name: "2x4 matrix",
			m:    2, n: 4,
			arr: []float64{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDense(tt.m, tt.n)(tt.arr...)
			if err != ErrNumElem {
				t.Errorf("New() error is nil, want err")
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDenseP(tt.m, tt.n)(tt.arr...)
			if err != ErrNumElem {
				t.Errorf("NewP() error is nil, want err")
			}
		})
	}
}

func TestMatrixAccess(t *testing.T) {
	A := Zero(m, n)
	A.Print()
	t.Run("Size()", func(t *testing.T) {
		mm, nn := A.Size()
		if mm != m || nn != n {
			t.Errorf("Size() error, expect: %d,%d, got: %d,%d", m, n, mm, nn)
		}
	})
	t.Run("Row()", func(t *testing.T) {
		mm := A.Row()
		if mm != m {
			t.Errorf("Row() error, expect: %d, got: %d", m, mm)
		}
	})
	t.Run("Col()", func(t *testing.T) {
		nn := A.Col()
		if nn != n {
			t.Errorf("Col() error, expect: %d, got: %d", n, nn)
		}
	})
	t.Run("At()", func(t *testing.T) {
		val := A.At(m-1, n-1)
		if val != 0.0 {
			t.Errorf("At() error, expect: 0.00, got: %.2f", val)
		}
	})
	t.Run("Set()", func(t *testing.T) {
		A.Set(m-1, n-1, 1.0)
		if A.At(m-1, n-1) != 1.0 {
			t.Errorf("Set() error, expect: 1.00, got: %.2f", A.At(m-1, n-1))
		}
	})
	t.Run("Inc()", func(t *testing.T) {
		A.Inc(m-1, n-1, 1.0)
		if A.At(m-1, n-1) != 2.0 {
			t.Errorf("Inc() error, expect: 2.00, got: %.2f", A.At(m-1, n-1))
		}
	})
	t.Run("Mult()", func(t *testing.T) {
		A.Mult(m-1, n-1, 3.0)
		if A.At(m-1, n-1) != 6.0 {
			t.Errorf("Mult() error, expect: 6.00, got: %.2f", A.At(m-1, n-1))
		}
	})
	t.Run("Pow()", func(t *testing.T) {
		A.Pow(m-1, n-1, 2.0)
		if A.At(m-1, n-1) != 36.0 {
			t.Errorf("Pow() error, expect: 36.00, got: %.2f", A.At(m-1, n-1))
		}
	})
	t.Run("Equal()", func(t *testing.T) {
		B := Zero(m, 1)
		if A.Equal(B) {
			t.Errorf("Equal() equal, expect: inequal, got: equal")
		}
		B = Zero(m, n)
		if A.Equal(B) {
			t.Errorf("Equal() equal, expect: inequal, got: equal")
		}
		B.Set(m-1, n-1, 36.0)
		if !A.Equal(B) {
			t.Errorf("Equal() not equal, expect: equal, got: inequal")
		}
	})
}

// -------------- benchmarks -----------------
func BenchmarkNewRand(b *testing.B) {
	fs := map[string]func(m, n int) *Dense{
		"Rand":         Rand,
		"RandParallel": RandP,
	}
	sizes := []int{
		5, 50, 500,
	}
	for k, f := range fs {
		b.Run(k, func(b *testing.B) {
			for _, s := range sizes {
				b.Run(fmt.Sprintf("size-%d", s), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						f(s, s)
					}
				})
			}
		})
	}
}
