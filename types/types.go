package types

import (
	"image"
	"image/color"
	"time"

	"github.com/setanarut/vec"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type ISystem interface {
	Init()
	Update()
	Draw(screen *ebiten.Image)
}

type ItemID uint8

type BlockSpawnFunc func(pos vec.Vec2, chunkCoord image.Point, blockType ItemID)

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

type DataIndex struct {
	Index int
}

type DataInventory struct {
	Slots [40]*DataItemStack
}

type DataItem struct {
	Item ItemID
	// to mark which chunk the block belongs to
	ChunkCoord image.Point
}
type DataHealth struct {
	Health    float64
	MaxHealth float64
}

type DataItemStack struct {
	Item     ItemID
	Quantity uint8
}
