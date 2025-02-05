package system

import (
	"kar/arc"
	"math"
)

// Overlaps checks if the rectangle overlaps with another rectangle
func Overlaps(p1 *arc.Position, s1 *arc.Size, p2 *arc.Position, s2 *arc.Size) bool {
	return p1.X+s1.W > p2.X && p2.X+s2.W > p1.X && p1.Y+s1.H > p2.Y && p2.Y+s2.H > p1.Y
}

func CheckCollision(p1 *arc.Position, s1 *arc.Size, v1 *arc.Velocity, p2 *arc.Position, s2 *arc.Size) CollisionInfo {
	info := CollisionInfo{
		Normal: [2]int{0, 0},
	}

	// Önce statik çarpışma kontrolü
	if Overlaps(p1, s1, p2, s2) {
		info.Collided = true

		// Statik çarpışmada itme yönünü belirle
		centerX1 := p1.X + s1.W/2
		centerX2 := p2.X + s2.W/2
		centerY1 := p1.Y + s1.H/2
		centerY2 := p2.Y + s2.H/2

		diffX := centerX1 - centerX2
		diffY := centerY1 - centerY2

		// Hangi eksende daha az örtüşme var ise o yönde it
		overlapX := s1.W/2 + s2.W/2 - math.Abs(diffX)
		overlapY := s1.H/2 + s2.H/2 - math.Abs(diffY)

		if overlapX < overlapY {
			if diffX > 0 {
				info.Normal[0] = 1
				info.DeltaX = overlapX
			} else {
				info.Normal[0] = -1
				info.DeltaX = -overlapX
			}
		} else {
			if diffY > 0 {
				info.Normal[1] = 1
				info.DeltaY = overlapY
			} else {
				info.Normal[1] = -1
				info.DeltaY = -overlapY
			}
		}

		return info
	}

	// Hareket varsa dinamik çarpışma kontrolü
	if math.Abs(v1.X) > 0 || math.Abs(v1.Y) > 0 {
		nextPos := &arc.Position{p1.X + v1.X, p1.Y + v1.Y}
		// Önce X ekseninde hareket et ve kontrol et
		if Overlaps(p2, s2, nextPos, s1) {
			info.Collided = true
			if v1.X > 0 {
				info.DeltaX = p2.X - (nextPos.X + s1.W)
				info.Normal[0] = -1
			} else if v1.X < 0 {
				info.DeltaX = (p2.X + s2.W) - nextPos.X
				info.Normal[0] = 1
			}
			return info
		}

		// Sonra Y ekseninde hareket et ve kontrol et
		if Overlaps(p2, s2, nextPos, s1) {
			info.Collided = true
			if v1.Y > 0 {
				info.DeltaY = p2.Y - (nextPos.Y + s1.H)
				info.Normal[1] = -1
			} else if v1.Y < 0 {
				info.DeltaY = (p2.Y + s2.H) - nextPos.Y
				info.Normal[1] = 1
			}
		}
	}
	return info
}

type CollisionInfo struct {
	Normal   [2]int
	DeltaX   float64
	DeltaY   float64
	Collided bool
}
