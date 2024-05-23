package system

import (
	"kar/comp"
	"kar/res"

	"github.com/yohamta/donburi"
)

type AISystem struct {
}

func NewAISystem() *AISystem {
	return &AISystem{}
}

func (s *AISystem) Init() {}

func (s *AISystem) Update() {

	if pla, ok := comp.PlayerTag.First(res.World); ok {
		playerBody := comp.Body.Get(pla)
		comp.AI.Each(res.World, func(e *donburi.Entry) {
			if e.HasComponent(comp.Body) && e.HasComponent(comp.Mobile) {
				ai := *comp.AI.Get(e)
				if ai.Follow {
					body := comp.Body.Get(e)
					mobile := comp.Mobile.Get(e)
					dist := playerBody.Position().Distance(body.Position())
					if dist < ai.FollowDistance {
						speed := body.Mass() * (mobile.Speed * 4)
						a := playerBody.Position().Sub(body.Position()).Normalize().Mult(speed)
						body.ApplyForceAtLocalPoint(a, body.CenterOfGravity())
					}
				}
			}
		})
	}

}

func (s *AISystem) Draw() {}
