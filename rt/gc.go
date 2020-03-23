package rt

import "runtime"

// GCSignal is a signal that allows you to receive a notification after
// GC cycle
type GCSignal struct {
	n *dummy
}

type dummy struct {
	done chan struct{}
	gc   chan struct{}
}

// NewGCSignal creates a GC signal for the complextion of GC cycle.
func NewGCSignal() *GCSignal {
	n := &dummy{
		done: make(chan struct{}, 1),
		gc:   make(chan struct{}),
	}
	runtime.SetFinalizer(&dummy{
		done: n.done,
		gc:   n.gc,
	}, fin)
	runtime.SetFinalizer(n, func(n *dummy) {
		select {
		case n.done <- struct{}{}:
		default: // prevent if n.done is closed
		}
	})
	return &GCSignal{n}
}

// Send sends a GC complextion signal to the given channel.
func (gc *GCSignal) Send(ch chan<- struct{}) {
	ch <- struct{}{}
}

func fin(n *dummy) {
	select {
	case <-n.done:
		close(n.gc)
		return
	default: // notifier is already closed.
	}

	select {
	case n.gc <- struct{}{}:
	default:
	}

	// prepare the next notification
	runtime.SetFinalizer(n, fin)
}
