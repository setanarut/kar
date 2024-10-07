package items

import (
	"image/color"
	"kar/engine/util"
	"kar/types"
)

var BlockColor = map[types.ItemID]color.RGBA{
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
var ItemName = map[types.ItemID]string{
	Air:                 "Air",
	Grass:               "Grass",
	Dirt:                "Dirt",
	Sand:                "Sand",
	Stone:               "Stone",
	CoalOre:             "CoalOre",
	GoldOre:             "GoldOre",
	IronOre:             "IronOre",
	DiamondOre:          "DiamondOre",
	DeepSlateStone:      "DeepSlateStone",
	DeepSlateCoalOre:    "DeepSlateCoalOre",
	DeepSlateGoldOre:    "DeepSlateGoldOre",
	DeepSlateIronOre:    "DeepSlateIronOre",
	DeepSlateDiamondOre: "DeepSlateDiamondOre",
	Coal:                "Coal",
	RawCopper:           "RawCopper",
	RawGold:             "RawGold",
	RawIron:             "RawIron",
	Diamond:             "Diamond",
	WoodShovel:          "WoodShovel",
	StoneShovel:         "StoneShovel",
	IronShovel:          "IronShovel",
	GoldenAxe:           "GoldenAxe",
	WoodAxe:             "WoodAxe",
	StoneAxe:            "StoneAxe",
	IronAxe:             "IronAxe",
	DiamondAxe:          "DiamondAxe",
	NetheriteAxe:        "NetheriteAxe",
	GoldenPickaxe:       "GoldenPickaxe",
	WoodPickaxe:         "WoodPickaxe",
	StonePickaxe:        "StonePickaxe",
	IronPickaxe:         "IronPickaxe",
	DiamondPickaxe:      "DiamondPickaxe",
	NetheritePickaxe:    "NetheritePickaxe",
	GoldenSword:         "GoldenSword",
	WoodSword:           "WoodSword",
	StoneSword:          "StoneSword",
	IronSword:           "IronSword",
	DiamondSword:        "DiamondSword",
	NetheriteSword:      "NetheriteSword",
}

const (
	Air types.ItemID = iota
	// Blocks
	Grass
	Dirt
	Sand
	Stone
	CoalOre
	GoldOre
	IronOre
	DiamondOre
	DeepSlateStone
	DeepSlateCoalOre
	DeepSlateGoldOre
	DeepSlateIronOre
	DeepSlateDiamondOre
	// Raw ores
	Coal
	RawCopper
	RawGold
	RawIron
	Diamond
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
	GoldenSword
	WoodSword
	StoneSword
	IronSword
	DiamondSword
	NetheriteSword
)
