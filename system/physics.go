package system

import (
	"image/color"
	"kar/comp"
	"kar/engine/cm"
	"kar/res"
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
	// res.Space.UseSpatialHash(200, 1000)
	res.Space.CollisionBias = math.Pow(0.5, 120)
	res.Space.CollisionSlop = 0.5
	// res.Space.Iterations = 10
	res.Space.SetGravity(cm.Vec2{0, -1500})
	// res.Space.Damping = 0.9
	// res.Space.NewCollisionHandler(types.CollSnowball, types.CollWall).BeginFunc = snowballBlockBegin
	// res.Space.NewCollisionHandler(types.CollSnowball, types.CollEnemy).BeginFunc = snowballEnemyBegin
	// res.Space.NewCollisionHandler(types.CollEnemy, types.CollPlayer).BeginFunc = enemyPlayerBegin
	// res.Space.NewCollisionHandler(types.CollEnemy, types.CollPlayer).PostSolveFunc = enemyPlayerPostSolve
	// res.Space.NewCollisionHandler(types.CollPlayer, types.CollEnemy).SeparateFunc = playerEnemySep
	res.Space.Step(ps.DT)
}

func (ps *PhysicsSystem) Update() {

	res.Space.Step(ps.DT)

}

func (ps *PhysicsSystem) Draw() {}

// Enemy <-> Player begin
func EnemyPlayerBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	_, playerBody := arb.Bodies()
	if CheckEntries(arb) {
		a, b := GetEntries(arb)
		if b.HasComponent(comp.Health) && a.HasComponent(comp.Damage) && b.HasComponent(comp.Render) {
			enemyDamage := comp.Damage.GetValue(a)
			playerHealth := comp.Health.Get(b)
			comp.Render.Get(b).ScaleColor = colornames.Red
			playerBody.ApplyImpulseAtLocalPoint(arb.Normal().Scale(500), playerBody.CenterOfGravity())
			*playerHealth -= enemyDamage
		}
	}
	return true
}

func SnowballBlockBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	if CheckEntries(arb) {
		snowball, block := GetEntries(arb)
		if block.HasComponent(comp.Health) && snowball.HasComponent(comp.Damage) {
			enemyHealth := comp.Health.Get(block)
			*enemyHealth -= comp.Damage.GetValue(snowball)
		}
		blockBody := comp.Body.Get(block)
		pos := blockBody.Position().Point().Div(50)

		res.Terrain.Set(pos.X, pos.Y, color.Gray{0})
		// DestroyEntryWithBody(snowball)
	}
	return true
}
func SnowballEnemyBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	if CheckEntries(arb) {
		snowball, enemy := GetEntries(arb)
		if enemy.HasComponent(comp.Health) && snowball.HasComponent(comp.Damage) && enemy.HasComponent(comp.Render) {
			enemyHealth := comp.Health.Get(enemy)
			*enemyHealth -= comp.Damage.GetValue(snowball)
		}
		DestroyEntryWithBody(snowball)
	}
	return true
}

// Enemy <-> Player postsolve
func EnemyPlayerPostSolve(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
	if CheckEntries(arb) {
		_, playerBody := arb.Bodies()
		enemyEntry, playerEntry := GetEntries(arb)
		if playerEntry.HasComponent(comp.Health) && enemyEntry.HasComponent(comp.Damage) && playerEntry.HasComponent(comp.Render) {
			enemyDamage := comp.Damage.GetValue(enemyEntry)
			playerHealth := comp.Health.Get(playerEntry)
			comp.Render.Get(playerEntry).ScaleColor = colornames.Red
			playerBody.ApplyImpulseAtLocalPoint(arb.Normal().Scale(500), playerBody.CenterOfGravity())
			*playerHealth -= enemyDamage / 60.0
		}
	}
}

// oyuncu düşmandan ayrılınca rengini beyaz yap
func PlayerEnemySep(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
	_, a := arb.Bodies()
	if CheckEntry(a) {
		e := GetEntry(a)
		if e.HasComponent(comp.Render) {
			if e.Valid() {
				r := comp.Render.Get(e)
				r.ScaleColor = colornames.White
			}
		}
	}
}
