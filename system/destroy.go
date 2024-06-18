package system

import (
	"kar/comp"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
)

type DestroySystem struct {
}

func NewDestroySystem() *DestroySystem {
	return &DestroySystem{}
}

func (s *DestroySystem) Init() {
}

func (s *DestroySystem) Update() {
	comp.BlockItemTag.Each(res.World, DestroyZeroHealthSetMapBlockState)
	// comp.DropItemTag.Each(res.World, DestroyZeroHealthSetMapBlockState)
	// comp.SnowballTag.Each(res.World, DestroyOnCollisionAndStopped)
}

func (s *DestroySystem) Draw(screen *ebiten.Image) {}
