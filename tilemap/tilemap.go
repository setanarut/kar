package tilemap

import (
	"fmt"
	"image"
	"kar/items"
	"math"
)

type TileMap struct {
	Grid         [][]uint16
	W, H         int
	TileW, TileH int
}

func MakeTileMap(w, h, tileW, tileH int) *TileMap {
	return &TileMap{
		Grid:  MakeGrid(w, h),
		W:     w,
		H:     h,
		TileW: tileW,
		TileH: tileH,
	}
}

func NewTileMap(tm [][]uint16, tileW, tileH int) *TileMap {
	return &TileMap{
		Grid:  tm,
		W:     len(tm[0]),
		H:     len(tm),
		TileW: tileW,
		TileH: tileH,
	}
}

func (t *TileMap) String() string {
	s := ""
	for _, row := range t.Grid {
		for _, cell := range row {
			s += fmt.Sprintf("%d ", cell)
		}
		s += "\n"
	}
	return s
}

func MakeGrid(width, height int) [][]uint16 {
	var tm [][]uint16
	for i := 0; i < height; i++ {
		tm = append(tm, make([]uint16, width))
	}
	return tm
}

func (t *TileMap) Raycast(pos, dir image.Point, dist int) (image.Point, bool) {
	// True if exactly one of the components is non-zero
	if (dir.X != 0 && dir.Y == 0) || (dir.X == 0 && dir.Y != 0) {
		for range dist {
			pos = pos.Add(dir)
			if t.TileID(pos) != 0 {
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

func (t *TileMap) TileToWorld(pos image.Point) (float64, float64) {
	a := float64((pos.X * t.TileW) + t.TileW/2)
	b := float64((pos.Y * t.TileH) + t.TileH/2)
	return a, b
}
func (t *TileMap) TileToWorldTopLeft(pos image.Point) (float64, float64) {
	return float64(pos.X * t.TileW), float64(pos.Y * t.TileH)
}
func (t *TileMap) TileToWorld2(x, y int) (float64, float64) {
	a := float64(x * t.TileW)
	b := float64(y * t.TileH)
	return a, b
}

func (t *TileMap) SetTile(pos image.Point, id uint16) {
	if pos.X < 0 || pos.X >= t.W || pos.Y < 0 || pos.Y >= t.H {
		return
	}
	t.Grid[pos.Y][pos.X] = id
}

func (t *TileMap) TileID(pos image.Point) uint16 {
	if pos.X < 0 || pos.X >= t.W || pos.Y < 0 || pos.Y >= t.H {
		return 0
	}
	return t.Grid[pos.Y][pos.X]
}
func (t *TileMap) Get(x, y int) uint16 {
	return t.Grid[y][x]
}
func (t *TileMap) Set(x, y int, v uint16) {
	t.Grid[y][x] = v
}

func (t *TileMap) GetTileRect(pos image.Point) (x, y, w, h float64) {
	return float64(pos.X * t.TileW), float64(pos.Y * t.TileH), float64(t.TileW), float64(t.TileH)
}

func (t *TileMap) FindSpawnPosition() (px, py float64) {
	x := 20 * 20
	for y := range t.H - 1 {
		upperTile := t.Get(x, y)
		downTile := t.Get(x, y+1)
		if downTile != items.Air && upperTile == items.Air {
			px, py = t.TileToWorld2(x, y-1)
			break
		}
	}
	return px, py

}
