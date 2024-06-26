package system

import (
	"fmt"
	"kar/comp"
	"kar/res"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Chipmunk Space draw system
type DrawHUDSystem struct {
}

func NewDrawHUDSystem() *DrawHUDSystem {
	return &DrawHUDSystem{}
}
func (hs *DrawHUDSystem) Init() {
	res.StatsTextOptions.GeoM.Translate(30, 25)
	res.CenterTextOptions.GeoM.Translate(400, 300)
}

func (hs *DrawHUDSystem) Update() {

}
func (hs *DrawHUDSystem) Draw(screen *ebiten.Image) {

	if ebiten.IsFocused() {
		// inventory
		if true {
			p, ok := comp.PlayerTag.First(res.World)
			if ok {
				// inv := comp.Inventory.Get(p)
				pos := comp.Body.Get(p).Position().Div(res.BlockSize).Point()
				fps := int(math.Round(ebiten.ActualFPS()))
				tps := int(math.Round(ebiten.ActualTPS()))
				// text.Draw(res.Screen, fmt.Sprintf("I %v | H %v", inv.Items, comp.Health.GetValue(p)), res.Futura, res.StatsTextOptions)

				text.Draw(screen,
					fmt.Sprintf(
						"FPS %v TPS %v\nPOS %v Chunk %v\nEntities %v",
						fps, tps, pos, playerChunkTemp, res.World.Len(),
					),
					res.Futura, res.StatsTextOptions)
			} else {
				text.Draw(screen, "You are dead \n Press Backspace key to restart", res.FuturaBig, res.CenterTextOptions)
			}
		}
	} else {

		// unfocused
		if true {
			text.Draw(screen, "PAUSED\n Click to resume", res.FuturaBig, res.CenterTextOptions)
		}

	}

	// debug
	if false {
		text.Draw(screen, fmt.Sprintf("FPS=%v TPS=%v", ebiten.ActualFPS(), ebiten.ActualTPS()), res.Futura, res.StatsTextOptions)

	}
}
