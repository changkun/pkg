// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package bo

import (
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/optimize"
)

var _ optimize.Method = BoundsMethod{}
var _ optimize.Statuser = BoundsMethod{}

// BoundsMethod is a hack that contains the target space.
type BoundsMethod struct {
	Method optimize.Method
	Bounds []Param
}

func (m BoundsMethod) constrain(loc *optimize.Location) {
	if loc == nil {
		return
	}

	for i, param := range m.Bounds {
		max := param.GetMax()
		min := param.GetMin()
		if loc.X[i] > max {
			loc.X[i] = max
		} else if loc.X[i] < min {
			loc.X[i] = min
		}
	}
}

// Init ...
func (m BoundsMethod) Init(dims, tasks int) int {
	return m.Method.Init(dims, tasks)
}

// Run ...
func (m BoundsMethod) Run(operation chan<- optimize.Task, result <-chan optimize.Task, tasks []optimize.Task) {
	op := make(chan optimize.Task)
	res := make(chan optimize.Task)

	go func() {
		defer close(res)

		for t := range result {
			m.constrain(t.Location)
			res <- t
		}
	}()

	go func() {
		defer close(operation)

		for t := range op {
			m.constrain(t.Location)
			operation <- t
		}
	}()

	for _, t := range tasks {
		m.constrain(t.Location)
	}
	m.Method.Run(op, res, tasks)
}

// Uses ...
func (m BoundsMethod) Uses(has optimize.Available) (uses optimize.Available, err error) {
	return m.Method.Uses(has)
}

// Status ...
func (m BoundsMethod) Status() (optimize.Status, error) {
	s, ok := m.Method.(optimize.Statuser)
	if ok {
		return s.Status()
	}
	return optimize.NotTerminated, errors.Errorf("not Statuser")
}
