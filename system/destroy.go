package system

import (
	"image/color"
	"kar/arche"
	"kar/comp"
	"kar/items"
	"kar/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type DestroySystem struct {
}

func NewDestroySystem() *DestroySystem {
	return &DestroySystem{}
}

func (s *DestroySystem) Init() {
	resources.ECSWorld.OnRemove(EntityOnRemoveCallback)
}

func (s *DestroySystem) Update() {
	comp.TagBlock.Each(resources.ECSWorld, DestroyDeadBlockCallback)
}

func (s *DestroySystem) Draw(screen *ebiten.Image) {}

func EntityOnRemoveCallback(world donburi.World, entity donburi.Entity) {
	entry := world.Entry(entity)
	if entry.HasComponent(comp.TagBlock) {
		body := comp.Body.Get(entry)
		pos := body.Position()
		itemData := comp.Item.Get(entry)
		healthData := comp.Health.Get(entry)
		if healthData.Health <= 0 {
			body.FirstShape().SetSensor(true)
			dropID := items.Property[itemData.ID].Drops
			arche.SpawnDropItem(pos, dropID, MainWorld.WorldPosToChunkCoord(pos))
		}
	}

}
func DestroyDeadBlockCallback(e *donburi.Entry) {
	if comp.Health.Get(e).Health <= 0 {
		body := comp.Body.Get(e)
		pos := body.Position()
		body.FirstShape().SetSensor(true)
		destroyEntry(e)
		blockPos := MainWorld.WorldSpaceToPixelSpace(pos)
		MainWorld.Image.SetGray16(blockPos.X, blockPos.Y, color.Gray16{items.Air})

	}
}
