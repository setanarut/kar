package system

import (
	"kar"
	"kar/comp"
	"kar/items"
	"kar/res"
	"kar/types"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/cm"

	"github.com/yohamta/donburi"
)

// index kontrol fonksiyonu, sadece bool değeri döndürüyor
func CheckIndex[T any](slice []T, index int) bool {
	if index < 0 || index >= len(slice) {
		return false
	}
	return true
}

func getSprite(id uint16) *ebiten.Image {
	im, ok := res.Images[id]
	if ok {
		return im
	} else {
		if len(res.Frames[id]) > 0 {
			return res.Frames[id][0]
		} else {
			return res.Images[items.Air]
		}
	}
}

func getDisplayName(e *donburi.Entry) string {
	return items.Property[comp.Item.Get(e).ID].DisplayName
}

func getArbiterDisplayNames(arb *cm.Arbiter) (string, string) {
	var aname, bname string
	if checkEntries(arb) {
		a, b := getEntries(arb)
		if a.HasComponent(comp.Item) {
			aname = items.Property[comp.Item.Get(a).ID].DisplayName
		} else if a.HasComponent(comp.TagPlayer) {
			aname = "Player"
		}
		if b.HasComponent(comp.Item) {
			bname = items.Property[comp.Item.Get(b).ID].DisplayName
		} else if a.HasComponent(comp.TagPlayer) {
			aname = "Player"
		}
	}
	return aname, bname
}

// destroy body with entry
func destroyBody(b *cm.Body) {
	if space.ContainsBody(b) {
		e := b.UserData.(*donburi.Entry)
		e.Remove()
		space.AddPostStepCallback(removeBodyPostStep, b, false)
	}
}

// destroy entry with body
func destroyEntry(entry *donburi.Entry) {
	if entry.Valid() {
		if entry.HasComponent(comp.Body) {
			body := comp.Body.Get(entry)
			destroyBody(body)
		}
	}
}

func removeBodyPostStep(space *cm.Space, body, data interface{}) {
	space.RemoveBodyWithShapes(body.(*cm.Body))
}

func resourcesetHealthComponent(e *donburi.Entry) {
	h := comp.Health.Get(e)
	h.Health = h.MaxHealth
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

func timerIsReady(t *types.Timer) bool {
	return t.Elapsed > t.Duration
}

func timerUpdate(timer *types.Timer) {
	if timer.Elapsed < timer.Duration {
		timer.Elapsed += kar.TimerTick
	}
}

// func timerRemaining(t *types.Timer) time.Duration {
// 	return t.Duration - t.Elapsed
// }

// func timerRemainingSecondsString(t *types.Timer) string {
// 	return fmt.Sprintf("%.1fs", timerRemaining(t).Abs().Seconds())
// }

// func timerReset(t *types.Timer) {
// 	t.Elapsed = 0
// }

// func timerIsStart(t *types.Timer) bool {
// 	return t.Elapsed == 0
// }

// func teleportBody(p vec.Vec2, entry *donburi.Entry) {
// 	if entry.Valid() {
// 		if entry.HasComponent(comp.Body) {
// 			body := comp.Body.Get(entry)
// 			body.SetVelocity(0, 0)
// 			body.SetPosition(p)
// 		}
// 	}
// }

// func explode(bomb *donburi.Entry) {
// 	bombBody := comp.Body.Get(bomb)
// 	space := bombBody.FirstShape().Space()
// 	comp.EnemyTag.Each(bomb.World, func(enemy *donburi.Entry) {
// 		enemyHealth := comp.Health.Get(enemy)
// 		enemyBody := comp.Body.Get(enemy)
// 		queryInfo := space.SegmentQueryFirst(bombBody.Position(), enemyBody.Position(), 0, resources.FilterBombRaycast)
// 		contactShape := queryInfo.Shape
// 		if contactShape != nil {
// 			if contactShape.Body() == enemyBody {
// 				applyRaycastImpulse(queryInfo, 1000)

// 				enemyHealth.Health -= mathutil.MapRange(queryInfo.Alpha, 0.5, 1, 200, 0)
// 				if enemyHealth.Health < 0 {
// 					destroyEntryWithBody(enemy)
// 				}

// 			}
// 		}

// 	})
// 	resources.Cam.AddTrauma(0.2)
// 	destroyEntryWithBody(bomb)
// }

// func applyRaycastImpulse(sqi cm.SegmentQueryInfo, power float64) {
// 	impulseVec2 := sqi.Normal.Neg().Scale(power * mathutil.MapRange(sqi.Alpha, 0.5, 1, 1, 0))
// 	sqi.Shape.Body().ApplyImpulseAtWorldPoint(impulseVec2, sqi.Point)
// }
