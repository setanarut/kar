package arc

import "math"

type ItemID struct {
	ID uint8
}

type Velocity struct {
	X float64
	Y float64
}

type ItemTimers struct {
	CollisionCountdown int
	AnimationIndex     int
}

type Durability struct {
	Durability int
}

type Mode struct {
	Mode int
}

type Health struct {
	Current int
	Max     int
}

type Rect struct {
	X, Y, W, H float64
}

// Overlaps checks if the rectangle overlaps with another rectangle
func (r *Rect) Overlaps(r2 *Rect) bool {
	return r.X+r.W > r2.X && r2.X+r2.W > r.X && r.Y+r.H > r2.Y && r2.Y+r2.H > r.Y
}

// Overlaps2 checks if the rectangle overlaps with another rectangle
func (r *Rect) Overlaps2(x, y, w, h float64) bool {
	return r.X+r.W > x && x+w > r.X && r.Y+r.H > y && y+h > r.Y

}

func (r *Rect) CheckCollision(r2 *Rect, velX, velY float64) CollisionInfo {
	info := CollisionInfo{
		Normal: [2]int{0, 0},
	}

	// Önce statik çarpışma kontrolü
	if r2.Overlaps2(r.X, r.Y, r.W, r.H) {
		info.Collided = true

		// Statik çarpışmada push direction'ı belirle
		centerX1 := r.X + r.W/2
		centerX2 := r2.X + r2.W/2
		centerY1 := r.Y + r.H/2
		centerY2 := r2.Y + r2.H/2

		diffX := centerX1 - centerX2
		diffY := centerY1 - centerY2

		// Hangi eksende daha az örtüşme var ise o yönde it
		overlapX := r.W/2 + r2.W/2 - math.Abs(diffX)
		overlapY := r.H/2 + r2.H/2 - math.Abs(diffY)

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
	if math.Abs(velX) > 0 || math.Abs(velY) > 0 {
		nextX := r.X + velX
		nextY := r.Y + velY

		// Önce X ekseninde hareket et ve kontrol et
		if r2.Overlaps2(nextX, r.Y, r.W, r.H) {
			info.Collided = true
			if velX > 0 {
				info.DeltaX = r2.X - (nextX + r.W)
				info.Normal[0] = -1
			} else if velX < 0 {
				info.DeltaX = (r2.X + r2.W) - nextX
				info.Normal[0] = 1
			}
			return info
		}

		// Sonra Y ekseninde hareket et ve kontrol et
		if r2.Overlaps2(nextX, nextY, r.W, r.H) {
			info.Collided = true
			if velY > 0 {
				info.DeltaY = r2.Y - (nextY + r.H)
				info.Normal[1] = -1
			} else if velY < 0 {
				info.DeltaY = (r2.Y + r2.H) - nextY
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
