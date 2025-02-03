package items

import (
	"image/color"
	"math/rand/v2"
)

var BlockIDs []uint8

func init() {
	for id := range Property {
		if HasTag(id, Block) {
			BlockIDs = append(BlockIDs, id)
		}
	}
}

func HasTag(id uint8, tag tag) bool {
	return Property[id].Tags&tag != 0
}

func IsBestTool(blockID, toolID uint8) bool {
	return Property[blockID].BestToolTag&Property[toolID].Tags != 0
}
func IsStackable(id uint8) bool {
	return Property[id].MaxStackSize > 1
}
func GetDefaultDurability(id uint8) int {
	if HasTag(id, Tool) {
		if HasTag(id, MaterialWooden) {
			return 25
		}
		if HasTag(id, MaterialStone) {
			return 50
		}
		if HasTag(id, MaterialGold) {
			return 100
		}
		if HasTag(id, MaterialIron) {
			return 200
		}
		if HasTag(id, MaterialDiamond) {
			return 400
		}
	}
	return 0
}

func RandomBlock() uint8 {
	return BlockIDs[rand.IntN(len(BlockIDs))]
}
func DisplayName(id uint8) string {
	return Property[id].DisplayName
}
func RandomItem() uint8 {
	max := len(Property) - 1
	return uint8(1 + rand.IntN(max-1+1))
}

var ColorMap = map[uint8]color.RGBA{
	Air:           rgb(0, 62, 161),
	CraftingTable: rgb(194, 137, 62),
	GrassBlock:    rgb(133, 75, 54),
	Dirt:          rgb(133, 75, 54),
	Sand:          rgb(199, 193, 158),
	Stone:         rgb(120, 120, 120),
	CoalOre:       rgb(0, 0, 0),
	GoldOre:       rgb(255, 221, 0),
	IronOre:       rgb(151, 176, 205),
	DiamondOre:    rgb(0, 247, 255),
	OakLog:        rgb(227, 131, 104),
	OakPlanks:     rgb(224, 153, 145),
	OakLeaves:     rgb(0, 160, 16),
}

func rgb(r, g, b uint8) color.RGBA {
	return color.RGBA{r, g, b, 255}
}
