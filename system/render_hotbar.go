package system

import (
	"fmt"
	"kar"
	"kar/items"
	"kar/res"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var hotbarPositionX = 8.
var hotbarPositionY = 8.
var itemQuantityTextDO = &text.DrawOptions{
	DrawImageOptions: ebiten.DrawImageOptions{},
	LayoutOptions: text.LayoutOptions{
		LineSpacing: 10,
	},
}

type DrawHotbar struct {
	itemQuantityTextDO *text.DrawOptions
}

func (gui *DrawHotbar) Init() {
	gui.itemQuantityTextDO = &text.DrawOptions{}
}

func (gui *DrawHotbar) Update() {
}
func (gui *DrawHotbar) Draw() {

	if kar.WorldECS.Alive(PlayerEntity) {

		// Background
		kar.GlobalColorMDIO.GeoM.Reset()
		kar.GlobalColorMDIO.GeoM.Scale(2, 2)
		kar.GlobalColorMDIO.GeoM.Translate(hotbarPositionX, hotbarPositionY)
		colorm.DrawImage(kar.Screen, res.Hotbar, kar.GlobalColorM, kar.GlobalColorMDIO)

		for x := range 9 {
			// draw item
			slotID := CTRL.Inventory.Slots[x].ID
			quantity := CTRL.Inventory.Slots[x].Quantity
			SlotOffsetX := float64(x) * 34
			SlotOffsetX += hotbarPositionX
			kar.GlobalColorMDIO.GeoM.Reset()
			kar.GlobalColorMDIO.GeoM.Scale(2, 2)
			kar.GlobalColorMDIO.GeoM.Translate(SlotOffsetX+8, hotbarPositionY+8)
			if slotID != items.Air && CTRL.Inventory.Slots[x].Quantity > 0 {
				colorm.DrawImage(kar.Screen, res.Icon8[slotID], kar.GlobalColorM, kar.GlobalColorMDIO)
			}

			if x == CTRL.Inventory.SelectedSlotIndex {
				// Draw border
				kar.GlobalColorMDIO.GeoM.Translate(-10, -10)
				colorm.DrawImage(kar.Screen, res.SelectionBar, kar.GlobalColorM, kar.GlobalColorMDIO)
				if !CTRL.Inventory.IsSelectedSlotEmpty() {
					// Draw display name
					itemQuantityTextDO.GeoM.Reset()
					itemQuantityTextDO.GeoM.Scale(2, 2)
					itemQuantityTextDO.GeoM.Translate(SlotOffsetX-2, hotbarPositionY+28)
					if items.HasTag(slotID, items.Tool) {
						text.Draw(kar.Screen, fmt.Sprintf(
							"%v\nDurability %v",
							items.Property[slotID].DisplayName,
							CTRL.Inventory.Slots[x].Durability,
						), res.Font, itemQuantityTextDO)
					} else {
						text.Draw(kar.Screen, items.Property[slotID].DisplayName, res.Font, itemQuantityTextDO)
					}
				}
			}
			// Draw item quantity number
			if quantity > 0 && items.IsStackable(slotID) {
				itemQuantityTextDO.GeoM.Reset()
				itemQuantityTextDO.GeoM.Scale(2, 2)
				itemQuantityTextDO.GeoM.Translate(SlotOffsetX+12, hotbarPositionY+6)
				num := strconv.FormatUint(uint64(quantity), 10)
				if quantity < 10 {
					num = " " + num
				}
				text.Draw(kar.Screen, num, res.Font, itemQuantityTextDO)
			}
		}

		// Draw player health text
		itemQuantityTextDO.GeoM.Reset()
		itemQuantityTextDO.GeoM.Scale(2, 2)
		itemQuantityTextDO.GeoM.Translate(400, 16)
		text.Draw(kar.Screen, fmt.Sprintf("Health %v", CTRL.Health.Health), res.Font, itemQuantityTextDO)

		if craftingState {
			kar.GlobalColorMDIO.GeoM.Reset()
			kar.GlobalColorMDIO.GeoM.Scale(2, 2)
			halfX := float64(res.CraftingTable.Bounds().Dx())
			halfY := float64(res.CraftingTable.Bounds().Dy())
			kar.GlobalColorMDIO.GeoM.Translate((kar.ScreenW/2 - halfX), kar.ScreenH/2-halfY)
			colorm.DrawImage(kar.Screen, res.CraftingTable, kar.GlobalColorM, kar.GlobalColorMDIO)
		}

		// Draw debug info
		if kar.DrawDebugTextEnabled {
			ebitenutil.DebugPrintAt(kar.Screen, fmt.Sprintf(
				stats,
				CTRL.CurrentState,
				CTRL.InputAxisLast,
			), 10, 50)
		}

	}
}

var stats = `state %v
inputAxisLast %v
`
