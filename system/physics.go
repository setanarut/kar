package system

import (
	"kar/comp"
	"kar/engine/cm"
	"kar/res"
	"math"

	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

type PhysicsSystem struct {
	DT float64
}

func NewPhysicsSystem() *PhysicsSystem {
	return &PhysicsSystem{
		DT: 1.0 / 60.0,
	}
}

func (ps *PhysicsSystem) Init() {
	res.Space.UseSpatialHash(50, 1000)
	res.Space.CollisionBias = math.Pow(0.3, 60)
	res.Space.CollisionSlop = 0.5
	res.Space.Damping = 0.03
	res.Space.NewCollisionHandler(res.CollEnemy, res.CollPlayer).PostSolveFunc = enemyPlayerPostSolve
	res.Space.NewCollisionHandler(res.CollPlayer, res.CollEnemy).SeparateFunc = playerEnemySep
	res.Space.NewCollisionHandler(res.CollEnemy, res.CollPlayer).BeginFunc = enemyPlayerBegin
	res.Space.NewCollisionHandler(res.CollSnowball, res.CollEnemy).BeginFunc = snowballEnemyBegin
	res.Space.Step(ps.DT)
}

func (ps *PhysicsSystem) Update() {

	comp.WASDControll.Each(res.World, func(e *donburi.Entry) {
		body := comp.Body.Get(e)
		mobileData := comp.Mobile.Get(e)
		velocity := res.Input.WASDDirection.Normalize().Mult(mobileData.Speed)
		body.SetVelocityVector(body.Velocity().LerpDistance(velocity, mobileData.Accel))
	})

	// düşman takip AI
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

func (ps *PhysicsSystem) Draw() {}

// Enemy <-> Player begin
func enemyPlayerBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	_, playerBody := arb.Bodies()
	if checkEntries(arb) {
		a, b := getEntries(arb)
		if b.HasComponent(comp.Health) && a.HasComponent(comp.Damage) && b.HasComponent(comp.Render) {
			enemyDamage := comp.Damage.GetValue(a)
			playerHealth := comp.Health.Get(b)
			comp.Render.Get(b).ScaleColor = colornames.Red
			playerBody.ApplyImpulseAtLocalPoint(arb.Normal().Mult(500), playerBody.CenterOfGravity())
			*playerHealth -= enemyDamage
		}
	}
	return true
}
func snowballEnemyBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	if checkEntries(arb) {
		snowball, enemy := getEntries(arb)
		if enemy.HasComponent(comp.Health) && snowball.HasComponent(comp.Damage) && enemy.HasComponent(comp.Render) {
			enemyDamage := comp.Damage.GetValue(snowball)
			enemyHealth := comp.Health.Get(enemy)
			*enemyHealth -= enemyDamage
		}
		destroyEntryWithBody(snowball)
	}
	return true
}

// Enemy <-> Player postsolve
func enemyPlayerPostSolve(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
	if checkEntries(arb) {
		_, playerBody := arb.Bodies()
		enemyEntry, playerEntry := getEntries(arb)
		if playerEntry.HasComponent(comp.Health) && enemyEntry.HasComponent(comp.Damage) && playerEntry.HasComponent(comp.Render) {
			enemyDamage := comp.Damage.GetValue(enemyEntry)
			playerHealth := comp.Health.Get(playerEntry)
			comp.Render.Get(playerEntry).ScaleColor = colornames.Red
			playerBody.ApplyImpulseAtLocalPoint(arb.Normal().Mult(500), playerBody.CenterOfGravity())
			*playerHealth -= enemyDamage / 60.0
		}
	}
}

// oyuncu düşmandan ayrılınca rengini beyaz yap
func playerEnemySep(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
	_, a := arb.Bodies()
	if checkEntry(a) {
		e := getEntry(a)
		if e.HasComponent(comp.Render) {
			if e.Valid() {
				r := comp.Render.Get(e)
				r.ScaleColor = colornames.White
			}
		}
	}
}
