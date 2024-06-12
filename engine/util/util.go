package util

import (
	"image"
	"image/draw"
	"kar/engine/vec"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

func DrawOver(src, dst image.Image) {
	draw.Draw(dst.(draw.Image), dst.Bounds(), src, image.Point{0, 0}, draw.Over)
}

func CloneImage(img image.Image) image.Image {
	copyImage := image.NewRGBA(img.Bounds())
	draw.Draw(copyImage, img.Bounds(), img, image.Point{0, 0}, draw.Src)
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
