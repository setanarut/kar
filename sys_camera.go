package kar

import (
	"image/color"
	"kar/items"
	"kar/res"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
)

type Camera struct{}

func (c *Camera) Init() {}

func (c *Camera) Update() {

	if world.Alive(currentPlayer) {
		playerAABB := MapAABB.Get(currentPlayer)
		// Toggle camera follow
		if inpututil.IsKeyJustPressed(ebiten.KeyL) {
			switch cameraRes.SmoothType {
			case kamera.Lerp:
				cameraRes.SetCenter(playerAABB.Pos.X, playerAABB.Pos.Y)
				cameraRes.SmoothType = kamera.SmoothDamp
			case kamera.SmoothDamp:
				cameraRes.SetCenter(playerAABB.Pos.X, playerAABB.Pos.Y)
				cameraRes.SmoothType = kamera.Lerp
			}
		}

		// Camera follow
		if MapHealth.Get(currentPlayer).Current > 0 {
			if cameraRes.SmoothType == kamera.Lerp {
				// if playerCenterX < CameraRes.X {
				// 	CameraRes.X -= CameraRes.Width
				// }
				// if playerCenterX > CameraRes.Right() {
				// 	CameraRes.X += CameraRes.Width
				// }
				if playerAABB.Pos.Y < cameraRes.Y {
					cameraRes.SetTopLeft(cameraRes.X, cameraRes.Y-cameraRes.Height)
				}
				if playerAABB.Pos.Y > cameraRes.Bottom() {
					cameraRes.SetTopLeft(cameraRes.X, cameraRes.Y+cameraRes.Height)
				}
				cameraRes.LookAt(math.Floor(playerAABB.Pos.X), math.Floor(playerAABB.Pos.Y))
			} else if cameraRes.SmoothType == kamera.SmoothDamp {
				cameraRes.LookAt(math.Floor(playerAABB.Pos.X), math.Floor(playerAABB.Pos.Y))
			}
		}
	}
}
func (c *Camera) Draw() {

	// DRAW TILEMAPs

	// clamp tilemap bounds
	camMin := tileMapRes.WorldToTile(cameraRes.X, cameraRes.Y)
	camMin.X = min(max(camMin.X, 0), tileMapRes.W)
	camMin.Y = min(max(camMin.Y, 0), tileMapRes.H)
	camMaxX := min(max(camMin.X+renderArea.X, 0), tileMapRes.W)
	camMaxY := min(max(camMin.Y+renderArea.Y, 0), tileMapRes.H)
	// draw tiles
	for y := camMin.Y; y < camMaxY; y++ {
		for x := camMin.X; x < camMaxX; x++ {
			tileID := tileMapRes.Grid[y][x]
			if tileID != 0 {
				px, py := float64(x*tileMapRes.TileW), float64(y*tileMapRes.TileH)
				ColorMDIO.GeoM.Reset()

				if x == ceilBlockCoord.X && y == ceilBlockCoord.Y {
					if tileID == items.Bedrock {
						if ceilBlockTick > 0 {
							ceilBlockTick -= 0.1
						}
						py -= ceilBlockTick
					}
				}

				ColorMDIO.GeoM.Translate(px, py)
				if items.HasTag(tileID, items.UnbreakableBlock) {
					cameraRes.DrawWithColorM(res.BlockUnbreakable[tileID], ColorM, ColorMDIO, Screen)
				} else {
					if x == gameDataRes.TargetBlockCoord.X && y == gameDataRes.TargetBlockCoord.Y {
						i := MapRange(gameDataRes.BlockHealth, 0, 180, 0, 5)
						cameraRes.DrawWithColorM(res.BlockCrackFrames[tileID][int(i)], ColorM, ColorMDIO, Screen)
					} else {
						cameraRes.DrawWithColorM(res.BlockCrackFrames[tileID][0], ColorM, ColorMDIO, Screen)
					}
				}
			}
		}
	}
	// Draw player
	if world.Alive(currentPlayer) {
		playerBox, _, _, _, pFacing := MapPlayer.Get(currentPlayer)
		ColorMDIO.GeoM.Reset()
		x := playerBox.Pos.X - playerBox.Half.X
		y := playerBox.Pos.Y - playerBox.Half.Y
		if pFacing.X == -1 {
			ColorMDIO.GeoM.Scale(-1, 1)
			ColorMDIO.GeoM.Translate(playerBox.Pos.X+playerBox.Half.X, y)
		} else {
			ColorMDIO.GeoM.Translate(x, y)
		}
		cameraRes.DrawWithColorM(animPlayer.CurrentFrame, ColorM, ColorMDIO, Screen)
		if DrawItemHitboxEnabled {
			x, y = cameraRes.ApplyCameraTransformToPoint(x, y)
			vector.DrawFilledRect(
				Screen,
				float32(x),
				float32(y),
				float32(playerBox.Half.X*2),
				float32(playerBox.Half.Y*2),
				color.RGBA{128, 0, 0, 10},
				false,
			)
		}
	}

	// Draw drop Items

	itemQuery := FilterDroppedItem.Query()
	for itemQuery.Next() {
		id, pos, animIndex := itemQuery.Get()
		ColorMDIO.GeoM.Reset()
		x := pos.X - dropItemHalfSize.X
		y := pos.Y - dropItemHalfSize.Y
		siny := y + Sinspace[animIndex.Index]
		ColorMDIO.GeoM.Translate(x, siny)
		if id.ID != items.Air {

			cameraRes.DrawWithColorM(res.Icon8[id.ID], ColorM, ColorMDIO, Screen)

			if DrawItemHitboxEnabled {
				x, y = cameraRes.ApplyCameraTransformToPoint(x, y)
				vector.DrawFilledRect(
					Screen,
					float32(x),
					float32(y),
					float32(dropItemHalfSize.X*2),
					float32(dropItemHalfSize.Y*2),
					color.RGBA{128, 0, 0, 10},
					false,
				)
			}
		}
	}

	// Draw snowball
	q := FilterProjectile.Query()
	for q.Next() {
		id, pos, _ := q.Get()
		ColorMDIO.GeoM.Reset()
		ColorMDIO.GeoM.Translate(pos.X, pos.Y)
		cameraRes.DrawWithColorM(res.Icon8[id.ID], ColorM, ColorMDIO, Screen)
	}

	// Draw target tile border
	if gameDataRes.IsRayHit {
		ColorMDIO.GeoM.Reset()
		ColorMDIO.GeoM.Translate(
			float64(gameDataRes.TargetBlockCoord.X*tileMapRes.TileW)-1,
			float64(gameDataRes.TargetBlockCoord.Y*tileMapRes.TileH)-1,
		)
		cameraRes.DrawWithColorM(res.BlockBorder, ColorM, ColorMDIO, Screen)
	}

}
