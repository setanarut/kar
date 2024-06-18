package main

import (
	"kar/res"
	"kar/system"
	"kar/types"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := NewGame()
	g.Init()
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(res.Screen.Bounds().Dx(), res.Screen.Bounds().Dy())
	gOpt := &ebiten.RunGameOptions{GraphicsLibrary: ebiten.GraphicsLibraryAuto, InitUnfocused: true}
	if err := ebiten.RunGameWithOptions(g, gOpt); err != nil {
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
	res.Screen = ebiten.NewImage(res.ScreenSize.X, res.ScreenSize.Y)
	g.systems = []types.ISystem{
		system.NewSpawnSystem(),
		// system.NewTimersSystem(),
		// system.NewAISystem(),
		system.NewPhysicsSystem(),
		system.NewPlayerControlSystem(),
		system.NewDrawCameraSystem(),
		system.NewDrawHUDSystem(),
		system.NewDestroySystem(),
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
	for _, s := range g.systems {
		s.Draw(screen)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return res.Screen.Bounds().Dx(), res.Screen.Bounds().Dy()
}
