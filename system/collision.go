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

	PlayerDropItemHandler := kar.Space.NewCollisionHandler(arc.Player, arc.DropItem)
	PlayerDropItemHandler.BeginFunc = PlayerDropItemBegin
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
