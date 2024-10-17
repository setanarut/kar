package res

import (
	"image"
	"kar/engine/util"
	"kar/items"

	"github.com/anthonynsimon/bild/blend"
	"github.com/hajimehoshi/ebiten/v2"
)

var Images = make(map[uint16]*ebiten.Image, 0)
var Frames = make(map[uint16][]*ebiten.Image, 0)
var BlockBorder *ebiten.Image

var (
	AtlasPlayer     = util.ReadEbImgFS(fs, "assets/img/player/player.png")
	Hotbar          = util.ReadEbImgFS(fs, "assets/img/gui/hotbar.png")
	HotbarSelection = util.ReadEbImgFS(fs, "assets/img/gui/hotbarBorder.png")
)

var cracks = util.ImgFromFS(fs, "assets/img/overlay/cracks.png")

func init() {

	switch BlockSize {
	case 32:
		BlockBorder = util.ReadEbImgFS(fs, "assets/img/overlay/border32.png")
	case 48:
		BlockBorder = util.ReadEbImgFS(fs, "assets/img/overlay/border48.png")
	case 64:
		BlockBorder = util.ReadEbImgFS(fs, "assets/img/overlay/border64.png")
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
	Images[items.Arrow] = ItemImg("arrow.png")
	Images[items.BeetrootSeeds] = ItemImg("beetroot_seeds.png")
	Images[items.Bow] = ItemImg("bow.png")
	Images[items.Bread] = ItemImg("bread.png")
	Images[items.Bucket] = ItemImg("bucket.png")
	Images[items.Charcoal] = ItemImg("charcoal.png")
	Images[items.Coal] = ItemImg("coal.png")
	Images[items.CopperIngot] = ItemImg("copper_ingot.png")
	Images[items.CrossbowStandby] = ItemImg("crossbow_standby.png")
	Images[items.Diamond] = ItemImg("diamond.png")
	Images[items.DiamondAxe] = ItemImg("diamond_axe.png")
	Images[items.DiamondHoe] = ItemImg("diamond_hoe.png")
	Images[items.DiamondPickaxe] = ItemImg("diamond_pickaxe.png")
	Images[items.DiamondShovel] = ItemImg("diamond_shovel.png")
	Images[items.DiamondSword] = ItemImg("diamond_sword.png")
	Images[items.Emerald] = ItemImg("emerald.png")
	Images[items.GoldIngot] = ItemImg("gold_ingot.png")
	Images[items.GoldenAxe] = ItemImg("golden_axe.png")
	Images[items.GoldenHoe] = ItemImg("golden_hoe.png")
	Images[items.GoldenPickaxe] = ItemImg("golden_pickaxe.png")
	Images[items.GoldenShovel] = ItemImg("golden_shovel.png")
	Images[items.GoldenSword] = ItemImg("golden_sword.png")
	Images[items.IronAxe] = ItemImg("iron_axe.png")
	Images[items.IronHoe] = ItemImg("iron_hoe.png")
	Images[items.IronIngot] = ItemImg("iron_ingot.png")
	Images[items.IronPickaxe] = ItemImg("iron_pickaxe.png")
	Images[items.IronShovel] = ItemImg("iron_shovel.png")
	Images[items.IronSword] = ItemImg("iron_sword.png")
	Images[items.LapisLazuli] = ItemImg("lapis_lazuli.png")
	Images[items.LavaBucket] = ItemImg("lava_bucket.png")
	Images[items.MelonSeeds] = ItemImg("melon_seeds.png")
	Images[items.MilkBucket] = ItemImg("milk_bucket.png")
	Images[items.NetheriteAxe] = ItemImg("netherite_axe.png")
	Images[items.NetheriteHoe] = ItemImg("netherite_hoe.png")
	Images[items.NetheriteIngot] = ItemImg("netherite_ingot.png")
	Images[items.NetheritePickaxe] = ItemImg("netherite_pickaxe.png")
	Images[items.NetheriteScrap] = ItemImg("netherite_scrap.png")
	Images[items.NetheriteShovel] = ItemImg("netherite_shovel.png")
	Images[items.NetheriteSword] = ItemImg("netherite_sword.png")
	Images[items.PowderSnowBucket] = ItemImg("powder_snow_bucket.png")
	Images[items.PumpkinSeeds] = ItemImg("pumpkin_seeds.png")
	Images[items.RawCopper] = ItemImg("raw_copper.png")
	Images[items.RawGold] = ItemImg("raw_gold.png")
	Images[items.RawIron] = ItemImg("raw_iron.png")
	Images[items.Redstone] = ItemImg("redstone.png")
	Images[items.Snowball] = ItemImg("snowball.png")
	Images[items.Stick] = ItemImg("stick.png")
	Images[items.StoneAxe] = ItemImg("stone_axe.png")
	Images[items.StoneHoe] = ItemImg("stone_hoe.png")
	Images[items.StonePickaxe] = ItemImg("stone_pickaxe.png")
	Images[items.StoneShovel] = ItemImg("stone_shovel.png")
	Images[items.StoneSword] = ItemImg("stone_sword.png")
	Images[items.TorchflowerSeeds] = ItemImg("torchflower_seeds.png")
	Images[items.WaterBucket] = ItemImg("water_bucket.png")
	Images[items.Wheat] = ItemImg("wheat.png")
	Images[items.WheatSeeds] = ItemImg("wheat_seeds.png")
	Images[items.WoodenAxe] = ItemImg("wooden_axe.png")
	Images[items.WoodenHoe] = ItemImg("wooden_hoe.png")
	Images[items.WoodenPickaxe] = ItemImg("wooden_pickaxe.png")
	Images[items.WoodenShovel] = ItemImg("wooden_shovel.png")
	Images[items.WoodenSword] = ItemImg("wooden_sword.png")

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
func ItemImg(f string) *ebiten.Image {
	return util.ReadEbImgFS(fs, "assets/img/items/"+f)
}
