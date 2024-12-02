package system

import (
	"image/color"
	"kar"
	"kar/arc"

	"github.com/hajimehoshi/ebiten/v2/vector"
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

	for y, row := range Map.Grid {
		for x, value := range row {
			if value != 0 {
				px, py := float64(x*Map.TileW), float64(y*Map.TileH)
				px, py = kar.Camera.ApplyCameraTransformToPoint(px, py)
				vector.DrawFilledRect(
					kar.Screen,
					float32(px),
					float32(py),
					float32(Map.TileW),
					float32(Map.TileH),
					color.Gray{127},
					false,
				)
			}
		}
	}

	// draw player
	q := arc.FilterDraw.Query(&kar.WorldECS)
	for q.Next() {
		dop, anim, rect := q.Get()
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
