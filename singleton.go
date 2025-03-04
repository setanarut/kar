package kar

import (
	"image"
	"image/color"
	"kar/tilemap"
	"log"
	"math"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	arkserde "github.com/mlange-42/ark-serde"
	"github.com/mlange-42/ark/ecs"
	"github.com/quasilyte/gdata"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/v"
)

type Vec = v.Vec

const (
	SnowballGravity         float64       = 0.5
	SnowballSpeedX          float64       = 3.5
	SnowballMaxFallVelocity float64       = 2.5
	SnowballBounceHeight    float64       = 9.0
	ItemGravity             float64       = 3.0
	PlayerBestToolDamage    float64       = 5.0
	PlayerDefaultDamage     float64       = 1.0
	RaycastDist             int           = 4 // block unit
	Tick                    time.Duration = time.Second / 60
	ItemCollisionDelay      time.Duration = time.Second / 2
)

var debugEnabled bool = false

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
	ceilBlockCoord image.Point
	ceilBlockTick  float64
	dropItemAABB   = &AABB{Half: Vec{4, 4}}
)
var (
	world                ecs.World = ecs.NewWorld(100)
	currentPlayer        ecs.Entity
	renderArea           = image.Point{(int(ScreenSize.X) / 20) + 3, (int(ScreenSize.Y) / 20) + 3}
	dataManager          *gdata.Manager
	sinspace             []float64  = Sinspace(0, 2*math.Pi, 3, 60)
	backgroundColor      color.RGBA = color.RGBA{36, 36, 39, 255}
	tileCollider         *Collider
	gameTileMapGenerator *tilemap.Generator
	colorMDIO            *colorm.DrawImageOptions = &colorm.DrawImageOptions{}
	colorM               colorm.ColorM            = colorm.ColorM{}
	textDO               *text.DrawOptions        = &text.DrawOptions{
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
	cameraRes.SmoothOptions.LerpSpeedX = 0.5
	cameraRes.SmoothOptions.LerpSpeedY = 0
	var err error
	dataManager, err = gdata.Open(gdata.Config{AppName: "kar"})
	if err != nil {
		panic(err)
	}
	gameTileMapGenerator = tilemap.NewGenerator(tileMapRes)
	tileCollider = NewCollider(
		tileMapRes.Grid,
		tileMapRes.TileW,
		tileMapRes.TileH,
	)
}

func NewGame() {
	debugEnabled = false
	world.Reset()
	inventoryRes.Reset()

	inventoryResMap.Add(inventoryRes)
	craftingtableResMap.Add(craftingTableRes)
	cameraResMap.Add(cameraRes)
	gameDataResMap.Add(gameDataRes)
	animPlaybackDataResMap.Add(animPlayer.Data)
	tilemapResMap.Add(tileMapRes)

	*animPlayer.Data = animDefaultPlaybackData
	gameDataRes = &gameData{}
	gameTileMapGenerator.SetSeed(rand.Int())
	gameTileMapGenerator.Generate()

	spawnCoord := tileMapRes.FindSpawnPosition()
	SpawnPos := tileMapRes.TileToWorld(spawnCoord)
	currentPlayer = SpawnPlayer(SpawnPos)
	box := mapAABB.Get(currentPlayer)
	box.SetBottom(tileMapRes.GetTileBottom(spawnCoord.X, spawnCoord.Y))
	cameraRes.SmoothType = kamera.SmoothDamp
	cameraRes.SetCenter(box.Pos.X, box.Pos.Y)
}

func SaveGame() {
	jsonData, err := arkserde.Serialize(&world)
	if err != nil {
		log.Fatal("Error serializing world:", err)
	}
	dataManager.SaveItem("01save", jsonData)

}

func LoadGame() {
	debugEnabled = false
	if dataManager.ItemExists("01save") {
		world.Reset()

		inventoryResMap.Add(inventoryRes)
		craftingtableResMap.Add(craftingTableRes)
		cameraResMap.Add(cameraRes)
		gameDataResMap.Add(gameDataRes)
		animPlaybackDataResMap.Add(animPlayer.Data)
		tilemapResMap.Add(tileMapRes)

		jsonData, err := dataManager.LoadItem("01save")
		if err != nil {
			log.Fatal("Error loading saved data:", err)
		}
		err = arkserde.Deserialize(jsonData, &world)
		if err != nil {
			log.Fatal("Error deserializing world:", err)
		}

		if !world.Alive(currentPlayer) {
			q := filterPlayer.Query()
			q.Next()
			currentPlayer = q.Entity()
			q.Close()
		}

		animPlayer.Update()
	}
}
