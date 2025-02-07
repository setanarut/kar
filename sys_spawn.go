package kar

import (
	"image"
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

	// Save TileMap image to desktop
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		playerTile := image.Point{}
		if ECWorld.Alive(CurrentPlayer) {
			playerPos := MapPosition.Get(CurrentPlayer)
			playerTile = TileMapRes.WorldToTile(playerPos.X+8, playerPos.Y+8)
		}
		TileMapRes.WriteAsImage(DesktopPath+"map.png", playerTile.X, playerTile.Y)
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		x, y := CameraRes.ScreenToWorld(ebiten.CursorPosition())
		SpawnEnemy(x, y, 0, 1)
	}
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
