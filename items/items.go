package items

import (
	"image/color"
	"kar/engine/util"
	"kar/types"
)

const (
	Air types.ItemID = iota
	// Breakable blocks
	Grass
	Snow
	Dirt
	Sand
	Stone
	CoalOre
	GoldOre
	IronOre
	DiamondOre
	CopperOre
	EmeraldOre
	LapisOre
	RedStoneOre
	DeepSlateStone
	DeepSlateCoalOre
	DeepSlateGoldOre
	DeepSlateIronOre
	DeepSlateDiamondOre
	DeepSlateCopperOre
	DeepSlateEmeraldOre
	DeepSlateLapisOre
	DeepSlateRedStoneOre
	Tree
	TreeLeaves
	TreePlank

	// items
	Sapling
	Torch
	Coal
	RawGold
	RawIron
	Diamond
	RawCopper

	// Tools
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

	// Weapons
	GoldenSword
	WoodSword
	StoneSword
	IronSword
	DiamondSword
	NetheriteSword
)

var DisplayName = map[types.ItemID]string{
	Air:                  "Air",
	Grass:                "Grass",
	Snow:                 "Snow",
	Dirt:                 "Dirt",
	Sand:                 "Sand",
	Stone:                "Stone",
	CoalOre:              "CoalOre",
	GoldOre:              "GoldOre",
	IronOre:              "IronOre",
	DiamondOre:           "DiamondOre",
	CopperOre:            "CopperOre",
	EmeraldOre:           "EmeraldOre",
	LapisOre:             "LapisOre",
	RedStoneOre:          "RedStoneOre",
	DeepSlateStone:       "DeepSlateStone",
	DeepSlateCoalOre:     "DeepSlateCoalOre",
	DeepSlateGoldOre:     "DeepSlateGoldOre",
	DeepSlateIronOre:     "DeepSlateIronOre",
	DeepSlateDiamondOre:  "DeepSlateDiamondOre",
	DeepSlateCopperOre:   "DeepSlateCopperOre",
	DeepSlateEmeraldOre:  "DeepSlateEmeraldOre",
	DeepSlateLapisOre:    "DeepSlateLapisOre",
	DeepSlateRedStoneOre: "DeepSlateRedStoneOre",
	Tree:                 "Tree",
	TreeLeaves:           "TreeLeaves",
	TreePlank:            "TreePlank",

	// items
	Sapling:   "Sapling",
	Torch:     "Torch",
	Coal:      "Coal",
	RawGold:   "RawGold",
	RawIron:   "RawIron",
	Diamond:   "Diamond",
	RawCopper: "RawCopper",

	// Tools
	WoodShovel:       "WoodShovel",
	StoneShovel:      "StoneShovel",
	IronShovel:       "IronShovel",
	GoldenAxe:        "GoldenAxe",
	WoodAxe:          "WoodAxe",
	StoneAxe:         "StoneAxe",
	IronAxe:          "IronAxe",
	DiamondAxe:       "DiamondAxe",
	NetheriteAxe:     "NetheriteAxe",
	GoldenPickaxe:    "GoldenPickaxe",
	WoodPickaxe:      "WoodPickaxe",
	StonePickaxe:     "StonePickaxe",
	IronPickaxe:      "IronPickaxe",
	DiamondPickaxe:   "DiamondPickaxe",
	NetheritePickaxe: "NetheritePickaxe",

	// Weapons
	GoldenSword:    "GoldenSword",
	WoodSword:      "WoodSword",
	StoneSword:     "StoneSword",
	IronSword:      "IronSword",
	DiamondSword:   "DiamondSword",
	NetheriteSword: "NetheriteSword",
}

var Colors = map[types.ItemID]color.RGBA{
	Air:                 util.HexToRGBA("#0099ff"),
	Dirt:                util.HexToRGBA("#74573E"),
	Sand:                util.HexToRGBA("#fff5cc"),
	Stone:               util.HexToRGBA("#949494"),
	CoalOre:             util.HexToRGBA("#372f2f"),
	GoldOre:             util.HexToRGBA("#ffe100"),
	IronOre:             util.HexToRGBA("#b8947d"),
	DiamondOre:          util.HexToRGBA("#40efd4"),
	DeepSlateStone:      util.HexToRGBA("#4c4c4c"),
	DeepSlateCoalOre:    util.HexToRGBA("#29344e"),
	DeepSlateGoldOre:    util.HexToRGBA("#ffe100"),
	DeepSlateIronOre:    util.HexToRGBA("#8a6548"),
	DeepSlateDiamondOre: util.HexToRGBA("#00ffe1"),
	Grass:               util.HexToRGBA("#00903f"),
}

var BlockMaxHealth = map[types.ItemID]float64{
	// Air:                  0.0,
	Grass:                5.0,
	Snow:                 5.0,
	Dirt:                 5.0,
	Sand:                 3.0,
	Stone:                10.0,
	CoalOre:              10.0,
	GoldOre:              10.0,
	IronOre:              10.0,
	DiamondOre:           10.0,
	CopperOre:            10.0,
	EmeraldOre:           10.0,
	LapisOre:             10.0,
	RedStoneOre:          10.0,
	DeepSlateStone:       15.0,
	DeepSlateCoalOre:     15.0,
	DeepSlateGoldOre:     15.0,
	DeepSlateIronOre:     15.0,
	DeepSlateDiamondOre:  15.0,
	DeepSlateCopperOre:   15.0,
	DeepSlateEmeraldOre:  15.0,
	DeepSlateLapisOre:    15.0,
	DeepSlateRedStoneOre: 15.0,
	Tree:                 10.0,
	TreeLeaves:           1.0,
	TreePlank:            10.0,
	Sapling:              1.0,
	// Torch:                0.0,
	// Coal:                 0.0,
	// RawGold:              0.0,
	// RawIron:              0.0,
	// Diamond:              0.0,
	// RawCopper:            0.0,
	// WoodShovel:           0.0,
	// StoneShovel:          0.0,
	// IronShovel:           0.0,
	// GoldenAxe:            0.0,
	// WoodAxe:              0.0,
	// StoneAxe:             0.0,
	// IronAxe:              0.0,
	// DiamondAxe:           0.0,
	// NetheriteAxe:         0.0,
	// GoldenPickaxe:        0.0,
	// WoodPickaxe:          0.0,
	// StonePickaxe:         0.0,
	// IronPickaxe:          0.0,
	// DiamondPickaxe:       0.0,
	// NetheritePickaxe:     0.0,
	// GoldenSword:          0.0,
	// WoodSword:            0.0,
	// StoneSword:           0.0,
	// IronSword:            0.0,
	// DiamondSword:         0.0,
	// NetheriteSword:       0.0,
}
