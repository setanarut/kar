package system

import (
	"kar/comp"
	"kar/constants"
	"kar/engine"
	"kar/engine/cm"
	"kar/model"
	"kar/res"
	"math"

	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

type CollisionSystem struct {
	DT float64
}

func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{
		DT: 1.0 / 60.0,
	}
}

func (ps *CollisionSystem) Init() {

	res.Space.UseSpatialHash(50, 1000)
	res.Space.CollisionBias = math.Pow(0.3, 60)
	res.Space.CollisionSlop = 0.5
	res.Space.Damping = 0.03
	// res.Space.Iterations = 1

	// Player
	// res.Space.NewCollisionHandler(constants.CollPlayer, constants.CollDoor).BeginFunc = playerDoorEnter
	res.Space.NewCollisionHandler(constants.CollPlayer, constants.CollDoor).SeparateFunc = playerDoorExit
	// res.Space.NewCollisionHandler(constants.CollPlayer, constants.CollCollectible).BeginFunc = playerCollectibleCollisionBegin

	// Enemy
	res.Space.NewCollisionHandler(constants.CollEnemy, constants.CollPlayer).PostSolveFunc = enemyPlayerPostSolve
	res.Space.NewCollisionHandler(constants.CollEnemy, constants.CollPlayer).BeginFunc = enemyPlayerBegin
	res.Space.NewCollisionHandler(constants.CollEnemy, constants.CollPlayer).SeparateFunc = enemyPlayerSep

	// Snowball
	res.Space.NewCollisionHandler(constants.CollSnowball, constants.CollEnemy).BeginFunc = snowballEnemyCollisionBegin
	res.Space.NewCollisionHandler(constants.CollSnowball, constants.CollBomb).BeginFunc = SnowballBombCollisionBegin
	res.Space.NewCollisionHandler(constants.CollSnowball, constants.CollWall).BeginFunc = snowballWallCollisionBegin
	res.Space.NewCollisionHandler(constants.CollSnowball, constants.CollDoor).BeginFunc = snowballDoorCollisionBegin

	res.Space.Step(ps.DT)

}

func (ps *CollisionSystem) Update() {

	comp.SnowballTag.Each(res.World, func(e *donburi.Entry) {
		b := comp.Body.Get(e)
		if engine.IsMoving(b.Velocity(), 80) {
			DestroyBodyWithEntry(b)
		}
	})

	if pla, ok := comp.PlayerTag.First(res.World); ok {
		playerBody := comp.Body.Get(pla)

		comp.EnemyTag.Each(res.World, func(e *donburi.Entry) {

			ene := comp.Body.Get(e)
			ai := *comp.AI.Get(e)
			mobile := comp.Mobile.Get(e)

			if ai.Follow {
				dist := playerBody.Position().Distance(ene.Position())
				if dist < ai.FollowDistance {
					speed := ene.Mass() * (mobile.Speed * 4)
					a := playerBody.Position().Sub(ene.Position()).Normalize().Mult(speed)
					ene.ApplyForceAtLocalPoint(a, ene.CenterOfGravity())
				}
			}

		})

	}
	res.Space.Step(ps.DT)
}

func (ps *CollisionSystem) Draw() {}

// Player <-> Door (exit)
func playerDoorExit(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
	playerBody, doorBody := arb.Bodies()
	doorEntry := doorBody.UserData.(*donburi.Entry)
	d := comp.Door.Get(doorEntry)
	d.Open = false
	doorBody.FirstShape().SetSensor(false)

	for _, room := range res.Rooms {
		if room.ContainsVect(playerBody.Position()) {
			res.CurrentRoom = room
		}
	}

}

// Enemy <-> Player postsolve
func enemyPlayerPostSolve(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
	enemyBody, playerBody := arb.Bodies()
	enemyEntry, eok := enemyBody.UserData.(*donburi.Entry)
	playerEntry, pok := playerBody.UserData.(*donburi.Entry)
	var charData *model.Mobile

	if eok && pok {

		if playerEntry.Valid() && enemyEntry.Valid() {
			if playerEntry.HasComponent(comp.Char) && enemyEntry.HasComponent(comp.Damage) && playerEntry.HasComponent(comp.Render) {
				charData = comp.Char.Get(playerEntry)
				comp.Render.Get(playerEntry).ScaleColor = colornames.Red
				charData.Health -= *comp.Damage.Get(enemyEntry) / 60.0
				if charData.Health < 0 {
					DestroyBodyWithEntry(playerBody)
				}
			}
		}

	}

}

