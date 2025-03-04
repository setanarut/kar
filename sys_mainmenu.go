package kar

import (
	"image/color"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type MainMenu struct {
	drawOpt    *text.DrawOptions
	line       int
	text       string
	menuOffset Vec
}

func (m *MainMenu) Init() {
	m.drawOpt = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions: text.LayoutOptions{
			LineSpacing: 18,
		},
	}

	m.text = "NEW GAME\nSAVE\nLOAD"
	m.menuOffset = ScreenSize.Scale(0.5).Sub(Vec{20, 30})
	m.drawOpt.ColorScale.ScaleWithColor(color.Gray{200})
}

func (m *MainMenu) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		switch m.line {
		case 0:
			NewGame()
			previousGameState = "menu"
			currentGameState = "playing"
		case 1:
			if previousGameState == "playing" {
				SaveGame() // TODO oyun yokken menude iken kayıt tuşu olmamalı
				previousGameState = "menu"
				currentGameState = "playing"
			}

		case 2:
			LoadGame()
			previousGameState = "menu"
			currentGameState = "playing"
		}

		colorM.Reset()
		textDO.ColorScale.Reset()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		m.line = (m.line - 1 + 3) % 3

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		m.line = (m.line + 1) % 3
	}
}
func (m *MainMenu) Draw() {
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
