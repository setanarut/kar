package system

import (
	"image/color"
	"kar"
	"kar/arc"
	"kar/engine/mathutil"
	"kar/res"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
)

type Game struct{}

func (d *Game) Init() {
}

func (d *Game) Update() {

	if kar.ECWorld.Alive(kar.CurrentPlayer) {

		playerPos, playerSize := arc.MapRect.Get(kar.CurrentPlayer)
		playerCenterX, playerCenterY := playerPos.X+playerSize.W, playerPos.Y+playerSize.H

		// Toggle camera follow
		if inpututil.IsKeyJustPressed(ebiten.KeyL) {
			switch kar.CameraRes.SmoothType {
			case kamera.None:
				kar.CameraRes.SetCenter(playerCenterX, playerCenterY-40)
				kar.CameraRes.SmoothType = kamera.SmoothDamp
			case kamera.SmoothDamp:
				kar.CameraRes.SetCenter(playerCenterX, playerCenterY-40)
				kar.CameraRes.SmoothType = kamera.None
			}
		}

		// Static camera logic
		if kar.CameraRes.SmoothType == kamera.None {
			if playerCenterX < kar.CameraRes.TopLeftX {
				kar.CameraRes.TopLeftX -= kar.CameraRes.Width()
			}
			if playerCenterX > kar.CameraRes.Right() {
				kar.CameraRes.TopLeftX += kar.CameraRes.Width()
			}
			if playerCenterY < kar.CameraRes.TopLeftY {
				kar.CameraRes.TopLeftY -= kar.CameraRes.Height()
			}
			if playerCenterY > kar.CameraRes.Bottom() {
				kar.CameraRes.TopLeftY += kar.CameraRes.Height()
			}
		} else {
			kar.CameraRes.LookAt(math.Floor(playerCenterX), math.Floor(playerCenterY-40))
		}
	}

}

func (d *Game) Draw() {

	// DRAW TILEMAPs

	// clamp tilemap bounds
	camMin := kar.TileMapRes.WorldToTile(kar.CameraRes.TopLeft())
	camMin.X = min(max(camMin.X, 0), kar.TileMapRes.W)
	camMin.Y = min(max(camMin.Y, 0), kar.TileMapRes.H)
	camMaxX := min(max(camMin.X+kar.RenderArea.X, 0), kar.TileMapRes.W)
	camMaxY := min(max(camMin.Y+kar.RenderArea.Y, 0), kar.TileMapRes.H)
	// draw tiles
	for y := camMin.Y; y < camMaxY; y++ {
		for x := camMin.X; x < camMaxX; x++ {
			tileID := kar.TileMapRes.Grid[y][x]
			if tileID != 0 {
				px, py := float64(x*kar.TileMapRes.TileW), float64(y*kar.TileMapRes.TileH)
				kar.ColorMDIO.GeoM.Reset()
				kar.ColorMDIO.GeoM.Translate(px, py)
				if x == kar.GameDataRes.TargetBlockCoord.X && y == kar.GameDataRes.TargetBlockCoord.Y {
					i := mathutil.MapRange(blockHealth, 0, 180, 0, 5)
					kar.CameraRes.DrawWithColorM(res.BlockCrackFrames[tileID][int(i)], kar.ColorM, kar.ColorMDIO, kar.Screen)
				} else {
					kar.CameraRes.DrawWithColorM(res.BlockCrackFrames[tileID][0], kar.ColorM, kar.ColorMDIO, kar.Screen)
				}
			}
		}
	}
	// Draw player
	if kar.ECWorld.Alive(kar.CurrentPlayer) {
		playerPos, playerSize, _, _, _, pFacing := arc.MapPlayer.Get(kar.CurrentPlayer)
		kar.ColorMDIO.GeoM.Reset()
		if pFacing.Dir.X == -1 {
			kar.ColorMDIO.GeoM.Scale(-1, 1)
			kar.ColorMDIO.GeoM.Translate(playerPos.X+playerSize.W, playerPos.Y)
		} else {
			kar.ColorMDIO.GeoM.Translate(playerPos.X, playerPos.Y)
		}
		kar.CameraRes.DrawWithColorM(kar.PlayerAnimPlayer.CurrentFrame, kar.ColorM, kar.ColorMDIO, kar.Screen)
	}

	// Draw drop Items
	itemQuery := arc.FilterDroppedItem.Query(&kar.ECWorld)
	for itemQuery.Next() {
		id, pos, timers, _, _ := itemQuery.Get()
		kar.ColorMDIO.GeoM.Reset()
		kar.ColorMDIO.GeoM.Translate(pos.X, pos.Y+kar.Sinspace[timers.AnimationIndex])
		kar.CameraRes.DrawWithColorM(res.Icon8[id.ID], kar.ColorM, kar.ColorMDIO, kar.Screen)
	}

	// Draw snowball
	q := arc.FilterProjectile.Query(&kar.ECWorld)
	for q.Next() {
		id, pos, _ := q.Get()
		kar.ColorMDIO.GeoM.Reset()
		kar.ColorMDIO.GeoM.Translate(pos.X, pos.Y)
		kar.CameraRes.DrawWithColorM(res.Icon8[id.ID], kar.ColorM, kar.ColorMDIO, kar.Screen)
	}

	// Draw target tile border
	if isRayHit {
		kar.ColorMDIO.GeoM.Reset()
		kar.ColorMDIO.GeoM.Translate(
			float64(kar.GameDataRes.TargetBlockCoord.X*kar.TileMapRes.TileW)-1,
			float64(kar.GameDataRes.TargetBlockCoord.Y*kar.TileMapRes.TileH)-1,
		)
		kar.CameraRes.DrawWithColorM(res.SelectionBlock, kar.ColorM, kar.ColorMDIO, kar.Screen)
	}

	// Draw all rects for debug
	if kar.DrawDebugHitboxesEnabled {

		itemQuery := arc.FilterDroppedItem.Query(&kar.ECWorld)
		for itemQuery.Next() {
			_, pos, _, _, _ := itemQuery.Get()
			vector.DrawFilledRect(
				kar.Screen,
				float32(pos.X),
				float32(pos.Y),
				float32(kar.GameDataRes.DropItemW),
				float32(kar.GameDataRes.DropItemH),
				color.RGBA{128, 0, 0, 10},
				false,
			)
		}

		q := arc.FilterRect.Query(&kar.ECWorld)
		for q.Next() {
			pos, size := q.Get()
			x, y := kar.CameraRes.ApplyCameraTransformToPoint(pos.X, pos.Y)
			vector.DrawFilledRect(
				kar.Screen,
				float32(x),
				float32(y),
				float32(size.W),
				float32(size.H),
				color.RGBA{128, 0, 0, 10},
				false,
			)

		}

	}
}
