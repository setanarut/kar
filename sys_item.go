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
		*delayer -= CollisionDelayer(Tick)
	}

	// dropped items collisions and animations
	if world.Alive(currentPlayer) {
		if gameDataRes.GameplayState == Playing {
			playerBox := mapAABB.GetUnchecked(currentPlayer)
			itemQuery := filterDroppedItem.Query()
			for itemQuery.Next() {
				itemID, itemPos, animIndex := itemQuery.Get()
				dropItemAABB.Pos = *(*Vec)(itemPos)
				itemEntity := itemQuery.Entity()

				if !mapCollisionDelayer.HasUnchecked(itemEntity) {
					// Check player-item collision
					if Overlap(playerBox, dropItemAABB, i.itemHit) {
						// if Durability component exists, get durability
						dur := 0
						if mapDurability.HasUnchecked(itemEntity) {
							dur = int(*mapDurability.GetUnchecked(itemEntity))
						}
						if inventoryRes.AddItemIfEmpty(uint8(*itemID), dur) {
							toRemove = append(toRemove, itemEntity)
						}
						onInventorySlotChanged()
					}
				} else {
					if *mapCollisionDelayer.GetUnchecked(itemEntity) < 0 {
						i.toRemoveComponent = append(i.toRemoveComponent, itemEntity)
					}
				}
				dropItemAABB.Pos.Y += 6
				// vertical item sine animation
				tileCollider.Collisions = tileCollider.Collisions[:0]
				dy := tileCollider.CollideY(dropItemAABB, ItemGravity)
				itemPos.Y += dy
				*animIndex = AnimationIndex((int(*animIndex) + 1) % len(sinspace))

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
		itemid, pos, animIndex := itemQuery.Get()
		id := uint8(*itemid)
		dropItemAABB.Pos = *(*Vec)(pos)
		colorMDIO.GeoM.Reset()
		colorMDIO.GeoM.Translate(dropItemAABB.Left(), dropItemAABB.Top()+sinspace[*animIndex])
		if id != items.Air {
			cameraRes.DrawWithColorM(res.Icon8[id], colorM, colorMDIO, Screen)
		}
	}
}
