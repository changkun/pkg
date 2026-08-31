// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package syscall

import (
	"net"
	"os"
	"syscall"
	"testing"
)

const message = "hello, world!"

var buffer = make([]byte, len(message))

// writeAll writes buf in full through the Go syscall wrapper.
func writeAll(fd int, buf []byte) error {
	for len(buf) > 0 {
		n, err := syscall.Write(fd, buf)
		if err != nil {
			return err
		}
		buf = buf[n:]
	}
	return nil
}

// rawPair returns a connected pair of raw descriptors, closed when the
// benchmark ends. The descriptors used to be leaked once per benchmark.
func rawPair(b *testing.B) [2]int {
	b.Helper()
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		b.Fatalf("socketpair: %v", err)
	}
	b.Cleanup(func() {
		syscall.Close(fds[0])
		syscall.Close(fds[1])
	})
	return [2]int{fds[0], fds[1]}
}

func BenchmarkReadWriteCgoCalls(b *testing.B) {
	fds := rawPair(b)
	for b.Loop() {
		if err := CwriteAll(fds[0], []byte(message)); err != nil {
			b.Fatalf("cgo write: %v", err)
		}
		if _, err := Cread(fds[1], buffer); err != nil {
			b.Fatalf("cgo read: %v", err)
		}
	}
}

func BenchmarkReadWriteGoCalls(b *testing.B) {
	fds := rawPair(b)
	for b.Loop() {
		if err := writeAll(fds[0], []byte(message)); err != nil {
			b.Fatalf("write: %v", err)
		}
		if _, err := syscall.Read(fds[1], buffer); err != nil {
			b.Fatalf("read: %v", err)
		}
	}
}

func BenchmarkReadWriteNetCalls(b *testing.B) {
	cs, err := socketpair()
	if err != nil {
		b.Fatalf("socketpair: %v", err)
	}
	b.Cleanup(func() {
		cs[0].Close()
		cs[1].Close()
	})

	for b.Loop() {
		if _, err := cs[0].Write([]byte(message)); err != nil {
			b.Fatalf("write: %v", err)
		}
		if _, err := cs[1].Read(buffer); err != nil {
			b.Fatalf("read: %v", err)
		}
	}
}

func socketpair() (conns [2]net.Conn, err error) {
	fds, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM, 0)
	if err != nil {
		return conns, err
	}
	conns[0], err = fdToFileConn(fds[0])
	if err != nil {
		syscall.Close(fds[1])
		return conns, err
	}
	conns[1], err = fdToFileConn(fds[1])
	if err != nil {
		conns[0].Close()
		return conns, err
	}
	return conns, nil
}

// fdToFileConn wraps a raw descriptor as a net.Conn. net.FileConn duplicates
// the descriptor, so closing the os.File here does not close the connection.
func fdToFileConn(fd int) (net.Conn, error) {
	f := os.NewFile(uintptr(fd), "")
	defer f.Close()
	return net.FileConn(f)
}
