package arc

import "math"

type ItemID struct {
	ID uint16
}

type Velocity struct {
	VelX float64
	VelY float64
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
	Health    int
	MaxHealth int
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

	// Hareket sonrası pozisyonu kontrol et
	nextX := r.X + velX
	nextY := r.Y + velY

	if !r2.Overlaps2(nextX, nextY, r.W, r.H) {
		return info
	}

	// Çarpışma var
	info.Collided = true

	// X ekseni çarpışması
	if velX > 0 {
		info.DeltaX = r2.X - (nextX + r.W)
		info.Normal[0] = -1
	} else if velX < 0 {
		info.DeltaX = (r2.X + r2.W) - nextX
		info.Normal[0] = 1
	}

	// Y ekseni çarpışması
	if velY > 0 {
		info.DeltaY = r2.Y - (nextY + r.H)
		info.Normal[1] = -1
	} else if velY < 0 {
		info.DeltaY = (r2.Y + r2.H) - nextY
		info.Normal[1] = 1
	}

	// Sadece bir eksende çarpışma olsun
	if math.Abs(velX) > 0 && math.Abs(velY) > 0 {
		if math.Abs(info.DeltaX) < math.Abs(info.DeltaY) {
			info.DeltaY = 0
			info.Normal[1] = 0
		} else {
			info.DeltaX = 0
			info.Normal[0] = 0
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
