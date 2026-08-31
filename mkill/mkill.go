// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mkill

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// threadCountTimeout bounds the shell out that counts threads. It is well
// under interval, so a slow count cannot pile up behind the next tick.
const threadCountTimeout = 500 * time.Millisecond

var (
	pid       = os.Getpid()
	maxThread = int32(runtime.NumCPU())
	interval  = time.Second
	debug     uint32
)

func checkwork() {
	_, err := getThreads()
	if err != nil {
		panic(fmt.Sprintf("mkill: failed to use the library: %v", err))
	}
}

func init() {
	checkwork()

	if atomic.LoadUint32(&debug) == 1 {
		fmt.Printf("mkill: pid %v, maxThread %v, interval %v\n", pid, maxThread, interval)
	}
	go func() {
		t := time.NewTicker(interval)
		for range t.C {
			n, err := getThreads()
			if err != nil {
				if atomic.LoadUint32(&debug) == 1 {
					fmt.Printf("mkill: %v\n", err)
				}
				continue
			}
			nkill := int32(n) - atomic.LoadInt32(&maxThread)
			if nkill <= 0 {
				if atomic.LoadUint32(&debug) == 1 {
					fmt.Printf("mkill: checked #threads total %v / max %v\n", n, maxThread)
				}
				continue
			}
			for i := int32(0); i < nkill; i++ {
				go func() {
					runtime.LockOSThread()
				}()
			}
			if atomic.LoadUint32(&debug) == 1 {
				fmt.Printf("mkill: killing #threads, remaining: %v\n", n)
			}
		}
	}()
}

// GOMAXTHREADS change the limits of the maximum threads in runtime
// and returns the previous number of threads limit
func GOMAXTHREADS(n int) int {
	return int(atomic.SwapInt32(&maxThread, int32(n)))
}

// SetDebug enables debug information for mkill
func SetDebug(flag bool) {
	if flag {
		atomic.StoreUint32(&debug, 1)
	} else {
		atomic.StoreUint32(&debug, 0)
	}
}

// getThreads returns the number of running threads. The shell out is bounded
// so a wedged command cannot stall the monitor goroutine forever.
func getThreads() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), threadCountTimeout)
	defer cancel()

	out, err := exec.CommandContext(ctx, "bash", "-c", cmdThreads).Output()
	if err != nil {
		return 0, fmt.Errorf("mkill: failed to fetch #threads: %w", err)
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, fmt.Errorf("mkill: failed to parse #threads: %w", err)
	}
	return n, nil
}
