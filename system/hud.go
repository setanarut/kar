package system

import (
	"fmt"
	"image/color"
	"kar/comp"
	"kar/itm"
	"kar/res"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var hudTextTemplate string
var selectedIm = ebiten.NewImage(16, 16)

type DrawHUDSystem struct {
}

func NewDrawHUDSystem() *DrawHUDSystem {
	return &DrawHUDSystem{}
}
func (hs *DrawHUDSystem) Init() {
	hudTextTemplate = `Player   %v
Selected %v
Chunk    %v
TPS/FPS  %v %v
Entities %v
SelectedSlot %v
`
	res.StatsTextOptions.GeoM.Translate(30, 26)
	selectedIm.Fill(color.White)
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

		txt := fmt.Sprintf(hudTextTemplate,
			playerPosMap,
			currentBlockPosMap,
			PlayerChunk,
			math.Round(ebiten.ActualTPS()),
			math.Round(ebiten.ActualFPS()),
			res.ECSWorld.Len(),
			selectedSlotDisplayName,
		)
		text.Draw(screen, txt, res.Font, res.StatsTextOptions)

		for x := range 10 {
			id := slots[x].ID
			quantity := slots[x].Quantity
			var im *ebiten.Image

			if id == itm.Air || quantity == 0 {
				im = res.Slot16
			} else {
				im = res.SpriteFrames[id][0]
			}

			offsetX := (float64(x) * 36) + 300
			offsetY := res.ScreenSize.Y - 100

			// gölge
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(-8, -8)
			op.GeoM.Scale(2.2, 2.2)
			op.GeoM.Translate(offsetX, float64(offsetY))
			if x == res.SelectedSlot {
				screen.DrawImage(selectedIm, op)
			} else {
				op.ColorScale.ScaleWithColor(color.Black)
				screen.DrawImage(im, op)
			}

			// öğeler
			op.ColorScale.Reset()
			op.GeoM.Reset()
			op.GeoM.Translate(-8, -8)
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(offsetX, float64(offsetY))
			screen.DrawImage(im, op)

			dop := &text.DrawOptions{}
			dop.GeoM.Translate(offsetX, offsetY)
			if quantity > 0 {
				text.Draw(screen, fmt.Sprintf("%d", quantity), res.Font, dop)
			}
		}
	}
}
