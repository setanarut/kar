package system

import (
	"kar/comp"
	"kar/engine"
	"kar/engine/mathutil"
	"kar/engine/vec"
	"kar/res"
	"kar/types"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

var dio *ebiten.DrawImageOptions

// DrawCameraSystem
type DrawCameraSystem struct {
}

func NewDrawCameraSystem() *DrawCameraSystem {
	return &DrawCameraSystem{}
}

func (ds *DrawCameraSystem) Init() {
	dio = &ebiten.DrawImageOptions{}

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

		ApplyDIO(drawOpt, pos)
		res.Camera.Draw(res.SpriteFrames[itemData.Item][blockSpriteFrameIndex], dio, screen)
	})

	comp.DropItemTag.Each(res.World, func(e *donburi.Entry) {
		pos := comp.Body.Get(e).Position()
		drawOpt := comp.DrawOptions.Get(e)
		itemData := comp.Item.Get(e)

		ApplyDIO(drawOpt, pos)
		res.Camera.Draw(res.SpriteFrames[itemData.Item][0], dio, screen)
	})

	e, ok := comp.PlayerTag.First(res.World)

	if ok {
		pos := comp.Body.Get(e).Position()
		drawOpt := comp.DrawOptions.Get(e)
		ap := comp.AnimationPlayer.Get(e)

		ApplyDIO(drawOpt, pos)
		if ap.CurrentFrame != nil {
			res.Camera.Draw(ap.CurrentFrame, dio, screen)
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
	dio.GeoM.Rotate(-drawOpt.Rotation)
	dio.GeoM.Translate(pos.X, pos.Y)
	dio.ColorScale.Reset()
}
