package main

import (
	"kar/engine"
	"kar/engine/cm"
	"kar/res"
	"kar/system"

	"github.com/hajimehoshi/ebiten/v2"
)

type System interface {
	Init()
	Update()
	Draw()
}

type Game struct {
	systems []System
}

func NewGame() *Game {
	return &Game{}
}

// var Start = time.Now()

func (g *Game) Init() {

	w, h := 800, 600

	res.Screen = ebiten.NewImage(w, h)
	res.ScreenRect = cm.NewBB(0, 0, float64(w), float64(h))
	res.Camera = engine.NewCamera(res.ScreenRect.Center(), res.ScreenRect.R, res.ScreenRect.T)
	res.Camera.Lerp = true
	g.systems = []System{
		system.NewSpawnSystem(),
		system.NewTimersSystem(),
		system.NewPhysicsSystem(),
		system.NewAISystem(),
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

func (g *Game) Draw(s *ebiten.Image) {

	for _, s := range g.systems {
		s.Draw()
	}
	s.DrawImage(res.Screen, nil)
}

func (g *Game) Layout(w, h int) (int, int) {
	return res.Screen.Bounds().Dx(), res.Screen.Bounds().Dy()
}
