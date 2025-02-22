package kar

import (
	"image"
	"kar/items"
	"kar/res"
	"kar/tilemap"

	"github.com/mlange-42/ark/ecs"
	"github.com/setanarut/anim"
	"github.com/setanarut/kamera/v2"
)

// ECS Resources
var (
	gameDataRes      *gameData
	tileMapRes       *tilemap.TileMap
	craftingTableRes *items.CraftTable
	inventoryRes     *items.Inventory
	cameraRes        *kamera.Camera

	inventoryResMap        = ecs.NewResource[items.Inventory](&world)
	tilemapResMap          = ecs.NewResource[tilemap.TileMap](&world)
	craftingtableResMap    = ecs.NewResource[items.CraftTable](&world)
	cameraResMap           = ecs.NewResource[kamera.Camera](&world)
	gameDataResMap         = ecs.NewResource[gameData](&world)
	animPlaybackDataResMap = ecs.NewResource[anim.PlaybackData](&world)
)
var animDefaultPlaybackData anim.PlaybackData

// GameplayStates
const (
	Playing int = iota
	CraftingTable3x3
	CraftingTable2x2
	Furnace
)

type gameData struct {
	GameplayState    int
	TargetBlockCoord image.Point
	IsRayHit         bool
	BlockHealth      float64
}

func init() {
	gameDataRes = &gameData{GameplayState: Playing}
	craftingTableRes = items.NewCraftTable()
	inventoryRes = items.NewInventory()
	tileMapRes = tilemap.MakeTileMap(512, 512, 20, 20)
	cameraRes = kamera.NewCamera(0, 0, ScreenSize.X, ScreenSize.Y)
	animPlayer = anim.NewAnimationPlayer(
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
	animPlayer.NewAnim("idleRight", 0, 0, 16, 16, 1, false, false, 1)
	animPlayer.NewAnim("idleUp", 208, 0, 16, 16, 1, false, false, 1)
	animPlayer.NewAnim("idleDown", 224, 0, 16, 16, 1, false, false, 1)
	animPlayer.NewAnim("walkRight", 16, 0, 16, 16, 4, false, false, 15)
	animPlayer.NewAnim("jump", 16*5, 0, 16, 16, 1, false, false, 15)
	animPlayer.NewAnim("skidding", 16*6, 0, 16, 16, 1, false, false, 15)
	animPlayer.NewAnim("attackDown", 16*7, 0, 16, 16, 2, false, false, 8)
	animPlayer.NewAnim("attackRight", 144, 0, 16, 16, 2, false, false, 8)
	animPlayer.NewAnim("attackWalk", 0, 16, 16, 16, 4, false, false, 8)
	animPlayer.NewAnim("attackUp", 16*11, 0, 16, 16, 2, false, false, 8)
	animPlayer.SetAnim("idleRight")
	animDefaultPlaybackData = *animPlayer.Data

	inventoryResMap.Add(inventoryRes)
	tilemapResMap.Add(tileMapRes)
	craftingtableResMap.Add(craftingTableRes)
	cameraResMap.Add(cameraRes)
	gameDataResMap.Add(gameDataRes)
	animPlaybackDataResMap.Add(animPlayer.Data)
}
