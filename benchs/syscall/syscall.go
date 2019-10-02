package syscall

/*
#include <unistd.h>
int write_all(int fd, void* buffer, size_t length) {
    while (length > 0) {
        int written = write(fd, buffer, length);
        if (written < 0)
            return -1;
        length -= written;
        buffer += written;
    }
    return length;
}
int read_call(int fd, void *buffer, size_t length) {
	return read(fd, buffer, length);
}
*/
import "C"
import (
	"unsafe"
)

// CwriteAll is a cgo call for write
func CwriteAll(fd int, buf []byte) error {
	_, err := C.write_all(C.int(fd), unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	return err
}

// Cread is a cgo call for read
func Cread(fd int, buf []byte) (int, error) {
	ret, err := C.read_call(C.int(fd), unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	return int(ret), err
}
