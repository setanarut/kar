package comp

import (
	"image/color"
	"kar/engine"
	"kar/engine/cm"
	"kar/res"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

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
	Items map[res.ItemType]int
}

// Components

var (
	Mobile = donburi.NewComponentType[DataMobile](DataMobile{
		Speed: 350,
		Accel: 80,
	})

	Render = donburi.NewComponentType[DataRender](DataRender{
		Offset:     cm.Vec2{},
		DrawScale:  cm.Vec2{1, 1},
		DrawAngle:  0.0,
		DIO:        &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear},
		ScaleColor: color.White,
	})

	Inventory = donburi.NewComponentType[DataInventory]()

	Damage      = donburi.NewComponentType[float64](1.0)
	Health      = donburi.NewComponentType[float64](8.0)
	Body        = donburi.NewComponentType[cm.Body]()
	AttackTimer = donburi.NewComponentType[DataTimer](DataTimer{TimerDuration: time.Second / 4})
	PoisonTimer = donburi.NewComponentType[DataTimer](DataTimer{TimerDuration: time.Second * 5})
	AI          = donburi.NewComponentType[DataAI](DataAI{Follow: false, FollowDistance: 300})
	Door        = donburi.NewComponentType[DataDoor]()
)

// Tags
var (
	WASDControll = donburi.NewTag()
	PlayerTag    = donburi.NewTag()
	WallTag      = donburi.NewTag()
	SnowballTag  = donburi.NewTag()
	BombTag      = donburi.NewTag()
	EnemyTag     = donburi.NewTag()
)
