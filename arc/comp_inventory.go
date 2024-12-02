package arc

import (
	"kar/items"
)

type ItemStack struct {
	ID       uint16
	Quantity uint8
}

type Inventory struct {
	Slots    [9]ItemStack
	HandSlot ItemStack
}

func NewInventory() *Inventory {
	inv := &Inventory{}
	inv.HandSlot = ItemStack{}
	for i := range inv.Slots {
		inv.Slots[i] = ItemStack{}
	}
	return inv
}

// Add item to inventory if empty
func (i *Inventory) AddItem(id uint16) bool {
	is, ok1 := i.HasItemStackSpace(id)
	if ok1 {
		i.Slots[is].Quantity++
		return true
	} else {
		i2, ok2 := i.HasEmptySlot()
		if ok2 {
			i.Slots[i2].Quantity++
			i.Slots[i2].ID = id
			return true
		}
	}
	return false
}

func (i *Inventory) RemoveHandItem(id uint16) bool {
	ok := i.HasHandItem(id)
	if ok {
		i.HandSlot.Quantity--
		return true
	} else {
		i.HandSlot.ID = items.Air
	}
	return false
}

func (i *Inventory) RemoveItem(id uint16) bool {
	is, ok := i.HasItem(id)
	if ok {
		i.Slots[is].Quantity--
		return true
	}
	return false
}

func (i *Inventory) DeleteSlot(index int) {
	i.Slots[index] = ItemStack{}
}

func (i *Inventory) Reset() {
	for si := range i.Slots {
		i.Slots[si].ID = items.Air
		i.Slots[si].Quantity = 0
	}
}

func (i *Inventory) HasEmptySlot() (index int, ok bool) {
	for i, v := range i.Slots {
		if v.Quantity == 0 {
			return i, true
		}
	}
	return -1, false
}
func (i *Inventory) HasItemStackSpace(id uint16) (index int, ok bool) {
	for i, v := range i.Slots {
		if v.ID == id && v.Quantity < 64 && v.Quantity > 0 {
			return i, true
		}
	}
	return -1, false
}

func (i *Inventory) HasItem(id uint16) (index int, ok bool) {
	for is, v := range i.Slots {
		if v.ID == id && v.Quantity > 0 {
			return is, true
		}
	}
	return -1, false
}
func (i *Inventory) HasHandItem(id uint16) bool {
	if i.HandSlot.ID == id && i.HandSlot.Quantity > 0 {
		return true
	}
	return false
}
