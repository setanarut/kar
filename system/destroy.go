package system

import (
	"image/color"
	"kar"
	"kar/arc"
	"kar/items"
	"kar/world"

	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/cm"
)

type Destroy struct {
	toRemove []ecs.Entity
}

func (s *Destroy) Init() {
	s.toRemove = make([]ecs.Entity, 0)
}

func (s *Destroy) Update() {

	// // Destroy counter for stucked drop item
	// arc.TagDropItem.Each(ecsWorld, func(dropEntry *donburi.Entry) {
	// 	dropShape := arc.Body.Get(dropEntry).Shapes[0]

	// 	Space.ShapeQuery(dropShape, func(shape *cm.Shape, points *cm.ContactPointSet) {
	// 		e := shape.Body.UserData.(*donburi.Entry)
	// 		if e.HasComponent(arc.TagBlock) {
	// 			if shape.BB.Contains(dropShape.BB.Offset(vec2{-3, -3})) {
	// 				ct := arc.StuckCountdown.Get(dropEntry)
	// 				ct.Duration -= 1
	// 			}
	// 		}
	// 	})
	// })

	// // Collision filter cooldown for item
	// comp.CollisionTimer.Each(ecsWorld, func(e *donburi.Entry) {
	// 	tm := comp.CollisionTimer.Get(e)
	// 	if timerIsReady(tm) {
	// 		shape := comp.Body.Get(e).Shapes[0]
	// 		shape.SetShapeFilter(dropItemFilterCooldown)
	// 	} else {
	// 	}
	// })

	DropItemQuery := arc.DropItemFilter.Query(&kar.WorldECS)
	for DropItemQuery.Next() {
		_, body, _, colltimer, stuckCountDown, _ := DropItemQuery.Get()
		dropShape := body.Shapes[0]

		kar.Space.ShapeQuery(dropShape, func(shape *cm.Shape, points *cm.ContactPointSet) {
			if shape.CollisionType == arc.Block {
				if shape.BB.Contains(dropShape.BB.Offset(vec2{-3, -3})) {
					stuckCountDown.Duration -= 1
				}
			}
		})

		if colltimer.Duration <= 0 {
			// destroy stucked item
			s.toRemove = append(s.toRemove, DropItemQuery.Entity())
		}
	}

	BlockQuery := arc.BlockFilter.Query(&kar.WorldECS)

	for BlockQuery.Next() {
		healthComponent, _, body, item := BlockQuery.Get()
		if healthComponent.Health <= 0 {
			pos := body.Position()
			body.Shapes[0].SetSensor(true)
			kar.Space.AddPostStepCallback(removeBodyPostStep, body, nil)
			s.toRemove = append(s.toRemove, BlockQuery.Entity())
			blockPos := world.WorldToPixel(pos)
			gameWorld.Image.SetGray16(blockPos.X, blockPos.Y, color.Gray16{items.Air})
			dropID := items.Property[item.ID].Drops
			arc.SpawnDropItem(pos, dropID)
		}
	}

	// The world is unlocked again.
	// Actually remove the collected entities.
	for _, e := range s.toRemove {
		body := kar.WorldECS.GetUnchecked(e, arc.BodyID)
		kar.Space.AddPostStepCallback(removeBodyPostStep, body, nil)
		kar.WorldECS.RemoveEntity(e)
	}
}

func (s *Destroy) Draw() {}
