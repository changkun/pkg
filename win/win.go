// Copyright 2022 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Command win opens a shiny window and prints every event it receives.
package main

import (
	"fmt"

	driver "golang.org/x/exp/shiny/driver/gldriver"
	"golang.org/x/exp/shiny/screen"
)

func main() {
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(nil)
		if err != nil {
			panic(err)
		}
		defer w.Release()

		for {
			fmt.Println(w.NextEvent())
		}
	})
}
