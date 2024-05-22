package system

import (
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/res"

	"github.com/yohamta/donburi"
)

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

func DestroyStopped(e *donburi.Entry) {
	b := comp.Body.Get(e)
	if engine.IsMoving(b.Velocity(), 80) {
		DestroyBodyWithEntry(b)
	}
}
func DestroyDead(e *donburi.Entry) {
	if comp.Health.GetValue(e) < 1 {
		DestroyEntryWithBody(e)
	}

}
