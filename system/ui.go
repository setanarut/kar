package system

import (
	"fmt"
	"kar"
	"kar/arc"
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
	hotbarPositionX        = kar.ScreenW/2 - float64(res.Hotbar.Bounds().Dx())/2
	hotbarPositionY        = 9.
	hotbarRightEdgePosX    = hotbarPositionX + float64(res.Hotbar.Bounds().Dx())
	craftingTablePositionX = hotbarPositionX + 49
	craftingTablePositionY = hotbarPositionY + 39
	debugInfo              = `state %v
Facing %v
`
)

type UI struct{}

func (ui *UI) Init() {}
func (ui *UI) Update() {

	if kar.ECWorld.Alive(kar.CurrentPlayer) {

		// Hotbar slot navigation
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			kar.InventoryRes.SelectPrevSlot()
			onInventorySlotChanged()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			kar.InventoryRes.SelectNextSlot()
			onInventorySlotChanged()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			if kar.InventoryRes.CurrentSlotIndex != kar.InventoryRes.QuickSlot2 {
				kar.InventoryRes.QuickSlot1 = kar.InventoryRes.CurrentSlotIndex
			}

		}
		if inpututil.IsKeyJustPressed(ebiten.KeyF) {
			if kar.InventoryRes.CurrentSlotIndex != kar.InventoryRes.QuickSlot1 {
				kar.InventoryRes.QuickSlot2 = kar.InventoryRes.CurrentSlotIndex
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
			switch kar.InventoryRes.CurrentSlotIndex {
			case kar.InventoryRes.QuickSlot1:
				kar.InventoryRes.CurrentSlotIndex = kar.InventoryRes.QuickSlot2
			case kar.InventoryRes.QuickSlot2:
				kar.InventoryRes.CurrentSlotIndex = kar.InventoryRes.QuickSlot1
			default:
				kar.InventoryRes.CurrentSlotIndex = kar.InventoryRes.QuickSlot1
			}
			onInventorySlotChanged()
		}

		// Toggle crafting state
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			if kar.TileMapRes.Get(kar.GameDataRes.TargetBlockCoord.X, kar.GameDataRes.TargetBlockCoord.Y) == items.CraftingTable {
				kar.GameDataRes.CraftingState4 = false
			} else {
				kar.GameDataRes.CraftingState4 = true
			}

			kar.CraftingTableRes.SlotPosX = 1
			kar.CraftingTableRes.SlotPosY = 1

			// clear crafting table when exit
			if kar.GameDataRes.CraftingState {
				for y := range 3 {
					for x := range 3 {
						itemID := kar.CraftingTableRes.Slots[y][x].ID
						if itemID != 0 {
							quantity := kar.CraftingTableRes.Slots[y][x].Quantity
							for range quantity {
								durability := kar.CraftingTableRes.Slots[y][x].Durability
								// move items from crafting table to hotbar if possible
								if kar.InventoryRes.AddItemIfEmpty(kar.CraftingTableRes.Slots[y][x].ID, durability) {
									kar.CraftingTableRes.RemoveItem(x, y)
								} else {
									// move items from crafting table to world if hotbar is full
									kar.CraftingTableRes.RemoveItem(x, y)
									playerPos := arc.MapPosition.GetUnchecked(kar.CurrentPlayer)
									arc.SpawnItem(
										playerPos.X+8,
										playerPos.Y+8,
										itemID,
										kar.CraftingTableRes.Slots[y][x].Durability,
									)
								}

							}
						}
					}
				}
				kar.CraftingTableRes.ResultSlot = items.Slot{}
			}

			kar.GameDataRes.CraftingState = !kar.GameDataRes.CraftingState

			onInventorySlotChanged()

		}

		if kar.GameDataRes.CraftingState {
			// Crafting table slot navigation
			if inpututil.IsKeyJustPressed(ebiten.KeyD) {
				if kar.GameDataRes.CraftingState4 {
					kar.CraftingTableRes.SlotPosX = min(kar.CraftingTableRes.SlotPosX+1, 1)
				} else {
					kar.CraftingTableRes.SlotPosX = min(kar.CraftingTableRes.SlotPosX+1, 2)
				}
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyA) {
				kar.CraftingTableRes.SlotPosX = max(kar.CraftingTableRes.SlotPosX-1, 0)
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyS) {
				if kar.GameDataRes.CraftingState4 {
					kar.CraftingTableRes.SlotPosY = min(kar.CraftingTableRes.SlotPosY+1, 1)
				} else {
					kar.CraftingTableRes.SlotPosY = min(kar.CraftingTableRes.SlotPosY+1, 2)
				}
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyW) {
				kar.CraftingTableRes.SlotPosY = max(kar.CraftingTableRes.SlotPosY-1, 0)
			}

			// Move items from hotbar to crafting table
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
				cs := kar.CraftingTableRes.CurrentSlot()
				if kar.InventoryRes.CurrentSlotID() != 0 {
					if kar.CraftingTableRes.CurrentSlot().ID == 0 {
						id, dur := kar.InventoryRes.RemoveItemFromSelectedSlot()
						cs.ID = id
						cs.Durability = dur
						cs.Quantity = 1
					} else if cs.ID == kar.InventoryRes.CurrentSlotID() {
						kar.InventoryRes.RemoveItemFromSelectedSlot()
						cs.Quantity++
					}
				}
				kar.CraftingTableRes.UpdateResultSlot()
				onInventorySlotChanged()
			}
			// Move items from crafting table to hotbar
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
				cs := kar.CraftingTableRes.CurrentSlot()
				if cs.ID != 0 {
					if cs.Quantity == 1 {
						if kar.InventoryRes.AddItemIfEmpty(cs.ID, cs.Durability) {
							kar.CraftingTableRes.ClearCurrenSlot()
						}
					} else if cs.Quantity > 1 {
						if kar.InventoryRes.AddItemIfEmpty(cs.ID, cs.Durability) {
							cs.Quantity--

						}
					}

				}
				kar.CraftingTableRes.UpdateResultSlot()
				onInventorySlotChanged()
			}
			// apply recipe
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
				minimum := kar.CraftingTableRes.UpdateResultSlot()
				resultID := kar.CraftingTableRes.ResultSlot.ID
				dur := items.GetDefaultDurability(resultID)
				if resultID != 0 {
					for range minimum {
						if kar.InventoryRes.AddItemIfEmpty(resultID, dur) {
							for y := range 3 {
								for x := range 3 {
									if kar.CraftingTableRes.Slots[y][x].Quantity > 0 {
										kar.CraftingTableRes.Slots[y][x].Quantity--
									}
									if kar.CraftingTableRes.Slots[y][x].Quantity == 0 {
										kar.CraftingTableRes.Slots[y][x].ID = 0
									}
								}
							}
							kar.CraftingTableRes.ResultSlot.ID = 0
						}
					}
				}
				kar.CraftingTableRes.UpdateResultSlot()
				onInventorySlotChanged()
			}
		}

		// Debug
		if inpututil.IsKeyJustPressed(ebiten.KeyV) {
			kar.DrawDebugTextEnabled = !kar.DrawDebugTextEnabled
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			kar.DrawDebugHitboxesEnabled = !kar.DrawDebugHitboxesEnabled
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			kar.InventoryRes.ClearCurrentSlot()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyK) {
			kar.InventoryRes.RandomFillAllSlots()
		}

	}
}

