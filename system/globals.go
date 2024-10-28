package system

import (
	"kar"
	"kar/engine/mathutil"
	"kar/res"
	"kar/world"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/kamera/v2"
	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

const itemAnimFrameCount int = 200

var (
	gameWorld            *world.World
	ecsWorld             = donburi.NewWorld()
	cmSpace              = cm.NewSpace()
	camera               *kamera.Camera
	selectedSlotIndex    = 0
	desktopDir           string
	blockCenterOffset    = vec2{(kar.BlockSize / 2), (kar.BlockSize / 2)}.Neg()
	globalDIO            = &ebiten.DrawImageOptions{}
	fontSmallDrawOptions = &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{LineSpacing: res.FontSmall.Size * 1.3},
	}
)

func init() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	desktopDir = homePath + "/Desktop/"
}

var (
	drawBlockBorderEnabled bool      = true
	debugDrawingEnabled    bool      = false
	sinSpaceFrames         []float64 = mathutil.SinSpace(
		0,
		2*math.Pi,
		4,
		itemAnimFrameCount+1,
	)
	cmDrawer = ebitencm.NewDrawer()
)

var (
	justPressed  = inpututil.IsKeyJustPressed
	justReleased = inpututil.IsKeyJustReleased
	pressed      = ebiten.IsKeyPressed
)
var (
	right = vec2{1, 0}
	left  = vec2{-1, 0}
	down  = vec2{0, 1}
	up    = vec2{0, -1}
	zero  = vec2{0, 0}
)

func init() {
	cmDrawer.StrokeDisabled = true
	cmDrawer.Theme.ShapeSleeping = colornames.Green
}
