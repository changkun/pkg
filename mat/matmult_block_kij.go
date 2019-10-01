// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

import "sync"

// DotBlockKIJ block matrix multiplication
func (A *Dense) DotBlockKIJ(blockSize int, B, C Matrix) (err error) {
	if (A.Col() != B.Row()) || (C.Row() != A.Row()) || (C.Col() != B.Col()) {
		return ErrMatSize
	}
	min := A.Row()
	if A.Col() < min {
		min = A.Col()
	}
	if B.Col() < min {
		min = B.Col()
	}
	var (
		kk, jj, i, j, k int
		r               float64
		en              = blockSize * (min / blockSize)
	)

	for kk = 0; kk < en; kk += blockSize {
		for jj = 0; jj < en; jj += blockSize {
			for k = kk; k < kk+blockSize; k++ {
				for i = 0; i < A.Row(); i++ {
					r = A.At(i, k)
					for j = jj; j < jj+blockSize; j++ {
						C.Inc(i, j, r*B.At(k, j))
					}
				}
			}
		}
		for k = kk; k < kk+blockSize; k++ {
			for i = 0; i < A.Row(); i++ {
				r = A.At(i, k)
				for j = en; j < B.Col(); j++ {
					C.Inc(i, j, r*B.At(k, j))
				}
			}
		}
	}

	// residule bottom
	for jj = 0; jj < en; jj += blockSize {
		for k = en; k < A.Col(); k++ {
			for i = 0; i < A.Row(); i++ {
				r = A.At(i, k)
				for j = jj; j < jj+blockSize; j++ {
					C.Inc(i, j, r*B.At(k, j))
				}
			}
		}
	}

	// residule bottom right
	for k = en; k < A.Col(); k++ {
		for i = 0; i < A.Row(); i++ {
			r = A.At(i, k)
			for j = en; j < B.Col(); j++ {
				C.Inc(i, j, r*B.At(k, j))
			}
		}
	}
	return
}

// DotBlockKIJP block matrix multiplication
func (A *Dense) DotBlockKIJP(blockSize int, B, C Matrix) (err error) {
	if (A.Col() != B.Row()) || (C.Row() != A.Row()) || (C.Col() != B.Col()) {
		return ErrMatSize
	}
	min := A.Row()
	if A.Col() < min {
		min = A.Col()
	}
	if B.Col() < min {
		min = B.Col()
	}
	en := blockSize * (min / blockSize)
	wg := sync.WaitGroup{}
	for kk := 0; kk < en; kk += blockSize {
		for jj := 0; jj < en; jj += blockSize {
			wg.Add(1)
			go func(kk, jj int) {
				for k := kk; k < kk+blockSize; k++ {
					for i := 0; i < A.Row(); i++ {
						r := A.At(i, k)
						for j := jj; j < jj+blockSize; j++ {
							C.Inc(i, j, r*B.At(k, j))
						}
					}
				}
				wg.Done()
			}(kk, jj)
		}
		wg.Wait()
		for k := kk; k < kk+blockSize; k++ {
			wg.Add(1)
			go func(k int) {
				for i := 0; i < A.Row(); i++ {
					r := A.At(i, k)
					for j := en; j < B.Col(); j++ {
						C.Inc(i, j, r*B.At(k, j))
					}
				}
				wg.Done()
			}(k)
		}
		wg.Wait()
	}

	// residule bottom
	for jj := 0; jj < en; jj += blockSize {
		for k := en; k < A.Col(); k++ {
			wg.Add(1)
			go func(jj, k int) {
				for i := 0; i < A.Row(); i++ {
					r := A.At(i, k)
					for j := jj; j < jj+blockSize; j++ {
						C.Inc(i, j, r*B.At(k, j))
					}
				}
				wg.Done()
			}(jj, k)
		}
	}
	wg.Wait()

	// residule bottom right
	for k := en; k < A.Col(); k++ {
		for i := 0; i < A.Row(); i++ {
			wg.Add(1)
			go func(k, i int) {
				r := A.At(i, k)
				for j := en; j < B.Col(); j++ {
					C.Inc(i, j, r*B.At(k, j))
				}
				wg.Done()
			}(k, i)
		}
	}
	wg.Wait()
	return
}
