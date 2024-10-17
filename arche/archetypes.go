package arche

import (
	"image"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/items"
	"kar/res"
	"kar/types"
	"time"

	"github.com/setanarut/anim"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"

	"github.com/yohamta/donburi"
)

var all = cm.AllCategories

type entry = *donburi.Entry

var (
	dropItemFilter = cm.ShapeFilter{
		Group:      cm.NoGroup,
		Categories: res.DropItemMask,
		Mask:       all &^ res.PlayerRayMask &^ res.DropItemMask &^ res.PlayerMask,
	}
	playerFilter = cm.ShapeFilter{
		Group:      cm.NoGroup,
		Categories: res.PlayerMask,
		Mask:       all &^ res.PlayerRayMask,
	}
	debugBoxBodyFilter = cm.NewShapeFilter(cm.NoGroup, res.EnemyMask, all)
)

func SpawnItem(pos vec.Vec2, id uint16) {
	if items.IsBlock(id) {
		SpawnBlock(pos, id)
	} else {
		SpawnDropItem(pos, id)
	}
}

func SpawnBlock(pos vec.Vec2, id uint16) {
	e := spawnStatic(pos, res.BlockSize, res.BlockSize)
	e.AddComponent(comp.Health)
	e.AddComponent(comp.Item)
	e.AddComponent(comp.DrawOptions)

	if items.IsHarvestable(id) {
		e.AddComponent(comp.TagHarvestable)

	} else if items.IsBlock(id) {
		e.AddComponent(comp.TagBlock)
	}

	// set max health
	comp.Health.Set(e, &types.Health{
		Health:    items.Property[id].MaxHealth,
		MaxHealth: items.Property[id].MaxHealth,
	})

	comp.DrawOptions.Set(e, &types.DrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        mathutil.GetRectScale(16, 16, res.BlockSize, res.BlockSize),
	})

	comp.Item.Set(e,
		&types.Item{
			Chunk: worldPosToChunkCoord(pos),
			ID:    id,
		})
}

func SpawnDropItem(pos vec.Vec2, itemID uint16) entry {
	DropItemEntry := res.ECSWorld.Entry(res.ECSWorld.Create(
		comp.DrawOptions,
		comp.Body,
		comp.Item,
		comp.TagItem,
		comp.SpawnTimer,
		comp.Index, // animasyon için zamanlayıcı
	))
	itemWidth := res.BlockSize / 3
	comp.DrawOptions.Set(DropItemEntry, &types.DrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        mathutil.GetRectScale(16, 16, itemWidth, itemWidth),
	})

	comp.SpawnTimer.Set(DropItemEntry, &types.Timer{Duration: time.Second / 2})
	comp.Index.Set(DropItemEntry, &types.Index{Index: 0})
	comp.Item.Set(DropItemEntry,
		&types.Item{
			Chunk: worldPosToChunkCoord(pos),
			ID:    itemID,
		})

	dropItemBody := cm.NewBody(0.8, cm.Infinity)
	dropItemShape := cm.NewCircle(dropItemBody, itemWidth, vec.Vec2{})
	dropItemShape.Filter = dropItemFilter
	dropItemShape.CollisionType = res.CollDropItem
	dropItemShape.SetElasticity(0)
	dropItemShape.SetFriction(1)
	dropItemBody.SetPosition(pos)
	res.Space.AddShape(dropItemShape)
	res.Space.AddBody(dropItemShape.Body())
	dropItemBody.UserData = DropItemEntry
	comp.Body.Set(DropItemEntry, dropItemBody)
	return DropItemEntry
}

