// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ds

import (
	"math/rand"

	"changkun.de/x/pkg/common"
)

// A SkipList maintains an ordered collection of key:valkue pairs.
// It support insertion, lookup, and deletion operations with O(log n) time complexity
// Paper: Pugh, William (June 1990). "Skip lists: a probabilistic alternative to balanced
// trees". Communications of the ACM 33 (6): 668–676
type SkipList struct {
	header   *skiplistitem
	len      int
	MaxLevel int
	less     common.Less
}

// NewSkipList returns a skiplist.
func NewSkipList(less common.Less) *SkipList {
	return &SkipList{
		header:   &skiplistitem{forward: []*skiplistitem{nil}},
		MaxLevel: 32,
		less:     less,
	}
}

// Len returns the length of given skiplist.
func (s *SkipList) Len() int {
	return s.len
}

// equal reports whether a and b are the same key under the list's ordering.
// The list only knows how to compare keys with less, so two keys are equal
// exactly when neither is smaller than the other.
func (s *SkipList) equal(a, b any) bool {
	return !s.less(a, b) && !s.less(b, a)
}

// Set sets given k and v pair into the skiplist.
func (s *SkipList) Set(k, v any) {
	// s.level starts from 0, we need to allocate one
	update := make([]*skiplistitem, s.level()+1, s.effectiveMaxLevel()+1) // make(type, len, cap)

	x := s.path(s.header, update, k)
	if x != nil && s.equal(x.k, k) { // if key exists, update
		x.v = v
		return
	}

	newl := s.randomLevel()
	if curl := s.level(); newl > curl {
		for i := curl + 1; i <= newl; i++ {
			update = append(update, s.header)
			s.header.forward = append(s.header.forward, nil)
		}
	}

	item := &skiplistitem{
		forward: make([]*skiplistitem, newl+1, s.effectiveMaxLevel()+1),
		k:       k,
		v:       v,
	}
	for i := 0; i <= newl; i++ {
		item.forward[i] = update[i].forward[i]
		update[i].forward[i] = item
	}

	s.len++
}

// path descends the list to the first item whose key is not smaller than k,
// recording in update the last item visited at each level. The returned
// candidate is the only item that can hold k, but its key may be larger, so
// callers must still compare it against k.
func (s *SkipList) path(x *skiplistitem, update []*skiplistitem, k any) (candidate *skiplistitem) {
	depth := len(x.forward) - 1
	for i := depth; i >= 0; i-- {
		for x.forward[i] != nil && s.less(x.forward[i].k, k) {
			x = x.forward[i]
		}
		if update != nil {
			update[i] = x
		}
	}
	return x.next()
}

func (s *SkipList) randomLevel() (n int) {
	for n = 0; n < s.effectiveMaxLevel() && rand.Float64() < 0.25; n++ {
	}
	return
}

// Get returns corresponding v with given k.
func (s *SkipList) Get(k any) (v any, ok bool) {
	x := s.path(s.header, nil, k)
	if x == nil || !s.equal(x.k, k) {
		return nil, false
	}
	return x.v, true
}

// Search returns true if k is found in the skiplist.
func (s *SkipList) Search(k any) bool {
	x := s.path(s.header, nil, k)
	return x != nil && s.equal(x.k, k)
}

// Range calls op with the value of every key in [from, to), in order. It
// stops at the end of the list, so a `to` past the last key is not an error.
func (s *SkipList) Range(from, to any, op func(v any)) {
	for x := s.path(s.header, nil, from); x != nil; x = x.next() {
		if !s.less(x.k, to) {
			return
		}
		op(x.v)
	}
}

// Del returns the deleted value if ok
func (s *SkipList) Del(k any) (v any, ok bool) {
	update := make([]*skiplistitem, s.level()+1, s.effectiveMaxLevel())

	x := s.path(s.header, update, k)
	if x == nil || !s.equal(x.k, k) {
		return nil, false
	}

	v = x.v
	for i := 0; i <= s.level() && update[i].forward[i] == x; i++ {
		update[i].forward[i] = x.forward[i]
	}
	for s.level() > 0 && s.header.forward[s.level()] == nil {
		s.header.forward = s.header.forward[:s.level()]
	}
	s.len--
	ok = true
	return
}

func (s *SkipList) level() int {
	return len(s.header.forward) - 1
}

func (s *SkipList) effectiveMaxLevel() int {
	if s.level() < s.MaxLevel {
		return s.MaxLevel
	}
	return s.level()
}

type skiplistitem struct {
	forward []*skiplistitem
	k       any
	v       any
}

func (s *skiplistitem) next() *skiplistitem {
	if len(s.forward) == 0 {
		return nil
	}
	return s.forward[0]
}
