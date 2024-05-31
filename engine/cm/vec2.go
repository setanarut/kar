package cm

import (
	"fmt"
	"image"
	"math"
)

type Vec2 struct {
	X, Y float64
}

// String returns string representation of this vector.
func (v Vec2) String() string {
	return fmt.Sprintf("Vector{X: %f, Y: %f}", v.X, v.Y)
}

// Equal checks if two vectors are equal. (Be careful when comparing floating point numbers!)
func (v Vec2) Equal(other Vec2) bool {
	return v.X == other.X && v.Y == other.Y
}

// Add two vector
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{v.X + other.X, v.Y + other.Y}
}

// Sub returns this - other
func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{v.X - other.X, v.Y - other.Y}
}

// Neg negates a vector.
func (v Vec2) Neg() Vec2 {
	return Vec2{-v.X, -v.Y}
}

// Mult scales vector
func (v Vec2) Mult(s float64) Vec2 {
	return Vec2{v.X * s, v.Y * s}
}

// Dot returns dot product
func (v Vec2) Dot(other Vec2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Cross calculates the 2D vector cross product analog.
// The cross product of 2D vectors results in a 3D vector with only a z component.
// This function returns the magnitude of the z value.
func (v Vec2) Cross(other Vec2) float64 {
	return v.X*other.Y - v.Y*other.X
}

// Perp returns a perpendicular vector. (90 degree rotation)
func (v Vec2) Perp() Vec2 {
	return Vec2{-v.Y, v.X}
}

// ReversePerp returns a perpendicular vector. (-90 degree rotation)
func (v Vec2) ReversePerp() Vec2 {
	return Vec2{v.Y, -v.X}
}

// Returns the vector projection onto other.
func (v Vec2) Project(other Vec2) Vec2 {
	return other.Mult(v.Dot(other) / other.Dot(other))
}

// ToAngle returns the angular direction v is pointing in (in radians).
func (v Vec2) ToAngle() float64 {
	return math.Atan2(v.Y, v.X)
}

// Rotate uses complex number multiplication to rotate this by other.
//
// Scaling will occur if this is not a unit vector.
func (v Vec2) Rotate(other Vec2) Vec2 {
	return Vec2{v.X*other.X - v.Y*other.Y, v.X*other.Y + v.Y*other.X}
}

// Unrotate is inverse of Vector.Rotate().
func (v Vec2) Unrotate(other Vec2) Vec2 {
	return Vec2{v.X*other.X + v.Y*other.Y, v.Y*other.X - v.X*other.Y}
}

// LengthSq returns the squared length of this vector.
//
// Faster than  Vector.Length() when you only need to compare lengths.
func (v Vec2) LengthSq() float64 {
	return v.Dot(v)
}

// Length returns the length of this vector
func (v Vec2) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

// Lerp linearly interpolates between this and other vector.
func (v Vec2) Lerp(other Vec2, t float64) Vec2 {
	return v.Mult(1.0 - t).Add(other.Mult(t))
}

// Normalize returns a normalized copy of this vector.
func (v Vec2) Normalize() Vec2 {
	// return v.Mult(1.0 / (v.Length() + math.SmallestNonzeroFloat64))
	return v.Mult(1.0 / (v.Length() + 1e-50))
}

// Spherical linearly interpolate between this and other.
func (v Vec2) LerpSpherical(other Vec2, t float64) Vec2 {
	dot := v.Normalize().Dot(other.Normalize())
	omega := math.Acos(Clamp(dot, -1, 1))

	if omega < 1e-3 {
		return v.Lerp(other, t)
	}

	denom := 1.0 / math.Sin(omega)
	return v.Mult(math.Sin((1.0-t)*omega) * denom).Add(other.Mult(math.Sin(t*omega) * denom))
}

// Spherical linearly interpolate between this towards other by no more than angle a radians.
func (v Vec2) LerpSphericalClamp(other Vec2, angle float64) Vec2 {
	dot := v.Normalize().Dot(other.Normalize())
	omega := math.Acos(Clamp(dot, -1, 1))
	return v.LerpSpherical(other, math.Min(angle, omega)/omega)
}

// ClampLenght clamps this vector lenght to len.
func (v Vec2) ClampLenght(len float64) Vec2 {
	if v.Dot(v) > len*len {
		return v.Normalize().Mult(len)
	}
	return Vec2{v.X, v.Y}
}

// LerpDistance linearly interpolates between this towards other by distance dist.
func (v Vec2) LerpDistance(other Vec2, dist float64) Vec2 {
	return v.Add(other.Sub(v).ClampLenght(dist))
}

// Returns the distance between this and other.
func (v Vec2) Distance(other Vec2) float64 {
	return v.Sub(other).Length()
}

// DistanceSq returns the squared distance between this and other.
//
// Faster than Vector.Distance() when you only need to compare distances.
func (v Vec2) DistanceSq(other Vec2) float64 {
	return v.Sub(other).LengthSq()
}

// Near returns true if the distance between this and other is less than dist.
func (v Vec2) Near(other Vec2, dist float64) bool {
	return v.DistanceSq(other) < dist*dist
}

// Collision related below

func (v Vec2) PointGreater(b, c Vec2) bool {
	return (b.Y-v.Y)*(v.X+b.X-2*c.X) > (b.X-v.X)*(v.Y+b.Y-2*c.Y)
}

func (v Vec2) CheckAxis(v1, p, n Vec2) bool {
	return p.Dot(n) <= math.Max(v.Dot(n), v1.Dot(n))
}

func (v Vec2) ClosestT(b Vec2) float64 {
	delta := b.Sub(v)
	return -Clamp(delta.Dot(v.Add(b))/delta.LengthSq(), -1.0, 1.0)
}

func (v Vec2) LerpT(b Vec2, t float64) Vec2 {
	ht := 0.5 * t
	return v.Mult(0.5 - ht).Add(b.Mult(0.5 + ht))
}

func (v Vec2) ClosestDist(v1 Vec2) float64 {
	return v.LerpT(v1, v.ClosestT(v1)).LengthSq()
}

func (v Vec2) ClosestPointOnSegment(a, b Vec2) Vec2 {
	delta := a.Sub(b)
	t := Clamp01(delta.Dot(v.Sub(b)) / delta.LengthSq())
	return b.Add(delta.Mult(t))
}

// Round returns the nearest integer Vector, rounding half away from zero.
func (v Vec2) Round() Vec2 {
	return Vec2{math.Round(v.X), math.Round(v.Y)}
}

// Point returns Vec2 as image.Point
func (v Vec2) Point() image.Point {
	return image.Point{int(v.X), int(v.Y)}
}

func Clamp(f, min, max float64) float64 {
	if f > min {
		return math.Min(f, max)
	} else {
		return math.Min(min, max)
	}
}

func Clamp01(f float64) float64 {
	return math.Max(0, math.Min(f, 1))
}

func Lerp(f1, f2, t float64) float64 {
	return f1*(1.0-t) + f2*t
}

func LerpConst(f1, f2, d float64) float64 {
	return f1 + Clamp(f2-f1, -d, d)
}

// ForAngle returns the unit length vector for the given angle (in radians).
func ForAngle(a float64) Vec2 {
	return Vec2{math.Cos(a), math.Sin(a)}
}
