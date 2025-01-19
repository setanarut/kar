package arc

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

	// İlk pozisyonları belirle
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
		// X çarpışması varsa Y hareketini koru ama çarpışmayı X olarak işaretle
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

	return info
}

type CollisionInfo struct {
	Normal   [2]int
	DeltaX   float64
	DeltaY   float64
	Collided bool
}
