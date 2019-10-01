// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gen_test

import (
	"testing"

	"github.com/changkun/gobase/gen"
)

func TestRandomString(t *testing.T) {
	s := gen.RandomString(12)
	if len(s) != 12 {
		t.Fatalf("want 12 chars, got: %v", s)
	}
}
