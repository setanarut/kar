package system

import (
	"kar"
	"kar/arc"
	"kar/items"
	"kar/res"

	"github.com/setanarut/tilecollider"
)

type Enemy struct {
	// maxFallVelocity float64
	// gravity         float64
}

func (en *Enemy) Init() {
	// en.maxFallVelocity = 2.5
	// en.gravity = 0.5
}

var normal [2]int

func (en *Enemy) Update() {
	// enemy system
	enemyQuery := arc.FilterEnemy.Query(&kar.WorldECS)
	for enemyQuery.Next() {
		rect, vel, _ := enemyQuery.Get()

		// vel.VelY += en.gravity
		// vel.VelY = min(vel.VelY, en.maxFallVelocity)
		normal = [2]int{0, 0}
		collider.Collide(
			rect.X,
			rect.Y,
			rect.W,
			rect.H,
			vel.X,
			vel.Y,
			func(infos []tilecollider.CollisionInfo[uint8], dx, dy float64) {

				// Apply tilemap collision response
				rect.X += dx
				rect.Y += dy

				for _, info := range infos {
					if normal[0] != info.Normal[0] {
						normal[0] += info.Normal[0]
					}
					if normal[1] != info.Normal[1] {
						normal[1] += info.Normal[1]
					}
				}

				switch normal {
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
	}
}
func (en *Enemy) Draw() {
	q := arc.FilterEnemy.Query(&kar.WorldECS)
	for q.Next() {
		rect, _, _ := q.Get()
		kar.ColorMDIO.GeoM.Reset()
		kar.ColorMDIO.GeoM.Translate(rect.X, rect.Y)
		kar.Camera.DrawWithColorM(res.Icon8[items.Sand], kar.ColorM, kar.ColorMDIO, kar.Screen)
	}
}
