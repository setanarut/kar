package system

import (
	"kar/comp"
	"kar/engine/mathutil"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

// var DebugImg = ebiten.NewImage(32, 32)
// var DebugImgCenterOffset = vec.Vec2{-16, -16}

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
	// DebugImg.Fill(color.RGBA{0, 0, 255, 40})
	res.Camera.ZoomFactor = 0
	res.Camera.Lerp = true
	res.Camera.TraumaEnabled = false
}

func (ds *DrawCameraSystem) Update() {
	p, ok := comp.PlayerTag.First(res.World)
	if ok {
		pos := comp.Body.Get(p).Position()
		res.Camera.LookAt(pos.FlipVertical(res.ScreenRect.T))
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

func (ds *DrawCameraSystem) Draw() {
	// clear color
	res.Screen.Fill(colornames.Black)
	comp.Block.Each(res.World, ds.DrawBlock)

	comp.AnimationPlayer.Each(res.World, ds.DrawAnimationPlayer)

}

func (ds *DrawCameraSystem) DrawAnimationPlayer(e *donburi.Entry) {

	body := comp.Body.Get(e)
	drawopt := comp.DrawOptions.Get(e)
	ap := comp.AnimationPlayer.Get(e)

	pos := body.Position().FlipVertical(res.ScreenRect.T)
	scl := drawopt.Scale
	if drawopt.FlipX {
		scl.X *= -1
	}

	ds.dio.GeoM.Reset()
	ds.dio.GeoM.Translate(drawopt.CenterOffset.X, drawopt.CenterOffset.Y)
	ds.dio.GeoM.Scale(scl.X, scl.Y)
	ds.dio.GeoM.Rotate(-drawopt.Rotation)
	ds.dio.GeoM.Translate(pos.X, pos.Y)

	res.Camera.Draw(ap.CurrentFrame, ds.dio, res.Screen)
	ds.dio.ColorScale.Reset()

	// //debug
	// dioDebug := &ebiten.DrawImageOptions{}
	// dioDebug.GeoM.Translate(DebugImgCenterOffset.X, DebugImgCenterOffset.Y)
	// dioDebug.GeoM.Translate(pos.X, pos.Y)
	// res.Camera.Draw(DebugImg, dioDebug, res.Screen)

}

func (ds *DrawCameraSystem) DrawBlock(e *donburi.Entry) {

	body := comp.Body.Get(e)
	blockData := comp.Block.Get(e)
	drawOpt := comp.DrawOptions.Get(e)
	healthData := comp.Health.Get(e)
	pos := body.Position().FlipVertical(res.ScreenRect.T)

	if blockData.BlockType == res.BlockStone {
		blockSpriteStageIndex := int(mathutil.MapRange(healthData.Health, healthData.MaxHealth, 0, 0, 8))
		ds.currentFrame = res.StoneStages[blockSpriteStageIndex]
	}
	if blockData.BlockType == res.BlockDirt {
		blockSpriteStageIndex := int(mathutil.MapRange(healthData.Health, healthData.MaxHealth, 0, 0, 8))
		ds.currentFrame = res.DirtStages[blockSpriteStageIndex]
	}

	ds.dio.GeoM.Reset()
	ds.dio.GeoM.Translate(drawOpt.CenterOffset.X, drawOpt.CenterOffset.Y)
	ds.dio.GeoM.Scale(drawOpt.Scale.X, drawOpt.Scale.Y)
	ds.dio.GeoM.Rotate(-drawOpt.Rotation)
	ds.dio.GeoM.Translate(pos.X, pos.Y)
	ds.dio.ColorScale.Reset()

	if ds.currentFrame != nil {
		res.Camera.Draw(ds.currentFrame, ds.dio, res.Screen)
	}
}
