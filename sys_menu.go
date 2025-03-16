package kar

import (
	"image/color"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Menu struct {
	drawOpt    *text.DrawOptions
	line       int
	text       string
	menuOffset Vec
}

func (m *Menu) Init() {
	m.drawOpt = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions: text.LayoutOptions{
			LineSpacing: 18,
		},
	}

	m.text = "SAVE\nMAIN MENU"
	m.menuOffset = ScreenSize.Scale(0.5).Sub(Vec{20, 30})
	m.drawOpt.ColorScale.ScaleWithColor(color.Gray{200})
}

func (m *Menu) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		switch m.line {
		case 0:
			SaveGame()
			// previousGameState = "menu"
			currentGameState = "playing"
			colorM = colorm.ColorM{}
			textDO.ColorScale.Reset()
		case 1:
			// previousGameState = "menu"
			currentGameState = "mainmenu"
			colorM = colorm.ColorM{}
			textDO.ColorScale.Reset()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		m.line = (m.line - 1 + 2) % 2

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		m.line = (m.line + 1) % 2
	}
}
func (m *Menu) Draw() {
	m.drawOpt.GeoM.Reset()
	m.drawOpt.GeoM.Translate(m.menuOffset.X, m.menuOffset.Y)

	// draw menu text
	text.Draw(Screen, m.text, res.Font, m.drawOpt)

	// draw selection box
	vector.DrawFilledRect(
		Screen,
		float32(m.menuOffset.X-8),
		float32(m.menuOffset.Y+float64(m.line*18))+5,
		3,
		7,
		color.White,
		false,
	)
}
