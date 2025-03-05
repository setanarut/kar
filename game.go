package kar

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	spawn      *Spawn
	enemy      *Enemy
	player     *Player
	item       *Item
	effects    *Effects
	camera     *Camera
	ui         *UI
	mainMenu   *MainMenu
	debug      *Debug
	projectile *Projectile
}

func (g *Game) Init() {
	g.spawn = &Spawn{}
	g.enemy = &Enemy{}
	g.player = &Player{}
	g.item = &Item{}
	g.effects = &Effects{}
	g.camera = &Camera{}
	g.ui = &UI{}
	g.mainMenu = &MainMenu{}
	g.debug = &Debug{}
	g.projectile = &Projectile{}

	g.spawn.Init()
	g.enemy.Init()
	g.player.Init()
	g.item.Init()
	g.effects.Init()
	g.camera.Init()
	g.ui.Init()
	g.mainMenu.Init()
	g.debug.Init()
	g.projectile.Init()

	colorM.ChangeHSV(1, 0, 0.5) // BW
	textDO.ColorScale.Scale(0.5, 0.5, 0.5, 1)
}

func (g *Game) Update() error {
	if ebiten.IsFocused() {

		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			if ebiten.IsKeyPressed(ebiten.KeyMeta) && ebiten.IsKeyPressed(ebiten.KeyShift) {
				debugEnabled = !debugEnabled
			}
		}

		// toggle menu
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			switch currentGameState {
			case "menu":
				currentGameState = "playing"
				previousGameState = "menu"
				colorM.Reset()
				textDO.ColorScale.Reset()
			case "playing":
				currentGameState = "menu"
				previousGameState = "playing"
				colorM.ChangeHSV(1, 0, 0.5) // BW
				textDO.ColorScale.Scale(0.5, 0.5, 0.5, 1)
			}
		}
		// Update systems
		switch currentGameState {
		case "menu":
			g.mainMenu.Update()
			if debugEnabled {
				g.debug.Update()
			}
		case "playing":
			g.spawn.Update()
			g.enemy.Update()
			g.player.Update()
			g.item.Update()
			g.effects.Update()
			g.camera.Update()
			g.ui.Update()
			g.projectile.Update()
			if debugEnabled {
				g.debug.Update()
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	Screen = screen
	Screen.Fill(backgroundColor)

	switch currentGameState {
	case "menu":
		g.camera.Draw()
		g.mainMenu.Draw()
		if debugEnabled {
			g.debug.Draw()
		}
	case "playing":
		g.spawn.Draw()
		g.enemy.Draw()
		g.player.Draw()
		g.item.Draw()
		g.effects.Draw()
		g.camera.Draw()
		g.ui.Draw()
		g.projectile.Draw()
		if debugEnabled {
			g.debug.Draw()
		}
	}
}

func (g *Game) LayoutF(w, h float64) (float64, float64) {
	return ScreenSize.X, ScreenSize.Y
}

func (g *Game) Layout(w, h int) (int, int) {
	return 0, 0
}
