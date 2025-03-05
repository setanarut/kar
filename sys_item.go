package kar

import (
	"kar/items"
	"kar/res"

	"github.com/mlange-42/ark/ecs"
)

type Item struct {
	toRemoveComponent []ecs.Entity
	itemHit           *HitInfo
}

func (i *Item) Init() {
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
				dropItemAABB.Pos = *(*Vec)(itemPos)
				itemEntity := itemQuery.Entity()

				if !mapCollisionDelayer.HasUnchecked(itemEntity) {
					// Check player-item collision
					if Overlap(playerBox, dropItemAABB, i.itemHit) {
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
				dropItemAABB.Pos.Y += 6
				// vertical item sine animation
				tileCollider.Collisions = tileCollider.Collisions[:0]
				dy := tileCollider.CollideY(dropItemAABB, ItemGravity)
				itemPos.Y += dy
				timers.Index = (timers.Index + 1) % len(sinspace)
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
	// Draw drop Items
	itemQuery := filterDroppedItem.Query()
	for itemQuery.Next() {
		id, pos, animIndex := itemQuery.Get()
		dropItemAABB.Pos = *(*Vec)(pos)
		colorMDIO.GeoM.Reset()
		colorMDIO.GeoM.Translate(dropItemAABB.Left(), dropItemAABB.Top()+sinspace[animIndex.Index])
		if id.ID != items.Air {
			cameraRes.DrawWithColorM(res.Icon8[id.ID], colorM, colorMDIO, Screen)
		}
	}
}
