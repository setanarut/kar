package system

import (
	"kar"
	"kar/arc"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/setanarut/tilecollider"
)

var state string

func (s *Movement) Init() {

}
func (s *Movement) Update() {
	q := arc.FilterMovement.Query(&kar.WorldECS)
	for q.Next() {
		controller, rect, anim, dop := q.Get()

		// dop.FlipX = math.Signbit(controller.VelX)
		if controller.VelX > 0 {
			dop.FlipX = false
		} else if controller.VelX < 0 {
			dop.FlipX = true
		}

		controller.UpdateInput()
		controller.UpdatePhysics()
		controller.UpdateState(anim)
		controller.IsOnFloor = false
		Collider.Collide(
			rect.X,
			rect.Y,
			rect.W,
			rect.H,
			controller.VelX,
			controller.VelY,
			func(ci []tilecollider.CollisionInfo[uint16], dx, dy float64) {
				for _, v := range ci {
					if v.Normal[1] == -1 {
						controller.VelY = 0
						controller.IsOnFloor = true
					}
					if v.Normal[1] == 1 {
						controller.VelY = 0
					}
					if v.Normal[0] == -1 {
						controller.VelX = 0
					}
					if v.Normal[0] == 1 {
						controller.VelX = 0
					}
				}
				rect.X += dx
				rect.Y += dy
			},
		)

	}
}
func (s *Movement) Draw() {
	ebitenutil.DebugPrint(kar.Screen, state)
}

type Movement struct{}
