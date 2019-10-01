// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

import "sync"

// DotBlockJKI block matrix multiplication
func (A *Dense) DotBlockJKI(blockSize int, B, C Matrix) (err error) {
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
			for j = jj; j < jj+blockSize; j++ {
				for k = kk; k < kk+blockSize; k++ {
					r = B.At(k, j)
					for i = 0; i < A.Row(); i++ {
						C.Inc(i, j, r*A.At(i, k))
					}
				}
			}
		}
		for j = en; j < B.Col(); j++ {
			for k = kk; k < kk+blockSize; k++ {
				r = B.At(k, j)
				for i = 0; i < A.Row(); i++ {
					C.Inc(i, j, r*A.At(i, k))
				}
			}
		}
	}

	// residule bottom
	for jj = 0; jj < en; jj += blockSize {
		for j = jj; j < jj+blockSize; j++ {
			for k = en; k < A.Col(); k++ {
				r = B.At(k, j)
				for i = 0; i < A.Row(); i++ {
					C.Inc(i, j, r*A.At(i, k))
				}
			}
		}
	}

	// residule bottom right
	for j = en; j < B.Col(); j++ {
		for k = en; k < A.Col(); k++ {
			r = B.At(k, j)
			for i = 0; i < A.Row(); i++ {
				C.Inc(i, j, r*A.At(i, k))
			}
		}
	}
	return
}

// DotBlockJKIP block matrix multiplication
func (A *Dense) DotBlockJKIP(blockSize int, B, C Matrix) (err error) {
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
				for j := jj; j < jj+blockSize; j++ {
					for k := kk; k < kk+blockSize; k++ {
						r := B.At(k, j)
						for i := 0; i < A.Row(); i++ {
							C.Inc(i, j, r*A.At(i, k))
						}
					}
				}
				wg.Done()
			}(kk, jj)
		}
		wg.Wait()
		for j := en; j < B.Col(); j++ {
			wg.Add(1)
			go func(j int) {
				for k := kk; k < kk+blockSize; k++ {
					r := B.At(k, j)
					for i := 0; i < A.Row(); i++ {
						C.Inc(i, j, r*A.At(i, k))
					}
				}
				wg.Done()
			}(j)
		}
		wg.Wait()
	}

	// residule bottom
	for jj := 0; jj < en; jj += blockSize {
		wg.Add(1)
		go func(jj int) {
			for j := jj; j < jj+blockSize; j++ {
				for k := en; k < A.Col(); k++ {
					r := B.At(k, j)
					for i := 0; i < A.Row(); i++ {
						C.Inc(i, j, r*A.At(i, k))
					}
				}
			}
			wg.Done()
		}(jj)
	}
	wg.Wait()

	// residule bottom right
	for j := en; j < B.Col(); j++ {
		wg.Add(1)
		go func(j int) {
			for k := en; k < A.Col(); k++ {
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
