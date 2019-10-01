// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

import "testing"

func TestMatrix_Dot(t *testing.T) {
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

	t.Run("Matrix_Dot()", func(t *testing.T) {
		C := Zero(5, 3)
		err := A.Dot(B, C)
		if err != nil {
			t.Errorf("Matrix_Dot(B, C) error, expect nil")
		}
		if !T.Equal(C) {
			t.Errorf("Matrix_Dot() not euqal, expect euqal, got:")
			C.Print()
		}
	})
	t.Run("Matrix_DotP()", func(t *testing.T) {
		C := Zero(5, 3)
		err := A.DotP(B, C)
		if err != nil {
			t.Errorf("Matrix_DotP(B, C) error, expect nil")
		}
		if !T.Equal(C) {
			t.Errorf("Matrix_DotP() not euqal, expect euqal, got:")
			C.Print()
		}
	})
	t.Run("Dot()", func(t *testing.T) {
		C, err := Dot(A, B)
		if err != nil {
			t.Errorf("Dot(B, C) error, expect nil")
		}
		if !T.Equal(C) {
			t.Errorf("Dot() not euqal, expect euqal, got:")
			C.Print()
		}
	})
	t.Run("DotP()", func(t *testing.T) {
		C, err := DotP(A, B)
		if err != nil {
			t.Errorf("DotP(B, C) error, expect nil")
		}
		if !T.Equal(C) {
			t.Errorf("DotP() not euqal, expect euqal, got:")
			C.Print()
		}
	})
}
