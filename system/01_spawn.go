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
	mainWorld = world.NewWorld(
		kar.WorldSize.X,
		kar.WorldSize.Y,
		kar.ChunkSize,
		kar.BlockSize,
	)
	var playerSpawnPos vec.Vec2

	for y := range 300 {
		posUp := mainWorld.Image.Gray16At(10, y).Y
		posDown := mainWorld.Image.Gray16At(10, y+1).Y
		if posDown != items.Air && posUp == items.Air {
			playerSpawnPos = world.PixelToWorldSpace(10, y)
			break
		}
	}
	playerChunkCoord := mainWorld.WorldPosToChunkCoord(playerSpawnPos)
	mainWorld.LoadedChunks = world.GetPlayerNeighborChunks(playerChunkCoord)
	for _, coord := range mainWorld.LoadedChunks {
		mainWorld.SpawnChunk(space, ecsWorld, coord)
	}
	playerEntry = arche.SpawnPlayer(space, ecsWorld, playerSpawnPos, 1, 0, 0)
}

func (s *Spawn) Update() {
	if iu.IsMouseButtonJustPressed(eb.MouseButton0) {

		if eb.IsKeyPressed(eb.KeyC) {
			x, y := Camera.ScreenToWorld(eb.CursorPosition())
			arche.SpawnDebugBox(space, ecsWorld, vec.Vec2{x, y})
		}
	}

	if player, ok := comp.TagPlayer.First(ecsWorld); ok {
		playerEntry = player
		inventory = comp.Inventory.Get(player)
		playerPos = comp.Body.Get(player).Position()
		playerChunkTemp := mainWorld.WorldPosToChunkCoord(playerPos)
		if playerChunk != playerChunkTemp {
			playerChunk = playerChunkTemp
			mainWorld.UpdateChunks(space, ecsWorld, playerChunkTemp)
		}
	}
}

func (s *Spawn) Draw() {
}
