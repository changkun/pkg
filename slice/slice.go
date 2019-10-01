// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package slice

// IsStringsContains check if string slice contains specfied target string.
func IsStringsContains(strings []string, target string) bool {
	for _, v := range strings {
		if v == target {
			return true
		}
	}
	return false
}
