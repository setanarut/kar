package arche

import (
	"image"
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/engine/vec"
	"kar/res"
	"kar/types"

	"github.com/yohamta/donburi"
)

func SpawnCircleBody(pos vec.Vec2, m, e, f, r float64, userData *donburi.Entry) *cm.Body {
	// body := cm.NewKinematicBody()
	body := cm.NewBody(m, cm.Infinity)
	shape := cm.NewCircle(body, r, vec.Vec2{})
	shape.SetElasticity(e)
	shape.SetFriction(f)
	body.SetPosition(pos)
	res.Space.AddShape(shape)
	res.Space.AddBody(shape.Body())
	body.UserData = userData
	return body
}
func SpawnBoxBody(pos vec.Vec2, m, e, f, w, h, r float64, userData *donburi.Entry) *cm.Body {
	// body := cm.NewKinematicBody()
	body := cm.NewBody(m, cm.Infinity)
	shape := cm.NewBox(body, w, h, r)
	shape.SetElasticity(e)
	shape.SetFriction(f)
	body.SetPosition(pos)
	res.Space.AddShape(shape)
	res.Space.AddBody(shape.Body())
	body.UserData = userData
	return body
}

func SpawnPlayer(pos vec.Vec2, mass, el, fr float64) *donburi.Entry {
	e := res.World.Entry(res.World.Create(
		comp.PlayerTag,
		comp.Health,
		comp.DrawOptions,
		comp.AnimationPlayer,
		comp.Body,
		comp.Mobile,
		comp.WASDTag,
	))

	ap := engine.NewAnimationPlayer(res.AtlasPlayer)
	ap.AddStateAnimation("idle", 0, 0, 16, 16, 3, false, false).FPS = 1
	ap.AddStateAnimation("dig_right", 0, 32, 16, 16, 6, false, false)
	ap.AddStateAnimation("dig_down", 0, 48, 16, 16, 6, false, false)
	ap.AddStateAnimation("walk_right", 0, 64, 16, 16, 6, false, false)
	comp.AnimationPlayer.Set(e, ap)

	comp.DrawOptions.Set(e, &types.DataDrawOptions{
		CenterOffset: util.EbitenImageCenterOffset(ap.CurrentFrame),
		Scale:        mathutil.CircleScaleFactor(res.BlockSize/2, 16),
	})
	// b := SpawnBoxBody(mass, el, fr, 16, 16, rad, e)
	b := SpawnCircleBody(pos, mass, el, fr, (res.BlockSize/2)*0.8, e)
	b.FirstShape().SetCollisionType(res.CollPlayer)
	b.FirstShape().Filter = cm.NewShapeFilter(0, res.BitmaskPlayer, cm.AllCategories&^res.BitmaskPlayerRaycast)
	comp.Body.Set(e, b)
	return e
}

func SpawnWall(pos vec.Vec2, boxW, boxH float64) *donburi.Entry {
	sbody := cm.NewStaticBody()
	wallShape := cm.NewBox(sbody, boxW, boxH, 0)
	wallShape.Filter = cm.NewShapeFilter(0, res.BitmaskWall, cm.AllCategories)
	wallShape.CollisionType = res.CollWall
	wallShape.SetElasticity(0)
	wallShape.SetFriction(0)
	sbody.SetPosition(pos)
	res.Space.AddShape(wallShape)
	res.Space.AddBody(wallShape.Body())
	// components
	entry := res.World.Entry(res.World.Create(
		comp.Body,
	))
	wallShape.Body().UserData = entry
	comp.Body.Set(entry, wallShape.Body())
	return entry
}

func SpawnBlock(pos vec.Vec2, chunkCoord image.Point, blockType types.ItemType) {
	e := SpawnWall(pos, float64(res.BlockSize), float64(res.BlockSize))
	e.AddComponent(comp.Health)
	e.AddComponent(comp.Item)
	e.AddComponent(comp.DrawOptions)
	e.AddComponent(comp.BlockItemTag)

	// set max health
	comp.Health.Set(e, &types.DataHealth{
		Health:    res.BlockMaxHealth[blockType],
		MaxHealth: res.BlockMaxHealth[blockType],
	})

	comp.DrawOptions.Set(e, &types.DataDrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        mathutil.RectangleScaleFactor(16, 16, res.BlockSize, res.BlockSize),
	})

	comp.Item.Set(e,
		&types.DataItem{
			ChunkCoord: chunkCoord,
			Item:       blockType,
		})

}

func SpawnDebug(pos vec.Vec2) {
	e := res.World.Entry(res.World.Create(
		comp.DrawOptions,
		comp.Body,
		comp.AnimationPlayer,
	))

	comp.DrawOptions.Set(e, &types.DataDrawOptions{
		CenterOffset: vec.Vec2{-16, -16},
		Scale:        vec.Vec2{1, 1},
	})
	ap := engine.NewAnimationPlayer(res.AtlasBlock)
	ap.AddStateAnimation("idle", 128, 128, 32, 32, 6, true, false)
	comp.AnimationPlayer.Set(e, ap)

	b := SpawnCircleBody(pos, 0.8, 0.5, 0, 16, e)
	b.FirstShape().Filter = cm.NewShapeFilter(
		0,
		res.BitmaskCollectible,
		cm.AllCategories&^res.BitmaskPlayerRaycast&^res.BitmaskCollectible)
	b.FirstShape().CollisionType = res.CollCollectible
	comp.Body.Set(e, b)

}
func SpawnItem(pos vec.Vec2, i types.ItemType, chunkCoord image.Point) {
	e := res.World.Entry(res.World.Create(
		comp.DrawOptions,
		comp.Body,
		comp.Item,
		comp.DropItemTag,
	))

	comp.DrawOptions.Set(e, &types.DataDrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        vec.Vec2{1, 1},
	})

	comp.Item.Set(e,
		&types.DataItem{
			ChunkCoord: chunkCoord,
			Item:       i,
		})

	b := SpawnCircleBody(pos, 0.8, 0.5, 0, 8, e)
	b.FirstShape().Filter = cm.NewShapeFilter(
		0,
		res.BitmaskCollectible,
		cm.AllCategories&^res.BitmaskPlayerRaycast&^res.BitmaskCollectible)
	b.FirstShape().CollisionType = res.CollCollectible
	comp.Body.Set(e, b)

}
