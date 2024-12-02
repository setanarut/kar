package system

import (
	"kar/arc"
	"kar/engine/v"
	"kar/tilemap"

	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/tilecollider"
)

type vec2 = v.Vec

var Mario ecs.Entity
var Map *tilemap.TileMap
var TCollider *tilecollider.Collider[uint16]

func (s *Spawn) Init() {
	Mario = arc.SpawnMario(50, 0)
	Map = tilemap.NewTileMap(7, 7, 48, 48)

	Map.Grid = [][]uint16{
		{1, 0, 1, 0, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 0, 0, 0, 1, 0, 1},
		{0, 0, 0, 1, 0, 1, 0, 1},
		{0, 0, 0, 0, 0, 1, 0, 1},
		{1, 0, 1, 1, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1}}

	TCollider = tilecollider.NewCollider(Map.Grid, 48, 48)

}
func (s *Spawn) Update() {}
func (s *Spawn) Draw()   {}

type Spawn struct{}
