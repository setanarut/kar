package kar

import (
	"time"

	"github.com/setanarut/cm"
)

const TimerTick time.Duration = time.Second / 60
const DeltaTime float64 = 1.0 / 60.0

// Collision Bitmask Category
const (
	PlayerMask uint = 1 << iota
	BlockMask
	DropItemMask
	PlayerRayMask
	EnemyMask
	BombMask
	BombRaycastMask
)
const GrabableMask uint = 1 << 31
const AllMask = ^uint(0)

// Collision type
const (
	PlayerCT cm.CollisionType = iota
	EnemyCT
	BlockCT
	DropItemCT
	BombCT
)
