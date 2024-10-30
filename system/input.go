package system

import (
	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/vec"
)

var inputAxis = vec.Vec2{}
var inputAxisLast = vec.Vec2{}

type Input struct{}

func (sys *Input) Init() {}
func (sys *Input) Draw() {}
func (sys *Input) Update() {

	if !inputAxis.Equal(vec.Vec2{}) {
		inputAxisLast = inputAxis
	}
	inputAxis = GetAxis()
}

func GetAxis() vec2 {
	axis := vec2{}
	if pressed(eb.KeyW) {
		axis.Y -= 1
	}
	if pressed(eb.KeyS) {
		axis.Y += 1
	}
	if pressed(eb.KeyA) {
		axis.X -= 1
	}
	if pressed(eb.KeyD) {
		axis.X += 1
	}
	return axis
}
