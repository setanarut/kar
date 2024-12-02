package system

import (
	"image/color"
	"kar"
	"kar/arc"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Render struct{}

func (rn *Render) Init() {
}

func (rn *Render) Update() {
	q := arc.FilterAnimPlayer.Query(&kar.WorldECS)

	for q.Next() {
		a := q.Get()
		a.Update()
	}

}

func (rn *Render) Draw() {
	kar.Screen.Fill(color.RGBA{64, 68, 108, 255})

	dio := &ebiten.DrawImageOptions{}
	kar.Camera.Draw(res.Border, dio, kar.Screen)

	q := arc.FilterDraw.Query(&kar.WorldECS)
	for q.Next() {
		dop, anim, rect := q.Get()

		ebitenutil.DebugPrint(kar.Screen, rect.String())

		anim.Update()
		kar.Camera.LookAt(rect.X, rect.Y)
		sclX := dop.Scale
		if dop.FlipX {
			sclX *= -1
		}
		kar.GlobalDIO.GeoM.Reset()
		kar.GlobalDIO.GeoM.Scale(sclX, dop.Scale)
		kar.GlobalDIO.GeoM.Translate(rect.X, rect.Y)

		kar.Camera.Draw(anim.CurrentFrame, kar.GlobalDIO, kar.Screen)

	}
}
