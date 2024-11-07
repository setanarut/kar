package main

import (
	"kar"
	"kar/system"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	game := NewGame()
	game.Init()
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(int(kar.ScreenSize.X), int(kar.ScreenSize.Y))

	// run
	if err := ebiten.RunGameWithOptions(
		game,
		&ebiten.RunGameOptions{
			GraphicsLibrary: ebiten.GraphicsLibraryAuto,
			InitUnfocused:   false},
	); err != nil {
		log.Fatal(err)
	}

}

type Game struct {
	systems []kar.ISystem
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Init() {
	g.systems = []kar.ISystem{
		&system.Input{},
		&system.Spawn{},
		&system.Player{},
		&system.DropItem{},
		&system.States{},
		&system.CollisionHandler{},
		&system.Render{},
		&system.RenderGUI{},
	}

	// Initalize systems
	for _, s := range g.systems {
		s.Init()
	}

}

func (g *Game) Update() error {

	if ebiten.IsFocused() {
		for _, s := range g.systems {
			s.Update()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	kar.Screen = screen
	for _, s := range g.systems {
		s.Draw()
	}
}

// func (g *Game) LayoutF(w, h float64) (float64, float64) {
// 	return resources.ScreenSize.X, resources.ScreenSize.Y
// }

func (g *Game) Layout(w, h int) (int, int) {
	return int(kar.ScreenSize.X), int(kar.ScreenSize.Y)
}