// Enemy <-> Player Begin
func enemyPlayerBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	enemyBody, playerBody := arb.Bodies()
	enemyEntry, eok := enemyBody.UserData.(*donburi.Entry)
	playerEntry, pok := playerBody.UserData.(*donburi.Entry)
	var charData *model.Mobile

	if eok && pok {

		if playerEntry.Valid() && enemyEntry.Valid() {
			if playerEntry.HasComponent(comp.Health) && enemyEntry.HasComponent(comp.Damage) && playerEntry.HasComponent(comp.Render) {
				pHealth := comp.Health.Get(playerEntry)
				comp.Render.Get(playerEntry).ScaleColor = colornames.Red
				playerBody.ApplyImpulseAtLocalPoint(arb.Normal().Mult(1000), playerBody.CenterOfGravity())
				*pHealth -= *comp.Damage.Get(enemyEntry)
				if charData.Health < 0 {
					DestroyBodyWithEntry(playerBody)
				}
			}
		}

	}
	return true
}

// Enemy <-> Player Sep
func enemyPlayerSep(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
	_, playerBody := arb.Bodies()
	playerEntry := playerBody.UserData.(*donburi.Entry)
	if playerEntry.Valid() {
		comp.Render.Get(playerEntry).ScaleColor = colornames.White
	}

}

// Snowball <-> Enemy
func snowballEnemyCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	bulletBody, enemyBody := arb.Bodies()
	bullet := bulletBody.UserData.(*donburi.Entry)
	enemy := enemyBody.UserData.(*donburi.Entry)

	if enemy.Valid() {

		if enemy.HasComponent(comp.Health) && bullet.HasComponent(comp.Damage) {
			enemyHealth := comp.Health.Get(enemy)
			*enemyHealth -= comp.Damage.GetValue(bullet)

			if *enemyHealth < 0 {
				DestroyBodyWithEntry(enemyBody)
			}
		}
	}
	// çarpan bulletı yok et
	DestroyEntryWithBody(bullet)
	return true

}

// Snowball <-> Wall
func snowballWallCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	Snowball, _ := arb.Bodies()
	DestroyBodyWithEntry(Snowball)
	return false
}

// Snowball <-> Bomb
func SnowballBombCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	snowball, bomb := arb.Bodies()
	DestroyBodyWithEntry(snowball)
	Explode(bomb.UserData.(*donburi.Entry))
	return false
}

// Snowball <-> Door
func snowballDoorCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	bodyA, _ := arb.Bodies()
	bulletEntry := bodyA.UserData.(*donburi.Entry)
	DestroyEntryWithBody(bulletEntry)
	return true
}

func DestroyBodyWithEntry(b *cm.Body) {
	s := b.FirstShape().Space()
	if s.ContainsBody(b) {
		e := b.UserData.(*donburi.Entry)
		e.Remove()
		s.AddPostStepCallback(removeBodyPostStep, b, false)
	}
}
func DestroyEntryWithBody(entry *donburi.Entry) {
	if entry.Valid() {
		if entry.HasComponent(comp.Body) {
			body := comp.Body.Get(entry)
			DestroyBodyWithEntry(body)
		}
	}
}

func Explode(bomb *donburi.Entry) {
	bombBody := comp.Body.Get(bomb)
	space := bombBody.FirstShape().Space()
	comp.EnemyTag.Each(bomb.World, func(enemy *donburi.Entry) {
		enemyHealth := comp.Health.GetValue(enemy)
		enemyBody := comp.Body.Get(enemy)
		queryInfo := space.SegmentQueryFirst(bombBody.Position(), enemyBody.Position(), 0, res.FilterBombRaycast)
		contactShape := queryInfo.Shape
		if contactShape != nil {
			if contactShape.Body() == enemyBody {
				ApplyRaycastImpulse(queryInfo, 1000)
				comp.Health.SetValue(enemy, enemyHealth-engine.MapRange(queryInfo.Alpha, 0.5, 1, 200, 0))
				if enemyHealth < 0 {
					DestroyEntryWithBody(enemy)
				}

			}
		}

	})
	res.Camera.AddTrauma(0.2)
	DestroyEntryWithBody(bomb)
}

func ApplyRaycastImpulse(sqi cm.SegmentQueryInfo, power float64) {
	impulseVec2 := sqi.Normal.Neg().Mult(power * engine.MapRange(sqi.Alpha, 0.5, 1, 1, 0))
	sqi.Shape.Body().ApplyImpulseAtWorldPoint(impulseVec2, sqi.Point)
}

func removeBodyPostStep(space *cm.Space, body, data interface{}) {
	space.RemoveBodyWithShapes(body.(*cm.Body))
}
