package system

import (
	"kar/comp"
	"kar/engine/cm"
	"kar/res"
	"kar/types"

	"github.com/yohamta/donburi"
)

type AISystem struct {
}

func NewAISystem() *AISystem {
	return &AISystem{}
}

func (s *AISystem) Init() {
}

func (s *AISystem) Update() {
	res.QueryAI.Each(res.World, followAI)
}

func (s *AISystem) Draw() {}

func followAI(e *donburi.Entry) {
	ai := comp.AI.Get(e)
	if ai.Follow {
		moveTo(comp.Body.Get(e), comp.Mobile.Get(e), ai)
	}
}

func moveTo(b *cm.Body, m *types.DataMobile, ai *types.DataAI) {
	if ai.Target != nil {
		if ai.Target.HasComponent(comp.Body) {
			targetPos := comp.Body.Get(ai.Target).Position()
			dist := targetPos.Distance(b.Position())
			if dist < ai.FollowDistance {
				force := targetPos.Sub(b.Position()).Normalize().Mult(b.Mass() * (m.Speed * 4))
				b.ApplyForceAtLocalPoint(force, b.CenterOfGravity())
			}
		}
	}
}
