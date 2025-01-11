package system

import (
	"fmt"
	"image/color"
	"kar"
	"kar/arc"
	"kar/items"
	"kar/res"
	"math/rand/v2"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	hotbarPositionX        = kar.ScreenW/2 - float64(res.Hotbar.Bounds().Dx())/2
	hotbarPositionY        = 8.
	hotbarRightEdgePosX    = hotbarPositionX + float64(res.Hotbar.Bounds().Dx())
	craftingTablePositionX = hotbarPositionX + 49
	craftingTablePositionY = hotbarPositionY + 39
	itemQuantityTextDO     = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions: text.LayoutOptions{
			LineSpacing: 10,
		},
	}
	debugInfo = `state %v
inputAxisLast %v
`
)

type UI struct{}

func (ui *UI) Init() {}
func (ui *UI) Update() {

	// Hotbar slot navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		ctrl.Inventory.SelectPrevSlot()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		ctrl.Inventory.SelectNextSlot()
	}

	// Toggle crafting state
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		id := tileMap.Get(targetBlockPos.X, targetBlockPos.Y)
		if id == items.CraftingTable {
			craftingTable.SlotPosX, craftingTable.SlotPosY = 1, 1

			// clear crafting table when exit
			if craftingState {
				for y := range 3 {
					for x := range 3 {
						id2 := craftingTable.Get(x, y).ID
						if id2 != 0 {
							quantity := craftingTable.Get(x, y).Quantity
							for range quantity {
								if ctrl.Inventory.AddItemIfEmpty(
									craftingTable.Get(x, y).ID,
									craftingTable.Get(x, y).Durability,
								) {
									craftingTable.RemoveItem(x, y)
								} else {
									craftingTable.RemoveItem(x, y)
									arc.SpawnItem(arc.SpawnData{
										X:          playerCenterX,
										Y:          playerCenterY,
										Id:         id2,
										Durability: craftingTable.Get(x, y).Durability,
									}, rand.IntN(sinspaceLen))
								}
							}
						}
					}
				}
				craftingTable.ResultSlot = items.Slot{}
			}

			craftingState = !craftingState
		}
	}

	// Crafting table slot navigation
	if craftingState {
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			craftingTable.SlotPosX = min(craftingTable.SlotPosX+1, 2)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			craftingTable.SlotPosX = max(craftingTable.SlotPosX-1, 0)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			craftingTable.SlotPosY = min(craftingTable.SlotPosY+1, 2)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			craftingTable.SlotPosY = max(craftingTable.SlotPosY-1, 0)
		}

		// move items from hotbar to crafting table
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			if ctrl.Inventory.CurrentSlotID() != 0 {
				Tableslot := craftingTable.CurrentSLot()
				if craftingTable.CurrentSlot().ID == 0 {
					id, dur := ctrl.Inventory.RemoveItemFromSelectedSlot()
					Tableslot.ID = id
					Tableslot.Durability = dur
					Tableslot.Quantity = 1
				} else if craftingTable.CurrentSlot().ID == ctrl.Inventory.CurrentSlotID() {
					ctrl.Inventory.RemoveItemFromSelectedSlot()
					Tableslot.Quantity++
				}
			}
			craftingTable.UpdateResultSlot()
		}
		// move items from crafting table to hotbar
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			if craftingTable.CurrentSlot().ID != 0 {
				if craftingTable.CurrentSlot().Quantity == 1 {
					if ctrl.Inventory.AddItemIfEmpty(craftingTable.CurrentSlot().ID, craftingTable.CurrentSlot().Durability) {
						craftingTable.ClearCurrenSlot()
					}
				} else if craftingTable.CurrentSlot().Quantity > 1 {
					if ctrl.Inventory.AddItemIfEmpty(craftingTable.CurrentSlot().ID, craftingTable.CurrentSlot().Durability) {
						craftingTable.SubCurrentSlotQuantity(1)

					}
				}

			}
			craftingTable.UpdateResultSlot()
		}
		// apply recipe
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			minimum := craftingTable.UpdateResultSlot()
			resultID := craftingTable.ResultSlot.ID
			dur := items.GetDefaultDurability(resultID)
			if resultID != 0 {
				for range minimum {
					if ctrl.Inventory.AddItemIfEmpty(resultID, dur) {
						for y := range 3 {
							for x := range 3 {
								if craftingTable.Slots[y][x].Quantity > 0 {
									craftingTable.Slots[y][x].Quantity--
								}
								if craftingTable.Slots[y][x].Quantity == 0 {
									craftingTable.Slots[y][x].ID = 0
								}
							}
						}
						craftingTable.ResultSlot.ID = 0
					}
				}
			}
			craftingTable.UpdateResultSlot()
		}
	}

	// Debug
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		kar.DrawDebugTextEnabled = !kar.DrawDebugTextEnabled
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		kar.DrawDebugHitboxesEnabled = !kar.DrawDebugHitboxesEnabled
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		ctrl.Inventory.ClearCurrentSlot()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ctrl.Inventory.RandomFillAllSlots()
	}

}
func (ui *UI) Draw() {
	if kar.WorldECS.Alive(player) {
		// Draw hotbar background
		kar.GlobalColorMDIO.GeoM.Reset()
		kar.GlobalColorMDIO.GeoM.Translate(hotbarPositionX, hotbarPositionY)
		colorm.DrawImage(kar.Screen, res.Hotbar, kar.GlobalColorM, kar.GlobalColorMDIO)

		// Draw slots
		for x := range 9 {
			slotID := ctrl.Inventory.Slots[x].ID
			quantity := ctrl.Inventory.Slots[x].Quantity
			SlotOffsetX := float64(x) * 17
			SlotOffsetX += hotbarPositionX

			// draw hotbar item icons
			kar.GlobalColorMDIO.GeoM.Reset()
			kar.GlobalColorMDIO.GeoM.Translate(SlotOffsetX+(5), hotbarPositionY+(5))
			if slotID != items.Air && ctrl.Inventory.Slots[x].Quantity > 0 {
				colorm.DrawImage(kar.Screen, res.Icon8[slotID], kar.GlobalColorM, kar.GlobalColorMDIO)
			}
			if x == ctrl.Inventory.CurrentSlotIndex {
				// Draw hotbar selected slot border
				kar.GlobalColorMDIO.GeoM.Translate(-5, -5)
				colorm.DrawImage(kar.Screen, res.SelectionBar, kar.GlobalColorM, kar.GlobalColorMDIO)

				// Draw hotbar slot item display name
				if !ctrl.Inventory.IsCurrentSlotEmpty() {
					itemQuantityTextDO.GeoM.Reset()
					itemQuantityTextDO.GeoM.Translate(SlotOffsetX-1, hotbarPositionY+14)
					if items.HasTag(slotID, items.Tool) {
						text.Draw(kar.Screen, fmt.Sprintf(
							"%v\nDurability %v",
							items.Property[slotID].DisplayName,
							ctrl.Inventory.Slots[x].Durability,
						), res.Font, itemQuantityTextDO)
					} else {
						text.Draw(kar.Screen, items.Property[slotID].DisplayName, res.Font, itemQuantityTextDO)
					}
				}
			}

			// Draw item quantity number
			if quantity > 1 && items.IsStackable(slotID) {
				itemQuantityTextDO.GeoM.Reset()
				itemQuantityTextDO.GeoM.Translate(SlotOffsetX+6, hotbarPositionY+4)
				num := strconv.FormatUint(uint64(quantity), 10)
				if quantity < 10 {
					num = " " + num
				}
				text.Draw(kar.Screen, num, res.Font, itemQuantityTextDO)
			}
		}

		// Draw player health text
		itemQuantityTextDO.GeoM.Reset()
		itemQuantityTextDO.GeoM.Translate(hotbarRightEdgePosX+8, hotbarPositionY)
		text.Draw(kar.Screen, fmt.Sprintf("Health %v", ctrl.Health.Health), res.Font, itemQuantityTextDO)

		// Draw crafting table
		if craftingState {

			// crafting table Background
			kar.GlobalColorMDIO.GeoM.Reset()
			kar.GlobalColorMDIO.GeoM.Translate(craftingTablePositionX, craftingTablePositionY)
			colorm.DrawImage(kar.Screen, res.CraftingTable, kar.GlobalColorM, kar.GlobalColorMDIO)

			// draw crafting table item icons
			for x := 0; x < 3; x++ {
				for y := 0; y < 3; y++ {
					if craftingTable.Slots[y][x].ID != items.Air {
						sx := craftingTablePositionX + float64(x*17)
						sy := craftingTablePositionY + float64(y*17)
						kar.GlobalColorMDIO.GeoM.Reset()
						kar.GlobalColorMDIO.GeoM.Translate(sx+6, sy+6)
						colorm.DrawImage(
							kar.Screen,
							res.Icon8[craftingTable.Slots[y][x].ID],
							kar.GlobalColorM,
							kar.GlobalColorMDIO,
						)

						// Draw item quantity number
						quantity := craftingTable.Slots[y][x].Quantity
						if quantity > 1 {
							itemQuantityTextDO.GeoM.Reset()
							itemQuantityTextDO.GeoM.Translate(sx+7, sy+5)
							num := strconv.FormatUint(uint64(quantity), 10)
							if quantity < 10 {
								num = " " + num
							}
							text.Draw(kar.Screen, num, res.Font, itemQuantityTextDO)
						}
					}

					// draw selected slot border of crqfting table
					if x == craftingTable.SlotPosX && y == craftingTable.SlotPosY {
						sx := craftingTablePositionX + float64(x*17)
						sy := craftingTablePositionY + float64(y*17)
						kar.GlobalColorMDIO.GeoM.Reset()
						kar.GlobalColorMDIO.GeoM.Translate(sx+1, sy+1)
						colorm.DrawImage(kar.Screen, res.SelectionBar, kar.GlobalColorM, kar.GlobalColorMDIO)
					}

				}
			}

			// draw crafting table result item icon
			if craftingTable.ResultSlot.ID != 0 {
				kar.GlobalColorMDIO.GeoM.Reset()
				kar.GlobalColorMDIO.GeoM.Translate(craftingTablePositionX+58, craftingTablePositionY+23)
				colorm.DrawImage(kar.Screen, res.Icon8[craftingTable.ResultSlot.ID], kar.GlobalColorM, kar.GlobalColorMDIO)

				// Draw result item quantity number
				quantity := craftingTable.ResultSlot.Quantity
				if quantity > 1 {
					itemQuantityTextDO.GeoM.Reset()
					itemQuantityTextDO.GeoM.Translate(craftingTablePositionX+58, craftingTablePositionY+22)
					num := strconv.FormatUint(uint64(quantity), 10)
					if quantity < 10 {
						num = " " + num
					}
					text.Draw(kar.Screen, num, res.Font, itemQuantityTextDO)
				}
			}

		}

		// Draw all rects for debug
		if kar.DrawDebugHitboxesEnabled {
			rectQ := arc.FilterRect.Query(&kar.WorldECS)
			for rectQ.Next() {
				rect := rectQ.Get()
				x, y := kar.Camera.ApplyCameraTransformToPoint(rect.X, rect.Y)
				vector.DrawFilledRect(
					kar.Screen,
					float32(x),
					float32(y),
					float32(rect.W),
					float32(rect.H),
					color.RGBA{128, 0, 0, 10},
					false,
				)
			}
		}

		// Draw debug info
		if kar.DrawDebugTextEnabled {
			ebitenutil.DebugPrintAt(kar.Screen, fmt.Sprintf(
				debugInfo,
				ctrl.CurrentState,
				ctrl.InputAxisLast,
			), 10, 50)
		}

	}

}
