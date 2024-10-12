package items

import (
	"image/color"
	"kar/engine/util"
)

type Category uint8

// Item Bitmask Category
const (
	CategoryNone  Category = 0
	CategoryBlock Category = 1 << iota
	CategoryOreBlock
	CategoryDropItem
	CategoryItem
	CategoryRawOre
	CategoryTool
	CategoryWeapon
	CategoryAll Category = 255
)

type ItemProp struct {
	DisplayName string
	Drops       uint16
	Stackable   uint16
	Breakable   bool
	MaxHealth   float64
	Category    Category
}

const (
	Air uint16 = iota
	Andesite
	Arrow
	Bedrock
	BeetrootSeeds
	Bow
	Bread
	BrewingStand
	Bucket
	CartographyTable
	Charcoal
	Clay
	Coal
	CoalBlock
	CoalOre
	CoarseDirt
	CobbledDeepslate
	Cobblestone
	CopperIngot
	CopperOre
	CraftingTable
	CrossbowStandby
	CryingObsidian
	Deepslate
	DeepslateCoalOre
	DeepslateCopperOre
	DeepslateDiamondOre
	DeepslateEmeraldOre
	DeepslateGoldOre
	DeepslateIronOre
	DeepslateLapisOre
	DeepslateRedstoneOre
	Diamond
	DiamondAxe
	DiamondHoe
	DiamondOre
	DiamondPickaxe
	DiamondShovel
	DiamondSword
	Dirt
	DirtPath
	Emerald
	EmeraldOre
	EnchantingTable
	EndPortalFrame
	FletchingTable
	Furnace
	FurnaceOn
	GoldIngot
	GoldOre
	GoldenAxe
	GoldenHoe
	GoldenPickaxe
	GoldenShovel
	GoldenSword
	GrassBlock
	GrassBlockSnow
	Gravel
	IronAxe
	IronHoe
	IronIngot
	IronOre
	IronPickaxe
	IronShovel
	IronSword
	LapisLazuli
	LapisOre
	LavaBucket
	MelonSeeds
	MilkBucket
	NetherBricks
	NetherGoldOre
	NetherQuartzOre
	NetheriteAxe
	NetheriteHoe
	NetheriteIngot
	NetheritePickaxe
	NetheriteScrap
	NetheriteShovel
	NetheriteSword
	Netherrack
	OakLeaves
	OakLog
	OakPlanks
	OakSapling
	Obsidian
	PowderSnowBucket
	PumpkinSeeds
	RawCopper
	RawGold
	RawIron
	RedNetherBricks
	RedSand
	RedSandstone
	Redstone
	RedstoneOre
	RedstoneTorch
	RedstoneTorchOff
	RootedDirt
	Sand
	Sandstone
	SmithingTable
	SmoothStone
	Snow
	Snowball
	SoulSand
	SoulSoil
	SoulTorch
	Stick
	Stone
	StoneAxe
	StoneBricks
	StoneHoe
	StonePickaxe
	StoneShovel
	StoneSword
	Tnt
	Torch
	TorchflowerSeeds
	WaterBucket
	Wheat
	WheatCrops
	WheatSeeds
	WoodenAxe
	WoodenHoe
	WoodenPickaxe
	WoodenShovel
	WoodenSword
)

