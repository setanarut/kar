package kar

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const TimerTick time.Duration = time.Second / 60
const DeltaTime float64 = 1.0 / 60.0

var (
	Screen *ebiten.Image
)
