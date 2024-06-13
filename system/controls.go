package system

import (
	"kar/comp"
	"kar/engine/cm"
	"kar/engine/vec"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

var HitShape *cm.Shape
var AttackSegmentQuery cm.SegmentQueryInfo

var (
	FacingLeft  bool
	FacingRight bool
	FacingDown  bool
	FacingUp    bool
	DigUp       bool
	IsGround    bool
	Idle        bool
	Walking     bool
	WalkRight   bool
	WalkLeft    bool
	Attacking   bool
	IdleAttack  bool
	NoWASD      bool
	DigDown     bool
)

type PlayerControlSystem struct {
}

func NewPlayerControlSystem() *PlayerControlSystem {
	return &PlayerControlSystem{}
}

func (sys *PlayerControlSystem) Init() {
}

func (sys *PlayerControlSystem) Update() {

	res.Input.UpdateWASDDirection()

	FacingRight = res.Input.LastPressedDirection.Equal(vec.Right) || res.Input.WASDDirection.Equal(vec.Right)
	FacingLeft = res.Input.LastPressedDirection.Equal(vec.Left) || res.Input.WASDDirection.Equal(vec.Left)
	FacingDown = res.Input.LastPressedDirection.Equal(vec.Down) || res.Input.WASDDirection.Equal(vec.Down)
	FacingUp = res.Input.LastPressedDirection.Equal(vec.Up) || res.Input.WASDDirection.Equal(vec.Up)
	NoWASD = res.Input.WASDDirection.Equal(vec.Zero)
	WalkRight = res.Input.WASDDirection.Equal(vec.Right)
	WalkLeft = res.Input.WASDDirection.Equal(vec.Left)
	Attacking = ebiten.IsKeyPressed(ebiten.KeyShiftRight)
	Walking = WalkLeft || WalkRight
	Idle = NoWASD && !Attacking && IsGround
	DigDown = FacingDown && Attacking
	DigUp = FacingUp && Attacking
	IdleAttack = NoWASD && Attacking && IsGround

	res.QueryWASDcontrollable.Each(res.World, WASDPlatformerForce)

	if player, ok := comp.PlayerTag.First(res.World); ok {

		playerBody := comp.Body.Get(player)
		playerAnimation := comp.AnimationPlayer.Get(player)
		playerDrawOptions := comp.DrawOptions.Get(player)
		p := playerBody.Position()
		AttackSegmentQuery = res.Space.SegmentQueryFirst(p, p.Add(res.Input.LastPressedDirection.Scale(50)), 0, res.FilterPlayerRaycast)

		// Reset block health
		if inpututil.IsKeyJustReleased(ebiten.KeyShiftRight) {
			if HitShape != nil {
				if CheckEntry(HitShape.Body()) {
					e := GetEntry(HitShape.Body())
					if e.HasComponent(comp.Block) && e.HasComponent(comp.Health) {
						ResetHealthComponent(e)
					}
				}
			}
		}

		// reset block health
		if AttackSegmentQuery.Shape == nil || AttackSegmentQuery.Shape != HitShape {
			if HitShape != nil {
				if CheckEntry(HitShape.Body()) {
					e := GetEntry(HitShape.Body())
					if e.HasComponent(comp.Block) && e.HasComponent(comp.Health) {
						ResetHealthComponent(e)
					}
				}
			}
		}

		// Attack
		if ebiten.IsKeyPressed(ebiten.KeyShiftRight) {
			if AttackSegmentQuery.Shape != nil && AttackSegmentQuery.Shape == HitShape {
				if HitShape != nil {
					if CheckEntry(HitShape.Body()) {
						e := GetEntry(HitShape.Body())
						if e.HasComponent(comp.Block) {
							h := comp.Health.Get(e)
							h.Health -= 0.3
						}
					}
				}
			}
		}

		if Idle {
			playerAnimation.SetState("idle")
		}

		if DigDown {
			playerAnimation.SetState("dig_down")
		}
		if DigUp {
			playerAnimation.SetState("dig_right")
		}

		if Attacking && FacingRight {
			playerAnimation.SetState("dig_right")
			playerDrawOptions.FlipX = false
		}
		if Attacking && FacingLeft {
			playerAnimation.SetState("dig_right")
			playerDrawOptions.FlipX = true
		}

		if WalkRight && !Attacking {
			playerAnimation.SetState("walk_right")
			playerDrawOptions.FlipX = false
		}
		if WalkLeft && !Attacking {
			playerAnimation.SetState("walk_right")
			playerDrawOptions.FlipX = true
		}

		HitShape = AttackSegmentQuery.Shape
	}

}

func (sys *PlayerControlSystem) Draw() {}

func WASDPlatformerForce(e *donburi.Entry) {

	body := comp.Body.Get(e)
	p := body.Position()
	// p.Add(vec.Vec2{0, -25}
	queryInfo := res.Space.SegmentQueryFirst(p, p.Add(vec.Vec2{0, -25}), 0, res.FilterPlayerRaycast)
	// queryInfoRight := res.Space.SegmentQueryFirst(p, p.Add(vec.Vec2{0, -25}), 0, res.FilterPlayerRaycast)
	contactShape := queryInfo.Shape

	bv := body.Velocity()
	body.SetVelocity(bv.X*0.9, bv.Y)
	// body.SetVelocityVector(body.Velocity().ClampLenght(500))
	if bv.X > 500 {
		body.SetVelocity(500, bv.Y)
	}
	if bv.X < -500 {
		body.SetVelocity(-500, bv.Y)
	}

	// yerde
	if contactShape != nil {
		IsGround = true
		// Zıpla
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			body.ApplyImpulseAtLocalPoint(vec.Vec2{0, 500}, body.CenterOfGravity())
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			body.ApplyForceAtLocalPoint(vec.Vec2{-1500, 0}, body.CenterOfGravity())
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			body.ApplyForceAtLocalPoint(vec.Vec2{1500, 0}, body.CenterOfGravity())
		}
	} else {
		IsGround = false
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			body.ApplyForceAtLocalPoint(vec.Vec2{-800, 0}, body.CenterOfGravity())
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			body.ApplyForceAtLocalPoint(vec.Vec2{800, 0}, body.CenterOfGravity())
		}
	}

}
func WASDPlatformer(e *donburi.Entry) {
	vel := vec.Vec2{}
	body := comp.Body.Get(e)
	p := body.Position()
	mobileData := comp.Mobile.Get(e)
	vel.X = res.Input.WASDDirection.X
	vel = vel.Scale(mobileData.Speed)
	vel = vel.Add(vec.Vec2{0, -500})
	body.SetVelocityVector(body.Velocity().LerpDistance(vel, mobileData.Accel))

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {

		queryInfo := res.Space.SegmentQueryFirst(p, p.Add(vec.Vec2{0, -50}), 0, res.FilterPlayerRaycast)
		contactShape := queryInfo.Shape

		if contactShape != nil {
			body.ApplyImpulseAtLocalPoint(vec.Vec2{0, 900}, body.CenterOfGravity())
		}

	}

}

func WASD4Directional(e *donburi.Entry) {
	body := comp.Body.Get(e)
	mobileData := comp.Mobile.Get(e)
	velocity := res.Input.WASDDirection.Normalize().Scale(mobileData.Speed)
	body.SetVelocityVector(body.Velocity().LerpDistance(velocity, mobileData.Accel))
}
