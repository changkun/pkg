// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package slice_test

import (
	"testing"

	"github.com/changkun/gobase/slice"
)

func TestIsStringsContains(t *testing.T) {
	testdata := []string{"a", "a", "a", "a", "a", "b"}
	if !slice.IsStringsContains(testdata, "b") {
		t.Fatalf("want true, got false")
	}
}
