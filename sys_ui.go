package kar

import (
	"fmt"
	"kar/items"
	"kar/res"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	// hotbarPositionX        = ScreenW/2 - float64(res.Hotbar.Bounds().Dx())/2
	hotbarPositionX        = 4.
	hotbarPositionY        = 9.
	hotbarRightEdgePosX    = hotbarPositionX + float64(res.Hotbar.Bounds().Dx())
	craftingTablePositionX = hotbarPositionX + 49
	craftingTablePositionY = hotbarPositionY + 39
)

type UI struct{}

func (ui *UI) Init() {}
func (ui *UI) Update() error {

	if ECWorld.Alive(CurrentPlayer) {

		// Hotbar slot navigation
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			InventoryRes.SelectPrevSlot()
			onInventorySlotChanged()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			InventoryRes.SelectNextSlot()
			onInventorySlotChanged()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			if InventoryRes.CurrentSlotIndex != InventoryRes.QuickSlot2 {
				InventoryRes.QuickSlot1 = InventoryRes.CurrentSlotIndex
			}

		}
		if inpututil.IsKeyJustPressed(ebiten.KeyT) {
			if InventoryRes.CurrentSlotIndex != InventoryRes.QuickSlot1 {
				InventoryRes.QuickSlot2 = InventoryRes.CurrentSlotIndex
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
			switch InventoryRes.CurrentSlotIndex {
			case InventoryRes.QuickSlot1:
				InventoryRes.CurrentSlotIndex = InventoryRes.QuickSlot2
			case InventoryRes.QuickSlot2:
				InventoryRes.CurrentSlotIndex = InventoryRes.QuickSlot1
			default:
				InventoryRes.CurrentSlotIndex = InventoryRes.QuickSlot1
			}
			onInventorySlotChanged()
		}

		// Toggle crafting state
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			if TileMapRes.Get(GameDataRes.TargetBlockCoord.X, GameDataRes.TargetBlockCoord.Y) == items.CraftingTable {
				GameDataRes.CraftingState4 = false
			} else {
				GameDataRes.CraftingState4 = true
			}

			CraftingTableRes.SlotPosX = 1
			CraftingTableRes.SlotPosY = 1

			// clear crafting table when exit
			if GameDataRes.CraftingState {
				for y := range 3 {
					for x := range 3 {
						itemID := CraftingTableRes.Slots[y][x].ID
						if itemID != 0 {
							quantity := CraftingTableRes.Slots[y][x].Quantity
							for range quantity {
								durability := CraftingTableRes.Slots[y][x].Durability
								// move items from crafting table to hotbar if possible
								if InventoryRes.AddItemIfEmpty(CraftingTableRes.Slots[y][x].ID, durability) {
									CraftingTableRes.RemoveItem(x, y)
								} else {
									// move items from crafting table to world if hotbar is full
									CraftingTableRes.RemoveItem(x, y)
									playerPos := MapPosition.GetUnchecked(CurrentPlayer)
									SpawnItem(
										playerPos.X+8,
										playerPos.Y+8,
										itemID,
										CraftingTableRes.Slots[y][x].Durability,
									)
								}

							}
						}
					}
				}
				CraftingTableRes.ResultSlot = items.Slot{}
			}

			GameDataRes.CraftingState = !GameDataRes.CraftingState

			onInventorySlotChanged()

		}

		if GameDataRes.CraftingState {
			// Crafting table slot navigation
			if inpututil.IsKeyJustPressed(ebiten.KeyD) {
				if GameDataRes.CraftingState4 {
					CraftingTableRes.SlotPosX = min(CraftingTableRes.SlotPosX+1, 1)
				} else {
					CraftingTableRes.SlotPosX = min(CraftingTableRes.SlotPosX+1, 2)
				}
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyA) {
				CraftingTableRes.SlotPosX = max(CraftingTableRes.SlotPosX-1, 0)
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyS) {
				if GameDataRes.CraftingState4 {
					CraftingTableRes.SlotPosY = min(CraftingTableRes.SlotPosY+1, 1)
				} else {
					CraftingTableRes.SlotPosY = min(CraftingTableRes.SlotPosY+1, 2)
				}
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyW) {
				CraftingTableRes.SlotPosY = max(CraftingTableRes.SlotPosY-1, 0)
			}

			// Move items from hotbar to crafting table
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
				cs := CraftingTableRes.CurrentSlot()
				if InventoryRes.CurrentSlotID() != 0 {
					if CraftingTableRes.CurrentSlot().ID == 0 {
						id, dur := InventoryRes.RemoveItemFromSelectedSlot()
						cs.ID = id
						cs.Durability = dur
						cs.Quantity = 1
					} else if cs.ID == InventoryRes.CurrentSlotID() {
						InventoryRes.RemoveItemFromSelectedSlot()
						cs.Quantity++
					}
				}
				CraftingTableRes.UpdateResultSlot()
				onInventorySlotChanged()
			}
			// Move items from crafting table to hotbar
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
				cs := CraftingTableRes.CurrentSlot()
				if cs.ID != 0 {
					if cs.Quantity == 1 {
						if InventoryRes.AddItemIfEmpty(cs.ID, cs.Durability) {
							CraftingTableRes.ClearCurrenSlot()
						}
					} else if cs.Quantity > 1 {
						if InventoryRes.AddItemIfEmpty(cs.ID, cs.Durability) {
							cs.Quantity--

						}
					}

				}
				CraftingTableRes.UpdateResultSlot()
				onInventorySlotChanged()
			}
			// apply recipe
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
				minimum := CraftingTableRes.UpdateResultSlot()
				resultID := CraftingTableRes.ResultSlot.ID
				dur := items.GetDefaultDurability(resultID)
				if resultID != 0 {
					for range minimum {
						if InventoryRes.AddItemIfEmpty(resultID, dur) {
							for y := range 3 {
								for x := range 3 {
									if CraftingTableRes.Slots[y][x].Quantity > 0 {
										CraftingTableRes.Slots[y][x].Quantity--
									}
									if CraftingTableRes.Slots[y][x].Quantity == 0 {
										CraftingTableRes.Slots[y][x].ID = 0
									}
								}
							}
							CraftingTableRes.ResultSlot.ID = 0
						}
					}
				}
				CraftingTableRes.UpdateResultSlot()
				onInventorySlotChanged()
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			InventoryRes.ClearCurrentSlot()
		}
	}
	return nil
}

