package resources

import (
	"time"

	"github.com/setanarut/cm"
)

const TimerTick time.Duration = time.Second / 60
const DeltaTime float64 = 1.0 / 60.0

// Collision Bitmask Category
const (
	BitmaskPlayer uint = 1 << iota
	BitmaskBlock
	BitmaskDropItem
	BitmaskPlayerRaycast
	BitmaskEnemy
	BitmaskBomb
	BitmaskBombRaycast
	BitmaskGrabable uint = 1 << 31
)

// Collision type
const (
	CollPlayer cm.CollisionType = iota
	CollEnemy
	CollBlock
	CollDropItem
	CollBomb
)
