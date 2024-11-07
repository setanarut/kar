package system

import (
	"image/color"
	"kar"
	"kar/arc"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/engine/vectorg"
	"kar/items"
	"kar/res"

	"github.com/setanarut/cm"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/vec"

	eb "github.com/hajimehoshi/ebiten/v2"
)

type Render struct{}

func (rn *Render) Init() {

	vectorg.GlobalTransform = &eb.GeoM{}
	if kar.WorldECS.Alive(playerEntity) {
		x, y := playerSpawnPos.X, playerSpawnPos.Y
		camera = kamera.NewCamera(x, y, kar.ScreenSize.X, kar.ScreenSize.Y)
	} else {
		camera = kamera.NewCamera(0, 0, kar.ScreenSize.X, kar.ScreenSize.Y)
	}
	camera.LerpEnabled = true
}

func (rn *Render) Update() {

	if justPressed(eb.KeyX) {
		debugDrawingEnabled = !debugDrawingEnabled
	}

	if debugDrawingEnabled {
		vectorg.GlobalTransform.Reset()
		cmDrawer.GeoM.Reset()
		camera.ApplyCameraTransform(cmDrawer.GeoM)
		camera.ApplyCameraTransform(vectorg.GlobalTransform)
	}

	camera.LookAt(playerPos.X, playerPos.Y)

	if pressed(eb.KeyP) {
		camera.ZoomFactor *= 1.02
	}

	if pressed(eb.KeyO) {
		camera.ZoomFactor /= 1.02
	}

	if justPressed(eb.KeyT) {
		camera.AddTrauma(1)
	}
	if justPressed(eb.KeyV) {
		drawBlockBorderEnabled = !drawBlockBorderEnabled
	}

	if justPressed(eb.KeyBackspace) {
		camera.ZoomFactor = 1
	}

	q := arc.FilterAnimPlayer.Query(&kar.WorldECS)
	for q.Next() {
		q.Get().Update()
	}

}

func (rn *Render) Draw() {

	// Clear color
	kar.Screen.Fill(color.RGBA{64, 68, 108, 255})

	drawDropItems()
	drawBlocks()
	drawPlayer()

	if drawBlockBorderEnabled {
		drawBlockBorder()
	}

}

func drawBlockBorder() {
	if hitShape != nil {
		dio := &eb.DrawImageOptions{}
		dio.GeoM.Translate(hitBlockPos.X+blockCenterOffset.X, hitBlockPos.Y+blockCenterOffset.Y)
		camera.Draw(res.Border, dio, kar.Screen)
	}
}

func drawPlayer() {
	if kar.WorldECS.Alive(playerEntity) {
		applyDIO(playerDrawOptions, playerBody.Position())
		if playerAnim.CurrentFrame != nil {
			camera.Draw(playerAnim.CurrentFrame, globalDIO, kar.Screen)
		}
		if debugDrawingEnabled {
			cm.DrawSpace(kar.Space, cmDrawer.WithScreen(kar.Screen))
			vectorg.Line(kar.Screen, playerPos, attackSegEnd, 1, color.White)
		}
	}
}

func drawDropItems() {
	q := arc.FilterDropItem.Query(&kar.WorldECS)
	for q.Next() {
		dop, bd, itm, _, _, idx := q.Get()
		pos := bd.Body.Position()
		pos.Y += sinSpaceFrames[idx.Index]
		applyDIO(dop, pos)
		camera.Draw(getSprite(itm.ID), globalDIO, kar.Screen)
	}
}
func drawBlocks() {
	q := arc.FilterBlock.Query(&kar.WorldECS)
	for q.Next() {
		h, dop, bd, itm := q.Get()
		imgIndex := int(mathutil.MapRange(h.Health, h.MaxHealth, 0, 0, 10))
		if util.CheckIndex(res.Frames[itm.ID], imgIndex) {
			applyDIO(dop, bd.Body.Position())
			if items.IsHarvestable(itm.ID) {
				camera.Draw(res.Frames[itm.ID][0], globalDIO, kar.Screen)
			} else {
				camera.Draw(res.Frames[itm.ID][imgIndex], globalDIO, kar.Screen)
			}
		}
	}
}

func applyDIO(drawOpt *arc.DrawOptions, pos vec.Vec2) {
	scl := drawOpt.Scale
	if drawOpt.FlipX {
		scl.X *= -1
	}
	globalDIO.GeoM.Reset()
	globalDIO.GeoM.Translate(
		drawOpt.CenterOffset.X,
		drawOpt.CenterOffset.Y,
	)
	globalDIO.GeoM.Scale(scl.X, scl.Y)
	// globalDIO.GeoM.Rotate(drawOpt.Rotation)
	globalDIO.GeoM.Translate(pos.X, pos.Y)
	globalDIO.ColorScale.Reset()
}
