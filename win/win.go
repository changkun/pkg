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
			switch e := w.NextEvent().(type) {
			default:
				fmt.Println(e)
			}
		}
	})
}
