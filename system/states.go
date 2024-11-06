package system

import (
	"fmt"
	"kar"
	"kar/engine/mathutil"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

const scale = 2.1

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
	isAttacking                 bool
	isCrouching                 bool
	IsFacingLeft, IsFacingRight bool
	isFacingUp, isFacingDown    bool
	isFalling                   bool
	isSkiding                   bool
	isDigDown, isDigUp          bool
	isWalking                   bool
	isMovingHorizontal          bool
	isMovingVertical            bool
	isMoving                    bool
	isRunning                   bool
	isIdle                      bool
)

var (
	minSpeedTemp       = MinSpeed
	maxSpeedTemp       = MaxWalkSpeed
	acceleration       = WalkAcceleration
	delta              = 1 / 60.0
	speedThreshold int = 0
)

var (
	right = vec.Vec2{1, 0}
	left  = vec.Vec2{-1, 0}
	down  = vec.Vec2{0, 1}
	up    = vec.Vec2{0, -1}
	zero  = vec.Vec2{0, 0}
)

var isOnFloor = true

const MovingThreshold float64 = 0.1

type States struct{}

func (sys *States) Init() {

	fsm.SetState(fsm.Idle)
	playerBody.SetVelocityUpdateFunc(VelocityFunc)

}
func (sys *States) Draw() {}
func (sys *States) Update() {

	if kar.WorldECS.Alive(playerEntity) {
		velocity := playerBody.Velocity()
		isOnFloor = OnFloor(playerBody)

		isCrouching = isOnFloor && ebiten.IsKeyPressed(ebiten.KeyDown)
		isRunning = isOnFloor && ebiten.IsKeyPressed(ebiten.KeyAltRight)

		isFacingDown = inputAxisLast.Equal(down) || inputAxis.Equal(down)
		isFacingUp = inputAxisLast.Equal(up) || inputAxis.Equal(up)
		IsFacingRight = inputAxisLast.Equal(right) || inputAxis.Equal(right)
		IsFacingLeft = inputAxisLast.Equal(left) || inputAxis.Equal(left)
		isDigDown = isFacingDown && isAttacking
		isDigUp = isFacingUp && isAttacking
		isMovingHorizontal = math.Abs(velocity.X) > MovingThreshold
		isMovingVertical = math.Abs(velocity.Y) > MovingThreshold
		isWalking = isOnFloor && isMovingHorizontal
		isIdle = !isAttacking && isOnFloor && !isMovingHorizontal
		isFalling = velocity.Y > MovingThreshold
		isSkiding = velocity.X < 0 != IsFacingLeft

		if IsFacingLeft {
			playerDrawOptions.FlipX = true
		} else {
			playerDrawOptions.FlipX = false
		}
		fsm.Update()
	}
}

var fsm = &FiniteStateMachine{
	Idle:      &Idle{},
	Jumping:   &Jumping{},
	Falling:   &Falling{},
	Walking:   &Walking{},
	Attacking: &Attacking{},
	Crouching: &Crouching{},
	Skidding:  &Skidding{},
}

// Finite State Machine for player
type FiniteStateMachine struct {
	Current State

	Idle      State
	Jumping   State
	Falling   State
	Walking   State
	Attacking State
	Crouching State
	Skidding  State
}

// SetState transitions to a new state.
func (f *FiniteStateMachine) SetState(s State) {
	if f.Current != nil {
		f.Current.Exit()
	}
	f.Current = s
	f.Current.Enter()
}

// Update the current state.
func (f *FiniteStateMachine) Update() {
	if f.Current != nil {
		f.Current.Update()
	}
}

// State interface to define required state functions.
type State interface {
	Enter()  // When entering a state
	Update() // Update logic for each state
	Exit()   // When exiting a state
}

type Idle struct{}

func (s *Idle) Enter() {
	fmt.Println("Entering Idle State")
}
func (s *Idle) Exit() {
}

func (s *Idle) Update() {
	playerAnim.SetState("idleRight")

	if isAttacking {
		fsm.SetState(fsm.Attacking)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && isOnFloor {
		fsm.SetState(fsm.Jumping)
	}
	if isWalking && isOnFloor {
		fsm.SetState(fsm.Walking)
	}

}

type Attacking struct{}

func (s *Attacking) Enter() {
	fmt.Println("Entering Attacking State")
}
func (s *Attacking) Exit() {
}

func (s *Attacking) Update() {
	// playerAnim.SetState("digRight")

	if !isAttacking && isIdle {
		fsm.SetState(fsm.Idle)
	}
	if !isAttacking && isWalking {
		fsm.SetState(fsm.Walking)
	}

}

type Jumping struct{}

func (s *Jumping) Enter() {
	fmt.Println("Entering Jumping State")
}
func (s *Jumping) Update() {

	playerAnim.SetState("jump")

	if isFalling {
		fsm.SetState(fsm.Falling)
	}
}
func (s *Jumping) Exit() {}

type Falling struct{}

func (s *Falling) Enter() {
	fmt.Println("Entering Falling State")
}
func (s *Falling) Update() {

	playerAnim.SetState("jump")

	if isOnFloor && isWalking {
		fsm.SetState(fsm.Walking)
	}
	if isOnFloor {
		fsm.SetState(fsm.Idle)
	}
}

func (s *Falling) Exit() {}

type Walking struct{}

func (s *Walking) Enter() {
	fmt.Println("Entering walking State")
}
func (s *Walking) Update() {
	fps := mathutil.MapRange(math.Abs(playerBody.Velocity().X), 0, 300, 0, 20)
	playerAnim.SetStateFPS("walkRight", fps)
	playerAnim.SetState("walkRight")

	if isSkiding {
		fsm.SetState(fsm.Skidding)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		fsm.SetState(fsm.Jumping)
	}
	if isIdle {
		fsm.SetState(fsm.Idle)
	}
}
func (s *Walking) Exit() {}

type Crouching struct{}

func (s *Crouching) Enter() {
	fmt.Println("Entering Crouching State")
}
func (s *Crouching) Update() {
}
func (s *Crouching) Exit() {}

type Skidding struct{}

func (s *Skidding) Enter() {
	fmt.Println("Entering Skidding State")
}
func (s *Skidding) Update() {

	playerAnim.SetState("skidding")

	if !isSkiding && isWalking {
		fsm.SetState(fsm.Walking)
	}

	if isIdle {
		fsm.SetState(fsm.Idle)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && isOnFloor {
		fsm.SetState(fsm.Jumping)
	}

}
func (s *Skidding) Exit() {}

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

func MoveToward(from, to, delta float64) float64 {
	if math.Abs(to-from) <= delta {
		return to
	}
	if to > from {
		return from + delta
	}
	return from - delta
}

func VelocityFunc(body *cm.Body, grav vec.Vec2, damping, dt float64) {
	velocity := body.Velocity()

	if isOnFloor {
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
