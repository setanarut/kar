package arche

import (
	"image"
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/res"
	"kar/types"

	"github.com/yohamta/donburi"
)

func SpawnDefaultPlayer(pos cm.Vec2) *donburi.Entry {
	return SpawnPlayer(1, 0, 0, 5, pos)

}
func SpawnBlock(pos cm.Vec2, chunkCoord image.Point, blockType types.BlockType) {
	e := SpawnWall(pos, 50, 50)
	e.AddComponent(comp.Health)
	e.AddComponent(comp.Block)
	e.AddComponent(comp.Sprite)

	comp.Block.Set(e,
		&types.DataBlock{
			ChunkCoord: chunkCoord,
			BlockType:  blockType,
		})

	sprite := comp.Sprite.Get(e)
	sprite.Image = res.StoneStages[0]
	sprite.Offset = engine.GetEbitenImageOffset(sprite.Image)
	sprite.DrawScale = engine.GetBoxScaleFactor(16, 16, 50, 50)
}
