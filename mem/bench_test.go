package main_test

import (
	"syscall"
	"testing"
)

var pageSize = syscall.Getpagesize()

func BenchmarkPrefetch(b *testing.B) {
	for j := 0; j < b.N; j++ {
		b.StopTimer()
		anonMB := 10 << 20 // MiB
		m, err := syscall.Mmap(-1, 0, anonMB, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
		if err != nil {
			panic(err)
		}

		err = syscall.Madvise(m, syscall.MADV_HUGEPAGE)
		if err != nil {
			panic(err)
		}
		b.StartTimer()
		for i := 0; i < len(m); i += pageSize {
			m[i] = 42
		}
		b.StopTimer()
		err = syscall.Madvise(m, syscall.MADV_DONTNEED)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkPageFault(b *testing.B) {
	for j := 0; j < b.N; j++ {
		b.StopTimer()
		anonMB := 10 << 20 // MiB
		m, err := syscall.Mmap(-1, 0, anonMB, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
		if err != nil {
			panic(err)
		}

		b.StartTimer()
		for i := 0; i < len(m); i += pageSize {
			m[i] = 42
		}
		b.StopTimer()
		err = syscall.Madvise(m, syscall.MADV_DONTNEED)
		if err != nil {
			panic(err)
		}
	}
}
