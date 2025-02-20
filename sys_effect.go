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
func (e *Effects) Update() error {

	q := FilterEffect.Query(&ECWorld)

	for q.Next() {
		_, p, v, _ := q.Get()
		v.Y += e.g
		p.Y += v.Y
		if p.Y > CameraRes.Y+CameraRes.Height {
			toRemove = append(toRemove, q.Entity())
		}
	}
	return nil
}
func (e *Effects) Draw() {
	q := FilterEffect.Query(&ECWorld)
	for q.Next() {
		id, p, v, r := q.Get()
		ColorMDIO.GeoM.Reset()
		ColorMDIO.GeoM.Translate(-4, -4)
		ColorMDIO.GeoM.Rotate(r.Angle)
		ColorMDIO.GeoM.Translate(4, 4)
		ColorMDIO.GeoM.Translate(p.X, p.Y)
		ColorM.Scale(1.3, 1.3, 1.3, 1)
		CameraRes.DrawWithColorM(res.Icon8[id.ID], ColorM, ColorMDIO, Screen)
		ColorM.Reset()

		if math.Signbit(r.Angle) {
			r.Angle -= 0.2
		} else {
			r.Angle += 0.2
		}

		p.X += v.X * 2
	}
}
