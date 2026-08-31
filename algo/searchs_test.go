// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package algo_test

import (
	"testing"

	"changkun.de/x/pkg/algo"
	"changkun.de/x/pkg/common"
)

func TestBinarySearch(t *testing.T) {
	tests := []struct {
		input []any
		x     any
		less  common.Less
		want  int
	}{
		{
			input: []any{1, 2, 3, 4, 5, 6, 7},
			x:     6,
			less: func(a, b any) bool {
				return a.(int) < b.(int)
			},
			want: 5,
		},
		{
			input: []any{1, 2, 3, 4, 5, 6, 7},
			x:     2,
			less: func(a, b any) bool {
				return a.(int) < b.(int)
			},
			want: 1,
		},
		{
			input: []any{},
			x:     2,
			less: func(a, b any) bool {
				return a.(int) < b.(int)
			},
			want: -1,
		},
	}

	for _, tt := range tests {
		r := algo.BinarySearch(tt.input, tt.x, tt.less)
		if r != tt.want {
			t.Fatalf("BinarySearch %v of %v: want %v, got %v", tt.x, tt.input, tt.want, r)
		}
	}
}
