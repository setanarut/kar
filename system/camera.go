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

var (
	drawTargetBoxOverlayEnabled bool      = true
	ChipmunkDebugSpaceDrawing             = false
	ItemAnimFrameCount                    = 100
	sinspace                    []float64 = mathutil.SinSpace(
		0,
		2*math.Pi,
		4,
		ItemAnimFrameCount+1,
	)
	cmdrawer        = ebitencm.NewDrawer()
	CamBB           cm.BB
	CamBoundenabled bool
)

type DrawCameraSystem struct{}

func (ds *DrawCameraSystem) Init() {
	vectorg.GlobalTransform = &ebiten.GeoM{}

	p, ok := comp.TagPlayer.First(res.ECSWorld)

	if ok {
		pos := comp.Body.Get(p).Position()
		res.Cam = kamera.NewCamera(
			pos.X,
			pos.Y,
			res.ScreenSize.X,
			res.ScreenSize.Y,
		)
	} else {
		pos := res.ScreenSize.Scale(0.5)
		res.Cam = kamera.NewCamera(
			pos.X,
			pos.Y,
			res.ScreenSize.X,
			res.ScreenSize.Y)
	}

	res.Cam.Lerp = true
	// resources.Cam.LerpSpeed = 0.2
}

func (ds *DrawCameraSystem) Update() {
	tx, ty := res.Cam.Target()
	CamBB = cm.NewBBForExtents(
		vec.Vec2{tx, ty},
		res.ScreenSize.X/1.8,
		res.ScreenSize.Y/1.8,
	)
	vectorg.GlobalTransform.Reset()
	cmdrawer.GeoM.Reset()
	res.Cam.ApplyCameraTransform(cmdrawer.GeoM)
	res.Cam.ApplyCameraTransform(vectorg.GlobalTransform)

	p, ok := comp.TagPlayer.First(res.ECSWorld)
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

	comp.AnimPlayer.Each(res.ECSWorld, func(e *donburi.Entry) {
		comp.AnimPlayer.Get(e).Update()
	})
}

func (ds *DrawCameraSystem) Draw(scr *ebiten.Image) {

	// Clear color
	scr.Fill(color.RGBA{64, 68, 108, 255})

	// Draw blocks
	comp.TagBlock.Each(res.ECSWorld, func(e *donburi.Entry) {
		body := comp.Body.Get(e)
		pos := body.Position()
		if CamBB.ContainsVect(pos) {
			itemData := comp.Item.Get(e)
			healthData := comp.Health.Get(e)
			imgIndex := int(
				mathutil.MapRange(healthData.Health, healthData.MaxHealth, 0, 0, 10),
			)
			if CheckIndex(res.Frames[itemData.ID], imgIndex) {
				drawOpt := comp.DrawOptions.Get(e)
				ApplyDIO(drawOpt, pos)
				res.Cam.Draw(
					res.Frames[itemData.ID][imgIndex],
					res.GlobalDIO,
					scr,
				)
			}
		}
	})

	// Draw player
	if playerEntry, ok := comp.TagPlayer.First(res.ECSWorld); ok {
		pBody := comp.Body.Get(playerEntry)
		drawOpt := comp.DrawOptions.Get(playerEntry)
		playerPos := pBody.Position()
		ApplyDIO(drawOpt, playerPos)
		ap := comp.AnimPlayer.Get(playerEntry)
		if ap.CurrentFrame != nil {
			res.Cam.Draw(ap.CurrentFrame, res.GlobalDIO, scr)
		}
		if ChipmunkDebugSpaceDrawing {
			cm.DrawSpace(res.Space, cmdrawer.WithScreen(scr))
			vectorg.Line(scr, playerPos, attackSegEnd, 1, color.White)
		}
	}

	// Draw harvestable blocks
	comp.TagHarvestable.Each(res.ECSWorld, func(e *donburi.Entry) {
		body := comp.Body.Get(e)
		pos := body.Position()
		if CamBB.ContainsVect(pos) {
			itemData := comp.Item.Get(e)
			drawOpt := comp.DrawOptions.Get(e)
			ApplyDIO(drawOpt, pos)
			res.Cam.Draw(
				res.Frames[itemData.ID][0],
				res.GlobalDIO,
				scr,
			)
		}
	})

	// Draw drop items
	comp.TagItem.Each(res.ECSWorld, func(e *donburi.Entry) {
		pos := comp.Body.Get(e).Position()
		if CamBB.ContainsVect(pos) {
			drawOpt := comp.DrawOptions.Get(e)
			itemData := comp.Item.Get(e)
			// Item sin animation
			datai := comp.Index.Get(e)
			pos.Y += sinspace[datai.Index]
			ApplyDIO(drawOpt, pos)
			res.Cam.Draw(
				getSprite(itemData.ID),
				res.GlobalDIO,
				scr,
			)
		}
	})

	comp.TagDebugBox.Each(res.ECSWorld, func(e *donburi.Entry) {
		b := comp.Body.Get(e)
		pos := b.Position()
		drawOpt := comp.DrawOptions.Get(e)
		drawOpt.Rotation = b.Angle()
		ApplyDIO(drawOpt, pos)
		res.Cam.Draw(
			res.Frames[items.Stone][0],
			res.GlobalDIO,
			scr,
		)
	})

	if drawTargetBoxOverlayEnabled {
		if hitShape != nil {
			if true {
				dio := &ebiten.DrawImageOptions{}
				dio.ColorScale.ScaleWithColor(colornames.Black)
				dio.GeoM.Translate(
					currentBlockPos.X+res.BlockCenterOffset.X,
					currentBlockPos.Y+res.BlockCenterOffset.Y)
				res.Cam.Draw(res.BlockBorder, dio, scr)
			}
		}
	}
}

func ApplyDIO(drawOpt *types.DrawOptions, pos vec.Vec2) {

	scl := drawOpt.Scale
	if drawOpt.FlipX {
		scl.X *= -1
	}
	res.GlobalDIO.GeoM.Reset()
	res.GlobalDIO.GeoM.Translate(
		drawOpt.CenterOffset.X,
		drawOpt.CenterOffset.Y,
	)
	res.GlobalDIO.GeoM.Scale(scl.X, scl.Y)
	res.GlobalDIO.GeoM.Rotate(drawOpt.Rotation)
	res.GlobalDIO.GeoM.Translate(pos.X, pos.Y)
	res.GlobalDIO.ColorScale.Reset()
}
