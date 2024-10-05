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
	ebiten.SetWindowSize(int(res.ScreenSize.X), int(res.ScreenSize.Y))
	gOpt := &ebiten.RunGameOptions{
		GraphicsLibrary: ebiten.GraphicsLibraryAuto,
		InitUnfocused:   false}
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
// 	return res.ScreenSize.X, res.ScreenSize.Y
// }

func (g *Game) Layout(w, h int) (int, int) {
	return int(res.ScreenSize.X), int(res.ScreenSize.Y)
}