var Property = map[uint16]ItemProp{
	Air: {
		DisplayName: "Air",
		Drops:       Air,
		Stackable:   0,
		Breakable:   false,
		MaxHealth:   0,
		Category:    CategoryNone,
	},
	GrassBlock: {
		DisplayName: "Grass",
		Drops:       Dirt,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock,
	},
	Snow: {
		DisplayName: "Snow",
		Drops:       Dirt,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock,
	},
	Dirt: {
		DisplayName: "Dirt",
		Drops:       Dirt,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock,
	},
	Sand: {
		DisplayName: "Sand",
		Drops:       Sand,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock,
	},
	Stone: {
		DisplayName: "Stone",
		Drops:       Cobblestone,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock,
	},
	Cobblestone: {
		DisplayName: "Cobblestone",
		Drops:       Cobblestone,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock,
	},
	CoalOre: {
		DisplayName: "Coal Ore",
		Drops:       Coal,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	GoldOre: {
		DisplayName: "Gold Ore",
		Drops:       RawGold,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	IronOre: {
		DisplayName: "Iron Ore",
		Drops:       RawIron,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	DiamondOre: {
		DisplayName: "Diamond Ore",
		Drops:       Diamond,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	CopperOre: {
		DisplayName: "Copper Ore",
		Drops:       RawCopper,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	EmeraldOre: {
		DisplayName: "Emerald Ore",
		Drops:       Emerald,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	LapisOre: {
		DisplayName: "Lapis Ore",
		Drops:       LapisLazuli,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	RedstoneOre: {
		DisplayName: "Redstone Ore",
		Drops:       Redstone,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},

	Deepslate: {
		DisplayName: "Deepslate Stone",
		Drops:       CobbledDeepslate,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock,
	},
	CobbledDeepslate: {
		DisplayName: "Cobbled Deepslate",
		Drops:       CobbledDeepslate,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock,
	},
	DeepslateCoalOre: {
		DisplayName: "Deepslate Coal Ore",
		Drops:       Coal,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	DeepslateGoldOre: {
		DisplayName: "Deepslate Gold Ore",
		Drops:       RawGold,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	DeepslateIronOre: {
		DisplayName: "Deepslate Iron Ore",
		Drops:       RawIron,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	DeepslateDiamondOre: {
		DisplayName: "Deepslate Diamond Ore",
		Drops:       Diamond,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	DeepslateCopperOre: {
		DisplayName: "Deepslate Copper Ore",
		Drops:       RawCopper,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	DeepslateEmeraldOre: {
		DisplayName: "Deepslate Emerald Ore",
		Drops:       Emerald,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},

	DeepslateLapisOre: {
		DisplayName: "Deepslate Lapis Ore",
		Drops:       LapisLazuli,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	DeepslateRedstoneOre: {
		DisplayName: "Deepslate Redstone Ore",
		Drops:       Redstone,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock | CategoryOreBlock,
	},
	Bedrock: {
		DisplayName: "Bedrock",
		Drops:       Air,
		Stackable:   0,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryBlock,
	},
	OakLog: {
		DisplayName: "Log",
		Drops:       OakLog,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock,
	},

	OakLeaves: {
		DisplayName: "Leaves",
		Drops:       OakLeaves,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock,
	},

	OakPlanks: {
		DisplayName: "Tree Plank",
		Drops:       OakPlanks,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CategoryBlock,
	},
	OakSapling: {
		DisplayName: "Sapling",
		Drops:       OakSapling,
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryBlock | CategoryItem,
	},
	Torch: {
		DisplayName: "Torch",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   10,
		Category:    CategoryItem,
	},

	Coal: {
		DisplayName: "Coal",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryDropItem | CategoryRawOre,
	},
	Charcoal: {
		DisplayName: "CharCoal",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryDropItem | CategoryRawOre,
	},
	RawGold: {
		DisplayName: "Raw Gold",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryDropItem | CategoryRawOre,
	},

	RawIron: {
		DisplayName: "Raw Iron",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryDropItem | CategoryRawOre,
	},

	Diamond: {
		DisplayName: "Diamond",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryDropItem | CategoryRawOre,
	},
	RawCopper: {
		DisplayName: "Raw Copper",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryDropItem | CategoryRawOre,
	},

	Emerald: {
		DisplayName: "Emerald",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryDropItem | CategoryRawOre,
	},
	LapisLazuli: {
		DisplayName: "Lapis Lazuli",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryDropItem | CategoryRawOre,
	},

	Redstone: {
		DisplayName: "Redstone",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryDropItem | CategoryRawOre,
	},

	WoodenAxe: {
		DisplayName: "Wooden Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	WoodenHoe: {
		DisplayName: "Wooden Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	WoodenPickaxe: {
		DisplayName: "Wooden Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	WoodenShovel: {
		DisplayName: "Wooden Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	WoodenSword: {
		DisplayName: "Wooden Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryWeapon,
	},
	StoneAxe: {
		DisplayName: "Stone Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	StoneHoe: {
		DisplayName: "Stone Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	StonePickaxe: {
		DisplayName: "Stone Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	StoneShovel: {
		DisplayName: "Stone Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	StoneSword: {
		DisplayName: "Stone Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryWeapon,
	},
	GoldenAxe: {
		DisplayName: "Golden Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	GoldenHoe: {
		DisplayName: "Golden Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	GoldenPickaxe: {
		DisplayName: "Golden Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	GoldenShovel: {
		DisplayName: "Golden Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	GoldenSword: {
		DisplayName: "Golden Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryWeapon,
	},
	IronAxe: {
		DisplayName: "Iron Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	IronHoe: {
		DisplayName: "Iron Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	IronPickaxe: {
		DisplayName: "Iron Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	IronShovel: {
		DisplayName: "Iron Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	IronSword: {
		DisplayName: "Iron Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryWeapon,
	},
	DiamondAxe: {
		DisplayName: "Diamond Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	DiamondHoe: {
		DisplayName: "Diamond Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	DiamondPickaxe: {
		DisplayName: "Diamond Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	DiamondShovel: {
		DisplayName: "Diamond Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	DiamondSword: {
		DisplayName: "Diamond Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryWeapon,
	},
	NetheriteAxe: {
		DisplayName: "Netherite Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	NetheriteHoe: {
		DisplayName: "Netherite Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	NetheritePickaxe: {
		DisplayName: "Netherite Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	NetheriteShovel: {
		DisplayName: "Netherite Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryTool,
	},
	NetheriteSword: {
		DisplayName: "Netherite Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryWeapon,
	},
	NetheriteScrap: {
		DisplayName: "Netherite Scrap",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem | CategoryRawOre | CategoryDropItem,
	},
	NetheriteIngot: {
		DisplayName: "Netherite Ingot",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem,
	},
	GoldIngot: {
		DisplayName: "Gold Ingot",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem,
	},
	IronIngot: {
		DisplayName: "Iron Ingot",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem,
	},
	CopperIngot: {
		DisplayName: "Copper Ingot",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem,
	},
	Stick: {
		DisplayName: "Stick",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CategoryItem,
	},
}

var ItemColorMap = map[uint16]color.RGBA{
	Air:                 util.HexToRGBA("#0099ff"),
	GrassBlock:          util.HexToRGBA("#00903f"),
	Dirt:                util.HexToRGBA("#74573E"),
	Sand:                util.HexToRGBA("#fff5cc"),
	Stone:               util.HexToRGBA("#949494"),
	CoalOre:             util.HexToRGBA("#372f2f"),
	GoldOre:             util.HexToRGBA("#ffe100"),
	IronOre:             util.HexToRGBA("#b8947d"),
	DiamondOre:          util.HexToRGBA("#40efd4"),
	Deepslate:           util.HexToRGBA("#4c4c4c"),
	DeepslateCoalOre:    util.HexToRGBA("#29344e"),
	DeepslateGoldOre:    util.HexToRGBA("#ffe100"),
	DeepslateIronOre:    util.HexToRGBA("#8a6548"),
	DeepslateDiamondOre: util.HexToRGBA("#00ffe1"),
}
