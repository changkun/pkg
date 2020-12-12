// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// +build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework Carbon
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import <Carbon/Carbon.h>
extern void go_hotkey_callback(void* handler);
static OSStatus _hotkey_handler(EventHandlerCallRef nextHandler, EventRef theEvent, void *userData);
int register_hotkey(void* go_hotkey_handler);
void app_run();
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

//export go_hotkey_callback
func go_hotkey_callback(c unsafe.Pointer) { (*gocallback)(c).call() }

type gocallback struct{ f func() }

func (c *gocallback) call() { c.f() }

func init() {
	runtime.LockOSThread()
}

var hkCallback func()

// RegisterHotKey registers fn as a callback for a configured global
// hotkey.
//
// This function must run on the main thread.
func RegisterHotKey(f func()) {
	hkCallback = f
	arg := unsafe.Pointer(&gocallback{func() {
		// This cannot be a direct call to f().
		// Because it can cause runtime error: cgo argument has
		// a Go pointer to Go pointer.
		hkCallback()
	}})
	ret := C.register_hotkey(arg)
	if ret == C.int(-1) {
		fmt.Println("register global system hotkey failed.")
	}
	fmt.Println("hotkey ctrl+option+s is registered")
	C.app_run()
}

func main() {
	hkCallback = func() {
		fmt.Println("hotkey is triggered")
	}
	RegisterHotKey(hkCallback)
}
