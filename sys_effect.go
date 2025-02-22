package kar

import (
	"kar/res"
	"math"
)

type Effects struct {
	g float64
}

func (e *Effects) Init() {
	e.g = 0.2

}
func (e *Effects) Update() {

	q := filterEffect.Query()

	for q.Next() {
		_, p, v, _ := q.Get()
		v.Y += e.g
		p.Y += v.Y
		if p.Y > cameraRes.Y+cameraRes.Height {
			toRemove = append(toRemove, q.Entity())
		}
	}
}
func (e *Effects) Draw() {
	q := filterEffect.Query()
	for q.Next() {
		id, p, v, r := q.Get()
		ColorMDIO.GeoM.Reset()
		ColorMDIO.GeoM.Translate(-4, -4)
		ColorMDIO.GeoM.Rotate(r.Angle)
		ColorMDIO.GeoM.Translate(4, 4)
		ColorMDIO.GeoM.Translate(p.X, p.Y)
		ColorM.Scale(1.3, 1.3, 1.3, 1)
		cameraRes.DrawWithColorM(res.Icon8[id.ID], ColorM, ColorMDIO, Screen)
		ColorM.Reset()

		if math.Signbit(r.Angle) {
			r.Angle -= 0.2
		} else {
			r.Angle += 0.2
		}

		p.X += v.X * 2
	}
}
