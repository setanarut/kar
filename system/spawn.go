package system

import (
	"kar"
	"kar/arc"
	"kar/items"
	"kar/tilemap"
	"math/rand/v2"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	archeserde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/fastnoise"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/tilecollider"
)

var (
	player         ecs.Entity
	ctrl           *Controller
	tileMap        *tilemap.TileMap
	generator      *tilemap.Generator
	craftingTable  = items.NewCraftTable()
	collider       *tilecollider.Collider[uint8]
	toSpawn        = []arc.SpawnData{}
	toRemove       []ecs.Entity
	craftingState  bool
	craftingState4 bool
)

func AppendToSpawnList(x, y float64, id uint8, dur int) {
	toSpawn = append(
		toSpawn,
		arc.SpawnData{X: x - 4, Y: y - 4, Id: id, Durability: dur},
	)
}

func (s *Spawn) Init() {
	tileMap = tilemap.MakeTileMap(512, 512, 20, 20)
	generator = tilemap.NewGenerator(tileMap)
	generator.Opts.HighestSurfaceLevel = 10
	generator.Opts.LowestSurfaceLevel = 30
	generator.SetSeed(12)
	generator.NoiseState.FractalType(fastnoise.FractalFBm)
	generator.NoiseState.Frequency = 0.01
	generator.Generate()
	collider = tilecollider.NewCollider(
		tileMap.Grid,
		tileMap.TileW,
		tileMap.TileH,
	)
	ctrl = NewController(0, 10, collider)
	ctrl.Collider = collider
	// ctrl.Collider.StaticCheck = true
	ctrl.SkiddingJumpEnabled = true
	x, y := tileMap.FindSpawnPosition()
	// tileMap.Set(x, y+2, items.CraftingTable)
	SpawnX, SpawnY := tileMap.TileToWorldCenter(x, y)
	kar.Camera = kamera.NewCamera(SpawnX, SpawnY, kar.ScreenW, kar.ScreenH)
	kar.Camera.SmoothType = kamera.SmoothDamp
	kar.Camera.SmoothOptions.LerpSpeedX = 0.5
	kar.Camera.SmoothOptions.LerpSpeedY = 0.05
	kar.Camera.SetTopLeft(tileMap.FloorToBlockCenter(kar.Camera.TopLeft()))

	player = arc.SpawnPlayer(SpawnX, SpawnY)

}

func (s *Spawn) Update() {

	// Save TileMap image to desktop
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		tileMap.WriteToDesktopAsImage(playerTile.X, playerTile.Y)
	}
	// Save TileMap data to desktop
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		tileMap.WriteToDisk(kar.DesktopPath + "map")
	}
	// read TileMap from desktop
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		copy(tileMap.Grid, tileMap.ReadFromDisk(kar.DesktopPath+"map"))
	}

	// Generate tilemap
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		generator.Generate()
	}

	// Save ECSWorld to desktop
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		b, _ := archeserde.Serialize(
			&kar.WorldECS,
			// archeserde.Opts.SkipComponents(generic.T[anim.AnimationPlayer]()),
		)
		os.WriteFile(kar.DesktopPath+"data.json", b, 0644)

	}

	// // Read ECSWorld to desktop
	// if inpututil.IsKeyJustPressed(ebiten.Key4) {
	// 	kar.WorldECS.Reset()
	// 	b, _ := os.ReadFile(kar.DesktopPath + "data.json")
	// 	archeserde.Deserialize(b, &kar.WorldECS)

	// }

	// Remove Player
	if inpututil.IsKeyJustPressed(ebiten.Key5) {
		kar.WorldECS.RemoveEntity(player)
	}
	// Spawn Player
	if inpututil.IsKeyJustPressed(ebiten.Key6) {
		e := arc.SpawnPlayer(playerCenterX, playerCenterY)
		player = e
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		x, y := kar.Camera.ScreenToWorld(ebiten.CursorPosition())
		arc.SpawnEnemy(x, y, 0, 1)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := kar.Camera.ScreenToWorld(ebiten.CursorPosition())
		p := tileMap.WorldToTile(x, y)
		tileMap.Set(p.X, p.Y, items.Stone)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		x, y := kar.Camera.ScreenToWorld(ebiten.CursorPosition())
		p := tileMap.WorldToTile(x, y)
		tileMap.Set(p.X, p.Y, items.Air)
	}

	// Spawn item
	for _, spawnData := range toSpawn {
		arc.SpawnItem(spawnData, rand.IntN(sinspaceLen))
	}
	toSpawn = toSpawn[:0]

	for _, e := range toRemove {
		kar.WorldECS.RemoveEntity(e)
	}
	toRemove = toRemove[:0]

}
func (s *Spawn) Draw() {
	kar.Screen.Fill(kar.BackgroundColor)
}

type Spawn struct{}
