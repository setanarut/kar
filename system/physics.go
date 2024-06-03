package system

import (
	"kar/comp"
	"kar/engine/cm"
	"kar/res"
	"kar/types"
	"math"

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
	// res.Space.UseSpatialHash(50, 1000)
	res.Space.CollisionBias = math.Pow(0.5, 60)
	res.Space.CollisionSlop = 0.3
	// res.Space.Iterations = 20
	// res.Space.SetGravity(cm.Vec2{0, -1200})
	res.Space.NewCollisionHandler(types.CollEnemy, types.CollPlayer).BeginFunc = enemyPlayerBegin
	res.Space.NewCollisionHandler(types.CollEnemy, types.CollPlayer).PostSolveFunc = enemyPlayerPostSolve
	res.Space.NewCollisionHandler(types.CollPlayer, types.CollEnemy).SeparateFunc = playerEnemySep
	res.Space.NewCollisionHandler(types.CollSnowball, types.CollEnemy).BeginFunc = snowballEnemyBegin
	res.Space.Step(ps.DT)
}

func (ps *PhysicsSystem) Update() {

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
			enemyHealth := comp.Health.Get(enemy)
			*enemyHealth -= comp.Damage.GetValue(snowball)
		}
		DestroyEntryWithBody(snowball)
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
