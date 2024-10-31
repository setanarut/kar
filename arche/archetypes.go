package arche

import (
	"image"
	"kar"
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

	db "github.com/yohamta/donburi"
)

// Collision Bitmask Category
const (
	PlayerBit uint = 1 << iota
	BlockBit
	DropItemBit
	PlayerRayBit
	EnemyBit
	BombBit
	BombRaycastBit
	GrabableBit uint = 1 << 31
	AllBits          = ^uint(0)
)

// Collision type
const (
	Player cm.CollisionType = iota
	Enemy
	Block
	DropItem
	Bomb
)

type vec2 = vec.Vec2

type entry = db.Entry

var (
	dropItemFilter = cm.ShapeFilter{
		Group:      2,
		Categories: DropItemBit,
		Mask:       AllBits &^ PlayerRayBit &^ PlayerBit,
	}
	playerFilter       = cm.ShapeFilter{0, PlayerBit, AllBits &^ PlayerRayBit}
	blockFilter        = cm.ShapeFilter{0, BlockBit, cm.AllCategories}
	debugBoxBodyFilter = cm.ShapeFilter{0, EnemyBit, AllBits}
)

func SpawnItem(s *cm.Space, w db.World, pos vec2, id uint16) {
	if items.IsBlock(id) {
		SpawnBlock(s, w, pos, id)
	} else {
		SpawnDropItem(s, w, pos, id)
	}
}

func SpawnBlock(s *cm.Space, w db.World, pos vec2, id uint16) *cm.Shape {
	e := w.Entry(w.Create(
		comp.Body,
		comp.Health,
		comp.Item,
		comp.DrawOptions,
	))

	if items.IsHarvestable(id) {
		e.AddComponent(comp.TagHarvestable)

	}
	if items.IsBlock(id) {
		e.AddComponent(comp.TagBlock)
	}

	if items.IsBreakable(id) {
		e.AddComponent(comp.TagBreakable)
	}

	// set max health
	comp.Health.Set(e, &types.Health{
		Health:    items.Property[id].MaxHealth,
		MaxHealth: items.Property[id].MaxHealth,
	})

	comp.DrawOptions.Set(e, &types.DrawOptions{
		CenterOffset: vec2{-8, -8},
		Scale:        mathutil.GetRectScale(16, 16, kar.BlockSize, kar.BlockSize),
	})

	comp.Item.Set(e,
		&types.Item{
			Chunk: worldToChunk(pos),
			ID:    id,
		})

	blockBody := cm.NewStaticBody()
	shape := cm.NewBoxShape(blockBody, kar.BlockSize, kar.BlockSize, 0)
	shape.SetShapeFilter(blockFilter)
	shape.SetElasticity(0)
	shape.SetFriction(0.1)
	shape.CollisionType = Block
	blockBody.SetPosition(pos)
	s.AddBodyWithShapes(blockBody)
	blockBody.UserData = e
	comp.Body.Set(e, blockBody)
	return shape
}

func SpawnDropItem(s *cm.Space, w db.World, pos vec2, id uint16) *entry {
	e := w.Entry(w.Create(
		comp.DrawOptions,
		comp.Body,
		comp.Item,
		comp.TagDropItem,
		comp.CollisionTimer,
		comp.StuckCountdown,
		comp.Index, // animasyon için zamanlayıcı
	))
	itemWidth := kar.BlockSize / 3
	comp.DrawOptions.Set(e, &types.DrawOptions{
		CenterOffset: vec2{-8, -8},
		Scale:        mathutil.GetRectScale(16, 16, itemWidth, itemWidth),
	})

	comp.CollisionTimer.Set(e, &types.Timer{Duration: time.Second / 2})
	comp.StuckCountdown.Set(e, &types.Countdown{Duration: 120})
	comp.Index.Set(e, &types.Index{Index: 0})
	comp.Item.Set(e,
		&types.Item{
			Chunk: worldToChunk(pos),
			ID:    id,
		})

	body := cm.NewBody(0.8, math.MaxFloat64)
	shape := cm.NewCircleShape(body, itemWidth, vec2{})
	shape.SetShapeFilter(dropItemFilter)
	shape.CollisionType = DropItem
	shape.SetElasticity(0)
	shape.SetFriction(1)

	body.SetPosition(pos)
	s.AddBodyWithShapes(body)
	body.UserData = e
	comp.Body.Set(e, body)
	return e
}

