package items

import (
	"image/color"
	"kar/engine/util"
	"kar/types"
)

var Colors = map[types.ItemType]color.RGBA{
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

const (
	Air types.ItemType = iota
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
