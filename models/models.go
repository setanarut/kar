package models

import (
	"image/color"
	"kar/engine"
	"kar/engine/cm"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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

type ItemType int

type DataAI struct {
	Follow         bool
	FollowDistance float64
}

type DataDoor struct {
	LockNumber   int
	Open         bool
	PlayerHasKey bool
}

type DataRender struct {
	Offset     cm.Vec2
	DrawScale  cm.Vec2
	DrawAngle  float64
	AnimPlayer *engine.AnimationPlayer
	DIO        *ebiten.DrawImageOptions
	ScaleColor color.Color
}
type DataMobile struct {
	Speed, Accel float64
}

type DataTimer struct {
	TimerDuration time.Duration
	Elapsed       time.Duration
}
type DataInventory struct {
	Items map[ItemType]int
}
