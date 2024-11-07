package arc

import (
	"image"
	"time"

	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

type CmBody struct {
	Body *cm.Body
}

type CollisionActivationCountdown struct {
	Tick uint8
}

type SelfDestuctionCountdown struct {
	Tick uint8
}

type DrawOptions struct {
	CenterOffset vec.Vec2
	Scale        vec.Vec2
	FlipX        bool
}

type Mobile struct {
	Speed, Accel float64
}

type Timer struct {
	Duration time.Duration
	Elapsed  time.Duration
}

type AnimationFrameIndex struct {
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

func NewInventory() *Inventory {
	inv := &Inventory{}
	inv.HandSlot = ItemStack{}
	for i := range inv.Slots {
		inv.Slots[i] = ItemStack{}
	}
	return inv
}
