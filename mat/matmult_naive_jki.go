// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

import "sync"

// DotNaiveJKI matrix multiplication O(n^3)
func (A *Dense) DotNaiveJKI(B, C Matrix) (err error) {
	var (
		i, j, k int
		r       float64
	)
	if (A.Col() != B.Row()) || (C.Row() != A.Row()) || (C.Col() != B.Col()) {
		return ErrMatSize
	}

	for j = 0; j < B.Col(); j++ {
		for k = 0; k < A.Col(); k++ {
			r = B.At(k, j)
			for i = 0; i < A.Row(); i++ {
				C.Inc(i, j, r*A.At(i, k))
			}
		}
	}
	return
}

// DotNaiveJKIP matrix multiplication O(n^3)
func (A *Dense) DotNaiveJKIP(B, C Matrix) (err error) {
	if A.Col() != B.Row() || C.Row() != A.Row() || C.Col() != B.Col() {
		return ErrMatSize
	}

	wg := sync.WaitGroup{}
	for j := 0; j < B.Col(); j++ {
		wg.Add(1)
		go func(j int) {
			for k := 0; k < A.Col(); k++ {
				r := B.At(k, j)
				for i := 0; i < A.Row(); i++ {
					C.Inc(i, j, r*A.At(i, k))
				}
			}
			wg.Done()
		}(j)
	}
	wg.Wait()
	return
}
