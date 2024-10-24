package system

import (
	"fmt"
	"kar"
	"kar/comp"
	"math"

	"github.com/setanarut/cm"
	"github.com/yohamta/donburi"

	"github.com/setanarut/vec"
)

var dropItemFilterCooldown = cm.ShapeFilter{
	Group:      2,
	Categories: kar.DropItemMask,
	Mask:       cm.AllCategories &^ kar.PlayerRayMask,
}

type Collision struct{}

func (cl *Collision) Init() {

	PlayerDropItemHandler := cmSpace.NewCollisionHandler(
		kar.PlayerCT,
		kar.DropItemCT)

	PlayerBlockHandler := cmSpace.NewCollisionHandler(
		kar.PlayerCT,
		kar.BlockCT)

	PlayerBlockHandler.PreSolveFunc = PlayerBlockPre

	PlayerDropItemHandler.BeginFunc = PlayerDropItemBegin
	// PlayerDropItemHandler.PreSolveFunc = PlayerDropItemPreCallback
	// PlayerDropItemHandler.PostSolveFunc = playerDropItemPostCallback

	DropItemBlockHandler := cmSpace.NewCollisionHandler(
		kar.DropItemCT,
		kar.BlockCT)
	DropItemBlockHandler.BeginFunc = DropItemBlockBegin
	DropItemBlockHandler.PostSolveFunc = DropItemBlockPost
}

func (cl *Collision) Update() {
	cmSpace.Step(kar.SpaceStep)

	// Destroy counter for stucked drop item
	comp.DropItem.Each(ecsWorld, func(dropEntry *donburi.Entry) {
		dropShape := comp.Body.Get(dropEntry).Shapes[0]
		cmSpace.ShapeQuery(dropShape, func(shape *cm.Shape, points *cm.ContactPointSet) {
			e := shape.Body.UserData.(*donburi.Entry)
			if e.HasComponent(comp.Block) {
				if shape.BB.Contains(dropShape.BB.Offset(vec.Vec2{-3, -3})) {
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

func (cl *Collision) Draw() {}

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

func DropItemBlockPost(arb *cm.Arbiter, _ *cm.Space, _ any) {
	dist := arb.ContactPointSet().Points[0].Distance
	if dist < -10 {
		fmt.Println("Post -10")
	}
}

func PlayerBlockPre(arb *cm.Arbiter, _ *cm.Space, _ any) bool {
	player, _ := arb.Bodies()
	vel := player.Velocity()

	if len(arb.ContactPointSet().Points) == 0 {
		return true
	}

	normal := arb.Normal()
	firstContactDepth := arb.ContactPointSet().Points[0].Distance
	dotProduct := vel.Dot(normal)
	// Eğer nesne zıt yönde hareket ediyorsa ve ilk temas noktasının kesişim
	//  mesafesi belirli bir eşik değerin altındaysa
	if dotProduct < 0 && math.Abs(firstContactDepth) < 0.1 {
		// Bu teması göz ardı et
		return false
	}

	return true
}
