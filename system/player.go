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
	isRayHit                     bool
)

func (c *Player) Init() {
	anim, hlt, rect := arc.MapPlayer.Get(player)
	ctrl.AnimPlayer = anim
	ctrl.Rect = rect
	ctrl.Health = hlt
	ctrl.Inventory = items.NewInventory()
	// ctrl.Inventory.SetSlot(2, items.OakLog, 64, 0)
	// ctrl.Inventory.SetSlot(5, items.IronIngot, 32, 0)
	ctrl.EnterFalling()
}

type Player struct {
}

func (c *Player) Update() {

	if kar.WorldECS.Alive(player) {

		if !craftingState {
			ctrl.UpdateInput()
			ctrl.UpdateState()
			ctrl.UpdatePhysics()
			playerCenterX, playerCenterY = ctrl.Rect.X+ctrl.Rect.W/2, ctrl.Rect.Y+ctrl.Rect.H/2
			playerTile = tileMap.WorldToTile(playerCenterX, playerCenterY)
			targetBlockTemp := targetBlockPos
			targetBlockPos, isRayHit = tileMap.Raycast(playerTile, ctrl.InputAxisLast, kar.RaycastDist)
			// eğer block odağı değiştiyse saldırıyı sıfırla
			if !targetBlockPos.Eq(targetBlockTemp) || !isRayHit {
				blockHealth = 0
			}

			// Drop Item
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
				currentSlot := ctrl.Inventory.CurrentSlotData()
				if currentSlot.ID != items.Air {
					AppendToSpawnList(playerCenterX, playerCenterY, currentSlot.ID, currentSlot.Durability)
					ctrl.Inventory.RemoveItemFromSelectedSlot()
				}
			}

			// Place block
			if ctrl.IsPlaceKeyJustPressed {
				anyItemOverlapsWithPlaceCoords := false
				if isRayHit && items.HasTag(ctrl.Inventory.CurrentSlotData().ID, items.Block) {
					placeBlock = targetBlockPos.Sub(ctrl.InputAxisLast)
					queryItem := arc.FilterItem.Query(&kar.WorldECS)
					for queryItem.Next() {
						_, itemRect, _, _ := queryItem.Get()
						anyItemOverlapsWithPlaceCoords = itemRect.Overlaps(tileMap.GetTileRect(placeBlock))
						if anyItemOverlapsWithPlaceCoords {
							queryItem.Close()
							break
						}
					}
					if !anyItemOverlapsWithPlaceCoords {
						if !ctrl.Rect.Overlaps(tileMap.GetTileRect(placeBlock)) {
							tileMap.Set(placeBlock.X, placeBlock.Y, ctrl.Inventory.CurrentSlot())
							ctrl.Inventory.RemoveItemFromSelectedSlot()
						}
					}
				}
			}

			// Remove dead player entity
			if ctrl.Health.Health <= 0 {
				toRemove = append(toRemove, player)
			}

		}
	}

}

func (c *Player) Draw() {

}
