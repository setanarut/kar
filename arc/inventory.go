package arc

import (
	"kar/items"
)

type ItemStack struct {
	ID         uint16
	Quantity   uint8
	Durability int
}

type Inventory struct {
	SelectedSlotIndex int
	Slots             [9]ItemStack
	HandSlot          ItemStack
}

func NewInventory() *Inventory {
	inv := &Inventory{}
	inv.HandSlot = ItemStack{}
	for i := range inv.Slots {
		inv.Slots[i] = ItemStack{}
	}
	return inv
}

// AddItemIfEmpty adds item to inventory if empty
func (i *Inventory) AddItemIfEmpty(id uint16, dura int) bool {
	idx, ok1 := i.HasItemStackSpace(id)
	if ok1 {
		i.Slots[idx].Quantity++
		return true
	} else {
		i2, ok2 := i.HasEmptySlot()
		if ok2 {
			i.Slots[i2].Quantity++
			i.Slots[i2].ID = id
			i.Slots[i2].Durability = dura
			return true
		}
	}
	return false
}

func (i *Inventory) SetSlot(slotIndex int, id uint16, quantity uint8, dur int) {
	if quantity > 0 {
		i.Slots[slotIndex] = ItemStack{
			ID:         id,
			Quantity:   quantity,
			Durability: dur,
		}
	}
}

func (i *Inventory) SelectNextSlot() {
	i.SelectedSlotIndex = (i.SelectedSlotIndex + 1) % 9
}
func (i *Inventory) SelectPrevSlot() {
	i.SelectedSlotIndex = (i.SelectedSlotIndex - 1 + 9) % 9
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
	idx, ok := i.HasItem(id)
	if ok {
		i.Slots[idx].Quantity--
		return true
	}
	return false
}
func (i *Inventory) RemoveItemFromSelectedSlot() {
	if i.Slots[i.SelectedSlotIndex].Quantity == 1 {
		i.ClearSelectedSlot()
		return
	}
	if i.Slots[i.SelectedSlotIndex].Quantity > 0 {
		i.Slots[i.SelectedSlotIndex].Quantity--
	}
}

func (i *Inventory) SelectedSlot() *ItemStack {
	return &i.Slots[i.SelectedSlotIndex]
}

func (i *Inventory) SelectedSlotID() uint16 {
	return i.Slots[i.SelectedSlotIndex].ID
}

func (i *Inventory) SelectedSlotQuantity() uint8 {
	return i.Slots[i.SelectedSlotIndex].Quantity
}

func (i *Inventory) ClearSlot(index int) {
	i.Slots[index] = ItemStack{
		ID:       items.Air,
		Quantity: 0,
	}
}
func (i *Inventory) IsSelectedSlotEmpty() bool {
	return i.Slots[i.SelectedSlotIndex].Quantity <= 0 || i.Slots[i.SelectedSlotIndex].ID == items.Air
}

func (i *Inventory) ClearAllSlots() {
	for idx := range i.Slots {
		i.Slots[idx] = ItemStack{}
	}
}
func (i *Inventory) RandomFillAllSlots() {
	for idx := range i.Slots {
		randItemID := items.RandomItem()
		dur := items.GetDefaultDurability(randItemID)
		i.SetSlot(idx, randItemID, items.Property[randItemID].Stackable, dur)
	}
}

func (i *Inventory) ClearSelectedSlot() {
	i.ClearSlot(i.SelectedSlotIndex)
}

func (i *Inventory) HasEmptySlot() (index int, ok bool) {
	// önce seçili slot boşsa tercih et
	if i.SelectedSlotQuantity() == 0 {
		return i.SelectedSlotIndex, true
	} else {
		for idx, v := range i.Slots {
			if v.Quantity == 0 {
				return idx, true
			}
		}
	}
	return -1, false
}

func (i *Inventory) HasItemStackSpace(id uint16) (index int, ok bool) {
	for idx, v := range i.Slots {
		s := items.Property[v.ID].Stackable
		if v.ID == id && v.Quantity < 64 && v.Quantity > 0 && s != 1 {
			return idx, true
		}
	}
	return -1, false
}

func (i *Inventory) HasItem(id uint16) (index int, ok bool) {
	for idx, v := range i.Slots {
		if v.ID == id && v.Quantity > 0 {
			return idx, true
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
