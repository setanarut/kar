package items

import (
	"fmt"
)

type Recipe [][]uint16

var Recipes map[uint16][]Recipe

type CraftTable struct {
	SlotPosX, SlotPosY int
	Slots              [][]uint16
	ResultSlot         uint16
}

func NewCraftTable() *CraftTable {
	return &CraftTable{
		Slots: [][]uint16{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0}},
	}
}

// tarif yoksa sıfır döndürür (air)
func (ct *CraftTable) CheckRecipe() uint16 {
	cropped := ct.cropRecipe(ct.Slots)
	for itemIDKey, subRecipes := range Recipes {
		for _, subRecipe := range subRecipes {
			if ct.Equal(cropped, subRecipe) {
				return itemIDKey
			}
		}
	}
	return 0
}

func (ct *CraftTable) UpdateResultSlot() {
	ct.ResultSlot = ct.CheckRecipe()
}
func (ct *CraftTable) ClearTable() {
	ct.Slots = [][]uint16{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0}}
}

func (ct *CraftTable) CurrentSlot() uint16 {
	return ct.Slots[ct.SlotPosY][ct.SlotPosX]
}
func (ct *CraftTable) SetCurrentSlot(id uint16) {
	ct.Slots[ct.SlotPosY][ct.SlotPosX] = id
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
			if recipeA[i][j] != recipeB[i][j] {
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
			if reci[i][j] != 0 {
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
	normalizedGrid := make([][]uint16, maxRow-minRow+1)
	for i := range normalizedGrid {
		normalizedGrid[i] = make([]uint16, maxCol-minCol+1)
		for j := range normalizedGrid[i] {
			normalizedGrid[i][j] = reci[minRow+i][minCol+j]
		}
	}
	return Recipe(normalizedGrid)
}

func (ct *CraftTable) PrintGrid(r Recipe) {
	for _, row := range r {
		fmt.Println(row)
	}
	fmt.Println()
}

func init() {
	Recipes = make(map[uint16][]Recipe)
	Recipes[Torch] = []Recipe{
		[][]uint16{
			{Coal},
			{Stick},
		},
	}
}