func (ui *UI) Draw() {
	if ECWorld.Alive(CurrentPlayer) {

		// Draw hotbar background
		ColorMDIO.GeoM.Reset()
		ColorMDIO.GeoM.Translate(hotbarPositionX, hotbarPositionY)
		colorm.DrawImage(Screen, res.Hotbar, ColorM, ColorMDIO)

		// Draw Quick-slot numbers
		dotX := float64(InventoryRes.QuickSlot1)*17 + float64(hotbarPositionX) + 7
		TextDO.GeoM.Reset()
		TextDO.GeoM.Translate(dotX, -4)
		text.Draw(Screen, strconv.Itoa(InventoryRes.QuickSlot1+1), res.Font, TextDO)
		dotX = float64(InventoryRes.QuickSlot2)*17 + float64(hotbarPositionX) + 7
		TextDO.GeoM.Reset()
		TextDO.GeoM.Translate(dotX, -4)
		text.Draw(Screen, strconv.Itoa(InventoryRes.QuickSlot2+1), res.Font, TextDO)

		// Draw slots
		for x := range 9 {
			slotID := InventoryRes.Slots[x].ID
			quantity := InventoryRes.Slots[x].Quantity
			SlotOffsetX := float64(x) * 17
			SlotOffsetX += hotbarPositionX

			// draw hotbar item icons
			ColorMDIO.GeoM.Reset()
			ColorMDIO.GeoM.Translate(SlotOffsetX+(5), hotbarPositionY+(5))
			if slotID != items.Air && InventoryRes.Slots[x].Quantity > 0 {
				colorm.DrawImage(Screen, res.Icon8[slotID], ColorM, ColorMDIO)
			}
			// Draw hotbar selected slot border
			if x == InventoryRes.CurrentSlotIndex {
				// border
				ColorMDIO.GeoM.Translate(-5, -5)
				colorm.DrawImage(Screen, res.SlotBorder, ColorM, ColorMDIO)
				// Draw hotbar slot item display name
				if !InventoryRes.IsCurrentSlotEmpty() {
					TextDO.GeoM.Reset()
					TextDO.GeoM.Translate(SlotOffsetX-1, hotbarPositionY+14)
					if items.HasTag(slotID, items.Tool) {
						text.Draw(Screen, fmt.Sprintf(
							"%v\nDurability %v",
							items.Property[slotID].DisplayName,
							InventoryRes.Slots[x].Durability,
						), res.Font, TextDO)
					} else {
						text.Draw(Screen, items.Property[slotID].DisplayName, res.Font, TextDO)
					}
				}
			}

			// Draw item quantity number
			if quantity > 1 && items.IsStackable(slotID) {
				TextDO.GeoM.Reset()
				TextDO.GeoM.Translate(SlotOffsetX+6, hotbarPositionY+4)
				num := strconv.FormatUint(uint64(quantity), 10)
				if quantity < 10 {
					num = " " + num
				}
				text.Draw(Screen, num, res.Font, TextDO)
			}
		}

		// Draw player health text
		TextDO.GeoM.Reset()
		TextDO.GeoM.Translate(hotbarRightEdgePosX+8, hotbarPositionY)
		playerHealth := MapHealth.GetUnchecked(CurrentPlayer)
		text.Draw(Screen, fmt.Sprintf("Health %v", playerHealth.Current), res.Font, TextDO)

		// Draw crafting table
		if GameDataRes.CraftingState {

			// crafting table Background
			ColorMDIO.GeoM.Reset()
			ColorMDIO.GeoM.Translate(craftingTablePositionX, craftingTablePositionY)

			if GameDataRes.CraftingState4 {
				colorm.DrawImage(Screen, res.CraftingTable4, ColorM, ColorMDIO)
			} else {
				colorm.DrawImage(Screen, res.CraftingTable, ColorM, ColorMDIO)
			}

			// draw crafting table item icons
			for x := 0; x < 3; x++ {
				for y := 0; y < 3; y++ {
					if CraftingTableRes.Slots[y][x].ID != items.Air {
						sx := craftingTablePositionX + float64(x*17)
						sy := craftingTablePositionY + float64(y*17)
						ColorMDIO.GeoM.Reset()
						ColorMDIO.GeoM.Translate(sx+6, sy+6)
						colorm.DrawImage(
							Screen,
							res.Icon8[CraftingTableRes.Slots[y][x].ID],
							ColorM,
							ColorMDIO,
						)

						// Draw item quantity number
						quantity := CraftingTableRes.Slots[y][x].Quantity
						if quantity > 1 {
							TextDO.GeoM.Reset()
							TextDO.GeoM.Translate(sx+7, sy+5)
							num := strconv.FormatUint(uint64(quantity), 10)
							if quantity < 10 {
								num = " " + num
							}
							text.Draw(Screen, num, res.Font, TextDO)
						}
					}

					// draw selected slot border of crqfting table
					if x == CraftingTableRes.SlotPosX && y == CraftingTableRes.SlotPosY {
						sx := craftingTablePositionX + float64(x*17)
						sy := craftingTablePositionY + float64(y*17)
						ColorMDIO.GeoM.Reset()
						ColorMDIO.GeoM.Translate(sx+1, sy+1)
						colorm.DrawImage(Screen, res.SlotBorder, ColorM, ColorMDIO)
					}

				}
			}

			// draw crafting table result item icon
			if CraftingTableRes.ResultSlot.ID != 0 {
				ColorMDIO.GeoM.Reset()

				if GameDataRes.CraftingState4 {
					ColorMDIO.GeoM.Translate(craftingTablePositionX+41, craftingTablePositionY+14)
				} else {
					ColorMDIO.GeoM.Translate(craftingTablePositionX+58, craftingTablePositionY+23)
				}

				colorm.DrawImage(Screen, res.Icon8[CraftingTableRes.ResultSlot.ID], ColorM, ColorMDIO)

				// Draw result item quantity number
				quantity := CraftingTableRes.ResultSlot.Quantity
				if quantity > 1 {
					TextDO.GeoM.Reset()
					if GameDataRes.CraftingState4 {
						TextDO.GeoM.Translate(craftingTablePositionX+42, craftingTablePositionY+13)
					} else {
						TextDO.GeoM.Translate(craftingTablePositionX+58, craftingTablePositionY+22)
					}
					num := strconv.FormatUint(uint64(quantity), 10)
					if quantity < 10 {
						num = " " + num
					}
					text.Draw(Screen, num, res.Font, TextDO)
				}
			}

		}

		// Draw debug info
		if DrawDebugTextEnabled {
			_, vel, _, playerController, _ := MapPlayer.Get(CurrentPlayer)
			ebitenutil.DebugPrintAt(Screen, fmt.Sprintf(
				"state %v\nAbsVelocity: %v",
				playerController.CurrentState,
				vel,
			), 10, 50)
		}

	} else if CurrentGameState != "menu" {
		ebitenutil.DebugPrintAt(Screen, "YOU ARE DEAD!", int(ScreenW/2)-30, int(ScreenH/2))
	}
}

