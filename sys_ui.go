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

type UI struct {
	hotbarPos           Vec
	craftingTablePos    Vec
	hotbarRightEdgePosX float64
}

func (ui *UI) Init() {
	ui.hotbarPos = Vec{4, 9}
	ui.hotbarRightEdgePosX = ui.hotbarPos.X + float64(res.Hotbar.Bounds().Dx())
	ui.craftingTablePos = ui.hotbarPos.Add(Vec{49, 39})
}
func (ui *UI) Update() {

	if world.Alive(currentPlayer) {

		// Toggle crafting state
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			switch gameDataRes.GameplayState {
			case Playing:
				if tileMapRes.Get(gameDataRes.TargetBlockCoord.X, gameDataRes.TargetBlockCoord.Y) == items.CraftingTable {
					gameDataRes.GameplayState = CraftingTable3x3
				} else {
					gameDataRes.GameplayState = CraftingTable2x2
				}
			case CraftingTable3x3, CraftingTable2x2:
				// clear crafting table when exit
				for y := range 3 {
					for x := range 3 {
						itemID := craftingTableRes.Slots[y][x].ID
						if itemID != 0 {
							quantity := craftingTableRes.Slots[y][x].Quantity
							for range quantity {
								durability := craftingTableRes.Slots[y][x].Durability
								// move items from crafting table to hotbar if possible
								if inventoryRes.AddItemIfEmpty(craftingTableRes.Slots[y][x].ID, durability) {
									craftingTableRes.RemoveItem(x, y)
								} else {
									// move items from crafting table to world if hotbar is full
									craftingTableRes.RemoveItem(x, y)
									playerPos := MapPosition.GetUnchecked(currentPlayer)
									SpawnItem(
										playerPos.X+8,
										playerPos.Y+8,
										itemID,
										craftingTableRes.Slots[y][x].Durability,
									)
								}
							}
						}
					}
				}
				craftingTableRes.ResultSlot = items.Slot{}
				gameDataRes.GameplayState = Playing
			}
			onInventorySlotChanged()
		}

		// Hotbar slot navigation
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			inventoryRes.SelectPrevSlot()
			onInventorySlotChanged()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			inventoryRes.SelectNextSlot()
			onInventorySlotChanged()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			if inventoryRes.CurrentSlotIndex != inventoryRes.QuickSlot2 {
				inventoryRes.QuickSlot1 = inventoryRes.CurrentSlotIndex
			}

		}
		if inpututil.IsKeyJustPressed(ebiten.KeyT) {
			if inventoryRes.CurrentSlotIndex != inventoryRes.QuickSlot1 {
				inventoryRes.QuickSlot2 = inventoryRes.CurrentSlotIndex
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
			switch inventoryRes.CurrentSlotIndex {
			case inventoryRes.QuickSlot1:
				inventoryRes.CurrentSlotIndex = inventoryRes.QuickSlot2
			case inventoryRes.QuickSlot2:
				inventoryRes.CurrentSlotIndex = inventoryRes.QuickSlot1
			default:
				inventoryRes.CurrentSlotIndex = inventoryRes.QuickSlot1
			}
			onInventorySlotChanged()
		}

		switch gameDataRes.GameplayState {
		case CraftingTable2x2, CraftingTable3x3:
			if inpututil.IsKeyJustPressed(ebiten.KeyD) {
				if gameDataRes.GameplayState == CraftingTable2x2 {
					craftingTableRes.SlotPosX = min(craftingTableRes.SlotPosX+1, 1)
				} else {
					craftingTableRes.SlotPosX = min(craftingTableRes.SlotPosX+1, 2)
				}
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyA) {
				craftingTableRes.SlotPosX = max(craftingTableRes.SlotPosX-1, 0)
			}

			if inpututil.IsKeyJustPressed(ebiten.KeyS) {
				if gameDataRes.GameplayState == CraftingTable2x2 {
					craftingTableRes.SlotPosY = min(craftingTableRes.SlotPosY+1, 1)
				} else {
					craftingTableRes.SlotPosY = min(craftingTableRes.SlotPosY+1, 2)
				}
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyW) {
				craftingTableRes.SlotPosY = max(craftingTableRes.SlotPosY-1, 0)
			}

			// Move items from hotbar to crafting table
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
				cs := craftingTableRes.CurrentSlot()
				if inventoryRes.CurrentSlotID() != 0 {
					if craftingTableRes.CurrentSlot().ID == 0 {
						id, dur := inventoryRes.RemoveItemFromSelectedSlot()
						cs.ID = id
						cs.Durability = dur
						cs.Quantity = 1
					} else if cs.ID == inventoryRes.CurrentSlotID() {
						inventoryRes.RemoveItemFromSelectedSlot()
						cs.Quantity++
					}
				}
				craftingTableRes.UpdateResultSlot()
				onInventorySlotChanged()
			}
			// Move items from crafting table to hotbar
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
				cs := craftingTableRes.CurrentSlot()
				if cs.ID != 0 {
					if cs.Quantity == 1 {
						if inventoryRes.AddItemIfEmpty(cs.ID, cs.Durability) {
							craftingTableRes.ClearCurrenSlot()
						}
					} else if cs.Quantity > 1 {
						if inventoryRes.AddItemIfEmpty(cs.ID, cs.Durability) {
							cs.Quantity--

						}
					}
				}
				craftingTableRes.UpdateResultSlot()
				onInventorySlotChanged()
			}
			// apply recipe
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
				minimum := craftingTableRes.UpdateResultSlot()
				resultID := craftingTableRes.ResultSlot.ID
				dur := items.GetDefaultDurability(resultID)
				if resultID != 0 {
					for range minimum {
						if inventoryRes.AddItemIfEmpty(resultID, dur) {
							for y := range 3 {
								for x := range 3 {
									if craftingTableRes.Slots[y][x].Quantity > 0 {
										craftingTableRes.Slots[y][x].Quantity--
									}
									if craftingTableRes.Slots[y][x].Quantity == 0 {
										craftingTableRes.Slots[y][x].ID = 0
									}
								}
							}
							craftingTableRes.ResultSlot.ID = 0
						}
					}
				}
				craftingTableRes.UpdateResultSlot()
				onInventorySlotChanged()
			}
		}

	}
}

