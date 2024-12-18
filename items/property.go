package items

type Category uint

// Item Bitmask Category
const (
	None  Category = 0
	Block Category = 1 << iota
	OreBlock
	Unbreakable
	Harvestable
	DropItem
	Item
	RawOre
	Tool
	Weapon
	Food
	All = ^uint(0)
)

type ItemProperty struct {
	DisplayName string
	DropID      uint16
	Stackable   uint8
	MaxHealth   float64
	Category    Category
}

var Property = map[uint16]ItemProperty{
	Air: {
		DisplayName: "",
		DropID:      Air,
		Stackable:   0,
		MaxHealth:   1000,
		Category:    None | Unbreakable,
	},
	Arrow: {
		DisplayName: "Arrow",
		DropID:      Arrow,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Weapon | Item | DropItem,
	},

	Bedrock: {
		DisplayName: "Bedrock",
		DropID:      0,
		Stackable:   0,
		MaxHealth:   1000,
		Category:    Block | Unbreakable,
	},
	BeetrootSeeds: {
		DisplayName: "Beetroot Seeds",
		DropID:      0,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Item | Unbreakable | DropItem,
	},
	Bow: {
		DisplayName: "Bow",
		DropID:      0,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Item | Weapon,
	},

	Bread: {
		DisplayName: "Bread",
		DropID:      0,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Item | Food,
	},
	BrewingStand: {
		DisplayName: "Brewing Stand",
		DropID:      0,
		Stackable:   1,
		MaxHealth:   10,
		Category:    Block,
	},
	Bucket: {
		DisplayName: "Bucket",
		DropID:      0,
		Stackable:   1,
		MaxHealth:   10,
		Category:    Item | Tool,
	},
	CartographyTable: {
		DisplayName: "Cartography Table",
		DropID:      0,
		Stackable:   1,
		MaxHealth:   10,
		Category:    Block,
	},
	Charcoal: {
		DisplayName: "Char Coal",
		Stackable:   64,
		MaxHealth:   1,
		Category:    Item | DropItem | RawOre,
	},
	Clay: {
		DisplayName: "Clay",
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},
	Coal: {
		DisplayName: "Coal",
		Stackable:   64,
		MaxHealth:   1,
		Category:    Item | DropItem | RawOre,
	},
	CoalBlock: {
		DisplayName: "Coal Block",
		Stackable:   64,
		MaxHealth:   1,
		Category:    Block,
	},
	CoalOre: {
		DisplayName: "Coal Ore",
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block | OreBlock,
	},
	Cobblestone: {
		DisplayName: "Cobblestone",
		DropID:      Cobblestone,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},
	CopperIngot: {
		DisplayName: "Copper Ingot",
		Stackable:   64,
		MaxHealth:   1,
		Category:    Item,
	},
	CopperOre: {
		DisplayName: "Copper Ore",
		DropID:      RawCopper,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block | OreBlock,
	},
	CraftingTable: {
		DisplayName: "Crafting Table",
		DropID:      Air,
		Stackable:   1,
		MaxHealth:   10,
		Category:    Block,
	},

	CrossbowStandby: {
		DisplayName: "CrossbowStandby",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Weapon,
	},

	Diamond: {
		DisplayName: "Diamond",
		Stackable:   64,
		MaxHealth:   1,
		Category:    Item | DropItem | RawOre,
	},
	DiamondAxe: {
		DisplayName: "Diamond Axe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},

	DiamondHoe: {
		DisplayName: "Diamond Hoe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	DiamondOre: {
		DisplayName: "Diamond Ore",
		DropID:      Diamond,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block | OreBlock,
	},
	DiamondPickaxe: {
		DisplayName: "Diamond Pickaxe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	DiamondShovel: {
		DisplayName: "Diamond Hoe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	DiamondSword: {
		DisplayName: "Diamond Sword",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Weapon,
	},
	Dirt: {
		DisplayName: "Dirt",
		DropID:      Dirt,
		Stackable:   64,
		MaxHealth:   5,
		Category:    Block,
	},

	EnchantingTable: {
		DisplayName: "Enchanting Table",
		DropID:      EnchantingTable,
		Stackable:   1,
		MaxHealth:   10,
		Category:    Block,
	},
	EndPortalFrame: {
		DisplayName: "EndPortalFrame",
		DropID:      EndPortalFrame,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},
	FletchingTable: {
		DisplayName: "Fletching Table",
		DropID:      FletchingTable,
		Stackable:   1,
		MaxHealth:   10,
		Category:    Block,
	},
	Furnace: {
		DisplayName: "Furnace",
		DropID:      Furnace,
		Stackable:   1,
		MaxHealth:   10,
		Category:    Block,
	},
	FurnaceOn: {
		DisplayName: "Furnace On",
		DropID:      Furnace,
		Stackable:   1,
		MaxHealth:   10,
		Category:    Block,
	},
	GoldIngot: {
		DisplayName: "Gold Ingot",
		Stackable:   64,
		MaxHealth:   1,
		Category:    Item,
	},
	GoldOre: {
		DisplayName: "Gold Ore",
		DropID:      RawGold,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block | OreBlock,
	},
	GoldenAxe: {
		DisplayName: "Golden Axe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	GoldenHoe: {
		DisplayName: "Golden Hoe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	GoldenPickaxe: {
		DisplayName: "Golden Pickaxe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	GoldenShovel: {
		DisplayName: "Golden Hoe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	GoldenSword: {
		DisplayName: "Golden Sword",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Weapon,
	},
	GrassBlock: {
		DisplayName: "Grass Block",
		DropID:      Dirt,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},
	GrassBlockSnow: {
		DisplayName: "Grass Block Snow",
		DropID:      Dirt,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},

	IronAxe: {
		DisplayName: "Iron Axe",
		Stackable:   1,
		MaxHealth:   1,
	},
	IronHoe: {
		DisplayName: "Iron Hoe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	IronIngot: {
		DisplayName: "Iron Ingot",
		Stackable:   64,
		MaxHealth:   1,
		Category:    Item,
	},
	IronOre: {
		DisplayName: "Iron Ore",
		DropID:      RawIron,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block | OreBlock,
	},
	IronPickaxe: {
		DisplayName: "Iron Pickaxe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	IronShovel: {
		DisplayName: "Iron Hoe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	IronSword: {
		DisplayName: "Iron Sword",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Weapon,
	},
	LapisLazuli: {
		DisplayName: "Lapis Lazuli",
		Stackable:   64,
		MaxHealth:   1,
		Category:    Item | DropItem | RawOre,
	},
	LapisOre: {
		DisplayName: "Lapis Ore",
		DropID:      LapisLazuli,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block | OreBlock,
	},
	LavaBucket: {
		DisplayName: "Lava Bucket",
		DropID:      Air,
		Stackable:   1,
		MaxHealth:   10,
		Category:    Item | Tool,
	},
	MelonSeeds: {
		DisplayName: "Melon Seeds",
		DropID:      Air,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Item | Unbreakable | DropItem,
	},
	MilkBucket: {
		DisplayName: "Milk Bucket",
		DropID:      Air,
		Stackable:   1,
		MaxHealth:   10,
		Category:    Item | Tool,
	},
	OakLeaves: {
		DisplayName: "Leaves",
		DropID:      OakLeaves,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},
	OakLog: {
		DisplayName: "Log",
		DropID:      OakLog,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},
	OakPlanks: {
		DisplayName: "Tree Plank",
		DropID:      OakPlanks,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},
	OakSapling: {
		DisplayName: "Sapling",
		DropID:      OakSapling,
		Stackable:   64,
		MaxHealth:   1,
		Category:    Block | Item,
	},
	Obsidian: {
		DisplayName: "Obsidian",
		DropID:      Obsidian,
		Stackable:   64,
		MaxHealth:   20,
		Category:    Block,
	},
	PowderSnowBucket: {
		DisplayName: "Powder Snow Bucket",
		DropID:      Air,
		Stackable:   1,
		MaxHealth:   1,
		Category:    Tool | Item,
	},
	PumpkinSeeds: {
		DisplayName: "Pumpkin Seeds",
		DropID:      Air,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Item | Unbreakable | DropItem,
	},
	RawCopper: {
		DisplayName: "Raw Copper",
		Stackable:   64,
		MaxHealth:   1,
		Category:    Item | DropItem | RawOre,
	},
	RawGold: {
		DisplayName: "Raw Gold",
		Stackable:   64,
		MaxHealth:   1,
		Category:    Item | DropItem | RawOre,
	},
	RawIron: {
		DisplayName: "Raw Iron",
		Stackable:   64,
		MaxHealth:   1,
		Category:    Item | DropItem | RawOre,
	},
	RootedDirt: {
		DisplayName: "Rooted Dirt",
		DropID:      RootedDirt,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},
	Sand: {
		DisplayName: "Sand",
		DropID:      Sand,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},

	SmoothStone: {
		DisplayName: "Smooth Stone",
		DropID:      SmoothStone,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},
	Snow: {
		DisplayName: "Snow",
		DropID:      Dirt,
		Stackable:   64,
		MaxHealth:   5,
		Category:    Block,
	},
	Snowball: {
		DisplayName: "Snowball",
		DropID:      Snowball,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Item,
	},
	Stick: {
		DisplayName: "Stick",
		Stackable:   64,
		MaxHealth:   1,
		Category:    Item,
	},

	Stone: {
		DisplayName: "Stone",
		DropID:      Cobblestone,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},
	StoneAxe: {
		DisplayName: "Stone Axe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	StoneBricks: {
		DisplayName: "Stone Bricks",
		DropID:      StoneBricks,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},
	StoneHoe: {
		DisplayName: "Stone Hoe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	StonePickaxe: {
		DisplayName: "Stone Pickaxe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	StoneShovel: {
		DisplayName: "Stone Shovel",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	StoneSword: {
		DisplayName: "Stone Sword",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Weapon,
	},
	Tnt: {
		DisplayName: "TNT",
		DropID:      Tnt,
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block,
	},
	Torch: {
		DisplayName: "Torch",
		Stackable:   64,
		MaxHealth:   10,
		Category:    Item,
	},
	TorchflowerSeeds: {
		DisplayName: "Torchflower Seeds",
		Stackable:   64,
		MaxHealth:   10,
		Category:    Item | Unbreakable | DropItem,
	},
	WaterBucket: {
		DisplayName: "Water Bucket",
		Stackable:   1,
		MaxHealth:   10,
		Category:    Item,
	},
	Wheat: {
		DisplayName: "Wheat",
		Stackable:   64,
		MaxHealth:   10,
		Category:    Item,
	},
	WheatCrops: {
		DisplayName: "Wheat Crops",
		Stackable:   64,
		MaxHealth:   10,
		Category:    Block | Harvestable,
	},
	WheatSeeds: {
		DisplayName: "Wheat Seeds",
		Stackable:   64,
		MaxHealth:   10,
		Category:    Item | DropItem,
	},
	WoodenAxe: {
		DisplayName: "Wooden Axe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	WoodenHoe: {
		DisplayName: "Wooden Hoe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	WoodenPickaxe: {
		DisplayName: "Wooden Pickaxe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	WoodenShovel: {
		DisplayName: "Wooden Hoe",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Tool,
	},
	WoodenSword: {
		DisplayName: "Wooden Sword",
		Stackable:   1,
		MaxHealth:   1,
		Category:    Item | Weapon,
	},
}
