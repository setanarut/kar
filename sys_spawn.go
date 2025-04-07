package kar

import (
	"time"

	"github.com/mlange-42/ark/ecs"
)

// spawnItemData is a helper for delaying spawn events
type spawnItemData struct {
	Pos        Vec
	Id         uint8
	Durability int
}
type spawnEffectData struct {
	Pos Vec
	Id  uint8
}

var (
	toSpawnItem   = []spawnItemData{}
	toSpawnEffect = []spawnEffectData{}
	toRemove      []ecs.Entity
)

type Spawn struct {
	spawnInterval time.Duration
}

func (s *Spawn) Init() {
	s.spawnInterval = time.Second * 4
}
func (s *Spawn) Update() {

	gameDataRes.Duration += Tick
	gameDataRes.SpawnElapsed += Tick

	if gameDataRes.SpawnElapsed > s.spawnInterval {
		gameDataRes.SpawnElapsed = 0
	}

	if gameDataRes.SpawnElapsed == 0 {
		if world.Alive(currentPlayer) {
			p := mapAABB.GetUnchecked(currentPlayer).Pos
			mapEnemy.NewEntity(
				&AABB{
					Pos:  p.Sub(Vec{50, 50}),
					Half: Vec{8, 4.5},
				},
				&Velocity{0.4, 0},
				ptr(CrabID),
				ptr(AnimationTick(0)),
			)
		}
	}

	// Spawn item
	for _, data := range toSpawnItem {
		SpawnItem(data.Pos, data.Id, data.Durability)
	}
	// Spawn effect
	for _, data := range toSpawnEffect {
		SpawnEffect(data.Pos, data.Id)
	}

	toSpawnItem = toSpawnItem[:0]
	toSpawnEffect = toSpawnEffect[:0]

	for _, e := range toRemove {
		world.RemoveEntity(e)
	}
	toRemove = toRemove[:0]
}
func (s *Spawn) Draw() {
}
