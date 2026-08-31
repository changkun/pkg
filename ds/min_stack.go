// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ds

import "container/list"

// MinStack ...
type MinStack struct {
	store *list.List
}
type node struct {
	min   int
	value int
}

// NewMinStack ...
func NewMinStack() MinStack {
	return MinStack{store: list.New()}
}

// Push ...
func (ms *MinStack) Push(x int) {
	// The new minimum is the smaller of x and the minimum carried by the item
	// below it. An empty stack has no item below, so the minimum is x.
	min := x
	if front := ms.store.Front(); front != nil && front.Value.(*node).min <= x {
		min = front.Value.(*node).min
	}
	ms.store.PushFront(&node{value: x, min: min})
}

// Pop ...
func (ms *MinStack) Pop() {
	ms.store.Remove(ms.store.Front())
}

// Top ...
func (ms *MinStack) Top() int {
	return ms.store.Front().Value.(*node).value
}

// GetMin ...
func (ms *MinStack) GetMin() int {
	return ms.store.Front().Value.(*node).min
}
