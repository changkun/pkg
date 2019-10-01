// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

// DotNaive matrix multiplication O(n^3)
// Use JIK version here, see ./benchmark/README.md
func (A *Dense) DotNaive(B, C Matrix) (err error) {
	return A.DotNaiveJIK(B, C)
}

// DotNaiveP matrix multiplication O(n^3)
// Use JIKP version here, see ./benchmark/README.md
func (A *Dense) DotNaiveP(B, C Matrix) (err error) {
	return A.DotNaiveJIKP(B, C)
}
