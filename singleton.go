package kar

import (
	"image"
	"image/color"
	"kar/items"
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
	currentGameState = "mainmenu"
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
	sinspaceOffsets      []float64  = sinspace(0, 2*math.Pi, 3, 60)
	backgroundColor      color.RGBA = color.RGBA{36, 36, 39, 255}
	gameTileMapGenerator tilemap.Generator
	colorMDIO            *colorm.DrawImageOptions = &colorm.DrawImageOptions{}
	colorM               colorm.ColorM            = colorm.ColorM{}
	textDO               *text.DrawOptions        = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions: text.LayoutOptions{
			LineSpacing: 10,
		},
	}
	tileCollider = Collider{
		TileMap:  tileMapRes.Grid,
		TileSize: image.Point{tileMapRes.TileW, tileMapRes.TileH},
	}
)

func init() {
	var err error
	dataManager, err = gdata.Open(gdata.Config{AppName: "kar"})
	if err != nil {
		panic(err)
	}
	gameTileMapGenerator = tilemap.NewGenerator(tileMapRes)

}

func NewGame() {
	world.Reset()
	inventoryRes.Reset()
	gameDataRes = gameData{}
	animPlayer.Data = animDefaultPlaybackData

	mapResInventory.Add(inventoryRes)
	mapResCraftingtable.Add(craftingTableRes)
	mapResCamera.Add(cameraRes)
	mapResGameData.Add(&gameDataRes)
	mapResAnimPlaybackData.Add(&animPlayer.Data)
	mapResTilemap.Add(&tileMapRes)

	gameTileMapGenerator.SetSeed(rand.Int())
	gameTileMapGenerator.Generate()
	spawnCoord := tileMapRes.FindSpawnPosition()
	SpawnPos := tileMapRes.TileToWorld(spawnCoord)
	currentPlayer = SpawnPlayer(SpawnPos)
	box := mapAABB.Get(currentPlayer)
	box.SetBottom(tileMapRes.GetTileBottom(spawnCoord.X, spawnCoord.Y))
	cameraRes.SmoothType = kamera.SmoothDamp
	cameraRes.SetCenter(box.Pos.X, box.Pos.Y)
	inventoryRes.SetSlot(8, items.Snowball, 64, 0)
}

func SaveGame() {
	jsonData, err := arkserde.Serialize(&world, arkserde.Opts.Compress())
	if err != nil {
		log.Fatal(err)
	}
	dataManager.SaveItem("01save", jsonData)

}

func LoadGame() {
	if dataManager.ItemExists("01save") {
		world.Reset()
		mapResInventory.Add(inventoryRes)
		mapResCraftingtable.Add(craftingTableRes)
		mapResCamera.Add(cameraRes)
		mapResGameData.Add(&gameDataRes)
		mapResAnimPlaybackData.Add(&animPlayer.Data)
		mapResTilemap.Add(&tileMapRes)

		jsonData, err := dataManager.LoadItem("01save")
		if err != nil {
			log.Fatal(err)
		}
		err = arkserde.Deserialize(jsonData, &world, arkserde.Opts.Compress())
		if err != nil {
			log.Fatal(err)
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
