package system

import (
	"kar/comp"
	"kar/items"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/cm"

	"github.com/yohamta/donburi"
)

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
func DisplayName(id uint16) string {
	return items.Property[id].DisplayName
}

// destroy body with entry
func destroyBody(b *cm.Body) {
	if Space.ContainsBody(b) {
		e := b.UserData.(*donburi.Entry)
		e.Remove()
		Space.AddPostStepCallback(removeBodyPostStep, b, false)
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

func resetHealthComponent(e *donburi.Entry) {
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
