package tilemap

import "fmt"

type TileMap struct {
	Grid         [][]uint16
	W, H         int
	TileW, TileH int
}

func NewTileMap(w, h, tileW, tileH int) *TileMap {
	return &TileMap{
		Grid:  MakeGrid(w, h),
		W:     w,
		H:     h,
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

func Raycast(tm [][]uint16, x, y int, dirX, dirY int) (pos [2]int, id uint16, ok bool) {
	cursorX, cursorY := x, y
	for range 3 {
		cursorX += dirX
		cursorY += dirY
		if tm[cursorY][cursorX] != 0 {
			return [2]int{cursorX, cursorY}, tm[cursorY][cursorX], true
		}
	}
	return [2]int{cursorX, cursorY}, 0, false
}
