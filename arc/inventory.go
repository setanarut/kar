package arc

import (
	"kar/items"
)

type SlotData struct {
	ItemID         uint16
	ItemQuantity   uint8
	ItemDurability int
}

type Inventory struct {
	CurrentSlotIndex int
	Slots            [9]SlotData
	HandSlot         SlotData
}

func NewInventory() *Inventory {
	inv := &Inventory{}
	inv.HandSlot = SlotData{}
	for i := range inv.Slots {
		inv.Slots[i] = SlotData{}
	}
	return inv
}

// AddItemIfEmpty adds item to inventory if empty
func (i *Inventory) AddItemIfEmpty(id uint16, dura int) bool {
	idx, ok1 := i.HasItemStackSpace(id)
	if ok1 {
		i.Slots[idx].ItemQuantity++
		return true
	} else {
		i2, ok2 := i.HasEmptySlot()
		if ok2 {
			i.Slots[i2].ItemQuantity++
			i.Slots[i2].ItemID = id
			i.Slots[i2].ItemDurability = dura
			return true
		}
	}
	return false
}

func (i *Inventory) SetSlot(slotIndex int, id uint16, quantity uint8, dur int) {
	if quantity > 0 {
		i.Slots[slotIndex] = SlotData{
			ItemID:         id,
			ItemQuantity:   quantity,
			ItemDurability: dur,
		}
	}
}

func (i *Inventory) SelectNextSlot() {
	i.CurrentSlotIndex = (i.CurrentSlotIndex + 1) % 9
}
func (i *Inventory) SelectPrevSlot() {
	i.CurrentSlotIndex = (i.CurrentSlotIndex - 1 + 9) % 9
}
func (i *Inventory) RemoveHandItem(id uint16) bool {
	ok := i.HasHandItem(id)
	if ok {
		i.HandSlot.ItemQuantity--
		return true
	} else {
		i.HandSlot.ItemID = items.Air
	}
	return false
}

func (i *Inventory) RemoveItem(id uint16) bool {
	idx, ok := i.HasItem(id)
	if ok {
		i.Slots[idx].ItemQuantity--
		return true
	}
	return false
}
func (i *Inventory) RemoveItemFromSelectedSlot() {
	if i.Slots[i.CurrentSlotIndex].ItemQuantity == 1 {
		i.ClearSelectedSlot()
		return
	}
	if i.Slots[i.CurrentSlotIndex].ItemQuantity > 0 {
		i.Slots[i.CurrentSlotIndex].ItemQuantity--
	}
}

func (i *Inventory) SelectedSlot() *SlotData {
	return &i.Slots[i.CurrentSlotIndex]
}

func (i *Inventory) SelectedSlotID() uint16 {
	return i.Slots[i.CurrentSlotIndex].ItemID
}

func (i *Inventory) SelectedSlotQuantity() uint8 {
	return i.Slots[i.CurrentSlotIndex].ItemQuantity
}

func (i *Inventory) ClearSlot(index int) {
	i.Slots[index] = SlotData{
		ItemID:       items.Air,
		ItemQuantity: 0,
	}
}
func (i *Inventory) IsSelectedSlotEmpty() bool {
	return i.Slots[i.CurrentSlotIndex].ItemQuantity <= 0 || i.Slots[i.CurrentSlotIndex].ItemID == items.Air
}

func (i *Inventory) ClearAllSlots() {
	for idx := range i.Slots {
		i.Slots[idx] = SlotData{}
	}
}
func (i *Inventory) RandomFillAllSlots() {
	for idx := range i.Slots {
		randItemID := items.RandomItem()
		dur := items.GetDefaultDurability(randItemID)
		i.SetSlot(idx, randItemID, items.Property[randItemID].MaxStackSize, dur)
	}
}

func (i *Inventory) ClearSelectedSlot() {
	i.ClearSlot(i.CurrentSlotIndex)
}

func (i *Inventory) HasEmptySlot() (index int, ok bool) {
	// önce seçili slot boşsa tercih et
	if i.SelectedSlotQuantity() == 0 {
		return i.CurrentSlotIndex, true
	} else {
		for idx, v := range i.Slots {
			if v.ItemQuantity == 0 {
				return idx, true
			}
		}
	}
	return -1, false
}

func (i *Inventory) HasItemStackSpace(id uint16) (index int, ok bool) {
	for idx, v := range i.Slots {
		s := items.Property[v.ItemID].MaxStackSize
		if v.ItemID == id && v.ItemQuantity < 64 && v.ItemQuantity > 0 && s != 1 {
			return idx, true
		}
	}
	return -1, false
}

func (i *Inventory) HasItem(id uint16) (index int, ok bool) {
	for idx, v := range i.Slots {
		if v.ItemID == id && v.ItemQuantity > 0 {
			return idx, true
		}
	}
	return -1, false
}
func (i *Inventory) HasHandItem(id uint16) bool {
	if i.HandSlot.ItemID == id && i.HandSlot.ItemQuantity > 0 {
		return true
	}
	return false
}
