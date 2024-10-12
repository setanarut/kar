package main

import (
	"kar/resources"
	"kar/system"
	"kar/types"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	game := NewGame()
	game.Init()
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(int(resources.ScreenSize.X), int(resources.ScreenSize.Y))

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
	systems []types.ISystem
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Init() {
	g.systems = []types.ISystem{
		system.NewSpawnSystem(),
		system.NewPhysicsSystem(),
		system.NewPlayerControlSystem(),
		system.NewDestroySystem(),
		system.NewDrawCameraSystem(),
		system.NewDrawHUDSystem(),
		system.NewTimersSystem(),
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

// dasdad
func (g *Game) Draw(screen *ebiten.Image) {
	for _, s := range g.systems {
		s.Draw(screen)
	}
}

// func (g *Game) LayoutF(w, h float64) (float64, float64) {
// 	return resources.ScreenSize.X, resources.ScreenSize.Y
// }

func (g *Game) Layout(w, h int) (int, int) {
	return int(resources.ScreenSize.X), int(resources.ScreenSize.Y)
}
