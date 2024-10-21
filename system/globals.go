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
	"github.com/setanarut/vec"
	"github.com/yohamta/donburi"
)

var (
	playerEntry        *donburi.Entry
	inventory          *types.Inventory
	Camera             *kamera.Camera
	ecsWorld           = donburi.NewWorld()
	space              = cm.NewSpace()
	selectedSlotItemID = items.Air
	selectedSlotIndex  = 0
	desktopDir         string
	blockCenterOffset  = vec.Vec2{(kar.BlockSize / 2), (kar.BlockSize / 2)}.Neg()
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
	attackSegQuery                                          cm.SegmentQueryInfo
	hitShape                                                *cm.Shape
	playerPos, placeBlockPos, currentBlockPos, attackSegEnd vec.Vec2
	playerPosMap, placeBlockPosMap, currentBlockPosMap      image.Point
	hitItemID                                               uint16
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
	cmDrawer     = ebitencm.NewDrawer()
	cameraBounds cm.BB
	// occlusionCulling bool
)

var (
	justPressed  = inpututil.IsKeyJustPressed
	justReleased = inpututil.IsKeyJustReleased
	keyPressed   = ebiten.IsKeyPressed
)
var playerChunk image.Point
var mainWorld *world.World

var wasdLast vec.Vec2
var wasd vec.Vec2

var (
	right = vec.Vec2{1, 0}
	left  = vec.Vec2{-1, 0}
	down  = vec.Vec2{0, 1}
	up    = vec.Vec2{0, -1}
	noDir = vec.Vec2{0, 0}
)
