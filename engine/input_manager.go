package engine

import (
	"kar/engine/vec"

	"github.com/hajimehoshi/ebiten/v2"
)

type InputManager struct {
	WASDDirection            vec.Vec2
	ArrowDirection           vec.Vec2
	LastPressedWASDDirection vec.Vec2
}

func (i *InputManager) UpdateArrowDirection() {

	i.ArrowDirection = vec.Vec2{}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		i.ArrowDirection.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		i.ArrowDirection.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		i.ArrowDirection.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		i.ArrowDirection.X += 1
	}

}

func (i *InputManager) UpdateWASDDirection() {
	i.WASDDirection = vec.Vec2{}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		i.WASDDirection.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		i.WASDDirection.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		i.WASDDirection.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		i.WASDDirection.X += 1
	}

	if !i.WASDDirection.Equal(vec.Vec2{}) {
		i.LastPressedWASDDirection = i.WASDDirection
	}

}

func (i *InputManager) AnyKeyDown(keys ...ebiten.Key) bool {
	for _, key := range keys {
		if ebiten.IsKeyPressed(key) {
			return true
		}
	}
	return false
}

func (i *InputManager) IsPressedAndNotABC(onlyKey, a, b, c ebiten.Key) bool {
	return ebiten.IsKeyPressed(onlyKey) && !ebiten.IsKeyPressed(a) && !ebiten.IsKeyPressed(b) && !ebiten.IsKeyPressed(c)
}
