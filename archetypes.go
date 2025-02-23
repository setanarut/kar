package kar

import (
	"kar/items"
	"kar/v"
	"math/rand/v2"

	"github.com/mlange-42/ark/ecs"
)

var (
	mapPosition         = ecs.NewMap1[Position](&world)
	mapAABB             = ecs.NewMap[AABB](&world)
	mapDurability       = ecs.NewMap[Durability](&world)
	mapHealth           = ecs.NewMap[Health](&world)
	mapCollisionDelayer = ecs.NewMap[CollisionDelayer](&world)
	mapEnemy            = ecs.NewMap3[Position, Velocity, AI](&world)
	mapProjectile       = ecs.NewMap3[ItemID, Position, Velocity](&world)
	mapDroppedItem      = ecs.NewMap4[ItemID, Position, AnimationIndex, CollisionDelayer](&world)
	mapEffect           = ecs.NewMap4[ItemID, Position, Velocity, Rotation](&world)
	mapPlayer           = ecs.NewMap5[AABB, Velocity, Health, Controller, Facing](&world)
)

// Query Filters
var (
	filterCollisionDelayer = ecs.NewFilter1[CollisionDelayer](&world)
	filterEnemy            = ecs.NewFilter3[Position, Velocity, AI](&world)
	filterProjectile       = ecs.NewFilter3[ItemID, Position, Velocity](&world).Without(ecs.C[Rotation]())
	filterDroppedItem      = ecs.NewFilter3[ItemID, Position, AnimationIndex](&world)
	filterEffect           = ecs.NewFilter4[ItemID, Position, Velocity, Rotation](&world)
)

func SpawnItem(pos Vec, id uint8, durability int) ecs.Entity {
	e := mapDroppedItem.NewEntity(
		&ItemID{id},
		&Position{pos},
		&AnimationIndex{rand.IntN(len(Sinspace) - 1)},
		&CollisionDelayer{ItemCollisionDelay},
	)
	if items.HasTag(id, items.Tool) {
		mapDurability.Add(e, &Durability{durability})
	}
	return e
}

func SpawnEnemy(pos, vel Vec) ecs.Entity {
	return mapEnemy.NewEntity(
		&Position{pos},
		&Velocity{vel},
		&AI{"worm"},
	)
}

func SpawnEffect(id uint8, pos Vec) {
	mapEffect.NewEntity(&ItemID{id}, &Position{Vec{pos.X - 10, pos.Y - 10}}, &Velocity{v.Left}, &Rotation{-0.1})
	mapEffect.NewEntity(&ItemID{id}, &Position{Vec{pos.X + 2, pos.Y - 10}}, &Velocity{v.Right}, &Rotation{0.1})
	mapEffect.NewEntity(&ItemID{id}, &Position{Vec{pos.X - 10, pos.Y + 2}}, &Velocity{Vec{-0.5, 0}}, &Rotation{-0.1})
	mapEffect.NewEntity(&ItemID{id}, &Position{Vec{pos.X + 2, pos.Y + 2}}, &Velocity{Vec{0.5, 0}}, &Rotation{0.1})
}

func SpawnProjectile(id uint8, pos, vel Vec) ecs.Entity {
	return mapProjectile.NewEntity(
		&ItemID{id},
		&Position{pos},
		&Velocity{vel},
	)
}

func SpawnPlayer(pos Vec) ecs.Entity {
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
	return mapPlayer.NewEntity(
		&AABB{
			Pos:  pos,
			Half: Vec{8, 8},
		},
		&Velocity{},
		&Health{20, 20},
		ctrl,
		&Facing{v.Down},
	)
}
