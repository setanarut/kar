package kar

import (
	"fmt"
	"image"
	"kar/items"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Debug struct {
	tile uint8
}

func (d *Debug) Init() {
}
func (d *Debug) Update() {

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		x, y := cameraRes.ScreenToWorld(ebiten.CursorPosition())
		SpawnEnemy(Vec{x, y}, Vec{0.5, 0})
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := cameraRes.ScreenToWorld(ebiten.CursorPosition())
		p := tileMapRes.WorldToTile(x, y)
		tileMapRes.Set(p.X, p.Y, d.tile)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		x, y := cameraRes.ScreenToWorld(ebiten.CursorPosition())
		p := tileMapRes.WorldToTile(x, y)
		d.tile = tileMapRes.GetID(p.X, p.Y)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		inventoryRes.ClearCurrentSlot()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		inventoryRes.RandomFillAllSlots()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		drawDebugTextEnabled = !drawDebugTextEnabled
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		drawItemHitboxEnabled = !drawItemHitboxEnabled
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		drawPlayerTileHitboxEnabled = !drawPlayerTileHitboxEnabled
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		dataManager.SaveItem("map.png", tileMapRes.GetImageByte())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		box := mapAABB.GetUnchecked(currentPlayer)
		tileMapRes.Set(tileMapRes.W/2, tileMapRes.H-3, items.Air)
		box.Pos = tileMapRes.TileToWorld(image.Point{tileMapRes.W / 2, tileMapRes.H - 3})
		cameraRes.SetCenter(box.Pos.X, box.Pos.Y)
	}

}
func (d *Debug) Draw() {
	// Draw debug info
	if drawDebugTextEnabled {
		if world.Alive(currentPlayer) {
			_, vel, _, playerController, _ := mapPlayer.GetUnchecked(currentPlayer)
			ebitenutil.DebugPrintAt(Screen, fmt.Sprintf(
				"state %v\nVel.X: %.2f\nVel.Y: %.2f\nCamera: %v",
				playerController.CurrentState,
				vel.X,
				vel.Y,
				cameraRes,
			), 0, 10)
		} else {
			ebitenutil.DebugPrintAt(Screen, fmt.Sprintf(
				"Camera: %v",
				cameraRes,
			), 0, 10)
		}
	}
}
