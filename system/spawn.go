package system

import (
	"image"
	"kar/arche"
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/engine/terr"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

type SpawnSystem struct {
	Terr *terr.Terrain
}

func NewSpawnSystem() *SpawnSystem {

	return &SpawnSystem{}
}

func (sys *SpawnSystem) Init() {
	sys.Terr = terr.NewTerrain(342)
	sys.Terr.NoiseOptions.Frequency = 0.2
	sys.Terr.Generate()
	ResetLevel(sys.Terr)

}

func (sys *SpawnSystem) Update() {

	// Reset Level
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		ResetLevel(sys.Terr)
	}

	// worldPos := res.Camera.ScreenToWorld(ebiten.CursorPosition())
	// cursor := engine.InvPosVectY(worldPos, res.CurrentRoom.T)

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {

		for range 4 {
			// arche.SpawnDefaultEnemy(engine.RandomPointInBB(res.CurrentRoom, 64))
			arche.SpawnDefaultMob(engine.RandomPointInBB(res.CurrentRoom, 64))
		}

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {

		for range 10 {
			arche.SpawnDefaultBomb(engine.RandomPointInBB(res.CurrentRoom, 64))
		}

	}

}

func (sys *SpawnSystem) Draw() {
}

func ResetLevel(tr *terr.Terrain) {

	res.Camera.Reset()
	playerSpawnPosition := cm.Vec2{-2 * 50, -2 * 50}

	if player, ok := comp.PlayerTag.First(res.World); ok {
		destroyEntryWithBody(player)

		arche.SpawnDefaultPlayer(playerSpawnPosition)
	} else {
		arche.SpawnDefaultPlayer(playerSpawnPosition)
	}

	comp.EnemyTag.Each(res.World, func(e *donburi.Entry) {
		destroyEntryWithBody(e)
	})

	comp.BombTag.Each(res.World, func(e *donburi.Entry) {
		destroyEntryWithBody(e)
	})

	chunkCoord := tr.ChunkCoord(playerSpawnPosition).Div(50)
	tr.SpawnChunk(chunkCoord, arche.SpawnBlock)
	tr.SpawnChunk(chunkCoord.Add(image.Point{0, 1}), arche.SpawnBlock)
	tr.SpawnChunk(chunkCoord.Add(image.Point{0, 2}), arche.SpawnBlock)
	tr.SpawnChunk(chunkCoord.Add(image.Point{1, 0}), arche.SpawnBlock)
	tr.SpawnChunk(chunkCoord.Add(image.Point{2, 0}), arche.SpawnBlock)
}
