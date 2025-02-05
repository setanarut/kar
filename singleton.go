package kar

import (
	"image"
	"image/color"
	"kar/engine/mathutil"
	"kar/items"
	"kar/res"
	"kar/tilemap"
	"log"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/anim"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/tilecollider"
)

const (
	SnowballGravity         = 0.5
	SnowballSpeedX          = 3.5
	SnowballMaxFallVelocity = 2.5
	SnowballBounceHeight    = 9

	Tick = time.Second / 60
)

var (
	CurrentGameState  = "menu"
	PreviousGameState = "menu"
)

type ISystem interface {
	Init()
	Update()
	Draw()
}

var (
	DesktopPath      string
	WindowScale      float64 = 2.0
	ScreenW, ScreenH float64 = 500.0, 340.0

	ItemCollisionDelay time.Duration = time.Second / 2
	RaycastDist        int           = 4 // block unit
	ItemGravity        float64       = 3
	Sinspace                         = mathutil.SinSpace(0, 2*math.Pi, 3, 60)

	DrawDebugHitboxesEnabled bool        = false
	DrawDebugTextEnabled     bool        = false
	PlayerBestToolDamage                 = 5.0
	PlayerDefaultDamage                  = 1.0
	RenderArea               image.Point = image.Point{
		(int(ScreenW) / 20) + 3,
		(int(ScreenH) / 20) + 3,
	}

	BackgroundColor      color.RGBA = color.RGBA{36, 36, 39, 255}
	GameTileMapGenerator *tilemap.Generator
	PlayerAnimPlayer     *anim.AnimationPlayer

	ECWorld       = ecs.NewWorld()
	CurrentPlayer ecs.Entity
	Collider      *tilecollider.Collider[uint8]

	// ECS Resources
	TileMapRes        *tilemap.TileMap
	GameDataRes       *GameData
	CraftingTableRes  *items.CraftTable
	InventoryRes      *items.Inventory
	AnimPlayerDataRes *AnimPlayerData
	CameraRes         *kamera.Camera

	Screen    *ebiten.Image
	ColorMDIO = &colorm.DrawImageOptions{}
	ColorM    = colorm.ColorM{}
	TextDO    = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions: text.LayoutOptions{
			LineSpacing: 10,
		},
	}
	// Debug
)

func init() {
	CraftingTableRes = items.NewCraftTable()
	InventoryRes = items.NewInventory()
	CameraRes = kamera.NewCamera(0, 0, ScreenW, ScreenH)
	TileMapRes = tilemap.MakeTileMap(512, 512, 20, 20)
	GameTileMapGenerator = tilemap.NewGenerator(TileMapRes)
	Collider = tilecollider.NewCollider(
		TileMapRes.Grid,
		TileMapRes.TileW,
		TileMapRes.TileH,
	)
	GameDataRes = &GameData{
		CraftingState:  false,
		CraftingState4: false,
		DropItemW:      8.0,
		DropItemH:      8.0,
	}
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	DesktopPath = homePath + "/Desktop/"

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

	// ECS Resources
	ecs.AddResource(&ECWorld, InventoryRes)
	ecs.AddResource(&ECWorld, CraftingTableRes)
	ecs.AddResource(&ECWorld, AnimPlayerDataRes)
	ecs.AddResource(&ECWorld, GameDataRes)
	ecs.AddResource(&ECWorld, TileMapRes)
	ecs.AddResource(&ECWorld, CameraRes)

	InventoryRes.SetSlot(0, items.Snowball, 64, 0)
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

type GameData struct {
	CraftingState        bool
	CraftingState4       bool
	TargetBlockCoord     image.Point
	DropItemW, DropItemH float64
}
