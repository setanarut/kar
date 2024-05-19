package comp

import (
	"image/color"
	"kar/engine"
	"kar/engine/cm"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mazznoer/colorgrad"
	"github.com/yohamta/donburi"
)

type ItemType int

const (
	Snowball ItemType = iota
	Bomb
	Key
	PowerUpItem
)

type InventoryData struct {
	Snowballs int
	Bombs     int
	PowerUp   int
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
	ShootCooldownTimer   engine.Timer
	SnowballPerCooldown  int
}
type EffectData struct {
	AddMovementSpeed, Accel, Health float64
	ShootCooldown                   time.Duration
	ExtraSnowball                   int
	EffectTimer                     engine.Timer
}

var Inventory = donburi.NewComponentType[InventoryData](InventoryData{
	Bombs:     100,
	PowerUp:   20,
	Snowballs: 5000,
	Keys:      make([]int, 0),
})
var Door = donburi.NewComponentType[DoorData]()

var Effect = donburi.NewComponentType[EffectData](EffectData{
	ShootCooldown:    -(time.Second / 10),
	ExtraSnowball:    2,
	AddMovementSpeed: -200,
	EffectTimer:      engine.NewTimer(time.Second * 6),
})

var Collectible = donburi.NewComponentType[CollectibleData]()

var Render = donburi.NewComponentType[RenderData](RenderData{
	Offset:     cm.Vec2{},
	DrawScale:  cm.Vec2{1, 1},
	DrawAngle:  0.0,
	DIO:        &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear},
	ScaleColor: color.White,
})

var Gradient = donburi.NewComponentType[colorgrad.Gradient](colorgrad.NewGradient().
	HtmlColors("rgb(0, 229, 255)", "rgb(93, 90, 193)", "rgb(255, 0, 123)").
	Domain(0, 100).
	Mode(colorgrad.BlendOklab).
	Interpolation(colorgrad.InterpolationBasis).
	Build())

var Body = donburi.NewComponentType[cm.Body]()
var AI = donburi.NewComponentType[AIData](AIData{Follow: true, FollowDistance: 300})

var Char = donburi.NewComponentType[CharacterData](CharacterData{
	Speed:               350,
	Accel:               80,
	Health:              100.,
	ShootCooldownTimer:  engine.NewTimer(time.Second / 5),
	SnowballPerCooldown: 1,
})

var Damage = donburi.NewComponentType[float64](25.0)

// Tags
var PlayerTag = donburi.NewTag()
var WallTag = donburi.NewTag()
var SnowballTag = donburi.NewTag()
var BombTag = donburi.NewTag()
var EnemyTag = donburi.NewTag()
