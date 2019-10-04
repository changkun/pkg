package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

const total int = 100

var (
	cho [total]int64
	per [total]int64
)

func main() {
	for i := range per {
		go func(i int) {
			left := i
			right := i + 1
			if right >= total {
				right = 0
			}
			println("left: ", left, ", right: ", right)
			time.Sleep(time.Second)

			// MOCK DEAD LOCK
			// for {
			// 	if atomic.CompareAndSwapInt64(&cho[right], 0, 1) {
			// 		for {
			// 			if atomic.CompareAndSwapInt64(&cho[left], 0, 1) {
			// 				// both success. put left and right back:
			// 				atomic.StoreInt64(&cho[left], 0)
			// 				atomic.StoreInt64(&cho[right], 0)
			// 				fmt.Printf("per %d: eat\n", i)
			// 				break
			// 			}
			// 			runtime.Gosched()
			// 		}
			// 	}
			// 	runtime.Gosched()
			// }

			for {
				if atomic.CompareAndSwapInt64(&cho[right], 0, 1) {
					// get right, now try left
					if atomic.CompareAndSwapInt64(&cho[left], 0, 1) {
						// both success. put left back:
						atomic.StoreInt64(&cho[left], 0)
						atomic.StoreInt64(&cho[right], 0)
						fmt.Printf("per %d: eat\n", i)
					} else {
						// fail to get left, put right back, try next time
						atomic.StoreInt64(&cho[right], 0)
					}
				}
				runtime.Gosched()
			}
		}(i)
	}
	select {}
}
