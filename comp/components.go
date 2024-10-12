package comp

import (
	"kar/types"
	"time"

	"github.com/setanarut/anim"
	"github.com/setanarut/cm"

	"github.com/setanarut/vec"

	"github.com/yohamta/donburi"
)

// Components

var (
	Mobile = donburi.NewComponentType[types.DataMobile](types.DataMobile{
		Speed: 500,
		Accel: 80,
	})

	AnimPlayer  = donburi.NewComponentType[anim.AnimationPlayer]()
	Body        = donburi.NewComponentType[cm.Body]()
	DrawOptions = donburi.NewComponentType[types.DataDrawOptions](types.DataDrawOptions{Scale: vec.Vec2{1, 1}})
	Inventory   = donburi.NewComponentType[types.DataInventory]()
	Damage      = donburi.NewComponentType[float64](1.0)
	Health      = donburi.NewComponentType[types.DataHealth](types.DataHealth{Health: 10.0, MaxHealth: 10.0})
	Timer       = donburi.NewComponentType[types.DataTimer](types.DataTimer{TimerDuration: time.Second / 3})
	Index       = donburi.NewComponentType[types.DataIndex]()
	AI          = donburi.NewComponentType[types.DataAI](types.DataAI{Follow: false, FollowDistance: 300})
	Item        = donburi.NewComponentType[types.DataItem]()
)

// Tags
var (
	TagWASD     = donburi.NewTag()
	TagWASDFly  = donburi.NewTag()
	TagPlayer   = donburi.NewTag()
	TagEnemy    = donburi.NewTag()
	TagItem     = donburi.NewTag()
	TagBlock    = donburi.NewTag()
	TagDebugBox = donburi.NewTag()
)
