package syscall

import (
	"net"
	"os"
	"syscall"
	"testing"
)

const message = "hello, world!"

var buffer = make([]byte, 13)

func writeAll(fd int, buf []byte) error {
	for len(buf) > 0 {
		n, err := syscall.Write(fd, buf)
		if n < 0 {
			return err
		}
		buf = buf[n:]
	}
	return nil
}

func BenchmarkReadWriteCgoCalls(b *testing.B) {
	fds, _ := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	for i := 0; i < b.N; i++ {
		CwriteAll(fds[0], []byte(message))
		Cread(fds[1], buffer)
	}
}

func BenchmarkReadWriteGoCalls(b *testing.B) {
	fds, _ := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	for i := 0; i < b.N; i++ {
		writeAll(fds[0], []byte(message))
		syscall.Read(fds[1], buffer)
	}
}

func BenchmarkReadWriteNetCalls(b *testing.B) {
	cs, _ := socketpair()
	for i := 0; i < b.N; i++ {
		cs[0].Write([]byte(message))
		cs[1].Read(buffer)
	}
}

func socketpair() (conns [2]net.Conn, err error) {
	fds, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM, 0)
	if err != nil {
		return
	}
	conns[0], err = fdToFileConn(fds[0])
	if err != nil {
		return
	}
	conns[1], err = fdToFileConn(fds[1])
	if err != nil {
		conns[0].Close()
		return
	}
	return
}

func fdToFileConn(fd int) (net.Conn, error) {
	f := os.NewFile(uintptr(fd), "")
	defer f.Close()
	return net.FileConn(f)
}
