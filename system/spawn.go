package system

import (
	"image"
	"kar/arche"
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/engine/terr"
	"kar/res"
	"kar/types"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// var spawnTick int
var playerChunkTemp image.Point

type SpawnSystem struct {
	Terr           *terr.Terrain
	spawnTimerData *types.DataTimer
}

func NewSpawnSystem() *SpawnSystem {

	return &SpawnSystem{}
}

func (sys *SpawnSystem) Init() {
	sys.spawnTimerData = &types.DataTimer{
		TimerDuration: time.Second * 2,
		Elapsed:       0,
	}
	sys.Terr = terr.NewTerrain(342, 16)
	sys.Terr.NoiseOptions.Frequency = 0.2
	sys.Terr.Generate(true, 2048)
	ResetLevel()

}

func (s *SpawnSystem) Update() {

	if player, ok := comp.PlayerTag.First(res.World); ok {
		pos := comp.Body.Get(player).Position()
		playerChunk := s.Terr.ChunkCoord(pos, 50)

		if playerChunkTemp != playerChunk {
			playerChunkTemp = playerChunk
			// Spawn Chunks
			s.Terr.SpawnChunks(playerChunk, arche.SpawnBlock)

		}
	}

	/* 	timerUpdate(s.spawnTimerData)
	   	if timerIsReady(s.spawnTimerData) {
	   		if spawnTick > -32 {
	   			spawnTick--
	   			s.Terr.SpawnChunk(image.Point{0, spawnTick}, arche.SpawnBlock)
	   		}
	   		timerReset(s.spawnTimerData)
	   	} */

	// Reset Level
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		// ResetLevel()
		// res.Camera.ZoomFactor = 0
		comp.WallTag.Each(res.World, DestroyEntryWithBody)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {

		if player, ok := comp.PlayerTag.First(res.World); ok {
			pos := comp.Body.Get(player).Position()
			chunkcoord := s.Terr.ChunkCoord(pos, 50)
			s.Terr.SpawnChunks(chunkcoord, arche.SpawnBlock)
		}

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {

		for range 10 {
			arche.SpawnDefaultBomb(engine.RandomPointInBB(res.CurrentRoom, 64))
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
