// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package bo_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/changkun/gobase/bo"
	"gonum.org/v1/gonum/floats"
)

func f(x, y float64) float64 {
	return math.Cos(x/2)/2 + math.Sin(y/4)
}

func gpAdd(gp *bo.GP, x, y float64) {
	gp.Add([]float64{x, y}, f(x, y))
}

func TestKnown(t *testing.T) {
	gp := bo.NewGP(bo.MaternCov{}, 0)

	gpAdd(gp, 0.25, 0.75)

	for i := 0; i < 20; i++ {
		gpAdd(gp, rand.Float64()*2*math.Pi-math.Pi, rand.Float64()*2*math.Pi-math.Pi)
	}

	if _, err := saveAll(gp); err != nil {
		t.Fatal(err)
	}
	mean, variance, err := gp.Estimate([]float64{0.25, 0.75})
	if err != nil {
		t.Fatal(err)
	}
	if !floats.EqualWithinAbs(mean, f(0.25, 0.75), 0.0001) {
		t.Fatalf("got mean = %f; not 1", mean)
	}
	if !floats.EqualWithinAbs(variance, 0, 0.0001) {
		t.Fatalf("got variance = %f; not 0", variance)
	}
}

func TestMaternCov(t *testing.T) {
	cases := []struct {
		a, b []float64
		want float64
	}{
		{
			[]float64{0},
			[]float64{0},
			1,
		},
		{
			[]float64{0, 1, 3},
			[]float64{0, 1, 2},
			0.828649,
		},
		{
			[]float64{0, 1, 4},
			[]float64{0, 1, 2},
			0.523994,
		},
	}
	for i, c := range cases {
		out := bo.MaternCov{}.Cov(c.a, c.b)
		if math.Abs(out-c.want) > 0.00001 {
			t.Errorf("%d. MaternCov(%+v, %+v) = %f; not %f", i, c.a, c.b, out, c.want)
		}
	}
}