func (ui *UI) Draw() {
	if world.Alive(currentPlayer) {

		// Draw hotbar background
		ColorMDIO.GeoM.Reset()
		ColorMDIO.GeoM.Translate(ui.hotbarPos.X, ui.hotbarPos.Y)
		colorm.DrawImage(Screen, res.Hotbar, ColorM, ColorMDIO)

		// Draw Quick-slot numbers
		dotX := float64(inventoryRes.QuickSlot1)*17 + float64(ui.hotbarPos.X) + 7
		TextDO.GeoM.Reset()
		TextDO.GeoM.Translate(dotX, -4)
		text.Draw(Screen, strconv.Itoa(inventoryRes.QuickSlot1+1), res.Font, TextDO)
		dotX = float64(inventoryRes.QuickSlot2)*17 + float64(ui.hotbarPos.X) + 7
		TextDO.GeoM.Reset()
		TextDO.GeoM.Translate(dotX, -4)
		text.Draw(Screen, strconv.Itoa(inventoryRes.QuickSlot2+1), res.Font, TextDO)

		// Draw slots
		for x := range 9 {
			slotID := inventoryRes.Slots[x].ID
			quantity := inventoryRes.Slots[x].Quantity
			SlotOffsetX := float64(x) * 17
			SlotOffsetX += ui.hotbarPos.X

			// draw hotbar item icons
			ColorMDIO.GeoM.Reset()
			ColorMDIO.GeoM.Translate(SlotOffsetX+(5), ui.hotbarPos.Y+(5))
			if slotID != items.Air && inventoryRes.Slots[x].Quantity > 0 {
				colorm.DrawImage(Screen, res.Icon8[slotID], ColorM, ColorMDIO)
			}
			// Draw hotbar selected slot border
			if x == inventoryRes.CurrentSlotIndex {
				// border
				ColorMDIO.GeoM.Translate(-5, -5)
				colorm.DrawImage(Screen, res.SlotBorder, ColorM, ColorMDIO)
				// Draw hotbar slot item display name
				if !inventoryRes.IsCurrentSlotEmpty() {
					TextDO.GeoM.Reset()
					TextDO.GeoM.Translate(SlotOffsetX-1, ui.hotbarPos.Y+14)
					if items.HasTag(slotID, items.Tool) {
						text.Draw(Screen, fmt.Sprintf(
							"%v\nDurability %v",
							items.Property[slotID].DisplayName,
							inventoryRes.Slots[x].Durability,
						), res.Font, TextDO)
					} else {
						text.Draw(Screen, items.Property[slotID].DisplayName, res.Font, TextDO)
					}
				}
			}

			// Draw item quantity number
			if quantity > 1 && items.IsStackable(slotID) {
				TextDO.GeoM.Reset()
				TextDO.GeoM.Translate(SlotOffsetX+6, ui.hotbarPos.Y+4)
				num := strconv.FormatUint(uint64(quantity), 10)
				if quantity < 10 {
					num = " " + num
				}
				text.Draw(Screen, num, res.Font, TextDO)
			}
		}

		// Draw player health text
		TextDO.GeoM.Reset()
		TextDO.GeoM.Translate(ui.hotbarRightEdgePosX+8, ui.hotbarPos.Y)
		playerHealth := MapHealth.GetUnchecked(currentPlayer)
		text.Draw(Screen, fmt.Sprintf("Health %v", playerHealth.Current), res.Font, TextDO)

		switch gameDataRes.GameplayState {
		case CraftingTable2x2, CraftingTable3x3:
			// crafting table Background
			ColorMDIO.GeoM.Reset()
			ColorMDIO.GeoM.Translate(ui.craftingTablePos.X, ui.craftingTablePos.Y)

			if gameDataRes.GameplayState == CraftingTable2x2 {
				colorm.DrawImage(Screen, res.CraftingTable4, ColorM, ColorMDIO)
			} else {
				colorm.DrawImage(Screen, res.CraftingTable, ColorM, ColorMDIO)
			}

			// draw crafting table item icons
			for x := 0; x < 3; x++ {
				for y := 0; y < 3; y++ {
					if craftingTableRes.Slots[y][x].ID != items.Air {
						sx := ui.craftingTablePos.X + float64(x*17)
						sy := ui.craftingTablePos.Y + float64(y*17)
						ColorMDIO.GeoM.Reset()
						ColorMDIO.GeoM.Translate(sx+6, sy+6)
						colorm.DrawImage(
							Screen,
							res.Icon8[craftingTableRes.Slots[y][x].ID],
							ColorM,
							ColorMDIO,
						)

						// Draw item quantity number
						quantity := craftingTableRes.Slots[y][x].Quantity
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
					if x == craftingTableRes.SlotPosX && y == craftingTableRes.SlotPosY {
						sx := ui.craftingTablePos.X + float64(x*17)
						sy := ui.craftingTablePos.Y + float64(y*17)
						ColorMDIO.GeoM.Reset()
						ColorMDIO.GeoM.Translate(sx+1, sy+1)
						colorm.DrawImage(Screen, res.SlotBorder, ColorM, ColorMDIO)
					}

				}
			}

			// draw crafting table result item icon
			if craftingTableRes.ResultSlot.ID != 0 {
				ColorMDIO.GeoM.Reset()

				if gameDataRes.GameplayState == CraftingTable2x2 {
					ColorMDIO.GeoM.Translate(ui.craftingTablePos.X+41, ui.craftingTablePos.Y+14)
				} else {
					ColorMDIO.GeoM.Translate(ui.craftingTablePos.X+58, ui.craftingTablePos.Y+23)
				}

				colorm.DrawImage(Screen, res.Icon8[craftingTableRes.ResultSlot.ID], ColorM, ColorMDIO)

				// Draw result item quantity number
				quantity := craftingTableRes.ResultSlot.Quantity
				if quantity > 1 {
					TextDO.GeoM.Reset()
					if gameDataRes.GameplayState == CraftingTable2x2 {
						TextDO.GeoM.Translate(ui.craftingTablePos.X+42, ui.craftingTablePos.Y+13)
					} else {
						TextDO.GeoM.Translate(ui.craftingTablePos.X+58, ui.craftingTablePos.Y+22)
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
			_, vel, _, playerController, _ := MapPlayer.Get(currentPlayer)
			ebitenutil.DebugPrintAt(Screen, fmt.Sprintf(
				"state %v\nVel.X: %.2f\nVel.Y: %.2f",
				playerController.CurrentState,
				vel.X,
				vel.Y,
			), 10, 50)
		}

	} else if currentGameState != "menu" {
		ebitenutil.DebugPrintAt(Screen, "YOU ARE DEAD!", int(ScreenSize.X/2)-30, int(ScreenSize.Y/2))
	}
}

func onInventorySlotChanged() {
	switch inventoryRes.CurrentSlotID() {
	case items.WoodenAxe:
		animPlayer.SetAtlas("WoodenAxe")
	case items.WoodenPickaxe:
		animPlayer.SetAtlas("WoodenPickaxe")
	case items.WoodenShovel:
		animPlayer.SetAtlas("WoodenShovel")
	case items.StoneAxe:
		animPlayer.SetAtlas("StoneAxe")
	case items.StonePickaxe:
		animPlayer.SetAtlas("StonePickaxe")
	case items.StoneShovel:
		animPlayer.SetAtlas("StoneShovel")
	case items.IronAxe:
		animPlayer.SetAtlas("IronAxe")
	case items.IronPickaxe:
		animPlayer.SetAtlas("IronPickaxe")
	case items.IronShovel:
		animPlayer.SetAtlas("IronShovel")
	case items.DiamondAxe:
		animPlayer.SetAtlas("DiamondAxe")
	case items.DiamondPickaxe:
		animPlayer.SetAtlas("DiamondPickaxe")
	case items.DiamondShovel:
		animPlayer.SetAtlas("DiamondShovel")
	default:
		animPlayer.SetAtlas("Default")
	}
}
