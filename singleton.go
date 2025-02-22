package kar

import (
	"image"
	"image/color"
	"kar/items"
	"kar/tilemap"
	"kar/v"
	"math"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	archeserde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/ark/ecs"
	"github.com/quasilyte/gdata"
	"github.com/setanarut/anim"
	"github.com/setanarut/kamera/v2"
)

type Vec = v.Vec

const (
	SnowballGravity             = 0.5
	SnowballSpeedX              = 3.5
	SnowballMaxFallVelocity     = 2.5
	SnowballBounceHeight        = 9.0
	ItemGravity                 = 3.0
	PlayerBestToolDamage        = 5.0
	PlayerDefaultDamage         = 1.0
	Tick                        = time.Second / 60
	ItemCollisionDelay          = time.Second / 2
	RaycastDist             int = 4 // block unit
)

var (
	WindowScale = 2.0
	ScreenSize  = Vec{500, 340}
)
var (
	currentGameState  = "menu"
	previousGameState = "menu"
)
var (
	ceilBlockCoord    image.Point
	ceilBlockTick     float64
	dropItemHalfSize  = Vec{4, 4}
	enemyWormHalfSize = Vec{6, 6}
)
var (
	world                       ecs.World = ecs.NewWorld(100)
	currentPlayer               ecs.Entity
	renderArea                  = image.Point{(int(ScreenSize.X) / 20) + 3, (int(ScreenSize.Y) / 20) + 3}
	dataManager                 *gdata.Manager
	serdeOpt                    archeserde.Option
	Sinspace                    []float64  = SinSpace(0, 2*math.Pi, 3, 60)
	DrawItemHitboxEnabled       bool       = false
	DrawPlayerTileHitboxEnabled bool       = false
	DrawDebugTextEnabled        bool       = false
	BackgroundColor             color.RGBA = color.RGBA{36, 36, 39, 255}
	TileCollider                *Collider
	GameTileMapGenerator        *tilemap.Generator
	animPlayer                  *anim.AnimationPlayer
	Screen                      *ebiten.Image
	ColorMDIO                   *colorm.DrawImageOptions = &colorm.DrawImageOptions{}
	ColorM                      colorm.ColorM            = colorm.ColorM{}
	TextDO                      *text.DrawOptions        = &text.DrawOptions{
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
	var err error
	dataManager, err = gdata.Open(gdata.Config{AppName: "kar"})
	if err != nil {
		panic(err)
	}

	// serdeOpt = archeserde.Opts.SkipComponents(generic.T[anim.AnimationPlayer]())
	GameTileMapGenerator = tilemap.NewGenerator(tileMapRes)
	TileCollider = NewCollider(
		tileMapRes.Grid,
		tileMapRes.TileW,
		tileMapRes.TileH,
	)

	inventoryRes.SetSlot(0, items.Snowball, 64, 0)
	currentPlayer = SpawnPlayer(-5000, -5000)
}

func NewGame() {
	// world.Reset()
	// inventoryRes.Reset()
	*animPlayer.Data = animDefaultPlaybackData
	gameDataRes = &gameData{}

	inventoryResMap.Add(inventoryRes)
	tilemapResMap.Add(tileMapRes)
	craftingtableResMap.Add(craftingTableRes)
	cameraResMap.Add(cameraRes)
	gameDataResMap.Add(gameDataRes)
	animPlaybackDataResMap.Add(animPlayer.Data)

	GameTileMapGenerator.SetSeed(rand.Int())
	GameTileMapGenerator.Generate()
	x, y := tileMapRes.FindSpawnPosition()
	SpawnX, SpawnY := tileMapRes.TileToWorldCenter(x, y)
	cameraRes.SmoothType = kamera.None
	cameraRes.SetCenter(SpawnX, SpawnY)
	currentPlayer = SpawnPlayer(SpawnX, SpawnY)
	cameraRes.SetTopLeft(tileMapRes.FloorToBlockCenter(cameraRes.X, cameraRes.Y))
	cameraRes.SmoothOptions.LerpSpeedX = 0.5
	cameraRes.SmoothOptions.LerpSpeedY = 0
	cameraRes.SmoothType = kamera.SmoothDamp
}

func SaveGame() {
	// jsonData, err := archeserde.Serialize(&world, serdeOpt)
	// if err != nil {
	// 	log.Fatal("Error serializing world:", err)
	// }
	// dataManager.SaveItem("01save", jsonData)

}

func LoadGame() {
	// if dataManager.ItemExists("01save") {
	// 	world.Reset()
	// 	ecs.AddResource(&world, gameDataRes)
	// 	ecs.AddResource(&world, inventoryRes)
	// 	ecs.AddResource(&world, craftingTableRes)
	// 	ecs.AddResource(&world, animPlayer.Data)
	// 	ecs.AddResource(&world, cameraRes)
	// 	ecs.AddResource(&world, tileMapRes)
	// 	jsonData, err := dataManager.LoadItem("01save")
	// 	if err != nil {
	// 		log.Fatal("Error loading saved data:", err)
	// 	}
	// 	err = archeserde.Deserialize(jsonData, &world)
	// 	if err != nil {
	// 		log.Fatal("Error deserializing world:", err)
	// 	}

	// 	animPlayer.Update()
	// }
}
