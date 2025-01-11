package items

import (
	"fmt"
)

type Recipe [][]Slot

var Recipes map[uint16]Recipe

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

func (ct *CraftTable) UpdateResultSlot() {
	ct.ResultSlot.ID = ct.CheckRecipe()
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

func (ct *CraftTable) SetCurrentSlot(id uint16) {
	ct.Slots[ct.SlotPosY][ct.SlotPosX].ID = id
	ct.UpdateResultSlot()
}
func (ct *CraftTable) SetSlot(x, y, id uint16) {
	ct.Slots[y][x].ID = id
	ct.UpdateResultSlot()
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
	for i := 0; i < len(reci); i++ {
		for j := 0; j < len(reci[i]); j++ {
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
	fmt.Println()
}

func init() {
	Recipes = make(map[uint16]Recipe)

	Recipes[OakPlanks] = [][]Slot{{Slot{ID: OakLog}}}
	Recipes[Stick] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{}, Slot{}},
		{Slot{ID: OakPlanks}, Slot{}, Slot{}},
		{Slot{}, Slot{}, Slot{}}}
	Recipes[IronPickaxe] = [][]Slot{
		{Slot{ID: IronIngot}, Slot{ID: IronIngot}, Slot{ID: IronIngot}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}}}
	Recipes[CraftingTable] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}},
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}},
	}
}
