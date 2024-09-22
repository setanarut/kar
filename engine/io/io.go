package io

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
