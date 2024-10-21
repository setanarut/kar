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
	Mobile         = donburi.NewComponentType[types.Mobile]()
	AnimPlayer     = donburi.NewComponentType[anim.AnimationPlayer]()
	Body           = donburi.NewComponentType[cm.Body]()
	DrawOptions    = donburi.NewComponentType[types.DrawOptions]()
	Inventory      = donburi.NewComponentType[types.Inventory]()
	Health         = donburi.NewComponentType[types.Health]()
	CollisionTimer = donburi.NewComponentType[types.Timer]()
	StuckCountdown = donburi.NewComponentType[types.Countdown]()
	Index          = donburi.NewComponentType[types.Index]()
	Item           = donburi.NewComponentType[types.Item]()
)

// Tags
var (
	TagWASD        = donburi.NewTag()
	TagWASDFly     = donburi.NewTag()
	TagPlayer      = donburi.NewTag()
	TagDropItem    = donburi.NewTag()
	TagBlock       = donburi.NewTag()
	TagHarvestable = donburi.NewTag()
	TagDebugBox    = donburi.NewTag()
)
