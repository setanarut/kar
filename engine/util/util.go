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
	return ebiten.NewImageFromImage(LoadImageFromFS(fs, filePath))
}

func LoadImageFromFS(fs embed.FS, filePath string) image.Image {
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

func ReadPNG(filePath string) image.Image {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return image
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

// FlipVertical inverts position Y axis beetween bottom-left and top-left coordinate systems
func FlipVertical(v vec.Vec2, screenbHeight float64) vec.Vec2 {
	return vec.Vec2{v.X, screenbHeight - v.Y}
}

// PointToVec2 converts image.Point to Vec2
func PointToVec2(p image.Point) vec.Vec2 {
	return vec.Vec2{float64(p.X), float64(p.Y)}
}

// Vec2ToPoint returns Vec2 as image.Point
func Vec2ToPoint(v vec.Vec2) image.Point {
	return image.Point{int(v.X), int(v.Y)}
}

// RotateSlice rotates a generic slice in place by n positions.
//
//	nums := []int{1, 2, 3, 4, 5}
//	RotateSlice(&nums, 2)
func RotateSlice[T any](slice *[]T, n int) {
	length := len(*slice)

	if length == 0 {
		return
	}

	// Handle cases where n is greater than the slice length or negative
	n = ((n % length) + length) % length

	// Perform the in-place rotation
	reverse(slice, 0, n-1)
	reverse(slice, n, length-1)
	reverse(slice, 0, length-1)
}

// reverse is a helper function to reverse a slice in place
func reverse[T any](slice *[]T, start, end int) {
	for start < end {
		(*slice)[start], (*slice)[end] = (*slice)[end], (*slice)[start]
		start++
		end--
	}
}

// RotateSlice2 rotates a generic slice by n positions.
func RotateSlice2[T any](slice []T, n int) []T {
	length := len(slice)

	if length == 0 {
		return slice
	}

	// Handle cases where n is greater than the slice length or negative
	n = ((n % length) + length) % length

	return append(slice[n:], slice[:n]...)
}

// SubImage returns sub-image from spriteSheet image
func SubImage(spriteSheet *ebiten.Image, x, y, w, h int, vertical bool) *ebiten.Image {
	frameRect := image.Rect(x, y, x+w, y+h)
	subImage := spriteSheet.SubImage(frameRect).(*ebiten.Image)
	return subImage

}
