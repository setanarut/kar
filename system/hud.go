package system

import (
	"fmt"
	"kar/comp"
	"kar/res"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var hudTextTemplate string

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
Inventory %v
SelectedItem %v
`
	res.StatsTextOptions.GeoM.Translate(30, 26)
}

func (hs *DrawHUDSystem) Update() {

}
func (hs *DrawHUDSystem) Draw(screen *ebiten.Image) {

	if player, ok := comp.PlayerTag.First(res.ECSWorld); ok {

		txt := fmt.Sprintf(hudTextTemplate,
			playerPosMap,
			currentBlockPosMap,
			PlayerChunk,
			math.Round(ebiten.ActualTPS()),
			math.Round(ebiten.ActualFPS()),
			res.ECSWorld.Len(),
			comp.Inventory.Get(player).Items,
			res.SelectedItem,
		)
		text.Draw(screen, txt, res.Font, res.StatsTextOptions)
	}

}
