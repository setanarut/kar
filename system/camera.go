package system

import (
	"kar/comp"
	"kar/engine"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

// var im *ebiten.Image

// DrawCameraSystem
type DrawCameraSystem struct {
}

func NewDrawCameraSystem() *DrawCameraSystem {
	return &DrawCameraSystem{}
}

func (ds *DrawCameraSystem) Init() {
	res.Camera.ZoomFactor = 0
	res.Camera.Lerp = true
	res.Camera.TraumaEnabled = false
	res.StoneBlockImage.DrawImage(res.BlockBreakingStagesImages[0], nil)
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

	comp.Render.Each(res.World, func(e *donburi.Entry) {
		comp.Render.Get(e).AnimPlayer.Update()
	})

}

func (ds *DrawCameraSystem) Draw() {

	// clear color
	res.Screen.Fill(colornames.Black)

	comp.BlockTag.Each(res.World, ds.DrawEntrySprite)
	// comp.SnowballTag.Each(res.World, ds.DrawEntrySprite)

	if e, ok := comp.PlayerTag.First(res.World); ok {
		ds.DrawEntrySprite(e)
	}

}

func (ds *DrawCameraSystem) DrawEntryAnimation(e *donburi.Entry) {

	body := comp.Body.Get(e)
	render := comp.Render.Get(e)
	pos := engine.InvPosVectY(body.Position(), res.ScreenRect.T)

	render.DIO.GeoM.Reset()
	render.DIO.GeoM.Translate(render.Offset.X, render.Offset.Y)
	render.DIO.GeoM.Scale(render.DrawScale.X, render.DrawScale.Y)
	render.DIO.GeoM.Rotate(engine.InvertAngle(render.DrawAngle))
	render.DIO.GeoM.Translate(pos.X, pos.Y)

	if e.HasComponent(comp.Health) {
		v := engine.MapRange(comp.Health.GetValue(e), 0, 3, 0, 1)
		render.DIO.ColorScale.ScaleWithColor(res.DamageGradient.At(v))
	} else {
		render.DIO.ColorScale.Reset()
	}

	// res.Screen.DrawImage(render.AnimPlayer.CurrentFrame, render.DIO)
	res.Camera.Draw(render.AnimPlayer.CurrentFrame, render.DIO, res.Screen)
	render.DIO.ColorScale.Reset()
}
func (ds *DrawCameraSystem) DrawEntrySprite(e *donburi.Entry) {

	body := comp.Body.Get(e)
	sprite := comp.Sprite.Get(e)
	pos := engine.InvPosVectY(body.Position(), res.ScreenRect.T)

	sprite.DIO.GeoM.Reset()
	sprite.DIO.GeoM.Translate(sprite.Offset.X, sprite.Offset.Y)
	sprite.DIO.GeoM.Scale(sprite.DrawScale.X, sprite.DrawScale.Y)
	sprite.DIO.GeoM.Rotate(engine.InvertAngle(sprite.DrawAngle))
	sprite.DIO.GeoM.Translate(pos.X, pos.Y)
	res.Camera.Draw(sprite.Image, sprite.DIO, res.Screen)
	sprite.DIO.ColorScale.Reset()
}
