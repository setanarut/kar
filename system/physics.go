package system

import (
	"fmt"
	"kar"
	"kar/arc"

	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

var dropItemFilterCooldown = cm.ShapeFilter{
	Group:      2,
	Categories: arc.DropItemBit,
	Mask:       cm.AllCategories &^ arc.PlayerRayBit,
}

type Physics struct{}

func (ps *Physics) Init() {
	kar.Space.SetGravity(vec.Vec2{0, kar.Gravity})
	kar.Space.CollisionBias = kar.CollisionBias
	kar.Space.CollisionSlop = kar.CollisionSlop
	kar.Space.Damping = kar.Damping
	kar.Space.Iterations = kar.Iterations

	if kar.UseSpatialHash {
		kar.Space.UseSpatialHash(kar.SpatialHashDim, kar.SpatialHashCount)
	}

	PlayerBlockHandler := kar.Space.NewCollisionHandler(arc.Player, arc.Block)
	PlayerBlockHandler.BeginFunc = PlayerBlockBegin
	if false {
		PlayerDropItemHandler := kar.Space.NewCollisionHandler(
			arc.Player,
			arc.DropItem)

		PlayerDropItemHandler.BeginFunc = PlayerDropItemBegin

		// PlayerDropItemHandler.PreSolveFunc = PlayerDropItemPreCallback
		// PlayerDropItemHandler.PostSolveFunc = playerDropItemPostCallback

	}
	// if false {
	// 	DropItemBlockHandler := kar.Space.NewCollisionHandler(
	// 		arc.DropItem,
	// 		arc.Block)
	// 	DropItemBlockHandler.BeginFunc = DropItemBlockBegin
	// 	DropItemBlockHandler.PreSolveFunc = DropItemBlockPre
	// 	DropItemBlockHandler.PostSolveFunc = DropItemBlockPost
	// }
}

func (ps *Physics) Update() {
	kar.Space.Step(kar.DeltaTime)
}

func (ps *Physics) Draw() {}

func PlayerDropItemBegin(arb *cm.Arbiter, s *cm.Space, dat any) bool {
	a, b := arb.Bodies()
	player := a.UserData.(ecs.Entity)
	dropItem := b.UserData.(ecs.Entity)
	inv := arc.InventoryMapper.Get(player)
	itm := arc.ItemMapper.Get(dropItem)
	ok := addItem(inv, itm.ID)
	if ok {
		kar.Space.AddPostStepCallback(removeBodyPostStep, b, nil)
		kar.WorldECS.RemoveEntity(dropItem)
	}
	return false
}

// func PlayerDropItemPre(arb *cm.Arbiter, _ *cm.Space, _ any) bool {
// 	if checkEntries(arb) {
// 		a, b := getEntries(arb)
// 		inv := comp.Inventory.Get(a)
// 		itemData := comp.Item.Get(b)
// 		ok := addItem(inv, itemData.ID)
// 		if ok {
// 			destroyEntry(b)
// 		}
// 	}
// 	return false
// }

// func PlayerDropItemPost(arb *cm.Arbiter, _ *cm.Space, dat any) {
// 	if checkEntries(arb) {
// 		a, b := getEntries(arb)
// 		inv := comp.Inventory.Get(a)
// 		itemData := comp.Item.Get(b)

// 		ok := addItem(inv, itemData.ID)
// 		if ok {
// 			destroyEntry(b)
// 		}
// 	}
// }

//	func DropItemBlockBegin(arb *cm.Arbiter, _ *cm.Space, _ any) bool {
//		dist := arb.ContactPointSet().Points[0].Distance
//		if dist < -32 {
//			return false
//		} else {
//			return true
//		}
//	}
func PlayerBlockBegin(arb *cm.Arbiter, _ *cm.Space, _ any) bool {
	fmt.Println(arb.Bodies())
	return true
}

// func DropItemBlockPre(arb *cm.Arbiter, _ *cm.Space, _ any) bool {
// 	dist := arb.ContactPointSet().Points[0].Distance
// 	if dist < -10 {
// 		fmt.Println("Pre -10")
// 	}
// 	return true
// }

// func DropItemBlockPost(arb *cm.Arbiter, _ *cm.Space, _ any) {
// 	dist := arb.ContactPointSet().Points[0].Distance
// 	if dist < -10 {
// 		fmt.Println("Post -10")
// 	}
// }
