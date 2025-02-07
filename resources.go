package kar

import (
	"image"
	"kar/items"
	"kar/res"
	"kar/tilemap"

	"github.com/setanarut/anim"
	"github.com/setanarut/kamera/v2"
)

// ECS Resources
var (
	GameDataRes       *GameData
	TileMapRes        *tilemap.TileMap
	CraftingTableRes  *items.CraftTable
	InventoryRes      *items.Inventory
	AnimPlayerDataRes *AnimPlayerData
	CameraRes         *kamera.Camera
)

func init() {
	GameDataRes = &GameData{CraftingState: false, CraftingState4: false}
	CraftingTableRes = items.NewCraftTable()
	InventoryRes = items.NewInventory()
	TileMapRes = tilemap.MakeTileMap(512, 512, 20, 20)
	CameraRes = kamera.NewCamera(0, 0, ScreenW, ScreenH)
	PlayerAnimPlayer = anim.NewAnimationPlayer(
		&anim.Atlas{"Default", res.Player},
		&anim.Atlas{"WoodenAxe", res.PlayerWoodenAxeAtlas},
		&anim.Atlas{"StoneAxe", res.PlayerStoneAxeAtlas},
		&anim.Atlas{"IronAxe", res.PlayerIronAxeAtlas},
		&anim.Atlas{"DiamondAxe", res.PlayerDiamondAxeAtlas},
		&anim.Atlas{"WoodenPickaxe", res.PlayerWoodenPickaxeAtlas},
		&anim.Atlas{"StonePickaxe", res.PlayerStonePickaxeAtlas},
		&anim.Atlas{"IronPickaxe", res.PlayerIronPickaxeAtlas},
		&anim.Atlas{"DiamondPickaxe", res.PlayerDiamondPickaxeAtlas},
		&anim.Atlas{"WoodenShovel", res.PlayerWoodenShovelAtlas},
		&anim.Atlas{"StoneShovel", res.PlayerStoneShovelAtlas},
		&anim.Atlas{"IronShovel", res.PlayerIronShovelAtlas},
		&anim.Atlas{"DiamondShovel", res.PlayerDiamondShovelAtlas},
	)
	PlayerAnimPlayer.NewState("idleRight", 0, 0, 16, 16, 1, false, false, 1)
	PlayerAnimPlayer.NewState("idleUp", 208, 0, 16, 16, 1, false, false, 1)
	PlayerAnimPlayer.NewState("idleDown", 224, 0, 16, 16, 1, false, false, 1)
	PlayerAnimPlayer.NewState("walkRight", 16, 0, 16, 16, 4, false, false, 15)
	PlayerAnimPlayer.NewState("jump", 16*5, 0, 16, 16, 1, false, false, 15)
	PlayerAnimPlayer.NewState("skidding", 16*6, 0, 16, 16, 1, false, false, 15)
	PlayerAnimPlayer.NewState("attackDown", 16*7, 0, 16, 16, 2, false, false, 8)
	PlayerAnimPlayer.NewState("attackRight", 144, 0, 16, 16, 2, false, false, 8)
	PlayerAnimPlayer.NewState("attackWalk", 0, 16, 16, 16, 4, false, false, 8)
	PlayerAnimPlayer.NewState("attackUp", 16*11, 0, 16, 16, 2, false, false, 8)
	PlayerAnimPlayer.CurrentState = "idleRight"
	AnimPlayerDataRes = NewAnimPlayerData(PlayerAnimPlayer)
}

type GameData struct {
	CraftingState    bool
	CraftingState4   bool
	TargetBlockCoord image.Point
}

func NewAnimPlayerData(ap *anim.AnimationPlayer) *AnimPlayerData {
	return &AnimPlayerData{
		CurrentState: ap.CurrentState,
		CurrentAtlas: ap.CurrentAtlas,
		Paused:       ap.Paused,
		Tick:         ap.Tick,
		CurrentIndex: ap.CurrentIndex,
	}
}

func FetchAnimPlayerData(ap *anim.AnimationPlayer, data *AnimPlayerData) {
	data.CurrentState = ap.CurrentState
	data.CurrentAtlas = ap.CurrentAtlas
	data.Paused = ap.Paused
	data.Tick = ap.Tick
	data.CurrentIndex = ap.CurrentIndex
}

func SetAnimPlayerData(ap *anim.AnimationPlayer, data *AnimPlayerData) {
	ap.CurrentState = data.CurrentState
	ap.CurrentAtlas = data.CurrentAtlas
	ap.Paused = data.Paused
	ap.Tick = data.Tick
	ap.CurrentIndex = data.CurrentIndex
}

type AnimPlayerData struct {
	CurrentState string
	CurrentAtlas string
	Paused       bool
	Tick         float64
	CurrentIndex int
}
