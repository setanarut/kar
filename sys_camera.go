package kar

import (
	"kar/items"
	"kar/res"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/kamera/v2"
)

type Camera struct{}

func (c *Camera) Init() {
	cameraRes.SmoothOptions.LerpSpeedX = 0.4
	cameraRes.SmoothOptions.LerpSpeedY = 0.0
}
func (c *Camera) Update() {
	if world.Alive(currentPlayer) {
		playerAABB := mapAABB.GetUnchecked(currentPlayer)
		// Toggle camera follow
		if inpututil.IsKeyJustPressed(ebiten.KeyL) {
			switch cameraRes.SmoothType {
			case kamera.Lerp:
				cameraRes.SetCenter(playerAABB.Pos.X, playerAABB.Pos.Y)
				cameraRes.SmoothType = kamera.SmoothDamp
			case kamera.SmoothDamp:
				cameraRes.SetCenter(playerAABB.Pos.X, playerAABB.Pos.Y)
				cameraRes.SmoothType = kamera.None
			case kamera.None:
				cameraRes.SetCenter(playerAABB.Pos.X, playerAABB.Pos.Y)
				cameraRes.SmoothType = kamera.Lerp
			}
		}
		// Camera follow
		if mapHealth.GetUnchecked(currentPlayer).Current > 0 {
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
				cameraRes.LookAt(playerAABB.Pos.X, playerAABB.Pos.Y)
				cameraRes.X = math.Floor(cameraRes.X)
				cameraRes.Y = math.Floor(cameraRes.Y)

			} else if cameraRes.SmoothType == kamera.SmoothDamp {
				cameraRes.LookAt(playerAABB.Pos.X, playerAABB.Pos.Y)
				// cameraRes.TempTargetX = math.Floor(cameraRes.TempTargetX)
				// cameraRes.TempTargetY = math.Floor(cameraRes.TempTargetY)
				cameraRes.X = math.Floor(cameraRes.X)
				cameraRes.Y = math.Floor(cameraRes.Y)
			} else if cameraRes.SmoothType == kamera.None {
				cameraRes.SetCenter(playerAABB.Pos.X, playerAABB.Pos.Y)
			}
		}
	}
}
func (c *Camera) Draw() {

	// DRAW TILEMAP

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
				if x == ceilBlockCoord.X && y == ceilBlockCoord.Y {
					if tileID == items.Bedrock {
						if ceilBlockTick > 0 {
							ceilBlockTick -= 0.1
						}
						py -= ceilBlockTick
					}
				}

				colorMDIO.GeoM.Reset()
				colorMDIO.GeoM.Translate(px, py)

				if items.HasTag(tileID, items.UnbreakableBlock) {
					cameraRes.DrawWithColorM(res.BlockUnbreakable[tileID], colorM, colorMDIO, Screen)
				} else {
					if x == gameDataRes.TargetBlockCoord.X && y == gameDataRes.TargetBlockCoord.Y {
						i := MapRange(gameDataRes.BlockHealth, 0, 180, 0, 5)
						cameraRes.DrawWithColorM(res.BlockCrackFrames[tileID][int(i)], colorM, colorMDIO, Screen)
					} else {
						cameraRes.DrawWithColorM(res.BlockCrackFrames[tileID][0], colorM, colorMDIO, Screen)
					}
				}
			}
		}
	}

	// Draw target tile border
	if gameDataRes.IsRayHit {
		colorMDIO.GeoM.Reset()
		colorMDIO.GeoM.Translate(
			float64(gameDataRes.TargetBlockCoord.X*tileMapRes.TileW)-1,
			float64(gameDataRes.TargetBlockCoord.Y*tileMapRes.TileH)-1,
		)
		cameraRes.DrawWithColorM(res.BlockBorder, colorM, colorMDIO, Screen)
	}

}
