package controller

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

const scale = 2.0

const (
	CooldownTimeSec  = 3.0 * scale
	MaxFallSpeed     = 270.0 * scale
	MaxFallSpeedCap  = 240.0 * scale
	MaxSpeed         = 153.75 * scale
	MaxWalkSpeed     = 93.75 * scale
	MinSlowDownSpeed = 33.75 * scale
	MinSpeed         = 4.453125 * scale
	RunAcceleration  = 200.390625 * scale
	SkidFriction     = 365.625 * scale
	StompSpeed       = 240.0 * scale
	StompSpeedCap    = -60.0 * scale
	WalkAcceleration = 133.59375 * scale
	WalkFriction     = 182.8125 * scale
)

var (
	jumpSpeeds        = [3]float64{-240.0 * scale, -240.0 * scale, -300.0 * scale}
	longJumpGravities = [3]float64{450.0 * scale, 421.875 * scale, 562.5 * scale}
	gravities         = [3]float64{1575.0 * scale, 1350.0 * scale, 2025.0 * scale}
	speedThresholds   = [2]float64{60 * scale, 138.75 * scale}
)

// States
var (
	isAttacking                   bool
	isCrouching                   bool
	IsFacingLeft, IsFacingRight   bool
	isFacingUp, isFacingDown      bool
	isFalling                     bool
	isIdle                        bool
	isOnFloor                     bool
	isRunning                     bool
	isSkiding                     bool
	isDigDown, isDigUp            bool
	isWalkingLeft, isWalkingRight bool
	// isJumping    bool
)
var (
	minSpeedTemp = MinSpeed
	maxSpeedTemp = MaxWalkSpeed
	acceleration = WalkAcceleration
	delta        = 1 / 60.0

	speedThreshold int = 0
)

var (
	right = vec.Vec2{1, 0}
	left  = vec.Vec2{-1, 0}
	down  = vec.Vec2{0, 1}
	up    = vec.Vec2{0, -1}
	zero  = vec.Vec2{0, 0}
)

func VelocityFunc(body *cm.Body, grav vec.Vec2, damping, dt float64) {

	velocity := body.Velocity()
	isOnFloor = OnFloor(body)
	inputAxis := GetAxis()

	inputAxisLast := vec.Vec2{}
	if !inputAxis.Equal(vec.Vec2{}) {
		inputAxisLast = inputAxis
	}

	isAttacking = ebiten.IsKeyPressed(ebiten.KeyShiftRight)
	isIdle = inputAxis.Equal(zero) && !isAttacking && isOnFloor

	isFacingDown = inputAxisLast.Equal(down) || inputAxis.Equal(down)
	isFacingUp = inputAxisLast.Equal(up) || inputAxis.Equal(up)
	IsFacingRight = inputAxisLast.Equal(right) || inputAxis.Equal(right)
	IsFacingLeft = inputAxisLast.Equal(left) || inputAxis.Equal(left)

	isWalkingLeft = !isIdle && isOnFloor && body.Velocity().X < 0.0
	isWalkingRight = !isIdle && isOnFloor && body.Velocity().X > 0.0

	isDigDown = isFacingDown && isAttacking
	isDigUp = isFacingUp && isAttacking

	if isOnFloor {
		isRunning = ebiten.IsKeyPressed(ebiten.KeyAltRight)
		isCrouching = ebiten.IsKeyPressed(ebiten.KeyDown)
		if isCrouching && inputAxis.X != 0 {
			isCrouching = false
			inputAxis.X = 0.0
		}
	}

	if isOnFloor {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			var speed = math.Abs(velocity.X)
			speedThreshold = len(speedThresholds)

			for i := 0; i < len(speedThresholds); i++ {
				if speed < speedThresholds[i] {
					speedThreshold = i
					break
				}
			}
			velocity.Y = jumpSpeeds[speedThreshold]

		}
	} else {
		var gravity = gravities[speedThreshold]
		if ebiten.IsKeyPressed(ebiten.KeySpace) && !isFalling {
			gravity = longJumpGravities[speedThreshold]
		}
		velocity.Y = velocity.Y + gravity*delta
		if velocity.Y > MaxFallSpeed {
			velocity.Y = MaxFallSpeedCap
		}
	}

	if velocity.Y > 0 {
		isFalling = true
	} else if isOnFloor {
		isFalling = false
	}

	if inputAxis.X != 0 {
		if isOnFloor {
			if velocity.X != 0 {
				IsFacingLeft = inputAxis.X < 0.0
				isSkiding = velocity.X < 0.0 != IsFacingLeft
			}
			if isSkiding {
				minSpeedTemp = MinSlowDownSpeed
				maxSpeedTemp = MaxWalkSpeed
				acceleration = SkidFriction
			} else if isRunning {
				minSpeedTemp = MinSpeed
				maxSpeedTemp = MaxSpeed
				acceleration = RunAcceleration
			} else {
				minSpeedTemp = MinSpeed
				maxSpeedTemp = MaxWalkSpeed
				acceleration = WalkAcceleration
			}
		} else if isRunning && math.Abs(velocity.X) > MaxWalkSpeed {
			maxSpeedTemp = MaxSpeed
		} else {
			maxSpeedTemp = MaxWalkSpeed
		}
		var target_speed = inputAxis.X * maxSpeedTemp
		velocity.X = MoveToward(velocity.X, target_speed, acceleration*delta)
	} else if isOnFloor && velocity.X != 0 {
		if !isSkiding {
			acceleration = WalkFriction
		}
		if inputAxis.Y != 0 {
			minSpeedTemp = MinSlowDownSpeed
		} else {
			minSpeedTemp = MinSpeed
		}
		if math.Abs(velocity.X) < minSpeedTemp {
			velocity.X = 0.0
		} else {
			velocity.X = MoveToward(velocity.X, 0.0, acceleration*delta)
		}
	}
	if math.Abs(velocity.X) < MinSlowDownSpeed {
		isSkiding = false
	}
	body.SetVelocityVector(velocity)
}

func OnFloor(b *cm.Body) bool {
	groundNormal := vec.Vec2{}
	b.EachArbiter(func(arb *cm.Arbiter) {
		n := arb.Normal().Neg()
		if n.Y < groundNormal.Y {
			groundNormal = n
		}
	})
	return groundNormal.Y < 0
}

func GetAxis() vec.Vec2 {
	axis := vec.Vec2{}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		axis.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		axis.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		axis.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		axis.X += 1
	}
	return axis
}

func MoveToward(from, to, delta float64) float64 {
	if math.Abs(to-from) <= delta {
		return to
	}
	if to > from {
		return from + delta
	}
	return from - delta
}
