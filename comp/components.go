package comp

import (
	"image/color"
	"kar/engine/cm"
	"kar/types"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// Components

var (
	Mobile = donburi.NewComponentType[types.DataMobile](types.DataMobile{
		Speed: 500,
		Accel: 80,
	})

	Render = donburi.NewComponentType[types.DataRender](types.DataRender{
		Offset:     cm.Vec2{},
		DrawScale:  cm.Vec2{1, 1},
		DrawAngle:  0.0,
		DIO:        &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest},
		ScaleColor: color.White,
	})

	Sprite = donburi.NewComponentType[types.DataSprite](types.DataSprite{
		Offset:     cm.Vec2{},
		DrawScale:  cm.Vec2{1, 1},
		DrawAngle:  0.0,
		DIO:        &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest},
		ScaleColor: color.White,
	})

	Inventory = donburi.NewComponentType[types.DataInventory]()

	Damage      = donburi.NewComponentType[float64](1.0)
	Health      = donburi.NewComponentType[types.DataHealth](types.DataHealth{Health: 10.0, MaxHealth: 10.0})
	Body        = donburi.NewComponentType[cm.Body]()
	AttackTimer = donburi.NewComponentType[types.DataTimer](types.DataTimer{TimerDuration: time.Second / 3})
	AI          = donburi.NewComponentType[types.DataAI](types.DataAI{Follow: false, FollowDistance: 300})
	Block       = donburi.NewComponentType[types.DataBlock]()
)

// Tags
var (
	WASDTag     = donburi.NewTag()
	PlayerTag   = donburi.NewTag().SetName("PlayerTag")
	SnowballTag = donburi.NewTag()
	BombTag     = donburi.NewTag()
	EnemyTag    = donburi.NewTag()
)
