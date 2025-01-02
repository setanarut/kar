package items

import (
	"kar/engine/mathutil"
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
	return Property[id].Stackable > 1
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
func RandomItem() uint16 {
	return uint16(mathutil.RandRangeInt(1, len(Property)-1))
}
