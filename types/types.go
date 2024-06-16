package types

import (
	"image"
	"image/color"
	"kar/engine/vec"
	"time"

	"github.com/yohamta/donburi"
)

type ISystem interface {
	Init()
	Update()
	Draw()
}

type ItemType int

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
	Items       []*donburi.Entry
}

type DataBlock struct {
	// to mark which chunk the block belongs to
	ChunkCoord image.Point
	BlockType  ItemType
}
type DataHealth struct {
	Health    float64
	MaxHealth float64
}
