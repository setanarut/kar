package comp

import (
	"kar/types"
	"time"

	"github.com/setanarut/cm"

	"github.com/setanarut/anim"
	"github.com/setanarut/vec"

	"github.com/yohamta/donburi"
)

// Components

var (
	Mobile = donburi.NewComponentType[types.DataMobile](types.DataMobile{
		Speed: 500,
		Accel: 80,
	})

	AnimationPlayer = donburi.NewComponentType[anim.AnimationPlayer]()
	Body            = donburi.NewComponentType[cm.Body]()
	DrawOptions     = donburi.NewComponentType[types.DataDrawOptions](types.DataDrawOptions{Scale: vec.Vec2{1, 1}})
	Inventory       = donburi.NewComponentType[types.DataInventory]()
	Damage          = donburi.NewComponentType[float64](1.0)
	Health          = donburi.NewComponentType[types.DataHealth](types.DataHealth{Health: 10.0, MaxHealth: 10.0})
	AttackTimer     = donburi.NewComponentType[types.DataTimer](types.DataTimer{TimerDuration: time.Second / 3})
	AI              = donburi.NewComponentType[types.DataAI](types.DataAI{Follow: false, FollowDistance: 300})
	Item            = donburi.NewComponentType[types.DataItem]()
)

// Tags
var (
	WASDTag      = donburi.NewTag()
	WASDFlyTag   = donburi.NewTag()
	PlayerTag    = donburi.NewTag().SetName("PlayerTag")
	EnemyTag     = donburi.NewTag()
	DropItemTag  = donburi.NewTag()
	BlockItemTag = donburi.NewTag()
	DebugBoxTag  = donburi.NewTag()
)
