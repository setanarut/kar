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

type Health struct {
	Health    int
	MaxHealth int
}

type Rect struct {
	X, Y, W, H float64
}

// Overlaps checks if the rectangle overlaps with another rectangle
func (r *Rect) Overlaps(x, y, w, h float64) bool {
	return r.X+r.W > x && x+w > r.X && r.Y+r.H > y && y+h > r.Y

}

// Overlaps2 checks if the rectangle overlaps with another rectangle
func (r *Rect) Overlaps2(box *Rect) bool {
	return r.X+r.W > box.X && box.X+box.W > r.X && r.Y+r.H > box.Y && box.Y+box.H > r.Y
}
