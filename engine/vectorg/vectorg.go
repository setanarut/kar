// vectorg is a vector drawing package for Ebitengine
package vectorg

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/vec"
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	// Transform for vector drawing (default is nil)
	Transform *ebiten.GeoM
)

func drawVerticesForUtil(dst *ebiten.Image, vs []ebiten.Vertex, is []uint16, clr color.Color, antialias bool) {
	r, g, b, a := clr.RGBA()
	for i := range vs {

		// Apply global Transfrom
		if Transform != nil {
			x, y := Transform.Apply(float64(vs[i].DstX), float64(vs[i].DstY))
			vs[i].DstX, vs[i].DstY = float32(x), float32(y)
		}

		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.ColorScaleMode = ebiten.ColorScaleModePremultipliedAlpha
	op.AntiAlias = antialias
	dst.DrawTriangles(vs, is, whiteSubImage, op)
}

func StrokeLine(dst *ebiten.Image, start, end vec.Vec2, strokeWidth float32, clr color.Color, antialias bool) {
	var path vector.Path
	path.MoveTo(float32(start.X), float32(start.Y))
	path.LineTo(float32(end.X), float32(end.Y))
	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = strokeWidth
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, strokeOp)
	drawVerticesForUtil(dst, vs, is, clr, antialias)
}

func FillRect(dst *ebiten.Image, topLeft vec.Vec2, w, h float64, clr color.Color, antialias bool) {
	x, y, width, height := float32(topLeft.X), float32(topLeft.Y), float32(w), float32(h)
	var path vector.Path
	path.MoveTo(x, y)
	path.LineTo(x, y+height)
	path.LineTo(x+width, y+height)
	path.LineTo(x+width, y)
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)

	drawVerticesForUtil(dst, vs, is, clr, antialias)
}

func StrokeRect(dst *ebiten.Image, pos vec.Vec2, w, h float64, strokeWidth float32, clr color.Color, antialias bool) {
	x, y, width, height := float32(pos.X), float32(pos.Y), float32(w), float32(h)
	var path vector.Path
	path.MoveTo(x, y)
	path.LineTo(x, y+height)
	path.LineTo(x+width, y+height)
	path.LineTo(x+width, y)
	path.Close()

	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = strokeWidth
	strokeOp.MiterLimit = 10
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, strokeOp)

	drawVerticesForUtil(dst, vs, is, clr, antialias)
}

// StrokeSquare draws centered square
func StrokeSquare(dst *ebiten.Image, center vec.Vec2, side float64, strokeWidth float32, clr color.Color, antialias bool) {
	offset := side * 0.5
	x, y, width, height := float32(center.X-offset), float32(center.Y-offset), float32(side), float32(side)

	var path vector.Path
	path.MoveTo(x, y)
	path.LineTo(x, y+height)
	path.LineTo(x+width, y+height)
	path.LineTo(x+width, y)
	path.Close()

	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = strokeWidth
	strokeOp.MiterLimit = 10
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, strokeOp)

	drawVerticesForUtil(dst, vs, is, clr, antialias)
}

func FillCircle(dst *ebiten.Image, origin vec.Vec2, r float64, clr color.Color, antialias bool) {
	var path vector.Path
	path.Arc(float32(origin.X), float32(origin.Y), float32(r), 0, 2*math.Pi, vector.Clockwise)
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)

	drawVerticesForUtil(dst, vs, is, clr, antialias)
}

func StrokeCircle(dst *ebiten.Image, origin vec.Vec2, r float64, strokeWidth float32, clr color.Color, antialias bool) {
	var path vector.Path
	path.Arc(float32(origin.X), float32(origin.Y), float32(r), 0, 2*math.Pi, vector.Clockwise)
	path.Close()

	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = strokeWidth
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, strokeOp)

	drawVerticesForUtil(dst, vs, is, clr, antialias)
}

func init() {
	b := whiteImage.Bounds()
	pix := make([]byte, 4*b.Dx()*b.Dy())
	for i := range pix {
		pix[i] = 0xff
	}
	// This is hacky, but WritePixels is better than Fill in term of automatic texture packing.
	whiteImage.WritePixels(pix)
}
