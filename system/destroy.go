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

	res.ECSWorld.OnRemove(func(world donburi.World, entity donburi.Entity) {
		entry := world.Entry(entity)
		if entry.HasComponent(comp.BlockTag) {
			body := comp.Body.Get(entry)
			pos := body.Position()
			itemData := comp.Item.Get(entry)
			healthData := comp.Health.Get(entry)
			if healthData.Health <= 0 {
				body.FirstShape().SetSensor(true)
				arche.SpawnDropItem(pos, itemData.Item, MainWorld.WorldPosToChunkCoord(pos))
			}
		}
	})
}

func (s *DestroySystem) Update() {
	comp.BlockTag.Each(res.ECSWorld, DestroyDeadBlock)
}

func (s *DestroySystem) Draw(screen *ebiten.Image) {}

func DestroyDeadBlock(e *donburi.Entry) {
	if comp.Health.Get(e).Health <= 0 {
		body := comp.Body.Get(e)
		pos := body.Position()
		body.FirstShape().SetSensor(true)
		DestroyEntryWithBody(e)
		blockPos := MainWorld.WorldSpaceToPixelSpace(pos)
		MainWorld.Image.SetGray(blockPos.X, blockPos.Y, color.Gray{uint8(items.Air)})

	}
}
