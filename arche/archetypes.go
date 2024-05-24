package arche

import (
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/res"
	"kar/types"

	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

func spawnBody(m, e, f, r float64, userData *donburi.Entry) *cm.Body {
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
	render.AnimPlayer.AddStateAnimation("right", 0, 0, w, w, 4, true)
	render.AnimPlayer.SetFPS(9)

	render.DrawScale = engine.GetCircleScaleFactor(rad, render.AnimPlayer.CurrentFrame)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.ScaleColor = colornames.White

	b := spawnBody(mass, el, fr, rad, e)
	b.FirstShape().SetCollisionType(types.CollPlayer)
	b.FirstShape().Filter = cm.NewShapeFilter(0, types.BitmaskPlayer, cm.AllCategories&^types.BitmaskSnowball)
	b.SetPosition(pos)
	// b.SetVelocityUpdateFunc(comp.PlayerVelocityFunc)
	comp.Body.Set(e, b)
	res.CurrentTool = types.ItemSnowball
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
	render.AnimPlayer = engine.NewAnimationPlayer(res.EnemyBody)
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
func SpawnBomb(m, e, f, r float64, pos cm.Vec2) *donburi.Entry {

	entry := res.World.Entry(res.World.Create(
		comp.EnemyTag,
		comp.AI,
		comp.Inventory,
		comp.Health,
		comp.Damage,
		comp.Render,
		comp.Body,
	))

	render := comp.Render.Get(entry)
	render.AnimPlayer = engine.NewAnimationPlayer(res.Items)
	render.AnimPlayer.AddStateAnimation("idle", 0, 0, 100, 100, 1, false)

	render.DrawScale = engine.GetCircleScaleFactor(r, render.AnimPlayer.CurrentFrame)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.AnimPlayer.Paused = true

	b := spawnBody(m, e, f, r, entry)
	b.FirstShape().SetCollisionType(types.CollBomb)
	b.FirstShape().Filter = cm.NewShapeFilter(0, types.BitmaskBomb, cm.AllCategories)
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
	render.AnimPlayer = engine.NewAnimationPlayer(res.Items)
	render.AnimPlayer.AddStateAnimation("idle", 200, 0, 100, 100, 1, false)
	render.DrawScale = engine.GetCircleScaleFactor(r, render.AnimPlayer.CurrentFrame)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.AnimPlayer.Paused = true

	b := spawnBody(m, e, f, r, entry)
	b.SetPosition(pos)
	b.FirstShape().SetCollisionType(types.CollSnowball)
	b.FirstShape().Filter = cm.NewShapeFilter(0, types.BitmaskSnowball, cm.AllCategories&^types.BitmaskPlayer)
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
		comp.WallTag,
		comp.Render,
	))

	wallShape.Body().UserData = entry
	comp.Body.Set(entry, wallShape.Body())

	render := comp.Render.Get(entry)

	render.AnimPlayer = engine.NewAnimationPlayer(res.Wall)
	imW := res.Wall.Bounds().Dx()
	render.AnimPlayer.AddStateAnimation("idle", 0, 0, imW, imW, 1, false)
	render.DrawScale = engine.GetBoxScaleFactor(float64(imW), float64(imW), boxW, boxH)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.AnimPlayer.Paused = true
	render.ScaleColor = colornames.Blue
	return entry
}

func SpawnDoor(boxCenter cm.Vec2, boxW, boxH float64, lockNumber int) *donburi.Entry {
	sbody := cm.NewStaticBody()
	shape := cm.NewBox(sbody, boxW, boxH, 0)
	shape.Filter = cm.NewShapeFilter(0, types.BitmaskDoor, cm.AllCategories)
	shape.SetSensor(false)
	shape.SetElasticity(0)
	shape.SetFriction(0)
	shape.CollisionType = types.CollDoor
	sbody.SetPosition(boxCenter)
	res.Space.AddShape(shape)
	res.Space.AddBody(shape.Body())

	// components
	entry := res.World.Entry(res.World.Create(
		comp.Body,
		comp.Door,
		comp.Render,
	))

	shape.Body().UserData = entry
	comp.Body.Set(entry, shape.Body())

	render := comp.Render.Get(entry)

	render.AnimPlayer = engine.NewAnimationPlayer(res.Wall)
	imW := res.Wall.Bounds().Dx()
	render.AnimPlayer.AddStateAnimation("idle", 0, 0, imW, imW, 1, false)
	render.DrawScale = engine.GetBoxScaleFactor(float64(imW), float64(imW), boxW, boxH)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.AnimPlayer.Paused = true
	return entry
}

// func SpawnCollectible(itemType model.ItemType, count, keyNumber int,
// 	r float64, pos cm.Vec2) *donburi.Entry {

