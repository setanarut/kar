package arc

import (
	"kar"
	"kar/res"

	"github.com/mlange-42/arche/ecs"
	gn "github.com/mlange-42/arche/generic"
	"github.com/setanarut/anim"
)

var (
	MapInventory = gn.NewMap1[Inventory](&kar.WorldECS)
	MapRect      = gn.NewMap1[Rect](&kar.WorldECS)
	MapHealth    = gn.NewMap1[Health](&kar.WorldECS)
	MapItem      = gn.NewMap4[ItemID, Durability, Rect, ItemTimers](&kar.WorldECS)
	MapPlayer    = gn.NewMap5[
		anim.AnimationPlayer,
		Health,
		DrawOptions,
		Rect,
		Inventory](&kar.WorldECS)
)

// Query Filters
var (
	FilterRect       = gn.NewFilter1[Rect]()
	FilterAnimPlayer = gn.NewFilter1[anim.AnimationPlayer]()
	FilterItem       = gn.NewFilter4[ItemID, Rect, ItemTimers, Durability]()
	FilterPlayer     = gn.NewFilter5[
		anim.AnimationPlayer,
		Health,
		DrawOptions,
		Rect,
		Inventory]()
)

func init() {
	FilterRect.Register(&kar.WorldECS)
	FilterAnimPlayer.Register(&kar.WorldECS)
	FilterItem.Register(&kar.WorldECS)
	FilterPlayer.Register(&kar.WorldECS)
}

func SpawnPlayer(x, y float64) ecs.Entity {
	AP := anim.NewAnimationPlayer(res.PlayerAtlas)
	AP.NewAnimationState("idleRight", 0, 0, 16, 16, 1, false, false).FPS = 1
	AP.NewAnimationState("idleUp", 208, 0, 16, 16, 1, false, false).FPS = 1
	AP.NewAnimationState("idleDown", 224, 0, 16, 16, 1, false, false).FPS = 1
	AP.NewAnimationState("walkRight", 16, 0, 16, 16, 4, false, false)
	AP.NewAnimationState("jump", 16*5, 0, 16, 16, 1, false, false)
	AP.NewAnimationState("skidding", 16*6, 0, 16, 16, 1, false, false)
	AP.NewAnimationState("attackDown", 16*7, 0, 16, 16, 2, false, false).FPS = 8
	AP.NewAnimationState("attackRight", 144, 0, 16, 16, 2, false, false).FPS = 8
	AP.NewAnimationState("attackWalk", 0, 16, 16, 16, 4, false, false).FPS = 8
	AP.NewAnimationState("attackUp", 16*11, 0, 16, 16, 2, false, false).FPS = 8
	AP.SetState("idleRight")

	inv := NewInventory()
	inv.RandomFillAllSlots()
	inv.ClearSlot(0)
	inv.ClearSlot(1)
	inv.ClearSlot(4)

	return MapPlayer.NewWith(
		AP,
		&Health{20, 20},
		&DrawOptions{Scale: kar.PlayerScale},
		&Rect{x, y, 16 * kar.PlayerScale, 16 * kar.PlayerScale},
		inv,
	)
}
func SpawnItem(data SpawnData) ecs.Entity {
	return MapItem.NewWith(
		&ItemID{data.Id},
		&Durability{data.Durability},
		&Rect{data.X, data.Y, 8 * kar.ItemScale, 8 * kar.ItemScale},
		&ItemTimers{kar.ItemCollisionDelay, 0},
	)
}

// SpawnData is a helper for delaying spawn events
type SpawnData struct {
	X, Y       float64
	Id         uint16
	Durability int
}
