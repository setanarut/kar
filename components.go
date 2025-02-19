package kar

import (
	"time"
)

type ItemID struct {
	ID uint8
}

type Velocity Vec
type Position Vec
type Facing Vec

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
	CurrentState                        string
	PreviousState                       string
	Gravity                             float64
	JumpPower                           float64
	MaxFallSpeed                        float64
	MaxRunSpeed                         float64
	MaxWalkSpeed                        float64
	Acceleration                        float64
	SkiddingFriction                    float64
	AirSkiddingDecel                    float64
	JumpHoldTime                        float64
	JumpBoost                           float64
	JumpTimer                           float64
	MinSpeedThresForJumpBoostMultiplier float64
	JumpBoostMultiplier                 float64
	SpeedJumpFactor                     float64
	ShortJumpVelocity                   float64
	JumpReleaseTimer                    float64
	WalkAcceleration                    float64
	WalkDeceleration                    float64
	RunAcceleration                     float64
	RunDeceleration                     float64
	pVelocityAbs                        float64
	FallingDamageTempPosY               float64
	IsOnFloor                           bool
	IsSkidding                          bool
	IsFalling                           bool
	SkiddingJumpEnabled                 bool
	IsBreakKeyPressed                   bool
	IsAttackKeyJustPressed              bool
	IsJumpKeyPressed                    bool
	IsJumpKeyJustPressed                bool
	IsRunKeyPressed                     bool
	InputAxis                           Vec
}
