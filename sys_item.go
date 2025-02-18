package kar

import (
	"github.com/mlange-42/arche/ecs"
)

var SinspaceLen int = len(Sinspace) - 1

type Item struct {
	toRemoveComponent []ecs.Entity
	itemBox           AABB
	itemHit           *HitInfo
}

func (i *Item) Init() {
	i.itemBox = AABB{Half: DropItemHalfSize}
	i.itemHit = &HitInfo{}
}
func (i *Item) Update() error {

	q := FilterCollisionDelayer.Query(&ECWorld)

	for q.Next() {
		delayer := q.Get()
		delayer.Duration -= Tick
	}

	// dropped items collisions and animations
	if ECWorld.Alive(CurrentPlayer) {
		if !GameDataRes.CraftingState {
			// playerPos, playerSize := MapAABB.GetUnchecked(CurrentPlayer)
			playerBox := MapAABB.GetUnchecked(CurrentPlayer)

			itemQuery := FilterDroppedItem.Query(&ECWorld)
			for itemQuery.Next() {

				itemID, itemPos, timers, delayer, durability := itemQuery.Get()
				i.itemBox.Pos = Vec(*itemPos)
				itemEntity := itemQuery.Entity()
				// Check player-item collision
				if delayer == nil {
					if playerBox.Overlap(i.itemBox, i.itemHit) {
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
				i.itemBox.Pos.Y += 10
				dy := TileCollider.CollideY(i.itemBox, ItemGravity)
				itemPos.Y += dy
				timers.Index = (timers.Index + 1) % SinspaceLen

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
