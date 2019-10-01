// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mat

import "sync"

// DotBlockIJK block matrix multiplication
func (A *Dense) DotBlockIJK(blockSize int, B, C Matrix) (err error) {
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
		sum             float64
		en              = blockSize * (min / blockSize)
	)

	for kk = 0; kk < en; kk += blockSize {
		for jj = 0; jj < en; jj += blockSize {
			for i = 0; i < A.Row(); i++ {
				for j = jj; j < jj+blockSize; j++ {
					sum = 0.0
					for k = kk; k < kk+blockSize; k++ {
						sum += A.At(i, k) * B.At(k, j)
					}
					C.Inc(i, j, sum)
				}
			}
		}

		// residue right
		for i = 0; i < A.Row(); i++ {
			for j = en; j < B.Col(); j++ {
				sum = 0.0
				for k = kk; k < kk+blockSize; k++ {
					sum += A.At(i, k) * B.At(k, j)
				}
				C.Inc(i, j, sum)
			}
		}
	}

	// residue bottom
	for jj = 0; jj < en; jj += blockSize {
		for i = 0; i < A.Row(); i++ {
			for j = jj; j < jj+blockSize; j++ {
				sum = 0.0
				for k = en; k < A.Col(); k++ {
					sum += A.At(i, k) * B.At(k, j)
				}
				C.Inc(i, j, sum)
			}
		}
	}

	// residule bottom right
	for i = 0; i < A.Row(); i++ {
		for j = en; j < B.Col(); j++ {
			sum = 0.0
			for k = en; k < A.Col(); k++ {
				sum += A.At(i, k) * B.At(k, j)
			}
			C.Inc(i, j, sum)
		}
	}
	return
}

// DotBlockIJKP block matrix multiplication
func (A *Dense) DotBlockIJKP(blockSize int, B, C Matrix) (err error) {
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
			wg.Add(1) // per block
			go func(kk, jj int) {
				for i := 0; i < A.Row(); i++ {
					for j := jj; j < jj+blockSize; j++ {
						sum := 0.0
						for k := kk; k < kk+blockSize; k++ {
							sum += A.At(i, k) * B.At(k, j)
						}
						C.Inc(i, j, sum)
					}
				}
				wg.Done()
			}(kk, jj)
		}
		wg.Wait()

		// residue right
		for i := 0; i < A.Row(); i++ {
			wg.Add(1) // per row
			go func(i int) {
				for j := en; j < B.Col(); j++ {
					sum := 0.0
					for k := kk; k < kk+blockSize; k++ {
						sum += A.At(i, k) * B.At(k, j)
					}
					C.Inc(i, j, sum)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
	}

	// residue bottom
	for jj := 0; jj < en; jj += blockSize {
		wg.Add(1) // per row
		go func(jj int) {
			for i := 0; i < A.Row(); i++ {
				for j := jj; j < jj+blockSize; j++ {
					sum := 0.0
					for k := en; k < A.Col(); k++ {
						sum += A.At(i, k) * B.At(k, j)
					}
					C.Inc(i, j, sum)
				}
			}
			wg.Done()
		}(jj)
	}
	wg.Wait()

	// residule bottom right
	for i := 0; i < A.Row(); i++ {
		wg.Add(1) // per row
		go func(i int) {
			for j := en; j < B.Col(); j++ {
				sum := 0.0
				for k := en; k < A.Col(); k++ {
					sum += A.At(i, k) * B.At(k, j)
				}
				C.Inc(i, j, sum)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return
}
