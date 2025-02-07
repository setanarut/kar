package main

import (
	"kar"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := &kar.Game{}
	game.Init()
	opts := &ebiten.RunGameOptions{
		DisableHiDPI:    true,
		GraphicsLibrary: ebiten.GraphicsLibraryAuto,
		InitUnfocused:   false,
	}
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(int(kar.ScreenW*kar.WindowScale), int(kar.ScreenH*kar.WindowScale))
	if err := ebiten.RunGameWithOptions(game, opts); err != nil {
		log.Fatal(err)
	}
}
