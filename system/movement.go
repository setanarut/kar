package system

import (
	"kar"
	"kar/arc"

	"github.com/setanarut/tilecollider"
)

var marioVel = vec2{0, 4}

func (s *Movement) Init() {

}
func (s *Movement) Update() {
	q := arc.FilterMovement.Query(&kar.WorldECS)

	for q.Next() {
		ctrl, rect := q.Get()

		if marioVel.Y < 0 {
			ctrl.IsOnFloor = false
		}
		marioVel := ctrl.ProcessVelocity(InputAxis, marioVel)

		dx, dy := TCollider.Collide(
			rect.X,
			rect.Y,
			rect.W,
			rect.H,
			marioVel.X,
			marioVel.Y,
			func(ci []tilecollider.CollisionInfo[uint16], f1, f2 float64) {
				for _, v := range ci {
					if v.Normal[1] == -1 {
						ctrl.IsOnFloor = true
					}
					if v.Normal[1] == 1 {
						ctrl.IsJumping = false
						marioVel.Y = 0
					}
					if v.Normal[0] == 1 || v.Normal[0] == -1 {
						marioVel.X = 0
					}
				}
			},
		)

		rect.X += dx
		rect.Y += dy
	}
}
func (s *Movement) Draw() {}

type Movement struct{}
