package res

import (
	"kar/engine/io"
	"kar/engine/util"
	"kar/items"
	"kar/types"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	AtlasPlayer  = io.LoadEbitenImageFromFS(assets, "assets/player.png")
	AtlasBlock   = io.LoadEbitenImageFromFS(assets, "assets/blocks.png")
	SpriteFrames = make(map[types.ItemType][]*ebiten.Image)
)

func init() {
	SpriteFrames[items.Dirt] = util.SubImages(AtlasBlock, 0, 0, 16, 16, 9, true)
	SpriteFrames[items.Sand] = util.SubImages(AtlasBlock, 16, 0, 16, 16, 9, true)
	SpriteFrames[items.Stone] = util.SubImages(AtlasBlock, 16*2, 0, 16, 16, 9, true)
	SpriteFrames[items.CoalOre] = util.SubImages(AtlasBlock, 16*3, 0, 16, 16, 9, true)
	SpriteFrames[items.GoldOre] = util.SubImages(AtlasBlock, 16*4, 0, 16, 16, 9, true)
	SpriteFrames[items.IronOre] = util.SubImages(AtlasBlock, 16*5, 0, 16, 16, 9, true)
	SpriteFrames[items.DiamondOre] = util.SubImages(AtlasBlock, 16*6, 0, 16, 16, 9, true)

	SpriteFrames[items.DeepSlateStone] = util.SubImages(AtlasBlock, 16*2, 144, 16, 16, 9, true)
	SpriteFrames[items.DeepSlateCoalOre] = util.SubImages(AtlasBlock, 16*3, 144, 16, 16, 9, true)
	SpriteFrames[items.DeepSlateGoldOre] = util.SubImages(AtlasBlock, 16*4, 144, 16, 16, 9, true)
	SpriteFrames[items.DeepSlateIronOre] = util.SubImages(AtlasBlock, 16*5, 144, 16, 16, 9, true)
	SpriteFrames[items.DeepSlateDiamondOre] = util.SubImages(AtlasBlock, 16*6, 144, 16, 16, 9, true)

	SpriteFrames[items.RawCopper] = util.SubImages(AtlasBlock, 16*2, 288, 16, 16, 1, false)
	SpriteFrames[items.Coal] = util.SubImages(AtlasBlock, 16*3, 288, 16, 16, 1, false)
	SpriteFrames[items.RawGold] = util.SubImages(AtlasBlock, 16*4, 288, 16, 16, 1, false)
	SpriteFrames[items.RawIron] = util.SubImages(AtlasBlock, 16*5, 288, 16, 16, 1, false)
	SpriteFrames[items.Diamond] = util.SubImages(AtlasBlock, 16*6, 288, 16, 16, 1, false)
}
