package kar

import (
	"kar/items"
	"math/rand/v2"

	"github.com/mlange-42/arche/ecs"
	gn "github.com/mlange-42/arche/generic"
)

var (
	MapPlayer           = gn.NewMap5[AABB, Velocity, Health, Controller, Facing](&world)
	MapEnemy            = gn.NewMap3[Position, Velocity, AI](&world)
	MapAABB             = gn.NewMap[AABB](&world)
	MapHealth           = gn.NewMap[Health](&world)
	MapDurability       = gn.NewMap[Durability](&world)
	MapPosition         = gn.NewMap[Position](&world)
	MapDroppedItem      = gn.NewMap4[ItemID, Position, AnimationIndex, CollisionDelayer](&world)
	MapDroppedToolItem  = gn.NewMap5[ItemID, Position, AnimationIndex, CollisionDelayer, Durability](&world)
	MapProjectile       = gn.NewMap3[ItemID, Position, Velocity](&world)
	MapCollisionDelayer = gn.NewMap1[CollisionDelayer](&world)
	MapEffect           = gn.NewMap4[ItemID, Position, Velocity, Rotation](&world)
)

// Query Filters
var (
	FilterPlayer = gn.NewFilter5[
		AABB,
		Velocity,
		Health,
		Controller,
		Facing]()
	FilterEnemy            = gn.NewFilter3[Position, Velocity, AI]()
	FilterProjectile       = gn.NewFilter3[ItemID, Position, Velocity]().Exclusive()
	FilterCollisionDelayer = gn.NewFilter1[CollisionDelayer]()
	FilterDroppedItem      = gn.NewFilter5[
		ItemID,
		Position,
		AnimationIndex,
		CollisionDelayer,
		Durability,
	]().
		Optional(gn.T[CollisionDelayer]()).
		Optional(gn.T[Durability]())
	FilterEffect = gn.NewFilter4[ItemID, Position, Velocity, Rotation]().Exclusive()
)

func SpawnItem(x, y float64, id uint8, durability int) ecs.Entity {

	if items.HasTag(id, items.Tool) {
		return MapDroppedToolItem.NewWith(
			&ItemID{id},
			&Position{x, y},
			&AnimationIndex{rand.IntN(len(Sinspace) - 1)},
			&CollisionDelayer{ItemCollisionDelay},
			&Durability{durability},
		)
	} else {
		return MapDroppedItem.NewWith(
			&ItemID{id},
			&Position{x, y},
			&AnimationIndex{rand.IntN(len(Sinspace) - 1)},
			&CollisionDelayer{ItemCollisionDelay},
		)
	}
}

func SpawnEnemy(x, y, vx, vy float64) ecs.Entity {
	return MapEnemy.NewWith(
		&Position{x, y},
		&Velocity{vx, vy},
		&AI{"worm"},
	)
}

func SpawnEffect(id uint8, x, y float64) {
	MapEffect.NewWith(&ItemID{id}, &Position{x - 10, y - 10}, &Velocity{-1, 0}, &Rotation{-0.1})
	MapEffect.NewWith(&ItemID{id}, &Position{x + 2, y - 10}, &Velocity{1, 0}, &Rotation{0.1})
	MapEffect.NewWith(&ItemID{id}, &Position{x - 10, y + 2}, &Velocity{-0.5, 0}, &Rotation{-0.1})
	MapEffect.NewWith(&ItemID{id}, &Position{x + 2, y + 2}, &Velocity{0.5, 0}, &Rotation{0.1})
}

func SpawnProjectile(id uint8, x, y, vx, vy float64) ecs.Entity {
	return MapProjectile.NewWith(
		&ItemID{id},
		&Position{x, y},
		&Velocity{vx, vy},
	)
}

func SpawnPlayer(centerX, centerY float64) ecs.Entity {
	ctrl := &Controller{
		CurrentState:                        "falling",
		Gravity:                             0.19,
		JumpPower:                           -3.7,
		MaxFallSpeed:                        100.0,
		MaxRunSpeed:                         3.0,
		MaxWalkSpeed:                        1.6,
		Acceleration:                        0.08,
		SkiddingFriction:                    0.08,
		AirSkiddingDecel:                    0.1,
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
	return MapPlayer.NewWith(
		&AABB{
			Pos:  Vec{centerX, centerY},
			Half: Vec{8, 8},
		},
		&Velocity{0, 0},
		&Health{20, 20},
		ctrl,
		&Facing{0, 1},
	)
}
