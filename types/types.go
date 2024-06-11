package types

import (
	"image"
	"image/color"
	"kar/engine/cm"
	"time"

	"github.com/yohamta/donburi"
)

type System interface {
	Init()
	Update()
	Draw()
}

type ItemType int
type BlockType int

type DataAI struct {
	Target         *donburi.Entry
	Follow         bool
	FollowDistance float64
}

type DataDrawOptions struct {
	CenterOffset cm.Vec2
	Scale        cm.Vec2
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
	BlockType  BlockType
}
type DataHealth struct {
	Health    float64
	MaxHealth float64
}
