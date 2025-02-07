package kar

import (
	"image"
	"image/color"
	"kar/items"
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
	"github.com/setanarut/tilecollider"
)

const (
	SnowballGravity                       = 0.5
	SnowballSpeedX                        = 3.5
	SnowballMaxFallVelocity               = 2.5
	SnowballBounceHeight                  = 9
	Tick                                  = time.Second / 60
	ItemGravity             float64       = 3
	PlayerBestToolDamage    float64       = 5.0
	PlayerDefaultDamage     float64       = 1.0
	ItemCollisionDelay      time.Duration = time.Second / 2
	RaycastDist             int           = 4 // block unit
)

var (
	CurrentGameState  = "menu"
	PreviousGameState = "menu"
)
var (
	DropItemSize  Size = Size{8, 8}
	EnemyWormSize Size = Size{8, 8}
)
var (
	ECWorld                  ecs.World = ecs.NewWorld()
	CurrentPlayer            ecs.Entity
	DesktopPath              string
	WindowScale              float64     = 2.0
	ScreenW, ScreenH         float64     = 500.0, 340.0
	Sinspace                 []float64   = SinSpace(0, 2*math.Pi, 3, 60)
	RenderArea               image.Point = image.Point{(int(ScreenW) / 20) + 3, (int(ScreenH) / 20) + 3}
	DrawDebugHitboxesEnabled bool        = false
	DrawDebugTextEnabled     bool        = false
	BackgroundColor          color.RGBA  = color.RGBA{36, 36, 39, 255}
	Collider                 *tilecollider.Collider[uint8]
	GameTileMapGenerator     *tilemap.Generator
	PlayerAnimPlayer         *anim.AnimationPlayer
	Screen                   *ebiten.Image
	ColorMDIO                = &colorm.DrawImageOptions{}
	ColorM                   = colorm.ColorM{}
	TextDO                   = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions: text.LayoutOptions{
			LineSpacing: 10,
		},
	}
)

type ISystem interface {
	Init()
	Update()
	Draw()
}

func init() {
	GameTileMapGenerator = tilemap.NewGenerator(TileMapRes)
	Collider = tilecollider.NewCollider(
		TileMapRes.Grid,
		TileMapRes.TileW,
		TileMapRes.TileH,
	)
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	DesktopPath = homePath + "/Desktop/"

	// ECS Resources
	ecs.AddResource(&ECWorld, InventoryRes)
	ecs.AddResource(&ECWorld, CraftingTableRes)
	ecs.AddResource(&ECWorld, AnimPlayerDataRes)
	ecs.AddResource(&ECWorld, GameDataRes)
	ecs.AddResource(&ECWorld, TileMapRes)
	ecs.AddResource(&ECWorld, CameraRes)

	InventoryRes.SetSlot(0, items.Snowball, 64, 0)
}
