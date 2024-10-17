package system

import (
	"fmt"
	"image/color"
	"kar/comp"
	"kar/items"
	"kar/res"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var selectedSlotDisplayName string
var hudTextTemplate string
var selectedIm = ebiten.NewImage(16, 16)

type DrawHUDSystem struct {
	hotbarDIO, itemsDIO *res.DIO
	itemQuantityTextDO  *text.DrawOptions
}

func (hs *DrawHUDSystem) Init() {
	hs.hotbarDIO = &res.DIO{}
	hs.itemsDIO = &res.DIO{}
	hs.itemQuantityTextDO = &text.DrawOptions{}
	res.FontDrawOptions.GeoM.Translate(30, 26)

	selectedIm.Fill(color.White)

	hudTextTemplate = `
Player   %d %d
Look     %d %d %s
Chunk    %d %d
TPS/FPS  %v %v
Entities %d
HandSlot %v %v
Slot     %v
`
}

func (hs *DrawHUDSystem) Update() {

}
func (hs *DrawHUDSystem) Draw(screen *ebiten.Image) {

	if player, ok := comp.TagPlayer.First(res.ECSWorld); ok {
		playerInv := comp.Inventory.Get(player)
		// Draw hotbar
		if playerInv.Slots[res.SelectedSlotIndex].ID == items.Air {
			selectedSlotDisplayName = ""
		} else {
			selectedSlotDisplayName = items.Property[playerInv.Slots[res.SelectedSlotIndex].ID].DisplayName
		}
		hs.hotbarDIO.GeoM.Reset()
		hs.hotbarDIO.GeoM.Translate(-91, -11)
		hs.hotbarDIO.GeoM.Scale(2, 2)
		hs.hotbarDIO.GeoM.Translate(res.ScreenSize.X/2, res.ScreenSize.Y-40)
		screen.DrawImage(res.Hotbar, hs.hotbarDIO)

		// Draw hotbar selected border
		hs.hotbarDIO.GeoM.Translate(-2, -2)
		selectedOffsetX := float64(res.SelectedSlotIndex) * 40
		hs.hotbarDIO.GeoM.Translate(selectedOffsetX, 0)
		screen.DrawImage(res.HotbarSelection, hs.hotbarDIO)

		// Draw hotbar slots
		for x := range 9 {
			quantity := playerInv.Slots[x].Quantity
			offsetX := (float64(x) * 40) + 320
			hs.itemsDIO.GeoM.Reset()
			hs.itemsDIO.GeoM.Translate(-8, -8)
			hs.itemsDIO.GeoM.Scale(2, 2)
			hs.itemsDIO.GeoM.Translate(offsetX, res.ScreenSize.Y-40)
			if playerInv.Slots[x].Quantity > 0 {
				screen.DrawImage(getSprite(playerInv.Slots[x].ID), hs.itemsDIO)
			}
			hs.itemQuantityTextDO.GeoM.Reset()
			hs.itemQuantityTextDO.GeoM.Translate(offsetX-8, res.ScreenSize.Y-45)
			if quantity > 0 {
				num := strconv.FormatUint(uint64(quantity), 10)
				if quantity < 10 {
					num = " " + num
				}
				text.Draw(screen, num, res.Font, hs.itemQuantityTextDO)
			}
		}

		// Draw stats text
		txt := fmt.Sprintf(hudTextTemplate,
			playerPosMap.X, playerPosMap.Y,
			currentBlockPosMap.X, currentBlockPosMap.Y,
			items.Property[hitItemID].DisplayName,
			playerChunk.X, playerChunk.X,
			int(ebiten.ActualTPS()), int(ebiten.ActualFPS()),
			res.ECSWorld.Len(),
			items.Property[playerInv.HandSlot.ID].DisplayName, playerInv.HandSlot.Quantity,
			selectedSlotDisplayName,
		)

		text.Draw(screen, txt, res.Font, res.FontDrawOptions)

	}
}

// hudTextTemplate = `
// PLAYER   %d %d
// SELECTED %d %d %s
// CHUNK    %d %d
// TPS/FPS  %d %d
// ENTITIES %d
// SLOT     %s
// `
