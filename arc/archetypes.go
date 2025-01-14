package arc

import (
	"kar"
	"kar/items"
	"kar/res"

	"github.com/mlange-42/arche/ecs"
	gn "github.com/mlange-42/arche/generic"
	"github.com/setanarut/anim"
)

var (
	MapRect     = gn.NewMap1[Rect](&kar.WorldECS)
	MapItem     = gn.NewMap4[ItemID, Durability, Rect, ItemTimers](&kar.WorldECS)
	MapSnowBall = gn.NewMap3[ItemID, Rect, Velocity](&kar.WorldECS)
	MapPlayer   = gn.NewMap3[anim.AnimationPlayer, Health, Rect](&kar.WorldECS)
)

// Query Filters
var (
	FilterMapSnowBall = gn.NewFilter3[ItemID, Rect, Velocity]().Exclusive()
	FilterRect        = gn.NewFilter1[Rect]()
	FilterAnimPlayer  = gn.NewFilter1[anim.AnimationPlayer]()
	FilterItem        = gn.NewFilter4[ItemID, Rect, ItemTimers, Durability]()
	FilterPlayer      = gn.NewFilter3[anim.AnimationPlayer, Health, Rect]()
)

func SpawnPlayer(x, y float64) ecs.Entity {
	AP := anim.NewAnimationPlayer(res.PlayerAtlas)
	AP.NewAnimationState("idleRight", 0, 0, 16, 16, 1, false, false).FPS = 1
	AP.NewAnimationState("idleUp", 208, 0, 16, 16, 1, false, false).FPS = 1
	AP.NewAnimationState("idleDown", 224, 0, 16, 16, 1, false, false).FPS = 1
	AP.NewAnimationState("walkRight", 16, 0, 16, 16, 4, false, false)
	AP.NewAnimationState("jump", 16*5, 0, 16, 16, 1, false, false)
	AP.NewAnimationState("skidding", 16*6, 0, 16, 16, 1, false, false)
	AP.NewAnimationState("attackDown", 16*7, 0, 16, 16, 2, false, false).FPS = 8
	AP.NewAnimationState("attackRight", 144, 0, 16, 16, 2, false, false).FPS = 8
	AP.NewAnimationState("attackWalk", 0, 16, 16, 16, 4, false, false).FPS = 8
	AP.NewAnimationState("attackUp", 16*11, 0, 16, 16, 2, false, false).FPS = 8
	AP.SetState("idleRight")
	return MapPlayer.NewWith(
		AP,
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
	Id         uint16
	Durability int
}
