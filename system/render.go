package system

import (
	"image/color"
	"kar"
	"kar/arc"

	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/vec"
)

type Render struct{}

func (rn *Render) Init() {
	kar.Camera = kamera.NewCamera(0, 0, kar.ScreenSize.X, kar.ScreenSize.Y)
}

func (rn *Render) Update() {
}

func (rn *Render) Draw() {
	kar.Screen.Fill(color.RGBA{64, 68, 108, 255})
	// drawBlocks()
	// drawPlayer()
}

func applyDIO(drawOpt *arc.DrawOptions, pos vec.Vec2) {
	sclX := drawOpt.Scale
	if drawOpt.FlipX {
		sclX *= -1
	}
	kar.GlobalDIO.GeoM.Reset()
	kar.GlobalDIO.GeoM.Scale(sclX, drawOpt.Scale)
	// globalDIO.GeoM.Rotate(drawOpt.Rotation)
	kar.GlobalDIO.GeoM.Translate(pos.X, pos.Y)
	kar.GlobalDIO.ColorScale.Reset()
}
