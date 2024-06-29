package util

import (
	"image"
	"image/color"
	"image/draw"
	"kar/engine/vec"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

func FillImage(img image.Image, c color.RGBA) {
	draw.Draw(img.(draw.Image), img.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
}
func DrawOver(src, dst image.Image) {
	draw.Draw(dst.(draw.Image), dst.Bounds(), src, image.Point{}, draw.Over)
}

func CloneImage(img image.Image) image.Image {
	copyImage := image.NewRGBA(img.Bounds())
	draw.Draw(copyImage, img.Bounds(), img, image.Point{}, draw.Src)
	return copyImage

}

func SubImage(spriteSheet *ebiten.Image, x, y, w, h int) *ebiten.Image {
	return spriteSheet.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image)
}

func SubImages(spriteSheet *ebiten.Image, x, y, w, h, subImageCount int, vertical bool) []*ebiten.Image {

	subImages := []*ebiten.Image{}
	frameRect := image.Rect(x, y, x+w, y+h)

	for i := 0; i < subImageCount; i++ {
		subImages = append(subImages, spriteSheet.SubImage(frameRect).(*ebiten.Image))
		if vertical {
			frameRect.Min.Y += h
			frameRect.Max.Y += h
		} else {
			frameRect.Min.X += w
			frameRect.Max.X += w
		}
	}
	return subImages

}
func SubImagesStd(spriteSheet *ebiten.Image, x, y, w, h, subImageCount int, vertical bool) []image.Image {

	subImages := []image.Image{}
	frameRect := image.Rect(x, y, x+w, y+h)

	for i := 0; i < subImageCount; i++ {
		subImages = append(subImages, spriteSheet.SubImage(frameRect))
		if vertical {
			frameRect.Min.Y += h
			frameRect.Max.Y += h
		} else {
			frameRect.Min.X += w
			frameRect.Max.X += w
		}
	}
	return subImages

}

func AddComponents(e *donburi.Entry, comps ...donburi.IComponentType) {
	for _, comp := range comps {
		e.AddComponent(comp)
	}
}

func EbitenImageCenterOffset(img *ebiten.Image) vec.Vec2 {
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
