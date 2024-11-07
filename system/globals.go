package system

import (
	"kar"
	"kar/engine/mathutil"
	"kar/world"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/kamera/v2"
	"golang.org/x/image/colornames"
)

const itemAnimFrameCount int = 200

var (
	gameWorld         *world.World
	Camera            *kamera.Camera
	selectedSlotIndex = 0
	desktopDir        string
	blockCenterOffset = vec2{kar.BlockSize / 2, kar.BlockSize / 2}.Neg()
	globalDIO         = &ebiten.DrawImageOptions{}
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

func init() {
	cmDrawer.StrokeDisabled = false
	cmDrawer.Theme.ShapeSleeping = colornames.Green
}