func onInventorySlotChanged() {
	switch InventoryRes.CurrentSlotID() {
	case items.WoodenAxe:
		PlayerAnimPlayer.SetAtlas("WoodenAxe")
	case items.WoodenPickaxe:
		PlayerAnimPlayer.SetAtlas("WoodenPickaxe")
	case items.WoodenShovel:
		PlayerAnimPlayer.SetAtlas("WoodenShovel")
	case items.StoneAxe:
		PlayerAnimPlayer.SetAtlas("StoneAxe")
	case items.StonePickaxe:
		PlayerAnimPlayer.SetAtlas("StonePickaxe")
	case items.StoneShovel:
		PlayerAnimPlayer.SetAtlas("StoneShovel")
	case items.IronAxe:
		PlayerAnimPlayer.SetAtlas("IronAxe")
	case items.IronPickaxe:
		PlayerAnimPlayer.SetAtlas("IronPickaxe")
	case items.IronShovel:
		PlayerAnimPlayer.SetAtlas("IronShovel")
	case items.DiamondAxe:
		PlayerAnimPlayer.SetAtlas("DiamondAxe")
	case items.DiamondPickaxe:
		PlayerAnimPlayer.SetAtlas("DiamondPickaxe")
	case items.DiamondShovel:
		PlayerAnimPlayer.SetAtlas("DiamondShovel")
	default:
		PlayerAnimPlayer.SetAtlas("Default")
	}
}
