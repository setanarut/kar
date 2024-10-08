package res

import (
	"kar/engine/util"
	"kar/items"

	"github.com/setanarut/anim"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	AtlasPlayer  = util.LoadEbitenImageFromFS(assets, "assets/player.png")
	Border48     = util.LoadEbitenImageFromFS(assets, "assets/border48.png")
	Border32     = util.LoadEbitenImageFromFS(assets, "assets/border32.png")
	Slot16       = util.LoadEbitenImageFromFS(assets, "assets/slot16.png")
	SpriteFrames = make(map[uint16][]*ebiten.Image)
)

func init() {
	blockAtlas := util.LoadEbitenImageFromFS(assets, "assets/blocks.png")
	itemAtlas := util.LoadEbitenImageFromFS(assets, "assets/items.png")
	s := 16

	// blocks
	SpriteFrames[items.Grass] = anim.SubImages(blockAtlas, 0, 0, s, s, 11, false)
	SpriteFrames[items.Snow] = anim.SubImages(blockAtlas, 0, s, s, s, 11, false)
	SpriteFrames[items.Dirt] = anim.SubImages(blockAtlas, 0, 2*s, s, s, 11, false)
	SpriteFrames[items.Sand] = anim.SubImages(blockAtlas, 0, 3*s, s, s, 11, false)
	SpriteFrames[items.Stone] = anim.SubImages(blockAtlas, 0, 4*s, s, s, 11, false)
	SpriteFrames[items.CoalOre] = anim.SubImages(blockAtlas, 0, 5*s, s, s, 11, false)
	SpriteFrames[items.GoldOre] = anim.SubImages(blockAtlas, 0, 6*s, s, s, 11, false)
	SpriteFrames[items.IronOre] = anim.SubImages(blockAtlas, 0, 7*s, s, s, 11, false)
	SpriteFrames[items.DiamondOre] = anim.SubImages(blockAtlas, 0, 8*s, s, s, 11, false)
	SpriteFrames[items.CopperOre] = anim.SubImages(blockAtlas, 0, 9*s, s, s, 11, false)
	SpriteFrames[items.EmeraldOre] = anim.SubImages(blockAtlas, 0, 10*s, s, s, 11, false)
	SpriteFrames[items.LapisOre] = anim.SubImages(blockAtlas, 0, 11*s, s, s, 11, false)
	SpriteFrames[items.RedstoneOre] = anim.SubImages(blockAtlas, 0, 12*s, s, s, 11, false)
	SpriteFrames[items.DeepslateStone] = anim.SubImages(blockAtlas, 0, 13*s, s, s, 11, false)
	SpriteFrames[items.DeepslateCoalOre] = anim.SubImages(blockAtlas, 0, 14*s, s, s, 11, false)
	SpriteFrames[items.DeepslateGoldOre] = anim.SubImages(blockAtlas, 0, 15*s, s, s, 11, false)
	SpriteFrames[items.DeepslateIronOre] = anim.SubImages(blockAtlas, 0, 16*s, s, s, 11, false)
	SpriteFrames[items.DeepslateDiamondOre] = anim.SubImages(blockAtlas, 0, 17*s, s, s, 11, false)
	SpriteFrames[items.DeepslateCopperOre] = anim.SubImages(blockAtlas, 0, 18*s, s, s, 11, false)
	SpriteFrames[items.DeepslateEmeraldOre] = anim.SubImages(blockAtlas, 0, 19*s, s, s, 11, false)
	SpriteFrames[items.DeepslateLapisOre] = anim.SubImages(blockAtlas, 0, 20*s, s, s, 11, false)
	SpriteFrames[items.DeepslateRedStoneOre] = anim.SubImages(blockAtlas, 0, 21*s, s, s, 11, false)
	SpriteFrames[items.Tree] = anim.SubImages(blockAtlas, 0, 22*s, s, s, 11, false)
	SpriteFrames[items.TreeLeaves] = anim.SubImages(blockAtlas, 0, 23*s, s, s, 11, false)
	SpriteFrames[items.TreePlank] = anim.SubImages(blockAtlas, 0, 24*s, s, s, 11, false)

	SpriteFrames[items.Sapling] = anim.SubImages(itemAtlas, 0, 0, s, s, 1, false)
	SpriteFrames[items.Torch] = anim.SubImages(itemAtlas, s, 0, s, s, 1, false)
	SpriteFrames[items.Coal] = anim.SubImages(itemAtlas, 2*s, 0, s, s, 1, false)
	SpriteFrames[items.RawGold] = anim.SubImages(itemAtlas, 3*s, 0, s, s, 1, false)
	SpriteFrames[items.RawIron] = anim.SubImages(itemAtlas, 4*s, 0, s, s, 1, false)
	SpriteFrames[items.Diamond] = anim.SubImages(itemAtlas, 5*s, 0, s, s, 1, false)
	SpriteFrames[items.RawCopper] = anim.SubImages(itemAtlas, 6*s, 0, s, s, 1, false)
}
