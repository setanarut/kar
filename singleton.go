package kar

import (
	"image/color"
	"kar/items"
	"time"

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

const TimerTick time.Duration = time.Second / 60
const DeltaTime float64 = 1.0 / 60.0

var (
	// DesktopPath      string
	ScreenW, ScreenH = 854.0, 480.0
	Screen           *ebiten.Image
	Camera           = kamera.NewCamera(400, 450, ScreenW, ScreenH)
	WorldECS         = ecs.NewWorld()
	GlobalDIO        = &colorm.DrawImageOptions{}
	GlobalColorM     = colorm.ColorM{}
	BackgroundColor  = rgb(38, 0, 121)
	ItemScale        = 1.0
	PlayerScale      = 2.0
	// Debug
	DrawDebugHitboxesEnabled = false
	DrawDebugTextEnabled     = false
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
