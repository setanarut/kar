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
var RenderScale cm.Vec2

type PlayerControlSystem struct {
}

func NewPlayerControlSystem() *PlayerControlSystem {
	return &PlayerControlSystem{}
}

func (sys *PlayerControlSystem) Init() {
}

func (sys *PlayerControlSystem) Update() {

	res.Input.UpdateWASDDirection()
	// res.Input.UpdateArrowDirection()
	// res.QueryWASDcontrollable.Each(res.World, WASD4Directional)
	res.QueryWASDcontrollable.Each(res.World, WASDPlatformerForce)
	// res.QueryWASDcontrollable.Each(res.World, WASDPlatformer)

	if player, ok := comp.PlayerTag.First(res.World); ok {
		playerBody := comp.Body.Get(player)
		r := comp.Render.Get(player)
		p := playerBody.Position()

		queryInfo := res.Space.SegmentQueryFirst(p, p.Add(res.Input.LastPressedDirection.Scale(50)), 0, res.FilterPlayerRaycast)

		// hedef değişti mi?
		if queryInfo.Shape == nil || queryInfo.Shape != HitShape {
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

		if ebiten.IsKeyPressed(ebiten.KeyShiftRight) {
			if queryInfo.Shape != nil && queryInfo.Shape == HitShape {
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

		if ebiten.IsKeyPressed(ebiten.KeyA) {
			r.CurrentScale = r.DrawScaleFlipX
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			r.CurrentScale = r.DrawScale
		}

		HitShape = queryInfo.Shape
	}

}

func (sys *PlayerControlSystem) Draw() {}

func WASDPlatformerForce(e *donburi.Entry) {

	body := comp.Body.Get(e)
	p := body.Position()

	queryInfo := res.Space.SegmentQueryFirst(p, p.Add(cm.Vec2{0, -30}), 0, res.FilterPlayerRaycast)
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
