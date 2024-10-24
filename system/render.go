package system

import (
	"image/color"
	"kar"
	"kar/comp"
	"kar/items"
	"kar/res"
	"kar/types"
	"kar/util"
	"kar/util/vectorg"

	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/vec"

	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

var cameraBounds cm.BB
var cmDrawer = ebitencm.NewDrawer()
var globalDIO = &eb.DrawImageOptions{}

type Render struct{}

func (r *Render) Init() {

	vectorg.GlobalTransform = &eb.GeoM{}
	if playerEntry.Valid() {
		x, y := playerSpawnPoint.X, playerSpawnPoint.Y
		camera = kamera.NewCamera(x, y, kar.ScreenSize.X, kar.ScreenSize.Y)
	} else {
		camera = kamera.NewCamera(0, 0, kar.ScreenSize.X, kar.ScreenSize.Y)
	}
	camera.Lerp = true
}

func (r *Render) Update() {

	// tx, ty := camera.Target()
	// cameraBounds = cm.NewBBForExtents(
	// 	vec.Vec2{tx, ty},
	// 	kar.ScreenSize.X/1.2,
	// 	kar.ScreenSize.Y/1.2,
	// )
	vectorg.GlobalTransform.Reset()
	cmDrawer.GeoM.Reset()
	camera.ApplyCameraTransform(cmDrawer.GeoM)
	camera.ApplyCameraTransform(vectorg.GlobalTransform)
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

	comp.AnimPlayer.Each(ecsWorld, func(e *donburi.Entry) {
		comp.AnimPlayer.Get(e).Update()
	})
}

func (r *Render) Draw() {
	// Clear color
	kar.Screen.Fill(color.RGBA{64, 68, 108, 255})
	comp.DropItem.Each(ecsWorld, drawDropItem)
	comp.Block.Each(ecsWorld, drawBlock)
	drawPlayer()
	comp.Harvestable.Each(ecsWorld, drawHarvestableBlock)
	comp.DebugBox.Each(ecsWorld, drawDebugBox)
	if drawBlockBorderEnabled {
		drawBlockBorder()
	}
}

func drawBlockBorder() {
	if hitShape != nil {
		if true {
			dio := &eb.DrawImageOptions{}
			dio.GeoM.Translate(-8, -8)
			dio.GeoM.Translate(hitBlockPos.X, hitBlockPos.Y)
			camera.Draw(res.Border16, dio, kar.Screen)
		}
	}
}

func drawPlayer() {
	if playerEntry.Valid() {
		pBody := comp.Body.Get(playerEntry)
		drawOpt := comp.DrawOptions.Get(playerEntry)
		playerPos := pBody.Position()
		applyDIO(drawOpt, playerPos)
		ap := comp.AnimPlayer.Get(playerEntry)
		if ap.CurrentFrame != nil {
			camera.Draw(ap.CurrentFrame, globalDIO, kar.Screen)
		}
		if debugDrawingEnabled {
			cm.DrawSpace(cmSpace, cmDrawer.WithScreen(kar.Screen))
			vectorg.Line(kar.Screen, playerPos, attackSegEnd, 1, color.White)
		}
	}
}

func drawDropItem(e *donburi.Entry) {
	pos := comp.Body.Get(e).Position()
	if cameraBounds.ContainsVect(pos) {
		drawOpt := comp.DrawOptions.Get(e)
		itemData := comp.Item.Get(e)
		// Item sin animation
		datai := comp.Index.Get(e)
		pos.Y += sinSpace[datai.Index]
		applyDIO(drawOpt, pos)
		camera.Draw(getSprite(itemData.ID), globalDIO, kar.Screen)
	}
}

func drawBlock(e *donburi.Entry) {
	body := comp.Body.Get(e)
	pos := body.Position()
	// if cameraBounds.ContainsVect(pos) {
	itemData := comp.Item.Get(e)
	healthData := comp.Health.Get(e)
	imgIndex := int(
		util.MapRange(healthData.Health, healthData.MaxHealth, 0, 0, 10),
	)
	if util.CheckIndex(res.Frames[itemData.ID], imgIndex) {
		drawOpt := comp.DrawOptions.Get(e)
		applyDIO(drawOpt, pos)
		camera.Draw(
			res.Frames[itemData.ID][imgIndex],
			globalDIO,
			kar.Screen,
		)
	}
	// }
}

func drawHarvestableBlock(e *donburi.Entry) {
	body := comp.Body.Get(e)
	pos := body.Position()
	if cameraBounds.ContainsVect(pos) {
		itemData := comp.Item.Get(e)
		drawOpt := comp.DrawOptions.Get(e)
		applyDIO(drawOpt, pos)
		camera.Draw(
			res.Frames[itemData.ID][0],
			globalDIO,
			kar.Screen,
		)
	}
}

func drawDebugBox(e *donburi.Entry) {
	b := comp.Body.Get(e)
	pos := b.Position()
	drawOpt := comp.DrawOptions.Get(e)
	drawOpt.Rotation = b.Angle()
	applyDIO(drawOpt, pos)
	camera.Draw(
		res.Frames[items.Stone][0],
		globalDIO,
		kar.Screen,
	)
}

func applyDIO(drawOpt *types.DrawOptions, pos vec.Vec2) {
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
	globalDIO.GeoM.Rotate(drawOpt.Rotation)
	globalDIO.GeoM.Translate(pos.X, pos.Y)
	globalDIO.ColorScale.Reset()
}
