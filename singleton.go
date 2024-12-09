package kar

import (
	"image/color"
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
	ScreenW, ScreenH = 854.0, 480.0
	Screen           *ebiten.Image
	Camera           = kamera.NewCamera(0, 0, ScreenW, ScreenH)
	WorldECS         = ecs.NewWorld()
	// DesktopPath      string
	GlobalDIO              = &ebiten.DrawImageOptions{}
	BackgroundColor        = rgb(30, 47, 91)
	TileColor              = rgb(181, 86, 35)
	TargetTileBorderColor  = rgb(30, 255, 0)
	BlockPlacementDistance = 4
	DrawPlayerDebugHitBox  = false
)

func init() {
	Camera.Smoothing = kamera.SmoothDamp
	Camera.SmoothingOptions.SmoothDampTimeX = 0.3
	Camera.SmoothingOptions.SmoothDampTimeY = 1
	Camera.SmoothingOptions.SmoothDampMaxSpeedX = 1000
	Camera.SmoothingOptions.SmoothDampMaxSpeedY = 1000

	// homePath, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// DesktopPath = homePath + "/Desktop/"
}

func rgb(r, g, b uint8) color.RGBA {
	return color.RGBA{r, g, b, 255}
}
