package kar

import (
	"github.com/mlange-42/ark/ecs"
)

type Item struct {
	toRemoveComponent []ecs.Entity
	itemBox           *AABB
	itemHit           *HitInfo
}

func (i *Item) Init() {
	i.itemBox = &AABB{Half: dropItemHalfSize}
	i.itemHit = &HitInfo{}
}
func (i *Item) Update() {

	q := filterCollisionDelayer.Query()

	for q.Next() {
		delayer := q.Get()
		delayer.Duration -= Tick
	}

	// dropped items collisions and animations
	if world.Alive(currentPlayer) {
		if gameDataRes.GameplayState == Playing {
			playerBox := mapAABB.GetUnchecked(currentPlayer)
			itemQuery := filterDroppedItem.Query()
			for itemQuery.Next() {
				itemID, itemPos, timers := itemQuery.Get()
				i.itemBox.Pos = Vec(*itemPos)
				itemEntity := itemQuery.Entity()

				if !mapCollisionDelayer.HasUnchecked(itemEntity) {
					// Check player-item collision
					if playerBox.Overlap(i.itemBox, i.itemHit) {
						// if Durability component exists, get durability
						dur := 0
						if mapDurability.HasUnchecked(itemEntity) {
							dur = mapDurability.GetUnchecked(itemEntity).Durability
						}
						if inventoryRes.AddItemIfEmpty(itemID.ID, dur) {
							toRemove = append(toRemove, itemEntity)
						}
						onInventorySlotChanged()
					}
				} else {
					if mapCollisionDelayer.GetUnchecked(itemEntity).Duration < 0 {
						i.toRemoveComponent = append(i.toRemoveComponent, itemEntity)
					}
				}
				i.itemBox.Pos.Y += 6
				// vertical item sine animation
				TileCollider.Collisions = TileCollider.Collisions[:0]
				dy := TileCollider.CollideY(*i.itemBox, ItemGravity)
				itemPos.Y += dy
				timers.Index = (timers.Index + 1) % len(Sinspace)
			}
		}
	}

	// Remove MapCollisionDelayer components
	for _, entity := range i.toRemoveComponent {
		mapCollisionDelayer.Remove(entity)
	}
	i.toRemoveComponent = i.toRemoveComponent[:0]
}
func (i *Item) Draw() {

	// itemQuery := FilterDroppedItem.Query(&ECWorld)
	// for itemQuery.Next() {

	// 	itemID, itemPos, timers, delayer, durability := itemQuery.Get()

	// }
}
