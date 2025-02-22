package kar

import (
	"kar/items"
	"math/rand/v2"

	"github.com/mlange-42/ark/ecs"
)

var (
	MapPlayer           = ecs.NewMap5[AABB, Velocity, Health, Controller, Facing](&world)
	MapEnemy            = ecs.NewMap3[Position, Velocity, AI](&world)
	MapAABB             = ecs.NewMap[AABB](&world)
	MapHealth           = ecs.NewMap[Health](&world)
	MapDurability       = ecs.NewMap[Durability](&world)
	MapPosition         = ecs.NewMap[Position](&world)
	MapDroppedItem      = ecs.NewMap4[ItemID, Position, AnimationIndex, CollisionDelayer](&world)
	MapDroppedToolItem  = ecs.NewMap5[ItemID, Position, AnimationIndex, CollisionDelayer, Durability](&world)
	MapProjectile       = ecs.NewMap3[ItemID, Position, Velocity](&world)
	MapCollisionDelayer = ecs.NewMap[CollisionDelayer](&world)
	MapEffect           = ecs.NewMap4[ItemID, Position, Velocity, Rotation](&world)
)

// Query Filters
var (
	FilterPlayer           = ecs.NewFilter5[AABB, Velocity, Health, Controller, Facing](&world)
	FilterEnemy            = ecs.NewFilter3[Position, Velocity, AI](&world)
	FilterProjectile       = ecs.NewFilter3[ItemID, Position, Velocity](&world).Without(ecs.C[Rotation]()) // Exclusive
	FilterCollisionDelayer = ecs.NewFilter1[CollisionDelayer](&world)

	FilterDroppedItem = ecs.NewFilter5[ItemID, Position, AnimationIndex, CollisionDelayer, Durability](&world)
	// Optional(gn.T[CollisionDelayer]())
	// Optional(gn.T[Durability]())

	FilterEffect = ecs.NewFilter4[ItemID, Position, Velocity, Rotation](&world)
)

func SpawnItem(x, y float64, id uint8, durability int) ecs.Entity {

	if items.HasTag(id, items.Tool) {
		return MapDroppedToolItem.NewEntity(
			&ItemID{id},
			&Position{x, y},
			&AnimationIndex{rand.IntN(len(Sinspace) - 1)},
			&CollisionDelayer{ItemCollisionDelay},
			&Durability{durability},
		)
	} else {
		return MapDroppedItem.NewEntity(
			&ItemID{id},
			&Position{x, y},
			&AnimationIndex{rand.IntN(len(Sinspace) - 1)},
			&CollisionDelayer{ItemCollisionDelay},
		)
	}
}

func SpawnEnemy(x, y, vx, vy float64) ecs.Entity {
	return MapEnemy.NewEntity(
		&Position{x, y},
		&Velocity{vx, vy},
		&AI{"worm"},
	)
}

func SpawnEffect(id uint8, x, y float64) {
	MapEffect.NewEntity(&ItemID{id}, &Position{x - 10, y - 10}, &Velocity{-1, 0}, &Rotation{-0.1})
	MapEffect.NewEntity(&ItemID{id}, &Position{x + 2, y - 10}, &Velocity{1, 0}, &Rotation{0.1})
	MapEffect.NewEntity(&ItemID{id}, &Position{x - 10, y + 2}, &Velocity{-0.5, 0}, &Rotation{-0.1})
	MapEffect.NewEntity(&ItemID{id}, &Position{x + 2, y + 2}, &Velocity{0.5, 0}, &Rotation{0.1})
}

func SpawnProjectile(id uint8, x, y, vx, vy float64) ecs.Entity {
	return MapProjectile.NewEntity(
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
	return MapPlayer.NewEntity(
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
