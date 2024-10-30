package system

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

// States
var (
	// isAttacking                   bool
	// isCrouching                   bool
	// IsFacingLeft, IsFacingRight   bool
	// isFacingUp, isFacingDown      bool
	// isFalling                     bool
	// isIdle                        bool
	// isRunning                     bool
	// isSkiding                     bool
	// isDigDown, isDigUp            bool
	// isWalkingLeft, isWalkingRight bool
	isWalking          bool
	isStanding         bool
	isMovingHorizontal bool
	isMovingVertical   bool
	isMoving           bool
	isRunning          bool
	// isJumping    bool
)
var velocity = vec.Vec2{}
var isOnFloor = true

const MovingThreshold float64 = 0.001

type States struct{}

func (sys *States) Init() {

	fsm.SetState(fsm.Idle)

}
func (sys *States) Draw() {}
func (sys *States) Update() {

	isRunning = ebiten.IsKeyPressed(ebiten.KeyAltRight)

	velocity = playerBody.Velocity()
	isOnFloor = OnFloor(playerBody)

	isMovingHorizontal = math.Abs(velocity.X) > MovingThreshold
	isMovingVertical = math.Abs(velocity.Y) > MovingThreshold
	isMoving = isMovingHorizontal || isMovingVertical

	isWalking = isOnFloor && isMoving
	isStanding = isOnFloor && !isMoving

	fsm.Update()
}

var fsm = &FiniteStateMachine{
	Idle:      &Idle{},
	Jumping:   &Jumping{},
	Falling:   &Falling{},
	Walking:   &Walking{},
	Crouching: &Crouching{},
}

// Finite State Machine for player
type FiniteStateMachine struct {
	Current State

	Idle      State
	Jumping   State
	Falling   State
	Walking   State
	Crouching State
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
	fmt.Println("Exit Idle State")
}

func (s *Idle) Update() {

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		fsm.SetState(fsm.Jumping)
	}

	if isOnFloor && math.Abs(velocity.X) > 0.001 {
		fsm.SetState(fsm.Walking)
	}

}

type Jumping struct{}

func (s *Jumping) Enter() {
	fmt.Println("Entering Jumping State")
}
func (s *Jumping) Update() {

	if velocity.Y > 0 {
		fsm.SetState(fsm.Falling)
	}
}
func (s *Jumping) Exit() {}

type Falling struct{}

func (s *Falling) Enter() {
	fmt.Println("Entering Falling State")
}
func (s *Falling) Update() {

	if isWalking {
		fsm.SetState(fsm.Walking)
	}
	if isStanding {
		fsm.SetState(fsm.Idle)
	}
}

func (s *Falling) Exit() {}

type Walking struct{}

func (s *Walking) Enter() {
	fmt.Println("Entering walking State")
}
func (s *Walking) Update() {

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		fsm.SetState(fsm.Jumping)
	}

	if isStanding {
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
