package gocall

/*
void empty() {
}
*/
import "C"

// Cempty is an empty cgo call
func Cempty() {
	C.empty()
}

// Empty is an empty go call
func Empty() {}
