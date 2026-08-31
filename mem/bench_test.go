// Copyright 2022 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// MADV_HUGEPAGE is a Linux extension to madvise, and syscall.Madvise itself is
// only defined on Linux. The comparison these benchmarks draw has no meaning
// elsewhere.
//go:build linux

package mem_test

import (
	"syscall"
	"testing"
)

var pageSize = syscall.Getpagesize()

const anonBytes = 10 << 20 // MiB

// BenchmarkPrefetch touches one byte per page over an anonymous mapping that
// was advised as a huge page, so the faults are served in large blocks.
func BenchmarkPrefetch(b *testing.B) {
	for j := 0; j < b.N; j++ {
		b.StopTimer()
		m, err := syscall.Mmap(-1, 0, anonBytes, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
		if err != nil {
			b.Fatalf("mmap: %v", err)
		}

		if err := syscall.Madvise(m, syscall.MADV_HUGEPAGE); err != nil {
			b.Fatalf("madvise MADV_HUGEPAGE: %v", err)
		}
		b.StartTimer()
		for i := 0; i < len(m); i += pageSize {
			m[i] = 42
		}
		b.StopTimer()

		if err := syscall.Madvise(m, syscall.MADV_DONTNEED); err != nil {
			b.Fatalf("madvise MADV_DONTNEED: %v", err)
		}
		if err := syscall.Munmap(m); err != nil {
			b.Fatalf("munmap: %v", err)
		}
	}
}

// BenchmarkPageFault is the same walk without the huge page hint, so every
// page costs a separate fault.
func BenchmarkPageFault(b *testing.B) {
	for j := 0; j < b.N; j++ {
		b.StopTimer()
		m, err := syscall.Mmap(-1, 0, anonBytes, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
		if err != nil {
			b.Fatalf("mmap: %v", err)
		}

		b.StartTimer()
		for i := 0; i < len(m); i += pageSize {
			m[i] = 42
		}
		b.StopTimer()

		if err := syscall.Madvise(m, syscall.MADV_DONTNEED); err != nil {
			b.Fatalf("madvise MADV_DONTNEED: %v", err)
		}
		if err := syscall.Munmap(m); err != nil {
			b.Fatalf("munmap: %v", err)
		}
	}
}
