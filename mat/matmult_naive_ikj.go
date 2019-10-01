// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

import "sync"

// DotNaiveIKJ matrix multiplication O(n^3)
func (A *Dense) DotNaiveIKJ(B, C Matrix) (err error) {
	var (
		i, j, k int
		r       float64
	)
	if (A.Col() != B.Row()) || (C.Row() != A.Row()) || (C.Col() != B.Col()) {
		return ErrMatSize
	}

	for i = 0; i < A.Row(); i++ {
		for k = 0; k < A.Col(); k++ {
			r = A.At(i, k)
			for j = 0; j < B.Col(); j++ {
				C.Inc(i, j, r*B.At(k, j))
			}
		}
	}
	return
}

// DotNaiveIKJP matrix multiplication O(n^3)
func (A *Dense) DotNaiveIKJP(B, C Matrix) (err error) {
	if A.Col() != B.Row() || C.Row() != A.Row() || C.Col() != B.Col() {
		return ErrMatSize
	}

	wg := sync.WaitGroup{}
	for i := 0; i < A.Row(); i++ {
		wg.Add(1)
		go func(i int) {
			for k := 0; k < A.Col(); k++ {
				r := A.At(i, k)
				for j := 0; j < B.Col(); j++ {
					C.Inc(i, j, r*B.At(k, j))
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return
}
