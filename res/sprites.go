package res

import (
	"kar/engine/util"
	"kar/items"
	"kar/types"

	"github.com/setanarut/anim"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	AtlasPlayer  = util.LoadEbitenImageFromFS(assets, "assets/player.png")
	AtlasBlock   = util.LoadEbitenImageFromFS(assets, "assets/blocks.png")
	SelectedBox  = util.LoadEbitenImageFromFS(assets, "assets/border48.png")
	SpriteFrames = make(map[types.ItemID][]*ebiten.Image)
)

func init() {
	SpriteFrames[items.Dirt] = anim.SubImages(AtlasBlock, 0, 0, 16, 16, 9, true)
	SpriteFrames[items.Sand] = anim.SubImages(AtlasBlock, 16, 0, 16, 16, 9, true)
	SpriteFrames[items.Stone] = anim.SubImages(AtlasBlock, 16*2, 0, 16, 16, 9, true)
	SpriteFrames[items.CoalOre] = anim.SubImages(AtlasBlock, 16*3, 0, 16, 16, 9, true)
	SpriteFrames[items.GoldOre] = anim.SubImages(AtlasBlock, 16*4, 0, 16, 16, 9, true)
	SpriteFrames[items.IronOre] = anim.SubImages(AtlasBlock, 16*5, 0, 16, 16, 9, true)
	SpriteFrames[items.DiamondOre] = anim.SubImages(AtlasBlock, 16*6, 0, 16, 16, 9, true)
	SpriteFrames[items.Grass] = anim.SubImages(AtlasBlock, 208, 0, 16, 16, 1, true)

	SpriteFrames[items.DeepSlateStone] = anim.SubImages(AtlasBlock, 16*2, 144, 16, 16, 9, true)
	SpriteFrames[items.DeepSlateCoalOre] = anim.SubImages(AtlasBlock, 16*3, 144, 16, 16, 9, true)
	SpriteFrames[items.DeepSlateGoldOre] = anim.SubImages(AtlasBlock, 16*4, 144, 16, 16, 9, true)
	SpriteFrames[items.DeepSlateIronOre] = anim.SubImages(AtlasBlock, 16*5, 144, 16, 16, 9, true)
	SpriteFrames[items.DeepSlateDiamondOre] = anim.SubImages(AtlasBlock, 16*6, 144, 16, 16, 9, true)

	SpriteFrames[items.RawCopper] = anim.SubImages(AtlasBlock, 16*2, 288, 16, 16, 1, false)
	SpriteFrames[items.Coal] = anim.SubImages(AtlasBlock, 16*3, 288, 16, 16, 1, false)
	SpriteFrames[items.RawGold] = anim.SubImages(AtlasBlock, 16*4, 288, 16, 16, 1, false)
	SpriteFrames[items.RawIron] = anim.SubImages(AtlasBlock, 16*5, 288, 16, 16, 1, false)
	SpriteFrames[items.Diamond] = anim.SubImages(AtlasBlock, 16*6, 288, 16, 16, 1, false)
}
