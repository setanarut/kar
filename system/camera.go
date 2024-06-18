package system

import (
	"kar/comp"
	"kar/engine"
	"kar/engine/mathutil"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

// DrawCameraSystem
type DrawCameraSystem struct {
	dio          *ebiten.DrawImageOptions
	currentFrame *ebiten.Image
}

func NewDrawCameraSystem() *DrawCameraSystem {
	return &DrawCameraSystem{
		dio: &ebiten.DrawImageOptions{},
	}
}

func (ds *DrawCameraSystem) Init() {

	p, ok := comp.PlayerTag.First(res.World)
	if ok {
		pos := comp.Body.Get(p).Position()
		res.Camera = engine.NewCamera(pos, res.ScreenSize.X, res.ScreenSize.Y)
	} else {
		res.Camera = engine.NewCamera(res.ScreenSizeF.Scale(0.5), res.ScreenSize.X, res.ScreenSize.Y)
	}
	res.Camera.ZoomFactor = 0
	res.Camera.Lerp = true
	res.Camera.TraumaEnabled = false
}

func (ds *DrawCameraSystem) Update() {
	p, ok := comp.PlayerTag.First(res.World)
	if ok {
		pos := comp.Body.Get(p).Position()
		res.Camera.LookAt(pos)
	}

	if ebiten.IsKeyPressed(ebiten.KeyO) {
		res.Camera.ZoomFactor -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		res.Camera.ZoomFactor += 5
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		res.Camera.ZoomFactor = 0
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
		ds.currentFrame = res.BlockFrames[itemData.Item][blockSpriteFrameIndex]

		ds.dio.GeoM.Reset()
		ds.dio.GeoM.Translate(drawOpt.CenterOffset.X, drawOpt.CenterOffset.Y)
		ds.dio.GeoM.Scale(drawOpt.Scale.X, drawOpt.Scale.Y)
		ds.dio.GeoM.Rotate(-drawOpt.Rotation)
		ds.dio.GeoM.Translate(pos.X, pos.Y)
		ds.dio.ColorScale.Reset()

		if ds.currentFrame != nil {
			res.Camera.Draw(ds.currentFrame, ds.dio, screen)
		}
	})
	comp.AnimationPlayer.Each(res.World, func(e *donburi.Entry) {
		body := comp.Body.Get(e)
		drawopt := comp.DrawOptions.Get(e)
		ap := comp.AnimationPlayer.Get(e)

		pos := body.Position()
		scl := drawopt.Scale
		if drawopt.FlipX {
			scl.X *= -1
		}

		ds.dio.GeoM.Reset()
		ds.dio.GeoM.Translate(drawopt.CenterOffset.X, drawopt.CenterOffset.Y)
		ds.dio.GeoM.Scale(scl.X, scl.Y)
		ds.dio.GeoM.Rotate(-drawopt.Rotation)
		ds.dio.GeoM.Translate(pos.X, pos.Y)
		if ap.CurrentFrame != nil {
			res.Camera.Draw(ap.CurrentFrame, ds.dio, screen)
		}
		ds.dio.ColorScale.Reset()
	})

}
