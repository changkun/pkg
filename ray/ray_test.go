package main

import (
	"testing"
)

// workaround: https://github.com/golang/go/issues/31859
var _ = func() bool {
	testing.Init()
	return true
}()

func BenchmarkRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		render()
	}
}
