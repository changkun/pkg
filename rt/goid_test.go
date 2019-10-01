// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rt_test

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/changkun/gobase/rt"
)

func slowGetGoID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, _ := strconv.ParseInt(idField, 10, 64) // very unlikely to be failed
	return id
}

func ExampleGoID() {
	cnums := make(chan int, 100)
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			cnums <- int(rt.GoID()) // down cast, wrong in large goid
			wg.Done()
		}()
	}
	wg.Wait()
	close(cnums)

	nums := []int{int(rt.GoID())}
	for v := range cnums {
		nums = append(nums, v)
	}
	sort.Ints(nums)
	fmt.Printf("%v", nums)
}

func TestGet(t *testing.T) {
	got := rt.GoID()
	want := slowGetGoID()
	if got != uint64(want) {
		t.Errorf("want %d, got: %d", want, got)
	}
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rt.GoID()
	}
}
