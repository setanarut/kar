package items

type Recipe [][]Slot

var Recipes map[uint16]Recipe

func init() {
	Recipes = make(map[uint16]Recipe)

	Recipes[OakPlanks] = [][]Slot{{Slot{ID: OakLog}}}
	Recipes[Stick] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{}},
		{Slot{ID: OakPlanks}, Slot{}},
	}
	Recipes[CraftingTable] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}},
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}},
	}
	Recipes[WoodenAxe] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}, Slot{ID: 0}},
		{Slot{ID: OakPlanks}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}
	Recipes[WoodenShovel] = [][]Slot{
		{Slot{ID: 0}, Slot{ID: OakPlanks}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}
	Recipes[WoodenPickaxe] = [][]Slot{
		{Slot{ID: OakPlanks}, Slot{ID: OakPlanks}, Slot{ID: OakPlanks}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}
	Recipes[IronPickaxe] = [][]Slot{
		{Slot{ID: IronIngot}, Slot{ID: IronIngot}, Slot{ID: IronIngot}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
		{Slot{ID: 0}, Slot{ID: Stick}, Slot{ID: 0}},
	}

	Recipes[OakPlanks][0][0].Quantity = 4
	Recipes[Stick][0][0].Quantity = 4
	Recipes[CraftingTable][0][0].Quantity = 1
	Recipes[WoodenAxe][0][0].Quantity = 1
	Recipes[WoodenShovel][0][0].Quantity = 1
	Recipes[WoodenPickaxe][0][0].Quantity = 1
	Recipes[IronPickaxe][0][0].Quantity = 1
}
