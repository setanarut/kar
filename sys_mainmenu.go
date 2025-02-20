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
	do   *text.DrawOptions
	line int
	text string
	x, y float64
}

func (m *MainMenu) Init() {
	m.do = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions: text.LayoutOptions{
			LineSpacing: 18,
		},
	}

	m.text = "NEW GAME\nSAVE\nLOAD"
	m.x = float64((int(ScreenW) / 2) - 20)
	m.y = float64((int(ScreenH) / 2) - 30)
	m.do.ColorScale.ScaleWithColor(color.Gray{200})
}

func (m *MainMenu) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		switch m.line {
		case 0:
			NewGame()
			PreviousGameState = "menu"
			CurrentGameState = "playing"
		case 1:
			if PreviousGameState == "playing" {
				SaveGame()
				PreviousGameState = "menu"
				CurrentGameState = "playing"
			}

		case 2:
			LoadGame()
			PreviousGameState = "menu"
			CurrentGameState = "playing"
		}

		ColorM.Reset()
		TextDO.ColorScale.Reset()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		m.line = (m.line - 1 + 3) % 3

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		m.line = (m.line + 1) % 3
	}
	return nil
}
func (m *MainMenu) Draw() {
	m.do.GeoM.Reset()
	m.do.GeoM.Translate(float64(m.x), float64(m.y))

	// draw menu text
	text.Draw(Screen, m.text, res.Font, m.do)

	// draw selection box
	vector.DrawFilledRect(
		Screen,
		float32(m.x)-8,
		float32(m.y+float64(m.line*18))+5,
		3,
		7,
		color.White,
		false,
	)
}
