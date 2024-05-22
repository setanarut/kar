package arche

import (
	"kar/constants"
	"kar/engine"
	"kar/engine/cm"
	"kar/model"
	"math/rand/v2"
	"time"

	"github.com/yohamta/donburi"
)

func PotionFreeze(speed float64) *model.EffectData {
	return &model.EffectData{
		AdditiveShootCooldown:  0,
		ExtraSnowballPerAttack: 0,
		AdditiveMovementSpeed:  -speed,
		EffectTimerData:        engine.NewTimer(time.Second * 3),
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
	SpawnCollectible(constants.ItemSnowball, 1, -1, 5, pos)
}
func SpawnDefaultEmeticCollectible(pos cm.Vec2) {
	SpawnCollectible(constants.ItemPotion, 1, -1, 13, pos)
}

func SpawnDefaultKeyCollectible(keyNumber int, pos cm.Vec2) {
	SpawnCollectible(constants.ItemKey, 1, keyNumber, 13, pos)
}

// Random

func SpawnRandomCollectible(pos cm.Vec2) {
	randomType := model.ItemType(rand.IntN(4))
	SpawnCollectible(randomType, 1, engine.RandRangeInt(1, 10), 10, pos)
}

func SpawnRandomKeyCollectible(pos cm.Vec2) {
	SpawnCollectible(constants.ItemKey, 1, engine.RandRangeInt(1, 10), 10, pos)
}

func SpawnRandomEnemy(pos cm.Vec2) {
	SpawnEnemy(1, 0.3, 0.5, engine.RandRange(5, 30), pos)
}
