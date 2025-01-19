package system

import (
	"fmt"
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
)

var (
	hotbarPositionX        = kar.ScreenW/2 - float64(res.Hotbar.Bounds().Dx())/2
	hotbarPositionY        = 8.
	hotbarRightEdgePosX    = hotbarPositionX + float64(res.Hotbar.Bounds().Dx())
	craftingTablePositionX = hotbarPositionX + 49
	craftingTablePositionY = hotbarPositionY + 39
	TextDO                 = &text.DrawOptions{
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
		if tileMap.Get(targetTile.X, targetTile.Y) == items.CraftingTable {
			craftingState4 = false
		} else {
			craftingState4 = true
		}

		craftingTable.SlotPosX = 1
		craftingTable.SlotPosY = 1

		// clear crafting table when exit
		if craftingState {
			for y := range 3 {
				for x := range 3 {
					itemID := craftingTable.Slots[y][x].ID
					if itemID != 0 {
						quantity := craftingTable.Slots[y][x].Quantity
						for range quantity {
							if ctrl.Inventory.AddItemIfEmpty(
								craftingTable.Slots[y][x].ID,
								craftingTable.Slots[y][x].Durability,
							) {
								craftingTable.RemoveItem(x, y)
							} else {
								craftingTable.RemoveItem(x, y)
								arc.SpawnItem(arc.SpawnData{
									X:          playerCenterX,
									Y:          playerCenterY,
									Id:         itemID,
									Durability: craftingTable.Slots[y][x].Durability,
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

	// Crafting table slot navigation
	if craftingState {
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			if craftingState4 {
				craftingTable.SlotPosX = min(craftingTable.SlotPosX+1, 1)
			} else {
				craftingTable.SlotPosX = min(craftingTable.SlotPosX+1, 2)
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			craftingTable.SlotPosX = max(craftingTable.SlotPosX-1, 0)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			if craftingState4 {
				craftingTable.SlotPosY = min(craftingTable.SlotPosY+1, 1)
			} else {
				craftingTable.SlotPosY = min(craftingTable.SlotPosY+1, 2)
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			craftingTable.SlotPosY = max(craftingTable.SlotPosY-1, 0)
		}

		// move items from hotbar to crafting table
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			cs := craftingTable.CurrentSlot()
			if ctrl.Inventory.CurrentSlotID() != 0 {
				if craftingTable.CurrentSlot().ID == 0 {
					id, dur := ctrl.Inventory.RemoveItemFromSelectedSlot()
					cs.ID = id
					cs.Durability = dur
					cs.Quantity = 1
				} else if cs.ID == ctrl.Inventory.CurrentSlotID() {
					ctrl.Inventory.RemoveItemFromSelectedSlot()
					cs.Quantity++
				}
			}
			craftingTable.UpdateResultSlot()
		}
		// move items from crafting table to hotbar
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			cs := craftingTable.CurrentSlot()
			if cs.ID != 0 {
				if cs.Quantity == 1 {
					if ctrl.Inventory.AddItemIfEmpty(cs.ID, cs.Durability) {
						craftingTable.ClearCurrenSlot()
					}
				} else if cs.Quantity > 1 {
					if ctrl.Inventory.AddItemIfEmpty(cs.ID, cs.Durability) {
						cs.Quantity--

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
		kar.ColorMDIO.GeoM.Reset()
		kar.ColorMDIO.GeoM.Translate(hotbarPositionX, hotbarPositionY)
		colorm.DrawImage(kar.Screen, res.Hotbar, kar.ColorM, kar.ColorMDIO)

		// Draw slots
		for x := range 9 {
			slotID := ctrl.Inventory.Slots[x].ID
			quantity := ctrl.Inventory.Slots[x].Quantity
			SlotOffsetX := float64(x) * 17
			SlotOffsetX += hotbarPositionX

			// draw hotbar item icons
			kar.ColorMDIO.GeoM.Reset()
			kar.ColorMDIO.GeoM.Translate(SlotOffsetX+(5), hotbarPositionY+(5))
			if slotID != items.Air && ctrl.Inventory.Slots[x].Quantity > 0 {
				colorm.DrawImage(kar.Screen, res.Icon8[slotID], kar.ColorM, kar.ColorMDIO)
			}
			if x == ctrl.Inventory.CurrentSlotIndex {
				// Draw hotbar selected slot border
				kar.ColorMDIO.GeoM.Translate(-5, -5)
				colorm.DrawImage(kar.Screen, res.SelectionBar, kar.ColorM, kar.ColorMDIO)

				// Draw hotbar slot item display name
				if !ctrl.Inventory.IsCurrentSlotEmpty() {
					TextDO.GeoM.Reset()
					TextDO.GeoM.Translate(SlotOffsetX-1, hotbarPositionY+14)
					if items.HasTag(slotID, items.Tool) {
						text.Draw(kar.Screen, fmt.Sprintf(
							"%v\nDurability %v",
							items.Property[slotID].DisplayName,
							ctrl.Inventory.Slots[x].Durability,
						), res.Font, TextDO)
					} else {
						text.Draw(kar.Screen, items.Property[slotID].DisplayName, res.Font, TextDO)
					}
				}
			}

			// Draw item quantity number
			if quantity > 1 && items.IsStackable(slotID) {
				TextDO.GeoM.Reset()
				TextDO.GeoM.Translate(SlotOffsetX+6, hotbarPositionY+4)
				num := strconv.FormatUint(uint64(quantity), 10)
				if quantity < 10 {
					num = " " + num
				}
				text.Draw(kar.Screen, num, res.Font, TextDO)
			}
		}

		// Draw player health text
		TextDO.GeoM.Reset()
		TextDO.GeoM.Translate(hotbarRightEdgePosX+8, hotbarPositionY)
		text.Draw(kar.Screen, fmt.Sprintf("Health %v", ctrl.Health.Health), res.Font, TextDO)

		// Draw crafting table
		if craftingState {

			// crafting table Background
			kar.ColorMDIO.GeoM.Reset()
			kar.ColorMDIO.GeoM.Translate(craftingTablePositionX, craftingTablePositionY)

			if craftingState4 {
				colorm.DrawImage(kar.Screen, res.CraftingTable4, kar.ColorM, kar.ColorMDIO)
			} else {
				colorm.DrawImage(kar.Screen, res.CraftingTable, kar.ColorM, kar.ColorMDIO)
			}

			// draw crafting table item icons
			for x := 0; x < 3; x++ {
				for y := 0; y < 3; y++ {
					if craftingTable.Slots[y][x].ID != items.Air {
						sx := craftingTablePositionX + float64(x*17)
						sy := craftingTablePositionY + float64(y*17)
						kar.ColorMDIO.GeoM.Reset()
						kar.ColorMDIO.GeoM.Translate(sx+6, sy+6)
						colorm.DrawImage(
							kar.Screen,
							res.Icon8[craftingTable.Slots[y][x].ID],
							kar.ColorM,
							kar.ColorMDIO,
						)

						// Draw item quantity number
						quantity := craftingTable.Slots[y][x].Quantity
						if quantity > 1 {
							TextDO.GeoM.Reset()
							TextDO.GeoM.Translate(sx+7, sy+5)
							num := strconv.FormatUint(uint64(quantity), 10)
							if quantity < 10 {
								num = " " + num
							}
							text.Draw(kar.Screen, num, res.Font, TextDO)
						}
					}

					// draw selected slot border of crqfting table
					if x == craftingTable.SlotPosX && y == craftingTable.SlotPosY {
						sx := craftingTablePositionX + float64(x*17)
						sy := craftingTablePositionY + float64(y*17)
						kar.ColorMDIO.GeoM.Reset()
						kar.ColorMDIO.GeoM.Translate(sx+1, sy+1)
						colorm.DrawImage(kar.Screen, res.SelectionBar, kar.ColorM, kar.ColorMDIO)
					}

				}
			}

			// draw crafting table result item icon
			if craftingTable.ResultSlot.ID != 0 {
				kar.ColorMDIO.GeoM.Reset()

				if craftingState4 {
					kar.ColorMDIO.GeoM.Translate(craftingTablePositionX+41, craftingTablePositionY+14)
				} else {
					kar.ColorMDIO.GeoM.Translate(craftingTablePositionX+58, craftingTablePositionY+23)
				}

				colorm.DrawImage(kar.Screen, res.Icon8[craftingTable.ResultSlot.ID], kar.ColorM, kar.ColorMDIO)

				// Draw result item quantity number
				quantity := craftingTable.ResultSlot.Quantity
				if quantity > 1 {
					TextDO.GeoM.Reset()
					if craftingState4 {
						TextDO.GeoM.Translate(craftingTablePositionX+42, craftingTablePositionY+13)
					} else {
						TextDO.GeoM.Translate(craftingTablePositionX+58, craftingTablePositionY+22)
					}
					num := strconv.FormatUint(uint64(quantity), 10)
					if quantity < 10 {
						num = " " + num
					}
					text.Draw(kar.Screen, num, res.Font, TextDO)
				}
			}

		}

		// Draw debug info
		if kar.DrawDebugTextEnabled {
			ebitenutil.DebugPrintAt(kar.Screen, fmt.Sprintf(
				debugInfo,
				ctrl.CurrentState,
				ctrl.AxisLast,
			), 10, 50)
		}

	}
}
