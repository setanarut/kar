package system

import (
	"kar/arche"
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/res"
	"kar/types"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

var bombDistance float64 = 40

type PlayerControlSystem struct {
}

func NewPlayerControlSystem() *PlayerControlSystem {
	return &PlayerControlSystem{}
}

func (sys *PlayerControlSystem) Init() {

}

func (sys *PlayerControlSystem) Update() {

	res.Input.UpdateArrowDirection()
	res.Input.UpdateWASDDirection()

	// WASD system
	res.QueryWASDcontrollable.Each(res.World, WASDMovementForce)
	// res.QueryWASDcontrollable.Each(res.World, WASDMovementVel)

	if player, ok := comp.PlayerTag.First(res.World); ok {

		playerAttackTimer := comp.AttackTimer.Get(player)
		inventory := comp.Inventory.Get(player)
		playerBody := comp.Body.Get(player)
		playerRender := comp.Render.Get(player)

		playerRender.AnimPlayer.SetState("right")
		playerRender.DrawAngle = res.Input.LastPressedDirection.ToAngle()

		if res.CurrentTool == types.ItemSnowball {

			if !res.Input.ArrowDirection.Equal(engine.NoDirection) {

				playerRender.AnimPlayer.SetState("shoot")
				playerRender.DrawAngle = res.Input.ArrowDirection.ToAngle()

				// SHOOTING
				if inventory.Items[types.ItemSnowball] > 0 {

					if timerIsReady(playerAttackTimer) {
						timerReset(playerAttackTimer)
					}

					if timerIsStart(playerAttackTimer) {
						// dir := engine.Rotate(res.Input.ArrowDirection.Mult(1000), engine.RandRange(0.2, -0.2))
						dir := res.Input.ArrowDirection.Normalize().Mult(1000)
						// spawn snowball
						bullet := arche.SpawnDefaultSnowball(playerBody.Position())
						inventory.Items[types.ItemSnowball] -= 1
						bulletBody := comp.Body.Get(bullet)
						bulletBody.ApplyImpulseAtWorldPoint(dir.Mult(bulletBody.Mass()), playerBody.Position())
					}

				}

			}

		}

		if res.CurrentTool == types.ItemBomb {
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
				if inventory.Items[types.ItemBomb] > 0 {
					arche.SpawnDefaultBomb(playerBody.Position().Add(res.Input.ArrowDirection.Mult(bombDistance)))
					inventory.Items[types.ItemBomb] -= 1
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyU) {
			comp.AI.Each(res.World, func(e *donburi.Entry) {
				e.RemoveComponent(comp.AI)
				e.AddComponent(comp.WASDTag)
			})
		}
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			res.CurrentTool = types.ItemSnowball
		}

		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			res.CurrentTool = types.ItemBomb
		}

	} // bitiş

	// Explode all bombs
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		comp.BombTag.Each(res.World, func(e *donburi.Entry) {
			explode(e)
		})
	}

	// AI on/off
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		comp.AI.Each(res.World, func(e *donburi.Entry) {
			ai := comp.AI.Get(e)
			ai.Follow = !ai.Follow
		})

	}

}

func (sys *PlayerControlSystem) Draw() {
}

/* func AddEffect(charData *model.Mobile, effectData *model.EffectData) {
	charData.SnowballPerCooldown += effectData.ExtraSnowballPerAttack
	charData.ShootCooldown.Target += effectData.AdditiveShootCooldown
	charData.Speed += effectData.AdditiveMovementSpeed
}
func RemoveEffect(charData *model.Mobile, effectData *model.EffectData) {
	charData.SnowballPerCooldown -= effectData.ExtraSnowballPerAttack
	charData.ShootCooldown.Target -= effectData.AdditiveShootCooldown
	charData.Speed -= effectData.AdditiveMovementSpeed
}
*/

func WASDMovementForce(e *donburi.Entry) {
	vel := cm.Vec2{}
	body := comp.Body.Get(e)
	mobileData := comp.Mobile.Get(e)

	vel.X = res.Input.WASDDirection.X
	vel = vel.Mult(mobileData.Speed)
	bv := body.Velocity()
	body.SetVelocity(bv.X*0.9, bv.Y)
	if bv.X > 500 {
		body.SetVelocity(500, bv.Y)
	}
	if bv.X < -500 {
		body.SetVelocity(-500, bv.Y)
	}

	body.ApplyForceAtLocalPoint(vel, body.CenterOfGravity())

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		body.ApplyImpulseAtLocalPoint(cm.Vec2{0, 400}, body.CenterOfGravity())
	}

}
func WASDMovementVel(e *donburi.Entry) {
	vel := cm.Vec2{}
	body := comp.Body.Get(e)
	mobileData := comp.Mobile.Get(e)
	vel.X = res.Input.WASDDirection.X
	vel = vel.Mult(mobileData.Speed)
	vel = vel.Add(cm.Vec2{0, -500})
	body.SetVelocityVector(body.Velocity().LerpDistance(vel, mobileData.Accel))

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		body.ApplyImpulseAtLocalPoint(cm.Vec2{0, 1200}, body.CenterOfGravity())
	}

}
