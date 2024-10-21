package system

import (
	"kar"
	"kar/arche"
	"kar/comp"
	"kar/items"
	"kar/world"

	"github.com/setanarut/vec"

	eb "github.com/hajimehoshi/ebiten/v2"
	iu "github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Spawn struct{}

func (s *Spawn) Init() {
	gameWorld = world.NewWorld(
		kar.WorldSize.X,
		kar.WorldSize.Y,
		kar.ChunkSize,
		kar.BlockSize,
	)
	findSpawnPosition()
	gameWorld.LoadChunks(cmSpace, ecsWorld, playerSpawnPos)
	playerEntry = arche.SpawnPlayer(cmSpace, ecsWorld, playerSpawnPos, 1, 0, 0)
}
func (s *Spawn) Update() {
	if iu.IsMouseButtonJustPressed(eb.MouseButton0) {
		if eb.IsKeyPressed(eb.KeyC) {
			x, y := camera.ScreenToWorld(eb.CursorPosition())
			arche.SpawnDebugBox(cmSpace, ecsWorld, vec.Vec2{x, y})
		}
	}

	if playerEntry.Valid() {
		inventory = comp.Inventory.Get(playerEntry)
		playerPos = comp.Body.Get(playerEntry).Position()
		gameWorld.UpdateChunks(cmSpace, ecsWorld, playerPos)
	}
}

func (s *Spawn) Draw() {
}

func findSpawnPosition() {
	for y := range 300 {
		posUp := gameWorld.Image.Gray16At(10, y).Y
		posDown := gameWorld.Image.Gray16At(10, y+1).Y
		if posDown != items.Air && posUp == items.Air {
			playerSpawnPos = world.PixelCoordToWorldPos(10, y)
			break
		}
	}
}
