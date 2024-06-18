package system

import (
	"fmt"
	"image"
	"image/color"
	"kar/arche"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/vec"
	"kar/items"
	"kar/res"
	"kar/terr"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// var spawnTick int
var playerChunkTemp image.Point
var Terr *terr.Terrain

type SpawnSystem struct {
}

func NewSpawnSystem() *SpawnSystem {

	return &SpawnSystem{}
}

func (s *SpawnSystem) Init() {

	Terr = terr.NewTerrain(int64(mathutil.RandRangeInt(0, 100)), res.MapSize, res.ChunkSize, res.BlockSize)
	Terr.NoiseOptions.Frequency = 0.2
	Terr.Generate()
	res.Terrain = Terr.TerrainImg
	res.Terrain.SetGray(0, 0, color.Gray{uint8(items.IronOre)})
	playerSpawnPosition := vec.Vec2{100, -100}

	playerChunk := Terr.WorldPosToChunkCoord(playerSpawnPosition)
	Terr.LoadedChunks = terr.GetPlayerChunks(playerChunk)
	for _, coord := range Terr.LoadedChunks {
		Terr.SpawnChunk(coord, arche.SpawnBlock)
	}
	arche.SpawnDefaultPlayer(playerSpawnPosition)

}

func (s *SpawnSystem) Update() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		cursor := res.Camera.ScreenToWorld(ebiten.CursorPosition())
		fmt.Println(cursor)
		arche.SpawnRawIron(cursor)
	}

	if player, ok := comp.PlayerTag.First(res.World); ok {
		pos := comp.Body.Get(player).Position()
		playerChunk := Terr.WorldPosToChunkCoord(pos)

		if playerChunkTemp != playerChunk {
			playerChunkTemp = playerChunk
			// Spawn/Destroy Chunks
			Terr.UpdateChunks(playerChunk, arche.SpawnBlock)

		}
	}

}

func (s *SpawnSystem) Draw(screen *ebiten.Image) {
}
