package system

import (
	"fmt"
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

func (sys *SpawnSystem) Init() {
	sys.Terr = terr.NewTerrain(342, 1024, 16, 100)
	sys.Terr.NoiseOptions.Frequency = 0.2
	sys.Terr.Generate()
	ResetLevel()

}

func (s *SpawnSystem) Update() {

	if player, ok := comp.PlayerTag.First(res.World); ok {
		pos := comp.Body.Get(player).Position()
		playerChunk := s.Terr.ChunkCoord(pos)

		if playerChunkTemp != playerChunk {
			playerChunkTemp = playerChunk
			// Spawn Chunks
			fmt.Println(s.Terr.LoadedChunks)
			s.Terr.SpawnChunks(playerChunk, arche.SpawnBlock)

		}
	}

}

func (s *SpawnSystem) Draw() {
}

func ResetLevel() {

	res.Camera.Reset()
	playerSpawnPosition := cm.Vec2{0, 0}

	if player, ok := comp.PlayerTag.First(res.World); ok {
		DestroyEntryWithBody(player)

		arche.SpawnDefaultPlayer(playerSpawnPosition)
	} else {
		arche.SpawnDefaultPlayer(playerSpawnPosition)
	}

}
