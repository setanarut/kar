package system

import (
	"kar/comp"
	"kar/res"
	"math"

	"github.com/setanarut/cm"

	"github.com/setanarut/vec"

	"github.com/hajimehoshi/ebiten/v2"
)

type PhysicsSystem struct {
}

func NewPhysicsSystem() *PhysicsSystem {
	return &PhysicsSystem{}
}

func (ps *PhysicsSystem) Init() {
	res.Space.SetGravity(vec.Vec2{0, (res.BlockSize * 20)})
	res.Space.CollisionBias = math.Pow(0.1, 60)
	res.Space.CollisionSlop = 0.4
	res.Space.UseSpatialHash(200, 800)
	res.Space.Iterations = 20
	res.Space.Damping = 0.9
	res.Space.NewCollisionHandler(res.CollPlayer, res.CollDropItem).BeginFunc = PlayerDropItemBegin
}

func (ps *PhysicsSystem) Update() {
	res.Space.Step(res.DeltaTime)
}

func (ps *PhysicsSystem) Draw(screen *ebiten.Image) {}

// Player <-> DropItem begin
func PlayerDropItemBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	if checkEntries(arb) {
		player, DropItemEntry := getEntries(arb)
		inv := comp.Inventory.Get(player)
		itemData := comp.Item.Get(DropItemEntry)

		ok := inventoryManager.addItem(inv, itemData.Item)
		if ok {
			destroyEntry(DropItemEntry)
		}
	}
	return false
}
