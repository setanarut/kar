package kar

import (
	"kar/items"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
)

var (
	toSpawn  = []SpawnData{}
	toRemove []ecs.Entity
)

func AppendToSpawnList(x, y float64, id uint8, dur int) {
	toSpawn = append(
		toSpawn,
		SpawnData{X: x - 4, Y: y - 4, Id: id, Durability: dur},
	)
}

type Spawn struct{}

func (s *Spawn) Init() {

}

func (s *Spawn) Update() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := CameraRes.ScreenToWorld(ebiten.CursorPosition())
		p := TileMapRes.WorldToTile(x, y)
		TileMapRes.Set(p.X, p.Y, items.Stone)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		x, y := CameraRes.ScreenToWorld(ebiten.CursorPosition())
		p := TileMapRes.WorldToTile(x, y)
		TileMapRes.Set(p.X, p.Y, items.Air)
	}

	// Spawn item
	for _, d := range toSpawn {
		SpawnItem(d.X, d.Y, d.Id, d.Durability)
	}

	toSpawn = toSpawn[:0]

	for _, e := range toRemove {
		ECWorld.RemoveEntity(e)
	}
	toRemove = toRemove[:0]

}
func (s *Spawn) Draw() {
}
