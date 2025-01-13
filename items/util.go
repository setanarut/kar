package items

import (
	"image/color"
	"math/rand/v2"
)

var BlockIDs []uint16

func init() {
	for id := range Property {
		if HasTag(id, Block) {
			BlockIDs = append(BlockIDs, id)
		}
	}
}

func HasTag(id uint16, tag tag) bool {
	return Property[id].Tags&tag != 0
}

func IsBestTool(blockID, toolID uint16) bool {
	return Property[blockID].BestToolTag&Property[toolID].Tags != 0
}
func IsStackable(id uint16) bool {
	return Property[id].MaxStackSize > 1
}
func GetDefaultDurability(id uint16) int {
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

func RandomBlock() uint16 {
	return BlockIDs[rand.IntN(len(BlockIDs))]
}
func DisplayName(id uint16) string {
	return Property[id].DisplayName
}
func RandomItem() uint16 {
	max := len(Property) - 1
	return uint16(1 + rand.IntN(max-1+1))
}

var ColorMap = map[uint16]color.RGBA{
	Air:           rgb(1, 1, 1),
	CraftingTable: rgb(194, 137, 62),
	GrassBlock:    rgb(133, 75, 54),
	Dirt:          rgb(133, 75, 54),
	Sand:          rgb(199, 193, 158),
	Stone:         rgb(139, 139, 139),
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
