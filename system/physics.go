package system

import (
	"kar/comp"
	"kar/res"
	"math"

	"github.com/setanarut/cm"
	"github.com/yohamta/donburi"

	"github.com/setanarut/vec"

	"github.com/hajimehoshi/ebiten/v2"
)

var dropItemFilterWithPlayer = cm.ShapeFilter{
	Group:      2,
	Categories: res.DropItemMask,
	Mask:       cm.AllCategories &^ res.PlayerRayMask,
}

type Physics struct{}

func (ps *Physics) Init() {
	res.Space.SetGravity(vec.Vec2{0, (res.BlockSize * 20)})
	res.Space.CollisionBias = math.Pow(0.2, 60)
	res.Space.CollisionSlop = 0.4
	res.Space.UseSpatialHash(128, 800)
	res.Space.Iterations = 20
	res.Space.Damping = 0.9

	playerDropItemHandler := res.Space.NewCollisionHandler(
		res.CollPlayer,
		res.CollDropItem)
	// DropItemBlockHandler := res.Space.NewCollisionHandler(
	// 	res.CollDropItem,
	// 	res.CollBlock)

	playerDropItemHandler.BeginFunc = playerDropItemBeginCallback
	// playerDropItemHandler.PreSolveFunc = playerDropItemPreCallback
	// playerDropItemHandler.PostSolveFunc = playerDropItemPostCallback
	// DropItemBlockHandler.PreSolveFunc = dropItemBlockPreCallback
}

func (ps *Physics) Update() {
	res.Space.Step(res.DeltaTime)

	comp.SpawnTimer.Each(res.ECSWorld, func(e *donburi.Entry) {
		tm := comp.SpawnTimer.Get(e)
		if TimerIsReady(tm) {
			f := comp.Body.Get(e).Shapes[0]
			f.SetShapeFilter(dropItemFilterWithPlayer)
		}
	})
}

func (ps *Physics) Draw(screen *ebiten.Image) {}

// Player <-> Drop Item Begin
func playerDropItemBeginCallback(arb *cm.Arbiter, s *cm.Space, dat any) bool {
	if checkEntries(arb) {
		a, b := getEntries(arb)
		inv := comp.Inventory.Get(a)
		itemData := comp.Item.Get(b)

		ok := inventoryManager.addItemIfEmpty(inv, itemData.ID)
		if ok {
			destroyEntry(b)
		}
	}
	return false
}

// Player <-> Drop Item Pre
func playerDropItemPreCallback(arb *cm.Arbiter, _ *cm.Space, _ any) bool {
	if checkEntries(arb) {
		a, b := getEntries(arb)
		inv := comp.Inventory.Get(a)
		itemData := comp.Item.Get(b)

		ok := inventoryManager.addItemIfEmpty(inv, itemData.ID)
		if ok {
			destroyEntry(b)
		}
	}
	return false
}

// Player <-> Drop Item Pots
func playerDropItemPostCallback(arb *cm.Arbiter, _ *cm.Space, dat any) {
	if checkEntries(arb) {
		a, b := getEntries(arb)
		inv := comp.Inventory.Get(a)
		itemData := comp.Item.Get(b)

		ok := inventoryManager.addItemIfEmpty(inv, itemData.ID)
		if ok {
			destroyEntry(b)
		}
	}
}

// Drop Item <-> Block Pre
func dropItemBlockPreCallback(arb *cm.Arbiter, s *cm.Space, dat any) bool {
	if arb.IsFirstContact() {
		dist := arb.ContactPointSet().Points[0].Distance
		if dist < -32 {
			return false
			// b, _ := arb.Bodies()
			// b.SetVelocity(0, 0)
			// b.SetPosition(b.Position().Add(vec.Vec2{0, dist}))
		}
	}
	return true
}
