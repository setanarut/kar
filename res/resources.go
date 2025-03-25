// res is resources
package res

import (
	"bytes"
	"embed"
	"image"
	"kar/items"
	"log"

	"github.com/anthonynsimon/bild/blend"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/setanarut/anim"
	"golang.org/x/text/language"
)

//go:embed assets/*
var fs embed.FS

var (
	Crab             = anim.SubImages(ReadEbImgFS(fs, "assets/img/crab.png"), 0, 0, 16, 9, 2, false)
	Icon8            = make(map[uint8]*ebiten.Image, 0)
	BlockCrackFrames = make(map[uint8][]*ebiten.Image, 0)
	BlockUnbreakable = make(map[uint8]*ebiten.Image, 0)
	HotbarEdge       = ReadEbImgFS(fs, "assets/img/gui/hotbar_edge.png")
	HotbarMid        = ReadEbImgFS(fs, "assets/img/gui/hotbar_mid.png")
	CraftingTable1x2 = ReadEbImgFS(fs, "assets/img/gui/table1x2.png")
	CraftingTable3x3 = ReadEbImgFS(fs, "assets/img/gui/table3x3.png")
	CraftingTable2x2 = ReadEbImgFS(fs, "assets/img/gui/table2x2.png")
	BlockBorder      = ReadEbImgFS(fs, "assets/img/gui/block_border.png")
	SlotBorder       = ReadEbImgFS(fs, "assets/img/gui/slot_border.png")
	cracks           = ImgFromFS(fs, "assets/img/cracks.png")
	Font             = LoadFontFromFS("assets/font/arkpixel10.ttf", 10, fs)
)

var (
	Player                    = ReadEbImgFS(fs, "assets/img/player/player.png")
	PlayerWoodenAxeAtlas      = ReadEbImgFS(fs, "assets/img/player/player_wood_axe.png")
	PlayerStoneAxeAtlas       = ReadEbImgFS(fs, "assets/img/player/player_stone_axe.png")
	PlayerIronAxeAtlas        = ReadEbImgFS(fs, "assets/img/player/player_iron_axe.png")
	PlayerDiamondAxeAtlas     = ReadEbImgFS(fs, "assets/img/player/player_diamond_axe.png")
	PlayerWoodenPickaxeAtlas  = ReadEbImgFS(fs, "assets/img/player/player_wood_pickaxe.png")
	PlayerStonePickaxeAtlas   = ReadEbImgFS(fs, "assets/img/player/player_stone_pickaxe.png")
	PlayerIronPickaxeAtlas    = ReadEbImgFS(fs, "assets/img/player/player_iron_pickaxe.png")
	PlayerDiamondPickaxeAtlas = ReadEbImgFS(fs, "assets/img/player/player_diamond_pickaxe.png")
	PlayerWoodenShovelAtlas   = ReadEbImgFS(fs, "assets/img/player/player_wood_shovel.png")
	PlayerStoneShovelAtlas    = ReadEbImgFS(fs, "assets/img/player/player_stone_shovel.png")
	PlayerIronShovelAtlas     = ReadEbImgFS(fs, "assets/img/player/player_iron_shovel.png")
	PlayerDiamondShovelAtlas  = ReadEbImgFS(fs, "assets/img/player/player_diamond_shovel.png")
)

