package system

import (
	"kar/comp"
	"kar/engine"
	"kar/res"
	"kar/types"

	"github.com/hajimehoshi/ebiten/v2"
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
	res.Camera.ZoomFactor = 0
	res.Camera.Lerp = true
	res.Camera.TraumaEnabled = false
}

func (ds *DrawCameraSystem) Update() {
	p, ok := comp.PlayerTag.First(res.World)
	if ok {
		pos := comp.Body.Get(p).Position()
		res.Camera.LookAt(engine.InvPosVectY(pos, res.ScreenRect.T))
	}

	if ebiten.IsKeyPressed(ebiten.KeyO) {
		res.Camera.ZoomFactor -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		res.Camera.ZoomFactor += 5
	}

	comp.AnimationPlayer.Each(res.World, func(e *donburi.Entry) {
		comp.AnimationPlayer.Get(e).Update()
	})

}

func (ds *DrawCameraSystem) Draw() {
	// clear color
	res.Screen.Fill(colornames.Black)
	comp.Block.Each(res.World, ds.DrawBlock)
	comp.AnimationPlayer.Each(res.World, ds.DrawAnimationPlayer)
}

func (ds *DrawCameraSystem) DrawAnimationPlayer(e *donburi.Entry) {

	body := comp.Body.Get(e)
	ap := comp.AnimationPlayer.Get(e)
	drawopt := comp.DrawOptions.Get(e)
	pos := engine.InvPosVectY(body.Position(), res.ScreenRect.T)

	ds.dio.GeoM.Reset()
	ds.dio.GeoM.Translate(drawopt.CenterOffset.X, drawopt.CenterOffset.Y)

	scl := drawopt.Scale
	if drawopt.FlipX {
		scl.X *= -1
	}
	ds.dio.GeoM.Scale(scl.X, scl.Y)
	ds.dio.GeoM.Rotate(engine.InvertAngle(drawopt.Rotation))
	ds.dio.GeoM.Translate(pos.X, pos.Y)
	res.Camera.Draw(ap.CurrentFrame, ds.dio, res.Screen)
	ds.dio.ColorScale.Reset()
}
func (ds *DrawCameraSystem) DrawBlock(e *donburi.Entry) {

	body := comp.Body.Get(e)
	b := comp.Block.Get(e)
	drawopt := comp.DrawOptions.Get(e)
	h := comp.Health.Get(e)
	pos := engine.InvPosVectY(body.Position(), res.ScreenRect.T)

	if b.BlockType == types.BlockStone {
		ds.currentFrame = res.StoneStages[int(engine.MapRange(h.Health, 10, 0, 0, 8))]
	}

	ds.dio.GeoM.Reset()
	ds.dio.GeoM.Translate(drawopt.CenterOffset.X, drawopt.CenterOffset.Y)
	ds.dio.GeoM.Scale(drawopt.Scale.X, drawopt.Scale.Y)
	ds.dio.GeoM.Rotate(engine.InvertAngle(drawopt.Rotation))
	ds.dio.GeoM.Translate(pos.X, pos.Y)
	ds.dio.ColorScale.Reset()

	if ds.currentFrame != nil {
		res.Camera.Draw(ds.currentFrame, ds.dio, res.Screen)
	}
}
