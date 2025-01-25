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

	// Toggle camera follow
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		switch kar.Camera.SmoothType {
		case kamera.None:
			kar.Camera.SetCenter(playerCenterX, playerCenterY-40)
			kar.Camera.SmoothType = kamera.SmoothDamp
		case kamera.SmoothDamp:
			kar.Camera.SetCenter(playerCenterX, playerCenterY-40)
			kar.Camera.SmoothType = kamera.None
		}
	}

	// Static camera logic
	if kar.Camera.SmoothType == kamera.None {
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
	} else {
		kar.Camera.LookAt(math.Floor(playerCenterX), math.Floor(playerCenterY-40))
	}

	// update animation players
	if !craftingState {
		q := arc.FilterAnimPlayer.Query(&kar.WorldECS)
		for q.Next() {
			a := q.Get()
			a.Update()
		}

	}
}

func (d *Game) Draw() {

	// DRAW TILEMAPs

	// clamp tilemap bounds
	camMin := tileMap.WorldToTile(kar.Camera.TopLeft())
	camMin.X = min(max(camMin.X, 0), tileMap.W)
	camMin.Y = min(max(camMin.Y, 0), tileMap.H)
	camMaxX := min(max(camMin.X+kar.RenderArea.X, 0), tileMap.W)
	camMaxY := min(max(camMin.Y+kar.RenderArea.Y, 0), tileMap.H)
	// draw tiles
	for y := camMin.Y; y < camMaxY; y++ {
		for x := camMin.X; x < camMaxX; x++ {
			tileID := tileMap.Grid[y][x]
			if tileID != 0 {
				px, py := float64(x*tileMap.TileW), float64(y*tileMap.TileH)
				kar.ColorMDIO.GeoM.Reset()
				kar.ColorMDIO.GeoM.Translate(px, py)
				if x == targetTile.X && y == targetTile.Y {
					i := mathutil.MapRange(blockHealth, 0, 180, 0, 5)
					kar.Camera.DrawWithColorM(res.BlockCrackFrames[tileID][int(i)], kar.ColorM, kar.ColorMDIO, kar.Screen)
				} else {
					kar.Camera.DrawWithColorM(res.BlockCrackFrames[tileID][0], kar.ColorM, kar.ColorMDIO, kar.Screen)
				}
			}
		}
	}

	// Draw Items
	itemQuery := arc.FilterItem.Query(&kar.WorldECS)
	for itemQuery.Next() {
		id, rect, timers, _ := itemQuery.Get()
		kar.ColorMDIO.GeoM.Reset()
		kar.ColorMDIO.GeoM.Translate(rect.X, rect.Y+sinspace[timers.AnimationIndex])
		kar.Camera.DrawWithColorM(res.Icon8[id.ID], kar.ColorM, kar.ColorMDIO, kar.Screen)
	}

	// Draw snowball
	q := arc.FilterMapSnowBall.Query(&kar.WorldECS)
	for q.Next() {
		id, rect, _ := q.Get()
		kar.ColorMDIO.GeoM.Reset()
		kar.ColorMDIO.GeoM.Translate(rect.X, rect.Y)
		kar.Camera.DrawWithColorM(res.Icon8[id.ID], kar.ColorM, kar.ColorMDIO, kar.Screen)
	}

	// Draw target tile border
	if isRayHit {
		kar.ColorMDIO.GeoM.Reset()
		kar.ColorMDIO.GeoM.Translate(
			float64(targetTile.X*tileMap.TileW)-1,
			float64(targetTile.Y*tileMap.TileH)-1,
		)
		kar.Camera.DrawWithColorM(res.SelectionBlock, kar.ColorM, kar.ColorMDIO, kar.Screen)
	}

	// Draw all rects for debug
	if kar.DrawDebugHitboxesEnabled {
		rectQ := arc.FilterRect.Query(&kar.WorldECS)
		for rectQ.Next() {
			rect := rectQ.Get()
			x, y := kar.Camera.ApplyCameraTransformToPoint(rect.X, rect.Y)
			vector.DrawFilledRect(
				kar.Screen,
				float32(x),
				float32(y),
				float32(rect.W),
				float32(rect.H),
				color.RGBA{128, 0, 0, 10},
				false,
			)
		}

		// Draw player tile for debug
		x, y, w, h := tileMap.GetTileRect(playerTile.X, playerTile.Y)
		x, y = kar.Camera.ApplyCameraTransformToPoint(x, y)
		vector.DrawFilledRect(
			kar.Screen,
			float32(x),
			float32(y),
			float32(w),
			float32(h),
			color.RGBA{0, 0, 128, 10},
			false,
		)
	}
}
