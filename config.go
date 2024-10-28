package kar

import "github.com/setanarut/vec"

var (
	ScreenSize = vec.Vec2{960, 540}
	WorldSize  = vec.Vec2{256, 256}
	ChunkSize  = vec.Vec2{11, 8} // {9, 5} {12, 9}
	BlockSize  = 80.0
)

var BlockCenterOffset = vec.Vec2{BlockSize / 2.0, BlockSize / 2.0}.Neg()
