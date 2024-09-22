package items

import (
	"image/color"
	"kar/engine/util"
	"kar/types"
)

var Colors map[types.ItemType]color.RGBA

const (
	Air types.ItemType = iota
	// Blocks
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

func init() {
	Colors = map[types.ItemType]color.RGBA{
		Dirt:                util.HexToRGBA("#ff7e57"),
		Sand:                util.HexToRGBA("#fff5cc"),
		Stone:               util.HexToRGBA("#c5c5c5"),
		CoalOre:             util.HexToRGBA("#000000"),
		GoldOre:             util.HexToRGBA("#ffe100"),
		IronOre:             util.HexToRGBA("#7d7d7d"),
		DiamondOre:          util.HexToRGBA("#00ffe1"),
		DeepSlateStone:      util.HexToRGBA("#c5c5c5"),
		DeepSlateCoalOre:    util.HexToRGBA("#000000"),
		DeepSlateGoldOre:    util.HexToRGBA("#ffe100"),
		DeepSlateIronOre:    util.HexToRGBA("#7d7d7d"),
		DeepSlateDiamondOre: util.HexToRGBA("#00ffe1"),
	}
}
