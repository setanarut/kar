package system

import (
	"github.com/hajimehoshi/ebiten/v2"
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

func GetAxis() vec.Vec2 {
	axis := vec.Vec2{}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		axis.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		axis.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		axis.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		axis.X += 1
	}
	return axis
}
