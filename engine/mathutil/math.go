package mathutil

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func MapRange[T Number](v, a, b, c, d T) T {
	return (v-a)/(b-a)*(d-c) + c
}

// degrees to radians
func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// radians to degrees
func Degrees(radians float64) float64 {
	return radians * 180.0 / math.Pi
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
