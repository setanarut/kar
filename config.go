package kar

import "github.com/setanarut/vec"

type vec2 = vec.Vec2

var (
	ScreenSize = vec2{960, 540}
	WorldSize  = vec2{256, 256}
	ChunkSize  = vec2{11, 8} // {9, 5} {12, 9}
	BlockSize  = 80.0
)
