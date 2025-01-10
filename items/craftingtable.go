package items

import (
	"fmt"
)

type Recipe [][]SlotData

var Recipes map[uint16]Recipe

type CraftTable struct {
	SlotPosX, SlotPosY int
	Slots              [][]SlotData
	ResultSlot         SlotData
}

func NewCraftTable() *CraftTable {
	return &CraftTable{
		Slots: [][]SlotData{
			{SlotData{}, SlotData{}, SlotData{}},
			{SlotData{}, SlotData{}, SlotData{}},
			{SlotData{}, SlotData{}, SlotData{}}},
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
	// if ct.ResultSlot.ID != 0 {
	// 	for y := 0; y < 3; y++ {
	// 		for x := 0; x < 3; x++ {
	// 			ct.Slots[y][x].Quantity--
	// 			if ct.Slots[y][x].Quantity == 0 {
	// 				ct.ClearCurrenSlot()
	// 			}
	// 		}
	// 	}
	// 	ct.ResultSlot.Quantity++
	// }
}
func (ct *CraftTable) ClearTable() {
	ct.Slots = [][]SlotData{
		{SlotData{}, SlotData{}, SlotData{}},
		{SlotData{}, SlotData{}, SlotData{}},
		{SlotData{}, SlotData{}, SlotData{}}}
}

func (ct *CraftTable) CurrentSlot() SlotData {
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
	ct.Slots[ct.SlotPosY][ct.SlotPosX] = SlotData{}
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
	normalizedGrid := make([][]SlotData, maxRow-minRow+1)
	for i := range normalizedGrid {
		normalizedGrid[i] = make([]SlotData, maxCol-minCol+1)
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

	Recipes[OakPlanks] = [][]SlotData{{SlotData{ID: OakLog}}}
	Recipes[Stick] = [][]SlotData{
		{SlotData{ID: OakPlanks}, SlotData{}, SlotData{}},
		{SlotData{ID: OakPlanks}, SlotData{}, SlotData{}},
		{SlotData{}, SlotData{}, SlotData{}}}
	Recipes[IronPickaxe] = [][]SlotData{
		{SlotData{ID: IronIngot}, SlotData{ID: IronIngot}, SlotData{ID: IronIngot}},
		{SlotData{ID: 0}, SlotData{ID: Stick}, SlotData{ID: 0}},
		{SlotData{ID: 0}, SlotData{ID: Stick}, SlotData{ID: 0}}}
}
