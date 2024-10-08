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

	Tree
	TreeLeaves
	TreePlank
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

	WoodShovel
	StoneShovel
	IronShovel
	GoldenAxe
	WoodAxe
	StoneAxe
	IronAxe
	DiamondAxe
	NetheriteAxe
	GoldenPickaxe
	WoodPickaxe
	StonePickaxe
	IronPickaxe
	DiamondPickaxe
	NetheritePickaxe
	GoldenSword
	WoodSword
	StoneSword
	IronSword
	DiamondSword
	NetheriteSword
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
		DisplayName: "Grass",
		Drops:       5,
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
