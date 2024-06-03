package engine

import (
	"bytes"
	"embed"
	"image"
	"kar/engine/cm"
	"log"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	"golang.org/x/text/language"
)

func GetEbitenImageOffset(img *ebiten.Image) cm.Vec2 {
	return cm.Vec2{float64(img.Bounds().Dx()), float64(img.Bounds().Dy())}.Mult(0.5).Neg()
}

func MapRange(v, a, b, c, d float64) float64 {
	return (v-a)/(b-a)*(d-c) + c
}

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func Degree(radians float64) float64 {
	return radians * 180.0 / math.Pi
}

// IsMoving hız vektörü hareket ediyor mu?
func IsMoving(velocityVector cm.Vec2, minSpeed float64) bool {
	if math.Abs(velocityVector.X) < minSpeed && math.Abs(velocityVector.Y) < minSpeed {
		return true
	} else {
		return false
	}
}

// InvDirVectY inverts direction unit vector Y axis beetween bottom-left and top-left coordinate systems
func InvDirVectY(v cm.Vec2) cm.Vec2 {
	v.Y = v.Y * -1
	return v
}

// InvPosVectY inverts position Y axis beetween bottom-left and top-left coordinate systems
func InvPosVectY(v cm.Vec2, screenbHeight float64) cm.Vec2 {
	v.Y = screenbHeight - v.Y
	return v
}

// InvertAngle invert angle
func InvertAngle(angle float64) float64 {
	return angle * -1
}

// Rotate a vector by an angle in radians
func Rotate(v cm.Vec2, angle float64) cm.Vec2 {
	return cm.Vec2{
		X: v.X*math.Cos(angle) - v.Y*math.Sin(angle),
		Y: v.X*math.Sin(angle) + v.Y*math.Cos(angle),
	}
}

// RotateAbout rotates point about origin
func RotateAbout(angle float64, point, origin cm.Vec2) cm.Vec2 {
	b := cm.Vec2{}
	b.X = math.Cos(angle)*(point.X-origin.X) - math.Sin(angle)*(point.Y-origin.Y) + origin.X
	b.Y = math.Sin(angle)*(point.X-origin.X) + math.Cos(angle)*(point.Y-origin.Y) + origin.Y
	return b
}

// PointOnCircle returns point at angle
func PointOnCircle(center cm.Vec2, radius float64, angle float64) cm.Vec2 {
	x := center.X + (radius * math.Cos(angle))
	y := center.Y + (radius * math.Sin(angle))
	return cm.Vec2{x, y}
}

func RandomPoint(minX, maxX, minY, maxY float64) cm.Vec2 {
	return cm.Vec2{X: minX + rand.Float64()*(maxX-minX), Y: minY + rand.Float64()*(maxY-minY)}
}
func RandomPointInBB(bb cm.BB, margin float64) cm.Vec2 {
	return RandomPoint(bb.L+margin, bb.R-margin, bb.T-margin, bb.B+margin)
}

func RandRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)

}
func RandRangeInt(min, max int) int {
	return rand.IntN(max-min+1) + min

}

func LoadTextFace(fileName string, size float64, assets embed.FS) *text.GoTextFace {
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

func LoadImage(name string, assets embed.FS) *ebiten.Image {

	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

// Linspace returns evenly spaced numbers over a specified closed interval.
func Linspace(start, stop float64, num int) (res []float64) {
	if num <= 0 {
		return []float64{}
	}
	if num == 1 {
		return []float64{start}
	}
	step := (stop - start) / float64(num-1)
	res = make([]float64, num)
	res[0] = start
	for i := 1; i < num; i++ {
		res[i] = start + float64(i)*step
	}
	res[num-1] = stop
	return
}

// Clamp returns f clamped to [low, high]
func Clamp(f, low, high float64) float64 {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}

func GetBoxScaleFactor(imageW, imageH, targetW, targetH float64) cm.Vec2 {
	return cm.Vec2{(targetW / imageW), (targetH / imageH)}
}
func GetCircleScaleFactor(radius float64, image *ebiten.Image) cm.Vec2 {
	scaleX := 2 * radius / float64(image.Bounds().Dx())
	return cm.Vec2{scaleX, scaleX}
}

func AddComponents(e *donburi.Entry, comps ...donburi.IComponentType) {
	for _, comp := range comps {
		e.AddComponent(comp)
	}
}
