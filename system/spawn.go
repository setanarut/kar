package system

import (
	"kar"
	"kar/arc"
	"kar/items"
	"kar/tilemap"

	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/tilecollider"
)

var (
	TempFallingY  float64
	PlayerEntity  ecs.Entity
	CTRL          *Controller
	Map           *tilemap.TileMap
	Collider      *tilecollider.Collider[uint16]
	ToSpawn       = []arc.SpawnData{}
	toRemove      []ecs.Entity
	craftingState bool
)

func AppendToSpawnList(x, y float64, id uint16, durability int) {
	ToSpawn = append(ToSpawn, arc.SpawnData{
		X:          x - 4*kar.ItemScale,
		Y:          y - 4*kar.ItemScale,
		Id:         id,
		Durability: durability,
	})
}

func (s *Spawn) Init() {

	Map = tilemap.MakeTileMap(512, 512, int(20*kar.ScreenScale), int(20*kar.ScreenScale))
	tilemap.Generate(Map)
	Collider = tilecollider.NewCollider(Map.Grid, Map.TileW, Map.TileH)
	CTRL = NewController(0, 10, Collider)
	CTRL.Collider = Collider
	CTRL.SetScale(kar.ScreenScale)
	CTRL.SkiddingJumpEnabled = true
	x, y := Map.FindSpawnPosition()
	Map.SetTileID(x, y+2, items.CraftingTable)
	SpawnX, SpawnY := Map.TileToWorldCenter(x, y)
	kar.Camera.LookAt(SpawnX, SpawnY)
	kar.Camera.SetTopLeft(Map.FloorToBlockCenter(kar.Camera.TopLeft()))
	PlayerEntity = arc.SpawnPlayer(SpawnX, SpawnY)
	anim, hlt, dop, rect, inv := arc.MapPlayer.Get(PlayerEntity)
	CTRL.AnimPlayer = anim
	CTRL.Rect = rect
	CTRL.Health = hlt
	CTRL.DOP = dop
	CTRL.Inventory = inv
	CTRL.EnterFalling()

}

func (s *Spawn) Update() {
	// Spawn item
	for _, spawnData := range ToSpawn {
		arc.SpawnItem(spawnData)
	}
	ToSpawn = ToSpawn[:0]

	for _, e := range toRemove {
		kar.WorldECS.RemoveEntity(e)
	}
	toRemove = toRemove[:0]

}
func (s *Spawn) Draw() {
	kar.Screen.Fill(kar.BackgroundColor)
}

type Spawn struct{}
