package kar

import (
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/kamera/v2"
)

const TimerTick time.Duration = time.Second / 60
const DeltaTime float64 = 1.0 / 60.0

var (
	ScreenW, ScreenH = 800.0, 600.0
	Screen           *ebiten.Image
	Camera           = kamera.NewCamera(0, 0, ScreenW, ScreenH)
	WorldECS         = ecs.NewWorld()
	DesktopPath      string
	GlobalDIO        = &ebiten.DrawImageOptions{}
)

type ISystem interface {
	Init()
	Update()
	Draw()
}

func init() {
	Camera.Smoothing = kamera.Lerp
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	DesktopPath = homePath + "/Desktop/"
}
