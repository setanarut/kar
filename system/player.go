package system

import (
	"image"
	"kar"
	"kar/arc"
	"kar/items"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	blockHealth                  float64
	targetBlockPos               image.Point
	placeBlock                   image.Point
	playerTile                   image.Point
	playerCenterX, playerCenterY float64
	IsRayHit                     bool
)

func (c *PlayerSys) Init() {
}

type PlayerSys struct {
}

func (c *PlayerSys) Update() {

	if kar.WorldECS.Alive(PlayerEntity) {

		CTRL.UpdateInput()
		CTRL.UpdateState()
		CTRL.UpdatePhysics()
		playerCenterX, playerCenterY = CTRL.Rect.X+CTRL.Rect.W/2, CTRL.Rect.Y+CTRL.Rect.H/2
		playerTile = Map.WorldToTile(playerCenterX, playerCenterY)
		targetBlockTemp := targetBlockPos
		targetBlockPos, IsRayHit = Map.Raycast(playerTile, CTRL.InputAxisLast, kar.RaycastDist)
		// eğer block odağı değiştiyse saldırıyı sıfırla
		if !targetBlockPos.Eq(targetBlockTemp) || !IsRayHit {
			blockHealth = 0
		}

		// Drop Item
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			currentSlot := CTRL.Inventory.SelectedSlot()
			if currentSlot.ID != items.Air {
				AppendToSpawnList(playerCenterX, playerCenterY, currentSlot.ID, currentSlot.Durability)
				CTRL.Inventory.RemoveItemFromSelectedSlot()
			}
		}

		// Place block
		if CTRL.IsPlaceKeyJustPressed {
			anyItemOverlapsWithPlaceCoords := false
			if IsRayHit && items.HasTag(CTRL.Inventory.SelectedSlot().ID, items.Block) {
				placeBlock = targetBlockPos.Sub(CTRL.InputAxisLast)
				queryItem := arc.FilterItem.Query(&kar.WorldECS)
				for queryItem.Next() {
					_, itemRect, _, _ := queryItem.Get()
					anyItemOverlapsWithPlaceCoords = itemRect.Overlaps(Map.GetTileRect(placeBlock))
					if anyItemOverlapsWithPlaceCoords {
						queryItem.Close()
						break
					}
				}
				if !anyItemOverlapsWithPlaceCoords {
					if !CTRL.Rect.Overlaps(Map.GetTileRect(placeBlock)) {
						Map.SetTile(placeBlock, CTRL.Inventory.SelectedSlotID())
						CTRL.Inventory.RemoveItemFromSelectedSlot()
					}
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			CTRL.Inventory.SelectPrevSlot()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			CTRL.Inventory.SelectNextSlot()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			CTRL.Inventory.ClearSelectedSlot()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			CTRL.Inventory.RandomFillAllSlots()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyV) {
			kar.DrawDebugTextEnabled = !kar.DrawDebugTextEnabled
		}

	}

}

func (c *PlayerSys) Draw() {

}
