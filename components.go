package kar

import (
	"time"
)

type ItemID struct {
	ID uint8
}
type Facing Vec
type Velocity Vec
type Position Vec

type Rotation struct {
	Angle float64
}

type AI struct {
	Name string
}

// AnimationIndex holds timing-related data for item animations.
// It tracks the current frame index for dropped item animations.
type AnimationIndex struct {
	Index int
}

type CollisionDelayer struct {
	Duration time.Duration
}
type Durability struct {
	Durability int
}

type Health struct {
	Current int
	Max     int
}

type Controller struct {
	Acceleration                        float64
	AirSkiddingDecel                    float64
	CurrentState                        string
	FallingDamageTempPosY               float64
	Gravity                             float64
	JumpBoost                           float64
	JumpBoostMultiplier                 float64
	JumpHoldTime                        float64
	JumpPower                           float64
	JumpReleaseTimer                    float64
	JumpTimer                           float64
	MaxFallSpeed                        float64
	MaxRunSpeed                         float64
	MaxWalkSpeed                        float64
	MinSpeedThresForJumpBoostMultiplier float64
	PreviousState                       string
	RunAcceleration                     float64
	RunDeceleration                     float64
	ShortJumpVelocity                   float64
	SkiddingFriction                    float64
	SkiddingJumpEnabled                 bool
	SpeedJumpFactor                     float64
	WalkAcceleration                    float64
	WalkDeceleration                    float64
}
