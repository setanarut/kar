package kar

import (
	"image/color"
	"kar/res"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
)

var CeilBlockCoord [2]int
var CeilBlockTick float64

type Camera struct{}

func (c *Camera) Init() {

}

func (c *Camera) Update() error {

	if ECWorld.Alive(CurrentPlayer) {
		playerPos, playerSize := MapRect.Get(CurrentPlayer)
		playerCenterX, playerCenterY := playerPos.X+playerSize.W, playerPos.Y+playerSize.H

		// Save TileMap image to desktop
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			playerPos := MapPosition.Get(CurrentPlayer)
			playerTile := TileMapRes.WorldToTile(playerPos.X+8, playerPos.Y+8)
			TileMapRes.WriteAsImage(DesktopPath+"map.png", playerTile.X, playerTile.Y)
		}

		// Toggle camera follow
		if inpututil.IsKeyJustPressed(ebiten.KeyL) {
			switch CameraRes.SmoothType {
			case kamera.None:
				CameraRes.SetCenter(playerCenterX, playerCenterY)
				CameraRes.SmoothType = kamera.Lerp
			case kamera.SmoothDamp, kamera.Lerp:
				CameraRes.SetCenter(playerCenterX, playerCenterY)
				CameraRes.SmoothType = kamera.None
			}
		}

		// Camera follow
		if MapHealth.Get(CurrentPlayer).Current > 0 {
			if CameraRes.SmoothType == kamera.Lerp {
				// if playerCenterX < CameraRes.X {
				// 	CameraRes.X -= CameraRes.Width
				// }
				// if playerCenterX > CameraRes.Right() {
				// 	CameraRes.X += CameraRes.Width
				// }
				if playerCenterY < CameraRes.Y {
					CameraRes.SetTopLeft(CameraRes.X, CameraRes.Y-CameraRes.Height)
				}
				if playerCenterY > CameraRes.Bottom() {
					CameraRes.SetTopLeft(CameraRes.X, CameraRes.Y+CameraRes.Height)
				}
				CameraRes.LookAt(math.Floor(playerCenterX), math.Floor(playerCenterY))
			}
		}
	}
	return nil
}

func (c *Camera) Draw() {

	// DRAW TILEMAPs

	// clamp tilemap bounds
	camMin := TileMapRes.WorldToTile(CameraRes.X, CameraRes.Y)
	camMin.X = min(max(camMin.X, 0), TileMapRes.W)
	camMin.Y = min(max(camMin.Y, 0), TileMapRes.H)
	camMaxX := min(max(camMin.X+RenderArea.X, 0), TileMapRes.W)
	camMaxY := min(max(camMin.Y+RenderArea.Y, 0), TileMapRes.H)
	// draw tiles
	for y := camMin.Y; y < camMaxY; y++ {
		for x := camMin.X; x < camMaxX; x++ {
			tileID := TileMapRes.Grid[y][x]
			if tileID != 0 {
				px, py := float64(x*TileMapRes.TileW), float64(y*TileMapRes.TileH)
				ColorMDIO.GeoM.Reset()

				if x == CeilBlockCoord[0] && y == CeilBlockCoord[1] {
					if CeilBlockTick > 0 {
						CeilBlockTick -= 0.1
					}
					py -= CeilBlockTick
				}

				ColorMDIO.GeoM.Translate(px, py)
				if x == GameDataRes.TargetBlockCoord.X && y == GameDataRes.TargetBlockCoord.Y {
					i := MapRange(blockHealth, 0, 180, 0, 5)
					CameraRes.DrawWithColorM(res.BlockCrackFrames[tileID][int(i)], ColorM, ColorMDIO, Screen)
				} else {
					CameraRes.DrawWithColorM(res.BlockCrackFrames[tileID][0], ColorM, ColorMDIO, Screen)
				}
			}
		}
	}
	// Draw player
	if ECWorld.Alive(CurrentPlayer) {
		playerPos, playerSize, _, _, _, pFacing := MapPlayer.Get(CurrentPlayer)
		ColorMDIO.GeoM.Reset()
		if pFacing.Dir.X == -1 {
			ColorMDIO.GeoM.Scale(-1, 1)
			ColorMDIO.GeoM.Translate(playerPos.X+playerSize.W, playerPos.Y)
		} else {
			ColorMDIO.GeoM.Translate(playerPos.X, playerPos.Y)
		}
		CameraRes.DrawWithColorM(PlayerAnimPlayer.CurrentFrame, ColorM, ColorMDIO, Screen)
	}

	// Draw drop Items
	itemQuery := FilterDroppedItem.Query(&ECWorld)
	for itemQuery.Next() {
		id, pos, animIndex, _, _ := itemQuery.Get()
		ColorMDIO.GeoM.Reset()
		ColorMDIO.GeoM.Translate(pos.X, pos.Y+Sinspace[animIndex.Index])
		CameraRes.DrawWithColorM(res.Icon8[id.ID], ColorM, ColorMDIO, Screen)
	}

	// Draw snowball
	q := FilterProjectile.Query(&ECWorld)
	for q.Next() {
		id, pos, _ := q.Get()
		ColorMDIO.GeoM.Reset()
		ColorMDIO.GeoM.Translate(pos.X, pos.Y)
		CameraRes.DrawWithColorM(res.Icon8[id.ID], ColorM, ColorMDIO, Screen)
	}

	// Draw target tile border
	if IsRayHit {
		ColorMDIO.GeoM.Reset()
		ColorMDIO.GeoM.Translate(
			float64(GameDataRes.TargetBlockCoord.X*TileMapRes.TileW)-1,
			float64(GameDataRes.TargetBlockCoord.Y*TileMapRes.TileH)-1,
		)
		CameraRes.DrawWithColorM(res.BlockBorder, ColorM, ColorMDIO, Screen)
	}

	// Draw all rects for debug
	if DrawDebugHitboxesEnabled {
		itemQuery := FilterDroppedItem.Query(&ECWorld)
		for itemQuery.Next() {
			_, pos, _, _, _ := itemQuery.Get()
			vector.DrawFilledRect(
				Screen,
				float32(pos.X),
				float32(pos.Y),
				float32(DropItemSize.W),
				float32(DropItemSize.H),
				color.RGBA{128, 0, 0, 10},
				false,
			)
		}

		q := FilterRect.Query(&ECWorld)
		for q.Next() {
			pos, size := q.Get()
			x, y := CameraRes.ApplyCameraTransformToPoint(pos.X, pos.Y)
			vector.DrawFilledRect(
				Screen,
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
