package kar

import (
	"kar/items"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	systems []ISystem
}

func (g *Game) Init() {
	g.systems = []ISystem{
		&Camera{},
		&Player{},
		&Enemy{},
		&Item{},
		&UI{},
		&Effects{},
		&MainMenu{},
		&Spawn{},
	}
	for _, sys := range g.systems {
		sys.Init()
	}
	ColorM.ChangeHSV(1, 0, 0.5) // BW
	TextDO.ColorScale.Scale(0.5, 0.5, 0.5, 1)
}

func (g *Game) Update() error {

	// Debug
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		InventoryRes.RandomFillAllSlots()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		DrawDebugTextEnabled = !DrawDebugTextEnabled
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		DrawItemHitboxEnabled = !DrawItemHitboxEnabled
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		DrawPlayerTileHitboxEnabled = !DrawPlayerTileHitboxEnabled
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		DataManager.SaveItem("map.png", TileMapRes.GetImageByte())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		box := MapAABB.Get(CurrentPlayer)
		TileMapRes.Set(TileMapRes.W/2, TileMapRes.H-3, items.Air)
		box.Pos.X, box.Pos.Y = TileMapRes.TileToWorldCenter(TileMapRes.W/2, TileMapRes.H-3)
		CameraRes.SetCenter(box.Pos.X, box.Pos.Y)
	}

	if ebiten.IsFocused() {
		switch CurrentGameState {
		case "menu":

			if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				if PreviousGameState == "playing" {
					PreviousGameState = "menu"
					CurrentGameState = "playing"
					ColorM.Reset()
					TextDO.ColorScale.Reset()
				}
			}
			// enter playing
			g.systems[6].Update() // MainMenu

		case "playing":
			if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				PreviousGameState = "playing"
				CurrentGameState = "menu"
				ColorM.ChangeHSV(1, 0, 0.5) // BW
				TextDO.ColorScale.Scale(0.5, 0.5, 0.5, 1)
			}
			g.systems[0].Update()
			g.systems[1].Update()
			g.systems[2].Update()
			g.systems[3].Update()
			g.systems[4].Update()
			g.systems[5].Update()
			// g.systems[6].Update()
			g.systems[7].Update()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	Screen = screen
	Screen.Fill(BackgroundColor)

	switch CurrentGameState {
	case "menu":
		g.systems[0].Draw()
		g.systems[6].Draw()
	case "playing":
		g.systems[0].Draw()
		g.systems[1].Draw()
		g.systems[2].Draw()
		g.systems[3].Draw()
		g.systems[4].Draw()
		g.systems[5].Draw()
		// g.systems[6].Draw()
		g.systems[7].Draw()
	}
}

func (g *Game) LayoutF(w, h float64) (float64, float64) {
	return ScreenW, ScreenH
}

func (g *Game) Layout(w, h int) (int, int) {
	return 0, 0
}
