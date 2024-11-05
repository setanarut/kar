package system

import (
	"kar/arc"
	"kar/items"
)

// Add item to inventory if empty
func addItem(inv *arc.Inventory, id uint16) bool {
	i, ok1 := hasItemStackSpace(inv, id)
	if ok1 {
		inv.Slots[i].Quantity++
		return true
	} else {
		i2, ok2 := hasEmptySlot(inv)
		if ok2 {
			inv.Slots[i2].Quantity++
			inv.Slots[i2].ID = id
			return true
		}
	}
	return false
}

func RemoveHandItem(inv *arc.Inventory, id uint16) bool {
	ok := HasHandItem(inv, id)
	if ok {
		inv.HandSlot.Quantity--
		return true
	} else {
		inv.HandSlot.ID = items.Air
	}
	return false
}

func removeItem(inv *arc.Inventory, id uint16) bool {
	i, ok := hasItem(inv, id)
	if ok {
		inv.Slots[i].Quantity--
		return true
	}
	return false
}

func deleteSlot(inv *arc.Inventory, index int) {
	inv.Slots[index] = arc.ItemStack{}
}

func resetInventory(inv *arc.Inventory) {
	for i := range inv.Slots {
		inv.Slots[i].ID = items.Air
		inv.Slots[i].Quantity = 0
	}
}

func hasEmptySlot(inv *arc.Inventory) (index int, ok bool) {
	for i, v := range inv.Slots {
		if v.Quantity == 0 {
			return i, true
		}
	}
	return -1, false
}
func hasItemStackSpace(inv *arc.Inventory, id uint16) (index int, ok bool) {
	for i, v := range inv.Slots {
		if v.ID == id && v.Quantity < 64 && v.Quantity > 0 {
			return i, true
		}
	}
	return -1, false
}

func hasItem(inv *arc.Inventory, id uint16) (index int, ok bool) {
	for i, v := range inv.Slots {
		if v.ID == id && v.Quantity > 0 {
			return i, true
		}
	}
	return -1, false
}
func HasHandItem(inv *arc.Inventory, id uint16) bool {
	if inv.HandSlot.ID == id && inv.HandSlot.Quantity > 0 {
		return true
	}
	return false
}
