package system

import (
	"kar/comp"
	"kar/res"
)

type DestroySystem struct {
}

func NewDestroySystem() *DestroySystem {
	return &DestroySystem{}
}

func (s *DestroySystem) Init() {
	// res.World.OnRemove(func(world donburi.World, entity donburi.Entity) {
	// 	e := world.Entry(entity)
	// 	if e.HasComponent(comp.EnemyTag) {
	// 		p := comp.Body.Get(e).Position()
	// 		i := comp.Inventory.Get(e)
	// 		for _, v := range i.Keys {
	// 			arche.SpawnDefaultKeyCollectible(v, p)
	// 		}
	// 		for range i.Bombs {
	// 			arche.SpawnDefaultBomb(p)
	// 		}
	// 		for range i.Snowballs {
	// 			arche.SpawnDefaultSnowball(p)
	// 		}
	// 	}

	// })

}

func (s *DestroySystem) Update() {
	comp.Health.Each(res.World, destroyDead)
	comp.SnowballTag.Each(res.World, destroyOnCollisionAndStopped)
}

func (s *DestroySystem) Draw() {}
