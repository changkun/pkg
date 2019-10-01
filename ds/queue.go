// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ds

// Queue is a FIFO queue
type Queue struct {
	v []interface{}
}

// NewQueue returns a queue
func NewQueue() *Queue {
	return &Queue{v: make([]interface{}, 0)}
}

// Enqueue enqueues a value to the tail of queue
func (q *Queue) Enqueue(v interface{}) {
	q.v = append(q.v, v)
}

// Dequeue dequeues a value from the head of queue
func (q *Queue) Dequeue() interface{} {
	if len(q.v) == 0 {
		return nil
	}
	v := q.v[0]
	q.v = q.v[1:]
	return v
}
