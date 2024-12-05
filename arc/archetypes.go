package arc

import (
	"kar"
	"kar/res"

	"github.com/mlange-42/arche/ecs"
	gn "github.com/mlange-42/arche/generic"
	"github.com/setanarut/anim"
)

var (
	// MapHealth               = gn.NewMap[Health](&kar.WorldECS)
	// MapAnimPlayer           = gn.NewMap1[anim.AnimationPlayer](&kar.WorldECS)
	MapRect = gn.NewMap1[Rect](&kar.WorldECS)
	// MapInventory            = gn.NewMap1[Inventory](&kar.WorldECS)
	// MapPlatformerController = gn.NewMap1[PlatformerController](&kar.WorldECS)
	// MapDraw                 = gn.NewMap3[DrawOptions, anim.AnimationPlayer, Rect](&kar.WorldECS)
	MapPlayer = gn.NewMap6[
		Health,
		DrawOptions,
		anim.AnimationPlayer,
		Rect,
		Inventory,
		Controller](&kar.WorldECS)
)

var (
	FilterPlayer = gn.NewFilter5[
		Health,
		DrawOptions,
		anim.AnimationPlayer,
		Rect,
		Inventory]()

	FilterDraw       = gn.NewFilter3[DrawOptions, anim.AnimationPlayer, Rect]()
	FilterMovement   = gn.NewFilter4[Controller, Rect, anim.AnimationPlayer, DrawOptions]()
	FilterAnimPlayer = gn.NewFilter1[anim.AnimationPlayer]()
)

func init() {
}

func SpawnMario(x, y float64) ecs.Entity {
	h := &Health{100, 100}
	a := anim.NewAnimationPlayer(res.Mario)
	a.NewAnimationState("idleRight", 0, 0, 16, 16, 1, false, false).FPS = 1
	a.NewAnimationState("walkRight", 16, 0, 16, 16, 4, false, false)
	a.NewAnimationState("jump", 16*5, 0, 16, 16, 1, false, false)
	a.NewAnimationState("skidding", 16*6, 0, 16, 16, 1, false, false)

	a.SetState("idleRight")
	i := NewInventory()
	d := &DrawOptions{Scale: 3}
	r := &Rect{X: x, Y: y, W: 48, H: 48}
	c := NewController(0, 3)
	c.SetScale(3)
	entity := MapPlayer.NewWith(h, d, a, r, i, c)
	return entity
}
