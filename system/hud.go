package system

import (
	"fmt"
	"image/color"
	"kar/comp"
	"kar/itm"
	"kar/res"

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
	res.StatsTextOptions.GeoM.Translate(30, 26)

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

	if player, ok := comp.PlayerTag.First(res.ECSWorld); ok {

		slots := comp.Inventory.Get(player).Slots
		var selectedSlotDisplayName string

		if slots[res.SelectedSlot].ID == itm.Air {
			selectedSlotDisplayName = ""
		} else {
			selectedSlotDisplayName = itm.Items[slots[res.SelectedSlot].ID].DisplayName
		}

		hs.hotbarDOP.GeoM.Reset()
		hs.hotbarDOP.GeoM.Translate(-91, -11)
		hs.hotbarDOP.GeoM.Scale(2, 2)
		hs.hotbarDOP.GeoM.Translate(res.ScreenSize.X/2, res.ScreenSize.Y-40)
		screen.DrawImage(res.Hotbar, hs.hotbarDOP)

		hs.hotbarDOP.GeoM.Translate(-2, -2)
		selectedOffsetX := float64(res.SelectedSlot) * 40
		hs.hotbarDOP.GeoM.Translate(selectedOffsetX, 0)
		screen.DrawImage(res.HotbarSelection, hs.hotbarDOP)

		for x := range 9 {
			id := slots[x].ID
			quantity := slots[x].Quantity
			var im *ebiten.Image

			if id == itm.Air || quantity == 0 {
				im = nil
			} else {
				im = res.SpriteFrames[id][0]
			}

			offsetX := (float64(x) * 40) + 320
			// öğeler
			hs.itemsDOP.GeoM.Reset()
			hs.itemsDOP.GeoM.Translate(-8, -8)
			hs.itemsDOP.GeoM.Scale(2, 2)
			hs.itemsDOP.GeoM.Translate(offsetX, res.ScreenSize.Y-40)

			if im != nil {
				screen.DrawImage(im, hs.itemsDOP)
			}

			hs.itemQuantityTextDOP.GeoM.Reset()
			hs.itemQuantityTextDOP.GeoM.Translate(offsetX, res.ScreenSize.Y-40)
			if quantity > 0 {
				text.Draw(screen, fmt.Sprintf("%d", quantity), res.Font, hs.itemQuantityTextDOP)
			}
		}

		txt := fmt.Sprintf(hudTextTemplate,
			playerPosMap.X, playerPosMap.Y,
			currentBlockPosMap.X, currentBlockPosMap.Y, itm.Items[HitItemID].DisplayName,
			PlayerChunk.X, PlayerChunk.X,
			int(ebiten.ActualTPS()), int(ebiten.ActualFPS()),
			res.ECSWorld.Len(),
			selectedSlotDisplayName,
		)

		text.Draw(screen, txt, res.Font, res.StatsTextOptions)

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
