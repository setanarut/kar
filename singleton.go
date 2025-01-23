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
	ScreenW, ScreenH         float64     = 500.0, 340.0
	ItemCollisionDelay       int         = 10
	RaycastDist              int         = 4 // block unit
	DrawDebugHitboxesEnabled bool        = false
	DrawDebugTextEnabled     bool        = false
	PlayerBestToolDamage                 = 5.0
	PlayerDefaultDamage                  = 1.0
	RenderArea               image.Point = image.Point{
		(int(ScreenW) / 20) + 3,
		(int(ScreenH) / 20) + 3,
	} // cam w/h blocks
	BackgroundColor color.RGBA = color.RGBA{36, 36, 39, 255}

	Screen    *ebiten.Image
	Camera    *kamera.Camera
	WorldECS  = ecs.NewWorld()
	ColorMDIO = &colorm.DrawImageOptions{}
	ColorM    = colorm.ColorM{}
	// Debug
)

func init() {
	// GlobalColorM.ChangeHSV(1, 0, 1)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	DesktopPath = homePath + "/Desktop/"
}
