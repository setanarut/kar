package kar

import (
	"image"
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
		&Enemy{},
		&Player{},
		&Item{},
		&UI{},
		&Effects{},
		&MainMenu{},
		&Spawn{},
	}
	for _, sys := range g.systems {
		sys.Init()
	}
	colorM.ChangeHSV(1, 0, 0.5) // BW
	textDO.ColorScale.Scale(0.5, 0.5, 0.5, 1)
}

func (g *Game) Update() error {

	if ebiten.IsFocused() {

		// Debug
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			inventoryRes.ClearCurrentSlot()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyK) {
			inventoryRes.RandomFillAllSlots()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyV) {
			drawDebugTextEnabled = !drawDebugTextEnabled
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			drawItemHitboxEnabled = !drawItemHitboxEnabled
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyB) {
			drawPlayerTileHitboxEnabled = !drawPlayerTileHitboxEnabled
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
			dataManager.SaveItem("map.png", tileMapRes.GetImageByte())
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
			box := mapAABB.GetUnchecked(currentPlayer)
			tileMapRes.Set(tileMapRes.W/2, tileMapRes.H-3, items.Air)
			box.Pos = tileMapRes.TileToWorld(image.Point{tileMapRes.W / 2, tileMapRes.H - 3})
			cameraRes.SetCenter(box.Pos.X, box.Pos.Y)
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
			g.systems[0].Update()
			g.systems[6].Update()
		case "playing":
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
	Screen.Fill(backgroundColor)

	switch currentGameState {
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
	return ScreenSize.X, ScreenSize.Y
}

func (g *Game) Layout(w, h int) (int, int) {
	return 0, 0
}
