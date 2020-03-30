// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// +build linux

package mkill

import "fmt"

var cmdThreads = fmt.Sprintf("ps hH p %d | wc -l", pid)
