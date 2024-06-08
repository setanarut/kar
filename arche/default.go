package arche

import (
	"image"
	"kar/comp"
	"kar/engine/cm"
	"kar/types"

	"github.com/yohamta/donburi"
)

func SpawnDefaultPlayer(pos cm.Vec2) *donburi.Entry {
	return SpawnPlayer(1, 0, 0, 5, pos)

}
func SpawnBlock(pos cm.Vec2, chunkCoord image.Point) {
	e := SpawnWall(pos, 50, 50)
	e.AddComponent(comp.BlockTag)
	e.AddComponent(comp.Health)
	comp.ChunkCoord.Set(e, &types.DataBlock{ChunkCoord: chunkCoord})

}
