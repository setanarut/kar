package items

func IsBreakable(id uint16) bool {
	return Property[id].Category&Unbreakable == 0
}

func IsHarvestable(id uint16) bool {
	return Property[id].Category&Harvestable != 0
}
func IsBlock(id uint16) bool {
	return Property[id].Category&Block != 0
}
