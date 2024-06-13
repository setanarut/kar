package vectordraw

import (
	"image"
	"image/color"
	"kar/engine/cm"
	"kar/engine/vec"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Vector drawing vars
var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	vertices      []ebiten.Vertex
	indices       []uint16
	dto           = &ebiten.DrawTrianglesOptions{AntiAlias: true}
	so            = &vector.StrokeOptions{
		Width:    5,
		LineCap:  vector.LineCapRound,
		LineJoin: vector.LineJoinRound,
	}
)

func init() {
	whiteImage.Fill(color.White)

}

func DrawChipmunkShape(screen *ebiten.Image, s *cm.Shape, clr color.Color, screenHeight float64) {
	switch s.Class.(type) {
	case *cm.Circle:
		StrokeCircle(screen, s.Class.(*cm.Circle).Radius(), 1, s.Body().Position().FlipVertical(screenHeight), clr)
		// FillCircle(screen, s.Class.(*cm.Circle).Radius(), InvPosVectY(s.Body().Position(), screenHeight), c)
	case *cm.Segment:
		r := s.Class.(*cm.Segment).Radius()
		a := s.Class.(*cm.Segment).TransformA().FlipVertical(screenHeight)
		b := s.Class.(*cm.Segment).TransformB().FlipVertical(screenHeight)
		if r < 1 {
			DrawLine(screen, a, b, 1, clr)
		} else {
			DrawLine2(screen, a, b, r*2)
		}
	case *cm.PolyShape:
		DrawChipmunkBB(screen, s.BB(), screenHeight, clr)
	}
}
func DrawChipmunkShapeGEOM(screen *ebiten.Image, s *cm.Shape, clr color.Color, screenHeight float64, geom *ebiten.GeoM) {

	switch s.Class.(type) {

	case *cm.Circle:
		pos := ApplyGeoM2Vec2(s.Body().Position().FlipVertical(screenHeight), geom)
		StrokeCircle(screen, s.Class.(*cm.Circle).Radius(), 1, pos, clr)
		// FillCircle(screen, s.Class.(*cm.Circle).Radius(), InvPosVectY(pos, screenHeight), c)

	case *cm.Segment:
		r := s.Class.(*cm.Segment).Radius()
		a := ApplyGeoM2Vec2(s.Class.(*cm.Segment).TransformA().FlipVertical(screenHeight), geom)
		b := ApplyGeoM2Vec2(s.Class.(*cm.Segment).TransformB().FlipVertical(screenHeight), geom)

		if r < 1 {
			DrawLine(screen, a, b, 1, clr)
		} else {
			DrawLine2(screen, a, b, r*2)
		}

	case *cm.PolyShape:
		DrawChipmunkBB(screen, s.BB(), screenHeight, clr)
	}
}

func DrawLine(screen *ebiten.Image, a, b vec.Vec2, strokeWidth float64, c color.Color) {
	vector.StrokeLine(
		screen,
		float32(a.X),
		float32(a.Y),
		float32(b.X),
		float32(b.Y),
		float32(strokeWidth),
		c,
		true)

}

func DrawLine2(screen *ebiten.Image, a, b vec.Vec2, strokeWidth float64) {
	var path vector.Path
	so.Width = float32(strokeWidth)
	path.MoveTo(float32(a.X), float32(a.Y))
	path.LineTo(float32(b.X), float32(b.Y))
	vs, is := path.AppendVerticesAndIndicesForStroke(vertices[:0], indices[:0], so)
	screen.DrawTriangles(vs, is, whiteSubImage, dto)

}

func StrokeCircle(screen *ebiten.Image, radius, strokeWidth float64, pos vec.Vec2, c color.Color) {
	vector.StrokeCircle(
		screen,
		float32(pos.X),
		float32(pos.Y),
		float32(radius),
		float32(strokeWidth), c, true)
}

func FillCircle(screen *ebiten.Image, radius float64, pos vec.Vec2, c color.Color) {
	vector.DrawFilledCircle(
		screen,
		float32(pos.X),
		float32(pos.Y),
		float32(radius), c, true)
}

func DrawChipmunkBB(screen *ebiten.Image, bb cm.BB, screenHeight float64, clr color.Color) {
	w := float32(bb.R - bb.L)
	h := float32(bb.T - bb.B)
	topLeft := vec.Vec2{bb.L, bb.T}.FlipVertical(screenHeight)
	vector.StrokeRect(screen, float32(topLeft.X), float32(topLeft.Y), w, h, 3, clr, false)
}
func DrawChipmunkBBGEOM(screen *ebiten.Image, bb cm.BB, screenHeight float64, geom *ebiten.GeoM) {
	w := float32(bb.R - bb.L)
	h := float32(bb.T - bb.B)
	topLeft := vec.Vec2{bb.L, bb.T}.FlipVertical(screenHeight)
	topLeft = ApplyGeoM2Vec2(topLeft, geom)
	vector.StrokeRect(screen, float32(topLeft.X), float32(topLeft.Y), w, h, 3, color.RGBA{255, 0, 0, 0}, false)
}

func ApplyGeoM2Vec2(pos vec.Vec2, geom *ebiten.GeoM) vec.Vec2 {
	x, y := geom.Apply(pos.X, pos.Y)
	return vec.Vec2{x, y}
}
