package main

import (
	"kar"
	"kar/system"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func main() {
	// ebiten.SetTPS(4)
	game := NewGame()
	game.Init()
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(int(kar.ScreenW*kar.WindowScale), int(kar.ScreenH*kar.WindowScale))

	// run
	if err := ebiten.RunGameWithOptions(
		game,
		&ebiten.RunGameOptions{
			DisableHiDPI:    true,
			GraphicsLibrary: ebiten.GraphicsLibraryAuto,
			InitUnfocused:   false},
	); err != nil {
		log.Fatal(err)
	}

}

func NewGame() *Game {
	return &Game{}
}

type Game struct {
	systems []kar.ISystem
}

func (g *Game) Init() {

	g.systems = []kar.ISystem{
		&system.Spawn{},    // 0
		&system.Game{},     // 1
		&system.Player{},   // 2
		&system.Enemy{},    // 3
		&system.Item{},     // 4
		&system.UI{},       // 5
		&system.MainMenu{}, // 6
	}

	// Initialize systems using a slice of systems
	for _, sys := range g.systems {
		sys.Init()
	}

	kar.ColorM.ChangeHSV(1, 0, 0.5) // BW
	kar.TextDO.ColorScale.Scale(0.5, 0.5, 0.5, 1)
}

func (g *Game) Update() error {
	if ebiten.IsFocused() {
		switch kar.CurrentGameState {
		case "menu":
			if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				kar.PreviousGameState = "menu"
				kar.CurrentGameState = "playing"
				kar.ColorM.Reset()
				kar.TextDO.ColorScale.Reset()
			}
			// enter playing
			g.systems[6].Update() // MainMenu

		case "playing":
			if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				kar.PreviousGameState = "playing"
				kar.CurrentGameState = "menu"
				kar.ColorM.ChangeHSV(1, 0, 0.5) // BW
				kar.TextDO.ColorScale.Scale(0.5, 0.5, 0.5, 1)
			}
			for i := 0; i < 6; i++ { // Update all systems except MainMenu
				g.systems[i].Update()
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	kar.Screen = screen
	kar.Screen.Fill(kar.BackgroundColor)

	switch kar.CurrentGameState {
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
	return kar.ScreenW, kar.ScreenH
}

func (g *Game) Layout(w, h int) (int, int) {
	return 0, 0
}
