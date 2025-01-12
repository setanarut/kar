package kar

import (
	"image"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/kamera/v2"
)

type ISystem interface {
	Init()
	Update()
	Draw()
}

var (
	DesktopPath              string
	WindowScale              float64     = 2.0
	ScreenW, ScreenH         float64     = 400.0, 240.0
	ItemCollisionDelay       int         = 10
	RaycastDist              int         = 4 // block unit
	DrawDebugHitboxesEnabled bool        = false
	DrawDebugTextEnabled     bool        = false
	RenderArea               image.Point = image.Point{23, 13} // cam w/h blocks
	BackgroundColor          color.RGBA  = color.RGBA{36, 36, 39, 255}

	Screen          *ebiten.Image
	Camera          = kamera.NewCamera(0, 0, ScreenW, ScreenH)
	WorldECS        = ecs.NewWorld()
	GlobalColorMDIO = &colorm.DrawImageOptions{}
	GlobalColorM    = colorm.ColorM{}
	// Debug
)

func init() {
	// GlobalColorM.ChangeHSV(1, 0, 1)
	// Camera.SmoothType = kamera.SmoothDamp
	Camera.SmoothType = kamera.Lerp
	Camera.SmoothOptions.LerpSpeedX = 0.5
	Camera.SmoothOptions.LerpSpeedY = 0.05
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	DesktopPath = homePath + "/Desktop/"
}
