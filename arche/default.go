package arche

import (
	"kar/comp"
	"kar/engine/cm"
	"kar/res"

	"github.com/yohamta/donburi"
)

func SpawnDefaultPlayer(pos cm.Vec2) *donburi.Entry {
	return SpawnPlayer(1, 0, 0, 24, pos)

}
func SpawnDefaultMob(pos cm.Vec2) {
	mob := SpawnMob(1, 0.3, 0.5, 20, pos)
	if p, ok := comp.PlayerTag.First(res.World); ok {
		ai := comp.AI.Get(mob)
		ai.Target = p
	}
}

func SpawnDefaultSnowball(pos cm.Vec2) *donburi.Entry {
	return SpawnSnowball(0.2, 0.3, 0.5, 7, pos)
}

func SpawnDefaultBomb(pos cm.Vec2) {
	SpawnBomb(1, 0.1, 0, 20, pos)
}
func SpawnBlock(pos cm.Vec2) {
	SpawnWall(pos, 50, 50)
}

// collectible
// func SpawnDefaultSnowballCollectible(pos cm.Vec2) {
// 	SpawnCollectible(res.ItemSnowball, 1, -1, 5, pos)
// }
// func SpawnDefaultEmeticCollectible(pos cm.Vec2) {
// 	SpawnCollectible(res.ItemPotion, 1, -1, 13, pos)
// }

// func SpawnDefaultKeyCollectible(keyNumber int, pos cm.Vec2) {
// 	SpawnCollectible(res.ItemKey, 1, keyNumber, 13, pos)
// }

// Random

// func SpawnRandomCollectible(pos cm.Vec2) {
// 	randomType := model.ItemType(rand.IntN(4))
// 	SpawnCollectible(randomType, 1, engine.RandRangeInt(1, 10), 10, pos)
// }

// func SpawnRandomKeyCollectible(pos cm.Vec2) {
// 	SpawnCollectible(res.ItemKey, 1, engine.RandRangeInt(1, 10), 10, pos)
// }

// func SpawnRandomEnemy(pos cm.Vec2) {
// 	SpawnEnemy(1, 0.3, 0.5, engine.RandRange(5, 30), pos)
// }
