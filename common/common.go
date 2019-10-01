// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package common

// Less defines a function that compares the order of a and b.
// Returns true if a < b
type Less func(a, b interface{}) bool
