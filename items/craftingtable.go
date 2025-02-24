package items

type CraftTable struct {
	SlotPosX, SlotPosY int
	Slots              [][]Slot
	ResultSlot         Slot
}

func NewCraftTable() *CraftTable {
	return &CraftTable{
		Slots: [][]Slot{{{}, {}, {}}, {{}, {}, {}}, {{}, {}, {}}},
	}
}

// tarif yoksa sıfır döndürür (air) tarif sonucu ve miktarını döndürür.
func (ct *CraftTable) CheckRecipe() (uint8, uint8) {
	cropped := ct.cropRecipe(ct.Slots)
	for itemIDKey, recipe := range Recipes { // TODO CraftTable Recipes'i argüman olarak alsın. fırın için farklı tarifler mümkün olsun.
		if ct.Equal(cropped, ct.cropRecipe(recipe)) {

			return itemIDKey, recipe[0][0].Quantity
		}
	}
	return 0, 0
}

// returns minimum result item quantity
func (ct *CraftTable) UpdateResultSlot() uint8 {
	id, quantity := ct.CheckRecipe()
	ct.ResultSlot.ID = id
	minimum := uint8(255)
	for y := range 3 {
		for x := range 3 {
			if ct.Slots[y][x].Quantity != 0 {
				minimum = min(minimum, ct.Slots[y][x].Quantity)
			}
		}
	}
	if minimum == 255 {
		minimum = 0
	}

	if ct.ResultSlot.ID != 0 {
		ct.ResultSlot.Quantity = minimum * quantity
	}
	return minimum * quantity
}
func (ct *CraftTable) ClearTable() {
	ct.Slots = [][]Slot{{{}, {}, {}}, {{}, {}, {}}, {{}, {}, {}}}
}

func (ct *CraftTable) CurrentSlot() *Slot {
	return &ct.Slots[ct.SlotPosY][ct.SlotPosX]
}

func (ct *CraftTable) ClearCurrenSlot() {
	ct.Slots[ct.SlotPosY][ct.SlotPosX] = Slot{}
}

func (ct *CraftTable) RemoveItem(x, y int) {
	if ct.Slots[y][x].Quantity == 1 {
		ct.Slots[y][x].ID = 0
		ct.Slots[y][x].Quantity = 0
	} else if ct.Slots[y][x].Quantity > 0 {
		ct.Slots[y][x].Quantity--
	}
}

func (ct *CraftTable) Equal(recipeA, recipeB Recipe) bool {
	// Öncelikle boyutlarını karşılaştır
	if len(recipeA) != len(recipeB) {
		return false
	}
	for i := range recipeA {
		if len(recipeA[i]) != len(recipeB[i]) {
			return false
		}
		// Her hücreyi tek tek karşılaştır
		for j := range recipeA[i] {
			if recipeA[i][j].ID != recipeB[i][j].ID {
				return false
			}
		}
	}
	return true
}

// cropRecipe normalizes grid
func (ct *CraftTable) cropRecipe(reci Recipe) Recipe {
	minRow, maxRow := len(reci), 0
	minCol, maxCol := len(reci[0]), 0
	for i := range len(reci) {
		for j := range len(reci[i]) {
			if reci[i][j].ID != 0 {
				if i < minRow {
					minRow = i
				}
				if i > maxRow {
					maxRow = i
				}
				if j < minCol {
					minCol = j
				}
				if j > maxCol {
					maxCol = j
				}
			}
		}
	}
	if minRow > maxRow || minCol > maxCol {
		return reci
	}
	normalizedGrid := make([][]Slot, maxRow-minRow+1)
	for i := range normalizedGrid {
		normalizedGrid[i] = make([]Slot, maxCol-minCol+1)
		for j := range normalizedGrid[i] {
			normalizedGrid[i][j] = reci[minRow+i][minCol+j]
		}
	}
	return Recipe(normalizedGrid)
}
