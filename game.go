package kar

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	systems []ISystem
}

func (g *Game) Init() {
	g.systems = []ISystem{
		&Spawn{},    // 0
		&Camera{},   // 1
		&Player{},   // 2
		&Enemy{},    // 3
		&Item{},     // 4
		&UI{},       // 5
		&MainMenu{}, // 6
	}
	for _, sys := range g.systems {
		sys.Init()
	}
	ColorM.ChangeHSV(1, 0, 0.5) // BW
	TextDO.ColorScale.Scale(0.5, 0.5, 0.5, 1)
}

func (g *Game) Update() error {
	if ebiten.IsFocused() {
		switch CurrentGameState {
		case "menu":
			if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				PreviousGameState = "menu"
				CurrentGameState = "playing"
				ColorM.Reset()
				TextDO.ColorScale.Reset()
			}
			// enter playing
			g.systems[6].Update() // MainMenu

		case "playing":
			if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				PreviousGameState = "playing"
				CurrentGameState = "menu"
				ColorM.ChangeHSV(1, 0, 0.5) // BW
				TextDO.ColorScale.Scale(0.5, 0.5, 0.5, 1)
			}
			for i := 0; i < 6; i++ { // Update all systems except MainMenu
				g.systems[i].Update()
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	Screen = screen
	Screen.Fill(BackgroundColor)

	switch CurrentGameState {
	case "menu":
		g.systems[1].Draw()
		g.systems[5].Draw()
		g.systems[6].Draw() // MainMenu
	case "playing":
		for i := 0; i < 6; i++ { // Draw all systems except MainMenu
			g.systems[i].Draw()
		}
	}
}

func (g *Game) LayoutF(w, h float64) (float64, float64) {
	return ScreenW, ScreenH
}

func (g *Game) Layout(w, h int) (int, int) {
	return 0, 0
}
