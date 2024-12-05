package kar

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/kamera/v2"
)

type ISystem interface {
	Init()
	Update()
	Draw()
}

const TimerTick time.Duration = time.Second / 60
const DeltaTime float64 = 1.0 / 60.0

var (
	ScreenW, ScreenH = 800.0, 600.0
	Screen           *ebiten.Image
	Camera           = kamera.NewCamera(0, 0, ScreenW, ScreenH)
	WorldECS         = ecs.NewWorld()
	// DesktopPath      string
	GlobalDIO = &ebiten.DrawImageOptions{}
)

func init() {
	Camera.Smoothing = kamera.SmoothDamp
	Camera.SmoothingOptions.SmoothDampTime = 0.2
	Camera.SmoothingOptions.SmoothDampMaxSpeed = 600
	// homePath, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// DesktopPath = homePath + "/Desktop/"
}
