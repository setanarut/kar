package system

import (
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/res"

	"github.com/yohamta/donburi"
)

func destroyBodyWithEntry(b *cm.Body) {
	s := b.FirstShape().Space()
	if s.ContainsBody(b) {
		e := b.UserData.(*donburi.Entry)
		e.Remove()
		s.AddPostStepCallback(removeBodyPostStep, b, false)
	}
}
func destroyEntryWithBody(entry *donburi.Entry) {
	if entry.Valid() {
		if entry.HasComponent(comp.Body) {
			body := comp.Body.Get(entry)
			destroyBodyWithEntry(body)
		}
	}
}

func explode(bomb *donburi.Entry) {
	bombBody := comp.Body.Get(bomb)
	space := bombBody.FirstShape().Space()
	comp.EnemyTag.Each(bomb.World, func(enemy *donburi.Entry) {
		enemyHealth := comp.Health.GetValue(enemy)
		enemyBody := comp.Body.Get(enemy)
		queryInfo := space.SegmentQueryFirst(bombBody.Position(), enemyBody.Position(), 0, res.FilterBombRaycast)
		contactShape := queryInfo.Shape
		if contactShape != nil {
			if contactShape.Body() == enemyBody {
				applyRaycastImpulse(queryInfo, 1000)
				comp.Health.SetValue(enemy, enemyHealth-engine.MapRange(queryInfo.Alpha, 0.5, 1, 200, 0))
				if enemyHealth < 0 {
					destroyEntryWithBody(enemy)
				}

			}
		}

	})
	res.Camera.AddTrauma(0.2)
	destroyEntryWithBody(bomb)
}

func applyRaycastImpulse(sqi cm.SegmentQueryInfo, power float64) {
	impulseVec2 := sqi.Normal.Neg().Mult(power * engine.MapRange(sqi.Alpha, 0.5, 1, 1, 0))
	sqi.Shape.Body().ApplyImpulseAtWorldPoint(impulseVec2, sqi.Point)
}

func removeBodyPostStep(space *cm.Space, body, data interface{}) {
	space.RemoveBodyWithShapes(body.(*cm.Body))
}

func destroyStopped(e *donburi.Entry) {
	if e.HasComponent(comp.Body) {
		b := comp.Body.Get(e)
		if engine.IsMoving(b.Velocity(), 80) {
			destroyBodyWithEntry(b)
		}
	}
}
func destroyDead(e *donburi.Entry) {
	if e.HasComponent(comp.Health) {
		if comp.Health.GetValue(e) < 1 {
			destroyEntryWithBody(e)
		}
	}
}

func getEntry(b *cm.Body) *donburi.Entry {
	return b.UserData.(*donburi.Entry)
}
func checkEntry(b *cm.Body) bool {
	e, ok := b.UserData.(*donburi.Entry)
	return ok && e.Valid()
}

func getEntries(arb *cm.Arbiter) (*donburi.Entry, *donburi.Entry) {
	a, b := arb.Bodies()
	return a.UserData.(*donburi.Entry), b.UserData.(*donburi.Entry)

}

func checkEntries(arb *cm.Arbiter) bool {
	aBody, bBody := arb.Bodies()
	return checkEntry(aBody) && checkEntry(bBody)
}
