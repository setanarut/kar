package arc

import (
	"image"
	"image/color"
	"time"

	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

// Collision Bitmask Category
const (
	PlayerBit uint = 1 << iota
	BlockBit
	DropItemBit
	PlayerRayBit
	EnemyBit
	BombBit
	BombRaycastBit
	GrabableBit uint = 1 << 31
	AllBits          = ^uint(0)
)

// Collision type
const (
	Player cm.CollisionType = iota
	Enemy
	Block
	DropItem
	Bomb
)

var (
	DropItemCollisionFilter = cm.ShapeFilter{
		Group:      2,
		Categories: DropItemBit,
		Mask:       AllBits &^ PlayerRayBit &^ PlayerBit,
	}
	PlayerCollisionFilter = cm.ShapeFilter{0, PlayerBit, AllBits &^ PlayerRayBit}
	BlockCollisionFilter  = cm.ShapeFilter{0, BlockBit, cm.AllCategories}
)

type CollisionTimer Timer

type Countdown struct {
	Duration uint8
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
