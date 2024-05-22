package constants

import (
	"kar/engine/cm"
	"kar/model"
	"time"
)

const TimerTick = time.Second / 60

const (
	ItemSnowball model.ItemType = iota
	ItemBomb
	ItemKey
	ItemPotion
	ItemAxe
	ItemShovel
)

// Collision Bitmask Category
const (
	BitmaskPlayer      uint = 1
	BitmaskEnemy       uint = 2
	BitmaskBomb        uint = 4
	BitmaskSnowball    uint = 8
	BitmaskWall        uint = 16
	BitmaskDoor        uint = 32
	BitmaskCollectible uint = 64
	BitmaskBombRaycast uint = 128
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

/* var (
	QueryEnemy  = donburi.NewQuery(filter.Contains(component.EnemyTagComp))
	QueryPlayer = donburi.NewQuery(filter.Contains(component.PlayerTagComp))
	QueryDoor   = donburi.NewQuery(filter.Contains(component.DoorComp))
	QuerySnowball   = donburi.NewQuery(filter.Contains(component.SnowballTagComp))
	QueryBomb   = donburi.NewQuery(filter.Contains(component.BombTagComp))
	QueryAI     = donburi.NewQuery(filter.Contains(component.AIComp))
	QueryCamera = donburi.NewQuery(filter.Contains(component.CameraComp))
)
*/
