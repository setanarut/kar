package itm

import (
	"image/color"
	"kar/engine/util"
)

type Category uint8

// Item Bitmask Category
const (
	CatNone  Category = 0
	CatBlock Category = 1 << iota
	CatOreBlock
	CatDropItem
	CatItem
	CatRawOre
	CatTool
	CatWeapon
	CatAll Category = 255
)

type Item struct {
	DisplayName string
	Drops       uint16
	Stackable   uint16
	Breakable   bool
	MaxHealth   float64
	Category    Category
}

const (
	Air uint16 = iota
	Grass
	Snow
	Dirt
	Sand
	Cobblestone
	CobbledDeepslate

	Stone
	CoalOre
	GoldOre
	IronOre
	DiamondOre
	CopperOre
	EmeraldOre
	LapisOre
	RedstoneOre

	DeepslateStone
	DeepslateCoalOre
	DeepslateGoldOre
	DeepslateIronOre
	DeepslateDiamondOre
	DeepslateCopperOre
	DeepslateEmeraldOre
	DeepslateLapisOre
	DeepslateRedStoneOre
	Bedrock

	Log
	Leaves
	Planks
	Sapling
	Torch

	Coal
	CharCoal
	RawGold
	RawIron
	Diamond
	RawCopper
	Emerald
	LapisLazuli
	Redstone
	WoodAxe
	WoodHoe
	WoodPickaxe
	WoodShovel
	WoodSword
	StoneAxe
	StoneHoe
	StonePickaxe
	StoneShovel
	StoneSword
	GoldenAxe
	GoldenHoe
	GoldenPickaxe
	GoldenShovel
	GoldenSword
	IronAxe
	IronHoe
	IronPickaxe
	IronShovel
	IronSword
	DiamondAxe
	DiamondHoe
	DiamondPickaxe
	DiamondShovel
	DiamondSword
	NetheriteAxe
	NetheriteHoe
	NetheritePickaxe
	NetheriteShovel
	NetheriteSword
	NetheriteScrap
	NetheriteIngot
	GoldIngot
	IronIngot
	CopperIngot
	Stick
)

