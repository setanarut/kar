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
	damage                       float64 = 0.3
	raycastDist                  int     = 3 // block unit
	blockHealth                  float64
	targetBlockPos               image.Point
	placeBlock                   image.Point
	playerTile                   image.Point
	playerCenterX, playerCenterY float64
	IsRayHit                     bool
)

type PlayerSys struct {
}

func (c *PlayerSys) Update() {
	q := arc.FilterPlayer.Query(&kar.WorldECS)
	for q.Next() {
		anim, _, dop, rect, _ := q.Get()
		PlayerController.DOP = dop
		PlayerController.AnimPlayer = anim
		playerCenterX, playerCenterY = rect.X+rect.W/2, rect.Y+rect.H/2
		PlayerController.UpdateInput()
		dx, dy := PlayerController.UpdatePhysics(rect.X, rect.Y, rect.W, rect.H)

		playerTile = Map.WorldToTile(playerCenterX, playerCenterY)
		targetBlockTemp := targetBlockPos
		targetBlockPos, IsRayHit = Map.Raycast(playerTile, PlayerController.InputAxisLast, raycastDist)
		// eğer block odağı değiştiyse saldırıyı sıfırla
		if !targetBlockPos.Eq(targetBlockTemp) || !IsRayHit {
			blockHealth = 0
		}
		rect.X += dx
		rect.Y += dy
		PlayerController.UpdateState()

		// Place block
		if PlayerController.IsPlaceKeyJustPressed {
			if items.IsBlock(PlayerInventory.SelectedSlotID()) {
				anyItemOverlaps := false
				playerOverlaps := false
				if IsRayHit {
					placeBlock = targetBlockPos.Sub(PlayerController.InputAxisLast)
					playerOverlaps = rect.Overlaps(Map.GetTileRect(placeBlock))
				}
				queryItem := arc.FilterItem.Query(&kar.WorldECS)
				for queryItem.Next() {
					_, _, itemRect, _ := queryItem.Get()
					anyItemOverlaps = itemRect.Overlaps(Map.GetTileRect(placeBlock))
				}
				IsBlockPlaceable := !(anyItemOverlaps || playerOverlaps)
				if IsRayHit && IsBlockPlaceable && PlayerInventory.SelectedSlotQuantity() > 0 {
					Map.SetTile(placeBlock, PlayerInventory.SelectedSlotID())
					PlayerInventory.RemoveItemFromSelectedSlot()
				}
			}
		}
	}

	// Drop Item
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		if PlayerInventory.SelectedSlotQuantity() > 0 {
			AppendToSpawnList(playerCenterX, playerCenterY, PlayerInventory.SelectedSlotID())
			PlayerInventory.RemoveItemFromSelectedSlot()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		PlayerInventory.SelectPrevSlot()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		PlayerInventory.SelectNextSlot()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		PlayerInventory.ClearSelectedSlot()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		PlayerInventory.RandomFillAllSlots()
	}
}

func (c *PlayerSys) Draw() {

}
func (c *PlayerSys) Init() {

}
