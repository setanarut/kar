package items

type Recipe [][]Slot

var CraftingRecipes map[uint8]Recipe
var FurnaceRecipes map[uint8]Recipe

func init() {

	CraftingRecipes = make(map[uint8]Recipe)
	FurnaceRecipes = make(map[uint8]Recipe)

	CraftingRecipes[OakPlanks] = [][]Slot{{Slot{ID: OakLog}}}

	CraftingRecipes[Stick] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{}},
		{Slot{ID: OakPlanks}, Slot{}},
	}

	CraftingRecipes[CraftingTable] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}},
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}},
	}

	// Axe
	CraftingRecipes[WoodenAxe] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}, Slot{ID: 0}},
		{Slot{ID: OakPlanks}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	CraftingRecipes[StoneAxe] = [][]Slot{
		{Slot{ID: Stone}, Slot{ID: Stone}, Slot{ID: 0}},
		{Slot{ID: Stone}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	CraftingRecipes[IronAxe] = [][]Slot{
		{Slot{ID: IronIngot}, Slot{ID: IronIngot}, Slot{ID: 0}},
		{Slot{ID: IronIngot}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	CraftingRecipes[DiamondAxe] = [][]Slot{
		{Slot{ID: Diamond}, Slot{ID: Diamond}, Slot{ID: 0}},
		{Slot{ID: Diamond}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	// Shovel
	CraftingRecipes[WoodenShovel] = [][]Slot{
		{Slot{ID: 0}, Slot{ID: OakPlanks}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	CraftingRecipes[StoneShovel] = [][]Slot{
		{Slot{ID: 0}, Slot{ID: Stone}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	CraftingRecipes[IronShovel] = [][]Slot{
		{Slot{ID: 0}, Slot{ID: IronIngot}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	CraftingRecipes[DiamondShovel] = [][]Slot{
		{Slot{ID: 0}, Slot{ID: Diamond}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	// Pickaxe
	CraftingRecipes[WoodenPickaxe] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}, Slot{ID: OakPlanks}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	CraftingRecipes[StonePickaxe] = [][]Slot{
		{Slot{ID: Stone}, Slot{ID: Stone}, Slot{ID: Stone}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	CraftingRecipes[IronPickaxe] = [][]Slot{
		{Slot{ID: IronIngot}, Slot{ID: IronIngot}, Slot{ID: IronIngot}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	CraftingRecipes[DiamondPickaxe] = [][]Slot{
		{Slot{ID: Diamond}, Slot{ID: Diamond}, Slot{ID: Diamond}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	CraftingRecipes[Furnace] = [][]Slot{
		{Slot{ID: Stone}, Slot{ID: Stone}, Slot{ID: Stone}},
		{Slot{ID: Stone}, Slot{ID: Air}, Slot{ID: Stone}},
		{Slot{ID: Stone}, Slot{ID: Stone}, Slot{ID: Stone}},
	}

	// output item multiplier
	for _, recipe := range CraftingRecipes {
		recipe[0][0].Quantity = 1
	}
	CraftingRecipes[OakPlanks][0][0].Quantity = 4
	CraftingRecipes[Stick][0][0].Quantity = 4

	// --- FURNACE RECIPES ---

	FurnaceRecipes[IronIngot] = [][]Slot{
		{Slot{ID: RawIron}, Slot{}},
		{Slot{ID: Coal}, Slot{}},
	}
	FurnaceRecipes[GoldIngot] = [][]Slot{
		{Slot{ID: RawGold}, Slot{}},
		{Slot{ID: Coal}, Slot{}},
	}

	// output item multiplier
	for _, recipe := range FurnaceRecipes {
		recipe[0][0].Quantity = 1
	}

}
