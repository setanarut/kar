package kar

import (
	"github.com/mlange-42/arche/ecs"
)

type Item struct {
	sinspaceLenght    int
	toRemoveComponent []ecs.Entity
}

func (i *Item) Init() {
	i.sinspaceLenght = len(Sinspace) - 1

}
func (i *Item) Update() error {

	q := FilterCollisionDelayer.Query(&ECWorld)

	for q.Next() {
		delayer := q.Get()
		delayer.Duration -= Tick
	}

	// dropped items collisions and animations
	if ECWorld.Alive(CurrentPlayer) {
		itemSize := &Size{DropItemSize.W, DropItemSize.H}
		if !GameDataRes.CraftingState {
			playerPos, playerSize := MapRect.GetUnchecked(CurrentPlayer)

			itemQuery := FilterDroppedItem.Query(&ECWorld)
			for itemQuery.Next() {

				itemID, itemPos, timers, delayer, durability := itemQuery.Get()
				itemEntity := itemQuery.Entity()
				// Check player-item collision
				if delayer == nil {
					if Overlaps(playerPos, playerSize, itemPos, itemSize) {
						// if Durability component exists, pass durability
						if durability != nil {
							if InventoryRes.AddItemIfEmpty(itemID.ID, durability.Durability) {
								toRemove = append(toRemove, itemEntity)
							}
						} else {
							if InventoryRes.AddItemIfEmpty(itemID.ID, 0) {
								toRemove = append(toRemove, itemEntity)
							}

						}
						onInventorySlotChanged()
					}
				} else {
					// delayer.Duration -= Tick
					if delayer.Duration < 0 {
						i.toRemoveComponent = append(i.toRemoveComponent, itemEntity)
					}
				}

				// vertical item sine animation
				dy := Collider.CollideY(itemPos.X, itemPos.Y+6, itemSize.W, itemSize.H, ItemGravity)
				itemPos.Y += dy
				timers.Index = (timers.Index + 1) % i.sinspaceLenght

			}
		}
	}

	// Remove MapCollisionDelayer components
	for _, entity := range i.toRemoveComponent {
		MapCollisionDelayer.Remove(entity)
	}
	i.toRemoveComponent = i.toRemoveComponent[:0]
	return nil
}
func (i *Item) Draw() {

}
