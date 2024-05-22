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
				playerInventory := comp.Inventory.Get(p)

				if p.HasComponent(comp.Effect) {
					eff := comp.Effect.Get(p)

					res.StatsTextOptions.GeoM.Translate(250, 10)

					text.Draw(res.Screen,
						fmt.Sprintf("Remaining %s\nSpeed: %v\nCooldown: %v\nExtra Snowball: %v",
							eff.EffectTimer.RemainingSecondsString(),
							eff.AddMovementSpeed,
							eff.ShootCooldown,
							eff.ExtraSnowball,
						),
						res.FuturaBig, res.StatsTextOptions)

					res.StatsTextOptions.GeoM.Translate(-250, -10)
				}

				liv := *comp.Char.Get(p)
				text.Draw(
					res.Screen,
					fmt.Sprintf(
						"Snowballs: %d\nBombs: %d\nKeys: %v\nPower-up: %v\nHealth: %.2f\nSpeed: %v",
						playerInventory.Snowballs,
						playerInventory.Bombs,
						playerInventory.Keys,
						playerInventory.Potion,
						liv.Health,
						liv.Speed,
					),
					res.Futura,
					res.StatsTextOptions)
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

	// FPS/TPS Debug text
	if false {
		text.Draw(
			res.Screen,
			fmt.Sprintf(
				"DynamicBodies : %d\nStaticBodies : %dEntities : %d\nActualTPS : %v\nActualFPS : %v",
				len(res.Space.DynamicBodies),
				len(res.Space.StaticBodies),
				res.World.Len(),
				ebiten.ActualTPS(),
				// ebiten.ActualFPS(),
				res.Input.ArrowDirection,
			),
			res.Futura,
			res.StatsTextOptions)
	}

	if true {
		if p, ok := comp.PlayerTag.First(res.World); ok {
			c := comp.Char.Get(p)
			res.StatsTextOptions.GeoM.Translate(250, 10)
			text.Draw(res.Screen, fmt.Sprintf("%v", c.ShootCooldown), res.FuturaBig, res.StatsTextOptions)
			res.StatsTextOptions.GeoM.Translate(-250, -10)
		}

	}
}
