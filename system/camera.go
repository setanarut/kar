package system

import (
	"image/color"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/vectorg"
	"kar/items"
	"kar/res"
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
	p, ok := comp.PlayerTag.First(res.ECSWorld)

	if ok {
		pos := comp.Body.Get(p).Position()
		res.Cam = kamera.NewCamera(pos.X, pos.Y, res.ScreenSize.X, res.ScreenSize.Y)
	} else {
		pos := res.ScreenSize.Scale(0.5)
		res.Cam = kamera.NewCamera(pos.X, pos.Y, res.ScreenSize.X, res.ScreenSize.Y)
	}

	res.Cam.Lerp = true
}

func (ds *DrawCameraSystem) Update() {

	vectorg.GlobalTransform.Reset()
	cmdrawer.GeoM.Reset()
	res.Cam.ApplyCameraTransform(cmdrawer.GeoM)
	res.Cam.ApplyCameraTransform(vectorg.GlobalTransform)

	p, ok := comp.PlayerTag.First(res.ECSWorld)
	if ok {
		pos := comp.Body.Get(p).Position()

		res.Cam.LookAt(pos.X, pos.Y)

	}

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		res.Cam.ZoomFactor *= 1.02
	}

	if ebiten.IsKeyPressed(ebiten.KeyO) {
		res.Cam.ZoomFactor /= 1.02
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		res.Cam.AddTrauma(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		drawTargetBoxOverlayEnabled = !drawTargetBoxOverlayEnabled
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		res.Cam.ZoomFactor = 1
	}

	comp.AnimationPlayer.Each(res.ECSWorld, func(e *donburi.Entry) {
		comp.AnimationPlayer.Get(e).Update()
	})
}

func (ds *DrawCameraSystem) Draw(scr *ebiten.Image) {

	// clear color
	scr.Fill(color.Gray{30})

	comp.BlockTag.Each(res.ECSWorld, func(e *donburi.Entry) {
		body := comp.Body.Get(e)
		itemData := comp.Item.Get(e)
		drawOpt := comp.DrawOptions.Get(e)
		healthData := comp.Health.Get(e)
		pos := body.Position()
		health := mathutil.Clamp(healthData.Health, 0, healthData.MaxHealth)
		l := float64(len(res.SpriteFrames[itemData.ID]))
		blockSpriteFrameIndex := int(mathutil.MapRange(health, healthData.MaxHealth, 0, 0, l))

		ApplyDIO(drawOpt, pos)
		res.Cam.Draw(res.SpriteFrames[itemData.ID][blockSpriteFrameIndex], res.GlobalDrawOptions, scr)
	})

	// Drop Item
	comp.DropItemTag.Each(res.ECSWorld, func(e *donburi.Entry) {
		pos := comp.Body.Get(e).Position()
		drawOpt := comp.DrawOptions.Get(e)
		itemData := comp.Item.Get(e)

		// Item sin animation
		datai := comp.Index.Get(e)
		pos.Y += sinspace[datai.Index]

		ApplyDIO(drawOpt, pos)
		res.Cam.Draw(res.SpriteFrames[itemData.ID][0], res.GlobalDrawOptions, scr)
	})

	comp.DebugBoxTag.Each(res.ECSWorld, func(e *donburi.Entry) {
		b := comp.Body.Get(e)
		pos := b.Position()
		drawOpt := comp.DrawOptions.Get(e)
		drawOpt.Rotation = b.Angle()
		ApplyDIO(drawOpt, pos)
		res.Cam.Draw(res.SpriteFrames[items.Stone][0], res.GlobalDrawOptions, scr)
	})

	playerEntry, ok := comp.PlayerTag.First(res.ECSWorld)
	if ok {

		// draw player
		pBody := comp.Body.Get(playerEntry)
		drawOpt := comp.DrawOptions.Get(playerEntry)
		playerPos := pBody.Position()
		ApplyDIO(drawOpt, playerPos)
		ap := comp.AnimationPlayer.Get(playerEntry)
		if ap.CurrentFrame != nil {
			res.Cam.Draw(ap.CurrentFrame, res.GlobalDrawOptions, scr)
		}

		if ChipmunkDebugSpaceDrawing {
			cm.DrawSpace(res.Space, cmdrawer.WithScreen(scr))
			vectorg.Line(scr, playerPos, attackSegmentEnd, 1, color.White)
		}

		if drawTargetBoxOverlayEnabled {
			if hitShape != nil {
				if true {
					dio := &ebiten.DrawImageOptions{}
					dio.ColorScale.ScaleWithColor(colornames.Black)
					dio.GeoM.Translate(
						currentBlockPos.X+res.BlockCenterOffset.X,
						currentBlockPos.Y+res.BlockCenterOffset.Y)
					res.Cam.Draw(res.Border48, dio, scr)
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
	res.GlobalDrawOptions.GeoM.Reset()
	res.GlobalDrawOptions.GeoM.Translate(drawOpt.CenterOffset.X, drawOpt.CenterOffset.Y)
	res.GlobalDrawOptions.GeoM.Scale(scl.X, scl.Y)
	res.GlobalDrawOptions.GeoM.Rotate(drawOpt.Rotation)
	res.GlobalDrawOptions.GeoM.Translate(pos.X, pos.Y)
	res.GlobalDrawOptions.ColorScale.Reset()
}