func SpawnPlayer(s *cm.Space, w db.World, pos vec2, mass, el, fr float64) *entry {
	e := w.Entry(w.Create(
		comp.TagPlayer,
		comp.Health,
		comp.DrawOptions,
		comp.AnimPlayer,
		comp.Body,
		comp.Inventory,
		comp.TagWASD,
	))

	inv := &types.Inventory{}
	for i := range inv.Slots {
		inv.Slots[i] = types.ItemStack{}
	}
	inv.HandSlot = types.ItemStack{}
	comp.Inventory.Set(e, inv)

	ap := anim.NewAnimationPlayer(res.AtlasPlayer)
	ap.NewAnimationState("idleRight", 0, 0, 16, 16, 1, false, false).FPS = 1
	ap.NewAnimationState("idleFront", 16*2, 0, 16, 16, 1, false, false).FPS = 1
	ap.NewAnimationState("digRight", 0, 16*2, 16, 16, 6, false, false)
	ap.NewAnimationState("digDown", 0, 16*3, 16, 16, 6, false, false)
	ap.NewAnimationState("walkRight", 0, 16*4, 16, 16, 6, false, false)
	ap.NewAnimationState("falling", 16, 16*5, 16, 16, 1, false, false)
	comp.AnimPlayer.Set(e, ap)

	playerScaleFactor := 3.0

	comp.DrawOptions.Set(e, &types.DrawOptions{
		CenterOffset: util.ImageCenterOffset(ap.CurrentFrame),
		Scale:        vec.Vec2{playerScaleFactor, playerScaleFactor},
	})

	b := cm.NewBody(mass, math.MaxFloat64)

	verts := []vec.Vec2{{0, -5}, {5, 6}, {4, 7}, {4, 7}, {-4, 7}, {-5, 6}}
	geom := cm.NewTransformTranslate(vec.Vec2{0, -4})
	geom.Scale(playerScaleFactor, playerScaleFactor)

	shape := cm.NewPolyShape(
		b,
		verts,
		geom,
		0,
	)

	// shape := cm.NewBoxShape(b, kar.BlockSize*0.7, (kar.BlockSize*2)*0.7, 2)
	shape.SetElasticity(el)
	shape.SetFriction(fr)
	b.SetPosition(pos)
	b.UserData = e
	b.ShapeAtIndex(0).SetCollisionType(Player)
	b.ShapeAtIndex(0).SetShapeFilter(playerFilter)
	s.AddBodyWithShapes(b)
	comp.Body.Set(e, b)
	return e
}

func SpawnDebugBox(s *cm.Space, w db.World, pos vec2) {

	e := w.Entry(w.Create(
		comp.DrawOptions,
		comp.Body,
		comp.TagDebugBox,
	))

	comp.DrawOptions.Set(e, &types.DrawOptions{
		CenterOffset: vec2{-8, -8},
		Scale:        mathutil.GetRectScale(16, 16, kar.BlockSize, kar.BlockSize),
	})

	b := cm.NewBody(1, cm.MomentForBox(1, kar.BlockSize, kar.BlockSize))
	shape := cm.NewBoxShape(b, kar.BlockSize, kar.BlockSize, 0)
	shape.Filter = debugBoxBodyFilter
	shape.SetElasticity(0.2)
	shape.SetFriction(0.2)
	shape.SetCollisionType(Enemy)
	b.SetPosition(pos)
	s.AddBodyWithShapes(b)
	b.UserData = e
	comp.Body.Set(e, b)
}

func worldToChunk(pos vec.Vec2) image.Point {
	// pos = pos.Add(kar.BlockCenterOffset)
	return image.Point{
		int(math.Floor((pos.X / kar.ChunkSize.X) / kar.BlockSize)),
		int(math.Floor((pos.Y / kar.ChunkSize.Y) / kar.BlockSize))}
}
