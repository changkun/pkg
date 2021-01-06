// Copyright 2021 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package convert

import "math"

// Linear2sRGB converts linear inputs to sRGB space.
func Linear2sRGB(v float64) float64 {
	i := v * lutSize
	ifloor := int(i) & (lutSize - 1)
	v0 := lin2sRGBLUT[ifloor]
	v1 := lin2sRGBLUT[ifloor+1]
	i -= float64(ifloor)
	return v0*(1.0-i) + v1*i
}

func linear2sRGB(v float64) float64 {
	if v <= 0.0031308 {
		v *= 12.92
	} else {
		v = 1.055*math.Pow(v, 1/2.4) - 0.055
	}
	return v
}

const lutSize = 1024 // keep a power of 2

var lin2sRGBLUT [lutSize + 1]float64

func init() {
	for i := range lin2sRGBLUT[:lutSize] {
		lin2sRGBLUT[i] = linear2sRGB(float64(i) / lutSize)
	}
	// add an extra entry to avoid an "if (idxFloor == lutSize-1)"
	lin2sRGBLUT[lutSize] = lin2sRGBLUT[lutSize-1]
}
