package system

import (
	"image/color"
	"kar/arche"
	"kar/comp"
	"kar/items"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
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

func DestroyZeroHealthSetMapBlockState(e *donburi.Entry) {
	if comp.Health.Get(e).Health <= 0 {
		body := comp.Body.Get(e)
		pos := body.Position()
		i := comp.Item.Get(e)
		body.FirstShape().SetSensor(true)
		DestroyEntryWithBody(e)
		if i.Item == items.IronOre {
			arche.SpawnDebug(pos)
		}
		blockPos := Terr.WorldSpaceToMapSpace(pos)
		res.Terrain.SetGray(blockPos.X, blockPos.Y, color.Gray{uint8(items.Air)})

	}
}
