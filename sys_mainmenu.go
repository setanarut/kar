package kar

import (
	"image/color"
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
	// m.newGame()
	m.do = &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
		LayoutOptions: text.LayoutOptions{
			LineSpacing: 18,
		},
	}

	m.serdeOpt = archeserde.Opts.SkipComponents(generic.T[anim.AnimationPlayer]())
	m.text = "SAVE\nLOAD\nNEW GAME"
	m.x = float64((int(ScreenW) / 2) - 10)
	m.y = float64((int(ScreenH) / 2) - 20)
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
		PreviousGameState = "menu"
		CurrentGameState = "playing"
		ColorM.Reset()
		TextDO.ColorScale.Reset()
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
	// 	Screen,
	// 	float32(m.x-12),
	// 	float32(m.y),
	// 	60,
	// 	50,
	// 	color.Black,
	// 	false,
	// )

	m.do.GeoM.Reset()
	m.do.GeoM.Translate(float64(m.x), float64(m.y))
	text.Draw(Screen, m.text, res.Font, m.do)

	vector.DrawFilledRect(
		Screen,
		float32(m.x)-8,
		float32(m.y+float64(m.line*18))+5,
		3,
		7,
		color.White,
		false,
	)

}
func (*MainMenu) newGame() {
	ECWorld.Reset()
	InventoryRes.Reset()
	ecs.AddResource(&ECWorld, GameDataRes)
	ecs.AddResource(&ECWorld, InventoryRes)
	ecs.AddResource(&ECWorld, CraftingTableRes)
	ecs.AddResource(&ECWorld, AnimPlayerDataRes)
	ecs.AddResource(&ECWorld, CameraRes)
	ecs.AddResource(&ECWorld, TileMapRes)

	GameTileMapGenerator.Opts.HighestSurfaceLevel = 10
	GameTileMapGenerator.Opts.LowestSurfaceLevel = 30
	GameTileMapGenerator.SetSeed(rand.Int())
	GameTileMapGenerator.NoiseState.FractalType(fastnoise.FractalFBm)
	GameTileMapGenerator.NoiseState.Frequency = 0.01
	GameTileMapGenerator.Generate()

	// ctrl.Collider.StaticCheck = true
	x, y := TileMapRes.FindSpawnPosition()
	// tileMap.Set(x, y+2, items.CraftingTable)
	SpawnX, SpawnY := TileMapRes.TileToWorldCenter(x, y)
	CameraRes.SmoothType = kamera.None
	CameraRes.SetCenter(SpawnX, SpawnY)
	CameraRes.SmoothOptions.LerpSpeedX = 0.5
	CameraRes.SmoothOptions.LerpSpeedY = 0.05
	CameraRes.SetTopLeft(TileMapRes.FloorToBlockCenter(CameraRes.X, CameraRes.Y))
	CurrentPlayer = SpawnPlayer(SpawnX, SpawnY)
	CameraRes.SmoothType = kamera.Lerp
}

func (m *MainMenu) saveGame() {
	FetchAnimPlayerData(PlayerAnimPlayer, AnimPlayerDataRes)
	jsonData, err := archeserde.Serialize(&ECWorld, m.serdeOpt)
	if err != nil {
		log.Fatal("Error serializing world:", err)
	}
	m.dataManager.SaveItem("01save", jsonData)
}

func (m *MainMenu) loadGame() {
	if m.dataManager.ItemExists("01save") {
		if !ECWorld.Alive(CurrentPlayer) {
			CurrentPlayer = SpawnPlayer(0, 0)
		}
		ECWorld.Reset()
		ecs.AddResource(&ECWorld, GameDataRes)
		ecs.AddResource(&ECWorld, InventoryRes)
		ecs.AddResource(&ECWorld, CraftingTableRes)
		ecs.AddResource(&ECWorld, AnimPlayerDataRes)
		ecs.AddResource(&ECWorld, CameraRes)
		ecs.AddResource(&ECWorld, TileMapRes)
		jsonData, err := m.dataManager.LoadItem("01save")
		if err != nil {
			log.Fatal("Error loading saved data:", err)
		}

		err = archeserde.Deserialize(jsonData, &ECWorld)
		if err != nil {
			log.Fatal("Error deserializing world:", err)
		}
		animData := ecs.GetResource[AnimPlayerData](&ECWorld)
		SetAnimPlayerData(PlayerAnimPlayer, animData)
	}
}
