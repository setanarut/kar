package main

import (
	"kar"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := &kar.Game{}
	game.Init()
	// ebiten.SetTPS(8)
	opts := &ebiten.RunGameOptions{
		DisableHiDPI:    true,
		GraphicsLibrary: ebiten.GraphicsLibraryAuto,
		InitUnfocused:   false,
	}
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(int(kar.ScreenSize.X*kar.WindowScale), int(kar.ScreenSize.Y*kar.WindowScale))
	if err := ebiten.RunGameWithOptions(game, opts); err != nil {
		log.Fatal(err)
	}
}
