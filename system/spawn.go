package system

import (
	"kar/arc"
	"kar/engine/v"

	"github.com/mlange-42/arche/ecs"
)

type vec2 = v.Vec

var Mario ecs.Entity

func (s *Spawn) Init() {
	Mario = arc.SpawnMario(10, 10)
}
func (s *Spawn) Update() {}
func (s *Spawn) Draw()   {}

type Spawn struct{}
