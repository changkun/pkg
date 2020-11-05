package convert

import "math"

// Linear2sRGB converts linear inputs to sRGB space.
func Linear2sRGB(v float64) float64 {
	if v <= 0.0031308 {
		v *= 12.92
	} else {
		v = 1.055*math.Pow(v, 1/2.4) - 0.055
	}
	return v
}
