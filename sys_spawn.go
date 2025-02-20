package kar

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
)

// SpawnData is a helper for delaying spawn events
type SpawnData struct {
	X, Y       float64
	Id         uint8
	Durability int
}

var (
	toSpawn  = []SpawnData{}
	toRemove []ecs.Entity
)

type Spawn struct {
	tile uint8
}

func (s *Spawn) Init() {

}

func (s *Spawn) Update() error {

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := CameraRes.ScreenToWorld(ebiten.CursorPosition())
		p := TileMapRes.WorldToTile(x, y)
		TileMapRes.Set(p.X, p.Y, s.tile)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		x, y := CameraRes.ScreenToWorld(ebiten.CursorPosition())
		p := TileMapRes.WorldToTile(x, y)
		s.tile = TileMapRes.Get(p.X, p.Y)
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
	return nil
}
func (s *Spawn) Draw() {
}
