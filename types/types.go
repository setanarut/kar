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

type DataAI struct {
	Target     *donburi.Entry
	Follow     bool
	FollowDist float64
}

type DrawOptions struct {
	CenterOffset vec.Vec2
	Scale        vec.Vec2
	Rotation     float64
	FlipX        bool
	ScaleColor   color.Color
}

type Mobile struct {
	Speed, Accel float64
}

type Timer struct {
	Duration time.Duration
	Elapsed  time.Duration
}

type Index struct {
	Index int
}

type Inventory struct {
	Slots    [9]ItemStack
	HandSlot ItemStack
}

type Item struct {
	ID    uint16
	Chunk image.Point // to mark which chunk the block belongs to
}

type ItemStack struct {
	ID       uint16
	Quantity uint8
}
type Health struct {
	Health    float64
	MaxHealth float64
}
