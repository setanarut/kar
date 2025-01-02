package util

import (
	"bytes"
	"embed"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
)

// index kontrol fonksiyonu, sadece bool değeri döndürüyor
func CheckIndex[T any](slice []T, index int) bool {
	if index < 0 || index >= len(slice) {
		return false
	}
	return true
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
