package res

import (
	"embed"
	"image"
	_ "image/png"
	"kar/comp"
	"kar/engine"
	"kar/engine/displayres"
	"kar/engine/io"
	"kar/engine/mathutil"
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

// GameSettings
var (
	MapSize    float64 = 1024
	BlockSize  float64 = 64.0
	ChunkSize  float64 = 8
	ScreenSize         = displayres.FullWideVGA

	ScreenSizeF = mathutil.FromPoint(ScreenSize)

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
		cm.ALL_CATEGORIES&^BitmaskBomb)
	FilterPlayerRaycast cm.ShapeFilter = cm.NewShapeFilter(
		0,
		BitmaskPlayerRaycast,
		cm.ALL_CATEGORIES&^BitmaskPlayer)
)
var (
	Screen  *ebiten.Image
	Terrain *image.Gray
)

// Donburi queries
var (
	QueryWASDcontrollable = donburi.NewQuery(filter.And(
		filter.Contains(comp.Mobile, comp.WASDTag, comp.Body),
		filter.Not(filter.Contains(comp.AI))))
	QueryAI = donburi.NewQuery(filter.And(
		filter.Contains(comp.Mobile, comp.AI, comp.Body),
		filter.Not(filter.Contains(comp.WASDTag))))

	QueryDraw = donburi.NewQuery(filter.Contains(comp.DrawOptions, comp.Body))
)

// text
var (
	Futura    = io.LoadGoTextFaceFromFS("assets/roboto-semi.ttf", 18, assets)
	FuturaBig = &text.GoTextFace{
		Source:   Futura.Source,
		Size:     28,
		Language: language.English,
	}
	StatsTextOptions = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{Filter: ebiten.FilterLinear},
		LayoutOptions:    text.LayoutOptions{LineSpacing: FuturaBig.Size * 1.2},
	}
	CenterTextOptions = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{Filter: ebiten.FilterLinear},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign:   text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
			LineSpacing:    FuturaBig.Size * 1.2},
	}
)

func init() {
	StatsTextOptions.ColorScale.ScaleWithColor(colornames.White)
}
