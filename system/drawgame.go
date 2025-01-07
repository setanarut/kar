package system

import (
	"kar"
	"kar/arc"
	"kar/engine/mathutil"
	"kar/res"
)

type DrawGame struct{}

func (d *DrawGame) Init() {}

func (d *DrawGame) Update() {

	if !craftingState {
		if playerCenterX < kar.Camera.TopLeftX {
			kar.Camera.TopLeftX -= kar.Camera.Width()
		}
		if playerCenterX > kar.Camera.Right() {
			kar.Camera.TopLeftX += kar.Camera.Width()
		}
		if playerCenterY < kar.Camera.TopLeftY {
			kar.Camera.TopLeftY -= kar.Camera.Height()
		}
		if playerCenterY > kar.Camera.Bottom() {
			kar.Camera.TopLeftY += kar.Camera.Height()
		}

		q := arc.FilterAnimPlayer.Query(&kar.WorldECS)
		for q.Next() {
			a := q.Get()
			a.Update()
		}

	}
}

func (d *DrawGame) Draw() {

	// Draw tilemap
	camMin := tileMap.WorldToTile(kar.Camera.TopLeft())
	camMin.X = min(max(camMin.X, 0), tileMap.W)
	camMin.Y = min(max(camMin.Y, 0), tileMap.H)
	camMaxX := min(max(camMin.X+kar.RenderArea.X, 0), tileMap.W)
	camMaxY := min(max(camMin.Y+kar.RenderArea.Y, 0), tileMap.H)
	for y := camMin.Y; y < camMaxY; y++ {
		for x := camMin.X; x < camMaxX; x++ {
			tileID := tileMap.Grid[y][x]
			if tileID != 0 {
				px, py := float64(x*tileMap.TileW), float64(y*tileMap.TileH)
				kar.GlobalColorMDIO.GeoM.Reset()
				kar.GlobalColorMDIO.GeoM.Translate(px, py)
				if x == targetBlockPos.X && y == targetBlockPos.Y {
					i := mathutil.MapRange(blockHealth, 0, 180, 0, 5)
					if res.BlockCrackFrames[tileID] != nil {
						kar.Camera.DrawWithColorM(res.BlockCrackFrames[tileID][int(i)], kar.GlobalColorM, kar.GlobalColorMDIO, kar.Screen)
					}
				} else {
					kar.Camera.DrawWithColorM(res.BlockCrackFrames[tileID][0], kar.GlobalColorM, kar.GlobalColorMDIO, kar.Screen)
				}
			}
		}
	}

	if kar.WorldECS.Alive(player) {
		// Draw target tile border
		if isRayHit {
			kar.GlobalColorMDIO.GeoM.Reset()
			kar.GlobalColorMDIO.GeoM.Translate(
				float64(targetBlockPos.X*tileMap.TileW)-1,
				float64(targetBlockPos.Y*tileMap.TileH)-1,
			)
			kar.Camera.DrawWithColorM(res.SelectionBlock, kar.GlobalColorM, kar.GlobalColorMDIO, kar.Screen)
		}

		// Draw player
		kar.GlobalColorMDIO.GeoM.Reset()
		kar.GlobalColorMDIO.GeoM.Scale(ctrl.FlipXFactor, 1)
		if ctrl.FlipXFactor == -1 {
			kar.GlobalColorMDIO.GeoM.Translate(ctrl.Rect.X+ctrl.Rect.W, ctrl.Rect.Y)
		} else {
			kar.GlobalColorMDIO.GeoM.Translate(ctrl.Rect.X, ctrl.Rect.Y)
		}
		kar.Camera.DrawWithColorM(ctrl.AnimPlayer.CurrentFrame, kar.GlobalColorM, kar.GlobalColorMDIO, kar.Screen)
		// }
	}

	// Draw Items
	itemQuery := arc.FilterItem.Query(&kar.WorldECS)
	for itemQuery.Next() {
		id, rect, timers, _ := itemQuery.Get()
		kar.GlobalColorMDIO.GeoM.Reset()
		kar.GlobalColorMDIO.GeoM.Translate(rect.X, rect.Y+sinspace[timers.AnimationIndex])
		kar.Camera.DrawWithColorM(res.Icon8[id.ID], kar.GlobalColorM, kar.GlobalColorMDIO, kar.Screen)
	}

}
