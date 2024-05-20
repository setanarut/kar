package comp

import (
	"image/color"
	"kar/engine"
	"kar/engine/cm"
	"kar/model"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mazznoer/colorgrad"
	"github.com/yohamta/donburi"
)

var (
	Char = donburi.NewComponentType[model.CharacterData](model.CharacterData{
		Speed:               350,
		Accel:               80,
		Health:              100.,
		ShootCooldownTimer:  engine.NewTimer(time.Second / 5),
		SnowballPerCooldown: 1,
	})
	Render = donburi.NewComponentType[model.RenderData](model.RenderData{
		Offset:     cm.Vec2{},
		DrawScale:  cm.Vec2{1, 1},
		DrawAngle:  0.0,
		DIO:        &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear},
		ScaleColor: color.White,
	})

	Inventory = donburi.NewComponentType[model.InventoryData](model.InventoryData{
		Bombs:     100,
		Potion:    20,
		Snowballs: 5000,
		Keys:      make([]int, 0),
	})

	Effect = donburi.NewComponentType[model.EffectData](model.EffectData{
		ShootCooldown:    -(time.Second / 10),
		ExtraSnowball:    2,
		AddMovementSpeed: -200,
		EffectTimer:      engine.NewTimer(time.Second * 6),
	})

	Collectible = donburi.NewComponentType[model.CollectibleData]()

	Body = donburi.NewComponentType[cm.Body]()
	AI   = donburi.NewComponentType[model.AIData](model.AIData{Follow: true, FollowDistance: 300})

	Door     = donburi.NewComponentType[model.DoorData]()
	Damage   = donburi.NewComponentType[float64](25.0)
	Gradient = donburi.NewComponentType[colorgrad.Gradient](colorgrad.NewGradient().
			HtmlColors("rgb(0, 229, 255)", "rgb(93, 90, 193)", "rgb(255, 0, 123)").
			Domain(0, 100).
			Mode(colorgrad.BlendOklab).
			Interpolation(colorgrad.InterpolationBasis).
			Build())
)

// Tags
var (
	PlayerTag   = donburi.NewTag()
	WallTag     = donburi.NewTag()
	SnowballTag = donburi.NewTag()
	BombTag     = donburi.NewTag()
	EnemyTag    = donburi.NewTag()
)
