package res

import (
	"embed"
	_ "image/png"
	"kar/engine/util"
	"kar/items"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/setanarut/cm"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/vec"
	"github.com/yohamta/donburi"
)

type DIO = ebiten.DrawImageOptions

//go:embed assets/*
var fs embed.FS

var (
	Cam                *kamera.Camera
	DesktopDir         string
	ScreenSize         = vec.Vec2{960, 540}
	WorldSize          = vec.Vec2{256, 256}
	ChunkSize          = vec.Vec2{11, 8} // {9, 5} {12, 9}
	BlockCenterOffset  = vec.Vec2{(BlockSize / 2), (BlockSize / 2)}.Neg()
	BlockSize          = 64.0
	GlobalDIO          = &DIO{Filter: ebiten.FilterNearest}
	SelectedSlotItemID = items.Air
	SelectedSlotIndex  = 0
	ECSWorld           = donburi.NewWorld()
	Space              = cm.NewSpace()
	Font               = util.LoadFontFromFS("assets/font/pixelcode.otf", 18, fs)
	FontSmall          = &text.GoTextFace{
		Source:    Font.Source,
		Direction: 0,
		Size:      9,
	}
	FilterPlayerRaycast = cm.ShapeFilter{
		Group:      cm.NoGroup,
		Categories: PlayerRayMask,
		Mask:       cm.AllCategories &^ PlayerMask}
	FontDrawOptions = &text.DrawOptions{
		DrawImageOptions: DIO{Filter: ebiten.FilterNearest},
		LayoutOptions:    text.LayoutOptions{LineSpacing: Font.Size * 1.3},
	}
	FontSmallDrawOptions = &text.DrawOptions{
		DrawImageOptions: DIO{Filter: ebiten.FilterNearest},
		LayoutOptions:    text.LayoutOptions{LineSpacing: FontSmall.Size * 1.3},
	}
)

func init() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	DesktopDir = homePath + "/Desktop/"
}
