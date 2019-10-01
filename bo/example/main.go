// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/aybabtme/rgbterm"
	"github.com/changkun/gobase/bo"
)

// predict gives a use case on how to predict next best point x
func predict() {
	N := 5

	// say a target space is N dimensions
	variables := []bo.Param{}
	for i := 0; i < N; i++ {
		variables = append(variables, bo.UniformParam{
			Max:  100,
			Min:  0,
			Name: fmt.Sprintf("V-%d", i),
		})
	}
	o := bo.NewOptimizer(variables, bo.WithMinimize(false))

	// generate random histXs
	genX := func() (map[bo.Param]float64, float64) {
		X := map[bo.Param]float64{}
		for i := 0; i < N; i++ {
			X[variables[i]] = rand.Float64() * 80
		}
		y := rand.Float64() * 5
		return X, y
	}

	histXs := []map[bo.Param]float64{}
	histYs := []float64{}
	// generate 10 history
	for i := 0; i < 10; i++ {
		X, y := genX()
		histXs = append(histXs, X)
		histYs = append(histYs, y)
	}

	// predict next best point.
	x, err := o.Predict(histXs, histYs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("predict: %+v", x)
}

// guess implements an active user modeling regarding guess your preferred color
func guess() {
	X := bo.UniformParam{
		Max:  1,
		Min:  0,
		Name: "X",
	}
	Y := bo.UniformParam{
		Max:  1,
		Min:  0,
		Name: "Y",
	}
	Z := bo.UniformParam{
		Max:  1,
		Min:  0,
		Name: "Z",
	}
	o := bo.NewOptimizer(
		[]bo.Param{
			X, Y, Z,
		},
		bo.WithMinimize(false),
		bo.WithRounds(10),
	)
	x, y, err := o.RunSerial(func(params map[bo.Param]float64) float64 {
		for {
			r, g, b := uint8(params[X]*255), uint8(params[Y]*255), uint8(params[Z]*255)
			word := "██████"
			coloredWord := rgbterm.FgString(word, r, g, b)
			fmt.Println("Is this your prefered color? ", coloredWord, " (", r, ", ", g, ", ", b, ")")
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Grade? (0 dislike to 5 like):")
			text, _ := reader.ReadString('\n')
			res, err := strconv.ParseFloat(text[:len(text)-1], 64)
			if err != nil {
				fmt.Println("err: ", err)
				continue
			}
			fmt.Println("res: ", res)
			if res >= 0 && res <= 5 {
				return res
			}
			fmt.Println("Wrong input, please input a number between 0 ~ 5!")
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	r, g, b := uint8(x[X]*255), uint8(x[Y]*255), uint8(x[Z]*255)
	word := "██████"
	coloredWord := rgbterm.FgString(word, r, g, b)
	fmt.Println("prefered color: ", coloredWord, " (", r, ", ", g, ", ", b, ")", " score: ", y)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// predict()
	guess()
}
