package types

import (
	"kar/engine/cm"
	"time"
)

const TimerTick = time.Second / 60

const (
	ItemSnowball ItemType = iota
	ItemBomb
	ItemKey
	ItemPotion
	ItemAxe
	ItemShovel
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
