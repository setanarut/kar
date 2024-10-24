package system

import (
	"fmt"
	"image/color"
	"kar"
	"kar/comp"
	"kar/items"
	"kar/res"
	"strconv"

	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	selectedSlotDisplayName string
	hudTextTemplate         string
	selectedIm              = eb.NewImage(16, 16)
)

type RenderGUI struct {
	hotbarDIO, itemsDIO *eb.DrawImageOptions
	itemQuantityTextDO  *text.DrawOptions
}

func (r *RenderGUI) Init() {
	r.hotbarDIO = &eb.DrawImageOptions{}
	r.itemsDIO = &eb.DrawImageOptions{}
	r.itemQuantityTextDO = &text.DrawOptions{}
	fontSmallDrawOptions.GeoM.Translate(4, 4)

	selectedIm.Fill(color.White)

	hudTextTemplate = `Player   %d %d
Look     %d %d %s
Chunk    %d %d
TPS/FPS  %v %v
Entities %d
HandSlot %v %v
Slot     %v
Vel      %.4f, %.4f
`
}

func (r *RenderGUI) Update() {
}
func (r *RenderGUI) Draw() {
	if playerEntry.Valid() {
		playerInv := comp.Inventory.Get(playerEntry)

		// Draw hotbar
		if playerInv.Slots[selectedSlotIndex].ID == items.Air {
			selectedSlotDisplayName = ""
		} else {
			id := playerInv.Slots[selectedSlotIndex].ID
			selectedSlotDisplayName = items.Property[id].DisplayName
		}
		r.hotbarDIO.GeoM.Reset()
		r.hotbarDIO.GeoM.Translate(-91, -11)
		r.hotbarDIO.GeoM.Scale(kar.GUIScale, kar.GUIScale)
		r.hotbarDIO.GeoM.Translate(kar.ScreenSize.X/2, kar.ScreenSize.Y-40)
		kar.Screen.DrawImage(res.Hotbar, r.hotbarDIO)

		// Draw hotbar selected border
		r.hotbarDIO.GeoM.Translate(-1, -1)
		selectedOffsetX := float64(selectedSlotIndex) * 20
		r.hotbarDIO.GeoM.Translate(selectedOffsetX, 0)
		kar.Screen.DrawImage(res.HotbarSelection, r.hotbarDIO)

		// Draw hotbar items
		for x := range 9 {

			// Draw item
			quantity := playerInv.Slots[x].Quantity
			offsetX := (float64(x) * 20) + 48
			r.itemsDIO.GeoM.Reset()
			r.itemsDIO.GeoM.Translate(-8, -8)
			r.itemsDIO.GeoM.Scale(kar.GUIScale, kar.GUIScale)
			r.itemsDIO.GeoM.Translate(offsetX, kar.ScreenSize.Y-40)
			if playerInv.Slots[x].Quantity > 0 {
				kar.Screen.DrawImage(getSprite(playerInv.Slots[x].ID), r.itemsDIO)
			}

			// Draw item quantity
			r.itemQuantityTextDO.GeoM.Reset()
			r.itemQuantityTextDO.GeoM.Translate(offsetX-3, kar.ScreenSize.Y-41)
			if quantity > 0 {
				num := strconv.FormatUint(uint64(quantity), 10)
				if quantity < 10 {
					num = " " + num
				}
				text.Draw(kar.Screen, num, res.FontSmall, r.itemQuantityTextDO)
			}
		}

		playerVel := playerBody.Velocity()

		// Draw HUD text
		txt := fmt.Sprintf(
			hudTextTemplate,
			playerPixelCoord.X,
			playerPixelCoord.Y,
			hitBlockPixelCoord.X,
			hitBlockPixelCoord.Y,
			items.Property[hitItemID].DisplayName,
			gameWorld.PlayerChunk.X,
			gameWorld.PlayerChunk.Y,
			int(eb.ActualTPS()),
			int(eb.ActualFPS()),
			ecsWorld.Len(),
			items.Property[playerInv.HandSlot.ID].DisplayName,
			playerInv.HandSlot.Quantity,
			selectedSlotDisplayName,
			playerVel.X, playerVel.Y,
		)

		text.Draw(kar.Screen, txt, res.FontSmall, fontSmallDrawOptions)

	}
}
