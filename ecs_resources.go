package kar

import (
	"image"
	"kar/items"
	"kar/res"
	"kar/tilemap"
	"time"

	"github.com/mlange-42/ark/ecs"
	"github.com/setanarut/anim"
	"github.com/setanarut/kamera/v2"
)

// ECS Resources
var (
	animPlayer              *anim.AnimationPlayer
	animDefaultPlaybackData anim.PlaybackData

	cameraRes        *kamera.Camera
	craftingTableRes *items.CraftTable
	gameDataRes      *gameData
	inventoryRes     *items.Inventory
	tileMapRes       *tilemap.TileMap

	mapResAnimPlaybackData = ecs.NewResource[anim.PlaybackData](&world)
	mapResCamera           = ecs.NewResource[kamera.Camera](&world)
	mapRescraftingtable    = ecs.NewResource[items.CraftTable](&world)
	mapResgameData         = ecs.NewResource[gameData](&world)
	mapResinventory        = ecs.NewResource[items.Inventory](&world)
	mapRestilemap          = ecs.NewResource[tilemap.TileMap](&world)
)

// GameplayStates
const (
	Playing int = iota
	CraftingTable3x3
	Crafting2x2
	Furnace1x2
)

type gameData struct {
	GameplayState    int
	TargetBlockCoord image.Point
	IsRayHit         bool
	BlockHealth      float64
	Duration         time.Duration // Gameplay duration
	SpawnElapsed     time.Duration // Entity spawn timer
}

func init() {
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
	gameDataRes = &gameData{GameplayState: Playing}
	craftingTableRes = items.NewCraftTable()
	inventoryRes = items.NewInventory()
	cameraRes = kamera.NewCamera(0, 0, ScreenSize.X, ScreenSize.Y)
	tileMapRes = tilemap.MakeTileMap(512, 512, 20, 20)
}
