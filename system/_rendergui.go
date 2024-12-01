package system

import (
	"fmt"
	"kar"
	"kar/items"
	"kar/res"
	"strconv"

	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	selectedSlotDisplayName string
	hudTextTemplate         string
)

type RenderGUI struct {
	hotbarDIO, itemsDIO *eb.DrawImageOptions
	itemQuantityTextDO  *text.DrawOptions
}

func (gui *RenderGUI) Init() {
	gui.hotbarDIO = &eb.DrawImageOptions{}
	gui.itemsDIO = &eb.DrawImageOptions{}
	gui.itemQuantityTextDO = &text.DrawOptions{}
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

func (gui *RenderGUI) Update() {
}
func (gui *RenderGUI) Draw() {

	if kar.WorldECS.Alive(Mario) {
		// Draw hotbar
		if playerInv.Slots[selectedSlotIndex].ID == items.Air {
			selectedSlotDisplayName = ""
		} else {
			id := playerInv.Slots[selectedSlotIndex].ID
			selectedSlotDisplayName = items.Property[id].DisplayName
		}
		gui.hotbarDIO.GeoM.Reset()
		gui.hotbarDIO.GeoM.Translate(-91, -11)
		gui.hotbarDIO.GeoM.Scale(2, 2)
		gui.hotbarDIO.GeoM.Translate(kar.ScreenSize.X/2, kar.ScreenSize.Y-40)
		kar.Screen.DrawImage(res.Hotbar, gui.hotbarDIO)

		// Draw hotbar selected border
		gui.hotbarDIO.GeoM.Translate(-2, -2)
		selectedOffsetX := float64(selectedSlotIndex) * 40
		gui.hotbarDIO.GeoM.Translate(selectedOffsetX, 0)
		kar.Screen.DrawImage(res.HotbarSelection, gui.hotbarDIO)

		// Draw hotbar slots
		for x := range 9 {
			quantity := playerInv.Slots[x].Quantity
			offsetX := (float64(x) * 40) + 320
			gui.itemsDIO.GeoM.Reset()
			gui.itemsDIO.GeoM.Translate(-8, -8)
			gui.itemsDIO.GeoM.Scale(2, 2)
			gui.itemsDIO.GeoM.Translate(offsetX, kar.ScreenSize.Y-40)
			if playerInv.Slots[x].Quantity > 0 {
				kar.Screen.DrawImage(getSprite(playerInv.Slots[x].ID), gui.itemsDIO)
			}
			gui.itemQuantityTextDO.GeoM.Reset()
			gui.itemQuantityTextDO.GeoM.Translate(offsetX-8, kar.ScreenSize.Y-45)
			if quantity > 0 {
				num := strconv.FormatUint(uint64(quantity), 10)
				if quantity < 10 {
					num = " " + num
				}
				text.Draw(kar.Screen, num, res.Font, gui.itemQuantityTextDO)
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
			kar.WorldECS.Stats().Entities.Used,
			items.Property[playerInv.HandSlot.ID].DisplayName,
			playerInv.HandSlot.Quantity,
			selectedSlotDisplayName,
			fmt.Sprintf("X %.0f Y %.0f", vel.X, vel.Y),
		)

		ebitenutil.DebugPrintAt(kar.Screen, txt, 20, 0)

		// text.Draw(kar.Screen, txt, res.FontSmall, fontSmallDrawOptions)

	}
}
