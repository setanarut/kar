package kar

import (
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

func MapRange(v, a, b, c, d float64) float64 {
	return (v-a)/(b-a)*(d-c) + c
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
	for i := range n {
		// Normalize t to [0, 2π] range
		t := MapRange(float64(i), 0, float64(n-1), start, end)
		points[i] = amplitude * math.Sin(t)
	}
	return points
}

func ReadPNG(filePath string) image.Image {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return image
}

func DrawAABB(aabb *AABB) {
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
