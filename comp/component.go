package comp

import (
	"kar/types"

	"github.com/setanarut/anim"
	"github.com/setanarut/cm"

	"github.com/yohamta/donburi"
)

// Components
var (
	// Damage      = donburi.NewComponentType[float64](1.0)
	// AI          = donburi.NewComponentType[types.DataAI]()
	// Timer       = donburi.NewComponentType[types.Timer]()
	// ID        = donburi.NewComponentType[types.Item]()
	AnimPlayer     = donburi.NewComponentType[anim.AnimationPlayer]()
	Body           = donburi.NewComponentType[cm.Body]()
	CollisionTimer = donburi.NewComponentType[types.Timer]()
	DrawOptions    = donburi.NewComponentType[types.DrawOptions]()
	Health         = donburi.NewComponentType[types.Health]()
	Index          = donburi.NewComponentType[types.Index]()
	Inventory      = donburi.NewComponentType[types.Inventory]()
	Item           = donburi.NewComponentType[types.Item]()
	Mobile         = donburi.NewComponentType[types.Mobile]()
	StuckCountdown = donburi.NewComponentType[types.Countdown]()
)

// Tags
var (
	Block       = donburi.NewTag()
	Breakable   = donburi.NewTag()
	DebugBox    = donburi.NewTag()
	DropItem    = donburi.NewTag()
	Harvestable = donburi.NewTag()
	Player      = donburi.NewTag()
)
