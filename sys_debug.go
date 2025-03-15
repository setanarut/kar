package kar

import (
	"fmt"
	"image"
	"kar/items"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/v"
)

type Debug struct {
	drawItemHitboxEnabled       bool
	drawPlayerTileHitboxEnabled bool
	drawDebugTextEnabled        bool
	tile                        uint8
}

func (d *Debug) Init() {}
func (d *Debug) Update() {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		pos := tileMapRes.FloorToBlockCenter(cameraRes.ScreenToWorld(ebiten.CursorPosition()))
		mapPlatform.NewEntity(
			&AABB{
				Pos:  v.Vec{pos.X, pos.Y},
				Half: v.Vec{10, 10},
			},
			&Velocity{1, 0},
			ptr(PlatformType("solid")),
		)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		pos := tileMapRes.FloorToBlockCenter(cameraRes.ScreenToWorld(ebiten.CursorPosition()))
		mapPlatform.NewEntity(
			&AABB{
				Pos:  v.Vec{pos.X, pos.Y},
				Half: v.Vec{10, 10},
			},
			&Velocity{1, 0},
			ptr(PlatformType("oneway")),
		)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		x, y := cameraRes.ScreenToWorld(ebiten.CursorPosition())
		mapEnemy.NewEntity(
			&AABB{
				Pos:  v.Vec{x, y},
				Half: v.Vec{10, 10},
			},
			&Velocity{0.5, 0.5},
			ptr(AI("worm")),
		)
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
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		inventoryRes.SetSlot(0, items.Coal, 64, 0)
		inventoryRes.SetSlot(1, items.RawGold, 64, 0)
		inventoryRes.SetSlot(2, items.RawIron, 64, 0)
		inventoryRes.SetSlot(3, items.Stick, 64, 0)
		inventoryRes.SetSlot(4, items.DiamondPickaxe, 1, items.GetDefaultDurability(items.DiamondPickaxe))
		inventoryRes.SetSlot(5, items.DiamondShovel, 1, items.GetDefaultDurability(items.DiamondShovel))
		inventoryRes.SetSlot(6, items.DiamondAxe, 1, items.GetDefaultDurability(items.DiamondAxe))
		inventoryRes.SetSlot(7, items.Diamond, 64, 0)
		inventoryRes.SetSlot(8, items.Snowball, 64, 0)

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		d.drawDebugTextEnabled = !d.drawDebugTextEnabled
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		d.drawItemHitboxEnabled = !d.drawItemHitboxEnabled
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		d.drawPlayerTileHitboxEnabled = !d.drawPlayerTileHitboxEnabled
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
	if world.Alive(currentPlayer) {
		box, vel, _, playerController, _ := mapPlayer.GetUnchecked(currentPlayer)
		if d.drawPlayerTileHitboxEnabled {
			DrawAABB(box)
		}
		if d.drawDebugTextEnabled {
			ebitenutil.DebugPrintAt(Screen, fmt.Sprintf(
				"state %v\nVel.X: %.2f\nVel.Y: %.2f\nCamera: %v",
				playerController.CurrentState,
				vel.X,
				vel.Y,
				cameraRes,
			), 0, 10)
		}

	} else {
		ebitenutil.DebugPrintAt(Screen, fmt.Sprintf("Camera: %v", cameraRes), 0, 10)
	}

	ebitenutil.DebugPrintAt(Screen, fmt.Sprintf("DEBUG MODE: %v", debugEnabled), int(ScreenSize.X)-60, 10)
}
