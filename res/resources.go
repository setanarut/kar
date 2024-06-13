package res

import (
	"embed"
	"image"
	_ "image/png"
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/engine/io"
	"kar/engine/util"
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
	BlockSize float64 = 48.0
)

var (
	PlayerAtlas         = io.LoadEbitenImageFromFS(assets, "assets/player_atlas.png")
	Atlas               = io.LoadEbitenImageFromFS(assets, "assets/atlas.png")
	CracksOverlayFrames = util.SubImages(io.LoadEbitenImageFromFS(assets, "assets/cracks.png"), 0, 0, 16, 16, 8, true)
	DirtStages          = makeBlockBreakingSpriteFrames(CracksOverlayFrames, 336, 208)
	StoneStages         = makeBlockBreakingSpriteFrames(CracksOverlayFrames, 96, 416)
)

var (
	World               donburi.World = donburi.NewWorld()
	Space               *cm.Space     = cm.NewSpace()
	Camera              *engine.Camera
	ScreenRect          cm.BB
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
	StatsTextOptions.ColorScale.ScaleWithColor(colornames.White)
}

func makeBlockBreakingSpriteFrames(cracksOverlayFrames []*ebiten.Image, x, y int) []*ebiten.Image {
	blockImage := util.SubImage(Atlas, x, y, 16, 16)
	blockStages := make([]*ebiten.Image, 9)
	for i := range blockStages {
		blockStages[i] = ebiten.NewImage(16, 16)
		blockStages[i].DrawImage(blockImage, nil)
	}
	dio := &ebiten.DrawImageOptions{
		Blend: ebiten.BlendSourceOver,
	}
	for i := 1; i < 9; i++ {
		blockStages[i].DrawImage(cracksOverlayFrames[i-1], dio)
	}
	return blockStages
}
