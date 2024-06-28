package system

import (
	"fmt"
	"kar/comp"
	"kar/engine/cm"
	"kar/engine/vec"
	"kar/res"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type PhysicsSystem struct {
}

func NewPhysicsSystem() *PhysicsSystem {
	return &PhysicsSystem{}
}

func (ps *PhysicsSystem) Init() {
	res.Space.SetGravity(vec.Vec2{0, (res.BlockSize * 20)})
	res.Space.CollisionBias = math.Pow(0.3, 60)
	res.Space.CollisionSlop = 0.4
	// res.Space.UseSpatialHash(200, 1000)
	// res.Space.Iterations = 20
	// res.Space.Damping = 0.9
	// res.Space.NewCollisionHandler(types.CollSnowball, types.CollWall).BeginFunc = snowballBlockBegin
	res.Space.NewCollisionHandler(res.CollPlayer, res.CollCollectible).BeginFunc = PlayerCollectibleBegin
	// res.Space.NewCollisionHandler(types.CollEnemy, types.CollPlayer).BeginFunc = enemyPlayerBegin
	// res.Space.NewCollisionHandler(types.CollEnemy, types.CollPlayer).PostSolveFunc = enemyPlayerPostSolve
	// res.Space.NewCollisionHandler(types.CollPlayer, types.CollEnemy).SeparateFunc = playerEnemySep
	res.Space.Step(res.DeltaTime)
}

func (ps *PhysicsSystem) Update() {
	res.Space.Step(res.DeltaTime)

}

func (ps *PhysicsSystem) Draw(screen *ebiten.Image) {}

// Player <-> Collectible begin
func PlayerCollectibleBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	_, collectible := arb.Bodies()
	if CheckEntries(arb) {
		p, i := GetEntries(arb)
		inv := comp.Inventory.Get(p)
		itemData := comp.Item.Get(i)
		inv.Items[itemData.Item]++
		fmt.Println(inv.Items)
		DestroyBodyWithEntry(collectible)
	}
	return true
}

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
