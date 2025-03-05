package items

import "image"

type CraftTable struct {
	Pos        image.Point // Slot position
	Slots      [][]Slot
	ResultSlot Slot
}

func NewCraftTable() *CraftTable {
	return &CraftTable{
		Slots: [][]Slot{{{}, {}, {}}, {{}, {}, {}}, {{}, {}, {}}},
	}
}

// returns zero (air) if no recipe is found, otherwise returns the recipe result and quantity.
func (c *CraftTable) CheckRecipe(recipes map[uint8]Recipe) (uint8, uint8) {
	cropped := c.cropRecipe(c.Slots)
	for itemIDKey, recipe := range recipes {
		if c.Equal(cropped, c.cropRecipe(recipe)) {

			return itemIDKey, recipe[0][0].Quantity
		}
	}
	return 0, 0
}

// returns minimum result item quantity
func (c *CraftTable) UpdateResultSlot(recipes map[uint8]Recipe) uint8 {
	id, quantity := c.CheckRecipe(recipes)
	c.ResultSlot.ID = id
	minimum := uint8(255)
	for y := range 3 {
		for x := range 3 {
			if c.Slots[y][x].Quantity != 0 {
				minimum = min(minimum, c.Slots[y][x].Quantity)
			}
		}
	}
	if minimum == 255 {
		minimum = 0
	}

	if c.ResultSlot.ID != 0 {
		c.ResultSlot.Quantity = minimum * quantity
	}
	return minimum * quantity
}
func (c *CraftTable) ClearTable() {
	c.Slots = [][]Slot{{{}, {}, {}}, {{}, {}, {}}, {{}, {}, {}}}
}

func (c *CraftTable) CurrentSlot() *Slot {
	return &c.Slots[c.Pos.Y][c.Pos.X]
}

func (c *CraftTable) ClearCurrenSlot() {
	c.Slots[c.Pos.Y][c.Pos.X] = Slot{}
}

func (c *CraftTable) RemoveItem(x, y int) {
	if c.Slots[y][x].Quantity == 1 {
		c.Slots[y][x].ID = 0
		c.Slots[y][x].Quantity = 0
	} else if c.Slots[y][x].Quantity > 0 {
		c.Slots[y][x].Quantity--
	}
}

func (c *CraftTable) Equal(recipeA, recipeB Recipe) bool {
	// First, compare their sizes
	if len(recipeA) != len(recipeB) {
		return false
	}
	for i := range recipeA {
		if len(recipeA[i]) != len(recipeB[i]) {
			return false
		}
		// Compare each cell one by one
		for j := range recipeA[i] {
			if recipeA[i][j].ID != recipeB[i][j].ID {
				return false
			}
		}
	}
	return true
}

// cropRecipe normalizes grid
func (c *CraftTable) cropRecipe(reci Recipe) Recipe {
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
