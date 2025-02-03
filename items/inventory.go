package items

type Slot struct {
	ID         uint8
	Quantity   uint8
	Durability int
}

type Inventory struct {
	CurrentSlotIndex int
	Slots            [9]Slot
	HandSlot         Slot
}

func NewInventory() *Inventory {
	inv := &Inventory{}
	inv.HandSlot = Slot{}
	for i := range inv.Slots {
		inv.Slots[i] = Slot{}
	}
	return inv
}

// AddItemIfEmpty adds item to inventory if empty
func (i *Inventory) AddItemIfEmpty(id uint8, dura int) bool {
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

func (i *Inventory) SetSlot(slotIndex int, id uint8, quantity uint8, dur int) {
	if quantity > 0 {
		i.Slots[slotIndex] = Slot{
			ID:         id,
			Quantity:   quantity,
			Durability: dur,
		}
	}
}

func (i *Inventory) SelectNextSlot() {
	i.CurrentSlotIndex = (i.CurrentSlotIndex + 1) % 9
}
func (i *Inventory) SelectPrevSlot() {
	i.CurrentSlotIndex = (i.CurrentSlotIndex - 1 + 9) % 9
}

func (i *Inventory) RemoveHandItem(id uint8) bool {
	ok := i.HasHandItem(id)
	if ok {
		i.HandSlot.Quantity--
		return true
	} else {
		i.HandSlot.ID = Air
	}
	return false
}

func (i *Inventory) RemoveItem(id uint8) bool {
	idx, ok := i.HasItem(id)
	if ok {
		i.Slots[idx].Quantity--
		return true
	}
	return false
}
func (i *Inventory) RemoveItemFromSelectedSlot() (uint8, int) {
	quantity := i.CurrentSlotQuantity()
	id := i.CurrentSlotID()
	dura := i.CurrentSlot().Durability
	if quantity == 1 {
		i.ClearCurrentSlot()
		return id, dura
	}
	if quantity > 0 {

		i.Slots[i.CurrentSlotIndex].Quantity--
		return id, dura
	}
	return 0, 0
}

func (i *Inventory) CurrentSlot() *Slot {
	return &i.Slots[i.CurrentSlotIndex]
}

func (i *Inventory) CurrentSlotID() uint8 {
	return i.Slots[i.CurrentSlotIndex].ID
}

func (i *Inventory) CurrentSlotQuantity() uint8 {
	return i.Slots[i.CurrentSlotIndex].Quantity
}

func (i *Inventory) ClearSlot(index int) {
	i.Slots[index] = Slot{}
}
func (i *Inventory) IsCurrentSlotEmpty() bool {
	return i.Slots[i.CurrentSlotIndex].Quantity <= 0 || i.Slots[i.CurrentSlotIndex].ID == Air
}

func (i *Inventory) ClearAllSlots() {
	for idx := range i.Slots {
		i.Slots[idx] = Slot{}
	}
}
func (i *Inventory) RandomFillAllSlots() {
	for idx := range i.Slots {
		randItemID := RandomItem()
		dur := GetDefaultDurability(randItemID)
		i.SetSlot(idx, randItemID, Property[randItemID].MaxStackSize, dur)
	}
}

func (i *Inventory) ClearCurrentSlot() {
	i.ClearSlot(i.CurrentSlotIndex)
}

func (i *Inventory) HasEmptySlot() (index int, ok bool) {
	// önce seçili slot boşsa tercih et
	if i.CurrentSlotQuantity() == 0 {
		return i.CurrentSlotIndex, true
	} else {
		for idx, v := range i.Slots {
			if v.Quantity == 0 {
				return idx, true
			}
		}
	}
	return -1, false
}

func (i *Inventory) HasItemStackSpace(id uint8) (index int, ok bool) {
	for idx, v := range i.Slots {
		s := Property[v.ID].MaxStackSize
		if v.ID == id && v.Quantity < 64 && v.Quantity > 0 && s != 1 {
			return idx, true
		}
	}
	return -1, false
}

func (i *Inventory) HasItem(id uint8) (index int, ok bool) {
	for idx, v := range i.Slots {
		if v.ID == id && v.Quantity > 0 {
			return idx, true
		}
	}
	return -1, false
}
func (i *Inventory) HasHandItem(id uint8) bool {
	if i.HandSlot.ID == id && i.HandSlot.Quantity > 0 {
		return true
	}
	return false
}
