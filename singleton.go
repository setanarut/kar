package kar

import (
	"image"
	"image/color"
	"kar/items"
	"kar/tilemap"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	archeserde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/quasilyte/gdata"
	"github.com/setanarut/anim"
	"github.com/setanarut/fastnoise"
	"github.com/setanarut/kamera/v2"
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
	DropItemHalfSize  = Vec{4, 4}
	EnemyWormHalfSize = Vec{6, 6}
)
var (
	ECWorld                  ecs.World = ecs.NewWorld()
	CurrentPlayer            ecs.Entity
	DataManager              *gdata.Manager
	SerdeOpt                 archeserde.Option
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
	Update() error
	Draw()
}

func init() {

	var err error
	DataManager, err = gdata.Open(gdata.Config{AppName: "kar"})
	if err != nil {
		panic(err)
	}
	SerdeOpt = archeserde.Opts.SkipComponents(generic.T[anim.AnimationPlayer]())
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

func NewGame() {
	ECWorld.Reset()
	InventoryRes.Reset()
	ecs.AddResource(&ECWorld, GameDataRes)
	ecs.AddResource(&ECWorld, InventoryRes)
	ecs.AddResource(&ECWorld, CraftingTableRes)
	ecs.AddResource(&ECWorld, AnimPlayerDataRes)
	ecs.AddResource(&ECWorld, CameraRes)
	ecs.AddResource(&ECWorld, TileMapRes)

	GameTileMapGenerator.Opts.HighestSurfaceLevel = 10
	GameTileMapGenerator.Opts.LowestSurfaceLevel = 30
	GameTileMapGenerator.SetSeed(rand.Int())
	GameTileMapGenerator.NoiseState.FractalType(fastnoise.FractalFBm)
	GameTileMapGenerator.NoiseState.Frequency = 0.01
	GameTileMapGenerator.Generate()

	// ctrl.Collider.StaticCheck = true
	x, y := TileMapRes.FindSpawnPosition()
	// tileMap.Set(x, y+2, items.CraftingTable)
	SpawnX, SpawnY := TileMapRes.TileToWorldCenter(x, y)
	CameraRes.SmoothType = kamera.None
	CameraRes.SetCenter(SpawnX, SpawnY)
	CurrentPlayer = SpawnPlayer(SpawnX, SpawnY)
	CameraRes.SetTopLeft(TileMapRes.FloorToBlockCenter(CameraRes.X, CameraRes.Y))

	CameraRes.SmoothOptions.LerpSpeedX = 0.5
	CameraRes.SmoothOptions.LerpSpeedY = 0
	// CameraRes.SmoothType = kamera.SmoothDamp
	// CameraRes.SmoothOptions.SmoothDampTimeY = 0.1
	CameraRes.SmoothType = kamera.Lerp
}

func SaveGame() {
	FetchAnimPlayerData(PlayerAnimPlayer, AnimPlayerDataRes)
	jsonData, err := archeserde.Serialize(&ECWorld, SerdeOpt)
	if err != nil {
		log.Fatal("Error serializing world:", err)
	}
	DataManager.SaveItem("01save", jsonData)
}

func LoadGame() {
	if DataManager.ItemExists("01save") {
		if !ECWorld.Alive(CurrentPlayer) {
			CurrentPlayer = SpawnPlayer(0, 0)
		}
		ECWorld.Reset()
		ecs.AddResource(&ECWorld, GameDataRes)
		ecs.AddResource(&ECWorld, InventoryRes)
		ecs.AddResource(&ECWorld, CraftingTableRes)
		ecs.AddResource(&ECWorld, AnimPlayerDataRes)
		ecs.AddResource(&ECWorld, CameraRes)
		ecs.AddResource(&ECWorld, TileMapRes)
		jsonData, err := DataManager.LoadItem("01save")
		if err != nil {
			log.Fatal("Error loading saved data:", err)
		}

		err = archeserde.Deserialize(jsonData, &ECWorld)
		if err != nil {
			log.Fatal("Error deserializing world:", err)
		}
		animData := ecs.GetResource[AnimPlayerData](&ECWorld)
		SetAnimPlayerData(PlayerAnimPlayer, animData)
	}
}
