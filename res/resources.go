package res

import (
	"embed"
	"image"
	_ "image/png"
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/engine/displayres"
	"kar/engine/io"
	"kar/engine/util"
	"kar/engine/vec"
	"kar/items"
	"kar/types"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/mazznoer/colorgrad"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"
	"golang.org/x/text/language"
)

//go:embed assets/*
var assets embed.FS

// GameSettings
var (
	MapSize    float64 = 1024
	BlockSize  float64 = 64.0
	ChunkSize  float64 = 8
	ScreenSize         = displayres.QuarterHD

	ScreenSizeF = vec.Vec2{float64(ScreenSize.X), float64(ScreenSize.Y)}

	BlockMaxHealth = map[types.ItemType]float64{
		items.Dirt:  5.0,
		items.Stone: 10.0,
	}
)

var (
	AtlasPlayer = io.LoadEbitenImageFromFS(assets, "assets/player.png")
	AtlasBlock  = io.LoadEbitenImageFromFS(assets, "assets/blocks.png")
	BlockFrames = make(map[types.ItemType][]*ebiten.Image)
)

var (
	World               donburi.World = donburi.NewWorld()
	Space               *cm.Space     = cm.NewSpace()
	Camera              *engine.Camera
	CurrentItem         types.ItemType
	Input               *engine.InputManager = &engine.InputManager{}
	FilterBombRaycast   cm.ShapeFilter       = cm.NewShapeFilter(0, BitmaskBombRaycast, cm.AllCategories&^BitmaskBomb)
	FilterPlayerRaycast cm.ShapeFilter       = cm.NewShapeFilter(0, BitmaskPlayerRaycast, cm.AllCategories&^BitmaskPlayer)
	DamageGradient, _                        = colorgrad.NewGradient().
				HtmlColors("rgb(255, 0, 0)", "rgb(255, 225, 0)", "rgb(111, 111, 111)").
				Domain(0, 1).
				Mode(colorgrad.BlendOklab).
				Interpolation(colorgrad.InterpolationBasis).
				Build()
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
)

// text
var (
	Futura    = io.LoadGoTextFaceFromFS("assets/futura.ttf", 18, assets)
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
	BlockFrames[items.Stone] = util.SubImages(AtlasBlock, 16, 0, 16, 16, 9, true)
	BlockFrames[items.Dirt] = util.SubImages(AtlasBlock, 0, 0, 16, 16, 9, true)
	StatsTextOptions.ColorScale.ScaleWithColor(colornames.White)
}
