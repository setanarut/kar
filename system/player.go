package system

import (
	"image"
	"kar"
	"kar/arc"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	damage float64 = 0.1
	// MinAttackBlockDist           float64 = 2.0
	BlockPlacementDistance       = 3
	blockHealth                  float64
	targetBlock                  image.Point
	placeBlock                   image.Point
	playerTile                   image.Point
	playerCenterX, playerCenterY float64
	isBlockPlaceable             bool
	IsRaycastHit                 bool
	DOP                          *arc.DrawOptions
)

type PlayerSys struct {
}

func (c *PlayerSys) Update() {
	q := arc.FilterPlayer.Query(&kar.WorldECS)

	for q.Next() {

		anim, _, dop, rect, _ := q.Get()

		playerCenterX, playerCenterY = rect.X+rect.W/2, rect.Y+rect.H/2

		PlayerController.AnimPlayer = anim
		DOP = dop
		PlayerController.UpdateInput()

		dx, dy := PlayerController.UpdatePhysics(rect.X, rect.Y, rect.W, rect.H)

		playerTile = Map.WorldToTile(playerCenterX, playerCenterY)
		targetBlockTemp := targetBlock

		targetBlock, IsRaycastHit = Map.Raycast(playerTile, PlayerController.InputAxisLast, BlockPlacementDistance)

		// eğer block odağı değiştiyse saldırıyı sıfırla
		if !targetBlock.Eq(targetBlockTemp) {
			blockHealth = 0
		}

		if IsRaycastHit {
			placeBlock = targetBlock.Sub(PlayerController.InputAxisLast)
			isBlockPlaceable = !rect.Overlaps(Map.GetTileRect(placeBlock))
		} else {
			isBlockPlaceable = false
			blockHealth = 0
		}

		rect.X += dx
		rect.Y += dy
		PlayerController.UpdateState()
	}

	// Place block
	if PlayerController.IsPlaceKeyJustPressed {
		if IsRaycastHit && isBlockPlaceable && PlayerInventory.SelectedSlotQuantity() > 0 {
			Map.SetTile(placeBlock, PlayerInventory.SelectedSlotID())
			PlayerInventory.RemoveItemFromSelectedSlot()
		}
	}
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
