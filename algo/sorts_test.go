// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package algo_test

import (
	"reflect"
	"testing"

	"github.com/changkun/gobase/algo"
)

func TestSorts(t *testing.T) {

	type TestCase struct {
		input []int
		want  []int
	}

	sortfunc := map[string]func([]int) []int{
		"insert": algo.InsertSort,
		"merge":  algo.MergeSort,
	}

	for name, sort := range sortfunc {
		t.Run(name, func(t *testing.T) {
			tc := []TestCase{
				TestCase{input: []int{1}, want: []int{1}},
				TestCase{input: []int{1, 2, 3}, want: []int{1, 2, 3}},
				TestCase{input: []int{3, 2, 1}, want: []int{1, 2, 3}},
				TestCase{input: []int{1, 6, 3, 7, 7, 34, 8, 9, 3, 9}, want: []int{1, 3, 3, 6, 7, 7, 8, 9, 9, 34}},
			}
			for _, tt := range tc {
				result := sort(tt.input)
				if !reflect.DeepEqual(result, tt.want) {
					t.Errorf("insert sort error, want: %v, got: %v", tt.want, result)
				}
			}
		})
	}

}
