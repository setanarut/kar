package arche

import (
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/res"
	"kar/types"

	"github.com/yohamta/donburi"
)

func SpawnCircleBody(m, e, f, r float64, userData *donburi.Entry) *cm.Body {
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
	// 26, 40
	entry := res.World.Entry(res.World.Create(
		comp.PlayerTag,
		comp.AttackTimer,
		comp.Health,
		comp.Render,
		comp.Body,
		comp.Mobile,
		comp.WASDTag,
	))

	render := comp.Render.Get(entry)
	render.AnimPlayer = engine.NewAnimationPlayer(res.PlayerAtlas)
	render.AnimPlayer.AddStateAnimation("idle", 0, 0, 16, 16, 3, false, false).FPS = 3
	render.AnimPlayer.AddStateAnimation("dig_right", 0, 32, 16, 16, 6, false, false)
	render.AnimPlayer.AddStateAnimation("dig_down", 0, 48, 16, 16, 6, false, false)
	render.AnimPlayer.AddStateAnimation("walk_right", 0, 64, 16, 16, 6, false, false)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.CurrentScale = engine.GetBoxScaleFactor(16, 16, 50, 50)
	render.DrawScale = engine.GetBoxScaleFactor(16, 16, 50, 50)
	render.DrawScaleFlipX = engine.GetBoxScaleFactorFlipX(16, 16, 50, 50)

	b := spawnBoxBody(mass, el, fr, 30, 40, rad, entry)
	b.FirstShape().SetCollisionType(types.CollPlayer)
	b.FirstShape().Filter = cm.NewShapeFilter(0, types.BitmaskPlayer, cm.AllCategories&^types.BitmaskPlayerRaycast)
	b.SetPosition(pos)
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
	))
	wallShape.Body().UserData = entry
	comp.Body.Set(entry, wallShape.Body())
	return entry
}
