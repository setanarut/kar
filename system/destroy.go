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
				switch itemData.ID {
				case items.CoalOre:
					arche.SpawnDropItem(pos, items.Coal, MainWorld.WorldPosToChunkCoord(pos))
				case items.GoldOre:
					arche.SpawnDropItem(pos, items.RawGold, MainWorld.WorldPosToChunkCoord(pos))
				case items.IronOre:
					arche.SpawnDropItem(pos, items.RawIron, MainWorld.WorldPosToChunkCoord(pos))
				case items.DiamondOre:
					arche.SpawnDropItem(pos, items.Diamond, MainWorld.WorldPosToChunkCoord(pos))
				case items.DeepSlateCoalOre:
					arche.SpawnDropItem(pos, items.Coal, MainWorld.WorldPosToChunkCoord(pos))
				case items.DeepSlateGoldOre:
					arche.SpawnDropItem(pos, items.RawGold, MainWorld.WorldPosToChunkCoord(pos))
				case items.DeepSlateIronOre:
					arche.SpawnDropItem(pos, items.RawIron, MainWorld.WorldPosToChunkCoord(pos))
				case items.DeepSlateDiamondOre:
					arche.SpawnDropItem(pos, items.Diamond, MainWorld.WorldPosToChunkCoord(pos))
				default:
					arche.SpawnDropItem(pos, itemData.ID, MainWorld.WorldPosToChunkCoord(pos))
				}

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
		destroyEntry(e)
		blockPos := MainWorld.WorldSpaceToPixelSpace(pos)
		MainWorld.Image.SetGray(blockPos.X, blockPos.Y, color.Gray{uint8(items.Air)})

	}
}
