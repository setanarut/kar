package arche

import (
	"image"
	"kar/comp"
	"kar/engine/cm"
	"kar/res"
	"kar/types"

	"github.com/yohamta/donburi"
)

func SpawnDefaultPlayer(pos cm.Vec2) *donburi.Entry {
	return SpawnPlayer(1, 0, 0, 10, pos)

}
func SpawnBlock(pos cm.Vec2, chunkCoord image.Point) {
	e := SpawnWall(pos, 50, 50)
	e.AddComponent(comp.BlockTag)
	e.AddComponent(comp.Health)
	comp.ChunkCoord.Set(e, &types.DataBlock{ChunkCoord: chunkCoord})

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
