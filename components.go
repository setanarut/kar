package kar

import (
	"time"
)

type ItemID uint8
type Rotation float64
type Facing Vec
type Velocity Vec
type Position Vec
type AI string
type PlatformType string
type Durability int
type AnimationIndex int // timing-related data for item animations.
type CollisionDelayer time.Duration

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

func ptr[T any](v T) *T { return &v }
