package items

import (
	"fmt"
)

type CraftTable struct {
	SlotPosX, SlotPosY int
	Slots              [][]Slot
	ResultSlot         Slot
}

func NewCraftTable() *CraftTable {
	return &CraftTable{
		Slots: [][]Slot{
			{Slot{}, Slot{}, Slot{}},
			{Slot{}, Slot{}, Slot{}},
			{Slot{}, Slot{}, Slot{}}},
	}
}

// tarif yoksa sıfır döndürür (air)
func (ct *CraftTable) CheckRecipe() uint16 {
	cropped := ct.cropRecipe(ct.Slots)
	for itemIDKey, recipe := range Recipes {
		if ct.Equal(cropped, ct.cropRecipe(recipe)) {
			return itemIDKey
		}
	}
	return 0
}

func (ct *CraftTable) UpdateResultSlot() uint8 {
	ct.ResultSlot.ID = ct.CheckRecipe()
	minimum := uint8(255)
	for y := range 3 {
		for x := range 3 {
			if ct.Get(x, y).Quantity != 0 {
				minimum = min(minimum, ct.Get(x, y).Quantity)
			}
		}
	}

	if minimum == 255 {
		minimum = 0
	}

	if ct.ResultSlot.ID != 0 {
		ct.ResultSlot.Quantity = minimum
	}

	// ct.PrintGrid()
	// fmt.Println("min:", minimum)
	// fmt.Println("result:", ct.ResultSlot)
	return minimum
}
func (ct *CraftTable) ClearTable() {
	ct.Slots = [][]Slot{
		{Slot{}, Slot{}, Slot{}},
		{Slot{}, Slot{}, Slot{}},
		{Slot{}, Slot{}, Slot{}}}
}

func (ct *CraftTable) CurrentSlot() Slot {
	return ct.Slots[ct.SlotPosY][ct.SlotPosX]
}
func (ct *CraftTable) SetCurrentSlotQuantity(q uint8) {
	ct.Slots[ct.SlotPosY][ct.SlotPosX].Quantity = q
}
func (ct *CraftTable) AddCurrentSlotQuantity(q uint8) {
	ct.Slots[ct.SlotPosY][ct.SlotPosX].Quantity += q
}
func (ct *CraftTable) SubCurrentSlotQuantity(q uint8) {
	ct.Slots[ct.SlotPosY][ct.SlotPosX].Quantity -= q
}
func (ct *CraftTable) ClearCurrenSlot() {
	ct.Slots[ct.SlotPosY][ct.SlotPosX] = Slot{}
}
func (ct *CraftTable) ClearSlot(x, y int) {
	ct.Slots[y][x] = Slot{}
}

func (ct *CraftTable) SetCurrentSlot(id uint16) {
	ct.Slots[ct.SlotPosY][ct.SlotPosX].ID = id
}

func (ct *CraftTable) CurrentSLot() *Slot {
	return &ct.Slots[ct.SlotPosY][ct.SlotPosX]
}

func (ct *CraftTable) Set(x, y int, id uint16) {
	ct.Slots[y][x].ID = id
}
func (ct *CraftTable) SetQuantity(x, y int, q uint8) {
	ct.Slots[y][x].Quantity = q
}

func (ct *CraftTable) Get(x, y int) *Slot {
	return &ct.Slots[y][x]
}
func (ct *CraftTable) RemoveItem(x, y int) {
	if ct.Slots[y][x].Quantity == 1 {
		ct.Slots[y][x].ID = 0
		ct.Slots[y][x].Quantity = 0
	} else {
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

func (ct *CraftTable) PrintGrid() {
	for _, row := range ct.Slots {
		fmt.Println(row)
	}
}