// 	entry := spawnBody(1, 0, 1, r, pos)
// 	body := comp.Body.Get(entry)
// 	// body.FirstShape().Filter = cm.NewShapeFilter(0, res.BitmaskCollectible, res.BitmaskPlayer|res.BitmaskWall|res.BitmaskCollectible)
// 	body.FirstShape().Filter = cm.NewShapeFilter(0, res.BitmaskCollectible, cm.AllCategories&^res.BitmaskSnowball)
// 	body.FirstShape().SetCollisionType(res.CollCollectible)

// 	entry.AddComponent(comp.Collectible)
// 	entry.AddComponent(comp.Render)

// 	comp.Collectible.SetValue(entry, model.CollectibleData{
// 		Type:      itemType,
// 		ItemCount: count,
// 		KeyNumber: keyNumber})

// 	var ap *engine.AnimationPlayer

// 	switch itemType {

// 	case res.ItemBomb:
// 		ap = engine.NewAnimationPlayer(res.Items)
// 		ap.AddStateAnimation("idle", 0, 0, 100, 100, 1, false)
// 	case res.ItemKey:
// 		ap = engine.NewAnimationPlayer(res.Items)
// 		ap.AddStateAnimation("idle", 100, 0, 100, 100, 1, false)
// 	case res.ItemSnowball:
// 		ap = engine.NewAnimationPlayer(res.Items)
// 		ap.AddStateAnimation("idle", 200, 0, 100, 100, 1, false)
// 	case res.ItemPotion:
// 		ap = engine.NewAnimationPlayer(res.Items)
// 		ap.AddStateAnimation("idle", 300, 0, 100, 100, 1, false)
// 	default:
// 		ap = engine.NewAnimationPlayer(res.Items)
// 		ap.AddStateAnimation("idle", 200, 0, 100, 100, 1, false)

// 	}

// 	render := comp.Render.Get(entry)
// 	render.AnimPlayer = ap

// 	render.AnimPlayer.Paused = true
// 	render.DrawScale = engine.GetCircleScaleFactor(r, render.AnimPlayer.CurrentFrame)
// 	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
// 	render.ScaleColor = colornames.Yellow

// 	return entry
// }

func SpawnRoom(roomBB cm.BB, opts RoomOptions) {

	topDoorLength := roomBB.Width() / 5
	leftDoorLength := roomBB.Height() / 5

	topDoorCenter := roomBB.LT().Lerp(roomBB.RT(), 0.5)
	bottomDoorCenter := roomBB.LB().Lerp(roomBB.RB(), 0.5)

	leftDoorCenter := roomBB.LT().Lerp(roomBB.LB(), 0.5)
	rightDoorCenter := roomBB.RT().Lerp(roomBB.RB(), 0.5)

	topLeftWallCenter := cm.Vec2{roomBB.L + topDoorLength, roomBB.T}
	topRightWallCenter := cm.Vec2{roomBB.R - topDoorLength, roomBB.T}
	bottomLeftWallCenter := cm.Vec2{roomBB.L + topDoorLength, roomBB.B}
	bottomRightWallCenter := cm.Vec2{roomBB.R - topDoorLength, roomBB.B}

	leftDoorBottom := cm.Vec2{roomBB.L, roomBB.B + leftDoorLength}
	leftDoorTop := cm.Vec2{roomBB.L, roomBB.T - leftDoorLength}

	rightDoorBottom := cm.Vec2{roomBB.R, roomBB.B + leftDoorLength}
	rightDoorTop := cm.Vec2{roomBB.R, roomBB.T - leftDoorLength}

	// Top Wall
	if opts.TopWall {
		SpawnWall(topLeftWallCenter, topDoorLength*2, 10)
		SpawnDoor(topDoorCenter, topDoorLength, 10, opts.TopDoorKeyNumber)
		SpawnWall(topRightWallCenter, topDoorLength*2, 10)
	}

	// Bottom Wall
	if opts.BottomWall {
		SpawnWall(bottomLeftWallCenter, topDoorLength*2, 10)
		SpawnDoor(bottomDoorCenter, topDoorLength, 10, opts.BottomDoorKeyNumber)
		SpawnWall(bottomRightWallCenter, topDoorLength*2, 10)
	}

	// Left Wall
	if opts.LeftWall {
		SpawnWall(leftDoorTop, 10, leftDoorLength*2)
		SpawnDoor(leftDoorCenter, 10, leftDoorLength, opts.LeftDoorKeyNumber)
		SpawnWall(leftDoorBottom, 10, leftDoorLength*2)
	}

	// Right Wall
	if opts.RightWall {
		SpawnWall(rightDoorTop, 10, leftDoorLength*2)
		SpawnDoor(rightDoorCenter, 10, leftDoorLength, opts.RightDoorKeyNumber)
		SpawnWall(rightDoorBottom, 10, leftDoorLength*2)
	}
}

type RoomOptions struct {
	TopWall, BottomWall, LeftWall, RightWall                                     bool
	TopDoorKeyNumber, BottomDoorKeyNumber, LeftDoorKeyNumber, RightDoorKeyNumber int
}
