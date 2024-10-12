package arche

import (
	"image"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/items"
	"kar/resources"
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
	resources.Space.AddShape(shape)
	resources.Space.AddBody(shape.Body())
	body.UserData = userData
	return body
}

func SpawnPlayer(pos vec.Vec2, mass, el, fr float64) *donburi.Entry {
	e := resources.ECSWorld.Entry(resources.ECSWorld.Create(
		comp.TagPlayer,
		comp.Health,
		comp.DrawOptions,
		comp.AnimPlayer,
		comp.Body,
		comp.Mobile,
		comp.Inventory,
		comp.TagWASD,
	))
	inv := &types.DataInventory{}

	for i := range inv.Slots {
		inv.Slots[i] = &types.ItemStack{0, 0}
	}

	comp.Inventory.Set(e, inv)

	ap := anim.NewAnimationPlayer(resources.AtlasPlayer)
	ap.NewAnimationState("idle", 0, 0, 16, 16, 3, false, false).FPS = 1
	ap.NewAnimationState("dig_right", 0, 16*2, 16, 16, 6, false, false)
	ap.NewAnimationState("dig_down", 0, 16*3, 16, 16, 6, false, false)
	ap.NewAnimationState("walk_right", 0, 16*4, 16, 16, 6, false, false)
	ap.NewAnimationState("jump", 16, 16*5, 16, 16, 1, false, false)
	comp.AnimPlayer.Set(e, ap)

	comp.DrawOptions.Set(e, &types.DataDrawOptions{
		CenterOffset: util.ImageCenterOffset(ap.CurrentFrame),
		Scale:        mathutil.CircleScaleFactor(resources.BlockSize/2, 16),
	})
	b := SpawnCircleBody(pos, mass, el, fr, (resources.BlockSize/2)*0.8, e)
	b.FirstShape().SetCollisionType(resources.CollPlayer)
	b.FirstShape().Filter = cm.NewShapeFilter(0, resources.BitmaskPlayer, cm.AllCategories&^resources.BitmaskPlayerRaycast)
	comp.Body.Set(e, b)
	return e
}

func SpawnStatic(pos vec.Vec2, w, h float64) *donburi.Entry {
	sbody := cm.NewStaticBody()
	wallShape := cm.NewBox(sbody, w, h, 0)
	wallShape.Filter = cm.NewShapeFilter(0, resources.BitmaskBlock, cm.AllCategories)
	wallShape.CollisionType = resources.CollBlock
	wallShape.SetElasticity(0)
	wallShape.SetFriction(0)
	sbody.SetPosition(pos)
	resources.Space.AddShape(wallShape)
	resources.Space.AddBody(wallShape.Body())
	// components
	entry := resources.ECSWorld.Entry(resources.ECSWorld.Create(
		comp.Body,
	))
	wallShape.Body().UserData = entry
	comp.Body.Set(entry, wallShape.Body())
	return entry
}

func SpawnBlock(pos vec.Vec2, chunkCoord image.Point, itemID uint16) {
	e := SpawnStatic(pos, resources.BlockSize, resources.BlockSize)
	e.AddComponent(comp.Health)
	e.AddComponent(comp.Item)
	e.AddComponent(comp.DrawOptions)
	e.AddComponent(comp.TagBlock)

	// set max health
	comp.Health.Set(e, &types.DataHealth{
		Health:    items.Property[itemID].MaxHealth,
		MaxHealth: items.Property[itemID].MaxHealth,
	})

	comp.DrawOptions.Set(e, &types.DataDrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        mathutil.RectangleScaleFactor(16, 16, resources.BlockSize, resources.BlockSize),
	})

	comp.Item.Set(e,
		&types.DataItem{
			ChunkCoord: chunkCoord,
			ID:         itemID,
		})

}

func SpawnBoxBody(pos vec.Vec2, m, e, f, w, h, r float64, userData *donburi.Entry) *cm.Body {
	body := cm.NewBody(m, cm.MomentForBox(m, w, h))
	shape := cm.NewBox(body, w, h, r)
	shape.SetElasticity(e)
	shape.SetFriction(f)
	shape.SetCollisionType(resources.CollEnemy)
	body.SetPosition(pos)
	resources.Space.AddShape(shape)
	resources.Space.AddBody(shape.Body())
	body.UserData = userData
	return body
}
func SpawnDebugBox(pos vec.Vec2) {
	e := resources.ECSWorld.Entry(resources.ECSWorld.Create(
		comp.DrawOptions,
		comp.Body,
		comp.TagDebugBox,
	))
	b := SpawnBoxBody(pos, 1, 0.2, 0.2, float64(resources.BlockSize), float64(resources.BlockSize), 0, e)
	b.FirstShape().Filter = cm.NewShapeFilter(0, resources.BitmaskEnemy, cm.AllCategories)

	comp.DrawOptions.Set(e, &types.DataDrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        mathutil.RectangleScaleFactor(16, 16, resources.BlockSize, resources.BlockSize),
	})

	comp.Body.Set(e, b)

}
func SpawnDropItem(pos vec.Vec2, itemID uint16, chunkCoord image.Point) {
	DropItemEntry := resources.ECSWorld.Entry(resources.ECSWorld.Create(
		comp.DrawOptions,
		comp.Body,
		comp.Item,
		comp.TagItem,
		comp.Index, // animasyon için zamanlayıcı
	))
	itemWidth := resources.BlockSize / 3
	comp.DrawOptions.Set(DropItemEntry, &types.DataDrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        mathutil.RectangleScaleFactor(16, 16, itemWidth, itemWidth),
	})

	comp.Index.Set(DropItemEntry, &types.DataIndex{Index: 0})

	comp.Item.Set(DropItemEntry,
		&types.DataItem{
			ChunkCoord: chunkCoord,
			ID:         itemID,
		})

	dropItemBody := cm.NewBody(0.8, cm.Infinity)
	dropItemShape := cm.NewCircle(dropItemBody, itemWidth, vec.Vec2{})
	dropItemShape.Filter = cm.NewShapeFilter(0,
		resources.BitmaskDropItem,
		cm.AllCategories&^resources.BitmaskPlayerRaycast&^resources.BitmaskDropItem)
	dropItemShape.CollisionType = resources.CollDropItem
	dropItemShape.SetElasticity(0.5)
	dropItemShape.SetFriction(0)
	dropItemBody.SetPosition(pos)
	resources.Space.AddShape(dropItemShape)
	resources.Space.AddBody(dropItemShape.Body())
	dropItemBody.UserData = DropItemEntry
	comp.Body.Set(DropItemEntry, dropItemBody)
}
