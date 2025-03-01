package kar

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/ark/ecs"
)

// spawnData is a helper for delaying spawn events
type spawnData struct {
	Pos        Vec
	Id         uint8
	Durability int
}

var (
	toSpawn  = []spawnData{}
	toRemove []ecs.Entity
)

type Spawn struct {
	tile uint8
}

func (s *Spawn) Init() {}
func (s *Spawn) Update() {

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		x, y := cameraRes.ScreenToWorld(ebiten.CursorPosition())
		SpawnEnemy(Vec{x, y}, Vec{0.3, 0})
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := cameraRes.ScreenToWorld(ebiten.CursorPosition())
		p := tileMapRes.WorldToTile(x, y)
		tileMapRes.Set(p.X, p.Y, s.tile)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		x, y := cameraRes.ScreenToWorld(ebiten.CursorPosition())
		p := tileMapRes.WorldToTile(x, y)
		s.tile = tileMapRes.GetID(p.X, p.Y)
	}

	// Spawn item
	for _, d := range toSpawn {
		SpawnItem(d.Pos, d.Id, d.Durability)
	}

	toSpawn = toSpawn[:0]

	for _, e := range toRemove {
		world.RemoveEntity(e)
	}
	toRemove = toRemove[:0]
}
func (s *Spawn) Draw() {
}
