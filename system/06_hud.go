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

type DrawHUD struct {
	hotbarDIO, itemsDIO *eb.DrawImageOptions
	itemQuantityTextDO  *text.DrawOptions
}

func (hs *DrawHUD) Init() {
	hs.hotbarDIO = &eb.DrawImageOptions{}
	hs.itemsDIO = &eb.DrawImageOptions{}
	hs.itemQuantityTextDO = &text.DrawOptions{}
	fontSmallDrawOptions.GeoM.Translate(30, 26)

	selectedIm.Fill(color.White)

	hudTextTemplate = `
Player   %d %d
Look     %d %d %s
Chunk    %d %d
TPS/FPS  %v %v
Entities %d
HandSlot %v %v
Slot     %v
Vel      %v
`
}

func (hs *DrawHUD) Update() {
}
func (hs *DrawHUD) Draw() {
	if playerEntry.Valid() {
		playerInv := comp.Inventory.Get(playerEntry)
		// Draw hotbar
		if playerInv.Slots[selectedSlotIndex].ID == items.Air {
			selectedSlotDisplayName = ""
		} else {
			id := playerInv.Slots[selectedSlotIndex].ID
			selectedSlotDisplayName = items.Property[id].DisplayName
		}
		hs.hotbarDIO.GeoM.Reset()
		hs.hotbarDIO.GeoM.Translate(-91, -11)
		hs.hotbarDIO.GeoM.Scale(2, 2)
		hs.hotbarDIO.GeoM.Translate(kar.ScreenSize.X/2, kar.ScreenSize.Y-40)
		kar.Screen.DrawImage(res.Hotbar, hs.hotbarDIO)

		// Draw hotbar selected border
		hs.hotbarDIO.GeoM.Translate(-2, -2)
		selectedOffsetX := float64(selectedSlotIndex) * 40
		hs.hotbarDIO.GeoM.Translate(selectedOffsetX, 0)
		kar.Screen.DrawImage(res.HotbarSelection, hs.hotbarDIO)

		// Draw hotbar slots
		for x := range 9 {
			quantity := playerInv.Slots[x].Quantity
			offsetX := (float64(x) * 40) + 320
			hs.itemsDIO.GeoM.Reset()
			hs.itemsDIO.GeoM.Translate(-8, -8)
			hs.itemsDIO.GeoM.Scale(2, 2)
			hs.itemsDIO.GeoM.Translate(offsetX, kar.ScreenSize.Y-40)
			if playerInv.Slots[x].Quantity > 0 {
				kar.Screen.DrawImage(getSprite(playerInv.Slots[x].ID), hs.itemsDIO)
			}
			hs.itemQuantityTextDO.GeoM.Reset()
			hs.itemQuantityTextDO.GeoM.Translate(offsetX-8, kar.ScreenSize.Y-45)
			if quantity > 0 {
				num := strconv.FormatUint(uint64(quantity), 10)
				if quantity < 10 {
					num = " " + num
				}
				text.Draw(kar.Screen, num, res.Font, hs.itemQuantityTextDO)
			}
		}

		// Draw stats text
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
			playerVel,
		)

		text.Draw(kar.Screen, txt, res.FontSmall, fontSmallDrawOptions)

	}
}
