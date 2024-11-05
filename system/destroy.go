package system

import (
	"kar"
	"kar/arc"

	"github.com/mlange-42/arche/ecs"
)

type Destroy struct {
	toRemove []ecs.Entity
}

func (s *Destroy) Init() {
	s.toRemove = make([]ecs.Entity, 0)
}

func (s *Destroy) Update() {

	CountdownQuery := arc.CountdownFilter.Query(&kar.WorldECS)
	for CountdownQuery.Next() {
		ct := CountdownQuery.Get()
		if ct.Duration <= 0 {
			s.toRemove = append(s.toRemove, CountdownQuery.Entity())
		}
	}

	// The world is unlocked again.
	// Actually remove the collected entities.
	for _, e := range s.toRemove {
		body := kar.WorldECS.GetUnchecked(e, arc.BodyID)
		kar.Space.AddPostStepCallback(removeBodyPostStep, body, false)
		kar.WorldECS.RemoveEntity(e)
	}

	// comp.TagBlock.Each(ecsWorld, DestroyDeadBlockCallback)
}

func (s *Destroy) Draw() {}

// func DestroyDeadBlockCallback(e *donburi.Entry) {
// 	if comp.Health.Get(e).Health <= 0 {
// 		body := comp.Body.Get(e)
// 		it := comp.Item.Get(e)
// 		pos := body.Position()
// 		body.Shapes[0].SetSensor(true)
// 		destroyEntry(e)
// 		blockPos := world.WorldToPixel(pos)
// 		gameWorld.Image.SetGray16(blockPos.X, blockPos.Y, color.Gray16{items.Air})
// 		dropID := items.Property[it.ID].Drops
// 		arc.SpawnDropItem(Space, ecsWorld, pos, dropID)
// 	}
// }

// func DestroyStuckedDropItem(e *donburi.Entry) {
// 	if comp.StuckCountdown.Get(e).Duration <= 0 {
// 		fmt.Println("Stucked item destoyed: ", getDisplayName(e))
// 		destroyEntry(e)
// 	}
// }
