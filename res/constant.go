package res

import (
	"kar/engine/cm"
	"kar/types"
	"time"
)

const TimerTick time.Duration = time.Second / 60
const DeltaTime float64 = 1.0 / 60.0

const (
	ItemSnowball types.ItemType = iota
	ItemBomb
	ItemKey
	ItemPotion
	ItemAxe
	ItemShovel
)
const (
	BlockAir types.BlockType = iota
	BlockDirt
	BlockStone
)

// Collision Bitmask Category
const (
	BitmaskPlayer        uint = 1
	BitmaskEnemy         uint = 2
	BitmaskBomb          uint = 4
	BitmaskSnowball      uint = 8
	BitmaskWall          uint = 16
	BitmaskDoor          uint = 32
	BitmaskCollectible   uint = 64
	BitmaskBombRaycast   uint = 128
	BitmaskPlayerRaycast uint = 256
)

// Collision type
const (
	CollPlayer cm.CollisionType = iota
	CollEnemy
	CollWall
	CollSnowball
	CollBomb
	CollCollectible
	CollDoor
)
