// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ds

import "math/rand"

// RandomizedSet implements O(1) randomized set
type RandomizedSet struct {
	store map[int]int
}

// NewRandomizedSet ...
func NewRandomizedSet() RandomizedSet {
	return RandomizedSet{store: map[int]int{}}
}

// Insert inserts a value to the set. Returns true if the set did not already contain the specified element.
func (rs *RandomizedSet) Insert(val int) bool {
	if _, ok := rs.store[val]; !ok {
		rs.store[val] = val
		return true
	}
	rs.store[val] = val
	return false
}

// Remove removes a value from the set. Returns true if the set contained the specified element.
func (rs *RandomizedSet) Remove(val int) bool {
	if _, ok := rs.store[val]; ok {
		delete(rs.store, val)
		return true
	}
	return false
}

// Get gets a random element from the set.
func (rs *RandomizedSet) Get() int {
	r := rand.Intn(len(rs.store))
	count := 0
	for k := range rs.store {
		if count == r {
			return k
		}
		count++
	}
	panic("never")
}

func randIntMapKey(m map[int]int) int {
	i := rand.Intn(len(m))
	for k := range m {
		if i == 0 {
			return k
		}
		i--
	}
	panic("never")
}
