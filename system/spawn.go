package system

import (
	"image"
	"kar/arche"
	"kar/comp"
	"kar/engine/cm"
	"kar/engine/terr"
	"kar/res"
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
	s.Terr = terr.NewTerrain(342, 1024, 16, 50)
	s.Terr.NoiseOptions.Frequency = 0.2
	s.Terr.Generate()
	res.Terrain = s.Terr.TerrainImg
	ResetLevel()

	if player, ok := comp.PlayerTag.First(res.World); ok {
		pos := comp.Body.Get(player).Position()
		playerChunk := s.Terr.ChunkCoord(pos)
		s.Terr.LoadedChunks = terr.GetPlayerChunks(playerChunk)

		for _, coord := range s.Terr.LoadedChunks {
			s.Terr.SpawnChunk(coord, arche.SpawnBlock)
		}
	}
}

func (s *SpawnSystem) Update() {

	if player, ok := comp.PlayerTag.First(res.World); ok {
		pos := comp.Body.Get(player).Position()
		playerChunk := s.Terr.ChunkCoord(pos)

		if playerChunkTemp != playerChunk {
			playerChunkTemp = playerChunk
			// Spawn/Destroy Chunks
			s.Terr.UpdateChunks(playerChunk, arche.SpawnBlock)

		}
	}

}

func (s *SpawnSystem) Draw() {
}

func ResetLevel() {

	res.Camera.Reset()
	playerSpawnPosition := cm.Vec2{512 * 50, 512 * 50}

	if player, ok := comp.PlayerTag.First(res.World); ok {
		DestroyEntryWithBody(player)
		e := arche.SpawnDefaultPlayer(playerSpawnPosition)
		r := comp.Render.Get(e)
		RenderScale = r.DrawScale
	} else {
		arche.SpawnDefaultPlayer(playerSpawnPosition)
	}

}
