// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

import "sync"

// DotNaiveKIJ matrix multiplication O(n^3)
func (A *Dense) DotNaiveKIJ(B, C Matrix) (err error) {
	var (
		i, j, k int
		r       float64
	)
	if (A.Col() != B.Row()) || (C.Row() != A.Row()) || (C.Col() != B.Col()) {
		return ErrMatSize
	}

	for k = 0; k < A.Col(); k++ {
		for i = 0; i < A.Row(); i++ {
			r = A.At(i, k)
			for j = 0; j < B.Col(); j++ {
				C.Inc(i, j, r*B.At(k, j))
			}
		}
	}
	return
}

// DotNaiveKIJP matrix multiplication O(n^3)
func (A *Dense) DotNaiveKIJP(B, C Matrix) (err error) {
	if A.Col() != B.Row() || C.Row() != A.Row() || C.Col() != B.Col() {
		return ErrMatSize
	}

	wg := sync.WaitGroup{}
	for k := 0; k < A.Col(); k++ {
		wg.Add(1)
		go func(k int) {
			for i := 0; i < A.Row(); i++ {
				r := A.At(i, k)
				for j := 0; j < B.Col(); j++ {
					C.Inc(i, j, r*B.At(k, j))
				}
			}
			wg.Done()
		}(k)
	}
	wg.Wait()
	return
}
