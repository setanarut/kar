package system

import (
	"kar/arche"
	"kar/comp"
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

	// Input.UpdateJustArrowDirection()
	res.Input.UpdateArrowDirection()
	res.Input.UpdateWASDDirection()

	if playerEntry, ok := comp.PlayerTag.First(res.World); ok {

		charData := comp.Char.Get(playerEntry)
		inventory := comp.Inventory.Get(playerEntry)
		playerBody := comp.Body.Get(playerEntry)
		playerRenderData := comp.Render.Get(playerEntry)
		playerPos := playerBody.Position()

		// playerBody.ApplyForceAtLocalPoint(res.Input.WASDDirection.Normalize().Mult(1800), playerBody.CenterOfGravity())

		if playerEntry.HasComponent(comp.PowerUp) {

			drugEffectData := comp.PowerUp.Get(playerEntry)

			if drugEffectData.EffectTimer.IsStart() {
				AddDrugEffect(charData, drugEffectData)

			}

			if drugEffectData.EffectTimer.IsReady() {
				RemoveDrugEffect(charData, drugEffectData)
				playerEntry.RemoveComponent(comp.PowerUp)
			}

			drugEffectData.EffectTimer.Update()
		}

		if !res.Input.ArrowDirection.Equal(engine.NoDirection) {
			playerRenderData.AnimPlayer.SetState("shoot")
			playerRenderData.DrawAngle = res.Input.ArrowDirection.ToAngle()

			// SHOOTING
			if inventory.Snowballs > 0 {
				if charData.ShootCooldownTimer.IsReady() {
					charData.ShootCooldownTimer.Reset()
				}

				if charData.ShootCooldownTimer.IsStart() {
					for range charData.SnowballPerCooldown {
						if inventory.Snowballs > 0 {
							inventory.Snowballs -= 1
						}
						dir := engine.Rotate(res.Input.ArrowDirection.Mult(1000), engine.RandRange(0.2, -0.2))
						bullet := arche.SpawnDefaultSnowball(playerPos)
						bulletBody := comp.Body.Get(bullet)
						bulletBody.ApplyImpulseAtWorldPoint(dir, playerPos)
					}
				}
			}

		} else {
			playerRenderData.AnimPlayer.SetState("right")
			playerRenderData.DrawAngle = res.Input.LastPressedDirection.ToAngle()

		}

		charData.ShootCooldownTimer.Update()

		if inventory.Bombs > 0 {

			// Bomba bırak
			if inpututil.IsKeyJustPressed(ebiten.KeyShiftRight) {
				bombPos := res.Input.LastPressedDirection.Neg().Mult(bombDistance)
				arche.SpawnDefaultBomb(playerPos.Add(bombPos))
				inventory.Bombs -= 1
			}

		}

		// ilac kullan
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {

			if inventory.PowerUp > 0 {
				if !playerEntry.HasComponent(comp.PowerUp) {
					playerEntry.AddComponent(comp.PowerUp)
					inventory.PowerUp -= 1
				}
			}
		}

	}

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

func AddDrugEffect(charData *comp.CharacterData, drugEffectData *comp.PowerUpData) {
	charData.SnowballPerCooldown += drugEffectData.ExtraSnowball
	charData.ShootCooldownTimer.Target += drugEffectData.ShootCooldown
	charData.Speed += drugEffectData.AddMovementSpeed
}
func RemoveDrugEffect(charData *comp.CharacterData, drugEffectData *comp.PowerUpData) {
	charData.SnowballPerCooldown -= drugEffectData.ExtraSnowball
	charData.ShootCooldownTimer.Target -= drugEffectData.ShootCooldown
	charData.Speed -= drugEffectData.AddMovementSpeed
}
