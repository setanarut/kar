package system

import (
	"image"
	"kar/arche"
	"kar/comp"
	"kar/items"
	"kar/res"
	"kar/world"

	"github.com/setanarut/vec"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SpawnSystem struct{}

var playerChunk image.Point
var MainWorld *world.World

func (s *SpawnSystem) Init() {
	MainWorld = world.NewWorld(
		res.WorldSize.X,
		res.WorldSize.Y,
		res.ChunkSize,
		res.BlockSize,
	)
	var playerSpawnPos vec.Vec2

	for y := range 300 {
		posUp := MainWorld.Image.Gray16At(10, y).Y
		posDown := MainWorld.Image.Gray16At(10, y+1).Y
		if posDown != items.Air && posUp == items.Air {
			playerSpawnPos = world.PixelToWorldSpace(10, y)
			break
		}
	}
	playerChunkCoord := MainWorld.WorldPosToChunkCoord(playerSpawnPos)
	MainWorld.LoadedChunks = world.GetPlayerNeighborChunks(playerChunkCoord)
	for _, coord := range MainWorld.LoadedChunks {
		MainWorld.SpawnChunk(coord)
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

	if player, ok := comp.TagPlayer.First(res.ECSWorld); ok {
		pos := comp.Body.Get(player).Position()
		playerChunkTemp := MainWorld.WorldPosToChunkCoord(pos)
		if playerChunk != playerChunkTemp {
			playerChunk = playerChunkTemp
			MainWorld.UpdateChunks(playerChunkTemp)
		}
	}

}

func (s *SpawnSystem) Draw(screen *ebiten.Image) {
}
