// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ds_test

import (
	"testing"

	"changkun.de/x/pkg/ds"
)

func newSkipList() *ds.SkipList {
	return ds.NewSkipList(func(a, b any) bool {
		return a.(int) < b.(int)
	})
}

func TestNewSkipList(t *testing.T) {
	if newSkipList() == nil {
		t.Fatalf("%v: got nil", t.Name())
	}
}

func TestSkipList_Len(t *testing.T) {
	sl := newSkipList()
	if sl == nil {
		t.Fatalf("%v: got nil", t.Name())
	}

	if got := sl.Len(); got != 0 {
		t.Fatalf("Len: got %d, want %d", got, 0)
	}

	for i := 0; i < 10000; i++ {
		sl.Set(i, i)
	}

	if got := sl.Len(); got != 10000 {
		t.Fatalf("Len: got %d, want %d", got, 10000)
	}
}

func TestSkipList_GetFail(t *testing.T) {
	sl := newSkipList()
	if sl == nil {
		t.Fatalf("%v: got nil", t.Name())
	}

	v, ok := sl.Get(-1)
	if ok {
		t.Fatalf("%v: suppose to fail, but got: %v, %v", t.Name(), v, ok)
	}
}

func TestSkipList_GetSuccess(t *testing.T) {
	sl := newSkipList()
	if sl == nil {
		t.Fatalf("%v: got nil", t.Name())
	}

	sl.Set(1, 2)
	if got, ok := sl.Get(1); got != 2 || ok != true {
		t.Fatalf("got %v, %v want %v, %v", got, ok, 2, true)
	}

	sl.Set(1, 3)
	if got, ok := sl.Get(1); got != 3 || ok != true {
		t.Fatalf("got %v, %v want %v, %v", got, ok, 3, true)
	}
}

func TestSkipList_Search(t *testing.T) {
	sl := newSkipList()
	if sl == nil {
		t.Fatalf("%v: got nil", t.Name())
	}

	if ok := sl.Search(1); ok {
		t.Fatalf("got %v want %v", ok, false)
	}

	sl.Set(1, 2)

	if got := sl.Len(); got != 1 {
		t.Fatalf("Len: got %d, want %d", got, 1)
	}

	if ok := sl.Search(1); !ok {
		t.Fatalf("got %v want %v", ok, true)
	}

	if v, ok := sl.Del(1); v != 2 || !ok {
		t.Fatalf("got %v,%v want %d", v, ok, 2)
	}

	if got := sl.Len(); got != 0 {
		t.Fatalf("Len: got %d, want %d", got, 1)
	}
}

func TestSkiplist_Del(t *testing.T) {
	sl := newSkipList()
	if sl == nil {
		t.Fatalf("%v: got nil", t.Name())
	}

	for i := 0; i < 10; i++ {
		sl.Set(i, i)
	}

	for i := 0; i < 100; i++ {
		if _, ok := sl.Del(i); i > 10 && ok {
			t.Fatalf("%v: should fail, got: %v", t.Name(), ok)
		}
	}

	if got := sl.Len(); got != 0 {
		t.Fatalf("Len: got %d, want %d", got, 0)
	}
}

func TestSkipList_Range(t *testing.T) {
	sl := newSkipList()
	if sl == nil {
		t.Fatalf("%v: got nil", t.Name())
	}

	for i := 0; i < 100; i++ {
		sl.Set(i, i)
	}

	current := 10
	sl.Range(10, 20, func(v any) {
		if v != current {
			t.Fatalf("range failed, want %v, got %v", current, v)
		}
		current++
	})

	current = 90
	sl.Range(90, 120, func(v any) {
		if v != current {
			t.Fatalf("range failed, want %v, got %v", current, v)
		}
		current++
	})
	// Range covers [from, to), clipped to the end of the list, so the last
	// key 99 is included and current lands one past it.
	if current != 100 {
		t.Fatalf("range out of bound, want %v, got %v", 100, current)
	}
}

// path returns the first item whose key is not smaller than the one looked
// up, so every lookup has to compare that candidate against the key it wanted.
// The comparison read s.less(x.k, k) on both sides of an or, which is always
// false for that candidate, so each of these operations answered about a
// neighbouring key instead.

func TestSkipList_GetMissingKeyWithSuccessor(t *testing.T) {
	sl := newSkipList()
	for _, k := range []int{10, 20, 30} {
		sl.Set(k, k)
	}

	// 15 is absent, and 20 is the candidate the descent lands on.
	if v, ok := sl.Get(15); ok {
		t.Errorf("Get(15) = %v, %v, want nil, false", v, ok)
	}
	// 5 is below every key, so the candidate is the first item.
	if v, ok := sl.Get(5); ok {
		t.Errorf("Get(5) = %v, %v, want nil, false", v, ok)
	}
	// 35 is past the end and has no candidate at all.
	if v, ok := sl.Get(35); ok {
		t.Errorf("Get(35) = %v, %v, want nil, false", v, ok)
	}
}

func TestSkipList_SearchMissingKeyWithSuccessor(t *testing.T) {
	sl := newSkipList()
	for _, k := range []int{10, 20, 30} {
		sl.Set(k, k)
	}

	if sl.Search(15) {
		t.Error("Search(15) = true, want false")
	}
	if !sl.Search(20) {
		t.Error("Search(20) = false, want true")
	}
}

func TestSkipList_DelMissingKeyKeepsNeighbour(t *testing.T) {
	sl := newSkipList()
	for _, k := range []int{10, 20, 30} {
		sl.Set(k, k)
	}

	if v, ok := sl.Del(15); ok {
		t.Errorf("Del(15) = %v, %v, want nil, false", v, ok)
	}
	if got := sl.Len(); got != 3 {
		t.Errorf("Len after deleting an absent key = %d, want 3", got)
	}
	if v, ok := sl.Get(20); !ok || v != 20 {
		t.Errorf("Get(20) = %v, %v after deleting an absent key, want 20, true", v, ok)
	}
}

func TestSkipList_SetExistingKeyReplaces(t *testing.T) {
	sl := newSkipList()
	sl.Set(1, "first")
	sl.Set(1, "second")

	if got := sl.Len(); got != 1 {
		t.Errorf("Len after setting the same key twice = %d, want 1", got)
	}
	if v, ok := sl.Get(1); !ok || v != "second" {
		t.Errorf("Get(1) = %v, %v, want second, true", v, ok)
	}

	// A duplicate hides rather than replaces, so deleting once must leave
	// nothing behind.
	if _, ok := sl.Del(1); !ok {
		t.Fatal("Del(1) = false, want true")
	}
	if v, ok := sl.Get(1); ok {
		t.Errorf("Get(1) = %v, %v after delete, want nil, false", v, ok)
	}
}

func TestSkipList_RangeIncludesLastKey(t *testing.T) {
	sl := newSkipList()
	for i := range 5 {
		sl.Set(i, i)
	}

	var got []any
	sl.Range(3, 100, func(v any) { got = append(got, v) })
	if len(got) != 2 || got[0] != 3 || got[1] != 4 {
		t.Errorf("Range(3, 100) visited %v, want [3 4]", got)
	}

	// A start past the end has no candidate; Range must return, not panic.
	got = nil
	sl.Range(99, 100, func(v any) { got = append(got, v) })
	if len(got) != 0 {
		t.Errorf("Range(99, 100) visited %v, want nothing", got)
	}
}
