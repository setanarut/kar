package system

import (
	"image/color"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/vectorg"
	"kar/items"
	"kar/resources"
	"kar/types"
	"math"

	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/vec"
	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

var drawTargetBoxOverlayEnabled bool = true
var ChipmunkDebugSpaceDrawing = false
var ItemAnimFrameCount = 100
var sinspace []float64 = mathutil.SinSpace(0, 2*math.Pi, 4, ItemAnimFrameCount+1)
var cmdrawer = ebitencm.NewDrawer()

// DrawCameraSystem
type DrawCameraSystem struct {
}

func NewDrawCameraSystem() *DrawCameraSystem {
	return &DrawCameraSystem{}
}

func (ds *DrawCameraSystem) Init() {
	vectorg.GlobalTransform = &ebiten.GeoM{}
	p, ok := comp.TagPlayer.First(resources.ECSWorld)

	if ok {
		pos := comp.Body.Get(p).Position()
		resources.Cam = kamera.NewCamera(pos.X, pos.Y, resources.ScreenSize.X, resources.ScreenSize.Y)
	} else {
		pos := resources.ScreenSize.Scale(0.5)
		resources.Cam = kamera.NewCamera(pos.X, pos.Y, resources.ScreenSize.X, resources.ScreenSize.Y)
	}

	resources.Cam.Lerp = true
}

func (ds *DrawCameraSystem) Update() {

	vectorg.GlobalTransform.Reset()
	cmdrawer.GeoM.Reset()
	resources.Cam.ApplyCameraTransform(cmdrawer.GeoM)
	resources.Cam.ApplyCameraTransform(vectorg.GlobalTransform)

	p, ok := comp.TagPlayer.First(resources.ECSWorld)
	if ok {
		pos := comp.Body.Get(p).Position()

		resources.Cam.LookAt(pos.X, pos.Y)

	}

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		resources.Cam.ZoomFactor *= 1.02
	}

	if ebiten.IsKeyPressed(ebiten.KeyO) {
		resources.Cam.ZoomFactor /= 1.02
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		resources.Cam.AddTrauma(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		drawTargetBoxOverlayEnabled = !drawTargetBoxOverlayEnabled
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		resources.Cam.ZoomFactor = 1
	}

	comp.AnimPlayer.Each(resources.ECSWorld, func(e *donburi.Entry) {
		comp.AnimPlayer.Get(e).Update()
	})
}

func (ds *DrawCameraSystem) Draw(scr *ebiten.Image) {

	// clear color
	scr.Fill(color.RGBA{64, 68, 108, 255})

	comp.TagBlock.Each(resources.ECSWorld, func(e *donburi.Entry) {
		body := comp.Body.Get(e)
		itemData := comp.Item.Get(e)
		drawOpt := comp.DrawOptions.Get(e)
		healthData := comp.Health.Get(e)
		pos := body.Position()
		health := mathutil.Clamp(healthData.Health, 0, healthData.MaxHealth)
		l := float64(len(resources.SpriteStages[itemData.ID]))
		blockSpriteFrameIndex := int(mathutil.MapRange(health, healthData.MaxHealth, 0, 0, l))

		ApplyDIO(drawOpt, pos)
		resources.Cam.Draw(resources.SpriteStages[itemData.ID][blockSpriteFrameIndex], resources.GlobalDrawOptions, scr)
	})

	// Drop Item
	comp.TagItem.Each(resources.ECSWorld, func(e *donburi.Entry) {
		pos := comp.Body.Get(e).Position()
		drawOpt := comp.DrawOptions.Get(e)
		itemData := comp.Item.Get(e)
		// Item sin animation
		datai := comp.Index.Get(e)
		pos.Y += sinspace[datai.Index]
		ApplyDIO(drawOpt, pos)
		resources.Cam.Draw(getSprite(itemData.ID), resources.GlobalDrawOptions, scr)
	})

	comp.TagDebugBox.Each(resources.ECSWorld, func(e *donburi.Entry) {
		b := comp.Body.Get(e)
		pos := b.Position()
		drawOpt := comp.DrawOptions.Get(e)
		drawOpt.Rotation = b.Angle()
		ApplyDIO(drawOpt, pos)
		resources.Cam.Draw(resources.SpriteStages[items.Stone][0], resources.GlobalDrawOptions, scr)
	})

	playerEntry, ok := comp.TagPlayer.First(resources.ECSWorld)
	if ok {

		// draw player
		pBody := comp.Body.Get(playerEntry)
		drawOpt := comp.DrawOptions.Get(playerEntry)
		playerPos := pBody.Position()
		ApplyDIO(drawOpt, playerPos)
		ap := comp.AnimPlayer.Get(playerEntry)
		if ap.CurrentFrame != nil {
			resources.Cam.Draw(ap.CurrentFrame, resources.GlobalDrawOptions, scr)
		}

		if ChipmunkDebugSpaceDrawing {
			cm.DrawSpace(resources.Space, cmdrawer.WithScreen(scr))
			vectorg.Line(scr, playerPos, attackSegmentEnd, 1, color.White)
		}

		if drawTargetBoxOverlayEnabled {
			if hitShape != nil {
				if true {
					dio := &ebiten.DrawImageOptions{}
					dio.ColorScale.ScaleWithColor(colornames.Black)
					dio.GeoM.Translate(
						currentBlockPos.X+resources.BlockCenterOffset.X,
						currentBlockPos.Y+resources.BlockCenterOffset.Y)
					resources.Cam.Draw(resources.BlockHighlightBorder, dio, scr)
				}
			}
		}
	}
}

func ApplyDIO(drawOpt *types.DataDrawOptions, pos vec.Vec2) {

	scl := drawOpt.Scale
	if drawOpt.FlipX {
		scl.X *= -1
	}
	resources.GlobalDrawOptions.GeoM.Reset()
	resources.GlobalDrawOptions.GeoM.Translate(drawOpt.CenterOffset.X, drawOpt.CenterOffset.Y)
	resources.GlobalDrawOptions.GeoM.Scale(scl.X, scl.Y)
	resources.GlobalDrawOptions.GeoM.Rotate(drawOpt.Rotation)
	resources.GlobalDrawOptions.GeoM.Translate(pos.X, pos.Y)
	resources.GlobalDrawOptions.ColorScale.Reset()
}
