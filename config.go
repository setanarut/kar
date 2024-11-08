package kar

import (
	"math"

	"github.com/setanarut/vec"
)

type vec2 = vec.Vec2

var (
	ScreenSize = vec2{960, 540}
	WorldSize  = vec2{256, 256}
	ChunkSize  = vec2{11, 8} // {9, 5} {12, 9}
	BlockSize  = 48.0
)

var (
	GUIScale    = 1.0
	ScreenScale = 1.0
	Gravity     = 800.
	SpaceStep   = 1.0 / 60.0

	// CollisionBias determines how fast overlapping shapes are pushed apart.
	// Expressed as a fraction of the error remaining after each second.
	// Defaults to math.Pow(0.9, 60) meaning that Chipmunk fixes 10% of overlap
	// each frame at 60Hz.
	CollisionBias = math.Pow(0.0000001, 60)

	// Çarpışma kayması çok küçük olduğunda, nesneler birbirine çok yakınlaşabilir
	// ve bu da sıkışma sorunlarına yol açabilir. Özellikle hızlı hareket eden
	// nesnelerde, çarpışma algılama başarısız olabilir ve nesneler geçiş yapamadıkları
	// noktalarda sıkışabilir.
	CollisionSlop = 0.0000000001

	// 0.9 would mean that each body's velocity will drop 10% per second.
	// The default value is 1.0,
	Damping         = 1.0
	Iterations uint = 30
)

var (
	UseSpatialHash           = false
	SpatialHashDim   float64 = 128.0
	SpatialHashCount int     = 256
)
