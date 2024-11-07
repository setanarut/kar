package arc

import "github.com/setanarut/cm"

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
	PlayerCollisionFilter   = cm.ShapeFilter{0, PlayerBit, AllBits &^ PlayerRayBit}
	BlockCollisionFilter    = cm.ShapeFilter{0, BlockBit, cm.AllCategories}
	DropItemCollisionFilter = cm.ShapeFilter{
		Group:      2,
		Categories: DropItemBit,
		Mask:       AllBits &^ PlayerRayBit &^ PlayerBit,
	}
)
