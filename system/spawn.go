package system

import (
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
	Terr terr.Terrain
}

func NewSpawnSystem() *SpawnSystem {

	return &SpawnSystem{}
}

func (sys *SpawnSystem) Init() {
	sys.Terr = *terr.NewTerrain(12, 128, 8)
	ResetLevel(&sys.Terr)

}

func (sys *SpawnSystem) Update() {

	// Reset Level
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		ResetLevel(&sys.Terr)
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

func ResetLevel(terra *terr.Terrain) {

	res.Camera.Reset()

	if p, ok := comp.PlayerTag.First(res.World); ok {
		destroyEntryWithBody(p)
		arche.SpawnDefaultPlayer(cm.Vec2{0, 100})
	} else {
		arche.SpawnDefaultPlayer(cm.Vec2{0, 100})
	}

	comp.EnemyTag.Each(res.World, func(e *donburi.Entry) {
		destroyEntryWithBody(e)
	})

	comp.BombTag.Each(res.World, func(e *donburi.Entry) {
		destroyEntryWithBody(e)
	})
	terra.NoiseOptions.Frequency = 0.1

	for y := 0; y > -terra.MapSize; y-- {
		for x := 0; x < terra.MapSize; x++ {
			if terra.GetBlockValue(x, y) > 0.5 {
				pos := cm.Vec2{float64(x) * 64, float64(y) * 64}
				arche.SpawnWall(pos.Round(), 64, 64)
			}
		}

	}
}
