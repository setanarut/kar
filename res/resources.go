package res

import (
	"embed"
	_ "image/png"
	"kar/comp"
	"kar/engine"
	"kar/engine/util"
	"kar/items"
	"kar/types"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/setanarut/cm"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/vec"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"
)

//go:embed assets/*
var assets embed.FS

// GameSettings
var (
	WorldSize         = vec.Vec2{1024, 512}
	BlockSize float64 = 16 * 3
	ChunkSize float64 = 12
	/*
		640,360  1280,720
		854,480  1366,768
		960,540  1600,900
		1024,576 1920,1080
	*/
	ScreenSize                                 = vec.Vec2{960, 540}
	GlobalDrawOptions *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
	DesktopDir        string
	BlockCenterOffset = vec.Vec2{(BlockSize / 2), (BlockSize / 2)}.Neg()

	BlockMaxHealth = map[types.ItemID]float64{
		items.Grass:               5.0,
		items.Dirt:                5.0,
		items.Sand:                5.0,
		items.Stone:               10.0,
		items.CoalOre:             10.0,
		items.GoldOre:             10.0,
		items.IronOre:             10.0,
		items.DiamondOre:          10.0,
		items.DeepSlateStone:      15.0,
		items.DeepSlateCoalOre:    15.0,
		items.DeepSlateGoldOre:    15.0,
		items.DeepSlateIronOre:    15.0,
		items.DeepSlateDiamondOre: 5.0,
	}
)

var (
	Cam               *kamera.Camera
	SelectedItem      types.ItemID         = items.Air
	ECSWorld          donburi.World        = donburi.NewWorld()
	Space             *cm.Space            = cm.NewSpace()
	Input             *engine.InputManager = &engine.InputManager{}
	FilterBombRaycast cm.ShapeFilter       = cm.NewShapeFilter(
		0,
		BitmaskBombRaycast,
		cm.AllCategories&^BitmaskBomb)
	FilterPlayerRaycast cm.ShapeFilter = cm.NewShapeFilter(
		0,
		BitmaskPlayerRaycast,
		cm.AllCategories&^BitmaskPlayer)
)

// Donburi queries
var (
	QueryWASDcontrollable = donburi.NewQuery(filter.And(
		filter.Contains(comp.Mobile, comp.WASDTag, comp.Body),
		filter.Not(filter.Contains(comp.AI))))
	QueryDraw = donburi.NewQuery(filter.Contains(comp.DrawOptions, comp.Body))
)

// text
var (
	Font             = util.LoadGoTextFaceFromFS("assets/iosevkaem.otf", 15, assets)
	StatsTextOptions = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{Filter: ebiten.FilterNearest},
		LayoutOptions:    text.LayoutOptions{LineSpacing: Font.Size * 1.3},
	}
)

var (
	Zero  = vec.Vec2{0, 0}
	Right = vec.Vec2{1, 0}
	Left  = vec.Vec2{-1, 0}
	Up    = vec.Vec2{0, -1}
	Down  = vec.Vec2{0, 1}
)

func init() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	DesktopDir = homePath + "/Desktop/"
	StatsTextOptions.ColorScale.ScaleWithColor(colornames.White)
}
