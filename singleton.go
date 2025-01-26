package kar

import (
	"image"
	"image/color"
	"kar/items"
	"kar/res"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/anim"
	"github.com/setanarut/kamera/v2"
)

type ISystem interface {
	Init()
	Update()
	Draw()
}

var (
	DesktopPath              string
	WindowScale              float64     = 2.0
	ScreenW, ScreenH         float64     = 500.0, 340.0
	ItemCollisionDelay       int         = 10
	RaycastDist              int         = 4 // block unit
	DrawDebugHitboxesEnabled bool        = false
	DrawDebugTextEnabled     bool        = false
	PlayerBestToolDamage                 = 5.0
	PlayerDefaultDamage                  = 1.0
	RenderArea               image.Point = image.Point{
		(int(ScreenW) / 20) + 3,
		(int(ScreenH) / 20) + 3,
	}
	BackgroundColor color.RGBA = color.RGBA{36, 36, 39, 255}

	Screen           *ebiten.Image
	WorldECS         = ecs.NewWorld()
	Camera           *kamera.Camera
	GopherAnimPlayer *anim.AnimationPlayer
	GopherInventory  *items.Inventory
	ColorMDIO        = &colorm.DrawImageOptions{}
	ColorM           = colorm.ColorM{}
	// Debug
)

func init() {
	// GlobalColorM.ChangeHSV(1, 0, 1)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	DesktopPath = homePath + "/Desktop/"

	GopherAnimPlayer = anim.NewAnimationPlayer(
		&anim.Atlas{"Default", res.Player},
		&anim.Atlas{"WoodenAxe", res.PlayerWoodenAxeAtlas},
		&anim.Atlas{"StoneAxe", res.PlayerStoneAxeAtlas},
		&anim.Atlas{"IronAxe", res.PlayerIronAxeAtlas},
		&anim.Atlas{"DiamondAxe", res.PlayerDiamondAxeAtlas},
		&anim.Atlas{"WoodenPickaxe", res.PlayerWoodenPickaxeAtlas},
		&anim.Atlas{"StonePickaxe", res.PlayerStonePickaxeAtlas},
		&anim.Atlas{"IronPickaxe", res.PlayerIronPickaxeAtlas},
		&anim.Atlas{"DiamondPickaxe", res.PlayerDiamondPickaxeAtlas},
		&anim.Atlas{"WoodenShovel", res.PlayerWoodenShovelAtlas},
		&anim.Atlas{"StoneShovel", res.PlayerStoneShovelAtlas},
		&anim.Atlas{"IronShovel", res.PlayerIronShovelAtlas},
		&anim.Atlas{"DiamondShovel", res.PlayerDiamondShovelAtlas},
	)

	GopherAnimPlayer.NewState("idleRight", 0, 0, 16, 16, 1, false, false, 1)
	GopherAnimPlayer.NewState("idleUp", 208, 0, 16, 16, 1, false, false, 1)
	GopherAnimPlayer.NewState("idleDown", 224, 0, 16, 16, 1, false, false, 1)
	GopherAnimPlayer.NewState("walkRight", 16, 0, 16, 16, 4, false, false, 15)
	GopherAnimPlayer.NewState("jump", 16*5, 0, 16, 16, 1, false, false, 15)
	GopherAnimPlayer.NewState("skidding", 16*6, 0, 16, 16, 1, false, false, 15)
	GopherAnimPlayer.NewState("attackDown", 16*7, 0, 16, 16, 2, false, false, 8)
	GopherAnimPlayer.NewState("attackRight", 144, 0, 16, 16, 2, false, false, 8)
	GopherAnimPlayer.NewState("attackWalk", 0, 16, 16, 16, 4, false, false, 8)
	GopherAnimPlayer.NewState("attackUp", 16*11, 0, 16, 16, 2, false, false, 8)
	GopherAnimPlayer.CurrentState = "idleRight"

	GopherInventory = items.NewInventory()
}
