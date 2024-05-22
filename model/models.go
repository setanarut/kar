package model

import (
	"image/color"
	"kar/engine"
	"kar/engine/cm"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type ItemType int

type InventoryData struct {
	Snowballs int
	Bombs     int
	Potion    int
	Keys      []int
}

type CollectibleData struct {
	Type      ItemType
	ItemCount int
	KeyNumber int
}

type AIData struct {
	Follow         bool
	FollowDistance float64
}

type DoorData struct {
	LockNumber   int
	Open         bool
	PlayerHasKey bool
}
type RenderData struct {
	Offset     cm.Vec2
	DrawScale  cm.Vec2
	DrawAngle  float64
	AnimPlayer *engine.AnimationPlayer
	DIO        *ebiten.DrawImageOptions
	ScaleColor color.Color
}
type CharacterData struct {
	Speed, Accel, Health float64
	ShootCooldown        *TimerData
	SnowballPerCooldown  int
	CurrentTool          ItemType
}
type EffectData struct {
	AddMovementSpeed, Accel, Health float64
	ShootCooldown                   time.Duration
	ExtraSnowball                   int
	EffectTimer                     engine.Timer
}

type TimerData struct {
	Target  time.Duration
	Elapsed time.Duration
}
