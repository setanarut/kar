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

func SpawnCircleBody(m, e, f, r float64, userData *donburi.Entry) *cm.Body {
	// body := cm.NewKinematicBody()
	body := cm.NewBody(m, cm.Infinity)
	shape := cm.NewCircle(body, r, vec.Vec2{})
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

func SpawnPlayer(mass, el, fr, rad float64, pos vec.Vec2) *donburi.Entry {
	// 26, 40
	e := res.World.Entry(res.World.Create(
		comp.PlayerTag,
		comp.Health,
		comp.DrawOptions,
		comp.AnimationPlayer,
		comp.Body,
		comp.Mobile,
		comp.WASDTag,
	))

	ap := engine.NewAnimationPlayer(res.PlayerAtlas)
	ap.AddStateAnimation("idle", 0, 0, 16, 16, 3, false, false).FPS = 1
	ap.AddStateAnimation("dig_right", 0, 32, 16, 16, 6, false, false)
	ap.AddStateAnimation("dig_down", 0, 48, 16, 16, 6, false, false)
	ap.AddStateAnimation("walk_right", 0, 64, 16, 16, 6, false, false)
	comp.AnimationPlayer.Set(e, ap)

	comp.DrawOptions.Set(e, &types.DataDrawOptions{
		CenterOffset: util.EbitenImageCenterOffset(ap.CurrentFrame),
		Scale:        mathutil.RectangleScaleFactor(16, 16, 50, 50),
	})

	b := spawnBoxBody(mass, el, fr, 30, 40, rad, e)
	b.FirstShape().SetCollisionType(res.CollPlayer)
	b.FirstShape().Filter = cm.NewShapeFilter(0, res.BitmaskPlayer, cm.AllCategories&^res.BitmaskPlayerRaycast)
	b.SetPosition(pos)
	comp.Body.Set(e, b)
	return e
}

func SpawnWall(boxCenter vec.Vec2, boxW, boxH float64) *donburi.Entry {
	sbody := cm.NewStaticBody()
	wallShape := cm.NewBox(sbody, boxW, boxH, 0)
	wallShape.Filter = cm.NewShapeFilter(0, res.BitmaskWall, cm.AllCategories)
	wallShape.CollisionType = res.CollWall
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

func SpawnBlock(pos vec.Vec2, chunkCoord image.Point, blockType types.BlockType) {
	e := SpawnWall(pos, 50, 50)
	e.AddComponent(comp.Health)
	e.AddComponent(comp.Block)
	e.AddComponent(comp.DrawOptions)

	comp.DrawOptions.Set(e, &types.DataDrawOptions{
		CenterOffset: vec.Vec2{-8, -8},
		Scale:        mathutil.RectangleScaleFactor(16, 16, 50, 50),
	})

	comp.Block.Set(e,
		&types.DataBlock{
			ChunkCoord: chunkCoord,
			BlockType:  blockType,
		})

}

func SpawnDefaultPlayer(pos vec.Vec2) *donburi.Entry {
	return SpawnPlayer(1, 0, 0, 5, pos)

}
