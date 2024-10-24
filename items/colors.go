package items

import (
	"image/color"
	"strconv"
)

var ItemColorMap = map[uint16]color.RGBA{
	Air:                 hexToRGBA("#0099ff"),
	GrassBlock:          hexToRGBA("#00903f"),
	Dirt:                hexToRGBA("#74573E"),
	Sand:                hexToRGBA("#fff5cc"),
	Stone:               hexToRGBA("#949494"),
	CoalOre:             hexToRGBA("#372f2f"),
	GoldOre:             hexToRGBA("#ffe100"),
	IronOre:             hexToRGBA("#b8947d"),
	DiamondOre:          hexToRGBA("#40efd4"),
	Deepslate:           hexToRGBA("#4c4c4c"),
	DeepslateCoalOre:    hexToRGBA("#29344e"),
	DeepslateGoldOre:    hexToRGBA("#ffe100"),
	DeepslateIronOre:    hexToRGBA("#8a6548"),
	DeepslateDiamondOre: hexToRGBA("#00ffe1"),
}

// hexToRGBA converts hex color to color.RGBA with "#FFFFFF" format
func hexToRGBA(hex string) color.RGBA {
	values, _ := strconv.ParseUint(string(hex[1:]), 16, 32)
	return color.RGBA{
		R: uint8(values >> 16),
		G: uint8((values >> 8) & 0xFF),
		B: uint8(values & 0xFF),
		A: 255,
	}
}
