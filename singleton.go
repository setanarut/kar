package kar

import (
	"image"
	"image/color"
	"kar/items"

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
	// DesktopPath      string
	ScreenW, ScreenH   = 860., 480.0
	Screen             *ebiten.Image
	Camera                 = kamera.NewCamera(0, 0, ScreenW, ScreenH)
	WorldECS               = ecs.NewWorld()
	GlobalDIO              = &colorm.DrawImageOptions{}
	GlobalColorM           = colorm.ColorM{}
	ItemScale              = 2.0
	PlayerScale            = 2.0
	ItemCollisionDelay     = 10
	RaycastDist        int = 4                   // block unit
	RenderArea             = image.Point{23, 13} // cam w/h blocks
	// Debug
	DrawDebugHitboxesEnabled = false
	DrawDebugTextEnabled     = false
	BackgroundColor          = rgb(36, 36, 39)
)
var ItemColorMap = map[uint16]color.RGBA{
	items.Air:        rgb(1, 1, 1),
	items.GrassBlock: rgb(0, 186, 53),
	items.Dirt:       rgb(133, 75, 54),
	items.Sand:       rgb(199, 193, 158),
	items.Stone:      rgb(139, 139, 139),
	items.CoalOre:    rgb(0, 0, 0),
	items.GoldOre:    rgb(255, 221, 0),
	items.IronOre:    rgb(171, 162, 147),
	items.DiamondOre: rgb(0, 247, 255),
}

func init() {
	Camera.SmoothType = kamera.None
	// GlobalColorM.ChangeHSV(1, 0, 1)

	// homePath, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// DesktopPath = homePath + "/Desktop/"
}

func rgb(r, g, b uint8) color.RGBA {
	return color.RGBA{r, g, b, 255}
}
