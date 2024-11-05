package system

import (
	"kar"
	"kar/arc"

	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/cm"
)

type Destroy struct {
	toRemove []ecs.Entity
}

func (s *Destroy) Init() {
}

func (s *Destroy) Update() {

	dropItemQuery := arc.DropItemFilter.Query(&kar.WorldECS)

	for dropItemQuery.Next() {

		_, body, _, collisionTimer, stuckCountDown, index := dropItemQuery.Get()
		dropShape := body.Shapes[0]

		UpdateItemAnimationIndex(index)
		TimerUpdate((*arc.Timer)(collisionTimer))
		if TimerIsReady((*arc.Timer)(collisionTimer)) {
			body.Shapes[0].SetShapeFilter(dropItemFilterCooldown)
		}

		kar.Space.ShapeQuery(dropShape, func(shape *cm.Shape, points *cm.ContactPointSet) {
			if shape.CollisionType == arc.Block {
				if shape.BB.Contains(dropShape.BB.Offset(vec2{-3, -3})) {
					stuckCountDown.Duration -= 1
				}
			}
		})

		if collisionTimer.Duration <= 0 {
			// destroy stucked item
			s.toRemove = append(s.toRemove, dropItemQuery.Entity())
		}

	}

	// BlockQuery := arc.BlockFilter.Query(&kar.WorldECS)
	// for BlockQuery.Next() {
	// 	healthComponent, _, body, _ := BlockQuery.Get()
	// 	if healthComponent.Health <= 0 {
	// 		pos := body.Position()
	// 		body.Shapes[0].SetSensor(true)
	// 		kar.Space.AddPostStepCallback(removeBodyPostStep, body, nil)
	// 		s.toRemove = append(s.toRemove, BlockQuery.Entity())
	// 		blockPos := world.WorldToPixel(pos)
	// 		gameWorld.Image.SetGray16(blockPos.X, blockPos.Y, color.Gray16{items.Air})

	// 		// dropID := items.Property[item.ID].Drops
	// 		// arc.SpawnDropItem(pos, dropID)
	// 	}
	// }

	// The world is unlocked again.
	// Actually remove the collected entities.
	for _, e := range s.toRemove {
		body := arc.BodyMapper.Get(e)
		kar.Space.AddPostStepCallback(removeBodyPostStep, body, nil)
		kar.WorldECS.RemoveEntity(e)
	}

	// Empty the slice, so we can reuse it in the next time step.
	s.toRemove = s.toRemove[:0]
}

func (s *Destroy) Draw() {}

func UpdateItemAnimationIndex(i *arc.Index) {
	if i.Index < itemAnimFrameCount {
		i.Index++
	} else {
		i.Index = 0
	}
}