func (ui *UI) Draw() {
	if kar.ECWorld.Alive(kar.CurrentPlayer) {

		// Draw hotbar background
		kar.ColorMDIO.GeoM.Reset()
		kar.ColorMDIO.GeoM.Translate(hotbarPositionX, hotbarPositionY)
		colorm.DrawImage(kar.Screen, res.Hotbar, kar.ColorM, kar.ColorMDIO)

		// Draw Quick-slot numbers
		dotX := float64(kar.InventoryRes.QuickSlot1)*17 + float64(hotbarPositionX) + 7
		kar.TextDO.GeoM.Reset()
		kar.TextDO.GeoM.Translate(dotX, -4)
		text.Draw(kar.Screen, strconv.Itoa(kar.InventoryRes.QuickSlot1+1), res.Font, kar.TextDO)
		dotX = float64(kar.InventoryRes.QuickSlot2)*17 + float64(hotbarPositionX) + 7
		kar.TextDO.GeoM.Reset()
		kar.TextDO.GeoM.Translate(dotX, -4)
		text.Draw(kar.Screen, strconv.Itoa(kar.InventoryRes.QuickSlot2+1), res.Font, kar.TextDO)

		// Draw slots
		for x := range 9 {
			slotID := kar.InventoryRes.Slots[x].ID
			quantity := kar.InventoryRes.Slots[x].Quantity
			SlotOffsetX := float64(x) * 17
			SlotOffsetX += hotbarPositionX

			// draw hotbar item icons
			kar.ColorMDIO.GeoM.Reset()
			kar.ColorMDIO.GeoM.Translate(SlotOffsetX+(5), hotbarPositionY+(5))
			if slotID != items.Air && kar.InventoryRes.Slots[x].Quantity > 0 {
				colorm.DrawImage(kar.Screen, res.Icon8[slotID], kar.ColorM, kar.ColorMDIO)
			}
			// Draw hotbar selected slot border
			if x == kar.InventoryRes.CurrentSlotIndex {
				// border
				kar.ColorMDIO.GeoM.Translate(-5, -5)
				colorm.DrawImage(kar.Screen, res.SlotBorder, kar.ColorM, kar.ColorMDIO)
				// Draw hotbar slot item display name
				if !kar.InventoryRes.IsCurrentSlotEmpty() {
					kar.TextDO.GeoM.Reset()
					kar.TextDO.GeoM.Translate(SlotOffsetX-1, hotbarPositionY+14)
					if items.HasTag(slotID, items.Tool) {
						text.Draw(kar.Screen, fmt.Sprintf(
							"%v\nDurability %v",
							items.Property[slotID].DisplayName,
							kar.InventoryRes.Slots[x].Durability,
						), res.Font, kar.TextDO)
					} else {
						text.Draw(kar.Screen, items.Property[slotID].DisplayName, res.Font, kar.TextDO)
					}
				}
			}

			// Draw item quantity number
			if quantity > 1 && items.IsStackable(slotID) {
				kar.TextDO.GeoM.Reset()
				kar.TextDO.GeoM.Translate(SlotOffsetX+6, hotbarPositionY+4)
				num := strconv.FormatUint(uint64(quantity), 10)
				if quantity < 10 {
					num = " " + num
				}
				text.Draw(kar.Screen, num, res.Font, kar.TextDO)
			}
		}

		// Draw player health text
		kar.TextDO.GeoM.Reset()
		kar.TextDO.GeoM.Translate(hotbarRightEdgePosX+8, hotbarPositionY)
		playerHealth := arc.MapHealth.GetUnchecked(kar.CurrentPlayer)
		text.Draw(kar.Screen, fmt.Sprintf("Health %v", playerHealth.Current), res.Font, kar.TextDO)

		// Draw crafting table
		if kar.GameDataRes.CraftingState {

			// crafting table Background
			kar.ColorMDIO.GeoM.Reset()
			kar.ColorMDIO.GeoM.Translate(craftingTablePositionX, craftingTablePositionY)

			if kar.GameDataRes.CraftingState4 {
				colorm.DrawImage(kar.Screen, res.CraftingTable4, kar.ColorM, kar.ColorMDIO)
			} else {
				colorm.DrawImage(kar.Screen, res.CraftingTable, kar.ColorM, kar.ColorMDIO)
			}

			// draw crafting table item icons
			for x := 0; x < 3; x++ {
				for y := 0; y < 3; y++ {
					if kar.CraftingTableRes.Slots[y][x].ID != items.Air {
						sx := craftingTablePositionX + float64(x*17)
						sy := craftingTablePositionY + float64(y*17)
						kar.ColorMDIO.GeoM.Reset()
						kar.ColorMDIO.GeoM.Translate(sx+6, sy+6)
						colorm.DrawImage(
							kar.Screen,
							res.Icon8[kar.CraftingTableRes.Slots[y][x].ID],
							kar.ColorM,
							kar.ColorMDIO,
						)

						// Draw item quantity number
						quantity := kar.CraftingTableRes.Slots[y][x].Quantity
						if quantity > 1 {
							kar.TextDO.GeoM.Reset()
							kar.TextDO.GeoM.Translate(sx+7, sy+5)
							num := strconv.FormatUint(uint64(quantity), 10)
							if quantity < 10 {
								num = " " + num
							}
							text.Draw(kar.Screen, num, res.Font, kar.TextDO)
						}
					}

					// draw selected slot border of crqfting table
					if x == kar.CraftingTableRes.SlotPosX && y == kar.CraftingTableRes.SlotPosY {
						sx := craftingTablePositionX + float64(x*17)
						sy := craftingTablePositionY + float64(y*17)
						kar.ColorMDIO.GeoM.Reset()
						kar.ColorMDIO.GeoM.Translate(sx+1, sy+1)
						colorm.DrawImage(kar.Screen, res.SlotBorder, kar.ColorM, kar.ColorMDIO)
					}

				}
			}

			// draw crafting table result item icon
			if kar.CraftingTableRes.ResultSlot.ID != 0 {
				kar.ColorMDIO.GeoM.Reset()

				if kar.GameDataRes.CraftingState4 {
					kar.ColorMDIO.GeoM.Translate(craftingTablePositionX+41, craftingTablePositionY+14)
				} else {
					kar.ColorMDIO.GeoM.Translate(craftingTablePositionX+58, craftingTablePositionY+23)
				}

				colorm.DrawImage(kar.Screen, res.Icon8[kar.CraftingTableRes.ResultSlot.ID], kar.ColorM, kar.ColorMDIO)

				// Draw result item quantity number
				quantity := kar.CraftingTableRes.ResultSlot.Quantity
				if quantity > 1 {
					kar.TextDO.GeoM.Reset()
					if kar.GameDataRes.CraftingState4 {
						kar.TextDO.GeoM.Translate(craftingTablePositionX+42, craftingTablePositionY+13)
					} else {
						kar.TextDO.GeoM.Translate(craftingTablePositionX+58, craftingTablePositionY+22)
					}
					num := strconv.FormatUint(uint64(quantity), 10)
					if quantity < 10 {
						num = " " + num
					}
					text.Draw(kar.Screen, num, res.Font, kar.TextDO)
				}
			}

		}

		// Draw debug info
		if kar.DrawDebugTextEnabled {
			_, _, _, _, playerController, pFacing := arc.MapPlayer.Get(kar.CurrentPlayer)
			ebitenutil.DebugPrintAt(kar.Screen, fmt.Sprintf(
				debugInfo,
				playerController.CurrentState,
				pFacing,
			), 10, 50)
		}

	} else {
		// Draw debug info
		ebitenutil.DebugPrintAt(kar.Screen, "YOU ARE DEAD!", int(kar.ScreenW/2)-30, int(kar.ScreenH/2))
	}
}

