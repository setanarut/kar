package system

import (
	"image"
	"kar/arche"
	"kar/comp"
	"kar/res"
	"kar/world"

	"github.com/setanarut/vec"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// var spawnTick int
var PlayerChunk image.Point
var MainWorld *world.World

type SpawnSystem struct {
}

func NewSpawnSystem() *SpawnSystem {

	return &SpawnSystem{}
}

func (s *SpawnSystem) Init() {
	MainWorld = world.NewWorld(res.WorldSize.X, res.WorldSize.Y, res.ChunkSize, res.BlockSize)
	var playerSpawnPos vec.Vec2

	for y := range 300 {
		if MainWorld.Image.Gray16At(256, y).Y != 0 {
			playerSpawnPos = vec.Vec2{255. * res.BlockSize, float64(y-1) * res.BlockSize}
			break
		}
	}
	playerChunk := MainWorld.WorldPosToChunkCoord(playerSpawnPos)
	MainWorld.LoadedChunks = world.GetPlayerNeighborChunks(playerChunk)
	for _, coord := range MainWorld.LoadedChunks {
		MainWorld.SpawnChunk(coord, arche.SpawnBlock)
	}
	arche.SpawnPlayer(playerSpawnPos, 1, 0, 0)
}

func (s *SpawnSystem) Update() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {

		if ebiten.IsKeyPressed(ebiten.KeyC) {
			x, y := res.Cam.ScreenToWorld(ebiten.CursorPosition())
			arche.SpawnDebugBox(vec.Vec2{x, y})
		}
	}

	if player, ok := comp.PlayerTag.First(res.ECSWorld); ok {
		pos := comp.Body.Get(player).Position()
		playerChunk := MainWorld.WorldPosToChunkCoord(pos)

		if PlayerChunk != playerChunk {
			PlayerChunk = playerChunk
			// Spawn/Destroy Chunks
			MainWorld.UpdateChunks(playerChunk, arche.SpawnBlock)

		}
	}

}

func (s *SpawnSystem) Draw(screen *ebiten.Image) {
}
