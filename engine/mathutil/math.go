package mathutil

import (
	"math"
	"math/rand/v2"

	"github.com/setanarut/cm"
	"golang.org/x/exp/constraints"

	"github.com/setanarut/vec"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func MapRange[T Number](v, a, b, c, d T) T {
	return (v-a)/(b-a)*(d-c) + c
}

// InRange checks if a number is within a given range
func InRange[T constraints.Ordered](num, min, max T) bool {
	return num >= min && num <= max
}

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func Degrees(radians float64) float64 {
	return radians * 180.0 / math.Pi
}

func RandRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)

}

func RandRangeInt(min, max int) int {
	return rand.IntN(max-min+1) + min
}

// Clamp returns f clamped to [low, high]
func Clamp(f, low, high float64) float64 {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}

// Linspace returns evenly spaced numbers over a specified closed interval.
func Linspace(start, stop float64, num int) (resources []float64) {
	if num <= 0 {
		return []float64{}
	}
	if num == 1 {
		return []float64{start}
	}
	step := (stop - start) / float64(num-1)
	resources = make([]float64, num)
	resources[0] = start
	for i := 1; i < num; i++ {
		resources[i] = start + float64(i)*step
	}
	resources[num-1] = stop
	return
}

// SinSpace returns n points between start and end based on a sinusoidal function with a given amplitude
//
//	start := 0.0       // Start of range
//	end := 2 * math.Pi // End of range (one full sine wave)
//	amplitude := 2.0   // Amplitude of the sine wave
//	n := 10            // Number of points
func SinSpace(start, end, amplitude float64, n int) []float64 {
	points := make([]float64, n)
	step := (end - start) / float64(n-1)

	for i := 0; i < n; i++ {
		// Normalize the step to be between 0 and 2*Pi for one full sinusoidal wave
		t := start + step*float64(i)
		// Apply the amplitude to the sine wave
		points[i] = amplitude * math.Sin(t)
	}

	return points
}

func RandomPoint(minX, maxX, minY, maxY float64) vec.Vec2 {
	return vec.Vec2{X: minX + rand.Float64()*(maxX-minX), Y: minY + rand.Float64()*(maxY-minY)}
}

func RandomPointInBB(bb cm.BB, margin float64) vec.Vec2 {
	return RandomPoint(bb.L+margin, bb.R-margin, bb.T-margin, bb.B+margin)
}

func GetRectScaleFactor(W, H, targetW, targetH float64) vec.Vec2 {
	return vec.Vec2{(targetW / W), (targetH / H)}
}

func GetCircleScaleFactor(radius float64, imageWidth int) vec.Vec2 {
	scaleX := 2 * radius / float64(imageWidth)
	return vec.Vec2{scaleX, scaleX}
}
