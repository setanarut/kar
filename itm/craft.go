package itm

import (
	"fmt"
)

type Recipe [][]uint16

type CraftingManager struct {
}

var recipes map[uint16][]Recipe

var emptyGrid = [][]uint16{
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
}

func NewCraftingManager() *CraftingManager {
	return &CraftingManager{}
}

func (ct *CraftingManager) CheckRecipe(pattern [][]uint16) (item uint16, ok bool) {
	cropped := ct.CropGrid(pattern)
	for itemIDKey, subRecipes := range recipes {
		for _, subRecipe := range subRecipes {
			if ct.Equal(cropped, subRecipe) {
				return itemIDKey, true
			}
		}
	}
	return 0, false
}

func (ct *CraftingManager) Equal(grid1, grid2 [][]uint16) bool {
	// Öncelikle boyutlarını karşılaştır
	if len(grid1) != len(grid2) {
		return false
	}
	for i := range grid1 {
		if len(grid1[i]) != len(grid2[i]) {
			return false
		}
		// Her hücreyi tek tek karşılaştır
		for j := range grid1[i] {
			if grid1[i][j] != grid2[i][j] {
				return false
			}
		}
	}
	return true
}

// CropGrid normalizes grid
func (ct *CraftingManager) CropGrid(grid [][]uint16) [][]uint16 {
	minRow, maxRow := len(grid), 0
	minCol, maxCol := len(grid[0]), 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] != 0 {
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
		return emptyGrid
	}
	normalizedGrid := make([][]uint16, maxRow-minRow+1)
	for i := range normalizedGrid {
		normalizedGrid[i] = make([]uint16, maxCol-minCol+1)
		for j := range normalizedGrid[i] {
			normalizedGrid[i][j] = grid[minRow+i][minCol+j]
		}
	}
	return normalizedGrid
}

func (ct *CraftingManager) PrintGrid(grid [][]uint16) {
	for _, row := range grid {
		fmt.Println(row)
	}
	fmt.Println()
}

func init() {
	recipes = make(map[uint16][]Recipe)
	recipes[Torch] = []Recipe{
		[][]uint16{
			{CharCoal},
			{Stick},
		},
		[][]uint16{
			{Coal},
			{Stick},
		},
	}
}
