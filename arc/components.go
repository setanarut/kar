package arc

import (
	"strconv"
	"time"
)

type Rect struct {
	X, Y, W, H float64
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

type Timer struct {
	Duration time.Duration
	Elapsed  time.Duration
}

type AnimationFrameIndex struct {
	Index int
}

type Item struct {
	ID uint16
}

type Health struct {
	Health    float64
	MaxHealth float64
}
