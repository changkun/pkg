package detect_test

import (
	"testing"

	"changkun.de/x/pkg/detect"
)

func TestMovingAverage(t *testing.T) {
	xs := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	t.Log(detect.MovingAverage(xs, 3))
	t.Log(detect.KolmogorovZurbenko(xs, 3, 3))
	t.Log(detect.AdaptiveKolmogorovZurbenko(xs, 3, 3))
}
