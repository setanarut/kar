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

		if effectData.EffectTimer.IsStart() {
			AddDrugEffect(charData, effectData)

		}

		if effectData.EffectTimer.IsReady() {
			RemoveDrugEffect(charData, effectData)
			e.RemoveComponent(comp.Effect)
		}

		effectData.EffectTimer.Update()
	})

	res.Input.UpdateArrowDirection()
	res.Input.UpdateWASDDirection()

	if player, ok := comp.PlayerTag.First(res.World); ok {

		playerCharData := comp.Char.Get(player)

		if playerCharData.CurrentTool == constants.ItemSnowball {

			playerInventoryData := comp.Inventory.Get(player)
			playerBody := comp.Body.Get(player)
			playerRenderData := comp.Render.Get(player)

			if !res.Input.ArrowDirection.Equal(engine.NoDirection) {
				playerRenderData.AnimPlayer.SetState("shoot")
				playerRenderData.DrawAngle = res.Input.ArrowDirection.ToAngle()

				// SHOOTING
				if playerInventoryData.Snowballs > 0 {
					if playerCharData.ShootCooldownTimer.IsReady() {
						playerCharData.ShootCooldownTimer.Reset()
					}

					if playerCharData.ShootCooldownTimer.IsStart() {
						for range playerCharData.SnowballPerCooldown {
							if playerInventoryData.Snowballs > 0 {
								playerInventoryData.Snowballs -= 1
							}
							dir := engine.Rotate(res.Input.ArrowDirection.Mult(1000), engine.RandRange(0.2, -0.2))
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

			playerCharData.ShootCooldownTimer.Update()

			if playerInventoryData.Bombs > 0 {

				// Bomba bırak
				if inpututil.IsKeyJustPressed(ebiten.KeyShiftRight) {
					bombPos := res.Input.LastPressedDirection.Neg().Mult(bombDistance)
					arche.SpawnDefaultBomb(playerBody.Position().Add(bombPos))
					playerInventoryData.Bombs -= 1
				}

			}

			// ilac kullan
			if inpututil.IsKeyJustPressed(ebiten.KeySpace) {

				if playerInventoryData.Potion > 0 {
					if !player.HasComponent(comp.Effect) {
						player.AddComponent(comp.Effect)
						playerInventoryData.Potion -= 1
					}
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

func AddDrugEffect(charData *model.CharacterData, effectData *model.EffectData) {
	charData.SnowballPerCooldown += effectData.ExtraSnowball
	charData.ShootCooldownTimer.Target += effectData.ShootCooldown
	charData.Speed += effectData.AddMovementSpeed
}
func RemoveDrugEffect(charData *model.CharacterData, effectData *model.EffectData) {
	charData.SnowballPerCooldown -= effectData.ExtraSnowball
	charData.ShootCooldownTimer.Target -= effectData.ShootCooldown
	charData.Speed -= effectData.AddMovementSpeed
}
