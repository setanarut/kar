package system

import (
	"image/color"
	"kar/arche"
	"kar/comp"
	"kar/items"
	"kar/res"
	world "kar/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type DestroySystem struct{}

func (s *DestroySystem) Init() {
	// res.ECSWorld.OnRemove(EntityOnRemoveCallback)
}

func (s *DestroySystem) Update() {
	comp.TagBlock.Each(res.ECSWorld, DestroyDeadBlockCallback)
}

func (s *DestroySystem) Draw(screen *ebiten.Image) {}

func EntityOnRemoveCallback(world donburi.World, entity donburi.Entity) {
	entry := world.Entry(entity)
	if entry.HasComponent(comp.TagBlock) && entry.HasComponent(comp.Health) {
	}

}
func DestroyDeadBlockCallback(e *donburi.Entry) {
	if comp.Health.Get(e).Health <= 0 {
		body := comp.Body.Get(e)
		it := comp.Item.Get(e)
		pos := body.Position()
		body.Shapes[0].SetSensor(true)
		destroyEntry(e)
		blockPos := world.WorldSpaceToPixelSpace(pos)
		MainWorld.Image.SetGray16(blockPos.X, blockPos.Y, color.Gray16{items.Air})
		dropID := items.Property[it.ID].Drops
		arche.SpawnDropItem(pos, dropID)
	}
}
