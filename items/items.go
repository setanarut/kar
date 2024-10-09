package items

import (
	"image/color"
	"kar/engine/util"
)

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

type Item struct {
	DisplayName string
	Drops       uint16
	Stackable   uint16
	Breakable   bool
	MaxHealth   float64
}

var Items = map[uint16]Item{
	Air: {
		DisplayName: "Air",
		Drops:       Air,
		Stackable:   0,
		Breakable:   false,
		MaxHealth:   0,
	},
	Grass: {
		DisplayName: "Grass",
		Drops:       Dirt,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	Snow: {
		DisplayName: "Snow",
		Drops:       Dirt,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	Dirt: {
		DisplayName: "Dirt",
		Drops:       Dirt,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	Sand: {
		DisplayName: "Sand",
		Drops:       Sand,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	Stone: {
		DisplayName: "Stone",
		Drops:       Cobblestone,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	Cobblestone: {
		DisplayName: "Cobblestone",
		Drops:       Cobblestone,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	CoalOre: {
		DisplayName: "Coal Ore",
		Drops:       Coal,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	GoldOre: {
		DisplayName: "Gold Ore",
		Drops:       RawGold,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	IronOre: {
		DisplayName: "Iron Ore",
		Drops:       RawIron,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	DiamondOre: {
		DisplayName: "Diamond Ore",
		Drops:       Diamond,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	CopperOre: {
		DisplayName: "Copper Ore",
		Drops:       RawCopper,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	EmeraldOre: {
		DisplayName: "Emerald Ore",
		Drops:       Emerald,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	LapisOre: {
		DisplayName: "Lapis Ore",
		Drops:       LapisLazuli,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	RedstoneOre: {
		DisplayName: "Redstone Ore",
		Drops:       Redstone,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	//DEEPSLATE
	DeepslateStone: {
		DisplayName: "Deepslate Stone",
		Drops:       CobbledDeepslate,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	DeepslateCoalOre: {
		DisplayName: "Deepslate Coal Ore",
		Drops:       Coal,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	DeepslateGoldOre: {
		DisplayName: "Deepslate Gold Ore",
		Drops:       RawGold,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	DeepslateIronOre: {
		DisplayName: "Deepslate Iron Ore",
		Drops:       RawIron,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	DeepslateDiamondOre: {
		DisplayName: "Deepslate Diamond Ore",
		Drops:       Diamond,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	DeepslateCopperOre: {
		DisplayName: "Deepslate Copper Ore",
		Drops:       RawCopper,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	DeepslateEmeraldOre: {
		DisplayName: "Deepslate Emerald Ore",
		Drops:       Emerald,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},

	DeepslateLapisOre: {
		DisplayName: "Deepslate Lapis Ore",
		Drops:       LapisLazuli,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	DeepslateRedStoneOre: {
		DisplayName: "Deepslate Red Stone Ore",
		Drops:       Redstone,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	Log: {
		DisplayName: "Log",
		Drops:       Log,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},

	Leaves: {
		DisplayName: "Leaves",
		Drops:       Leaves,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},

	Planks: {
		DisplayName: "Tree Plank",
		Drops:       Planks,
		Stackable:   64,
		Breakable:   true,
		MaxHealth:   10,
	},
	Sapling: {
		DisplayName: "Sapling",
		Drops:       Sapling,
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   10,
	},
	Torch: {
		DisplayName: "Torch",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   10,
	},

	Coal: {
		DisplayName: "Coal",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},
	CharCoal: {
		DisplayName: "CharCoal",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},
	RawGold: {
		DisplayName: "Raw Gold",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},

	RawIron: {
		DisplayName: "Raw Iron",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},

	Diamond: {
		DisplayName: "Diamond",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},
	RawCopper: {
		DisplayName: "Raw Copper",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},

	Emerald: {
		DisplayName: "Emerald",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},
	LapisLazuli: {
		DisplayName: "Lapis Lazuli",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},

	Redstone: {
		DisplayName: "Redstone",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},

	WoodAxe: {
		DisplayName: "Wood Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	WoodHoe: {
		DisplayName: "Wood Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	WoodPickaxe: {
		DisplayName: "Wood Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	WoodShovel: {
		DisplayName: "Wood Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	WoodSword: {
		DisplayName: "Wood Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	StoneAxe: {
		DisplayName: "Stone Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	StoneHoe: {
		DisplayName: "Stone Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	StonePickaxe: {
		DisplayName: "Stone Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	StoneShovel: {
		DisplayName: "Stone Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	StoneSword: {
		DisplayName: "Stone Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	GoldenAxe: {
		DisplayName: "Golden Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	GoldenHoe: {
		DisplayName: "Golden Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	GoldenPickaxe: {
		DisplayName: "Golden Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	GoldenShovel: {
		DisplayName: "Golden Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	GoldenSword: {
		DisplayName: "Golden Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
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
	},
	IronPickaxe: {
		DisplayName: "Iron Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	IronShovel: {
		DisplayName: "Iron Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	IronSword: {
		DisplayName: "Iron Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	DiamondAxe: {
		DisplayName: "Diamond Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	DiamondHoe: {
		DisplayName: "Diamond Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	DiamondPickaxe: {
		DisplayName: "Diamond Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	DiamondShovel: {
		DisplayName: "Diamond Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	DiamondSword: {
		DisplayName: "Diamond Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	NetheriteAxe: {
		DisplayName: "Netherite Axe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	NetheriteHoe: {
		DisplayName: "Netherite Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	NetheritePickaxe: {
		DisplayName: "Netherite Pickaxe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	NetheriteShovel: {
		DisplayName: "Netherite Hoe",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	NetheriteSword: {
		DisplayName: "Netherite Sword",
		Stackable:   1,
		Breakable:   false,
		MaxHealth:   1,
	},
	NetheriteScrap: {
		DisplayName: "Netherite Scrap",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},
	NetheriteIngot: {
		DisplayName: "Netherite Ingot",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},
	GoldIngot: {
		DisplayName: "Gold Ingot",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},
	IronIngot: {
		DisplayName: "Iron Ingot",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},
	CopperIngot: {
		DisplayName: "Copper Ingot",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},
	Stick: {
		DisplayName: "Stick",
		Stackable:   64,
		Breakable:   false,
		MaxHealth:   1,
	},
}

var ItemColorMap = map[uint16]color.RGBA{
	Air:                 util.HexToRGBA("#0099ff"),
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
	Grass:               util.HexToRGBA("#00903f"),
}
