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

// index kontrol fonksiyonu, sadece bool değeri döndürüyor
func CheckIndex[T any](slice []T, index int) bool {
	if index < 0 || index >= len(slice) {
		return false
	}
	return true
}

func ImageCenterOffset(img image.Image) vec.Vec2 {
	o := vec.Vec2{float64(img.Bounds().Dx()), float64(img.Bounds().Dy())}
	return o.Scale(0.5).Neg()
}

// HexToRGBA converts hex color to color.RGBA with "#FFFFFF" format
func HexToRGBA(hex string) color.RGBA {
	values, _ := strconv.ParseUint(string(hex[1:]), 16, 32)
	return color.RGBA{
		R: uint8(values >> 16),
		G: uint8((values >> 8) & 0xFF),
		B: uint8(values & 0xFF),
		A: 255,
	}
}

func ReadEbImgFS(fs embed.FS, filePath string) *ebiten.Image {
	return ebiten.NewImageFromImage(ImgFromFS(fs, filePath))
}

func ImgFromFS(fs embed.FS, filePath string) image.Image {
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

func LoadFontFromFS(file string, size float64, fs embed.FS) *text.GoTextFace {
	f, err := fs.ReadFile(file)
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
