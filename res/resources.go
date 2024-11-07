// res is resources
package res

import (
	"embed"
	"image"
	"kar"
	"kar/engine/util"
	"kar/items"

	"github.com/anthonynsimon/bild/blend"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var fs embed.FS

var (
	Images          = make(map[uint16]*ebiten.Image, 0)
	Frames          = make(map[uint16][]*ebiten.Image, 0)
	AtlasPlayer     = util.ReadEbImgFS(fs, "assets/img/player/player.png")
	Hotbar          = util.ReadEbImgFS(fs, "assets/img/gui/hotbar.png")
	HotbarSelection = util.ReadEbImgFS(fs, "assets/img/gui/hotbarBorder.png")
	Font            = util.LoadFontFromFS("assets/font/pixelcode.otf", 18, fs)
)

var (
	Mario = util.ReadEbImgFS(fs, "assets/img/player/mario.png")
)

var cracks = util.ImgFromFS(fs, "assets/img/overlay/cracks.png")
var Border *ebiten.Image

func init() {

	switch kar.BlockSize {
	case 32:
		Border = util.ReadEbImgFS(fs, "assets/img/overlay/border32.png")
	case 48:
		Border = util.ReadEbImgFS(fs, "assets/img/overlay/border48.png")
	case 64:
		Border = util.ReadEbImgFS(fs, "assets/img/overlay/border64.png")
	case 80:
		Border = util.ReadEbImgFS(fs, "assets/img/overlay/border80.png")
	}

	Images[items.Air] = util.ReadEbImgFS(fs, "assets/img/air.png")
	Frames[items.Andesite] = blockImgs("andesite.png")
	Frames[items.Bedrock] = blockImgs("bedrock.png")
	Frames[items.BrewingStand] = blockImgs("brewing_stand.png")
	Frames[items.CartographyTable] = blockImgs("cartography_table.png")
	Frames[items.Clay] = blockImgs("clay.png")
	Frames[items.CoalBlock] = blockImgs("coal_block.png")
	Frames[items.CoalOre] = blockImgs("coal_ore.png")
	Frames[items.CoarseDirt] = blockImgs("coarse_dirt.png")
	Frames[items.CobbledDeepslate] = blockImgs("cobbled_deepslate.png")
	Frames[items.Cobblestone] = blockImgs("cobblestone.png")
	Frames[items.CopperOre] = blockImgs("copper_ore.png")
	Frames[items.CraftingTable] = blockImgs("crafting_table.png")
	Frames[items.CryingObsidian] = blockImgs("crying_obsidian.png")
	Frames[items.Deepslate] = blockImgs("deepslate.png")
	Frames[items.DeepslateCoalOre] = blockImgs("deepslate_coal_ore.png")
	Frames[items.DeepslateCopperOre] = blockImgs("deepslate_copper_ore.png")
	Frames[items.DeepslateDiamondOre] = blockImgs("deepslate_diamond_ore.png")
	Frames[items.DeepslateEmeraldOre] = blockImgs("deepslate_emerald_ore.png")
	Frames[items.DeepslateGoldOre] = blockImgs("deepslate_gold_ore.png")
	Frames[items.DeepslateIronOre] = blockImgs("deepslate_iron_ore.png")
	Frames[items.DeepslateLapisOre] = blockImgs("deepslate_lapis_ore.png")
	Frames[items.DeepslateRedstoneOre] = blockImgs("deepslate_redstone_ore.png")
	Frames[items.DiamondOre] = blockImgs("diamond_ore.png")
	Frames[items.Dirt] = blockImgs("dirt.png")
	Frames[items.DirtPath] = blockImgs("dirt_path.png")
	Frames[items.EmeraldOre] = blockImgs("emerald_ore.png")
	Frames[items.EnchantingTable] = blockImgs("enchanting_table.png")
	Frames[items.EndPortalFrame] = blockImgs("end_portal_frame.png")
	Frames[items.FletchingTable] = blockImgs("fletching_table.png")
	Frames[items.Furnace] = blockImgs("furnace.png")
	Frames[items.FurnaceOn] = blockImgs("furnace_on.png")
	Frames[items.GoldOre] = blockImgs("gold_ore.png")
	Frames[items.GrassBlock] = blockImgs("grass_block.png")
	Frames[items.GrassBlockSnow] = blockImgs("grass_block_snow.png")
	Frames[items.Gravel] = blockImgs("gravel.png")
	Frames[items.IronOre] = blockImgs("iron_ore.png")
	Frames[items.LapisOre] = blockImgs("lapis_ore.png")
	Frames[items.NetherBricks] = blockImgs("nether_bricks.png")
	Frames[items.NetherGoldOre] = blockImgs("nether_gold_ore.png")
	Frames[items.NetherQuartzOre] = blockImgs("nether_quartz_ore.png")
	Frames[items.Netherrack] = blockImgs("netherrack.png")
	Frames[items.OakLeaves] = blockImgs("oak_leaves.png")
	Frames[items.OakLog] = blockImgs("oak_log.png")
	Frames[items.OakPlanks] = blockImgs("oak_planks.png")
	Frames[items.OakSapling] = blockImgs("oak_sapling.png")
	Frames[items.Obsidian] = blockImgs("obsidian.png")
	Frames[items.RedNetherBricks] = blockImgs("red_nether_bricks.png")
	Frames[items.RedSand] = blockImgs("red_sand.png")
	Frames[items.RedSandstone] = blockImgs("red_sandstone.png")
	Frames[items.RedstoneOre] = blockImgs("redstone_ore.png")
	Frames[items.RedstoneTorch] = blockImgs("redstone_torch.png")
	Frames[items.RedstoneTorchOff] = blockImgs("redstone_torch_off.png")
	Frames[items.RootedDirt] = blockImgs("rooted_dirt.png")
	Frames[items.Sand] = blockImgs("sand.png")
	Frames[items.Sandstone] = blockImgs("sandstone.png")
	Frames[items.SmithingTable] = blockImgs("smithing_table.png")
	Frames[items.SmoothStone] = blockImgs("smooth_stone.png")
	Frames[items.Snow] = blockImgs("snow.png")
	Frames[items.SoulSand] = blockImgs("soul_sand.png")
	Frames[items.SoulSoil] = blockImgs("soul_soil.png")
	Frames[items.SoulTorch] = blockImgs("soul_torch.png")
	Frames[items.Stone] = blockImgs("stone.png")
	Frames[items.StoneBricks] = blockImgs("stone_bricks.png")
	Frames[items.Tnt] = blockImgs("tnt.png")
	Frames[items.Torch] = blockImgs("torch.png")

	Frames[items.WheatCrops] = []*ebiten.Image{
		util.ReadEbImgFS(fs, "assets/img/blocks_stages/wheat_crops/wheat0.png"),
		util.ReadEbImgFS(fs, "assets/img/blocks_stages/wheat_crops/wheat1.png"),
		util.ReadEbImgFS(fs, "assets/img/blocks_stages/wheat_crops/wheat2.png"),
		util.ReadEbImgFS(fs, "assets/img/blocks_stages/wheat_crops/wheat3.png"),
		util.ReadEbImgFS(fs, "assets/img/blocks_stages/wheat_crops/wheat4.png"),
		util.ReadEbImgFS(fs, "assets/img/blocks_stages/wheat_crops/wheat5.png"),
		util.ReadEbImgFS(fs, "assets/img/blocks_stages/wheat_crops/wheat6.png"),
		util.ReadEbImgFS(fs, "assets/img/blocks_stages/wheat_crops/wheat7.png"),
	}
	Images[items.Arrow] = itemImg("arrow.png")
	Images[items.BeetrootSeeds] = itemImg("beetroot_seeds.png")
	Images[items.Bow] = itemImg("bow.png")
	Images[items.Bread] = itemImg("bread.png")
	Images[items.Bucket] = itemImg("bucket.png")
	Images[items.Charcoal] = itemImg("charcoal.png")
	Images[items.Coal] = itemImg("coal.png")
	Images[items.CopperIngot] = itemImg("copper_ingot.png")
	Images[items.CrossbowStandby] = itemImg("crossbow_standby.png")
	Images[items.Diamond] = itemImg("diamond.png")
	Images[items.DiamondAxe] = itemImg("diamond_axe.png")
	Images[items.DiamondHoe] = itemImg("diamond_hoe.png")
	Images[items.DiamondPickaxe] = itemImg("diamond_pickaxe.png")
	Images[items.DiamondShovel] = itemImg("diamond_shovel.png")
	Images[items.DiamondSword] = itemImg("diamond_sword.png")
	Images[items.Emerald] = itemImg("emerald.png")
	Images[items.GoldIngot] = itemImg("gold_ingot.png")
	Images[items.GoldenAxe] = itemImg("golden_axe.png")
	Images[items.GoldenHoe] = itemImg("golden_hoe.png")
	Images[items.GoldenPickaxe] = itemImg("golden_pickaxe.png")
	Images[items.GoldenShovel] = itemImg("golden_shovel.png")
	Images[items.GoldenSword] = itemImg("golden_sword.png")
	Images[items.IronAxe] = itemImg("iron_axe.png")
	Images[items.IronHoe] = itemImg("iron_hoe.png")
	Images[items.IronIngot] = itemImg("iron_ingot.png")
	Images[items.IronPickaxe] = itemImg("iron_pickaxe.png")
	Images[items.IronShovel] = itemImg("iron_shovel.png")
	Images[items.IronSword] = itemImg("iron_sword.png")
	Images[items.LapisLazuli] = itemImg("lapis_lazuli.png")
	Images[items.LavaBucket] = itemImg("lava_bucket.png")
	Images[items.MelonSeeds] = itemImg("melon_seeds.png")
	Images[items.MilkBucket] = itemImg("milk_bucket.png")
	Images[items.NetheriteAxe] = itemImg("netherite_axe.png")
	Images[items.NetheriteHoe] = itemImg("netherite_hoe.png")
	Images[items.NetheriteIngot] = itemImg("netherite_ingot.png")
	Images[items.NetheritePickaxe] = itemImg("netherite_pickaxe.png")
	Images[items.NetheriteScrap] = itemImg("netherite_scrap.png")
	Images[items.NetheriteShovel] = itemImg("netherite_shovel.png")
	Images[items.NetheriteSword] = itemImg("netherite_sword.png")
	Images[items.PowderSnowBucket] = itemImg("powder_snow_bucket.png")
	Images[items.PumpkinSeeds] = itemImg("pumpkin_seeds.png")
	Images[items.RawCopper] = itemImg("raw_copper.png")
	Images[items.RawGold] = itemImg("raw_gold.png")
	Images[items.RawIron] = itemImg("raw_iron.png")
	Images[items.Redstone] = itemImg("redstone.png")
	Images[items.Snowball] = itemImg("snowball.png")
	Images[items.Stick] = itemImg("stick.png")
	Images[items.StoneAxe] = itemImg("stone_axe.png")
	Images[items.StoneHoe] = itemImg("stone_hoe.png")
	Images[items.StonePickaxe] = itemImg("stone_pickaxe.png")
	Images[items.StoneShovel] = itemImg("stone_shovel.png")
	Images[items.StoneSword] = itemImg("stone_sword.png")
	Images[items.TorchflowerSeeds] = itemImg("torchflower_seeds.png")
	Images[items.WaterBucket] = itemImg("water_bucket.png")
	Images[items.Wheat] = itemImg("wheat.png")
	Images[items.WheatSeeds] = itemImg("wheat_seeds.png")
	Images[items.WoodenAxe] = itemImg("wooden_axe.png")
	Images[items.WoodenHoe] = itemImg("wooden_hoe.png")
	Images[items.WoodenPickaxe] = itemImg("wooden_pickaxe.png")
	Images[items.WoodenShovel] = itemImg("wooden_shovel.png")
	Images[items.WoodenSword] = itemImg("wooden_sword.png")

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
	for i := range 10 {
		x := i * 16
		rec := image.Rect(x, 0, x+16, x+16)
		si := stages.(*image.NRGBA).SubImage(rec)
		frames = append(frames, blend.Overlay(block, si))
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