var Items = map[uint16]Item{
	Air: {
		DisplayName: "Air",
		Drops:       Air,
		Stackable:   0,
		Breakable:   false,
		MaxHealth:   0,
		Category:    CatNone,
	},
	Grass: {
		DisplayName: "Grass",
		Drops:       Dirt,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock,
	},
	Snow: {
		DisplayName: "Snow",
		Drops:       Dirt,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock,
	},
	Dirt: {
		DisplayName: "Dirt",
		Drops:       Dirt,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock,
	},
	Sand: {
		DisplayName: "Sand",
		Drops:       Sand,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock,
	},
	Stone: {
		DisplayName: "Stone",
		Drops:       Cobblestone,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock,
	},
	Cobblestone: {
		DisplayName: "Cobblestone",
		Drops:       Cobblestone,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock,
	},
	CoalOre: {
		DisplayName: "Coal Ore",
		Drops:       Coal,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	GoldOre: {
		DisplayName: "Gold Ore",
		Drops:       RawGold,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	IronOre: {
		DisplayName: "Iron Ore",
		Drops:       RawIron,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	DiamondOre: {
		DisplayName: "Diamond Ore",
		Drops:       Diamond,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	CopperOre: {
		DisplayName: "Copper Ore",
		Drops:       RawCopper,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	EmeraldOre: {
		DisplayName: "Emerald Ore",
		Drops:       Emerald,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	LapisOre: {
		DisplayName: "Lapis Ore",
		Drops:       LapisLazuli,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	RedstoneOre: {
		DisplayName: "Redstone Ore",
		Drops:       Redstone,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	//DEEPSLATE
	DeepslateStone: {
		DisplayName: "Deepslate Stone",
		Drops:       CobbledDeepslate,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock,
	},
	CobbledDeepslate: {
		DisplayName: "Cobbled Deepslate",
		Drops:       CobbledDeepslate,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock,
	},
	DeepslateCoalOre: {
		DisplayName: "Deepslate Coal Ore",
		Drops:       Coal,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	DeepslateGoldOre: {
		DisplayName: "Deepslate Gold Ore",
		Drops:       RawGold,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	DeepslateIronOre: {
		DisplayName: "Deepslate Iron Ore",
		Drops:       RawIron,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	DeepslateDiamondOre: {
		DisplayName: "Deepslate Diamond Ore",
		Drops:       Diamond,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	DeepslateCopperOre: {
		DisplayName: "Deepslate Copper Ore",
		Drops:       RawCopper,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	DeepslateEmeraldOre: {
		DisplayName: "Deepslate Emerald Ore",
		Drops:       Emerald,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},

	DeepslateLapisOre: {
		DisplayName: "Deepslate Lapis Ore",
		Drops:       LapisLazuli,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	DeepslateRedStoneOre: {
		DisplayName: "Deepslate Red Stone Ore",
		Drops:       Redstone,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock | CatOreBlock,
	},
	Bedrock: {
		DisplayName: "Bedrock",
		Drops:       Air,
		Stackable:   0,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatBlock,
	},
	Log: {
		DisplayName: "Log",
		Drops:       Log,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock,
	},

	Leaves: {
		DisplayName: "Leaves",
		Drops:       Leaves,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock,
	},

	Planks: {
		DisplayName: "Tree Plank",
		Drops:       Planks,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
		Category:    CatBlock,
	},
	Sapling: {
		DisplayName: "Sapling",
		Drops:       Sapling,
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatBlock | CatItem,
	},
	Torch: {
		DisplayName: "Torch",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   10,
		Category:    CatItem,
	},

	Coal: {
		DisplayName: "Coal",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatDropItem | CatRawOre,
	},
	CharCoal: {
		DisplayName: "CharCoal",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatDropItem | CatRawOre,
	},
	RawGold: {
		DisplayName: "Raw Gold",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatDropItem | CatRawOre,
	},

	RawIron: {
		DisplayName: "Raw Iron",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatDropItem | CatRawOre,
	},

	Diamond: {
		DisplayName: "Diamond",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatDropItem | CatRawOre,
	},
	RawCopper: {
		DisplayName: "Raw Copper",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatDropItem | CatRawOre,
	},

	Emerald: {
		DisplayName: "Emerald",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatDropItem | CatRawOre,
	},
	LapisLazuli: {
		DisplayName: "Lapis Lazuli",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatDropItem | CatRawOre,
	},

	Redstone: {
		DisplayName: "Redstone",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatDropItem | CatRawOre,
	},

	WoodAxe: {
		DisplayName: "Wood Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	WoodHoe: {
		DisplayName: "Wood Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	WoodPickaxe: {
		DisplayName: "Wood Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	WoodShovel: {
		DisplayName: "Wood Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	WoodSword: {
		DisplayName: "Wood Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatWeapon,
	},
	StoneAxe: {
		DisplayName: "Stone Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	StoneHoe: {
		DisplayName: "Stone Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	StonePickaxe: {
		DisplayName: "Stone Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	StoneShovel: {
		DisplayName: "Stone Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	StoneSword: {
		DisplayName: "Stone Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatWeapon,
	},
	GoldenAxe: {
		DisplayName: "Golden Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	GoldenHoe: {
		DisplayName: "Golden Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	GoldenPickaxe: {
		DisplayName: "Golden Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	GoldenShovel: {
		DisplayName: "Golden Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	GoldenSword: {
		DisplayName: "Golden Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatWeapon,
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
		Category:    CatItem | CatTool,
	},
	IronPickaxe: {
		DisplayName: "Iron Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	IronShovel: {
		DisplayName: "Iron Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	IronSword: {
		DisplayName: "Iron Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatWeapon,
	},
	DiamondAxe: {
		DisplayName: "Diamond Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	DiamondHoe: {
		DisplayName: "Diamond Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	DiamondPickaxe: {
		DisplayName: "Diamond Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	DiamondShovel: {
		DisplayName: "Diamond Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	DiamondSword: {
		DisplayName: "Diamond Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatWeapon,
	},
	NetheriteAxe: {
		DisplayName: "Netherite Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	NetheriteHoe: {
		DisplayName: "Netherite Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	NetheritePickaxe: {
		DisplayName: "Netherite Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	NetheriteShovel: {
		DisplayName: "Netherite Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatTool,
	},
	NetheriteSword: {
		DisplayName: "Netherite Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatWeapon,
	},
	NetheriteScrap: {
		DisplayName: "Netherite Scrap",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem | CatRawOre | CatDropItem,
	},
	NetheriteIngot: {
		DisplayName: "Netherite Ingot",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem,
	},
	GoldIngot: {
		DisplayName: "Gold Ingot",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem,
	},
	IronIngot: {
		DisplayName: "Iron Ingot",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem,
	},
	CopperIngot: {
		DisplayName: "Copper Ingot",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem,
	},
	Stick: {
		DisplayName: "Stick",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
		Category:    CatItem,
	},
}

var ItemColorMap = map[uint16]color.RGBA{
	Air:                 util.HexToRGBA("#0099ff"),
	Grass:               util.HexToRGBA("#00903f"),
	Dirt:                util.HexToRGBA("#74573E"),
	Sand:                util.HexToRGBA("#fff5cc"),
	Stone:               util.HexToRGBA("#949494"),
	CoalOre:             util.HexToRGBA("#372f2f"),
	GoldOre:             util.HexToRGBA("#ffe100"),
	IronOre:             util.HexToRGBA("#b8947d"),
	DiamondOre:          util.HexToRGBA("#40efd4"),
	DeepslateStone:      util.HexToRGBA("#4c4c4c"),
	DeepslateCoalOre:    util.HexToRGBA("#29344e"),
	DeepslateGoldOre:    util.HexToRGBA("#ffe100"),
	DeepslateIronOre:    util.HexToRGBA("#8a6548"),
	DeepslateDiamondOre: util.HexToRGBA("#00ffe1"),
}
