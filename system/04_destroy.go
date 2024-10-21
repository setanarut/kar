package system

import (
	"fmt"
	"image/color"
	"kar/arche"
	"kar/comp"
	"kar/items"
	world "kar/world"

	"github.com/yohamta/donburi"
)

type Destroy struct{}

func (s *Destroy) Init() {
	// res.ECSWorld.OnRemove(EntityOnRemoveCallback)
}

func (s *Destroy) Update() {
	comp.TagBlock.Each(ecsWorld, DestroyDeadBlockCallback)
	comp.StuckCountdown.Each(ecsWorld, DestroyStuckedDropItem)
}

func (s *Destroy) Draw() {}

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
		mainWorld.Image.SetGray16(blockPos.X, blockPos.Y, color.Gray16{items.Air})
		dropID := items.Property[it.ID].Drops
		arche.SpawnDropItem(space, ecsWorld, pos, dropID)
	}
}

func DestroyStuckedDropItem(e *donburi.Entry) {
	if comp.StuckCountdown.Get(e).Duration <= 0 {
		fmt.Println("Stucked item destoyed: ", getDisplayName(e))
		destroyEntry(e)
	}
}
