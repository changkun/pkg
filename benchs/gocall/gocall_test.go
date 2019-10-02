package gocall_test

import (
	"testing"

	"github.com/changkun/gobase/benchs/gocall"
)

func BenchmarkEmptyCgoCalls(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gocall.Cempty()
	}
}

func BenchmarkEmptyGoCalls(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gocall.Empty()
	}
}
