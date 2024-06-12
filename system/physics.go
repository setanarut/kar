package system

import (
	"kar/engine/cm"
	"kar/res"
	"math"
)

type PhysicsSystem struct {
}

func NewPhysicsSystem() *PhysicsSystem {
	return &PhysicsSystem{}
}

func (ps *PhysicsSystem) Init() {
	res.Space.SetGravity(cm.Vec2{0, -1500})
	res.Space.CollisionBias = math.Pow(0.5, 120)
	res.Space.CollisionSlop = 0.5
	// res.Space.UseSpatialHash(200, 1000)
	// res.Space.Iterations = 10
	// res.Space.Damping = 0.9
	// res.Space.NewCollisionHandler(types.CollSnowball, types.CollWall).BeginFunc = snowballBlockBegin
	// res.Space.NewCollisionHandler(types.CollSnowball, types.CollEnemy).BeginFunc = snowballEnemyBegin
	// res.Space.NewCollisionHandler(types.CollEnemy, types.CollPlayer).BeginFunc = enemyPlayerBegin
	// res.Space.NewCollisionHandler(types.CollEnemy, types.CollPlayer).PostSolveFunc = enemyPlayerPostSolve
	// res.Space.NewCollisionHandler(types.CollPlayer, types.CollEnemy).SeparateFunc = playerEnemySep
	res.Space.Step(res.DeltaTime)
}

func (ps *PhysicsSystem) Update() {
	res.Space.Step(res.DeltaTime)

}

func (ps *PhysicsSystem) Draw() {}

// // Enemy <-> Player begin
// func EnemyPlayerBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
// 	_, playerBody := arb.Bodies()
// 	if CheckEntries(arb) {
// 		a, b := GetEntries(arb)
// 		if b.HasComponent(comp.Health) && a.HasComponent(comp.Damage) && b.HasComponent(comp.DrawOptions) {
// 			enemyDamage := comp.Damage.GetValue(a)
// 			playerHealth := comp.Health.Get(b)
// 			comp.DrawOptions.Get(b).ScaleColor = colornames.Red
// 			playerBody.ApplyImpulseAtLocalPoint(arb.Normal().Scale(500), playerBody.CenterOfGravity())
// 			playerHealth.Health -= enemyDamage
// 		}
// 	}
// 	return true
// }

// func SnowballEnemyBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
// 	if CheckEntries(arb) {
// 		snowball, enemy := GetEntries(arb)
// 		if enemy.HasComponent(comp.Health) && snowball.HasComponent(comp.Damage) {
// 			enemyHealth := comp.Health.Get(enemy)
// 			enemyHealth.Health -= comp.Damage.GetValue(snowball)
// 		}
// 		DestroyEntryWithBody(snowball)
// 	}
// 	return true
// }

// // Enemy <-> Player postsolve
// func EnemyPlayerPostSolve(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
// 	if CheckEntries(arb) {
// 		_, playerBody := arb.Bodies()
// 		enemyEntry, playerEntry := GetEntries(arb)
// 		if playerEntry.HasComponent(comp.Health) && enemyEntry.HasComponent(comp.Damage) {
// 			fmt.Println(playerBody)
// 		}
// 	}
// }

// // oyuncu düşmandan ayrılınca rengini beyaz yap
// func PlayerEnemySep(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
// 	_, a := arb.Bodies()
// 	if CheckEntry(a) {
// 		e := GetEntry(a)
// 		if e.HasComponent(comp.DrawOptions) {
// 			if e.Valid() {
// 				do := comp.DrawOptions.Get(e)
// 				do.ScaleColor = colornames.White
// 			}
// 		}
// 	}
// }
