// Copyright 2021 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package convert_test

import (
	"testing"

	"changkun.de/x/pkg/convert"
)

func BenchmarkLinear2sRGB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0.0; j <= 1.0; j += 0.01 {
			convert.Linear2sRGB(j)
		}
	}
}
