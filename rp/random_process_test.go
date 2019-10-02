// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rp

import (
	"testing"
)

func TestRandomProcess(t *testing.T) {

	t.Run("store", func(t *testing.T) {
		counts := []float64{
			10, 2, 3, 4, 5, 6, 6, 5, 4, 3, 100,
			10, 2, 3, 4, 5, 6, 6, 5, 4, 3, 100,
			10, 2, 3, 4, 5, 6, 6, 5, 4, 3, 100,
			10, 2, 3, 4, 5, 6, 6, 5, 4, 3, 100,
			10, 2, 3, 4, 5, 6, 6, 5, 4, 3, 100,
			10, 2, 3, 4, 5, 6, 6, 5, 4, 3, 100,
			10, 2, 3, 4, 5, 6, 6, 5, 4, 3, 100,
			10, 2, 3, 4, 5, 6, 6, 5, 4, 3, 100,
			10, 2, 3, 4, 5, 6, 6, 5, 4, 3, 100,
			10, 2, 3, 4, 5, 6, 6, 5, 4, 3, 100,
		}
		window := &countProcess{maxsize: 60}
		lc := len(counts)
		for i := 0; i < len(counts); i++ {
			window.Store(counts[i])
		}

		if window.maxsize != float64(len(window.nevents)) {
			t.Fatalf("want %v, got %v", window.maxsize, float64(len(window.nevents)))
		}

		// check all elements
		lw := len(window.nevents)
		for i := 0; i < int(window.maxsize); i++ {
			if window.nevents[lw-i-1] != counts[lc-i-1] {
				t.Fatalf("want %v, got %v, index: %v", counts[lc-i-1], window.nevents[lw-i-1], i)
			}
		}
	})

	t.Run("Significant", func(t *testing.T) {
		ws := map[*countProcess]bool{
			// 100 is a significant bigger than previous one
			&countProcess{
				maxsize:    60,
				nevents:    []float64{10, 2, 3, 4, 5, 6, 6, 5, 4, 3, 2, 100},
				confidence: 0.05,
			}: true,
			// 10 is not significant bigger enough than previous one
			&countProcess{
				maxsize:    60,
				nevents:    []float64{10, 2, 3, 4, 5, 6, 6, 5, 4, 3, 2, 10},
				confidence: 0.05,
			}: false,
		}

		for w, sig := range ws {
			if sig != w.Significant() {
				t.Fatalf("want %v, got %v", sig, w.Significant())
			}
		}
	})

	t.Run("Acceptable", func(t *testing.T) {
		// we have 44 remaining slots for the events.
		var remains float64 = 44

		ws := map[*countProcess]bool{
			// 44 remaining slots isn't able to handle next sample
			// because there is a significant increasing on #events
			&countProcess{
				maxsize:    60,
				nevents:    []float64{1, 2, 3, 4, 8, 16, 32, 64, 100},
				confidence: 0.05,
			}: false,
			&countProcess{
				maxsize:    60,
				nevents:    []float64{1, 2, 4, 8, 32, 64},
				confidence: 0.05,
			}: false,
			// 44 remaining slots is able to handle next sample
			&countProcess{
				maxsize:    60,
				nevents:    []float64{10, 2, 3, 4, 5, 6, 6, 5, 4, 3, 2, 10},
				confidence: 0.05,
			}: true,
			&countProcess{
				maxsize:    60,
				nevents:    []float64{0},
				confidence: 0.05,
			}: true,
		}

		for w, s := range ws {
			if _, ss := w.Acceptable(remains); s != ss {
				t.Fatalf("want %v, got %v", s, ss)
			}
		}

		// if emaining slots are very low, but there are (literaly)
		// no requests at all, we think this is acceptable.
		if _, ok := (&countProcess{
			maxsize:    60,
			nevents:    []float64{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
			confidence: 0.05,
		}).Acceptable(2); !ok {
			t.Fatalf("want true, got false")
		}
	})

}
