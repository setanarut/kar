package arc

import (
	"image"
	"kar"
	"kar/items"
	"math/rand/v2"

	"github.com/mlange-42/arche/ecs"
	gn "github.com/mlange-42/arche/generic"
	"github.com/setanarut/anim"
)

var (
	MapPlayer           = gn.NewMap6[Position, Size, Velocity, Health, Controller, Facing](&kar.ECWorld)
	MapEnemy            = gn.NewMap4[Position, Size, Velocity, Health](&kar.ECWorld)
	MapRect             = gn.NewMap2[Position, Size](&kar.ECWorld)
	MapHealth           = gn.NewMap[Health](&kar.ECWorld)
	MapDurability       = gn.NewMap[Durability](&kar.ECWorld)
	MapPosition         = gn.NewMap[Position](&kar.ECWorld)
	MapSize             = gn.NewMap[Size](&kar.ECWorld)
	MapDroppedItem      = gn.NewMap4[ItemID, Position, AnimationIndex, CollisionDelayer](&kar.ECWorld)
	MapDroppedToolItem  = gn.NewMap5[ItemID, Position, AnimationIndex, CollisionDelayer, Durability](&kar.ECWorld)
	MapProjectile       = gn.NewMap3[ItemID, Position, Velocity](&kar.ECWorld)
	MapCollisionDelayer = gn.NewMap1[CollisionDelayer](&kar.ECWorld)
)

// Query Filters
var (
	FilterPlayer           = gn.NewFilter6[Position, Size, Velocity, Health, Controller, Facing]()
	FilterEnemy            = gn.NewFilter4[Position, Size, Velocity, Health]().Exclusive()
	FilterProjectile       = gn.NewFilter3[ItemID, Position, Velocity]().Exclusive()
	FilterRect             = gn.NewFilter2[Position, Size]()
	FilterPosition         = gn.NewFilter1[Position]()
	FilterCollisionDelayer = gn.NewFilter1[CollisionDelayer]()
	FilterAnimPlayer       = gn.NewFilter1[anim.AnimationPlayer]()
	FilterDroppedItem      = gn.NewFilter5[
		ItemID,
		Position,
		AnimationIndex,
		CollisionDelayer,
		Durability,
	]().
		Optional(gn.T[CollisionDelayer]()).
		Optional(gn.T[Durability]())
)

func SpawnItem(x, y float64, id uint8, durability int) ecs.Entity {

	if items.HasTag(id, items.Tool) {
		return MapDroppedToolItem.NewWith(
			&ItemID{id},
			&Position{x, y},
			&AnimationIndex{rand.IntN(len(kar.Sinspace) - 1)},
			&CollisionDelayer{kar.ItemCollisionDelay},
			&Durability{durability},
		)
	} else {
		return MapDroppedItem.NewWith(
			&ItemID{id},
			&Position{x, y},
			&AnimationIndex{rand.IntN(len(kar.Sinspace) - 1)},
			&CollisionDelayer{kar.ItemCollisionDelay},
		)
	}
}

func SpawnEnemy(x, y, vx, vy float64) ecs.Entity {
	return MapEnemy.NewWith(
		&Position{x, y},
		&Size{8, 8},
		&Velocity{vx, vy},
		&Health{
			Current: 0,
			Max:     20,
		},
	)
}

func SpawnProjectile(id uint8, x, y, vx, vy float64) ecs.Entity {
	return MapProjectile.NewWith(
		&ItemID{id},
		&Position{x, y},
		&Velocity{vx, vy},
	)
}

// SpawnData is a helper for delaying spawn events
type SpawnData struct {
	X, Y       float64
	Id         uint8
	Durability int
}

func SpawnPlayer(centerX, centerY float64) ecs.Entity {
	return MapPlayer.NewWith(
		&Position{centerX - 16*0.5, centerY - 16*0.5},
		&Size{16, 16},
		&Velocity{0, 0},
		&Health{20, 20},
		DefaultController(),
		&Facing{image.Point{0, 1}},
	)
}

func DefaultController() *Controller {
	return &Controller{
		CurrentState:                        "falling",
		Gravity:                             0.19,
		JumpPower:                           -3.7,
		MaxFallSpeed:                        100.0,
		MaxRunSpeed:                         3.0,
		MaxWalkSpeed:                        1.6,
		Acceleration:                        0.08,
		Deceleration:                        0.1,
		JumpHoldTime:                        20.0,
		JumpBoost:                           -0.1,
		MinSpeedThresForJumpBoostMultiplier: 0.1,
		JumpBoostMultiplier:                 1.01,
		SpeedJumpFactor:                     0.3,
		ShortJumpVelocity:                   -2.0,
		JumpReleaseTimer:                    5,
		WalkAcceleration:                    0.04,
		WalkDeceleration:                    0.04,
		RunAcceleration:                     0.04,
		RunDeceleration:                     0.04,
		SkiddingJumpEnabled:                 true,
	}
}
