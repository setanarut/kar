package arche

import (
	"image"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/items"
	"kar/res"
	"kar/types"
	"math"
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
		Group:      2,
		Categories: res.DropItemMask,
		Mask:       all &^ res.PlayerRayMask &^ res.PlayerMask,
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

	ibody := cm.NewBody(0.8, math.MaxFloat64)
	cm.NewCircleShapeWithBody(ibody, itemWidth, vec.Vec2{})
	ibody.Shapes[0].SetShapeFilter(dropItemFilter)
	ibody.Shapes[0].CollisionType = res.CollDropItem
	ibody.Shapes[0].SetElasticity(0)
	ibody.Shapes[0].SetFriction(1)

	ibody.SetPosition(pos)
	res.Space.AddBodyWithShapes(ibody)
	ibody.UserData = DropItemEntry
	comp.Body.Set(DropItemEntry, ibody)
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
	b.ShapeAtIndex(0).SetCollisionType(res.CollPlayer)
	b.ShapeAtIndex(0).SetShapeFilter(playerFilter)
	comp.Body.Set(e, b)
	return e
}

func SpawnBoxBody(pos vec.Vec2, m, e, f, w, h, r float64, en entry) *cm.Body {
	boxBody := cm.NewBody(m, cm.MomentForBox(m, w, h))
	cm.NewBoxShapeWithBody(boxBody, w, h, r)
	boxBody.Shapes[0].SetElasticity(e)
	boxBody.Shapes[0].SetFriction(f)
	boxBody.Shapes[0].SetCollisionType(res.CollEnemy)
	boxBody.SetPosition(pos)
	res.Space.AddBodyWithShapes(boxBody)
	boxBody.UserData = en
	return boxBody
}
func SpawnDebugBox(pos vec.Vec2) {
	e := res.ECSWorld.Entry(res.ECSWorld.Create(
		comp.DrawOptions,
		comp.Body,
		comp.TagDebugBox,
	))
	b := SpawnBoxBody(pos, 1, 0.2, 0.2, res.BlockSize, res.BlockSize, 0, e)
	b.ShapeAtIndex(0).SetShapeFilter(debugBoxBodyFilter)

	comp.DrawOptions.Set(e, &types.DrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        mathutil.GetRectScale(16, 16, res.BlockSize, res.BlockSize),
	})

	comp.Body.Set(e, b)

}

func spawnStatic(pos vec.Vec2, w, h float64) entry {
	sbody := cm.NewStaticBody()
	cm.NewBoxShapeWithBody(sbody, w, h, 0)
	sbody.Shapes[0].Filter = cm.NewShapeFilter(0, res.BlockMask, cm.AllCategories)
	sbody.Shapes[0].CollisionType = res.CollBlock
	sbody.Shapes[0].SetElasticity(0)
	sbody.Shapes[0].SetFriction(0.1)
	sbody.SetPosition(pos)
	res.Space.AddBodyWithShapes(sbody)
	// components
	entry := res.ECSWorld.Entry(res.ECSWorld.Create(
		comp.Body,
	))
	sbody.UserData = entry
	comp.Body.Set(entry, sbody)
	return entry
}
func spawnCircleBody(pos vec.Vec2, m, e, f, r float64, en entry) *cm.Body {
	body := cm.NewBody(m, math.MaxFloat64)
	cm.NewCircleShapeWithBody(body, r, vec.Vec2{})
	body.Shapes[0].SetElasticity(e)
	body.Shapes[0].SetFriction(f)
	body.SetPosition(pos)
	res.Space.AddBodyWithShapes(body)
	body.UserData = en
	return body
}
func worldPosToChunkCoord(worldPos vec.Vec2) image.Point {
	x := int((worldPos.X / res.ChunkSize.X) / res.BlockSize)
	y := int((worldPos.Y / res.ChunkSize.Y) / res.BlockSize)
	return image.Point{x, y}

}
