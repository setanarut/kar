package res

import (
	"time"

	"github.com/setanarut/cm"
)

const TimerTick time.Duration = time.Second / 60
const DeltaTime float64 = 1.0 / 60.0

// Collision Bitmask Category
const (
	BitmaskPlayer        uint = 1 << 0
	BitmaskBlock         uint = 1 << 1
	BitmaskDropItem      uint = 1 << 2
	BitmaskPlayerRaycast uint = 1 << 4
	BitmaskEnemy         uint = 1 << 5
	BitmaskBomb          uint = 1 << 6
	BitmaskBombRaycast   uint = 1 << 7
	BitmaskGrabable      uint = 1 << 31
)

// Collision type
const (
	CollPlayer cm.CollisionType = iota
	CollEnemy
	CollBlock
	CollDropItem
	CollBomb
)
