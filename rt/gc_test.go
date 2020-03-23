package rt_test

import (
	"sync/atomic"
	"testing"

	"github.com/changkun/gobase/rt"
)

func TestGCNotification(t *testing.T) {
	go func() {
		for {
			_ = make([]byte, 1024)
		}
	}()

	s := rt.NewGCSignal()
	ch := make(chan struct{})
	var done int64
	go func() {
		for atomic.LoadInt64(&done) < 100 {
			s.Send(ch)
		}
	}()

	for atomic.LoadInt64(&done) < 100 {
		t.Logf("GC completion: %v", <-ch)
		atomic.AddInt64(&done, 1)
	}
}
