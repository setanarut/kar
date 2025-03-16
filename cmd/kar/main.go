package main

import (
	"kar"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var game = &kar.Game{
	Spawn:      kar.Spawn{},
	Platform:   kar.Platform{},
	Enemy:      kar.Enemy{},
	Player:     kar.Player{},
	Item:       kar.Item{},
	Effects:    kar.Effects{},
	Camera:     kar.Camera{},
	Ui:         kar.UI{},
	MainMenu:   kar.MainMenu{},
	Menu:       kar.Menu{},
	Debug:      kar.Debug{},
	Projectile: kar.Projectile{},
}
var opts = &ebiten.RunGameOptions{
	DisableHiDPI:    true,
	GraphicsLibrary: ebiten.GraphicsLibraryAuto,
	InitUnfocused:   false,
}

func main() {
	game.Init()
	// ebiten.SetTPS(8)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(int(kar.ScreenSize.X*kar.WindowScale), int(kar.ScreenSize.Y*kar.WindowScale))
	if err := ebiten.RunGameWithOptions(game, opts); err != nil {
		log.Fatal(err)
	}
}
