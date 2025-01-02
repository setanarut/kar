package arc

import (
	"strconv"
)

type Rect struct {
	X, Y, W, H float64
}

// Overlaps checks if the rectangle overlaps with another rectangle
func (r *Rect) Overlaps(x, y, w, h float64) bool {
	return r.X+r.W > x && x+w > r.X && r.Y+r.H > y && y+h > r.Y

}

// Overlaps checks if the rectangle overlaps with another rectangle
func (r *Rect) OverlapsRect(box *Rect) bool {
	return r.X+r.W > box.X && box.X+box.W > r.X && r.Y+r.H > box.Y && box.Y+box.H > r.Y
}

func (r *Rect) String() string {
	x := "X: " + strconv.FormatFloat(r.X, 'f', -1, 64)
	y := "Y: " + strconv.FormatFloat(r.X, 'f', -1, 64)
	return x + y
}

type DrawOptions struct {
	Scale float64
	FlipX bool
}

type ItemTimers struct {
	CollisionCountdown int
	AnimationIndex     int
}

type AnimationFrameIndex struct {
	Index int
}

type ItemID struct {
	ID uint16
}

type Durability struct {
	Durability int
}

type Health struct {
	Health    int
	MaxHealth int
}
