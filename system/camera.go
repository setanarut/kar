package system

import (
	"image/color"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/vectorg"
	"kar/items"
	"kar/res"
	"kar/types"

	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/vec"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

var dio *ebiten.DrawImageOptions
var debugDrawEnabled bool

// DrawCameraSystem
type DrawCameraSystem struct {
}

func NewDrawCameraSystem() *DrawCameraSystem {
	return &DrawCameraSystem{}
}

func (ds *DrawCameraSystem) Init() {
	debugDrawEnabled = true
	vectorg.Transform = &ebiten.GeoM{}
	dio = &ebiten.DrawImageOptions{}
	p, ok := comp.PlayerTag.First(res.World)

	if ok {
		pos := comp.Body.Get(p).Position()
		res.Cam = kamera.NewCamera(pos.X, pos.Y, res.ScreenSizeF.X, res.ScreenSizeF.Y)
	} else {
		pos := res.ScreenSizeF.Scale(0.5)
		res.Cam = kamera.NewCamera(pos.X, pos.Y, res.ScreenSizeF.X, res.ScreenSizeF.Y)
	}

	res.Cam.ZoomFactor = 0
	res.Cam.Lerp = true
}

func (ds *DrawCameraSystem) Update() {

	if debugDrawEnabled {
		vectorg.Transform.Reset()
		res.Cam.ApplyCameraTransform(vectorg.Transform)
	}

	p, ok := comp.PlayerTag.First(res.World)
	if ok {
		pos := comp.Body.Get(p).Position()
		res.Cam.LookAt(pos.X, pos.Y)

	}

	if ebiten.IsKeyPressed(ebiten.KeyO) {
		res.Cam.ZoomFactor -= 5
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		res.Cam.AddTrauma(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		debugDrawEnabled = !debugDrawEnabled
	}

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		res.Cam.ZoomFactor += 5
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		res.Cam.ZoomFactor = 0
	}

	comp.AnimationPlayer.Each(res.World, func(e *donburi.Entry) {
		comp.AnimationPlayer.Get(e).Update()
	})

}

func (ds *DrawCameraSystem) Draw(screen *ebiten.Image) {
	// clear color
	screen.Fill(colornames.Black)

	comp.BlockItemTag.Each(res.World, func(e *donburi.Entry) {
		body := comp.Body.Get(e)
		itemData := comp.Item.Get(e)
		drawOpt := comp.DrawOptions.Get(e)
		healthData := comp.Health.Get(e)
		pos := body.Position()
		health := mathutil.Clamp(healthData.Health, 0, healthData.MaxHealth)
		blockSpriteFrameIndex := int(mathutil.MapRange(health, healthData.MaxHealth, 0, 0, 8))

		ApplyDIO(drawOpt, pos)
		res.Cam.Draw(res.SpriteFrames[itemData.Item][blockSpriteFrameIndex], dio, screen)
	})

	comp.DropItemTag.Each(res.World, func(e *donburi.Entry) {
		pos := comp.Body.Get(e).Position()
		drawOpt := comp.DrawOptions.Get(e)
		itemData := comp.Item.Get(e)

		ApplyDIO(drawOpt, pos)
		res.Cam.Draw(res.SpriteFrames[itemData.Item][0], dio, screen)
	})
	comp.DebugBoxTag.Each(res.World, func(e *donburi.Entry) {
		b := comp.Body.Get(e)
		pos := b.Position()
		drawOpt := comp.DrawOptions.Get(e)
		drawOpt.Rotation = b.Angle()
		ApplyDIO(drawOpt, pos)
		res.Cam.Draw(res.SpriteFrames[items.Stone][0], dio, screen)
	})

	playerEntry, ok := comp.PlayerTag.First(res.World)

	if ok {
		pos := comp.Body.Get(playerEntry).Position()
		drawOpt := comp.DrawOptions.Get(playerEntry)
		ap := comp.AnimationPlayer.Get(playerEntry)

		if debugDrawEnabled {
			// Draw attack raycast line
			vectorg.StrokeLine(screen, pos, attackSegmentEnd, 2, color.White, true)
			if hitShape != nil {
				hitBlockPos := hitShape.Body().Position()
				vectorg.StrokeSquare(screen, hitBlockPos, res.BlockSize, 2, color.White, false)
				placeBlockPos := hitBlockPos.Add(attackSegmentQuery.Normal.Scale(res.BlockSize))
				vectorg.StrokeSquare(screen, placeBlockPos, res.BlockSize, 1, colornames.Orange, true)
			}
		}

		ApplyDIO(drawOpt, pos)

		if ap.CurrentFrame != nil {
			res.Cam.Draw(ap.CurrentFrame, dio, screen)
		}

	}

}

func ApplyDIO(drawOpt *types.DataDrawOptions, pos vec.Vec2) {

	scl := drawOpt.Scale
	if drawOpt.FlipX {
		scl.X *= -1
	}
	dio.GeoM.Reset()
	dio.GeoM.Translate(drawOpt.CenterOffset.X, drawOpt.CenterOffset.Y)
	dio.GeoM.Scale(scl.X, scl.Y)
	dio.GeoM.Rotate(drawOpt.Rotation)
	dio.GeoM.Translate(pos.X, pos.Y)
	dio.ColorScale.Reset()
}
