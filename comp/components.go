package comp

import (
	"image/color"
	"kar/engine/cm"
	"kar/models"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// Components

var (
	Mobile = donburi.NewComponentType[models.DataMobile](models.DataMobile{
		Speed: 350,
		Accel: 80,
	})

	Render = donburi.NewComponentType[models.DataRender](models.DataRender{
		Offset:     cm.Vec2{},
		DrawScale:  cm.Vec2{1, 1},
		DrawAngle:  0.0,
		DIO:        &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear},
		ScaleColor: color.White,
	})

	Inventory = donburi.NewComponentType[models.DataInventory]()

	Damage      = donburi.NewComponentType[float64](1.0)
	Health      = donburi.NewComponentType[float64](8.0)
	Body        = donburi.NewComponentType[cm.Body]()
	AttackTimer = donburi.NewComponentType[models.DataTimer](models.DataTimer{TimerDuration: time.Second / 4})
	PoisonTimer = donburi.NewComponentType[models.DataTimer](models.DataTimer{TimerDuration: time.Second * 5})
	AI          = donburi.NewComponentType[models.DataAI](models.DataAI{Follow: false, FollowDistance: 300})
	Door        = donburi.NewComponentType[models.DataDoor]()
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
