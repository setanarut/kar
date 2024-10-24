package kar

import (
	"math"

	"github.com/setanarut/vec"
)

var (
	ScreenSize  = vec.Vec2{256, 240}
	WorldSize   = vec.Vec2{256, 256}
	BlockSize   = 16.0
	ChunkSize   = vec.Vec2{10, 10} // {9, 5} {12, 9}
	GUIScale    = 1.0
	ScreenScale = 1.0
	Gravity     = 1300.0
	SpaceStep   = 1.0 / 60.0

	// CollisionBias determines how fast overlapping shapes are pushed apart.
	// Expressed as a fraction of the error remaining after each second.
	// Defaults to math.Pow(0.9, 60) meaning that Chipmunk fixes 10% of overlap
	// each frame at 60Hz.
	CollisionBias = math.Pow(0.9, 60)

	// CollisionSlop is amount of encouraged penetration between colliding shapes.
	// Used to reduce oscillating contacts and keep the collision cache warm.
	// Defaults to 0.1. If you have poor simulation quality, increase this number
	// as much as possible without allowing visible amounts of overlap
	CollisionSlop = 0.01

	// 0.9 would mean that each body's velocity will drop 10% per second.
	// The default value is 1.0,
	Damping         = 0.9
	Iterations uint = 30
)

var (
	UseSpatialHash           = false
	SpatialHashDim   float64 = 128.0
	SpatialHashCount int     = 256
)
