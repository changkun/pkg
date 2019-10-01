// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lockfree

import (
	"sync/atomic"
	"unsafe"
)

type directItem struct {
	next unsafe.Pointer
	v    interface{}
}

func loaditem(p *unsafe.Pointer) *directItem {
	return (*directItem)(atomic.LoadPointer(p))
}
func casitem(p *unsafe.Pointer, old, new *directItem) bool {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}
