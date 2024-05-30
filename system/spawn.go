package system

import (
	"kar/arche"
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

type SpawnSystem struct {
}

func NewSpawnSystem() *SpawnSystem {

	return &SpawnSystem{}
}

func (sys *SpawnSystem) Init() {
	ResetLevel()

}

func (sys *SpawnSystem) Update() {

	// Reset Level
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		ResetLevel()
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

func ResetLevel() {

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

	for y := 0; y > -1024; y -= 64 {
		for x := 0; x < 1024; x += 64 {
			p := cm.Vec2{float64(x), float64(y)}
			if engine.RandRangeInt(0, 1) == 1 {
				arche.SpawnWall(p.Round(), 64, 64)

			}
		}

	}
}
