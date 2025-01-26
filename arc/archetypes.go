package arc

import (
	"kar"
	"kar/items"

	"github.com/mlange-42/arche/ecs"
	gn "github.com/mlange-42/arche/generic"
	"github.com/setanarut/anim"
)

var (
	MapRect     = gn.NewMap1[Rect](&kar.WorldECS)
	MapHealth   = gn.NewMap[Health](&kar.WorldECS)
	MapEnemy    = gn.NewMap3[Rect, Velocity, Health](&kar.WorldECS)
	MapItem     = gn.NewMap4[ItemID, Durability, Rect, ItemTimers](&kar.WorldECS)
	MapSnowBall = gn.NewMap3[ItemID, Rect, Velocity](&kar.WorldECS)
	MapPlayer   = gn.NewMap2[Health, Rect](&kar.WorldECS)
)

// Query Filters
var (
	FilterPlayer      = gn.NewFilter2[Health, Rect]().Exclusive()
	FilterEnemy       = gn.NewFilter3[Rect, Velocity, Health]().Exclusive()
	FilterMapSnowBall = gn.NewFilter3[ItemID, Rect, Velocity]().Exclusive()
	FilterRect        = gn.NewFilter1[Rect]()
	FilterAnimPlayer  = gn.NewFilter1[anim.AnimationPlayer]()
	FilterItem        = gn.NewFilter4[ItemID, Rect, ItemTimers, Durability]()
)

type AnimPlayerData struct {
	CurrentState string
	CurrentAtlas string
	Paused       bool
	Tick         float64
	CurrentIndex int
}

func SpawnPlayer(x, y float64) ecs.Entity {

	return MapPlayer.NewWith(
		&Health{20, 20},
		&Rect{x - 16*0.5, y - 16*0.5, 16, 16},
	)

}

func SpawnItem(data SpawnData, animIndex int) ecs.Entity {
	return MapItem.NewWith(
		&ItemID{data.Id},
		&Durability{data.Durability},
		&Rect{data.X, data.Y, 8, 8},
		&ItemTimers{kar.ItemCollisionDelay, animIndex},
	)
}

func SpawnEnemy(x, y, vx, vy float64) ecs.Entity {
	return MapEnemy.NewWith(
		&Rect{x, y, 8, 8},
		&Velocity{vx, vy},
		&Health{
			Current: 0,
			Max:     20,
		},
	)
}
func SpawnSnowBall(x, y, vx, vy float64) ecs.Entity {
	return MapSnowBall.NewWith(
		&ItemID{items.Snowball},
		&Rect{x, y, 8, 8},
		&Velocity{vx, vy},
	)
}

// SpawnData is a helper for delaying spawn events
type SpawnData struct {
	X, Y       float64
	Id         uint8
	Durability int
}
