package system

import (
	"fmt"
	"kar"
	"kar/arche"
	"kar/comp"

	"github.com/setanarut/cm"
	"github.com/yohamta/donburi"
)

var dropItemFilterCooldown = cm.ShapeFilter{
	Group:      2,
	Categories: arche.DropItemBit,
	Mask:       cm.AllCategories &^ arche.PlayerRayBit,
}

type Physics struct{}

func (ps *Physics) Init() {
	Space.SetGravity(vec2{0, (kar.BlockSize * 20)})
	Space.CollisionBias = kar.CollisionBias
	Space.CollisionSlop = kar.CollisionSlop
	Space.Damping = kar.Damping
	if kar.UseSpatialHash {
		Space.UseSpatialHash(128, 800)
		Space.Iterations = 10

	}

	if true {
		PlayerDropItemHandler := Space.NewCollisionHandler(
			arche.Player,
			arche.DropItem)

		PlayerDropItemHandler.BeginFunc = PlayerDropItemBegin
		// PlayerDropItemHandler.PreSolveFunc = PlayerDropItemPreCallback
		// PlayerDropItemHandler.PostSolveFunc = playerDropItemPostCallback

	}
	if false {
		DropItemBlockHandler := Space.NewCollisionHandler(
			arche.DropItem,
			arche.Block)
		DropItemBlockHandler.BeginFunc = DropItemBlockBegin
		DropItemBlockHandler.PreSolveFunc = DropItemBlockPre
		DropItemBlockHandler.PostSolveFunc = DropItemBlockPost
	}
}

func (ps *Physics) Update() {
	Space.Step(kar.DeltaTime)

	// Destroy counter for stucked drop item
	comp.TagDropItem.Each(ecsWorld, func(dropEntry *donburi.Entry) {
		dropShape := comp.Body.Get(dropEntry).Shapes[0]
		Space.ShapeQuery(dropShape, func(shape *cm.Shape, points *cm.ContactPointSet) {
			e := shape.Body.UserData.(*donburi.Entry)
			if e.HasComponent(comp.TagBlock) {
				if shape.BB.Contains(dropShape.BB.Offset(vec2{-3, -3})) {
					ct := comp.StuckCountdown.Get(dropEntry)
					ct.Duration -= 1
				}
			}
		})
	})

	// Collision filter cooldown for item
	comp.CollisionTimer.Each(ecsWorld, func(e *donburi.Entry) {
		tm := comp.CollisionTimer.Get(e)
		if timerIsReady(tm) {
			shape := comp.Body.Get(e).Shapes[0]
			shape.SetShapeFilter(dropItemFilterCooldown)
		} else {
		}
	})
}

func (ps *Physics) Draw() {}

func PlayerDropItemBegin(arb *cm.Arbiter, s *cm.Space, dat any) bool {
	if checkEntries(arb) {
		a, b := getEntries(arb)
		inv := comp.Inventory.Get(a)
		itemData := comp.Item.Get(b)
		ok := addItem(inv, itemData.ID)
		if ok {
			destroyEntry(b)
		}
	}
	return false
}

func PlayerDropItemPre(arb *cm.Arbiter, _ *cm.Space, _ any) bool {
	if checkEntries(arb) {
		a, b := getEntries(arb)
		inv := comp.Inventory.Get(a)
		itemData := comp.Item.Get(b)
		ok := addItem(inv, itemData.ID)
		if ok {
			destroyEntry(b)
		}
	}
	return false
}

func PlayerDropItemPost(arb *cm.Arbiter, _ *cm.Space, dat any) {
	if checkEntries(arb) {
		a, b := getEntries(arb)
		inv := comp.Inventory.Get(a)
		itemData := comp.Item.Get(b)

		ok := addItem(inv, itemData.ID)
		if ok {
			destroyEntry(b)
		}
	}
}

func DropItemBlockBegin(arb *cm.Arbiter, _ *cm.Space, _ any) bool {
	dist := arb.ContactPointSet().Points[0].Distance
	if dist < -32 {
		return false
	} else {
		return true
	}
}

func DropItemBlockPre(arb *cm.Arbiter, _ *cm.Space, _ any) bool {
	dist := arb.ContactPointSet().Points[0].Distance
	if dist < -10 {
		fmt.Println("Pre -10")
	}
	return true
}

func DropItemBlockPost(arb *cm.Arbiter, _ *cm.Space, _ any) {
	dist := arb.ContactPointSet().Points[0].Distance
	if dist < -10 {
		fmt.Println("Post -10")
	}
}
