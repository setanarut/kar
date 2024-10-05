package arche

import (
	"image"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/res"
	"kar/types"

	"github.com/setanarut/anim"
	"github.com/setanarut/cm"

	"github.com/setanarut/vec"

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

func SpawnPlayer(pos vec.Vec2, mass, el, fr float64) *donburi.Entry {
	e := res.ECSWorld.Entry(res.ECSWorld.Create(
		comp.PlayerTag,
		comp.Health,
		comp.DrawOptions,
		comp.AnimationPlayer,
		comp.Body,
		comp.Mobile,
		comp.Inventory,
		comp.WASDTag,
	))
	inv := &types.DataInventory{
		Items: make(map[types.ItemType]int),
	}
	comp.Inventory.Set(e, inv)

	ap := anim.NewAnimationPlayer(res.AtlasPlayer)
	ap.NewAnimationState("idle", 0, 0, 16, 16, 3, false, false).FPS = 1
	ap.NewAnimationState("dig_right", 0, 16*2, 16, 16, 6, false, false)
	ap.NewAnimationState("dig_down", 0, 16*3, 16, 16, 6, false, false)
	ap.NewAnimationState("walk_right", 0, 16*4, 16, 16, 6, false, false)
	ap.NewAnimationState("jump", 16, 16*5, 16, 16, 1, false, false)
	comp.AnimationPlayer.Set(e, ap)

	comp.DrawOptions.Set(e, &types.DataDrawOptions{
		CenterOffset: util.ImageCenterOffset(ap.CurrentFrame),
		Scale:        mathutil.CircleScaleFactor(res.BlockSize/2, 16),
	})
	b := SpawnCircleBody(pos, mass, el, fr, (res.BlockSize/2)*0.8, e)
	b.FirstShape().SetCollisionType(res.CollPlayer)
	b.FirstShape().Filter = cm.NewShapeFilter(0, res.BitmaskPlayer, cm.AllCategories&^res.BitmaskPlayerRaycast)
	comp.Body.Set(e, b)
	return e
}

func SpawnStatic(pos vec.Vec2, w, h float64) *donburi.Entry {
	sbody := cm.NewStaticBody()
	wallShape := cm.NewBox(sbody, w, h, 0)
	wallShape.Filter = cm.NewShapeFilter(0, res.BitmaskBlock, cm.AllCategories)
	wallShape.CollisionType = res.CollBlock
	wallShape.SetElasticity(0)
	wallShape.SetFriction(0)
	sbody.SetPosition(pos)
	res.Space.AddShape(wallShape)
	res.Space.AddBody(wallShape.Body())
	// components
	entry := res.ECSWorld.Entry(res.ECSWorld.Create(
		comp.Body,
	))
	wallShape.Body().UserData = entry
	comp.Body.Set(entry, wallShape.Body())
	return entry
}

func SpawnBlock(pos vec.Vec2, chunkCoord image.Point, blockType types.ItemType) {
	e := SpawnStatic(pos, res.BlockSize, res.BlockSize)
	e.AddComponent(comp.Health)
	e.AddComponent(comp.Item)
	e.AddComponent(comp.DrawOptions)
	e.AddComponent(comp.BlockTag)

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

func SpawnBoxBody(pos vec.Vec2, m, e, f, w, h, r float64, userData *donburi.Entry) *cm.Body {
	body := cm.NewBody(m, cm.MomentForBox(m, w, h))
	shape := cm.NewBox(body, w, h, r)
	shape.SetElasticity(e)
	shape.SetFriction(f)
	shape.SetCollisionType(res.CollEnemy)
	body.SetPosition(pos)
	res.Space.AddShape(shape)
	res.Space.AddBody(shape.Body())
	body.UserData = userData
	return body
}
func SpawnDebugBox(pos vec.Vec2) {
	e := res.ECSWorld.Entry(res.ECSWorld.Create(
		comp.DrawOptions,
		comp.Body,
		comp.DebugBoxTag,
	))
	b := SpawnBoxBody(pos, 1, 0.2, 0.2, float64(res.BlockSize), float64(res.BlockSize), 0, e)
	b.FirstShape().Filter = cm.NewShapeFilter(0, res.BitmaskEnemy, cm.AllCategories)

	comp.DrawOptions.Set(e, &types.DataDrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        mathutil.RectangleScaleFactor(16, 16, res.BlockSize, res.BlockSize),
	})

	comp.Body.Set(e, b)

}
func SpawnDropItem(pos vec.Vec2, item types.ItemType, chunkCoord image.Point) {
	DropItemEntry := res.ECSWorld.Entry(res.ECSWorld.Create(
		comp.DrawOptions,
		comp.Body,
		comp.Item,
		comp.DropItemTag,
		comp.Index, // animasyon için zamanlayıcı
	))
	itemWidth := res.BlockSize / 4
	comp.DrawOptions.Set(DropItemEntry, &types.DataDrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        mathutil.RectangleScaleFactor(16, 16, itemWidth, itemWidth),
	})

	comp.Index.Set(DropItemEntry, &types.DataIndex{Index: 0})

	comp.Item.Set(DropItemEntry,
		&types.DataItem{
			ChunkCoord: chunkCoord,
			Item:       item,
		})

	dropItemBody := cm.NewBody(0.8, cm.Infinity)
	dropItemShape := cm.NewCircle(dropItemBody, itemWidth, vec.Vec2{})
	dropItemShape.Filter = cm.NewShapeFilter(0,
		res.BitmaskDropItem,
		cm.AllCategories&^res.BitmaskPlayerRaycast&^res.BitmaskDropItem)
	dropItemShape.CollisionType = res.CollDropItem
	dropItemShape.SetElasticity(0.5)
	dropItemShape.SetFriction(0)
	dropItemBody.SetPosition(pos)
	res.Space.AddShape(dropItemShape)
	res.Space.AddBody(dropItemShape.Body())
	dropItemBody.UserData = DropItemEntry
	comp.Body.Set(DropItemEntry, dropItemBody)
}
