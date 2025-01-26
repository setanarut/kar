package system

import (
	"kar"
	"kar/arc"
	"kar/engine/mathutil"
	"math"
)

var (
	itemGravity float64 = 3
	sinspace            = mathutil.SinSpace(0, 2*math.Pi, 3, 60)
	sinspaceLen         = len(sinspace) - 1
)

type Item struct{}

func (itm *Item) Init() {}
func (itm *Item) Update() {

	if kar.WorldECS.Alive(player) {
		if !craftingState {
			itemQuery := arc.FilterItem.Query(&kar.WorldECS)
			for itemQuery.Next() {
				itemID, rect, timers, durability := itemQuery.Get()

				if timers.CollisionCountdown != 0 {
					timers.CollisionCountdown--
				} else {
					if ctrl.Rect.Overlaps(rect) {
						// Çarpan öğeyi envantere ekle
						if kar.GopherInventory.AddItemIfEmpty(itemID.ID, durability.Durability) {
							toRemove = append(toRemove, itemQuery.Entity())
							onInventorySlotChanged()
						}
					}
				}
				dy := collider.CollideY(rect.X, rect.Y+6, rect.W, rect.H, itemGravity)
				rect.Y += dy
				timers.AnimationIndex = (timers.AnimationIndex + 1) % sinspaceLen

			}
		}
	}
}
func (itm *Item) Draw() {

}
