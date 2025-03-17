package kar

import (
	"kar/res"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// block breaking effect system
type Effects struct {
	g float64
}

func (e *Effects) Init() {
	e.g = 0.2

}
func (e *Effects) Update() {

	q := filterEffect.Query()

	for q.Next() {
		_, pos, vel, angle := q.Get()
		vel.Y += e.g
		pos.Y += vel.Y

		if math.Signbit(float64(*angle)) {
			*angle -= 0.2
		} else {
			*angle += 0.2
		}

		pos.X += vel.X * 2

		if pos.Y > cameraRes.Y+cameraRes.Height {
			toRemove = append(toRemove, q.Entity())
		}
	}
}
func (e *Effects) Draw() {
	q := filterEffect.Query()
	for q.Next() {
		id, pos, _, angle := q.Get()
		colorMDIO.GeoM = ebiten.GeoM{}
		colorMDIO.GeoM.Translate(-4, -4)
		colorMDIO.GeoM.Rotate(float64(*angle))
		colorMDIO.GeoM.Translate(4, 4)
		colorMDIO.GeoM.Translate(pos.X, pos.Y)
		colorM.Scale(1.3, 1.3, 1.3, 1)
		if *id != 0 {
			cameraRes.DrawWithColorM(res.Icon8[uint8(*id)], colorM, colorMDIO, Screen)
		}
		colorM.Reset()
	}
}
