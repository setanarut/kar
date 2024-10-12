package system

import (
	"fmt"
	"image/color"
	"kar/comp"
	"kar/items"
	"kar/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var hudTextTemplate string
var selectedIm = ebiten.NewImage(16, 16)

type DrawHUDSystem struct {
	hotbarDOP, itemsDOP *ebiten.DrawImageOptions
	itemQuantityTextDOP *text.DrawOptions
}

func NewDrawHUDSystem() *DrawHUDSystem {
	return &DrawHUDSystem{}
}
func (hs *DrawHUDSystem) Init() {
	hs.hotbarDOP = &ebiten.DrawImageOptions{}
	hs.itemsDOP = &ebiten.DrawImageOptions{}
	hs.itemQuantityTextDOP = &text.DrawOptions{}
	resources.StatsTextOptions.GeoM.Translate(30, 26)

	selectedIm.Fill(color.White)

	hudTextTemplate = `
PLAYER   %d %d
LOOKAT   %d %d %s
CHUNK    %d %d
TPS/FPS  %v %v
ENTITIES %d
SELECTED %v
`
}

func (hs *DrawHUDSystem) Update() {

}
func (hs *DrawHUDSystem) Draw(screen *ebiten.Image) {

	if player, ok := comp.TagPlayer.First(resources.ECSWorld); ok {

		slots := comp.Inventory.Get(player).Slots
		var selectedSlotDisplayName string

		if slots[resources.SelectedSlot].ID == items.Air {
			selectedSlotDisplayName = ""
		} else {
			selectedSlotDisplayName = items.Property[slots[resources.SelectedSlot].ID].DisplayName
		}

		hs.hotbarDOP.GeoM.Reset()
		hs.hotbarDOP.GeoM.Translate(-91, -11)
		hs.hotbarDOP.GeoM.Scale(2, 2)
		hs.hotbarDOP.GeoM.Translate(resources.ScreenSize.X/2, resources.ScreenSize.Y-40)
		screen.DrawImage(resources.Hotbar, hs.hotbarDOP)

		hs.hotbarDOP.GeoM.Translate(-2, -2)
		selectedOffsetX := float64(resources.SelectedSlot) * 40
		hs.hotbarDOP.GeoM.Translate(selectedOffsetX, 0)
		screen.DrawImage(resources.HotbarSelection, hs.hotbarDOP)

		for x := range 9 {
			id := slots[x].ID
			quantity := slots[x].Quantity
			var im *ebiten.Image

			if id == items.Air || quantity == 0 {
				im = nil

			} else {
				im = getSprite(id)
			}
			offsetX := (float64(x) * 40) + 320
			// öğeler
			hs.itemsDOP.GeoM.Reset()
			hs.itemsDOP.GeoM.Translate(-8, -8)
			hs.itemsDOP.GeoM.Scale(2, 2)
			hs.itemsDOP.GeoM.Translate(offsetX, resources.ScreenSize.Y-40)

			if im != nil {
				screen.DrawImage(im, hs.itemsDOP)
			}

			hs.itemQuantityTextDOP.GeoM.Reset()
			hs.itemQuantityTextDOP.GeoM.Translate(offsetX, resources.ScreenSize.Y-40)
			if quantity > 0 {
				text.Draw(screen, fmt.Sprintf("%d", quantity), resources.Font, hs.itemQuantityTextDOP)
			}
		}

		txt := fmt.Sprintf(hudTextTemplate,
			playerPosMap.X, playerPosMap.Y,
			currentBlockPosMap.X, currentBlockPosMap.Y, items.Property[HitItemID].DisplayName,
			PlayerChunk.X, PlayerChunk.X,
			int(ebiten.ActualTPS()), int(ebiten.ActualFPS()),
			resources.ECSWorld.Len(),
			selectedSlotDisplayName,
		)

		text.Draw(screen, txt, resources.Font, resources.StatsTextOptions)

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
