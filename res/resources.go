// res is resources
package res

import (
	"embed"
	"image"
	"kar/engine/util"
	"kar/items"

	"github.com/anthonynsimon/bild/blend"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var fs embed.FS

var (
	Images         = make(map[uint16]*ebiten.Image, 0)
	Frames         = make(map[uint16][]*ebiten.Image, 0)
	Hotbar         = util.ReadEbImgFS(fs, "assets/img/gui/hotbar.png")
	SelectionBar   = util.ReadEbImgFS(fs, "assets/img/gui/selection_bar.png")
	SelectionBlock = util.ReadEbImgFS(fs, "assets/img/gui/selection_block.png")
	Font           = util.LoadFontFromFS("assets/font/pixelcode.otf", 18, fs)
	PlayerAtlas    = util.ReadEbImgFS(fs, "assets/img/player/player.png")
	cracks         = util.ImgFromFS(fs, "assets/img/cracks.png")
)

func init() {
	Images[items.Air] = util.ReadEbImgFS(fs, "assets/img/air.png")
	Frames[items.Bedrock] = blockImgs("bedrock.png")
	Frames[items.CoalBlock] = blockImgs("coal_block.png")
	Frames[items.CoalOre] = blockImgs("coal_ore.png")
	Frames[items.Cobblestone] = blockImgs("cobblestone.png")
	Frames[items.CraftingTable] = blockImgs("crafting_table.png")
	Frames[items.DiamondOre] = blockImgs("diamond_ore.png")
	Frames[items.Dirt] = blockImgs("dirt.png")
	Frames[items.Furnace] = blockImgs("furnace.png")
	Frames[items.FurnaceOn] = blockImgs("furnace_on.png")
	Frames[items.GoldOre] = blockImgs("gold_ore.png")
	Frames[items.GrassBlock] = blockImgs("grass_block.png")
	Frames[items.GrassBlockSnow] = blockImgs("grass_block_snow.png")
	Frames[items.IronOre] = blockImgs("iron_ore.png")
	Frames[items.OakLeaves] = blockImgs("oak_leaves.png")
	Frames[items.OakLog] = blockImgs("oak_log.png")
	Frames[items.OakPlanks] = blockImgs("oak_planks.png")
	Frames[items.OakSapling] = blockImgs("oak_sapling.png")
	Frames[items.Obsidian] = blockImgs("obsidian.png")
	Frames[items.Sand] = blockImgs("sand.png")
	Frames[items.SmoothStone] = blockImgs("smooth_stone.png")
	Frames[items.Snow] = blockImgs("snow.png")
	Frames[items.Stone] = blockImgs("stone.png")
	Frames[items.StoneBricks] = blockImgs("stone_bricks.png")
	Frames[items.Tnt] = blockImgs("tnt.png")
	Frames[items.Torch] = blockImgs("torch.png")

	Images[items.Bread] = itemImg("bread.png")
	Images[items.Bucket] = itemImg("bucket.png")
	Images[items.Coal] = itemImg("coal.png")
	Images[items.Diamond] = itemImg("diamond.png")
	Images[items.DiamondAxe] = itemImg("diamond_axe.png")
	Images[items.DiamondPickaxe] = itemImg("diamond_pickaxe.png")
	Images[items.DiamondShovel] = itemImg("diamond_shovel.png")
	Images[items.GoldIngot] = itemImg("gold_ingot.png")
	Images[items.GoldenAxe] = itemImg("golden_axe.png")
	Images[items.GoldenPickaxe] = itemImg("golden_pickaxe.png")
	Images[items.GoldenShovel] = itemImg("golden_shovel.png")
	Images[items.IronAxe] = itemImg("iron_axe.png")
	Images[items.IronIngot] = itemImg("iron_ingot.png")
	Images[items.IronPickaxe] = itemImg("iron_pickaxe.png")
	Images[items.IronShovel] = itemImg("iron_shovel.png")
	Images[items.RawGold] = itemImg("raw_gold.png")
	Images[items.RawIron] = itemImg("raw_iron.png")
	Images[items.Snowball] = itemImg("snowball.png")
	Images[items.Stick] = itemImg("stick.png")
	Images[items.StoneAxe] = itemImg("stone_axe.png")
	Images[items.StonePickaxe] = itemImg("stone_pickaxe.png")
	Images[items.StoneShovel] = itemImg("stone_shovel.png")
	Images[items.WaterBucket] = itemImg("water_bucket.png")
	Images[items.WoodenAxe] = itemImg("wooden_axe.png")
	Images[items.WoodenPickaxe] = itemImg("wooden_pickaxe.png")
	Images[items.WoodenShovel] = itemImg("wooden_shovel.png")

}

func toEbiten(st []image.Image) []*ebiten.Image {
	l := make([]*ebiten.Image, 0)
	for _, v := range st {
		l = append(l, ebiten.NewImageFromImage(v))
	}
	return l
}

func makeStages(block, stages image.Image) []image.Image {
	frames := make([]image.Image, 0)
	frames = append(frames, block)
	for i := range 4 {
		x := i * 16
		rec := image.Rect(x, 0, x+16, x+16)
		si := stages.(*image.NRGBA).SubImage(rec)
		frames = append(frames, blend.Normal(block, si))
	}
	return frames
}
func blockImgs(f string) []*ebiten.Image {
	frames := makeStages(util.ImgFromFS(fs, "assets/img/blocks/"+f), cracks)
	return toEbiten(frames)
}
func itemImg(f string) *ebiten.Image {
	return util.ReadEbImgFS(fs, "assets/img/items/"+f)
}
