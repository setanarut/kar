package system

import (
	"kar/arche"
	"kar/comp"
	"kar/constants"
	"kar/engine"
	"kar/model"
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

	comp.Effect.Each(res.World, func(e *donburi.Entry) {
		effectData := comp.Effect.Get(e)
		charData := comp.Char.Get(e)

		if effectData.EffectTimerData.IsStart() {
			AddEffect(charData, effectData)

		}

		if effectData.EffectTimerData.IsReady() {
			RemoveEffect(charData, effectData)
			e.RemoveComponent(comp.Effect)
		}

		effectData.EffectTimerData.Update()
	})

	res.Input.UpdateArrowDirection()
	res.Input.UpdateWASDDirection()

	if player, ok := comp.PlayerTag.First(res.World); ok {

		playerCharData := comp.Char.Get(player)
		playerInventoryData := comp.Inventory.Get(player)
		playerBody := comp.Body.Get(player)
		playerRenderData := comp.Render.Get(player)

		if playerCharData.CurrentTool == constants.ItemSnowball {

			if !res.Input.ArrowDirection.Equal(engine.NoDirection) {
				playerRenderData.AnimPlayer.SetState("shoot")
				playerRenderData.DrawAngle = res.Input.ArrowDirection.ToAngle()

				// SHOOTING
				if playerInventoryData.Snowballs > 0 {

					if IsTimerReady(playerCharData.ShootCooldown) {
						ResetTimer(playerCharData.ShootCooldown)
					}

					if IsTimerStart(playerCharData.ShootCooldown) {
						for range playerCharData.SnowballPerCooldown {
							if playerInventoryData.Snowballs > 0 {
								playerInventoryData.Snowballs -= 1
							}
							// dir := engine.Rotate(res.Input.ArrowDirection.Mult(1000), engine.RandRange(0.2, -0.2))
							dir := res.Input.ArrowDirection.Mult(1000)
							// spawn snowball
							bullet := arche.SpawnDefaultSnowball(playerBody.Position())
							bulletBody := comp.Body.Get(bullet)

							bulletBody.ApplyImpulseAtWorldPoint(dir.Mult(bulletBody.Mass()), playerBody.Position())
						}
					}

				}

			} else {
				playerRenderData.AnimPlayer.SetState("right")
				playerRenderData.DrawAngle = res.Input.LastPressedDirection.ToAngle()

			}

		}

		if playerCharData.CurrentTool == constants.ItemBomb {
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
				if playerInventoryData.Bombs > 0 {
					arche.SpawnDefaultBomb(playerBody.Position().Add(res.Input.ArrowDirection.Mult(bombDistance)))
					playerInventoryData.Bombs -= 1
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			playerCharData.CurrentTool = constants.ItemSnowball
		}
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			playerCharData.CurrentTool = constants.ItemSnowball
		}

		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			playerCharData.CurrentTool = constants.ItemBomb
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			if playerInventoryData.Potion > 0 {
				if !player.HasComponent(comp.Effect) {
					player.AddComponent(comp.Effect)
					playerInventoryData.Potion -= 1
				}
			}
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

func AddEffect(charData *model.Mobile, effectData *model.EffectData) {
	charData.SnowballPerCooldown += effectData.ExtraSnowballPerAttack
	charData.ShootCooldown.Target += effectData.AdditiveShootCooldown
	charData.Speed += effectData.AdditiveMovementSpeed
}
func RemoveEffect(charData *model.Mobile, effectData *model.EffectData) {
	charData.SnowballPerCooldown -= effectData.ExtraSnowballPerAttack
	charData.ShootCooldown.Target -= effectData.AdditiveShootCooldown
	charData.Speed -= effectData.AdditiveMovementSpeed
}
