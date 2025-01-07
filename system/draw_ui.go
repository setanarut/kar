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
	hotbarPositionX     = 8.
	hotbarPositionY     = 8.
	hotbarRightEdgePosX = hotbarPositionX + float64(res.Hotbar.Bounds().Dx())
	itemQuantityTextDO  = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions: text.LayoutOptions{
			LineSpacing: 10,
		},
	}
	debugInfo = `state %v
inputAxisLast %v
`
)

type DrawUI struct{}

func (ui *DrawUI) Init() {}
func (ui *DrawUI) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		kar.DrawDebugHitboxesEnabled = !kar.DrawDebugHitboxesEnabled
	}
}
func (ui *DrawUI) Draw() {
	if kar.WorldECS.Alive(player) {
		// Draw hotbar background
		kar.GlobalColorMDIO.GeoM.Reset()
		kar.GlobalColorMDIO.GeoM.Translate(hotbarPositionX, hotbarPositionY)
		colorm.DrawImage(kar.Screen, res.Hotbar, kar.GlobalColorM, kar.GlobalColorMDIO)

		for x := range 9 {

			// Draw hotbar slot items
			slotID := ctrl.Inventory.Slots[x].ItemID
			quantity := ctrl.Inventory.Slots[x].ItemQuantity
			SlotOffsetX := float64(x) * 17
			SlotOffsetX += hotbarPositionX
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

		// Draw crafting table GUI
		if craftingState {
			kar.GlobalColorMDIO.GeoM.Reset()
			halfX := float64(res.CraftingTable.Bounds().Dx())
			halfY := float64(res.CraftingTable.Bounds().Dy())
			kar.GlobalColorMDIO.GeoM.Translate((kar.ScreenW/2 - halfX), kar.ScreenH/2-halfY)
			colorm.DrawImage(kar.Screen, res.CraftingTable, kar.GlobalColorM, kar.GlobalColorMDIO)
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
