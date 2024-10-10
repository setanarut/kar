package res

import (
	"kar/engine/util"
	"kar/itm"

	"github.com/setanarut/anim"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	AtlasPlayer     = util.LoadEbitenImageFromFS(assets, "assets/player.png")
	Border          = util.LoadEbitenImageFromFS(assets, "assets/border64.png")
	Hotbar          = util.LoadEbitenImageFromFS(assets, "assets/hotbar.png")
	HotbarSelection = util.LoadEbitenImageFromFS(assets, "assets/hotbarselect.png")
	SpriteFrames    = make(map[uint16][]*ebiten.Image)
)

func init() {
	blockAtlas := util.LoadEbitenImageFromFS(assets, "assets/blocks.png")
	itemAtlas := util.LoadEbitenImageFromFS(assets, "assets/items.png")
	s := 16

	// blocks
	SpriteFrames[itm.Bedrock] = anim.SubImages(blockAtlas, 176, 0, s, s, 1, false)
	SpriteFrames[itm.Grass] = anim.SubImages(blockAtlas, 0, 0, s, s, 11, false)
	SpriteFrames[itm.Snow] = anim.SubImages(blockAtlas, 0, s, s, s, 11, false)
	SpriteFrames[itm.Dirt] = anim.SubImages(blockAtlas, 0, 2*s, s, s, 11, false)
	SpriteFrames[itm.Sand] = anim.SubImages(blockAtlas, 0, 3*s, s, s, 11, false)
	SpriteFrames[itm.Stone] = anim.SubImages(blockAtlas, 0, 4*s, s, s, 11, false)
	SpriteFrames[itm.CoalOre] = anim.SubImages(blockAtlas, 0, 5*s, s, s, 11, false)
	SpriteFrames[itm.GoldOre] = anim.SubImages(blockAtlas, 0, 6*s, s, s, 11, false)
	SpriteFrames[itm.IronOre] = anim.SubImages(blockAtlas, 0, 7*s, s, s, 11, false)
	SpriteFrames[itm.DiamondOre] = anim.SubImages(blockAtlas, 0, 8*s, s, s, 11, false)
	SpriteFrames[itm.CopperOre] = anim.SubImages(blockAtlas, 0, 9*s, s, s, 11, false)
	SpriteFrames[itm.EmeraldOre] = anim.SubImages(blockAtlas, 0, 10*s, s, s, 11, false)
	SpriteFrames[itm.LapisOre] = anim.SubImages(blockAtlas, 0, 11*s, s, s, 11, false)
	SpriteFrames[itm.RedstoneOre] = anim.SubImages(blockAtlas, 0, 12*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateStone] = anim.SubImages(blockAtlas, 0, 13*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateCoalOre] = anim.SubImages(blockAtlas, 0, 14*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateGoldOre] = anim.SubImages(blockAtlas, 0, 15*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateIronOre] = anim.SubImages(blockAtlas, 0, 16*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateDiamondOre] = anim.SubImages(blockAtlas, 0, 17*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateCopperOre] = anim.SubImages(blockAtlas, 0, 18*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateEmeraldOre] = anim.SubImages(blockAtlas, 0, 19*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateLapisOre] = anim.SubImages(blockAtlas, 0, 20*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateRedStoneOre] = anim.SubImages(blockAtlas, 0, 21*s, s, s, 11, false)
	SpriteFrames[itm.Log] = anim.SubImages(blockAtlas, 0, 22*s, s, s, 11, false)
	SpriteFrames[itm.Leaves] = anim.SubImages(blockAtlas, 0, 23*s, s, s, 11, false)
	SpriteFrames[itm.Planks] = anim.SubImages(blockAtlas, 0, 24*s, s, s, 11, false)
	SpriteFrames[itm.Cobblestone] = anim.SubImages(blockAtlas, 0, 25*s, s, s, 11, false)
	SpriteFrames[itm.CobbledDeepslate] = anim.SubImages(blockAtlas, 0, 26*s, s, s, 11, false)
	// items
	SpriteFrames[itm.Sapling] = anim.SubImages(itemAtlas, 0, 0, s, s, 1, false)
	SpriteFrames[itm.Torch] = anim.SubImages(itemAtlas, s, 0, s, s, 1, false)
	SpriteFrames[itm.Coal] = anim.SubImages(itemAtlas, 2*s, 0, s, s, 1, false)
	SpriteFrames[itm.RawGold] = anim.SubImages(itemAtlas, 3*s, 0, s, s, 1, false)
	SpriteFrames[itm.RawIron] = anim.SubImages(itemAtlas, 4*s, 0, s, s, 1, false)
	SpriteFrames[itm.Diamond] = anim.SubImages(itemAtlas, 5*s, 0, s, s, 1, false)
	SpriteFrames[itm.RawCopper] = anim.SubImages(itemAtlas, 6*s, 0, s, s, 1, false)
	SpriteFrames[itm.CharCoal] = anim.SubImages(itemAtlas, 7*s, 0, s, s, 1, false)
	SpriteFrames[itm.Emerald] = anim.SubImages(itemAtlas, 8*s, 0, s, s, 1, false)
	SpriteFrames[itm.LapisLazuli] = anim.SubImages(itemAtlas, 9*s, 0, s, s, 1, false)
	SpriteFrames[itm.Redstone] = anim.SubImages(itemAtlas, 10*s, 0, s, s, 1, false)
	SpriteFrames[itm.WoodAxe] = anim.SubImages(itemAtlas, 11*s, 0, s, s, 1, false)
	SpriteFrames[itm.WoodHoe] = anim.SubImages(itemAtlas, 12*s, 0, s, s, 1, false)
	SpriteFrames[itm.WoodPickaxe] = anim.SubImages(itemAtlas, 13*s, 0, s, s, 1, false)
	SpriteFrames[itm.WoodShovel] = anim.SubImages(itemAtlas, 14*s, 0, s, s, 1, false)
	SpriteFrames[itm.WoodSword] = anim.SubImages(itemAtlas, 15*s, 0, s, s, 1, false)
	SpriteFrames[itm.StoneAxe] = anim.SubImages(itemAtlas, 16*s, 0, s, s, 1, false)
	SpriteFrames[itm.StoneHoe] = anim.SubImages(itemAtlas, 17*s, 0, s, s, 1, false)
	SpriteFrames[itm.StonePickaxe] = anim.SubImages(itemAtlas, 18*s, 0, s, s, 1, false)
	SpriteFrames[itm.StoneShovel] = anim.SubImages(itemAtlas, 19*s, 0, s, s, 1, false)
	SpriteFrames[itm.StoneSword] = anim.SubImages(itemAtlas, 20*s, 0, s, s, 1, false)
	SpriteFrames[itm.GoldenAxe] = anim.SubImages(itemAtlas, 21*s, 0, s, s, 1, false)
	SpriteFrames[itm.GoldenHoe] = anim.SubImages(itemAtlas, 22*s, 0, s, s, 1, false)
	SpriteFrames[itm.GoldenPickaxe] = anim.SubImages(itemAtlas, 23*s, 0, s, s, 1, false)
	SpriteFrames[itm.GoldenShovel] = anim.SubImages(itemAtlas, 24*s, 0, s, s, 1, false)
	SpriteFrames[itm.GoldenSword] = anim.SubImages(itemAtlas, 25*s, 0, s, s, 1, false)
	SpriteFrames[itm.IronAxe] = anim.SubImages(itemAtlas, 26*s, 0, s, s, 1, false)
	SpriteFrames[itm.IronHoe] = anim.SubImages(itemAtlas, 27*s, 0, s, s, 1, false)
	SpriteFrames[itm.IronPickaxe] = anim.SubImages(itemAtlas, 28*s, 0, s, s, 1, false)
	SpriteFrames[itm.IronShovel] = anim.SubImages(itemAtlas, 29*s, 0, s, s, 1, false)
	SpriteFrames[itm.IronSword] = anim.SubImages(itemAtlas, 30*s, 0, s, s, 1, false)
	SpriteFrames[itm.DiamondAxe] = anim.SubImages(itemAtlas, 31*s, 0, s, s, 1, false)
	SpriteFrames[itm.DiamondHoe] = anim.SubImages(itemAtlas, 0, s, s, s, 1, false)
	SpriteFrames[itm.DiamondPickaxe] = anim.SubImages(itemAtlas, s, s, s, s, 1, false)
	SpriteFrames[itm.DiamondShovel] = anim.SubImages(itemAtlas, 2*s, s, s, s, 1, false)
	SpriteFrames[itm.DiamondSword] = anim.SubImages(itemAtlas, 3*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheriteAxe] = anim.SubImages(itemAtlas, 4*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheriteHoe] = anim.SubImages(itemAtlas, 5*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheritePickaxe] = anim.SubImages(itemAtlas, 6*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheriteShovel] = anim.SubImages(itemAtlas, 7*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheriteSword] = anim.SubImages(itemAtlas, 8*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheriteScrap] = anim.SubImages(itemAtlas, 9*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheriteIngot] = anim.SubImages(itemAtlas, 10*s, s, s, s, 1, false)
	SpriteFrames[itm.GoldIngot] = anim.SubImages(itemAtlas, 11*s, s, s, s, 1, false)
	SpriteFrames[itm.IronIngot] = anim.SubImages(itemAtlas, 12*s, s, s, s, 1, false)
	SpriteFrames[itm.CopperIngot] = anim.SubImages(itemAtlas, 13*s, s, s, s, 1, false)
	SpriteFrames[itm.Stick] = anim.SubImages(itemAtlas, 14*s, s, s, s, 1, false)

}
