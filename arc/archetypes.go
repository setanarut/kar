package arc

import (
	"image"
	"kar"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/items"
	"kar/res"
	"math"

	"github.com/mlange-42/arche/ecs"
	gn "github.com/mlange-42/arche/generic"
	"github.com/setanarut/anim"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

var (
	MapHealth     = gn.NewMap[Health](&kar.WorldECS)
	MapAnimPlayer = gn.NewMap1[anim.AnimationPlayer](&kar.WorldECS)
	MapBody       = gn.NewMap1[CmBody](&kar.WorldECS)
	MapInventory  = gn.NewMap1[Inventory](&kar.WorldECS)
	MapItem       = gn.NewMap1[Item](&kar.WorldECS)
	MapBlock      = gn.NewMap4[Health, DrawOptions, CmBody, Item](&kar.WorldECS)
	MapPlayer     = gn.NewMap5[Health, DrawOptions, anim.AnimationPlayer, CmBody, Inventory](&kar.WorldECS)
	MapDropItem   = gn.NewMap6[
		DrawOptions,
		CmBody,
		Item,
		CollisionActivationCountdown,
		SelfDestuctionCountdown,
		AnimationFrameIndex](&kar.WorldECS)
)

var (
	FilterItem     = gn.NewFilter1[Item]()
	FilterBlock    = gn.NewFilter4[Health, DrawOptions, CmBody, Item]()
	FilterDropItem = gn.NewFilter6[
		DrawOptions,
		CmBody,
		Item,
		CollisionActivationCountdown,
		SelfDestuctionCountdown,
		AnimationFrameIndex]()
)

func SpawnBlock(pos vec.Vec2, id uint16) {

	hlt := &Health{
		Health:    items.Property[id].MaxHealth,
		MaxHealth: items.Property[id].MaxHealth,
	}

	dop := &DrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        vec.Vec2{3, 3},
		// Scale:        mathutil.GetRectScaleFactor(16, 16, kar.BlockSize, kar.BlockSize),
	}

	itm := &Item{
		Chunk: WorldToChunk(pos),
		ID:    id,
	}

	body := cm.NewStaticBody()
	shape := cm.NewBoxShape(body, kar.BlockSize, kar.BlockSize, 0)
	shape.SetShapeFilter(BlockCollisionFilter)
	shape.SetElasticity(0)
	shape.SetFriction(0)
	shape.CollisionType = Block
	body.SetPosition(pos)

	b := &CmBody{Body: body}
	kar.Space.AddBodyWithShapes(b.Body)

	e := MapBlock.NewWith(hlt, dop, b, itm)
	body.UserData = e
}

func SpawnDropItem(pos vec.Vec2, id uint16) ecs.Entity {

	itemWidth := kar.BlockSize / 3
	dop := &DrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        mathutil.GetRectScaleFactor(16, 16, itemWidth, itemWidth),
	}

	cac := &CollisionActivationCountdown{Tick: 30}
	ct := &SelfDestuctionCountdown{Tick: 120}
	idx := &AnimationFrameIndex{Index: 0}
	itm := &Item{
		Chunk: WorldToChunk(pos),
		ID:    id,
	}

	body := cm.NewBody(0.8, math.MaxFloat64)
	shape := cm.NewCircleShape(body, itemWidth, vec.Vec2{})
	shape.SetShapeFilter(DropItemCollisionFilter)
	shape.CollisionType = DropItem
	shape.SetElasticity(0)
	shape.SetFriction(1)
	body.SetPosition(pos)

	b := &CmBody{Body: body}
	kar.Space.AddBodyWithShapes(b.Body)

	e := MapDropItem.NewWith(dop, b, itm, cac, ct, idx)
	body.UserData = e
	return e
}

func SpawnMario(pos vec.Vec2) ecs.Entity {
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

	body := cm.NewBody(0.0001, math.MaxFloat64)
	// body := cm.NewKinematicBody()

	// verts := []vec.Vec2{{0, -5}, {5, 6}, {4, 7}, {4, 7}, {-4, 7}, {-5, 6}}
	// geom := cm.NewTransformTranslate(vec.Vec2{0, -2})
	// geom.Scale(playerScaleFactor+0.4, playerScaleFactor+0.4)
	// shape := cm.NewPolyShape(body, verts, geom, 0)

	shape := cm.NewBoxShape(body, 12*2, 16*2, 0)
	// shape := cm.NewCircleShape(body, 16, vec.Vec2{})

	shape.CollisionType = Player
	shape.Filter = PlayerCollisionFilter

	b := &CmBody{Body: body}
	kar.Space.AddBodyWithShapes(b.Body)

	e := MapPlayer.NewWith(hlt, dop, ap, b, NewInventory())
	b.Body.SetPosition(pos)
	b.Body.UserData = e
	return e
}

func WorldToChunk(pos vec.Vec2) image.Point {
	return image.Point{
		int(math.Floor((pos.X / kar.ChunkSize.X) / kar.BlockSize)),
		int(math.Floor((pos.Y / kar.ChunkSize.Y) / kar.BlockSize))}
}
