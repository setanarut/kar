package kar

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

func mapRange(v, a, b, c, d float64) float64 {
	return (v-a)/(b-a)*(d-c) + c
}

func linspace(min, max float64, n int) []float64 {
	if n == 1 {
		return []float64{min}
	}
	d := max - min
	l := float64(n) - 1
	res := make([]float64, n)
	for i := range res {
		res[i] = (min + (float64(i)*d)/l)
	}
	return res
}

// sinspace returns n points between start and end based on a sinusoidal function with a given amplitude
//
//	start := 0.0       // Start of range
//	end := 2 * math.Pi // End of range (one full sine wave)
//	amplitude := 2.0   // Amplitude of the sine wave
//	n := 10            // Number of points
func sinspace(start, end, amplitude float64, n int) []float64 {
	tValues := linspace(start, end, n)
	for i, t := range tValues {
		tValues[i] = amplitude * math.Sin(t)
	}
	return tValues
}

func drawAABB(aabb *AABB) {
	x, y := cameraRes.ApplyCameraTransformToPoint(aabb.Pos.X, aabb.Pos.Y)
	vector.DrawFilledRect(
		Screen,
		float32(x-aabb.Half.X),
		float32(y-aabb.Half.Y),
		float32(aabb.Half.X*2),
		float32(aabb.Half.Y*2),
		color.RGBA{128, 0, 0, 10},
		false,
	)
}

func formatDuration(d time.Duration) string {
	s := int(d.Seconds())
	return fmt.Sprintf("%02d:%02d:%02d:%02d", (s/3600)/24, (s/3600)%24, (s/60)%60, s%60)
}
