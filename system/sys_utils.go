package system

import (
	"image/color"
	"kar/comp"
	"kar/engine/cm"
	"kar/engine/mathutil"
	"kar/res"

	"github.com/yohamta/donburi"
)

func DestroyBodyWithEntry(b *cm.Body) {
	s := b.FirstShape().Space()
	if s.ContainsBody(b) {
		e := b.UserData.(*donburi.Entry)
		e.Remove()
		s.AddPostStepCallback(RemoveBodyPostStep, b, false)
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
		enemyHealth := comp.Health.Get(enemy)
		enemyBody := comp.Body.Get(enemy)
		queryInfo := space.SegmentQueryFirst(bombBody.Position(), enemyBody.Position(), 0, res.FilterBombRaycast)
		contactShape := queryInfo.Shape
		if contactShape != nil {
			if contactShape.Body() == enemyBody {
				ApplyRaycastImpulse(queryInfo, 1000)

				enemyHealth.Health -= mathutil.MapRange(queryInfo.Alpha, 0.5, 1, 200, 0)
				if enemyHealth.Health < 0 {
					DestroyEntryWithBody(enemy)
				}

			}
		}

	})
	res.Camera.AddTrauma(0.2)
	DestroyEntryWithBody(bomb)
}

func ApplyRaycastImpulse(sqi cm.SegmentQueryInfo, power float64) {
	impulseVec2 := sqi.Normal.Neg().Scale(power * mathutil.MapRange(sqi.Alpha, 0.5, 1, 1, 0))
	sqi.Shape.Body().ApplyImpulseAtWorldPoint(impulseVec2, sqi.Point)
}

func RemoveBodyPostStep(space *cm.Space, body, data interface{}) {
	space.RemoveBodyWithShapes(body.(*cm.Body))
}

func DestroyStopped(e *donburi.Entry) {
	if e.HasComponent(comp.Body) {
		b := comp.Body.Get(e)
		if mathutil.IsMoving(b.Velocity(), 80) {
			DestroyBodyWithEntry(b)
		}
	}
}
func ResetHealthComponent(e *donburi.Entry) {
	h := comp.Health.Get(e)
	h.Health = h.MaxHealth
}

func DestroyOnCollisionAndStopped(e *donburi.Entry) {
	b := comp.Body.Get(e)
	b.EachArbiter(func(a *cm.Arbiter) {
		if CheckEntries(a) {
			snow, _ := GetEntries(a)
			DestroyEntryWithBody(snow)
		}
	})
	if e.Valid() {
		DestroyStopped(e)
	}
}

func DestroyDead(e *donburi.Entry) {
	if e.HasComponent(comp.Health) {

		if comp.Health.Get(e).Health <= 0 {
			if e.HasComponent(comp.Block) {
				blockPos := comp.Body.Get(e).Position().Point().Div(50)
				res.Terrain.SetGray(blockPos.X, blockPos.Y, color.Gray{0})
				DestroyEntryWithBody(e)
			}

			DestroyEntryWithBody(e)
		}
	}
}

func GetEntry(b *cm.Body) *donburi.Entry {
	return b.UserData.(*donburi.Entry)
}
func CheckEntry(b *cm.Body) bool {
	e, ok := b.UserData.(*donburi.Entry)
	return ok && e.Valid()
}

func GetEntries(arb *cm.Arbiter) (*donburi.Entry, *donburi.Entry) {
	a, b := arb.Bodies()
	return a.UserData.(*donburi.Entry), b.UserData.(*donburi.Entry)

}

func CheckEntries(arb *cm.Arbiter) bool {
	aBody, bBody := arb.Bodies()
	return CheckEntry(aBody) && CheckEntry(bBody)
}
