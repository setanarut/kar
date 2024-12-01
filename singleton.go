package kar

import (
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/vec"
)

type vec2 = vec.Vec2

const TimerTick time.Duration = time.Second / 60
const DeltaTime float64 = 1.0 / 60.0

var (
	Screen      *ebiten.Image
	Camera      *kamera.Camera
	WorldECS    = ecs.NewWorld()
	DesktopPath string
	GlobalDIO   = &ebiten.DrawImageOptions{}
)

type ISystem interface {
	Init()
	Update()
	Draw()
}

func init() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	DesktopPath = homePath + "/Desktop/"
}
