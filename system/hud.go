package system

import (
	"fmt"
	"kar/comp"
	"kar/items"
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
Inventory %v%v%v%v
SelectedItem %v
`
	res.StatsTextOptions.GeoM.Translate(30, 26)
}

func (hs *DrawHUDSystem) Update() {

}
func (hs *DrawHUDSystem) Draw(screen *ebiten.Image) {

	if player, ok := comp.PlayerTag.First(res.ECSWorld); ok {
		slots := comp.Inventory.Get(player).Slots
		txt := fmt.Sprintf(hudTextTemplate,
			playerPosMap,
			currentBlockPosMap,
			PlayerChunk,
			math.Round(ebiten.ActualTPS()),
			math.Round(ebiten.ActualFPS()),
			res.ECSWorld.Len(),
			slots[0],
			slots[1],
			slots[2],
			slots[3],
			res.SelectedItem,
		)
		text.Draw(screen, txt, res.Font, res.StatsTextOptions)
		if res.SelectedItem != items.Air {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(-8, -8)
			op.GeoM.Scale(3, 3)
			op.GeoM.Translate(res.ScreenSize.X/2, res.ScreenSize.Y-64)
			screen.DrawImage(res.SpriteFrames[res.SelectedItem][0], op)
		}
	}

}
