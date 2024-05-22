package system

import (
	"fmt"
	"kar/comp"
	"kar/constants"
	"kar/engine/cm"
	"kar/res"
	"math"

	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

type CollisionSystem struct {
	DT float64
}

func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{
		DT: 1.0 / 60.0,
	}
}

func (ps *CollisionSystem) Init() {
	res.Space.UseSpatialHash(50, 1000)
	res.Space.CollisionBias = math.Pow(0.3, 60)
	res.Space.CollisionSlop = 0.5
	res.Space.Damping = 0.03
	res.Space.NewCollisionHandler(constants.CollEnemy, constants.CollPlayer).PostSolveFunc = enemyPlayerPostSolve
	res.Space.NewCollisionHandler(constants.CollPlayer, constants.CollEnemy).SeparateFunc = playerEnemySep
	res.Space.Step(ps.DT)
}

func (ps *CollisionSystem) Update() {

	// Snowball arbiters
	comp.SnowballTag.Each(res.World, func(e *donburi.Entry) {
		b := comp.Body.Get(e)
		snowBallEntry := b.UserData.(*donburi.Entry)
		snowBallDamage := comp.Damage.GetValue(snowBallEntry)

		b.EachArbiter(func(a *cm.Arbiter) {
			if a.IsFirstContact() {
				_, other := a.Bodies()
				otherEntry := other.UserData.(*donburi.Entry)

				if otherEntry.Valid() && otherEntry.HasComponent(comp.Health) {
					h := comp.Health.Get(otherEntry)
					*h -= snowBallDamage

				}
				// fmt.Println(otherEntry.Archetype().Layout())
				DestroyBodyWithEntry(b)
			}

		})

	})

	if pla, ok := comp.PlayerTag.First(res.World); ok {
		playerBody := comp.Body.Get(pla)

		comp.EnemyTag.Each(res.World, func(e *donburi.Entry) {

			ene := comp.Body.Get(e)
			ai := *comp.AI.Get(e)
			mobile := comp.Mobile.Get(e)

			if ai.Follow {
				dist := playerBody.Position().Distance(ene.Position())
				if dist < ai.FollowDistance {
					speed := ene.Mass() * (mobile.Speed * 4)
					a := playerBody.Position().Sub(ene.Position()).Normalize().Mult(speed)
					ene.ApplyForceAtLocalPoint(a, ene.CenterOfGravity())
				}
			}

		})

	}

	res.Space.Step(ps.DT)

}

func (ps *CollisionSystem) Draw() {}

// Enemy <-> Player postsolve
func enemyPlayerPostSolve(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
	enemyBody, playerBody := arb.Bodies()
	enemyEntry, eok := enemyBody.UserData.(*donburi.Entry)
	fmt.Println(enemyEntry.HasComponent(comp.EnemyTag))
	playerEntry, pok := playerBody.UserData.(*donburi.Entry)

	if eok && pok {

		if playerEntry.Valid() && enemyEntry.Valid() {
			if playerEntry.HasComponent(comp.Health) && enemyEntry.HasComponent(comp.Damage) && playerEntry.HasComponent(comp.Render) {
				enemyDamage := comp.Damage.GetValue(enemyEntry)
				playerHealth := comp.Health.Get(playerEntry)
				comp.Render.Get(playerEntry).ScaleColor = colornames.Red
				if arb.IsFirstContact() {
					// playerBody.ApplyImpulseAtLocalPoint(arb.Normal().Mult(500), playerBody.CenterOfGravity())
					*playerHealth -= enemyDamage
				}
				*playerHealth -= enemyDamage / 60.0
				if *playerHealth < 0 {
					DestroyBodyWithEntry(playerBody)
				}
			}
		}

	}

}

// oyuncu düşmandan ayrılınca rengini beyaz yap
func playerEnemySep(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
	_, a := arb.Bodies()
	e := a.UserData.(*donburi.Entry)
	comp.Render.Get(e).ScaleColor = colornames.White
}
