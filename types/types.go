package types

import (
	"image"
	"image/color"
	"kar/engine"
	"kar/engine/cm"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type System interface {
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

type DataRender struct {
	Offset     cm.Vec2
	DrawScale  cm.Vec2
	DrawAngle  float64
	AnimPlayer *engine.AnimationPlayer
	DIO        *ebiten.DrawImageOptions
	ScaleColor color.Color
}

type DataSprite struct {
	Image      *ebiten.Image
	Offset     cm.Vec2
	DrawScale  cm.Vec2
	ScaleColor color.Color
	DrawAngle  float64
	DIO        *ebiten.DrawImageOptions
}

type DataMobile struct {
	Speed, Accel float64
}

type DataTimer struct {
	TimerDuration time.Duration
	Elapsed       time.Duration
}
type DataInventory struct {
	Items map[ItemType]int
}

type DataBlock struct {
	// to mark which chunk the block belongs to
	ChunkCoord image.Point
}
