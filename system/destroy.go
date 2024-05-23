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

func (s *DestroySystem) Init() {}

func (s *DestroySystem) Update() {
	comp.Health.Each(res.World, destroyDead)
	comp.SnowballTag.Each(res.World, destroyStopped)
}

func (s *DestroySystem) Draw() {}