func init() {

	// Unbreakable blocks
	BlockUnbreakable[items.Bedrock] = ReadEbImgFS(fs, "assets/img/blocks/bedrock.png")
	BlockUnbreakable[items.Random] = ReadEbImgFS(fs, "assets/img/blocks/random.png")

	// Breakable blocks
	BlockCrackFrames[items.CoalOre] = blockImgs("coal_ore.png")
	BlockCrackFrames[items.CraftingTable] = blockImgs("crafting_table.png")
	BlockCrackFrames[items.DiamondOre] = blockImgs("diamond_ore.png")
	BlockCrackFrames[items.Dirt] = blockImgs("dirt.png")
	BlockCrackFrames[items.Furnace] = blockImgs("furnace.png")
	BlockCrackFrames[items.GoldOre] = blockImgs("gold_ore.png")
	BlockCrackFrames[items.GrassBlock] = blockImgs("grass_block.png")
	BlockCrackFrames[items.GrassBlockSnow] = blockImgs("grass_block_snow.png")
	BlockCrackFrames[items.IronOre] = blockImgs("iron_ore.png")
	BlockCrackFrames[items.OakLeaves] = blockImgs("oak_leaves.png")
	BlockCrackFrames[items.OakLog] = blockImgs("oak_log.png")
	BlockCrackFrames[items.OakPlanks] = blockImgs("oak_planks.png")
	BlockCrackFrames[items.OakSapling] = blockImgs("oak_sapling.png")
	BlockCrackFrames[items.Obsidian] = blockImgs("obsidian.png")
	BlockCrackFrames[items.Sand] = blockImgs("sand.png")
	BlockCrackFrames[items.SmoothStone] = blockImgs("smooth_stone.png")
	BlockCrackFrames[items.Snow] = blockImgs("snow.png")
	BlockCrackFrames[items.Stone] = blockImgs("stone.png")
	BlockCrackFrames[items.StoneBricks] = blockImgs("stone_bricks.png")
	BlockCrackFrames[items.Tnt] = blockImgs("tnt.png")
	BlockCrackFrames[items.Torch] = blockImgs("torch.png")

	// Block icons
	Icon8[items.Bedrock] = blockIconImg("bedrock.png")
	Icon8[items.Random] = blockIconImg("random.png")
	Icon8[items.CoalOre] = blockIconImg("coal_ore.png")
	Icon8[items.CraftingTable] = blockIconImg("crafting_table.png")
	Icon8[items.DiamondOre] = blockIconImg("diamond_ore.png")
	Icon8[items.Dirt] = blockIconImg("dirt.png")
	Icon8[items.Furnace] = blockIconImg("furnace.png")
	Icon8[items.GoldOre] = blockIconImg("gold_ore.png")
	Icon8[items.GrassBlock] = blockIconImg("grass_block.png")
	Icon8[items.GrassBlockSnow] = blockIconImg("grass_block_snow.png")
	Icon8[items.IronOre] = blockIconImg("iron_ore.png")
	Icon8[items.OakLeaves] = blockIconImg("oak_leaves.png")
	Icon8[items.OakLog] = blockIconImg("oak_log.png")
	Icon8[items.OakPlanks] = blockIconImg("oak_planks.png")
	Icon8[items.OakSapling] = blockIconImg("oak_sapling.png")
	Icon8[items.Obsidian] = blockIconImg("obsidian.png")
	Icon8[items.Sand] = blockIconImg("sand.png")
	Icon8[items.SmoothStone] = blockIconImg("smooth_stone.png")
	Icon8[items.Snow] = blockIconImg("snow.png")
	Icon8[items.Stone] = blockIconImg("stone.png")
	Icon8[items.StoneBricks] = blockIconImg("stone_bricks.png")
	Icon8[items.Tnt] = blockIconImg("tnt.png")
	Icon8[items.Torch] = blockIconImg("torch.png")

	// Item icons
	Icon8[items.Bread] = itemIconImg("bread.png")
	Icon8[items.Bucket] = itemIconImg("bucket.png")
	Icon8[items.Coal] = itemIconImg("coal.png")
	Icon8[items.Diamond] = itemIconImg("diamond.png")
	Icon8[items.DiamondAxe] = itemIconImg("diamond_axe.png")
	Icon8[items.DiamondPickaxe] = itemIconImg("diamond_pickaxe.png")
	Icon8[items.DiamondShovel] = itemIconImg("diamond_shovel.png")
	Icon8[items.GoldIngot] = itemIconImg("gold_ingot.png")
	Icon8[items.IronAxe] = itemIconImg("iron_axe.png")
	Icon8[items.IronIngot] = itemIconImg("iron_ingot.png")
	Icon8[items.IronPickaxe] = itemIconImg("iron_pickaxe.png")
	Icon8[items.IronShovel] = itemIconImg("iron_shovel.png")
	Icon8[items.RawGold] = itemIconImg("raw_gold.png")
	Icon8[items.RawIron] = itemIconImg("raw_iron.png")
	Icon8[items.Snowball] = itemIconImg("snowball.png")
	Icon8[items.Stick] = itemIconImg("stick.png")
	Icon8[items.StoneAxe] = itemIconImg("stone_axe.png")
	Icon8[items.StonePickaxe] = itemIconImg("stone_pickaxe.png")
	Icon8[items.StoneShovel] = itemIconImg("stone_shovel.png")
	Icon8[items.WaterBucket] = itemIconImg("water_bucket.png")
	Icon8[items.WoodenAxe] = itemIconImg("wooden_axe.png")
	Icon8[items.WoodenPickaxe] = itemIconImg("wooden_pickaxe.png")
	Icon8[items.WoodenShovel] = itemIconImg("wooden_shovel.png")
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
		x := i * 20
		rec := image.Rect(x, 0, x+20, x+20)
		si := stages.(*image.NRGBA).SubImage(rec)
		frames = append(frames, blend.Normal(block, si))
	}
	return frames
}
func blockImgs(f string) []*ebiten.Image {
	frames := makeStages(ImgFromFS(fs, "assets/img/blocks/"+f), cracks)
	return toEbiten(frames)
}
func itemIconImg(f string) *ebiten.Image {
	return ReadEbImgFS(fs, "assets/img/items_icon/"+f)
}
func blockIconImg(f string) *ebiten.Image {
	return ReadEbImgFS(fs, "assets/img/blocks_icon/"+f)
}

func ReadEbImgFS(fs embed.FS, filePath string) *ebiten.Image {
	return ebiten.NewImageFromImage(ImgFromFS(fs, filePath))
}

func ImgFromFS(fs embed.FS, filePath string) image.Image {
	f, err := fs.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return image
}

func LoadFontFromFS(file string, size float64, fs embed.FS) *text.GoTextFace {
	f, err := fs.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	src, err := text.NewGoTextFaceSource(bytes.NewReader(f))
	if err != nil {
		log.Fatal(err)
	}
	gtf := &text.GoTextFace{
		Source:   src,
		Size:     size,
		Language: language.English,
	}
	return gtf
}
