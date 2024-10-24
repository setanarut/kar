package system

import (
	"image"
	"kar"
	"kar/items"
	"kar/res"
	"kar/util"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

var (
	selectedSlotItemID    = items.Air
	selectedSlotIndex     = 0
	desktopDir            string
	playerFlyModeDisabled bool
	filterPlayerRaycast   = cm.ShapeFilter{
		Group:      cm.NoGroup,
		Categories: kar.PlayerRayMask,
		Mask:       cm.AllCategories &^ kar.PlayerMask &^ kar.DropItemMask}
	// fontDrawOptions = &text.DrawOptions{
	// 	LayoutOptions: text.LayoutOptions{LineSpacing: res.Font.Size * 1.3},
	// }
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
	attackSegQuery                                             cm.SegmentQueryInfo
	hitShape                                                   *cm.Shape
	placeBlockPos, hitBlockPos, attackSegEnd                   vec.Vec2
	playerPixelCoord, placeBlockPixelCoord, hitBlockPixelCoord image.Point
	hitItemID                                                  uint16
)

var (
	drawBlockBorderEnabled bool      = true
	debugDrawingEnabled    bool      = false
	itemAnimFrameCount     int       = 100
	sinSpace               []float64 = util.SinSpace(0, 2*math.Pi, 2, itemAnimFrameCount+1)

	// occlusionCulling bool
)

var (
	justPressed  = inpututil.IsKeyJustPressed
	justReleased = inpututil.IsKeyJustReleased
	pressed      = ebiten.IsKeyPressed
)

var (
	right = vec.Vec2{1, 0}
	left  = vec.Vec2{-1, 0}
	down  = vec.Vec2{0, 1}
	up    = vec.Vec2{0, -1}
	zero  = vec.Vec2{0, 0}
)
