package items

type Recipe [][]Slot

var Recipes map[uint16]Recipe

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
