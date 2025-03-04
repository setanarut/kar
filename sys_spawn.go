package kar

import (
	"github.com/mlange-42/ark/ecs"
)

// spawnData is a helper for delaying spawn events
type spawnData struct {
	Pos        Vec
	Id         uint8
	Durability int
}

var (
	toSpawn  = []spawnData{}
	toRemove []ecs.Entity
)

type Spawn struct {
}

func (s *Spawn) Init() {}
func (s *Spawn) Update() {

	// Spawn item
	for _, data := range toSpawn {
		SpawnItem(data.Pos, data.Id, data.Durability)
	}

	toSpawn = toSpawn[:0]

	for _, e := range toRemove {
		world.RemoveEntity(e)
	}
	toRemove = toRemove[:0]
}
func (s *Spawn) Draw() {
}
