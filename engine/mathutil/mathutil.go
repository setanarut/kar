package mathutil

import (
	"image"
	"math"
	"math/rand/v2"

	"github.com/setanarut/cm"

	"github.com/setanarut/vec"
)

func MapRange(v, a, b, c, d float64) float64 {
	return (v-a)/(b-a)*(d-c) + c
}

// InRange returns true if v is in range [min..max], else false
func InRange(v, min, max float64) bool {
	return v >= min && v <= max
}

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func Degrees(radians float64) float64 {
	return radians * 180.0 / math.Pi
}

// InvertAngle invert angle
func InvertAngle(angle float64) float64 {
	return -angle
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
func Linspace(start, stop float64, num int) (res []float64) {
	if num <= 0 {
		return []float64{}
	}
	if num == 1 {
		return []float64{start}
	}
	step := (stop - start) / float64(num-1)
	res = make([]float64, num)
	res[0] = start
	for i := 1; i < num; i++ {
		res[i] = start + float64(i)*step
	}
	res[num-1] = stop
	return
}

// RotateAbout rotates point about origin
func RotateAbout(angle float64, point, origin vec.Vec2) vec.Vec2 {
	b := vec.Vec2{}
	b.X = math.Cos(angle)*(point.X-origin.X) - math.Sin(angle)*(point.Y-origin.Y) + origin.X
	b.Y = math.Sin(angle)*(point.X-origin.X) + math.Cos(angle)*(point.Y-origin.Y) + origin.Y
	return b
}

// PointOnCircle returns point at angle
func PointOnCircle(center vec.Vec2, radius float64, angle float64) vec.Vec2 {
	x := center.X + (radius * math.Cos(angle))
	y := center.Y + (radius * math.Sin(angle))
	return vec.Vec2{x, y}
}

func RandomPoint(minX, maxX, minY, maxY float64) vec.Vec2 {
	return vec.Vec2{X: minX + rand.Float64()*(maxX-minX), Y: minY + rand.Float64()*(maxY-minY)}
}

func RandomPointInBB(bb cm.BB, margin float64) vec.Vec2 {
	return RandomPoint(bb.L+margin, bb.R-margin, bb.T-margin, bb.B+margin)
}

// IsMoving is velocity vector moving?
func IsMoving(velocityVector vec.Vec2, minSpeed float64) bool {
	if math.Abs(velocityVector.X) < minSpeed && math.Abs(velocityVector.Y) < minSpeed {
		return true
	} else {
		return false
	}
}

func RectangleScaleFactor(W, H, targetW, targetH float64) vec.Vec2 {
	return vec.Vec2{(targetW / W), (targetH / H)}
}

func CircleScaleFactor(radius float64, imageWidth int) vec.Vec2 {
	scaleX := 2 * radius / float64(imageWidth)
	return vec.Vec2{scaleX, scaleX}
}

func DistanceSq(a, b image.Point) float64 {
	return vec.FromPoint(a).DistanceSq(vec.FromPoint(b))
}
