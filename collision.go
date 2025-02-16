package kar

import (
	"math"
)

const EPSILON = 1e-8

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func sign(value float64) float64 {
	if value < 0 {
		return -1
	}
	return 1
}

type AABB struct {
	Pos  Vec
	Half Vec
}

type Hit struct {
	Pos    Vec
	Delta  Vec
	Normal Vec
	Time   float64
}

func (a AABB) Segment(pos, delta, padding Vec, hit *Hit) bool {
	scaleX := 1.0 / delta.X
	scaleY := 1.0 / delta.Y
	signX := sign(scaleX)
	signY := sign(scaleY)
	nearTimeX := (a.Pos.X - signX*(a.Half.X+padding.X) - pos.X) * scaleX
	nearTimeY := (a.Pos.Y - signY*(a.Half.Y+padding.Y) - pos.Y) * scaleY
	farTimeX := (a.Pos.X + signX*(a.Half.X+padding.X) - pos.X) * scaleX
	farTimeY := (a.Pos.Y + signY*(a.Half.Y+padding.Y) - pos.Y) * scaleY
	if math.IsNaN(nearTimeY) {
		nearTimeY = math.Inf(1)
	}
	if math.IsNaN(farTimeY) {
		farTimeY = math.Inf(1)
	}
	if nearTimeX > farTimeY || nearTimeY > farTimeX {
		return false
	}
	nearTime := math.Max(nearTimeX, nearTimeY)
	farTime := math.Min(farTimeX, farTimeY)
	if nearTime >= 1 || farTime <= 0 {
		return false
	}
	if hit == nil {
		return true
	}
	hit.Time = clamp(nearTime, 0, 1)
	if nearTimeX > nearTimeY {
		hit.Normal.X = -signX
		hit.Normal.Y = 0
	} else {
		hit.Normal.X = 0
		hit.Normal.Y = -signY
	}
	hit.Delta.X = (1.0 - hit.Time) * -delta.X
	hit.Delta.Y = (1.0 - hit.Time) * -delta.Y
	hit.Pos.X = pos.X + delta.X*hit.Time
	hit.Pos.Y = pos.Y + delta.Y*hit.Time
	return true
}

func (a AABB) Overlap(a2 AABB, hit *Hit) bool {
	dx := a2.Pos.X - a.Pos.X
	px := a2.Half.X + a.Half.X - math.Abs(dx)
	if px <= 0 {
		return false
	}

	dy := a2.Pos.Y - a.Pos.Y
	py := a2.Half.Y + a.Half.Y - math.Abs(dy)
	if py <= 0 {
		return false
	}

	if hit == nil {
		return true
	}

	// hit.Collider = box1
	hit.Delta = Vec{}
	hit.Normal = Vec{}
	hit.Time = 0 // boxes overlap
	if px < py {
		sx := sign(dx)
		hit.Delta.X = px * sx
		hit.Normal.X = sx
		hit.Pos.X = a.Pos.X + a.Half.X*sx
		hit.Pos.Y = a2.Pos.Y
	} else {
		sy := sign(dy)
		hit.Delta.Y = py * sy
		hit.Normal.Y = sy
		hit.Pos.X = a2.Pos.X
		hit.Pos.Y = a.Pos.Y + a.Half.Y*sy
	}
	return true
}

func (a AABB) OverlapSweep(a2 AABB, delta Vec, hit *Hit) bool {
	if delta.IsZero() {
		return a.Overlap(a2, hit)
	}
	result := a.Segment(a2.Pos, delta, a2.Half, hit)
	if result {
		// hit.Time = 1.0
		hit.Time = clamp(hit.Time-EPSILON, 0, 1)
		direction := delta.Unit()
		hit.Pos.X = clamp(
			hit.Pos.X+direction.X*a2.Half.X,
			a.Pos.X-a.Half.X,
			a.Pos.X+a.Half.X,
		)
		hit.Pos.Y = clamp(
			hit.Pos.Y+direction.Y*a2.Half.Y,
			a.Pos.Y-a.Half.Y,
			a.Pos.Y+a.Half.Y,
		)
	}
	return result
}
