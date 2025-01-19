package system

import (
	"kar"
	"kar/arc"
	"kar/items"
	"kar/res"

	"github.com/setanarut/tilecollider"
)

type Enemy struct {
	maxFallVelocity float64
	gravity         float64
	speedX          float64
	jumpH           float64
}

func (en *Enemy) Init() {
	en.maxFallVelocity = 2.5
	en.gravity = 0.5
	en.speedX = 3.5
	en.jumpH = 9
}
func (en *Enemy) Update() {
	// enemy system
	enemyQuery := arc.FilterEnemy.Query(&kar.WorldECS)
	for enemyQuery.Next() {
		rect, vel, _ := enemyQuery.Get()

		vel.VelY += en.gravity
		vel.VelY = min(vel.VelY, en.maxFallVelocity)

		collider.Collide(
			rect.X,
			rect.Y,
			rect.W,
			rect.H,
			vel.VelX,
			vel.VelY,
			func(collisionInfos []tilecollider.CollisionInfo[uint16], dx, dy float64) {

				// Apply tilemap collision response
				rect.X += dx
				rect.Y += dy

				// change direction when wall collision
				for _, collisionInfo := range collisionInfos {
					if collisionInfo.Normal[0] == -1 || collisionInfo.Normal[0] == 1 {
						vel.VelX *= -1
					}
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
