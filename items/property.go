package items

type tag uint

// Item Bitmask tag
const (
	None  tag = 0
	Block tag = 1 << iota
	OreBlock
	Unbreakable
	Harvestable
	DropItem
	Item
	RawOre
	Tool
	Weapon
	Food
	MaterialWooden
	MaterialGold
	MaterialStone
	MaterialIron
	MaterialDiamond
	ToolHand
	ToolAxe
	ToolPickaxe
	ToolShovel
	ToolBucket
)
const All = tag(^uint(0))

type ItemProperty struct {
	DisplayName  string
	DropID       uint8
	MaxStackSize uint8
	Tags         tag
	BestToolTag  tag
}

var Property = map[uint8]ItemProperty{
	Air: {
		DisplayName:  "Air",
		DropID:       Air,
		MaxStackSize: 1,
		Tags:         None | Unbreakable,
	},
	Bedrock: {
		DisplayName:  "Bedrock",
		DropID:       0,
		MaxStackSize: 1,
		Tags:         Block | Unbreakable,
	},
	Bread: {
		DisplayName:  "Bread",
		DropID:       0,
		MaxStackSize: 64,
		Tags:         Item | Food,
	},
	Bucket: {
		DisplayName:  "Bucket",
		DropID:       0,
		MaxStackSize: 1,
		Tags:         Item | Tool,
	},
	Coal: {
		DisplayName:  "Coal",
		MaxStackSize: 64,
		Tags:         Item | DropItem | RawOre,
	},
	CoalOre: {
		DisplayName:  "Coal Ore",
		DropID:       Coal,
		MaxStackSize: 64,
		Tags:         Block | OreBlock,
		BestToolTag:  ToolPickaxe,
	},
	CraftingTable: {
		DisplayName:  "Crafting Table",
		DropID:       CraftingTable,
		MaxStackSize: 1,
		Tags:         Block,
		BestToolTag:  ToolAxe,
	},
	Diamond: {
		DisplayName:  "Diamond",
		MaxStackSize: 64,
		Tags:         Item | DropItem | RawOre,
	},
	DiamondAxe: {
		DisplayName:  "Diamond Axe",
		MaxStackSize: 1,
		Tags:         Item | Tool | ToolAxe | Weapon | MaterialDiamond,
	},
	DiamondOre: {
		DisplayName:  "Diamond Ore",
		DropID:       Diamond,
		MaxStackSize: 64,
		Tags:         Block | OreBlock,
		BestToolTag:  ToolPickaxe,
	},
	DiamondPickaxe: {
		DisplayName:  "Diamond Pickaxe",
		MaxStackSize: 1,
		Tags:         Item | Tool | ToolPickaxe | MaterialDiamond,
	},
	DiamondShovel: {
		DisplayName:  "Diamond Shovel",
		MaxStackSize: 1,
		Tags:         Item | Tool | ToolShovel | MaterialDiamond,
	},
	Dirt: {
		DisplayName:  "Dirt",
		DropID:       Dirt,
		MaxStackSize: 64,
		Tags:         Block,
		BestToolTag:  ToolShovel,
	},
	Furnace: {
		DisplayName:  "Furnace",
		DropID:       Furnace,
		MaxStackSize: 1,
		Tags:         Block,
		BestToolTag:  ToolPickaxe,
	},
	FurnaceOn: {
		DisplayName:  "Furnace On",
		DropID:       Furnace,
		MaxStackSize: 1,
		BestToolTag:  ToolPickaxe,
	},
	GoldIngot: {
		DisplayName:  "Gold Ingot",
		MaxStackSize: 64,
		Tags:         Item,
	},
	GoldOre: {
		DisplayName:  "Gold Ore",
		DropID:       RawGold,
		MaxStackSize: 64,
		Tags:         Block | OreBlock,
		BestToolTag:  ToolPickaxe,
	},
	GrassBlock: {
		DisplayName:  "Grass Block",
		DropID:       Dirt,
		MaxStackSize: 64,
		Tags:         Block,
		BestToolTag:  ToolShovel,
	},
	GrassBlockSnow: {
		DisplayName:  "Grass Block Snow",
		DropID:       Dirt,
		MaxStackSize: 64,
		Tags:         Block,
		BestToolTag:  ToolShovel,
	},
	IronAxe: {
		DisplayName:  "Iron Axe",
		MaxStackSize: 1,
		Tags:         Item | Tool | ToolAxe | MaterialIron,
	},
	IronIngot: {
		DisplayName:  "Iron Ingot",
		MaxStackSize: 64,
		Tags:         Item,
	},
	IronOre: {
		DisplayName:  "Iron Ore",
		DropID:       RawIron,
		MaxStackSize: 64,
		Tags:         Block | OreBlock,
		BestToolTag:  ToolPickaxe,
	},
	IronPickaxe: {
		DisplayName:  "Iron Pickaxe",
		MaxStackSize: 1,
		Tags:         Item | Tool | ToolPickaxe | MaterialIron,
	},
	IronShovel: {
		DisplayName:  "Iron Shovel",
		MaxStackSize: 1,
		Tags:         Item | Tool | ToolShovel | MaterialIron,
	},
	OakLeaves: {
		DisplayName:  "Oak Leaves",
		DropID:       OakSapling,
		MaxStackSize: 64,
		Tags:         Block,
		BestToolTag:  All,
	},
	OakLog: {
		DisplayName:  "Oak Log",
		DropID:       OakLog,
		MaxStackSize: 64,
		Tags:         Block,
		BestToolTag:  ToolAxe,
	},
	OakPlanks: {
		DisplayName:  "Oak Planks",
		DropID:       OakPlanks,
		MaxStackSize: 64,
		Tags:         Block,
		BestToolTag:  ToolAxe,
	},
	OakSapling: {
		DisplayName:  "Oak Sapling",
		DropID:       OakSapling,
		MaxStackSize: 64,
		Tags:         Block | Item,
		BestToolTag:  ToolAxe,
	},
	Obsidian: {
		DisplayName:  "Obsidian",
		DropID:       Obsidian,
		MaxStackSize: 64,
		Tags:         Block,
		BestToolTag:  ToolPickaxe,
	},
	RawGold: {
		DisplayName:  "Raw Gold",
		MaxStackSize: 64,
		Tags:         Item | DropItem | RawOre,
	},
	RawIron: {
		DisplayName:  "Raw Iron",
		MaxStackSize: 64,
		Tags:         Item | DropItem | RawOre,
	},
	Sand: {
		DisplayName:  "Sand",
		DropID:       Sand,
		MaxStackSize: 64,
		Tags:         Block,
		BestToolTag:  ToolShovel,
	},
	SmoothStone: {
		DisplayName:  "Smooth Stone",
		DropID:       SmoothStone,
		MaxStackSize: 64,
		Tags:         Block,
	},
	Snow: {
		DisplayName:  "Snow",
		DropID:       Dirt,
		MaxStackSize: 64,
		Tags:         Block,
		BestToolTag:  ToolShovel,
	},
	Snowball: {
		DisplayName:  "Snowball",
		DropID:       Snowball,
		MaxStackSize: 64,
		Tags:         Item,
	},
	Stick: {
		DisplayName:  "Stick",
		MaxStackSize: 64,
		Tags:         Item,
	},
	Stone: {
		DisplayName:  "Stone",
		DropID:       Stone,
		MaxStackSize: 64,
		Tags:         Block,
		BestToolTag:  ToolPickaxe,
	},
	StoneAxe: {
		DisplayName:  "Stone Axe",
		MaxStackSize: 1,
		Tags:         Item | Tool | ToolAxe | MaterialStone,
	},
	StoneBricks: {
		DisplayName:  "Stone Bricks",
		DropID:       StoneBricks,
		MaxStackSize: 64,
		Tags:         Block,
		BestToolTag:  ToolPickaxe,
	},
	StonePickaxe: {
		DisplayName:  "Stone Pickaxe",
		MaxStackSize: 1,
		Tags:         Item | Tool | ToolPickaxe | MaterialStone,
	},
	StoneShovel: {
		DisplayName:  "Stone Shovel",
		MaxStackSize: 1,
		Tags:         Item | Tool | ToolShovel | MaterialStone,
	},
	Tnt: {
		DisplayName:  "TNT",
		DropID:       Tnt,
		MaxStackSize: 64,
		Tags:         Block,
		BestToolTag:  All,
	},
	Torch: {
		DisplayName:  "Torch",
		DropID:       Torch,
		MaxStackSize: 64,
		Tags:         Item | Block,
		BestToolTag:  All,
	},
	WaterBucket: {
		DisplayName:  "Water Bucket",
		MaxStackSize: 1,
		Tags:         Item | Tool | ToolBucket | MaterialIron,
	},
	WoodenAxe: {
		DisplayName:  "Wooden Axe",
		MaxStackSize: 1,
		Tags:         Item | Tool | ToolAxe | MaterialWooden,
	},
	WoodenPickaxe: {
		DisplayName:  "Wooden Pickaxe",
		MaxStackSize: 1,
		Tags:         Item | Tool | ToolPickaxe | MaterialWooden,
	},
	WoodenShovel: {
		DisplayName:  "Wooden Shovel",
		MaxStackSize: 1,
		Tags:         Item | Tool | ToolShovel | MaterialWooden,
	},
}
