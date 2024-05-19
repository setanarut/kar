package arche

import (
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"math/rand/v2"
	"time"

	"github.com/yohamta/donburi"
)

func FreezeEffect(speed float64) *comp.EffectData {
	return &comp.EffectData{
		ShootCooldown:    0,
		ExtraSnowball:    0,
		AddMovementSpeed: -speed,
		EffectTimer:      engine.NewTimer(time.Second * 3),
	}
}

func SpawnDefaultPlayer(pos cm.Vec2) *donburi.Entry {
	return SpawnPlayer(1, 0.3, 0.5, 22, pos)

}
func SpawnDefaultEnemy(pos cm.Vec2) *donburi.Entry {
	return SpawnEnemy(1, 0.3, 0.5, 20, pos)
}

func SpawnDefaultSnowball(pos cm.Vec2) *donburi.Entry {
	return SpawnSnowball(0.2, 0.3, 0.5, 7, pos)
}

func SpawnDefaultBomb(pos cm.Vec2) {
	SpawnBomb(1, 0.1, 0, 20, pos)
}

// collectible
func SpawnDefaultSnowballCollectible(pos cm.Vec2) {
	SpawnCollectible(comp.Snowball, 1, -1, 5, pos)
}
func SpawnDefaultEmeticCollectible(pos cm.Vec2) {
	SpawnCollectible(comp.PowerUpItem, 1, -1, 13, pos)
}

func SpawnDefaultKeyCollectible(keyNumber int, pos cm.Vec2) {
	SpawnCollectible(comp.Key, 1, keyNumber, 13, pos)
}

// Random

func SpawnRandomCollectible(pos cm.Vec2) {
	randomType := comp.ItemType(rand.IntN(4))
	SpawnCollectible(randomType, 1, engine.RandRangeInt(1, 10), 10, pos)
}

func SpawnRandomKeyCollectible(pos cm.Vec2) {
	SpawnCollectible(comp.Key, 1, engine.RandRangeInt(1, 10), 10, pos)
}

func SpawnRandomEnemy(pos cm.Vec2) {
	SpawnEnemy(1, 0.3, 0.5, engine.RandRange(5, 30), pos)
}
