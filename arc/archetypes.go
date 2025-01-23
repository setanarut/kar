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
	MapHealth   = gn.NewMap[Health](&kar.WorldECS)
	MapEnemy    = gn.NewMap3[Rect, Velocity, Health](&kar.WorldECS)
	MapItem     = gn.NewMap4[ItemID, Durability, Rect, ItemTimers](&kar.WorldECS)
	MapSnowBall = gn.NewMap3[ItemID, Rect, Velocity](&kar.WorldECS)
	MapPlayer   = gn.NewMap3[anim.AnimationPlayer, Health, Rect](&kar.WorldECS)
)

// Query Filters
var (
	FilterEnemy       = gn.NewFilter3[Rect, Velocity, Health]().Exclusive()
	FilterMapSnowBall = gn.NewFilter3[ItemID, Rect, Velocity]().Exclusive()
	FilterRect        = gn.NewFilter1[Rect]()
	FilterAnimPlayer  = gn.NewFilter1[anim.AnimationPlayer]()
	FilterItem        = gn.NewFilter4[ItemID, Rect, ItemTimers, Durability]()
	FilterPlayer      = gn.NewFilter3[anim.AnimationPlayer, Health, Rect]()
)

func SpawnPlayer(x, y float64) ecs.Entity {
	ap := anim.NewAnimationPlayer(
		&anim.Atlas{"Default", res.Player},
		&anim.Atlas{"WoodenAxe", res.PlayerWoodenAxeAtlas},
		&anim.Atlas{"StoneAxe", res.PlayerStoneAxeAtlas},
		&anim.Atlas{"IronAxe", res.PlayerIronAxeAtlas},
		&anim.Atlas{"DiamondAxe", res.PlayerDiamondAxeAtlas},
		&anim.Atlas{"WoodenPickaxe", res.PlayerWoodenPickaxeAtlas},
		&anim.Atlas{"StonePickaxe", res.PlayerStonePickaxeAtlas},
		&anim.Atlas{"IronPickaxe", res.PlayerIronPickaxeAtlas},
		&anim.Atlas{"DiamondPickaxe", res.PlayerDiamondPickaxeAtlas},
		&anim.Atlas{"WoodenShovel", res.PlayerWoodenShovelAtlas},
		&anim.Atlas{"StoneShovel", res.PlayerStoneShovelAtlas},
		&anim.Atlas{"IronShovel", res.PlayerIronShovelAtlas},
		&anim.Atlas{"DiamondShovel", res.PlayerDiamondShovelAtlas},
	)

	ap.NewState("idleRight", 0, 0, 16, 16, 1, false, false, 1)
	ap.NewState("idleUp", 208, 0, 16, 16, 1, false, false, 1)
	ap.NewState("idleDown", 224, 0, 16, 16, 1, false, false, 1)
	ap.NewState("walkRight", 16, 0, 16, 16, 4, false, false, 15)
	ap.NewState("jump", 16*5, 0, 16, 16, 1, false, false, 15)
	ap.NewState("skidding", 16*6, 0, 16, 16, 1, false, false, 15)
	ap.NewState("attackDown", 16*7, 0, 16, 16, 2, false, false, 8)
	ap.NewState("attackRight", 144, 0, 16, 16, 2, false, false, 8)
	ap.NewState("attackWalk", 0, 16, 16, 16, 4, false, false, 8)
	ap.NewState("attackUp", 16*11, 0, 16, 16, 2, false, false, 8)
	ap.CurrentState = "idleRight"
	return MapPlayer.NewWith(
		ap,
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
	Id         uint16
	Durability int
}
