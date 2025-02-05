package system

import (
	"image/color"
	"kar"
	"kar/arc"
	"kar/res"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	archeserde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/quasilyte/gdata"
	"github.com/setanarut/anim"
	"github.com/setanarut/fastnoise"
	"github.com/setanarut/kamera/v2"
)

type MainMenu struct {
	serdeOpt    archeserde.Option
	do          *text.DrawOptions
	line        int
	text        string
	x, y        float64
	dataManager *gdata.Manager
}

func (m *MainMenu) Init() {

	m.newGame()

	m.do = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions: text.LayoutOptions{
			LineSpacing: 18,
		},
	}

	m.serdeOpt = archeserde.Opts.SkipComponents(generic.T[anim.AnimationPlayer]())
	m.text = "SAVE\nLOAD\nNEW GAME"
	m.x = float64((int(kar.ScreenW) / 2) - 10)
	m.y = float64((int(kar.ScreenH) / 2) - 20)
	m.do.ColorScale.ScaleWithColor(color.Gray{200})
	var err error
	m.dataManager, err = gdata.Open(gdata.Config{AppName: "kar"})
	if err != nil {
		panic(err)
	}

}

func (m *MainMenu) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		switch m.line {
		case 0:
			m.saveGame()
		case 1:
			m.loadGame()
		case 2:
			m.newGame()
		}
		kar.PreviousGameState = "menu"
		kar.CurrentGameState = "playing"
		kar.ColorM.Reset()
		kar.TextDO.ColorScale.Reset()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		// m.line = max(m.line-1, 0)
		m.line = (m.line - 1 + 3) % 3

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		// m.line = min(m.line+1, 2)
		m.line = (m.line + 1) % 3
	}

}
func (m *MainMenu) Draw() {

	// vector.DrawFilledRect(
	// 	kar.Screen,
	// 	float32(m.x-12),
	// 	float32(m.y),
	// 	60,
	// 	50,
	// 	color.Black,
	// 	false,
	// )

	m.do.GeoM.Reset()
	m.do.GeoM.Translate(float64(m.x), float64(m.y))
	text.Draw(kar.Screen, m.text, res.Font, m.do)

	vector.DrawFilledRect(
		kar.Screen,
		float32(m.x)-8,
		float32(m.y+float64(m.line*18))+5,
		3,
		7,
		color.White,
		false,
	)

}
func (*MainMenu) newGame() {
	kar.ECWorld.Reset()
	kar.InventoryRes.Reset()
	ecs.AddResource(&kar.ECWorld, kar.GameDataRes)
	ecs.AddResource(&kar.ECWorld, kar.InventoryRes)
	ecs.AddResource(&kar.ECWorld, kar.CraftingTableRes)
	ecs.AddResource(&kar.ECWorld, kar.AnimPlayerDataRes)
	ecs.AddResource(&kar.ECWorld, kar.CameraRes)
	ecs.AddResource(&kar.ECWorld, kar.TileMapRes)

	kar.GameTileMapGenerator.Opts.HighestSurfaceLevel = 10
	kar.GameTileMapGenerator.Opts.LowestSurfaceLevel = 30
	kar.GameTileMapGenerator.SetSeed(rand.Int())
	kar.GameTileMapGenerator.NoiseState.FractalType(fastnoise.FractalFBm)
	kar.GameTileMapGenerator.NoiseState.Frequency = 0.01
	kar.GameTileMapGenerator.Generate()

	// ctrl.Collider.StaticCheck = true
	x, y := kar.TileMapRes.FindSpawnPosition()
	// tileMap.Set(x, y+2, items.CraftingTable)
	SpawnX, SpawnY := kar.TileMapRes.TileToWorldCenter(x, y)
	kar.CameraRes.SmoothType = kamera.None
	kar.CameraRes.SetCenter(SpawnX, SpawnY)
	kar.CameraRes.SmoothOptions.LerpSpeedX = 0.5
	kar.CameraRes.SmoothOptions.LerpSpeedY = 0.05
	kar.CameraRes.SetTopLeft(kar.TileMapRes.FloorToBlockCenter(kar.CameraRes.X, kar.CameraRes.Y))
	kar.CurrentPlayer = arc.SpawnPlayer(SpawnX, SpawnY)
	kar.CameraRes.SmoothType = kamera.Lerp
}

func (m *MainMenu) saveGame() {
	kar.FetchAnimPlayerData(kar.PlayerAnimPlayer, kar.AnimPlayerDataRes)
	jsonData, err := archeserde.Serialize(&kar.ECWorld, m.serdeOpt)
	if err != nil {
		log.Fatal("Error serializing world:", err)
	}
	m.dataManager.SaveItem("01save", jsonData)
}

func (m *MainMenu) loadGame() {
	if m.dataManager.ItemExists("01save") {
		kar.ECWorld.Reset()
		ecs.AddResource(&kar.ECWorld, kar.GameDataRes)
		ecs.AddResource(&kar.ECWorld, kar.InventoryRes)
		ecs.AddResource(&kar.ECWorld, kar.CraftingTableRes)
		ecs.AddResource(&kar.ECWorld, kar.AnimPlayerDataRes)
		ecs.AddResource(&kar.ECWorld, kar.CameraRes)
		ecs.AddResource(&kar.ECWorld, kar.TileMapRes)

		jsonData, err := m.dataManager.LoadItem("01save")
		if err != nil {
			log.Fatal("Error loading saved data:", err)
		}

		err = archeserde.Deserialize(jsonData, &kar.ECWorld)
		if err != nil {
			log.Fatal("Error deserializing world:", err)
		}
		animData := ecs.GetResource[kar.AnimPlayerData](&kar.ECWorld)
		kar.SetAnimPlayerData(kar.PlayerAnimPlayer, animData)
	}
}
