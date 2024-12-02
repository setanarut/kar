package system

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var InputAxis vec2
var InputAxisLast vec2

type Input struct{}

func (sys *Input) Init() {}
func (sys *Input) Draw() {}
func (sys *Input) Update() {

	if !InputAxis.Equals(vec2{}) {
		InputAxisLast = InputAxis
	}
	InputAxis = Axis()
}

func Axis() (axis vec2) {
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
