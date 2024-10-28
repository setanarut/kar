package system

import (
	"fmt"
	"image/color"
	"kar"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/engine/vectorg"
	"kar/items"
	"kar/res"
	"kar/types"
	"kar/world"

	"github.com/setanarut/cm"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/vec"
	"golang.org/x/image/colornames"

	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

type Render struct{}

func (rn *Render) Init() {

	vectorg.GlobalTransform = &eb.GeoM{}
	if playerEntry.Valid() {
		x, y := playerSpawnPos.X, playerSpawnPos.Y
		camera = kamera.NewCamera(x, y, kar.ScreenSize.X, kar.ScreenSize.Y)
	} else {
		camera = kamera.NewCamera(0, 0, kar.ScreenSize.X, kar.ScreenSize.Y)
	}
	camera.Lerp = true
}

func (rn *Render) Update() {
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

func (rn *Render) Draw() {
	// Clear color
	kar.Screen.Fill(color.RGBA{64, 68, 108, 255})
	comp.TagDropItem.Each(ecsWorld, drawDropItem)
	comp.TagBlock.Each(ecsWorld, drawBlock)
	drawPlayer()
	comp.TagHarvestable.Each(ecsWorld, drawHarvestableBlock)
	comp.TagDebugBox.Each(ecsWorld, drawDebugBox)
	if drawBlockBorderEnabled {
		drawBlockBorder()
		drawBlockBorder2()
	}
}

func drawBlockBorder2() {
	if inpututil.IsMouseButtonJustPressed(eb.MouseButton0) {
		x, y := camera.ScreenToWorld(eb.CursorPosition())
		worldPos := vec.Vec2{x, y}
		pixPos := world.WorldToPixel(worldPos)
		fmt.Println("WorldToPixel", pixPos)
		gameWorld.Image.SetGray16(pixPos.X, pixPos.Y, color.Gray16{items.GoldOre})
		fmt.Println("PixelToWorld", world.PixelToWorld(pixPos.X, pixPos.Y))

		if hitShape != nil {
			dio := &eb.DrawImageOptions{}
			dio.ColorScale.ScaleWithColor(colornames.Black)
			dio.GeoM.Translate(
				placeBlockPos.X+blockCenterOffset.X,
				placeBlockPos.Y+blockCenterOffset.Y)
			camera.Draw(res.BlockBorder, dio, kar.Screen)
		}
	}
}
func drawBlockBorder() {
	if hitShape != nil {
		dio := &eb.DrawImageOptions{}
		dio.ColorScale.ScaleWithColor(colornames.Black)
		dio.GeoM.Translate(
			hitBlockPos.X+blockCenterOffset.X,
			hitBlockPos.Y+blockCenterOffset.Y)
		camera.Draw(res.BlockBorder, dio, kar.Screen)
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
	drawOpt := comp.DrawOptions.Get(e)
	itemData := comp.Item.Get(e)
	// Item sin animation
	datai := comp.Index.Get(e)
	pos.Y += sinSpace[datai.Index]
	applyDIO(drawOpt, pos)
	camera.Draw(getSprite(itemData.ID), globalDIO, kar.Screen)
}

func drawBlock(e *donburi.Entry) {
	body := comp.Body.Get(e)
	pos := body.Position()
	itemData := comp.Item.Get(e)
	healthData := comp.Health.Get(e)
	imgIndex := int(
		mathutil.MapRange(healthData.Health, healthData.MaxHealth, 0, 0, 10),
	)
	if util.CheckIndex(res.Frames[itemData.ID], imgIndex) {
		drawOpt := comp.DrawOptions.Get(e)
		applyDIO(drawOpt, pos)
		camera.Draw(res.Frames[itemData.ID][imgIndex], globalDIO, kar.Screen)
	}
}

func drawHarvestableBlock(e *donburi.Entry) {
	body := comp.Body.Get(e)
	pos := body.Position()
	itemData := comp.Item.Get(e)
	drawOpt := comp.DrawOptions.Get(e)
	applyDIO(drawOpt, pos)
	camera.Draw(
		res.Frames[itemData.ID][0],
		globalDIO,
		kar.Screen,
	)
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
