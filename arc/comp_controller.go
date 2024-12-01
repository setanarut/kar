package arc

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Controller struct {
	MinSpeed         float64
	MaxSpeed         float64
	MaxWalkSpeed     float64
	MaxFallSpeed     float64
	MaxFallSpeedCap  float64
	MinSlowDownSpeed float64
	WalkAcceleration float64
	RunAcceleration  float64
	WalkFriction     float64
	SkidFriction     float64
	StompSpeed       float64
	StompSpeedCap    float64
	JumpSpeed        [3]float64
	LongJumpGravity  [3]float64
	Gravity          float64
	SpeedThresholds  [2]float64
	// states
	IsFacingLeft bool
	IsRunning    bool
	IsJumping    bool
	IsFalling    bool
	IsSkidding   bool
	IsCrouching  bool
	IsOnFloor    bool
	// private
	minSpeedValue       float64
	maxSpeedValue       float64
	accel               float64
	speedThresholdIndex int
}

func NewController() *Controller {
	pc := &Controller{
		MinSpeed:            0.07421875,
		MaxSpeed:            2.5625,
		MaxWalkSpeed:        1.5625,
		MaxFallSpeed:        4.5,
		MaxFallSpeedCap:     4,
		MinSlowDownSpeed:    0.5625,
		WalkAcceleration:    0.037109375,
		RunAcceleration:     0.0556640625,
		WalkFriction:        0.05078125,
		SkidFriction:        0.1015625,
		StompSpeed:          4,
		StompSpeedCap:       4,
		JumpSpeed:           [3]float64{-4, -4, -5},
		LongJumpGravity:     [3]float64{0.12, 0.11, 0.15},
		Gravity:             0.43,
		SpeedThresholds:     [2]float64{1, 2.3125},
		IsFacingLeft:        false,
		IsRunning:           false,
		IsJumping:           false,
		IsFalling:           false,
		IsSkidding:          false,
		IsCrouching:         false,
		IsOnFloor:           false,
		speedThresholdIndex: 0,
	}

	pc.minSpeedValue = pc.MinSpeed
	pc.maxSpeedValue = pc.MaxSpeed
	pc.accel = pc.WalkAcceleration

	return pc
}

func (pc *Controller) SetPhyicsScale(s float64) {
	pc.MinSpeed *= s
	pc.MaxSpeed *= s
	pc.MaxWalkSpeed *= s
	pc.MaxFallSpeed *= s
	pc.MaxFallSpeedCap *= s
	pc.MinSlowDownSpeed *= s
	pc.WalkAcceleration *= s
	pc.RunAcceleration *= s
	pc.WalkFriction *= s
	pc.SkidFriction *= s
	pc.StompSpeed *= s
	pc.StompSpeedCap *= s
	pc.JumpSpeed[0] *= s
	pc.JumpSpeed[1] *= s
	pc.JumpSpeed[2] *= s
	pc.LongJumpGravity[0] *= s
	pc.LongJumpGravity[1] *= s
	pc.LongJumpGravity[2] *= s
	pc.Gravity *= s
	pc.SpeedThresholds[0] *= s
	pc.SpeedThresholds[1] *= s
}

func (pc *Controller) ProcessVelocity(inputAxis, vel vec2) vec2 {

	if pc.IsOnFloor {
		pc.IsRunning = ebiten.IsKeyPressed(ebiten.KeyShift)
		pc.IsCrouching = ebiten.IsKeyPressed(ebiten.KeyDown)
		if pc.IsCrouching && inputAxis.X != 0 {
			pc.IsCrouching = false
			inputAxis.X = 0.0
		}
	}

	if pc.IsOnFloor {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			pc.IsJumping = true
			speed := math.Abs(vel.X)
			pc.speedThresholdIndex = 0
			if speed >= pc.SpeedThresholds[1] {
				pc.speedThresholdIndex = 2
			} else if speed >= pc.SpeedThresholds[0] {
				pc.speedThresholdIndex = 1
			}

			vel.Y = pc.JumpSpeed[pc.speedThresholdIndex]

		}
	} else {
		gravityValue := pc.Gravity
		if ebiten.IsKeyPressed(ebiten.KeySpace) && pc.IsJumping && vel.Y < 0 {
			gravityValue = pc.LongJumpGravity[pc.speedThresholdIndex]
		}
		vel.Y += gravityValue
		if vel.Y > pc.MaxFallSpeedCap {
			vel.Y = pc.MaxFallSpeedCap
		}
	}

	// Update states
	if vel.Y > 0 {
		pc.IsJumping = false
		pc.IsFalling = true
	} else if pc.IsOnFloor {
		pc.IsFalling = false
	}

	if inputAxis.X != 0 {
		if pc.IsOnFloor {
			if vel.X != 0 {
				pc.IsFacingLeft = inputAxis.X < 0.0
				pc.IsSkidding = vel.X < 0.0 != pc.IsFacingLeft
			}
			if pc.IsSkidding {
				pc.minSpeedValue = pc.MinSlowDownSpeed
				pc.maxSpeedValue = pc.MaxWalkSpeed
				pc.accel = pc.SkidFriction
			} else if pc.IsRunning {
				pc.minSpeedValue = pc.MinSpeed
				pc.maxSpeedValue = pc.MaxSpeed
				pc.accel = pc.RunAcceleration
			} else {
				pc.minSpeedValue = pc.MinSpeed
				pc.maxSpeedValue = pc.MaxWalkSpeed
				pc.accel = pc.WalkAcceleration
			}
		} else if pc.IsRunning && math.Abs(vel.X) > pc.MaxWalkSpeed {
			pc.maxSpeedValue = pc.MaxSpeed
		} else {
			pc.maxSpeedValue = pc.MaxWalkSpeed
		}
		targetSpeed := inputAxis.X * pc.maxSpeedValue

		// Manually implementing moveToward()
		if vel.X < targetSpeed {
			vel.X += pc.accel
			if vel.X > targetSpeed {
				vel.X = targetSpeed
			}
		} else if vel.X > targetSpeed {
			vel.X -= pc.accel
			if vel.X < targetSpeed {
				vel.X = targetSpeed
			}
		}

	} else if pc.IsOnFloor && vel.X != 0 {
		if !pc.IsSkidding {
			pc.accel = pc.WalkFriction
		}
		if inputAxis.Y != 0 {
			pc.minSpeedValue = pc.MinSlowDownSpeed
		} else {
			pc.minSpeedValue = pc.MinSpeed
		}
		if math.Abs(vel.X) < pc.minSpeedValue {
			vel.X = 0.0
		} else {
			// Manually implementing moveToward() for deceleration
			if vel.X > 0 {
				vel.X -= pc.accel
				if vel.X < 0 {
					vel.X = 0
				}
			} else {
				vel.X += pc.accel
				if vel.X > 0 {
					vel.X = 0
				}
			}
		}
	}
	if math.Abs(vel.X) < pc.MinSlowDownSpeed {
		pc.IsSkidding = false
	}

	return vel
}
