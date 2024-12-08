package system

import (
	"kar"
	"kar/arc"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Render struct{}

func (rn *Render) Init() {
}

func (rn *Render) Update() {

	kar.Camera.LookAt(playerCenterX, playerCenterY)

	q := arc.FilterAnimPlayer.Query(&kar.WorldECS)

	for q.Next() {
		a := q.Get()
		a.Update()
	}

}

func (rn *Render) Draw() {

	// Draw tilemap
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
					kar.TileColor,
					false,
				)
			}
		}
	}

	// Draw player
	q := arc.FilterDraw.Query(&kar.WorldECS)
	for q.Next() {
		dop, anim, rect := q.Get()
		sclX := dop.Scale
		kar.GlobalDIO.GeoM.Reset()
		if dop.FlipX {
			sclX *= -1
			kar.GlobalDIO.GeoM.Scale(sclX, dop.Scale)
			kar.GlobalDIO.GeoM.Translate(rect.X+rect.W, rect.Y)
		} else {
			kar.GlobalDIO.GeoM.Scale(sclX, dop.Scale)
			kar.GlobalDIO.GeoM.Translate(rect.X, rect.Y)
		}
		kar.Camera.Draw(anim.CurrentFrame, kar.GlobalDIO, kar.Screen)
	}

	// Draw target tile
	px, py := float64(targetBlock.X*Map.TileW), float64(targetBlock.Y*Map.TileH)
	px, py = kar.Camera.ApplyCameraTransformToPoint(px, py)
	vector.StrokeRect(
		kar.Screen,
		float32(px),
		float32(py),
		float32(Map.TileW),
		float32(Map.TileH),
		1,
		kar.TargetTileBorderColor,
		false,
	)

	// Draw debug info
	ebitenutil.DebugPrintAt(kar.Screen, PlayerController.CurrentState, 10, 10)
	ebitenutil.DebugPrintAt(kar.Screen, "InputLast"+PlayerController.InputAxisLast.String(), 10, 20)
	ebitenutil.DebugPrintAt(kar.Screen, "Target Block"+targetBlock.String(), 10, 30)
}
