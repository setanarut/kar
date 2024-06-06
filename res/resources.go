package res

import (
	"embed"
	"image"
	"image/color"
	_ "image/png"
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
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

var (
	World donburi.World = donburi.NewWorld()
	Space *cm.Space     = cm.NewSpace()

	ScreenRect, CurrentRoom cm.BB
	Camera                  *engine.Camera

	CurrentTool         types.ItemType
	Rooms               []cm.BB              = make([]cm.BB, 0)
	Input               *engine.InputManager = &engine.InputManager{}
	FilterBombRaycast   cm.ShapeFilter       = cm.NewShapeFilter(0, types.BitmaskBombRaycast, cm.AllCategories&^types.BitmaskBomb)
	FilterPlayerRaycast cm.ShapeFilter       = cm.NewShapeFilter(0, types.BitmaskPlayerRaycast, cm.AllCategories&^types.BitmaskPlayer)
	DamageGradient, _                        = colorgrad.NewGradient().
				HtmlColors("rgb(255, 0, 0)", "rgb(255, 225, 0)", "rgb(111, 111, 111)").
				Domain(0, 1).
				Mode(colorgrad.BlendOklab).
				Interpolation(colorgrad.InterpolationBasis).
				Build()
	QueryWASDcontrollable = donburi.NewQuery(filter.And(
		filter.Contains(comp.Mobile, comp.WASDTag, comp.Body),
		filter.Not(filter.Contains(comp.AI))))
	QueryAI = donburi.NewQuery(filter.And(
		filter.Contains(comp.Mobile, comp.AI, comp.Body),
		filter.Not(filter.Contains(comp.WASDTag))))
)

var (
	Screen  *ebiten.Image
	Terrain *image.Gray
	Wall    = ebiten.NewImage(50, 50)
	Player  = engine.LoadImage("assets/player.png", assets)
)

var (
	Futura = engine.LoadTextFace("assets/futura.ttf", 20, assets)
	// Miglia    = engine.LoadTextFace("assets/miglia_fur.ttf", 20, assets)
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
	Wall.Fill(color.Gray{128})
	StatsTextOptions.ColorScale.ScaleWithColor(colornames.White)
}
