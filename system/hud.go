package system

import (
	"fmt"
	"kar/comp"
	"kar/res"

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
func (hs *DrawHUDSystem) Draw() {

	if ebiten.IsFocused() {
		// inventory
		if false {
			p, ok := comp.PlayerTag.First(res.World)
			if ok {
				inv := comp.Inventory.Get(p)

				text.Draw(res.Screen, fmt.Sprintf("I %v | H %v", inv.Items, comp.Health.GetValue(p)), res.Futura, res.StatsTextOptions)
			} else {

				text.Draw(res.Screen, "You are dead \n Press Backspace key to restart", res.FuturaBig, res.CenterTextOptions)
			}
		}
	} else {

		// unfocused
		if true {
			text.Draw(res.Screen, "PAUSED\n Click to resume", res.FuturaBig, res.CenterTextOptions)
		}

	}

	// debug
	if true {
		text.Draw(res.Screen, fmt.Sprintf("FPS%v | TPS%v", ebiten.ActualFPS(), ebiten.ActualTPS()), res.Futura, res.StatsTextOptions)

	}
}
