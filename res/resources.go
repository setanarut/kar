package res

import (
	"embed"
	"image"
	_ "image/png"
	"kar/comp"
	"kar/engine"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/items"
	"kar/types"

	"github.com/setanarut/cm"
	"github.com/setanarut/kamera/v2"

	"github.com/setanarut/vec"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"
	"golang.org/x/text/language"
)

//go:embed assets/*
var assets embed.FS

var (
	Zero  = vec.Vec2{0, 0}
	Right = vec.Vec2{1, 0}
	Left  = vec.Vec2{-1, 0}
	Up    = vec.Vec2{0, -1}
	Down  = vec.Vec2{0, 1}
)

var (
	NinthHD      = image.Point{640, 360}   // 640x360
	FullWideVGA  = image.Point{854, 480}   // 854x480
	QuarterHD    = image.Point{960, 540}   // 960x540
	WideSuperVGA = image.Point{1024, 576}  // 1024x576
	HD           = image.Point{1280, 720}  // 1280x720
	FullWideXGA  = image.Point{1366, 768}  // 1366x768
	HDPlus       = image.Point{1600, 900}  // 1600x900
	FullHD       = image.Point{1920, 1080} // 1920x1080
)

// GameSettings
var (
	MapSize        float64                  = 1024
	BlockSize      float64                  = 64
	ChunkSize      float64                  = 8
	ScreenSize                              = FullWideVGA
	ScreenSizeF                             = mathutil.PointToVec2(ScreenSize)
	CameraDrawOpts *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}

	BlockMaxHealth = map[types.ItemType]float64{
		items.Dirt:    5.0,
		items.Stone:   10.0,
		items.IronOre: 10.0,
	}
)

var (
	World             donburi.World = donburi.NewWorld()
	Space             *cm.Space     = cm.NewSpace()
	Cam               *kamera.Camera
	CurrentItem       types.ItemType
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
var (
	Terrain *image.Gray
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
	Font    = util.LoadGoTextFaceFromFS("assets/roboto-semi.ttf", 18, assets)
	FontBig = &text.GoTextFace{
		Source:   Font.Source,
		Size:     28,
		Language: language.English,
	}
	StatsTextOptions = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{Filter: ebiten.FilterLinear},
		LayoutOptions:    text.LayoutOptions{LineSpacing: Font.Size * 1.5},
	}
	CenterTextOptions = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{Filter: ebiten.FilterLinear},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign:   text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
			LineSpacing:    FontBig.Size * 1.2},
	}
)

func init() {
	StatsTextOptions.ColorScale.ScaleWithColor(colornames.White)
}
