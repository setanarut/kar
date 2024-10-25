package system

import (
	"image"
	"kar"
	"kar/engine/mathutil"
	"kar/items"
	"kar/res"
	"kar/types"
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

var (
	gameWorld          *world.World
	ecsWorld           = donburi.NewWorld()
	cmSpace            = cm.NewSpace()
	playerEntry        *donburi.Entry
	playerVel          vec2
	playerSpawnPos     vec2
	playerBody         *cm.Body
	inventory          *types.Inventory
	camera             *kamera.Camera
	selectedSlotItemID = items.Air
	selectedSlotIndex  = 0
	desktopDir         string
	blockCenterOffset  = vec2{(kar.BlockSize / 2), (kar.BlockSize / 2)}.Neg()
	globalDIO          = &ebiten.DrawImageOptions{}

	filterPlayerRaycast = cm.ShapeFilter{
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
	playerPos, placeBlockPos, hitBlockPos, attackSegEnd        vec2
	playerPixelCoord, placeBlockPixelCoord, hitBlockPixelCoord image.Point
	hitItemID                                                  uint16
)
var (
	attacking, digDown, digUp, facingDown, facingLeft, facingRight bool
	facingUp, idle, isGround, noWASD, walking, walkLeft, walkRight bool
)

var (
	drawBlockBorderEnabled bool      = true
	debugDrawingEnabled    bool      = false
	itemAnimFrameCount     int       = 100
	sinSpace               []float64 = mathutil.SinSpace(
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

var wasdLast vec2
var wasd vec2

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
