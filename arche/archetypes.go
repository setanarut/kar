package arche

import (
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/res"
	"kar/types"

	"github.com/yohamta/donburi"
)

func spawnCircleBody(m, e, f, r float64, userData *donburi.Entry) *cm.Body {
	// body := cm.NewKinematicBody()
	body := cm.NewBody(m, cm.Infinity)
	shape := cm.NewCircle(body, r, cm.Vec2{})
	shape.SetElasticity(e)
	shape.SetFriction(f)
	res.Space.AddShape(shape)
	res.Space.AddBody(shape.Body())
	body.UserData = userData
	return body
}
func spawnBoxBody(m, e, f, w, h, r float64, userData *donburi.Entry) *cm.Body {
	// body := cm.NewKinematicBody()
	body := cm.NewBody(m, cm.Infinity)
	shape := cm.NewBox(body, w, h, r)
	shape.SetElasticity(e)
	shape.SetFriction(f)
	res.Space.AddShape(shape)
	res.Space.AddBody(shape.Body())
	body.UserData = userData
	return body
}

func SpawnPlayer(mass, el, fr, rad float64, pos cm.Vec2) *donburi.Entry {

	entry := res.World.Entry(res.World.Create(
		comp.PlayerTag,
		comp.Inventory,
		comp.AttackTimer,
		comp.Health,
		comp.Damage,
		comp.Sprite,
		comp.Body,
		comp.Mobile,
		comp.WASDTag,
	))
	comp.Health.SetValue(entry, 100000)
	comp.Inventory.Set(entry, &types.DataInventory{
		Items: make(map[types.ItemType]int),
	})

	i := comp.Inventory.Get(entry)
	i.Items[types.ItemSnowball] = 1000

	sprite := comp.Sprite.Get(entry)
	sprite.Image = res.StoneBlockImage
	sprite.Offset = engine.GetEbitenImageOffset(sprite.Image)
	sprite.DrawScale = engine.GetBoxScaleFactor(16, 16, 25, 25)

	b := spawnBoxBody(mass, el, fr, 25, 25, rad, entry)
	b.FirstShape().SetCollisionType(types.CollPlayer)
	b.FirstShape().Filter = cm.NewShapeFilter(0, types.BitmaskPlayer, cm.AllCategories&^types.BitmaskSnowball)
	b.SetPosition(pos)
	comp.Body.Set(entry, b)
	res.CurrentTool = types.ItemSnowball
	// b.FirstShape().SetSensor(true)
	return entry
}

// func SpawnSnowball(m, e, f, r float64, pos cm.Vec2) *donburi.Entry {
// 	entry := res.World.Entry(res.World.Create(
// 		comp.SnowballTag,
// 		comp.Render,
// 		comp.Damage,
// 		comp.Body,
// 	))

// 	render := comp.Render.Get(entry)
// 	render.AnimPlayer = engine.NewAnimationPlayer(res.Player)
// 	render.AnimPlayer.AddStateAnimation("idle", 0, 0, 100, 100, 1, false, false)
// 	render.DrawScale = engine.GetCircleScaleFactor(r, render.AnimPlayer.CurrentFrame)
// 	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
// 	render.AnimPlayer.Paused = true

// 	b := spawnCircleBody(m, e, f, r, entry)
// 	b.SetPosition(pos)
// 	b.FirstShape().SetCollisionType(types.CollSnowball)
// 	b.FirstShape().Filter = cm.NewShapeFilter(0, types.BitmaskSnowball, cm.AllCategories&^types.BitmaskPlayer&^types.BitmaskSnowball)
// 	// b.FirstShape().SetSensor(true)
// 	comp.Body.Set(entry, b)

// 	return entry
// }

func SpawnWall(boxCenter cm.Vec2, boxW, boxH float64) *donburi.Entry {
	sbody := cm.NewStaticBody()
	wallShape := cm.NewBox(sbody, boxW, boxH, 0)
	wallShape.Filter = cm.NewShapeFilter(0, types.BitmaskWall, cm.AllCategories)
	wallShape.CollisionType = types.CollWall
	wallShape.SetElasticity(0)
	wallShape.SetFriction(0)
	sbody.SetPosition(boxCenter)
	res.Space.AddShape(wallShape)
	res.Space.AddBody(wallShape.Body())

	// components
	entry := res.World.Entry(res.World.Create(
		comp.Body,
		comp.Sprite,
		comp.ChunkCoord,
	))

	wallShape.Body().UserData = entry
	comp.Body.Set(entry, wallShape.Body())

	sprite := comp.Sprite.Get(entry)
	sprite.Image = res.StoneBlockImage
	sprite.Offset = engine.GetEbitenImageOffset(sprite.Image)
	sprite.DrawScale = engine.GetBoxScaleFactor(16, 16, 50, 50)

	return entry
}
