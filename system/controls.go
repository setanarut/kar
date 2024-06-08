package system

import (
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

var currentBlock *donburi.Entry

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

		// playerAttackTimer := comp.AttackTimer.Get(player)
		// inventory := comp.Inventory.Get(player)
		playerBody := comp.Body.Get(player)

		// if res.CurrentTool == types.ItemSnowball {

		// 	if !res.Input.ArrowDirection.Equal(engine.NoDirection) {

		// 		// SHOOTING
		// 		if inventory.Items[types.ItemSnowball] > 0 {

		// 			if TimerIsReady(playerAttackTimer) {
		// 				TimerReset(playerAttackTimer)
		// 			}

		// 			if TimerIsStart(playerAttackTimer) {
		// 				dir := res.Input.ArrowDirection.Normalize().Mult(1000)
		// 				// spawn snowball
		// 				bullet := arche.SpawnDefaultSnowball(playerBody.Position())
		// 				inventory.Items[types.ItemSnowball] -= 1
		// 				bulletBody := comp.Body.Get(bullet)
		// 				bulletBody.ApplyImpulseAtWorldPoint(dir.Mult(bulletBody.Mass()), playerBody.Position())
		// 			}

		// 		}

		// 	}

		// }
		if ebiten.IsKeyPressed(ebiten.KeyShiftRight) {
			p := playerBody.Position()
			queryInfo := res.Space.SegmentQueryFirst(p, p.Add(res.Input.LastPressedDirection.Mult(50)), 0, res.FilterPlayerRaycast)
			contactShape := queryInfo.Shape
			if contactShape != nil {
				if CheckEntry(contactShape.Body()) {
					if currentBlock != nil {
						if currentBlock != GetEntry(contactShape.Body()) {
							comp.Health.SetValue(currentBlock, 3.0)
							s := comp.Sprite.Get(currentBlock)
							s.Image = res.StoneBlockImage
						}
					}
					currentBlock = GetEntry(contactShape.Body())
					if currentBlock.HasComponent(comp.BlockTag) && currentBlock.HasComponent(comp.Health) {
						h := comp.Health.GetValue(currentBlock)
						s := comp.Sprite.Get(currentBlock)
						i := int(engine.MapRange(h, 3, 0, 0, 7))
						s.Image = res.BlockBreakingStagesImages[i]
						comp.Health.SetValue(currentBlock, h-0.06)

					}
				}
			}

			if contactShape == nil {
				comp.Health.Each(res.World, ResetHealth)
				comp.BlockTag.Each(res.World, func(e *donburi.Entry) {
					s := comp.Sprite.Get(e)
					s.Image = res.StoneBlockImage
				})
			}
		}

		if inpututil.IsKeyJustReleased(ebiten.KeyShiftRight) {
			comp.Health.Each(res.World, ResetHealth)
			comp.BlockTag.Each(res.World, func(e *donburi.Entry) {
				s := comp.Sprite.Get(e)
				s.Image = res.StoneBlockImage
			})
		}
	}
}

func (sys *PlayerControlSystem) Draw() {
}

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
	vel = vel.Mult(mobileData.Speed)
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
	velocity := res.Input.WASDDirection.Normalize().Mult(mobileData.Speed)
	body.SetVelocityVector(body.Velocity().LerpDistance(velocity, mobileData.Accel))
}
