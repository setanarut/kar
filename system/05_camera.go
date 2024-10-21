package system

import (
	"image/color"
	"kar"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/vectorg"
	"kar/items"
	"kar/res"
	"kar/types"

	"github.com/setanarut/cm"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/vec"
	"golang.org/x/image/colornames"

	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type DrawCamera struct{}

func (ds *DrawCamera) Init() {

	vectorg.GlobalTransform = &eb.GeoM{}
	if playerEntry.Valid() {
		// pos := comp.Body.Get(PlayerEntry).Position()

		Camera = kamera.NewCamera(
			playerPos.X,
			playerPos.Y,
			kar.ScreenSize.X,
			kar.ScreenSize.Y,
		)
	} else {
		pos := kar.ScreenSize.Scale(0.5)
		Camera = kamera.NewCamera(
			pos.X,
			pos.Y,
			kar.ScreenSize.X,
			kar.ScreenSize.Y)
	}
	Camera.Lerp = true
}

func (ds *DrawCamera) Update() {
	tx, ty := Camera.Target()
	cameraBounds = cm.NewBBForExtents(
		vec.Vec2{tx, ty},
		kar.ScreenSize.X/1.8,
		kar.ScreenSize.Y/1.8,
	)
	vectorg.GlobalTransform.Reset()
	cmDrawer.GeoM.Reset()
	Camera.ApplyCameraTransform(cmDrawer.GeoM)
	Camera.ApplyCameraTransform(vectorg.GlobalTransform)

	if playerEntry, ok := comp.TagPlayer.First(ecsWorld); ok {
		pos := comp.Body.Get(playerEntry).Position()
		Camera.LookAt(pos.X, pos.Y)
	}

	if keyPressed(eb.KeyP) {
		Camera.ZoomFactor *= 1.02
	}

	if keyPressed(eb.KeyO) {
		Camera.ZoomFactor /= 1.02
	}

	if justPressed(eb.KeyT) {
		Camera.AddTrauma(1)
	}
	if justPressed(eb.KeyV) {
		drawBlockBorderEnabled = !drawBlockBorderEnabled
	}

	if justPressed(eb.KeyBackspace) {
		Camera.ZoomFactor = 1
	}

	comp.AnimPlayer.Each(ecsWorld, func(e *donburi.Entry) {
		comp.AnimPlayer.Get(e).Update()
	})
}

func (ds *DrawCamera) Draw() {
	// Clear color
	kar.Screen.Fill(color.RGBA{64, 68, 108, 255})
	comp.TagDropItem.Each(ecsWorld, drawDropItem)
	comp.TagBlock.Each(ecsWorld, drawBlock)
	drawPlayer()
	comp.TagHarvestable.Each(ecsWorld, drawHarvestableBlock)
	comp.TagDebugBox.Each(ecsWorld, drawDebugBox)
	if drawBlockBorderEnabled {
		drawBlockBorder()
	}
}

func drawBlockBorder() {
	if hitShape != nil {
		if true {
			dio := &eb.DrawImageOptions{}
			dio.ColorScale.ScaleWithColor(colornames.Black)
			dio.GeoM.Translate(
				currentBlockPos.X+blockCenterOffset.X,
				currentBlockPos.Y+blockCenterOffset.Y)
			Camera.Draw(res.BlockBorder64, dio, kar.Screen)
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
			Camera.Draw(ap.CurrentFrame, globalDIO, kar.Screen)
		}
		if debugDrawingEnabled {
			cm.DrawSpace(space, cmDrawer.WithScreen(kar.Screen))
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
		Camera.Draw(getSprite(itemData.ID), globalDIO, kar.Screen)
	}
}

func drawBlock(e *donburi.Entry) {
	body := comp.Body.Get(e)
	pos := body.Position()
	if cameraBounds.ContainsVect(pos) {
		itemData := comp.Item.Get(e)
		healthData := comp.Health.Get(e)
		imgIndex := int(
			mathutil.MapRange(healthData.Health, healthData.MaxHealth, 0, 0, 10),
		)
		if CheckIndex(res.Frames[itemData.ID], imgIndex) {
			drawOpt := comp.DrawOptions.Get(e)
			applyDIO(drawOpt, pos)
			Camera.Draw(
				res.Frames[itemData.ID][imgIndex],
				globalDIO,
				kar.Screen,
			)
		}
	}
}

func drawHarvestableBlock(e *donburi.Entry) {
	body := comp.Body.Get(e)
	pos := body.Position()
	if cameraBounds.ContainsVect(pos) {
		itemData := comp.Item.Get(e)
		drawOpt := comp.DrawOptions.Get(e)
		applyDIO(drawOpt, pos)
		Camera.Draw(
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
	Camera.Draw(
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
