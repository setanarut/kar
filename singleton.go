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
	Screen      *ebiten.Image
	ScreenSize  = Vec{500, 340}
	WindowScale = 2.0
)
var (
	currentGameState  = "menu"
	previousGameState = "menu"
)
var (
	ceilBlockCoord   image.Point
	ceilBlockTick    float64
	dropItemHalfSize = Vec{4, 4}
)
var (
	world         ecs.World = ecs.NewWorld(100)
	currentPlayer ecs.Entity
	renderArea    = image.Point{(int(ScreenSize.X) / 20) + 3, (int(ScreenSize.Y) / 20) + 3}
	dataManager   *gdata.Manager
	// serdeOpt                    archeserde.Option
	sinspace                    []float64  = SinSpace(0, 2*math.Pi, 3, 60)
	drawItemHitboxEnabled       bool       = false
	drawPlayerTileHitboxEnabled bool       = false
	drawDebugTextEnabled        bool       = false
	backgroundColor             color.RGBA = color.RGBA{36, 36, 39, 255}
	tileCollider                *Collider
	gameTileMapGenerator        *tilemap.Generator
	animPlayer                  *anim.AnimationPlayer
	colorMDIO                   *colorm.DrawImageOptions = &colorm.DrawImageOptions{}
	colorM                      colorm.ColorM            = colorm.ColorM{}
	textDO                      *text.DrawOptions        = &text.DrawOptions{
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
	gameTileMapGenerator = tilemap.NewGenerator(tileMapRes)
	tileCollider = NewCollider(
		tileMapRes.Grid,
		tileMapRes.TileW,
		tileMapRes.TileH,
	)
	currentPlayer = SpawnPlayer(Vec{-5000, -5000})
}

func NewGame() {
	inventoryRes.Reset()

	*animPlayer.Data = animDefaultPlaybackData
	gameDataRes = &gameData{}
	gameTileMapGenerator.SetSeed(rand.Int())
	gameTileMapGenerator.Generate()

	spawnCoord := tileMapRes.FindSpawnPosition()
	SpawnPos := tileMapRes.TileToWorld(spawnCoord)

	cameraRes.SmoothOptions.LerpSpeedX = 0.5
	cameraRes.SmoothOptions.LerpSpeedY = 0
	cameraRes.SmoothType = kamera.SmoothDamp
	cameraRes.SetCenter(SpawnPos.X, SpawnPos.Y)

	currentPlayer = SpawnPlayer(SpawnPos)

	// debug
	spawnCoord.X++
	tileMapRes.Set(spawnCoord.X, spawnCoord.Y, items.Furnace)
	spawnCoord.X -= 2
	tileMapRes.Set(spawnCoord.X, spawnCoord.Y, items.CraftingTable)

	inventoryRes.SetSlot(0, items.Coal, 64, 0)
	inventoryRes.SetSlot(1, items.RawGold, 64, 0)
	inventoryRes.SetSlot(2, items.RawIron, 64, 0)
	inventoryRes.SetSlot(3, items.Stick, 64, 0)
	inventoryRes.SetSlot(4, items.DiamondPickaxe, 1, items.GetDefaultDurability(items.DiamondPickaxe))
	inventoryRes.SetSlot(5, items.DiamondShovel, 1, items.GetDefaultDurability(items.DiamondShovel))
	inventoryRes.SetSlot(6, items.DiamondAxe, 1, items.GetDefaultDurability(items.DiamondAxe))
	inventoryRes.SetSlot(7, items.Diamond, 64, 0)

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
