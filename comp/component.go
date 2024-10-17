package comp

import (
	"kar/types"

	"github.com/setanarut/anim"
	"github.com/setanarut/cm"

	"github.com/yohamta/donburi"
)

// Components
var (
	Mobile      = donburi.NewComponentType[types.Mobile]()
	AnimPlayer  = donburi.NewComponentType[anim.AnimationPlayer]()
	Body        = donburi.NewComponentType[cm.Body]()
	DrawOptions = donburi.NewComponentType[types.DrawOptions]()
	Inventory   = donburi.NewComponentType[types.Inventory]()
	Damage      = donburi.NewComponentType[float64](1.0)
	Health      = donburi.NewComponentType[types.Health]()
	Timer       = donburi.NewComponentType[types.Timer]()
	SpawnTimer  = donburi.NewComponentType[types.Timer]()
	Index       = donburi.NewComponentType[types.Index]()
	AI          = donburi.NewComponentType[types.DataAI]()
	Item        = donburi.NewComponentType[types.Item]()
)

// Tags
var (
	TagWASD        = donburi.NewTag()
	TagWASDFly     = donburi.NewTag()
	TagPlayer      = donburi.NewTag()
	TagEnemy       = donburi.NewTag()
	TagItem        = donburi.NewTag()
	TagBlock       = donburi.NewTag()
	TagHarvestable = donburi.NewTag()
	TagDebugBox    = donburi.NewTag()
)
