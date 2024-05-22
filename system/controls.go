package system

import (
	"kar/arche"
	"kar/comp"
	"kar/constants"
	"kar/engine"
	"kar/res"

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

	if player, ok := comp.PlayerTag.First(res.World); ok {

		playerAttackTimer := comp.AttackTimer.Get(player)
		inventory := comp.Inventory.Get(player)
		playerBody := comp.Body.Get(player)
		playerRender := comp.Render.Get(player)

		playerRender.AnimPlayer.SetState("right")
		playerRender.DrawAngle = res.Input.LastPressedDirection.ToAngle()

		if res.CurrentTool == constants.ItemSnowball {

			if !res.Input.ArrowDirection.Equal(engine.NoDirection) {

				playerRender.AnimPlayer.SetState("shoot")
				playerRender.DrawAngle = res.Input.ArrowDirection.ToAngle()

				// SHOOTING
				if inventory.Items[constants.ItemSnowball] > 0 {

					if IsTimerReady(playerAttackTimer) {
						ResetTimer(playerAttackTimer)
					}

					if IsTimerStart(playerAttackTimer) {
						// dir := engine.Rotate(res.Input.ArrowDirection.Mult(1000), engine.RandRange(0.2, -0.2))
						dir := res.Input.ArrowDirection.Mult(1000)
						// spawn snowball
						bullet := arche.SpawnDefaultSnowball(playerBody.Position())
						bulletBody := comp.Body.Get(bullet)

						bulletBody.ApplyImpulseAtWorldPoint(dir.Mult(bulletBody.Mass()), playerBody.Position())
					}

				}

			}

		}

		if res.CurrentTool == constants.ItemBomb {
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
				if inventory.Items[constants.ItemBomb] > 0 {
					arche.SpawnDefaultBomb(playerBody.Position().Add(res.Input.ArrowDirection.Mult(bombDistance)))
					inventory.Items[constants.ItemBomb] -= 1
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			res.CurrentTool = constants.ItemSnowball
		}
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			res.CurrentTool = constants.ItemSnowball
		}

		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			res.CurrentTool = constants.ItemBomb
		}

	} // bitiş

	// Explode all bombs
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		comp.BombTag.Each(res.World, func(e *donburi.Entry) {
			Explode(e)
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
