package arche

import (
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/res"
	"kar/types"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

func spawnBody(m, e, f, r float64, userData *donburi.Entry) *cm.Body {
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

func SpawnPlayer(mass, el, fr, rad float64, pos cm.Vec2) *donburi.Entry {

	e := res.World.Entry(res.World.Create(
		comp.PlayerTag,
		comp.Inventory,
		comp.AttackTimer,
		comp.Health,
		comp.Damage,
		comp.Render,
		comp.Body,
		comp.Mobile,
		comp.WASDTag,
	))
	comp.Health.SetValue(e, 100000)
	comp.Inventory.Set(e, &types.DataInventory{
		Items: make(map[types.ItemType]int),
	})

	i := comp.Inventory.Get(e)
	i.Items[types.ItemSnowball] = 1000

	w := 100
	render := comp.Render.Get(e)
	render.AnimPlayer = engine.NewAnimationPlayer(res.Player)
	render.AnimPlayer.AddStateAnimation("shoot", 0, 0, w, w, 1, false)
	render.AnimPlayer.AddStateAnimation("right", 0, 0, w, w, 1, false)
	render.AnimPlayer.SetFPS(9)
	render.DIO.Filter = ebiten.FilterLinear
	render.DrawScale = engine.GetCircleScaleFactor(rad, render.AnimPlayer.CurrentFrame)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.ScaleColor = colornames.White

	b := spawnBody(mass, el, fr, rad, e)
	b.FirstShape().SetCollisionType(types.CollPlayer)
	b.FirstShape().Filter = cm.NewShapeFilter(0, types.BitmaskPlayer, cm.AllCategories&^types.BitmaskSnowball)
	b.SetPosition(pos)
	comp.Body.Set(e, b)
	res.CurrentTool = types.ItemSnowball
	b.FirstShape().SetSensor(true)
	return e
}
func SpawnMob(m, e, f, r float64, pos cm.Vec2) *donburi.Entry {

	entry := res.World.Entry(res.World.Create(
		comp.EnemyTag,
		comp.AI,
		comp.Health,
		comp.Damage,
		comp.Render,
		comp.Mobile,
		comp.Body,
	))

	render := comp.Render.Get(entry)
	render.AnimPlayer = engine.NewAnimationPlayer(res.Player)
	render.AnimPlayer.AddStateAnimation("idle", 0, 0, 100, 100, 1, false)

	render.DrawScale = engine.GetCircleScaleFactor(r, render.AnimPlayer.CurrentFrame)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.AnimPlayer.Paused = true

	b := spawnBody(m, e, f, r, entry)
	b.SetPosition(pos)
	b.FirstShape().SetCollisionType(types.CollEnemy)
	b.FirstShape().Filter = cm.NewShapeFilter(0, types.BitmaskEnemy, cm.AllCategories)
	comp.Body.Set(entry, b)

	return entry
}

func SpawnSnowball(m, e, f, r float64, pos cm.Vec2) *donburi.Entry {
	entry := res.World.Entry(res.World.Create(
		comp.SnowballTag,
		comp.Render,
		comp.Damage,
		comp.Body,
	))

	render := comp.Render.Get(entry)
	render.AnimPlayer = engine.NewAnimationPlayer(res.Player)
	render.AnimPlayer.AddStateAnimation("idle", 0, 0, 100, 100, 1, false)
	render.DrawScale = engine.GetCircleScaleFactor(r, render.AnimPlayer.CurrentFrame)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.AnimPlayer.Paused = true

	b := spawnBody(m, e, f, r, entry)
	b.SetPosition(pos)
	b.FirstShape().SetCollisionType(types.CollSnowball)
	b.FirstShape().Filter = cm.NewShapeFilter(0, types.BitmaskSnowball, cm.AllCategories&^types.BitmaskPlayer&^types.BitmaskSnowball)
	// b.FirstShape().SetSensor(true)
	comp.Body.Set(entry, b)

	return entry
}

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
		comp.Render,
		comp.ChunkCoord,
	))

	wallShape.Body().UserData = entry
	comp.Body.Set(entry, wallShape.Body())

	render := comp.Render.Get(entry)
	render.DIO.Filter = ebiten.FilterNearest
	render.AnimPlayer = engine.NewAnimationPlayer(res.Wall)
	imW := res.Wall.Bounds().Dx()
	render.AnimPlayer.AddStateAnimation("idle", 0, 0, imW, imW, 1, false)
	render.DrawScale = engine.GetBoxScaleFactor(float64(imW), float64(imW), boxW, boxH)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.AnimPlayer.Paused = true
	render.ScaleColor = colornames.Gray
	return entry
}
