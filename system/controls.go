package system

import (
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/res"
	"kar/types"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

var HitShape *cm.Shape
var HitBlockHealth float64
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

	FacingRight = res.Input.LastPressedDirection.Equal(engine.RightDirection) || res.Input.WASDDirection.Equal(engine.RightDirection)
	FacingLeft = res.Input.LastPressedDirection.Equal(engine.LeftDirection) || res.Input.WASDDirection.Equal(engine.LeftDirection)
	FacingDown = res.Input.LastPressedDirection.Equal(engine.DownDirection) || res.Input.WASDDirection.Equal(engine.DownDirection)
	FacingUp = res.Input.LastPressedDirection.Equal(engine.UpDirection) || res.Input.WASDDirection.Equal(engine.UpDirection)
	NoWASD = res.Input.WASDDirection.Equal(engine.NoDirection)
	WalkRight = res.Input.WASDDirection.Equal(engine.RightDirection)
	WalkLeft = res.Input.WASDDirection.Equal(engine.LeftDirection)
	Attacking = ebiten.IsKeyPressed(ebiten.KeyShiftRight)
	Walking = WalkLeft || WalkRight
	Idle = NoWASD && !Attacking && IsGround
	DigDown = FacingDown && Attacking
	DigUp = FacingUp && Attacking
	IdleAttack = NoWASD && Attacking && IsGround

	res.QueryWASDcontrollable.Each(res.World, WASDPlatformerForce)

	if player, ok := comp.PlayerTag.First(res.World); ok {

		playerBody := comp.Body.Get(player)
		r := comp.Render.Get(player)
		p := playerBody.Position()

		AttackSegmentQuery = res.Space.SegmentQueryFirst(p, p.Add(res.Input.LastPressedDirection.Scale(50)), 0, res.FilterPlayerRaycast)

		// hedef değişti mi?
		if AttackSegmentQuery.Shape == nil || AttackSegmentQuery.Shape != HitShape {
			if HitShape != nil {
				if CheckEntry(HitShape.Body()) {
					e := GetEntry(HitShape.Body())
					if e.HasComponent(comp.Block) {
						comp.Health.SetValue(e, 3.0)
						s := comp.Sprite.Get(e)
						b := comp.Block.Get(e)
						if b.BlockType == types.BlockStone {
							s.Image = res.StoneStages[0]
						}
					}
				}
			}
		}

		// kaz
		if ebiten.IsKeyPressed(ebiten.KeyShiftRight) {
			if AttackSegmentQuery.Shape != nil && AttackSegmentQuery.Shape == HitShape {
				if HitShape != nil {
					if CheckEntry(HitShape.Body()) {
						e := GetEntry(HitShape.Body())
						if e.HasComponent(comp.Block) {
							h := comp.Health.GetValue(e)
							b := comp.Block.Get(e)
							s := comp.Sprite.Get(e)
							if b.BlockType == types.BlockStone {
								s.Image = res.StoneStages[int(engine.MapRange(h, 3, 0, 0, 8))]
							}

							HitBlockHealth = h
							comp.Health.SetValue(e, h-0.06)
						}
					}
				}
			}
		}

		if inpututil.IsKeyJustReleased(ebiten.KeyShiftRight) {
			if HitShape != nil {
				if CheckEntry(HitShape.Body()) {
					e := GetEntry(HitShape.Body())
					if e.HasComponent(comp.Block) {
						comp.Health.SetValue(e, 3.0)
						s := comp.Sprite.Get(e)
						b := comp.Block.Get(e)
						if b.BlockType == types.BlockStone {
							s.Image = res.StoneStages[0]
						}
					}
				}
			}
		}

		if Idle {
			r.AnimPlayer.SetState("idle")
		}

		if DigDown {
			r.AnimPlayer.SetState("dig_down")
		}
		if DigUp {
			r.AnimPlayer.SetState("dig_right")
		}

		if Attacking && FacingRight {
			r.AnimPlayer.SetState("dig_right")
			r.CurrentScale = r.DrawScale
		}
		if Attacking && FacingLeft {
			r.AnimPlayer.SetState("dig_right")
			r.CurrentScale = r.DrawScaleFlipX
		}

		if WalkRight && !Attacking {
			r.AnimPlayer.SetState("walk_right")
			r.CurrentScale = r.DrawScale
		}
		if WalkLeft && !Attacking {
			r.AnimPlayer.SetState("walk_right")
			r.CurrentScale = r.DrawScaleFlipX
		}

		HitShape = AttackSegmentQuery.Shape
	}

}

func (sys *PlayerControlSystem) Draw() {}

func WASDPlatformerForce(e *donburi.Entry) {

	body := comp.Body.Get(e)
	p := body.Position()
	// p.Add(cm.Vec2{0, -25}
	queryInfo := res.Space.SegmentQueryFirst(p, p.Add(cm.Vec2{0, -25}), 0, res.FilterPlayerRaycast)
	// queryInfoRight := res.Space.SegmentQueryFirst(p, p.Add(cm.Vec2{0, -25}), 0, res.FilterPlayerRaycast)
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
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			body.ApplyForceAtLocalPoint(cm.Vec2{-1500, 0}, body.CenterOfGravity())
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			body.ApplyForceAtLocalPoint(cm.Vec2{1500, 0}, body.CenterOfGravity())
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			body.ApplyImpulseAtLocalPoint(cm.Vec2{0, 500}, body.CenterOfGravity())
		}
	} else {
		IsGround = false
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			body.ApplyForceAtLocalPoint(cm.Vec2{-800, 0}, body.CenterOfGravity())
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			body.ApplyForceAtLocalPoint(cm.Vec2{800, 0}, body.CenterOfGravity())
		}
	}

}
func WASDPlatformer(e *donburi.Entry) {
	vel := cm.Vec2{}
	body := comp.Body.Get(e)
	p := body.Position()
	mobileData := comp.Mobile.Get(e)
	vel.X = res.Input.WASDDirection.X
	vel = vel.Scale(mobileData.Speed)
	vel = vel.Add(cm.Vec2{0, -500})
	body.SetVelocityVector(body.Velocity().LerpDistance(vel, mobileData.Accel))

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {

		queryInfo := res.Space.SegmentQueryFirst(p, p.Add(cm.Vec2{0, -50}), 0, res.FilterPlayerRaycast)
		contactShape := queryInfo.Shape

		if contactShape != nil {
			body.ApplyImpulseAtLocalPoint(cm.Vec2{0, 900}, body.CenterOfGravity())
		}

	}

}

func WASD4Directional(e *donburi.Entry) {
	body := comp.Body.Get(e)
	mobileData := comp.Mobile.Get(e)
	velocity := res.Input.WASDDirection.Normalize().Scale(mobileData.Speed)
	body.SetVelocityVector(body.Velocity().LerpDistance(velocity, mobileData.Accel))
}
