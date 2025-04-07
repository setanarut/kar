package tilemap

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"kar/items"
	"log"
	"math"

	"github.com/setanarut/v"
)

var (
	Up    = image.Point{0, -1} // {0, -1}
	Down  = image.Point{0, 1}  // {0, 1}
	Left  = image.Point{-1, 0} // {-1, 0}
	Right = image.Point{1, 0}  // {1, 0}
)

type TileMap struct {
	Grid         [][]uint8
	W, H         int
	TileW, TileH int
}

func MakeTileMap(w, h, tileW, tileH int) TileMap {
	return TileMap{
		Grid:  MakeGrid(w, h),
		W:     w,
		H:     h,
		TileW: tileW,
		TileH: tileH,
	}
}

func MakeGrid(width, height int) [][]uint8 {
	var tm [][]uint8
	for range height {
		tm = append(tm, make([]uint8, width))
	}
	return tm
}

func (t *TileMap) Raycast(pos image.Point, dirX, dirY, dist int) (image.Point, bool) {
	// True if exactly one of the components is non-zero
	if (dirX != 0 && dirY == 0) || (dirX == 0 && dirY != 0) {
		for range dist {
			pos.X += dirX
			pos.Y += dirY
			if t.GetID(pos.X, pos.Y) != items.Air {
				return pos, true
			}
		}
	} else {
		return image.Point{}, false
	}
	return image.Point{}, false
}

func (t *TileMap) WorldToTile(x, y float64) image.Point {
	return image.Point{int(math.Floor(x / float64(t.TileW))), int(math.Floor(y / float64(t.TileH)))}
}
func (t *TileMap) WorldToTile2(x, y float64) (int, int) {
	return int(math.Floor(x / float64(t.TileW))), int(math.Floor(y / float64(t.TileH)))
}

func (t *TileMap) FloorToBlockCenter(x, y float64) v.Vec {
	return t.TileToWorld(t.WorldToTile(x, y))
}

// Tile coords to block center
func (t *TileMap) TileToWorld(p image.Point) v.Vec {
	return v.Vec{float64((p.X * t.TileW) + t.TileW/2), float64((p.Y * t.TileH) + t.TileH/2)}
}

func (t *TileMap) GetID(x, y int) uint8 {
	if x < 0 || x >= t.W || y < 0 || y >= t.H {
		return 0
	}
	return t.Grid[y][x]
}

func (t *TileMap) GetIDUnchecked(coords image.Point) uint8 {
	return t.Grid[coords.Y][coords.X]
}

func (t *TileMap) TileIDProperty(x, y int) items.ItemProperty {
	return items.Property[t.GetID(x, y)]
}

func (t *TileMap) Set(x, y int, id uint8) {
	if x < 0 || x >= t.W || y < 0 || y >= t.H {
		return
	}
	t.Grid[y][x] = id
}
func (t *TileMap) SetUnchecked(coord image.Point, id uint8) {
	t.Grid[coord.Y][coord.X] = id
}

func (t *TileMap) GetTileRect(x, y int) (rectX, rectY, rectW, rectH float64) {
	return float64(x * t.TileW), float64(y * t.TileH), float64(t.TileW), float64(t.TileH)
}

func (t *TileMap) GetTileBottom(x, y int) float64 {
	return float64((y + 1) * t.TileH)
}

func (t *TileMap) FindSpawnPosition() image.Point {
	x := 20 * 20
	for y := range t.H - 1 {
		upperTile := t.GetID(x, y)
		downTile := t.GetID(x, y+1)
		if downTile != items.Air && upperTile == items.Air {
			return image.Point{x, y}
		}
	}
	return image.Point{}
}
func (t *TileMap) FindSpawnPosition2(x, y int) image.Point {
	for i := range 10 {
		if t.GetID(x, y+i) == items.Air && t.GetID(x, (y+i)+1) != items.Air {
			return image.Point{x, y + i}
		}
	}

	return image.Point{}
}

func (t *TileMap) GetImageByte() []byte {
	buf := &bytes.Buffer{}
	if err := png.Encode(buf, t.GetImage()); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func (tm *TileMap) GetImage() *image.RGBA {
	im := image.NewRGBA(image.Rect(0, 0, tm.W, tm.H))
	for y := range tm.H {
		for x := range tm.W {
			id := tm.Grid[y][x]
			v, ok := items.ColorMap[id]
			if ok {
				im.Set(x, y, v)
			} else {
				im.Set(x, y, color.Black)
			}
		}
	}
	return im
}
