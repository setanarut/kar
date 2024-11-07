package system

import (
	"fmt"
	"kar"
	"kar/arc"

	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/cm"
)

var dropItemFilterCooldown = cm.ShapeFilter{
	Group:      2,
	Categories: arc.DropItemBit,
	Mask:       cm.AllCategories &^ arc.PlayerRayBit,
}

type DropItem struct {
	toRemove []ecs.Entity
}

func (s *DropItem) Init() {
}

func (s *DropItem) Update() {

	dropItemQuery := arc.FilterDropItem.Query(&kar.WorldECS)

	for dropItemQuery.Next() {

		_, bd, _, cac, stuckDestuctionCountdown, idx := dropItemQuery.Get()
		dropShape := bd.Body.Shapes[0]

		// Update drop item aimation frames
		if idx.Index < itemAnimFrameCount {
			idx.Index++
		} else {
			idx.Index = 0
		}

		// Collision Activation Countdown
		cac.Tick -= 1
		if cac.Tick <= 0 {
			bd.Body.Shapes[0].SetShapeFilter(dropItemFilterCooldown)
		}

		// eğer Item blok içinde sıkışmışsa yok etme sayacını ilerlet
		kar.Space.ShapeQuery(dropShape, func(shape *cm.Shape, points *cm.ContactPointSet) {
			if shape.CollisionType == arc.Block {
				if shape.BB.Contains(dropShape.BB.Offset(vec2{-3, -3})) {
					stuckDestuctionCountdown.Tick -= 1
				}
			}
		})

		if stuckDestuctionCountdown.Tick <= 0 {
			fmt.Println("Stucked item destroyed!")
			s.toRemove = append(s.toRemove, dropItemQuery.Entity())
		}

	}

	// The world is unlocked again.
	// Actually remove the collected entities.
	for _, e := range s.toRemove {
		bd := arc.MapBody.Get(e)
		kar.Space.AddPostStepCallback(removeBodyPostStep, bd.Body, nil)
		kar.WorldECS.RemoveEntity(e)
	}

	// Empty the slice, so we can reuse it in the next time step.
	s.toRemove = s.toRemove[:0]
}

func (s *DropItem) Draw() {}
