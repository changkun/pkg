// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ds_test

import (
	"testing"

	"github.com/changkun/gobase/ds"
)

func TestLRU(t *testing.T) {
	lru := ds.NewLRU(2)
	lru.Put(1, 1)
	lru.Put(2, 2) // 2, 1

	if lru.Get(1) != 1 { // 1, 2
		t.Fatalf("want 1")
	}
	lru.Put(3, 3)         // 3, 1
	if lru.Get(2) != -1 { // 3, 1
		t.Fatalf("want -1")
	}
	if lru.Get(1) != 1 { // 1, 3
		t.Fatalf("want 1")
	}
	if lru.Get(3) != 3 { // 3, 1
		t.Fatalf("want 3")
	}
	if lru.Get(4) != -1 { // 3, 1
		t.Fatalf("want 4")
	}
}
