package kar

import (
	"kar/items"
	"kar/res"

	"github.com/setanarut/tilecollider"
)

type Enemy struct {
	normal   [2]int
	wormSize Size
}

func (e *Enemy) Init() {
	e.wormSize = Size{8, 8}
}

func (e *Enemy) Update() {
	enemyQuery := FilterEnemy.Query(&ECWorld)
	for enemyQuery.Next() {
		pos, vel, ai := enemyQuery.Get()
		switch ai.Name {
		case "worm":
			e.normal = [2]int{0, 0}
			Collider.Collide(
				pos.X,
				pos.Y,
				e.wormSize.W,
				e.wormSize.H,
				vel.X,
				vel.Y,
				func(infos []tilecollider.CollisionInfo[uint8], dx, dy float64) {

					// Apply tilemap collision response
					pos.X += dx
					pos.Y += dy

					for _, info := range infos {
						if e.normal[0] != info.Normal[0] {
							e.normal[0] += info.Normal[0]
						}
						if e.normal[1] != info.Normal[1] {
							e.normal[1] += info.Normal[1]
						}
					}

					switch e.normal {
					case [2]int{0, -1}: // bottom
						vel.X = 1
					case [2]int{1, -1}: // bottomleft
						vel.X = 1
					case [2]int{-1, -1}: // bottomright
						vel.Y = -1
					case [2]int{-1, 1}: // topright
						vel.X = -1
					case [2]int{1, 1}: // topleft
						vel.Y = 1
					}
				},
			)
		case "other":
			// other
		}
	}
}
func (e *Enemy) Draw() {
	q := FilterEnemy.Query(&ECWorld)
	for q.Next() {
		pos, _, _ := q.Get()
		ColorMDIO.GeoM.Reset()
		ColorMDIO.GeoM.Translate(pos.X, pos.Y)
		CameraRes.DrawWithColorM(res.Icon8[items.Sand], ColorM, ColorMDIO, Screen)
	}
}
