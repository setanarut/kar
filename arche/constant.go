package arche

import (
	"kar/engine/cm"
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
	CollisionTypePlayer cm.CollisionType = iota
	CollisionTypeEnemy
	CollisionTypeWall
	CollisionTypeSnowball
	CollisionTypeBomb
	CollisionTypeCollectible
	CollisionTypeDoor
)

var FilterBombRaycast cm.ShapeFilter = cm.NewShapeFilter(0, BitmaskBombRaycast, cm.AllCategories&^BitmaskBomb)

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