func onInventorySlotChanged() {
	switch kar.InventoryRes.CurrentSlotID() {
	case items.WoodenAxe:
		kar.PlayerAnimPlayer.CurrentAtlas = "WoodenAxe"
	case items.WoodenPickaxe:
		kar.PlayerAnimPlayer.CurrentAtlas = "WoodenPickaxe"
	case items.WoodenShovel:
		kar.PlayerAnimPlayer.CurrentAtlas = "WoodenShovel"
	case items.StoneAxe:
		kar.PlayerAnimPlayer.CurrentAtlas = "StoneAxe"
	case items.StonePickaxe:
		kar.PlayerAnimPlayer.CurrentAtlas = "StonePickaxe"
	case items.StoneShovel:
		kar.PlayerAnimPlayer.CurrentAtlas = "StoneShovel"
	case items.IronAxe:
		kar.PlayerAnimPlayer.CurrentAtlas = "IronAxe"
	case items.IronPickaxe:
		kar.PlayerAnimPlayer.CurrentAtlas = "IronPickaxe"
	case items.IronShovel:
		kar.PlayerAnimPlayer.CurrentAtlas = "IronShovel"
	case items.DiamondAxe:
		kar.PlayerAnimPlayer.CurrentAtlas = "DiamondAxe"
	case items.DiamondPickaxe:
		kar.PlayerAnimPlayer.CurrentAtlas = "DiamondPickaxe"
	case items.DiamondShovel:
		kar.PlayerAnimPlayer.CurrentAtlas = "DiamondShovel"
	default:
		kar.PlayerAnimPlayer.CurrentAtlas = "Default"
	}
}
