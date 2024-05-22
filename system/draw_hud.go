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
}

func (hs *DrawHUDSystem) Update() {

}
func (hs *DrawHUDSystem) Draw() {

	if ebiten.IsFocused() {
		// inventory
		if true {
			p, ok := comp.PlayerTag.First(res.World)
			if ok {
				inv := comp.Inventory.Get(p)

				text.Draw(res.Screen, fmt.Sprintf("%v", inv), res.Futura, res.StatsTextOptions)
			} else {
				res.CenterTextOptions.GeoM.Translate(res.ScreenRect.Center().X, res.ScreenRect.Center().X)
				text.Draw(res.Screen, "You are dead \n Press Backspace key to restart", res.FuturaBig, res.CenterTextOptions)
				res.CenterTextOptions.GeoM.Translate(-res.ScreenRect.Center().X, -res.ScreenRect.Center().X)
			}
		}
	} else {

		// unfocused
		if true {
			text.Draw(res.Screen, "PAUSED\n Click to resume", res.FuturaBig, res.CenterTextOptions)
		}

	}

	// debug
	if false {
		text.Draw(res.Screen, "debug", res.FuturaBig, res.CenterTextOptions)

	}
}
