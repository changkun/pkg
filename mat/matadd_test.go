// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

import "testing"

func TestAdd(t *testing.T) {
	A := RandP(100, 100)
	B := Zero(100, 100)
	Bb := Zero(100, 1)
	_, err := Add(A, Bb)
	if err == nil {
		t.Errorf("Add() not error, expect err")
	}

	C, err := Add(A, B)
	if err != nil {
		t.Errorf("Add() error, expect not err")
	}
	if !C.Equal(A) {
		t.Errorf("Add() error, expect euqal")
	}
}

func TestMatrix_Add(t *testing.T) {
	A := RandP(100, 100)
	oldA := *A
	B := Zero(100, 100)
	Bb := Zero(100, 1)
	err := A.Add(Bb)
	if err == nil {
		t.Errorf("Add() not error, expect err")
	}

	err = A.Add(B)
	if err != nil {
		t.Errorf("Add() error, expect not err")
	}
	if !oldA.Equal(A) {
		t.Errorf("Add() error, expect euqal")
	}
}
