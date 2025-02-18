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
	GameDataRes      *GameData
	TileMapRes       *tilemap.TileMap
	CraftingTableRes *items.CraftTable
	InventoryRes     *items.Inventory
	// AnimPlayerDataRes *anim.PlaybackData
	CameraRes *kamera.Camera
)
var AnimDefaultPlaybackData anim.PlaybackData

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
	PlayerAnimPlayer.NewAnim("idleRight", 0, 0, 16, 16, 1, false, false, 1)
	PlayerAnimPlayer.NewAnim("idleUp", 208, 0, 16, 16, 1, false, false, 1)
	PlayerAnimPlayer.NewAnim("idleDown", 224, 0, 16, 16, 1, false, false, 1)
	PlayerAnimPlayer.NewAnim("walkRight", 16, 0, 16, 16, 4, false, false, 15)
	PlayerAnimPlayer.NewAnim("jump", 16*5, 0, 16, 16, 1, false, false, 15)
	PlayerAnimPlayer.NewAnim("skidding", 16*6, 0, 16, 16, 1, false, false, 15)
	PlayerAnimPlayer.NewAnim("attackDown", 16*7, 0, 16, 16, 2, false, false, 8)
	PlayerAnimPlayer.NewAnim("attackRight", 144, 0, 16, 16, 2, false, false, 8)
	PlayerAnimPlayer.NewAnim("attackWalk", 0, 16, 16, 16, 4, false, false, 8)
	PlayerAnimPlayer.NewAnim("attackUp", 16*11, 0, 16, 16, 2, false, false, 8)
	PlayerAnimPlayer.SetAnim("idleRight")
	AnimDefaultPlaybackData = *PlayerAnimPlayer.Data
	// AnimPlayerDataRes = &anim.PlaybackData{}
	// *AnimPlayerDataRes = *PlayerAnimPlayer.Data
}

type GameData struct {
	CraftingState    bool
	CraftingState4   bool
	TargetBlockCoord image.Point
}
