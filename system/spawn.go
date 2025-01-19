package system

import (
	"kar"
	"kar/arc"
	"kar/items"
	"kar/tilemap"
	"math/rand/v2"

	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/tilecollider"
)

var (
	player         ecs.Entity
	ctrl           *Controller
	tileMap        *tilemap.TileMap
	craftingTable  = items.NewCraftTable()
	collider       *tilecollider.Collider[uint16]
	toSpawn        = []arc.SpawnData{}
	toRemove       []ecs.Entity
	craftingState  bool
	craftingState4 bool
)

func AppendToSpawnList(x, y float64, id uint16, dur int) {
	toSpawn = append(
		toSpawn,
		arc.SpawnData{X: x - 4, Y: y - 4, Id: id, Durability: dur},
	)
}

func (s *Spawn) Init() {
	tileMap = tilemap.MakeTileMap(512, 512, 20, 20)
	tilemap.Generate(tileMap)
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

	arc.SpawnEnemy(SpawnX, SpawnY-40, -1, 0)

}

func (s *Spawn) Update() {
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
