package arc

import (
	"image"
	"kar"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/items"
	"kar/res"
	"math"
	"time"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/setanarut/anim"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

type vec2 = vec.Vec2

var PlayerMapper5 = generic.NewMap5[Health, DrawOptions, anim.AnimationPlayer,
	cm.Body, Inventory](&kar.WorldECS)
var PlayerFilter5 = generic.NewFilter5[Health, DrawOptions, anim.AnimationPlayer,
	cm.Body, Inventory]()

var BlockMapper4 = generic.NewMap4[
	Health, DrawOptions, cm.Body, Item](&kar.WorldECS)
var BlockFilter = generic.NewFilter4[Health, DrawOptions, cm.Body, Item]()

var DropItemMapper6 = generic.NewMap6[
	DrawOptions, cm.Body, Item, CollisionTimer, Countdown, Index](&kar.WorldECS)

var DropItemFilter = generic.NewFilter6[DrawOptions, cm.Body, Item,
	CollisionTimer, Countdown, Index]()

var CountdownFilter = generic.NewFilter1[Countdown]()

var BodyID = ecs.ComponentID[cm.Body](&kar.WorldECS)

func NewInventory() *Inventory {
	inv := &Inventory{}
	inv.HandSlot = ItemStack{}
	for i := range inv.Slots {
		inv.Slots[i] = ItemStack{}
	}
	return inv
}

func SpawnBlock(pos vec2, id uint16) *cm.Shape {
	hlt := &Health{
		Health:    items.Property[id].MaxHealth,
		MaxHealth: items.Property[id].MaxHealth,
	}

	dop := &DrawOptions{
		CenterOffset: vec2{-8, -8},
		Scale:        mathutil.GetRectScale(16, 16, kar.BlockSize, kar.BlockSize),
	}

	itm := &Item{
		Chunk: WorldToChunk(pos),
		ID:    id,
	}

	body := cm.NewStaticBody()
	shape := cm.NewBoxShape(body, kar.BlockSize, kar.BlockSize, 0)
	shape.SetShapeFilter(BlockCollisionFilter)
	shape.SetElasticity(0)
	shape.SetFriction(0.1)
	shape.CollisionType = Block
	body.SetPosition(pos)
	kar.Space.AddBodyWithShapes(body)

	e := BlockMapper4.NewWith(hlt, dop, body, itm)
	body.UserData = e
	return shape
}

func SpawnDropItem(pos vec2, id uint16) ecs.Entity {

	itemWidth := kar.BlockSize / 3
	dop := &DrawOptions{
		CenterOffset: vec2{-8, -8},
		Scale:        mathutil.GetRectScale(16, 16, itemWidth, itemWidth),
	}

	collt := &CollisionTimer{Duration: time.Second / 2}
	ct := &Countdown{Duration: 120}
	idx := &Index{Index: 0}
	itm := &Item{
		Chunk: WorldToChunk(pos),
		ID:    id,
	}

	body := cm.NewBody(0.8, math.MaxFloat64)
	shape := cm.NewCircleShape(body, itemWidth, vec2{})
	shape.SetShapeFilter(DropItemCollisionFilter)
	shape.CollisionType = DropItem
	shape.SetElasticity(0)
	shape.SetFriction(1)

	body.SetPosition(pos)
	kar.Space.AddBodyWithShapes(body)
	e := DropItemMapper6.NewWith(dop, body, itm, collt, ct, idx)
	body.UserData = e
	return e
}

func SpawnPlayer(pos vec2) ecs.Entity {
	hlt := &Health{100, 100}

	ap := anim.NewAnimationPlayer(res.AtlasPlayer)
	ap.NewAnimationState("idleRight", 0, 0, 16, 16, 1, false, false).FPS = 1
	ap.NewAnimationState("idleFront", 16*2, 0, 16, 16, 1, false, false).FPS = 1
	ap.NewAnimationState("digRight", 0, 16*2, 16, 16, 6, false, false)
	ap.NewAnimationState("digDown", 0, 16*3, 16, 16, 6, false, false)
	ap.NewAnimationState("walkRight", 0, 16*4, 16, 16, 6, false, false)
	ap.NewAnimationState("falling", 16, 16*5, 16, 16, 1, false, false)

	playerScaleFactor := 3.0

	dop := &DrawOptions{
		CenterOffset: util.ImageCenterOffset(ap.CurrentFrame),
		Scale:        vec.Vec2{playerScaleFactor, playerScaleFactor},
	}

	body := cm.NewBody(1, math.MaxFloat64)
	verts := []vec.Vec2{{0, -5}, {5, 6}, {4, 7}, {4, 7}, {-4, 7}, {-5, 6}}
	geom := cm.NewTransformTranslate(vec.Vec2{0, -4})
	geom.Scale(playerScaleFactor, playerScaleFactor)

	shape := cm.NewPolyShape(body, verts, geom, 3)
	shape.Elasticity, shape.Friction = 0, 0
	shape.SetCollisionType(Player)
	shape.SetShapeFilter(PlayerCollisionFilter)

	body.SetPosition(pos)
	kar.Space.AddBodyWithShapes(body)

	e := PlayerMapper5.NewWith(hlt, dop, ap, body, NewInventory())
	body.UserData = e
	return e
}

func SpawnMario(pos vec2) ecs.Entity {
	hlt := &Health{100, 100}

	ap := anim.NewAnimationPlayer(res.Mario)
	ap.NewAnimationState("idleRight", 0, 0, 16, 16, 1, false, false).FPS = 1
	ap.NewAnimationState("walkRight", 16, 0, 16, 16, 4, false, false)
	ap.NewAnimationState("jump", 16*6, 0, 16, 16, 1, false, false)
	ap.NewAnimationState("skidding", 16*7, 0, 16, 16, 1, false, false)
	playerScaleFactor := 2.0

	dop := &DrawOptions{
		CenterOffset: util.ImageCenterOffset(ap.CurrentFrame),
		Scale:        vec.Vec2{playerScaleFactor, playerScaleFactor},
	}

	body := cm.NewBody(1, math.MaxFloat64)
	verts := []vec.Vec2{{0, -5}, {5, 6}, {4, 7}, {4, 7}, {-4, 7}, {-5, 6}}
	geom := cm.NewTransformTranslate(vec.Vec2{0, 0})
	geom.Scale(playerScaleFactor, playerScaleFactor)

	shape := cm.NewPolyShape(body, verts, geom, 3)
	shape.Elasticity, shape.Friction = 0, 0
	shape.SetCollisionType(Player)
	shape.SetShapeFilter(PlayerCollisionFilter)

	body.SetPosition(pos)
	kar.Space.AddBodyWithShapes(body)

	e := PlayerMapper5.NewWith(hlt, dop, ap, body, NewInventory())
	body.UserData = e
	return e
}

func WorldToChunk(pos vec.Vec2) image.Point {
	// pos = pos.Add(kar.BlockCenterOffset)
	return image.Point{
		int(math.Floor((pos.X / kar.ChunkSize.X) / kar.BlockSize)),
		int(math.Floor((pos.Y / kar.ChunkSize.Y) / kar.BlockSize))}
}
