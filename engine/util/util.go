package util

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/setanarut/vec"
	"golang.org/x/text/language"
)

func RectPoints(ori vec.Vec2, side float64) []vec.Vec2 {
	ori.X += -(side / 2)
	ori.Y += -(side / 2)
	points := []vec.Vec2{{0, 0}, {side, 0}, {side, side}, {0, side}}
	for i, p := range points {
		points[i] = p.Add(ori)
	}
	return points
}

func ImageCenterOffset(img image.Image) vec.Vec2 {
	return vec.Vec2{float64(img.Bounds().Dx()), float64(img.Bounds().Dy())}.Scale(0.5).Neg()
}

func UnpackPoint(p image.Point) (int, int) {
	return p.X, p.Y
}
func UnpackVec2(v vec.Vec2) (float64, float64) {
	return v.X, v.Y
}

// HexToRGBA converts hex color to color.RGBA with "#FFFFFF" format
func HexToRGBA(hex string) color.RGBA {
	values, _ := strconv.ParseUint(string(hex[1:]), 16, 32)
	return color.RGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}

func LoadEbitenImageFromFS(fs embed.FS, filePath string) *ebiten.Image {
	return ebiten.NewImageFromImage(LoadStandartImageFromFS(fs, filePath))
}

func LoadStandartImageFromFS(fs embed.FS, filePath string) image.Image {
	f, err := fs.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return image
}

func WriteImage(img image.Image, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	png.Encode(file, img)
}

func LoadGoTextFaceFromFS(fileName string, size float64, assets embed.FS) *text.GoTextFace {
	f, err := assets.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	src, err := text.NewGoTextFaceSource(bytes.NewReader(f))
	if err != nil {
		log.Fatal(err)
	}
	gtf := &text.GoTextFace{
		Source:   src,
		Size:     size,
		Language: language.English,
	}
	return gtf
}
