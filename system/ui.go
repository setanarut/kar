package system

import (
	"fmt"
	"image/color"
	"kar"
	"kar/arc"
	"kar/items"
	"kar/res"
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
	craftingTablePositionX = 50 + hotbarPositionX
	craftingTablePositionY = 40 + hotbarPositionY
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
		id := tileMap.TileID(targetBlockPos.X, targetBlockPos.Y)
		if id == items.CraftingTable {
			craftingTable.PosX, craftingTable.PosY = 1, 1
			craftingState = !craftingState
		}
	}

	// Crafting table slot navigation
	if craftingState {
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			craftingTable.PosX = min(craftingTable.PosX+1, 2)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			craftingTable.PosX = max(craftingTable.PosX-1, 0)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			craftingTable.PosY = min(craftingTable.PosY+1, 2)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			craftingTable.PosY = max(craftingTable.PosY-1, 0)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {

			id := ctrl.Inventory.RemoveItemFromSelectedSlot()

			if id > 0 {
				craftingTable.SetCurrentSlot(id)
			}

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
		ctrl.Inventory.ClearSelectedSlot()
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
			slotID := ctrl.Inventory.Slots[x].ItemID
			quantity := ctrl.Inventory.Slots[x].ItemQuantity
			SlotOffsetX := float64(x) * 17
			SlotOffsetX += hotbarPositionX

			// draw item icon
			kar.GlobalColorMDIO.GeoM.Reset()
			kar.GlobalColorMDIO.GeoM.Translate(SlotOffsetX+(4), hotbarPositionY+(4))
			if slotID != items.Air && ctrl.Inventory.Slots[x].ItemQuantity > 0 {
				colorm.DrawImage(kar.Screen, res.Icon8[slotID], kar.GlobalColorM, kar.GlobalColorMDIO)
			}
			if x == ctrl.Inventory.CurrentSlotIndex {
				// Draw selected slot border
				kar.GlobalColorMDIO.GeoM.Translate(-5, -5)
				colorm.DrawImage(kar.Screen, res.SelectionBar, kar.GlobalColorM, kar.GlobalColorMDIO)

				// Draw slot item display name
				if !ctrl.Inventory.IsSelectedSlotEmpty() {
					itemQuantityTextDO.GeoM.Reset()
					itemQuantityTextDO.GeoM.Translate(SlotOffsetX-1, hotbarPositionY+14)
					if items.HasTag(slotID, items.Tool) {
						text.Draw(kar.Screen, fmt.Sprintf(
							"%v\nDurability %v",
							items.Property[slotID].DisplayName,
							ctrl.Inventory.Slots[x].ItemDurability,
						), res.Font, itemQuantityTextDO)
					} else {
						text.Draw(kar.Screen, items.Property[slotID].DisplayName, res.Font, itemQuantityTextDO)
					}
				}
			}

			// Draw item quantity number
			if quantity > 0 && items.IsStackable(slotID) {
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

			// Background
			kar.GlobalColorMDIO.GeoM.Reset()
			kar.GlobalColorMDIO.GeoM.Translate(craftingTablePositionX, craftingTablePositionY)
			colorm.DrawImage(kar.Screen, res.CraftingTable, kar.GlobalColorM, kar.GlobalColorMDIO)

			for x := 0; x < 3; x++ {
				for y := 0; y < 3; y++ {
					if craftingTable.Slots[y][x] != items.Air {
						sx := craftingTablePositionX + float64(x*17)
						sy := craftingTablePositionY + float64(y*17)
						kar.GlobalColorMDIO.GeoM.Reset()
						kar.GlobalColorMDIO.GeoM.Translate(sx+5, sy+5)
						colorm.DrawImage(kar.Screen, res.Icon8[craftingTable.Slots[y][x]], kar.GlobalColorM, kar.GlobalColorMDIO)

					}

					if x == craftingTable.PosX && y == craftingTable.PosY {
						sx := craftingTablePositionX + float64(x*17)
						sy := craftingTablePositionY + float64(y*17)
						kar.GlobalColorMDIO.GeoM.Reset()
						kar.GlobalColorMDIO.GeoM.Translate(sx, sy)
						colorm.DrawImage(kar.Screen, res.SelectionBar, kar.GlobalColorM, kar.GlobalColorMDIO)
					}

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
