// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

// DotBlock matrix multiplication
// Use JIK block 36 version here, see ./benchmark/README.md
func (A *Dense) DotBlock(B, C Matrix) (err error) {
	return A.DotBlockJIK(36, B, C)
}

// DotBlockP matrix multiplication
// Use JIKP block 36 version here, see ./benchmark/README.md
func (A *Dense) DotBlockP(B, C Matrix) (err error) {
	return A.DotBlockJIKP(36, B, C)
}
