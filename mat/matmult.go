// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

// Dot matrix multiplication
func (A *Dense) Dot(B, C Matrix) (err error) {
	return A.DotNaive(B, C)
}

// DotP matrix multiplication
func (A *Dense) DotP(B, C Matrix) (err error) {
	return A.DotNaiveP(B, C)
}

// Dot matrix multiplication
func Dot(A, B *Dense) (*Dense, error) {
	C := Zero(A.Row(), B.Col())
	if err := A.Dot(B, C); err != nil {
		return nil, err
	}
	return C, nil
}

// DotP matrix multiplication
func DotP(A, B *Dense) (*Dense, error) {
	C := Zero(A.Row(), B.Col())
	if err := A.DotP(B, C); err != nil {
		return nil, err
	}
	return C, nil
}
