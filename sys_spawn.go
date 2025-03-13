package kar

import (
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
}

func (s *Spawn) Init() {}
func (s *Spawn) Update() {

	gameDataRes.Duration += Tick

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
