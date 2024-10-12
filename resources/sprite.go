package resources

import (
	"image"
	"kar/engine/util"
	"kar/items"

	"github.com/anthonynsimon/bild/blend"
	"github.com/hajimehoshi/ebiten/v2"
)

var Sprite = make(map[uint16]*ebiten.Image, 0)
var SpriteStages = make(map[uint16][]*ebiten.Image, 0)
var BlockHighlightBorder = util.LoadEbitenImageFromFS(assets, "assets/sprite/overlay/border_64.png")
var blockBreakStagesOverlay image.Image = util.LoadImageFromFS(assets, "assets/sprite/overlay/break_stages.png")

var (
	AtlasPlayer     = util.LoadEbitenImageFromFS(assets, "assets/sprite/player/player_atlas.png")
	Hotbar          = util.LoadEbitenImageFromFS(assets, "assets/sprite/gui/hotbar.png")
	HotbarSelection = util.LoadEbitenImageFromFS(assets, "assets/sprite/gui/hotbar_selected.png")
)

func init() {
	SpriteStages[items.Air] = getFrames("air.png")
	SpriteStages[items.Andesite] = getFrames("andesite.png")
	SpriteStages[items.Bedrock] = getFrames("bedrock.png")
	SpriteStages[items.BrewingStand] = getFrames("brewing_stand.png")
	SpriteStages[items.CartographyTable] = getFrames("cartography_table.png")
	SpriteStages[items.Clay] = getFrames("clay.png")
	SpriteStages[items.CoalBlock] = getFrames("coal_block.png")
	SpriteStages[items.CoalOre] = getFrames("coal_ore.png")
	SpriteStages[items.CoarseDirt] = getFrames("coarse_dirt.png")
	SpriteStages[items.CobbledDeepslate] = getFrames("cobbled_deepslate.png")
	SpriteStages[items.Cobblestone] = getFrames("cobblestone.png")
	SpriteStages[items.CopperOre] = getFrames("copper_ore.png")
	SpriteStages[items.CraftingTable] = getFrames("crafting_table.png")
	SpriteStages[items.CryingObsidian] = getFrames("crying_obsidian.png")
	SpriteStages[items.Deepslate] = getFrames("deepslate.png")
	SpriteStages[items.DeepslateCoalOre] = getFrames("deepslate_coal_ore.png")
	SpriteStages[items.DeepslateCopperOre] = getFrames("deepslate_copper_ore.png")
	SpriteStages[items.DeepslateDiamondOre] = getFrames("deepslate_diamond_ore.png")
	SpriteStages[items.DeepslateEmeraldOre] = getFrames("deepslate_emerald_ore.png")
	SpriteStages[items.DeepslateGoldOre] = getFrames("deepslate_gold_ore.png")
	SpriteStages[items.DeepslateIronOre] = getFrames("deepslate_iron_ore.png")
	SpriteStages[items.DeepslateLapisOre] = getFrames("deepslate_lapis_ore.png")
	SpriteStages[items.DeepslateRedstoneOre] = getFrames("deepslate_redstone_ore.png")
	SpriteStages[items.DiamondOre] = getFrames("diamond_ore.png")
	SpriteStages[items.Dirt] = getFrames("dirt.png")
	SpriteStages[items.DirtPath] = getFrames("dirt_path.png")
	SpriteStages[items.EmeraldOre] = getFrames("emerald_ore.png")
	SpriteStages[items.EnchantingTable] = getFrames("enchanting_table.png")
	SpriteStages[items.EndPortalFrame] = getFrames("end_portal_frame.png")
	SpriteStages[items.FletchingTable] = getFrames("fletching_table.png")
	SpriteStages[items.Furnace] = getFrames("furnace.png")
	SpriteStages[items.FurnaceOn] = getFrames("furnace_on.png")
	SpriteStages[items.GoldOre] = getFrames("gold_ore.png")
	SpriteStages[items.GrassBlock] = getFrames("grass_block.png")
	SpriteStages[items.GrassBlockSnow] = getFrames("grass_block_snow.png")
	SpriteStages[items.Gravel] = getFrames("gravel.png")
	SpriteStages[items.IronOre] = getFrames("iron_ore.png")
	SpriteStages[items.LapisOre] = getFrames("lapis_ore.png")
	SpriteStages[items.NetherBricks] = getFrames("nether_bricks.png")
	SpriteStages[items.NetherGoldOre] = getFrames("nether_gold_ore.png")
	SpriteStages[items.NetherQuartzOre] = getFrames("nether_quartz_ore.png")
	SpriteStages[items.Netherrack] = getFrames("netherrack.png")
	SpriteStages[items.OakLeaves] = getFrames("oak_leaves.png")
	SpriteStages[items.OakLog] = getFrames("oak_log.png")
	SpriteStages[items.OakPlanks] = getFrames("oak_planks.png")
	SpriteStages[items.OakSapling] = getFrames("oak_sapling.png")
	SpriteStages[items.Obsidian] = getFrames("obsidian.png")
	SpriteStages[items.RedNetherBricks] = getFrames("red_nether_bricks.png")
	SpriteStages[items.RedSand] = getFrames("red_sand.png")
	SpriteStages[items.RedSandstone] = getFrames("red_sandstone.png")
	SpriteStages[items.RedstoneOre] = getFrames("redstone_ore.png")
	SpriteStages[items.RedstoneTorch] = getFrames("redstone_torch.png")
	SpriteStages[items.RedstoneTorchOff] = getFrames("redstone_torch_off.png")
	SpriteStages[items.RootedDirt] = getFrames("rooted_dirt.png")
	SpriteStages[items.Sand] = getFrames("sand.png")
	SpriteStages[items.Sandstone] = getFrames("sandstone.png")
	SpriteStages[items.SmithingTable] = getFrames("smithing_table.png")
	SpriteStages[items.SmoothStone] = getFrames("smooth_stone.png")
	SpriteStages[items.Snow] = getFrames("snow.png")
	SpriteStages[items.SoulSand] = getFrames("soul_sand.png")
	SpriteStages[items.SoulSoil] = getFrames("soul_soil.png")
	SpriteStages[items.SoulTorch] = getFrames("soul_torch.png")
	SpriteStages[items.Stone] = getFrames("stone.png")
	SpriteStages[items.StoneBricks] = getFrames("stone_bricks.png")
	SpriteStages[items.Tnt] = getFrames("tnt.png")
	SpriteStages[items.Torch] = getFrames("torch.png")

	SpriteStages[items.WheatCrops] = []*ebiten.Image{
		util.LoadEbitenImageFromFS(assets, "assets/sprite/blocks_stages/wheat_crops/wheat_stage0.png"),
		util.LoadEbitenImageFromFS(assets, "assets/sprite/blocks_stages/wheat_crops/wheat_stage1.png"),
		util.LoadEbitenImageFromFS(assets, "assets/sprite/blocks_stages/wheat_crops/wheat_stage2.png"),
		util.LoadEbitenImageFromFS(assets, "assets/sprite/blocks_stages/wheat_crops/wheat_stage3.png"),
		util.LoadEbitenImageFromFS(assets, "assets/sprite/blocks_stages/wheat_crops/wheat_stage4.png"),
		util.LoadEbitenImageFromFS(assets, "assets/sprite/blocks_stages/wheat_crops/wheat_stage5.png"),
		util.LoadEbitenImageFromFS(assets, "assets/sprite/blocks_stages/wheat_crops/wheat_stage6.png"),
		util.LoadEbitenImageFromFS(assets, "assets/sprite/blocks_stages/wheat_crops/wheat_stage7.png"),
	}

	Sprite[items.Arrow] = getItemSprite("arrow.png")
	Sprite[items.BeetrootSeeds] = getItemSprite("beetroot_seeds.png")
	Sprite[items.Bow] = getItemSprite("bow.png")
	Sprite[items.Bread] = getItemSprite("bread.png")
	Sprite[items.Bucket] = getItemSprite("bucket.png")
	Sprite[items.Charcoal] = getItemSprite("charcoal.png")
	Sprite[items.Coal] = getItemSprite("coal.png")
	Sprite[items.CopperIngot] = getItemSprite("copper_ingot.png")
	Sprite[items.CrossbowStandby] = getItemSprite("crossbow_standby.png")
	Sprite[items.Diamond] = getItemSprite("diamond.png")
	Sprite[items.DiamondAxe] = getItemSprite("diamond_axe.png")
	Sprite[items.DiamondHoe] = getItemSprite("diamond_hoe.png")
	Sprite[items.DiamondPickaxe] = getItemSprite("diamond_pickaxe.png")
	Sprite[items.DiamondShovel] = getItemSprite("diamond_shovel.png")
	Sprite[items.DiamondSword] = getItemSprite("diamond_sword.png")
	Sprite[items.Emerald] = getItemSprite("emerald.png")
	Sprite[items.GoldIngot] = getItemSprite("gold_ingot.png")
	Sprite[items.GoldenAxe] = getItemSprite("golden_axe.png")
	Sprite[items.GoldenHoe] = getItemSprite("golden_hoe.png")
	Sprite[items.GoldenPickaxe] = getItemSprite("golden_pickaxe.png")
	Sprite[items.GoldenShovel] = getItemSprite("golden_shovel.png")
	Sprite[items.GoldenSword] = getItemSprite("golden_sword.png")
	Sprite[items.IronAxe] = getItemSprite("iron_axe.png")
	Sprite[items.IronHoe] = getItemSprite("iron_hoe.png")
	Sprite[items.IronIngot] = getItemSprite("iron_ingot.png")
	Sprite[items.IronPickaxe] = getItemSprite("iron_pickaxe.png")
	Sprite[items.IronShovel] = getItemSprite("iron_shovel.png")
	Sprite[items.IronSword] = getItemSprite("iron_sword.png")
	Sprite[items.LapisLazuli] = getItemSprite("lapis_lazuli.png")
	Sprite[items.LavaBucket] = getItemSprite("lava_bucket.png")
	Sprite[items.MelonSeeds] = getItemSprite("melon_seeds.png")
	Sprite[items.MilkBucket] = getItemSprite("milk_bucket.png")
	Sprite[items.NetheriteAxe] = getItemSprite("netherite_axe.png")
	Sprite[items.NetheriteHoe] = getItemSprite("netherite_hoe.png")
	Sprite[items.NetheriteIngot] = getItemSprite("netherite_ingot.png")
	Sprite[items.NetheritePickaxe] = getItemSprite("netherite_pickaxe.png")
	Sprite[items.NetheriteScrap] = getItemSprite("netherite_scrap.png")
	Sprite[items.NetheriteShovel] = getItemSprite("netherite_shovel.png")
	Sprite[items.NetheriteSword] = getItemSprite("netherite_sword.png")
	Sprite[items.PowderSnowBucket] = getItemSprite("powder_snow_bucket.png")
	Sprite[items.PumpkinSeeds] = getItemSprite("pumpkin_seeds.png")
	Sprite[items.RawCopper] = getItemSprite("raw_copper.png")
	Sprite[items.RawGold] = getItemSprite("raw_gold.png")
	Sprite[items.RawIron] = getItemSprite("raw_iron.png")
	Sprite[items.Redstone] = getItemSprite("redstone.png")
	Sprite[items.Snowball] = getItemSprite("snowball.png")
	Sprite[items.Stick] = getItemSprite("stick.png")
	Sprite[items.StoneAxe] = getItemSprite("stone_axe.png")
	Sprite[items.StoneHoe] = getItemSprite("stone_hoe.png")
	Sprite[items.StonePickaxe] = getItemSprite("stone_pickaxe.png")
	Sprite[items.StoneShovel] = getItemSprite("stone_shovel.png")
	Sprite[items.StoneSword] = getItemSprite("stone_sword.png")
	Sprite[items.TorchflowerSeeds] = getItemSprite("torchflower_seeds.png")
	Sprite[items.WaterBucket] = getItemSprite("water_bucket.png")
	Sprite[items.Wheat] = getItemSprite("wheat.png")
	Sprite[items.WheatSeeds] = getItemSprite("wheat_seeds.png")
	Sprite[items.WoodenAxe] = getItemSprite("wooden_axe.png")
	Sprite[items.WoodenHoe] = getItemSprite("wooden_hoe.png")
	Sprite[items.WoodenPickaxe] = getItemSprite("wooden_pickaxe.png")
	Sprite[items.WoodenShovel] = getItemSprite("wooden_shovel.png")
	Sprite[items.WoodenSword] = getItemSprite("wooden_sword.png")

}

func convert2ebiten(st []image.Image) []*ebiten.Image {
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
func getFrames(f string) []*ebiten.Image {
	return convert2ebiten(makeStages(util.LoadImageFromFS(assets, "assets/sprite/blocks/"+f), blockBreakStagesOverlay))
}
func getItemSprite(f string) *ebiten.Image {
	return util.LoadEbitenImageFromFS(assets, "assets/sprite/items/"+f)
}