func SpawnPlayer(pos vec.Vec2, mass, el, fr float64) entry {
	e := res.ECSWorld.Entry(res.ECSWorld.Create(
		comp.TagPlayer,
		comp.Health,
		comp.DrawOptions,
		comp.AnimPlayer,
		comp.Body,
		comp.Mobile,
		comp.Inventory,
		comp.TagWASD,
	))

	inv := &types.Inventory{}
	for i := range inv.Slots {
		inv.Slots[i] = types.ItemStack{}
	}
	comp.Mobile.Set(e, &types.Mobile{Speed: 500, Accel: 80})
	inv.HandSlot = types.ItemStack{}
	comp.Inventory.Set(e, inv)

	ap := anim.NewAnimationPlayer(res.AtlasPlayer)
	ap.NewAnimationState("idle_right", 0, 0, 16, 16, 1, false, false).FPS = 1
	ap.NewAnimationState("idle_left", 16, 0, 16, 16, 1, false, false).FPS = 1
	ap.NewAnimationState("idle_front", 16*2, 0, 16, 16, 1, false, false).FPS = 1
	ap.NewAnimationState("dig_right", 0, 16*2, 16, 16, 6, false, false)
	ap.NewAnimationState("dig_down", 0, 16*3, 16, 16, 6, false, false)
	ap.NewAnimationState("walk_right", 0, 16*4, 16, 16, 6, false, false)
	ap.NewAnimationState("jump", 16, 16*5, 16, 16, 1, false, false)
	comp.AnimPlayer.Set(e, ap)

	comp.DrawOptions.Set(e, &types.DrawOptions{
		CenterOffset: util.ImageCenterOffset(ap.CurrentFrame),
		Scale:        mathutil.CircleScaleFactor(res.BlockSize/2, 16),
	})

	b := spawnCircleBody(pos, mass, el, fr, (res.BlockSize/2)*0.8, e)
	b.FirstShape().SetCollisionType(res.CollPlayer)
	b.FirstShape().SetShapeFilter(playerFilter)
	comp.Body.Set(e, b)
	return e
}

func SpawnBoxBody(pos vec.Vec2, m, e, f, w, h, r float64, en entry) *cm.Body {
	body := cm.NewBody(m, cm.MomentForBox(m, w, h))
	shape := cm.NewBox(body, w, h, r)
	shape.SetElasticity(e)
	shape.SetFriction(f)
	shape.SetCollisionType(res.CollEnemy)
	body.SetPosition(pos)
	res.Space.AddShape(shape)
	res.Space.AddBody(shape.Body())
	body.UserData = en
	return body
}
func SpawnDebugBox(pos vec.Vec2) {
	e := res.ECSWorld.Entry(res.ECSWorld.Create(
		comp.DrawOptions,
		comp.Body,
		comp.TagDebugBox,
	))
	b := SpawnBoxBody(pos, 1, 0.2, 0.2, res.BlockSize, res.BlockSize, 0, e)
	b.FirstShape().SetShapeFilter(debugBoxBodyFilter)

	comp.DrawOptions.Set(e, &types.DrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        mathutil.GetRectScale(16, 16, res.BlockSize, res.BlockSize),
	})

	comp.Body.Set(e, b)

}

func spawnStatic(pos vec.Vec2, w, h float64) entry {
	sbody := cm.NewStaticBody()
	wallShape := cm.NewBox(sbody, w, h, 0)
	wallShape.Filter = cm.NewShapeFilter(0, res.BlockMask, cm.AllCategories)
	wallShape.CollisionType = res.CollBlock
	wallShape.SetElasticity(0)
	wallShape.SetFriction(0.1)
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
func spawnCircleBody(pos vec.Vec2, m, e, f, r float64, en entry) *cm.Body {
	body := cm.NewBody(m, cm.Infinity)
	shape := cm.NewCircle(body, r, vec.Vec2{})
	shape.SetElasticity(e)
	shape.SetFriction(f)
	body.SetPosition(pos)
	res.Space.AddShape(shape)
	res.Space.AddBody(shape.Body())
	body.UserData = en
	return body
}
func worldPosToChunkCoord(worldPos vec.Vec2) image.Point {
	x := int((worldPos.X / res.ChunkSize.X) / res.BlockSize)
	y := int((worldPos.Y / res.ChunkSize.Y) / res.BlockSize)
	return image.Point{x, y}

}
