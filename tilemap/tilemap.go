package tilemap

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"kar/items"
	"log"
	"math"
	"os"
)

type TileMap struct {
	Grid         [][]uint8
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

func NewTileMap(tm [][]uint8, tileW, tileH int) *TileMap {
	return &TileMap{
		Grid:  tm,
		W:     len(tm[0]),
		H:     len(tm),
		TileW: tileW,
		TileH: tileH,
	}
}

func (tm *TileMap) WriteAsImage(path string, playerX, playerY int) {
	im := tm.GetImage()
	im.Set(playerX, playerY, color.RGBA{255, 0, 255, 255})
	WritePNG(im, path)
}

func (tm *TileMap) CloneEmpty() *TileMap {
	return MakeTileMap(tm.W, tm.H, tm.TileW, tm.TileH)
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

func MakeGrid(width, height int) [][]uint8 {
	var tm [][]uint8
	for i := 0; i < height; i++ {
		tm = append(tm, make([]uint8, width))
	}
	return tm
}

func (t *TileMap) Raycast(pos, dir image.Point, dist int) (image.Point, bool) {
	// True if exactly one of the components is non-zero
	if (dir.X != 0 && dir.Y == 0) || (dir.X == 0 && dir.Y != 0) {
		for range dist {
			pos = pos.Add(dir)
			if t.Get(pos.X, pos.Y) != items.Air {
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

func (t *TileMap) FloorToBlockCenter(x, y float64) (float64, float64) {
	p := t.WorldToTile(x, y)
	return t.TileToWorldCenter(p.X, p.Y)
}

// Tile coords to block center
func (t *TileMap) TileToWorldCenter(x, y int) (float64, float64) {
	a := float64((x * t.TileW) + t.TileW/2)
	b := float64((y * t.TileH) + t.TileH/2)
	return a, b
}

func (t *TileMap) Get(x, y int) uint8 {
	if x < 0 || x >= t.W || y < 0 || y >= t.H {
		return 0
	}
	return t.Grid[y][x]
}

func (t *TileMap) GetUnchecked(coords image.Point) uint8 {
	return t.Grid[coords.Y][coords.X]
}

func (t *TileMap) TileIDProperty(x, y int) items.ItemProperty {
	return items.Property[t.Get(x, y)]
}

func (t *TileMap) Set(x, y int, id uint8) {
	if x < 0 || x >= t.W || y < 0 || y >= t.H {
		return
	}
	t.Grid[y][x] = id
}

func (t *TileMap) GetTileRect(x, y int) (rectX, rectY, rectW, rectH float64) {
	return float64(x * t.TileW), float64(y * t.TileH), float64(t.TileW), float64(t.TileH)
}

// func (t *TileMap) GetTileRect2(x, y int) (Position, Size) {
// 	return float64(x * t.TileW), float64(y * t.TileH), float64(t.TileW), float64(t.TileH)
// }

func (t *TileMap) FindSpawnPosition() (px, py int) {
	x := 20 * 20
	for y := range t.H - 1 {
		upperTile := t.Get(x, y)
		downTile := t.Get(x, y+1)
		if downTile != items.Air && upperTile == items.Air {
			// px, py = t.TileToWorldCenter(x, y-1)
			px, py = x, y-1
			break
		}
	}
	return px, py
}

// Veriyi diske yazan fonksiyon (hata durumunda log kullanılıyor)
func (t *TileMap) WriteToDisk(filename string) {
	// Veriyi byte array'e dönüştürmek için gob kullanıyoruz
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(t.Grid); err != nil {
		log.Fatalf("Veriyi encode ederken hata: %v", err)
		return
	}

	// Dosya oluştur ve yaz
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Dosya oluşturulamadı: %v", err)
		return
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		log.Fatalf("Dosyaya yazarken hata: %v", err)
		return
	}
}

// Diske yazılmış veriyi okuyan fonksiyon (hata durumunda log kullanılıyor)
func (t *TileMap) ReadFromDisk(filename string) [][]uint8 {
	// Dosyayı aç
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Dosya açılamadı: %v", err)
		return nil
	}
	defer file.Close()

	// Dosya içeriğini oku
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Dosya bilgisi alınamadı: %v", err)
		return nil
	}
	data := make([]byte, fileInfo.Size())
	_, err = file.Read(data)
	if err != nil {
		log.Fatalf("Dosyadan okuma hatası: %v", err)
		return nil
	}

	// Byte array'i geri çözerek [][]uint8'ya dönüştür
	var result [][]uint8
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&result); err != nil {
		log.Fatalf("Veriyi decode ederken hata: %v", err)
		return nil
	}

	return result
}

func WritePNG(im image.Image, name string) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, im); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
