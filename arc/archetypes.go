package arc

import (
	"kar"
	"kar/engine/v"
	"kar/res"

	"github.com/mlange-42/arche/ecs"
	gn "github.com/mlange-42/arche/generic"
	"github.com/setanarut/anim"
)

type vec2 = v.Vec

var (
	MapHealth     = gn.NewMap[Health](&kar.WorldECS)
	MapAnimPlayer = gn.NewMap1[anim.AnimationPlayer](&kar.WorldECS)
	MapRect       = gn.NewMap1[Rect](&kar.WorldECS)
	MapInventory  = gn.NewMap1[Inventory](&kar.WorldECS)
	MapPlayer     = gn.NewMap6[
		Health,
		DrawOptions,
		anim.AnimationPlayer,
		Rect,
		Inventory,
		Controller](&kar.WorldECS)
)

func SpawnMario(x, y float64) ecs.Entity {
	h := &Health{100, 100}
	a := anim.NewAnimationPlayer(res.Mario)
	a.NewAnimationState("idleRight", 0, 0, 16, 16, 1, false, false).FPS = 1
	a.NewAnimationState("walkRight", 16, 0, 16, 16, 4, false, false)
	a.NewAnimationState("jump", 16*6, 0, 16, 16, 1, false, false)
	a.NewAnimationState("skidding", 16*7, 0, 16, 16, 1, false, false)
	i := NewInventory()
	d := &DrawOptions{Scale: 2}
	r := &Rect{X: x, Y: y, W: 16, H: 16}
	c := NewController()
	entity := MapPlayer.NewWith(h, d, a, r, i, c)
	return entity
}
