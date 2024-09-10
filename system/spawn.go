package system

import (
	"image"
	"kar/arche"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/res"
	"kar/terr"

	"github.com/setanarut/vec"

	"github.com/hajimehoshi/ebiten/v2"
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
	seed := mathutil.RandRangeInt(0, 1000)
	Terr = terr.NewTerrain(seed, res.MapSize, res.ChunkSize, res.BlockSize)
	Terr.NoiseOptions.Frequency = 0.2
	Terr.Generate()
	res.Terrain = Terr.TerrainImg
	playerSpawnPosition := vec.Vec2{100, -100}

	playerChunk := Terr.WorldPosToChunkCoord(playerSpawnPosition)
	Terr.LoadedChunks = terr.GetPlayerChunks(playerChunk)
	for _, coord := range Terr.LoadedChunks {
		Terr.SpawnChunk(coord, arche.SpawnBlock)
	}
	arche.SpawnPlayer(playerSpawnPosition, 1, 0, 0)

}

func (s *SpawnSystem) Update() {

	// if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
	// 	cursor := res.Camera.ScreenToWorld(ebiten.CursorPosition())
	// 	fmt.Println(cursor)
	// 	// arche.SpawnDebug(cursor)
	// }

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
