package system

import (
	"kar/comp"
	"kar/resources"
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
	resources.Space.SetGravity(vec.Vec2{0, (resources.BlockSize * 20)})
	resources.Space.CollisionBias = math.Pow(0.1, 60)
	resources.Space.CollisionSlop = 0.4
	resources.Space.UseSpatialHash(200, 800)
	resources.Space.Iterations = 20
	resources.Space.Damping = 0.9
	resources.Space.NewCollisionHandler(resources.CollPlayer, resources.CollDropItem).BeginFunc = PlayerDropItemBegin
}

func (ps *PhysicsSystem) Update() {
	resources.Space.Step(resources.DeltaTime)
}

func (ps *PhysicsSystem) Draw(screen *ebiten.Image) {}

// Player <-> DropItem begin
func PlayerDropItemBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	if checkEntries(arb) {
		player, DropItemEntry := getEntries(arb)
		inv := comp.Inventory.Get(player)
		itemData := comp.Item.Get(DropItemEntry)

		ok := inventoryManager.addItem(inv, itemData.ID)
		if ok {
			destroyEntry(DropItemEntry)
		}
	}
	return false
}
