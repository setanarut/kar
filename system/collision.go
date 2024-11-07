package system

import (
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

	PlayerDropItemHandler := &cm.CollisionHandler{
		TypeA: arc.Player,
		TypeB: arc.DropItem,
		BeginFunc: func(arb *cm.Arbiter, _ *cm.Space, _ any) bool {
			a, b := arb.Bodies()
			dropItem := b.UserData.(ecs.Entity)
			itm := arc.MapItem.Get(dropItem)
			ok := addItem(arc.MapInventory.Get(a.UserData.(ecs.Entity)), itm.ID)
			if ok {
				kar.Space.AddPostStepCallback(removeBodyPostStep, b, nil)
				kar.WorldECS.RemoveEntity(dropItem)
			}
			return false
		},
		PreSolveFunc:  cm.DefaultPreSolve,
		PostSolveFunc: cm.DefaultPostSolve,
		SeparateFunc:  cm.DefaultSeparate,
		UserData:      nil,
	}
	kar.Space.AddCollisionHandler2(PlayerDropItemHandler)
}

func (ps *Physics) Update() {
	kar.Space.Step(kar.DeltaTime)
}
func (ps *Physics) Draw() {}
