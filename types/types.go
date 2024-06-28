package types

import (
	"image"
	"image/color"
	"kar/engine/vec"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type ISystem interface {
	Init()
	Update()
	Draw(screen *ebiten.Image)
}

type ItemType uint8

type BlockSpawnFunc func(pos vec.Vec2, chunkCoord image.Point, blockType ItemType)

type DataAI struct {
	Target         *donburi.Entry
	Follow         bool
	FollowDistance float64
}

type DataDrawOptions struct {
	CenterOffset vec.Vec2
	Scale        vec.Vec2
	Rotation     float64
	FlipX        bool
	ScaleColor   color.Color
}

type DataMobile struct {
	Speed, Accel float64
}

type DataTimer struct {
	TimerDuration time.Duration
	Elapsed       time.Duration
}
type DataInventory struct {
	CurrentItem *donburi.Entry
	Items       map[ItemType]int
}

type DataItem struct {
	// to mark which chunk the block belongs to
	ChunkCoord image.Point
	Item       ItemType
}
type DataHealth struct {
	Health    float64
	MaxHealth float64
}
