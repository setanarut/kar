package items

type Recipe [][]Slot

var Recipes map[uint8]Recipe

func init() {
	Recipes = make(map[uint8]Recipe)

	Recipes[OakPlanks] = [][]Slot{{Slot{ID: OakLog}}}
	Recipes[Stick] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{}},
		{Slot{ID: OakPlanks}, Slot{}},
	}
	Recipes[CraftingTable] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}},
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}},
	}

	// Axe
	Recipes[WoodenAxe] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}, Slot{ID: 0}},
		{Slot{ID: OakPlanks}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	Recipes[StoneAxe] = [][]Slot{
		{Slot{ID: Stone}, Slot{ID: Stone}, Slot{ID: 0}},
		{Slot{ID: Stone}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	Recipes[IronAxe] = [][]Slot{
		{Slot{ID: IronIngot}, Slot{ID: IronIngot}, Slot{ID: 0}},
		{Slot{ID: IronIngot}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}
	Recipes[DiamondAxe] = [][]Slot{
		{Slot{ID: Diamond}, Slot{ID: Diamond}, Slot{ID: 0}},
		{Slot{ID: Diamond}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	// Shovel
	Recipes[WoodenShovel] = [][]Slot{
		{Slot{ID: 0}, Slot{ID: OakPlanks}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}
	Recipes[StoneShovel] = [][]Slot{
		{Slot{ID: 0}, Slot{ID: Stone}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}
	Recipes[IronShovel] = [][]Slot{
		{Slot{ID: 0}, Slot{ID: IronIngot}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}
	Recipes[DiamondShovel] = [][]Slot{
		{Slot{ID: 0}, Slot{ID: Diamond}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	// Pickaxe
	Recipes[WoodenPickaxe] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}, Slot{ID: OakPlanks}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}
	Recipes[StonePickaxe] = [][]Slot{
		{Slot{ID: Stone}, Slot{ID: Stone}, Slot{ID: Stone}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	Recipes[IronPickaxe] = [][]Slot{
		{Slot{ID: IronIngot}, Slot{ID: IronIngot}, Slot{ID: IronIngot}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}
	Recipes[DiamondPickaxe] = [][]Slot{
		{Slot{ID: Diamond}, Slot{ID: Diamond}, Slot{ID: Diamond}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	Recipes[Furnace] = [][]Slot{
		{Slot{ID: Stone}, Slot{ID: Stone}, Slot{ID: Stone}},
		{Slot{ID: Stone}, Slot{ID: Air}, Slot{ID: Stone}},
		{Slot{ID: Stone}, Slot{ID: Stone}, Slot{ID: Stone}},
	}

	// output item multiplier
	for _, recipe := range Recipes {
		recipe[0][0].Quantity = 1
	}
	Recipes[OakPlanks][0][0].Quantity = 4
	Recipes[Stick][0][0].Quantity = 4
}
