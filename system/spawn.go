package system

import (
	"image"
	"kar/arche"
	"kar/comp"
	"kar/engine/vec"
	"kar/res"
	"kar/terr"
)

// var spawnTick int
var playerChunkTemp image.Point

type SpawnSystem struct {
	Terr *terr.Terrain
}

func NewSpawnSystem() *SpawnSystem {

	return &SpawnSystem{}
}

func (s *SpawnSystem) Init() {

	s.Terr = terr.NewTerrain(342, 1000, 13, res.BlockSize)
	s.Terr.NoiseOptions.Frequency = 0.2
	s.Terr.Generate()
	res.Terrain = s.Terr.TerrainImg

	playerSpawnPosition := vec.Vec2{
		(float64(s.Terr.MapSize) / 2) * s.Terr.BlockSize,
		s.Terr.BlockSize * 2}

	playerChunk := s.Terr.WorldPosToChunkCoord(playerSpawnPosition)
	s.Terr.LoadedChunks = terr.GetPlayerChunks(playerChunk)
	for _, coord := range s.Terr.LoadedChunks {
		s.Terr.SpawnChunk(coord, arche.SpawnBlock)
	}
	arche.SpawnDefaultPlayer(playerSpawnPosition)
}

func (s *SpawnSystem) Update() {

	if player, ok := comp.PlayerTag.First(res.World); ok {
		pos := comp.Body.Get(player).Position()
		playerChunk := s.Terr.WorldPosToChunkCoord(pos)

		if playerChunkTemp != playerChunk {
			playerChunkTemp = playerChunk
			// Spawn/Destroy Chunks
			s.Terr.UpdateChunks(playerChunk, arche.SpawnBlock)

		}
	}

}

func (s *SpawnSystem) Draw() {
}
