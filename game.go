package kar

import (
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Spawn      *Spawn
	Enemy      *Enemy
	Player     *Player
	Item       *Item
	Effects    *Effects
	Camera     *Camera
	Ui         *UI
	MainMenu   *MainMenu
	Menu       *Menu
	Debug      *Debug
	Projectile *Projectile
}

func (g *Game) Init() {
	g.Spawn = &Spawn{}
	g.Enemy = &Enemy{}
	g.Player = &Player{}
	g.Item = &Item{}
	g.Effects = &Effects{}
	g.Camera = &Camera{}
	g.Ui = &UI{}
	g.MainMenu = &MainMenu{}
	g.Menu = &Menu{}
	g.Debug = &Debug{}
	g.Projectile = &Projectile{}

	// Initalize systems
	val := reflect.ValueOf(g).Elem()
	for i := range val.NumField() {
		val.Field(i).MethodByName("Init").Call(nil)
	}

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
				// previousGameState = "menu"
				colorM.Reset()
				textDO.ColorScale.Reset()
			case "playing":
				currentGameState = "menu"
				// previousGameState = "playing"
				colorM.ChangeHSV(1, 0, 0.5) // BW
				textDO.ColorScale.Scale(0.5, 0.5, 0.5, 1)
			}
		}
		// Update systems
		switch currentGameState {
		case "mainmenu":
			g.MainMenu.Update()
			if debugEnabled {
				g.Debug.Update()
			}
		case "menu":
			g.Menu.Update()
			if debugEnabled {
				g.Debug.Update()
			}
		case "playing":
			g.Camera.Update()
			g.Player.Update()
			g.Enemy.Update()
			g.Item.Update()
			g.Effects.Update()
			g.Ui.Update()
			g.Projectile.Update()
			g.Spawn.Update()
			if debugEnabled {
				g.Debug.Update()
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	Screen = screen
	Screen.Fill(backgroundColor)

	switch currentGameState {
	case "mainmenu":
		g.MainMenu.Draw()
		if debugEnabled {
			g.Debug.Draw()
		}
	case "menu":
		g.Camera.Draw()
		g.Menu.Draw()
		if debugEnabled {
			g.Debug.Draw()
		}
	case "playing":
		g.Camera.Draw()
		g.Player.Draw()
		g.Enemy.Draw()
		g.Item.Draw()
		g.Projectile.Draw()
		g.Effects.Draw()
		g.Ui.Draw()
		if debugEnabled {
			g.Debug.Draw()
		}
	}
}

func (g *Game) LayoutF(w, h float64) (float64, float64) {
	return ScreenSize.X, ScreenSize.Y
}

func (g *Game) Layout(w, h int) (int, int) {
	return 0, 0
}
