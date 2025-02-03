package system

import (
	"kar"
	"kar/arc"

	"github.com/mlange-42/arche/ecs"
)

type Item struct {
	sinspaceLenght    int
	toRemoveComponent []ecs.Entity
}

func (itm *Item) Init() {
	itm.sinspaceLenght = len(kar.Sinspace) - 1

}
func (itm *Item) Update() {

	q := arc.FilterCollisionDelayer.Query(&kar.ECWorld)

	for q.Next() {
		delayer := q.Get()
		delayer.Duration -= kar.Tick
	}

	// dropped items collisions and animations
	if kar.ECWorld.Alive(kar.CurrentPlayer) {
		itemSize := &arc.Size{kar.GameDataRes.DropItemW, kar.GameDataRes.DropItemH}
		if !kar.GameDataRes.CraftingState {
			playerPos, playerSize := arc.MapRect.GetUnchecked(kar.CurrentPlayer)

			itemQuery := arc.FilterDroppedItem.Query(&kar.ECWorld)
			for itemQuery.Next() {

				itemID, itemPos, timers, delayer, durability := itemQuery.Get()
				itemEntity := itemQuery.Entity()
				// Check player-item collision
				if delayer == nil {
					if Overlaps(playerPos, playerSize, itemPos, itemSize) {
						// if Durability component exists, pass durability
						if durability != nil {
							if kar.InventoryRes.AddItemIfEmpty(itemID.ID, durability.Durability) {
								toRemove = append(toRemove, itemEntity)
							}
						} else {
							if kar.InventoryRes.AddItemIfEmpty(itemID.ID, 0) {
								toRemove = append(toRemove, itemEntity)
							}

						}
						onInventorySlotChanged()
					}
				} else {
					// delayer.Duration -= kar.Tick
					if delayer.Duration < 0 {
						itm.toRemoveComponent = append(itm.toRemoveComponent, itemEntity)
					}
				}

				// vertical item sine animation
				dy := kar.Collider.CollideY(itemPos.X, itemPos.Y+6, itemSize.W, itemSize.H, kar.ItemGravity)
				itemPos.Y += dy
				timers.AnimationIndex = (timers.AnimationIndex + 1) % itm.sinspaceLenght

			}
		}
	}

	// Remove MapCollisionDelayer components
	for _, entity := range itm.toRemoveComponent {
		arc.MapCollisionDelayer.Remove(entity)
	}
	itm.toRemoveComponent = itm.toRemoveComponent[:0]

}
func (itm *Item) Draw() {

}
