package system

import (
	"kar"
	"kar/arc"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
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
				kar.GlobalDIO.GeoM.Reset()
				kar.GlobalDIO.GeoM.Scale(3, 3)
				kar.GlobalDIO.GeoM.Translate(px, py)
				kar.Camera.Draw(GetSprite(value), kar.GlobalDIO, kar.Screen)
			}
		}
	}

	q := arc.FilterDraw.Query(&kar.WorldECS)
	for q.Next() {
		dop, anim, rect := q.Get()

		// Draw player
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

		// Draw player hit box for debug
		if kar.DrawPlayerDebugHitBox {
			x, y := kar.Camera.ApplyCameraTransformToPoint(rect.X, rect.Y)
			vector.StrokeRect(
				kar.Screen,
				float32(x),
				float32(y),
				float32(rect.W),
				float32(rect.H),
				1,
				colornames.Magenta,
				false,
			)
			// Draw player center
			vector.DrawFilledCircle(
				kar.Screen,
				float32(x+rect.W*0.5),
				float32(y+rect.H*0.5),
				2,
				colornames.Magenta,
				false,
			)
		}
	}

	// Draw target tile
	if isBlockPlaceable {
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
	}

	// Draw debug info
	ebitenutil.DebugPrintAt(kar.Screen, PlayerController.CurrentState, 10, 10)
	ebitenutil.DebugPrintAt(kar.Screen, "InputLast"+PlayerController.InputAxisLast.String(), 10, 20)
	ebitenutil.DebugPrintAt(kar.Screen, "Target Block"+targetBlock.String(), 10, 30)
}
