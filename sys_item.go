package kar

import (
	"github.com/mlange-42/ark/ecs"
)

type Item struct {
	toRemoveComponent []ecs.Entity
	itemBox           AABB
	itemHit           *HitInfo
}

func (i *Item) Init() {
	i.itemBox = AABB{Half: dropItemHalfSize}
	i.itemHit = &HitInfo{}
}
func (i *Item) Update() {

	q := FilterCollisionDelayer.Query()

	for q.Next() {
		delayer := q.Get()
		delayer.Duration -= Tick
	}

	// dropped items collisions and animations
	if world.Alive(currentPlayer) {
		if gameDataRes.GameplayState == Playing {
			playerBox := MapAABB.GetUnchecked(currentPlayer)
			itemQuery := FilterDroppedItem.Query()
			for itemQuery.Next() {
				itemID, itemPos, timers := itemQuery.Get()
				i.itemBox.Pos = Vec(*itemPos)
				itemEntity := itemQuery.Entity()
				// Check player-item collision
				if !MapCollisionDelayer.HasUnchecked(itemEntity) {
					if playerBox.Overlap(i.itemBox, i.itemHit) {
						// if Durability component exists, pass durability
						if MapDurability.HasUnchecked(itemEntity) {
							d := MapDurability.GetUnchecked(itemEntity)
							if inventoryRes.AddItemIfEmpty(itemID.ID, d.Durability) {
								toRemove = append(toRemove, itemEntity)
							}
						} else {
							if inventoryRes.AddItemIfEmpty(itemID.ID, 0) {
								toRemove = append(toRemove, itemEntity)
							}

						}
						onInventorySlotChanged()
					}
				} else {
					if MapCollisionDelayer.GetUnchecked(itemEntity).Duration < 0 {
						i.toRemoveComponent = append(i.toRemoveComponent, itemEntity)
					}
				}
				i.itemBox.Pos.Y += 6
				// vertical item sine animation
				TileCollider.Collisions = TileCollider.Collisions[:0]
				dy := TileCollider.CollideY(i.itemBox, ItemGravity)
				itemPos.Y += dy
				timers.Index = (timers.Index + 1) % len(Sinspace)
			}
		}
	}

	// Remove MapCollisionDelayer components
	for _, entity := range i.toRemoveComponent {
		MapCollisionDelayer.Remove(entity)
	}
	i.toRemoveComponent = i.toRemoveComponent[:0]
}
func (i *Item) Draw() {

	// itemQuery := FilterDroppedItem.Query(&ECWorld)
	// for itemQuery.Next() {

	// 	itemID, itemPos, timers, delayer, durability := itemQuery.Get()

	// }
}
