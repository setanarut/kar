package kar

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/cm"
)

const TimerTick time.Duration = time.Second / 60
const DeltaTime float64 = 1.0 / 60.0

var (
	Screen   *ebiten.Image
	WorldECS = ecs.NewWorld()
	Space    = cm.NewSpace()
)

type ISystem interface {
	Init()
	Update()
	Draw()
}
