package system

import (
	"image"
	"kar/arche"
	"kar/comp"
	"kar/engine/vec"
	"kar/items"
	"kar/res"
	"kar/terr"
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

	Terr = terr.NewTerrain(342, res.MapSize, res.ChunkSize, res.BlockSize)
	Terr.NoiseOptions.Frequency = 0.2
	Terr.Generate()
	res.Terrain = Terr.TerrainImg

	playerSpawnPosition := vec.Vec2{(res.MapSize * 0.5) * res.BlockSize, Terr.BlockSize * 2}

	playerChunk := Terr.WorldPosToChunkCoord(playerSpawnPosition)
	Terr.LoadedChunks = terr.GetPlayerChunks(playerChunk)
	for _, coord := range Terr.LoadedChunks {
		Terr.SpawnChunk(coord, arche.SpawnItem)
	}
	arche.SpawnDefaultPlayer(playerSpawnPosition)
	arche.SpawnItem(vec.Vec2{}, playerChunk, items.RawIron)
}

func (s *SpawnSystem) Update() {

	if player, ok := comp.PlayerTag.First(res.World); ok {
		pos := comp.Body.Get(player).Position()
		playerChunk := Terr.WorldPosToChunkCoord(pos)

		if playerChunkTemp != playerChunk {
			playerChunkTemp = playerChunk
			// Spawn/Destroy Chunks
			Terr.UpdateChunks(playerChunk, arche.SpawnItem)

		}
	}

}

func (s *SpawnSystem) Draw() {
}
