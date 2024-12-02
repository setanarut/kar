package system

import (
	"kar"
	"kar/arc"
)

func (s *Movement) Init() {
}
func (s *Movement) Update() {
	q := arc.FilterMovement.Query(&kar.WorldECS)

	for q.Next() {
		_, rect := q.Get()
		// vel = InputAxis.Scale(3)
		// vel = ctrl.ProcessVelocity(InputAxis, vel)

		rect.X += InputAxis.X * 3
		rect.Y += InputAxis.Y * 3

	}

}
func (s *Movement) Draw() {}

type Movement struct{}
