package system

import (
	"image"
	"image/color"
	"kar"
	"kar/arche"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/items"
	"kar/res"
	"kar/world"

	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

type vec2 = vec.Vec2

var (
	attackSegQuery                                             cm.SegmentQueryInfo
	hitShape                                                   *cm.Shape
	playerPos, placeBlockPos, hitBlockPos, attackSegEnd        vec2
	playerPixelCoord, placeBlockPixelCoord, hitBlockPixelCoord image.Point
	hitItemID                                                  uint16
)

var (
	filterPlayerRaycast = cm.ShapeFilter{
		0,
		arche.PlayerRayBit,
		cm.AllCategories &^ arche.PlayerBit &^ arche.DropItemBit,
	}
)
var playerFlyModeDisabled = true

type Player struct {
}

func (plr *Player) Init() {
}
func (plr *Player) Draw() {

}
func (plr *Player) Update() {

	if playerEntry.Valid() {
		playerPixelCoord = world.WorldToPixel(playerPos)
		attackSegEnd = playerPos.Add(inputAxisLast.Scale(kar.BlockSize * 3.5))
		hitShape = attackSegQuery.Shape

		if hitShape != nil {
			if checkEntry(hitShape.Body) {
				e := getEntry(hitShape.Body)
				if e.HasComponent(comp.TagBlock) {
					hitItemID = comp.Item.Get(e).ID
				} else {
					hitItemID = items.Air
				}
			}
			hitBlockPos = hitShape.Body.Position()
			hitBlockPixelCoord = world.WorldToPixel(hitBlockPos)
			placeBlockPos = hitBlockPos.Add(attackSegQuery.Normal.Scale(kar.BlockSize))
			placeBlockPixelCoord = world.WorldToPixel(placeBlockPos)
		}

		attackSegQuery = Space.SegmentQueryFirst(
			playerPos,
			attackSegEnd,
			0,
			filterPlayerRaycast)

		// Fly Mode
		if justPressed(eb.KeyG) {
			toggleFlyMode()
		}

		if justReleased(eb.KeyShiftRight) {
			isAttacking = false
			ResetHitBlockHealth()
		}
		if attackSegQuery.Shape == nil || attackSegQuery.Shape != hitShape {
			ResetHitBlockHealth()
		}

		// Give damage to block
		if pressed(eb.KeyShiftRight) {
			GiveDamageToBlock()
		}

		// Place block
		if justPressed(eb.KeySlash) {
			PlaceBlock()
		}

		// Eğer boş slot varsa eline al
		if justPressed(eb.KeyE) {
			TakeInHand()
		}
		// Drop Item
		if justPressed(eb.KeyQ) {
			DropSlotItem()
		}

		// Adds random items to inventory
		if justPressed(eb.KeyR) {
			RandomFillInventory()

		}
		if justPressed(eb.KeyTab) {
			GoToNextSlot()
		}
		if justPressed(eb.Key0) {
			deleteSlot(playerInv, selectedSlotIndex)
		}

		UpdateSlotInput()

	}

	UpdateFunctionKeys()
}
func UpdateFunctionKeys() {
	if justPressed(eb.KeyX) {
		debugDrawingEnabled = !debugDrawingEnabled
	}
	if justPressed(eb.KeyF4) {
		go util.WritePNG(
			res.Frames[items.Dirt][0],
			desktopDir+"map.png",
		)
	}
	if justPressed(eb.KeyF6) {
		go util.WritePNG(
			world.ApplyColorMap(gameWorld.Image, items.ItemColorMap),
			desktopDir+"map.png",
		)
	}
	if justPressed(eb.KeyF5) {
		go util.WritePNG(
			world.ApplyColorMap(
				gameWorld.ChunkImage(gameWorld.PlayerChunk),
				items.ItemColorMap,
			),
			desktopDir+"playerChunk.png",
		)
	}
}

func UpdateSlotInput() {
	if justPressed(eb.Key1) {
		selectedSlotIndex = 0
	}
	if justPressed(eb.Key2) {
		selectedSlotIndex = 1
	}
	if justPressed(eb.Key3) {
		selectedSlotIndex = 2
	}
	if justPressed(eb.Key4) {
		selectedSlotIndex = 3
	}
	if justPressed(eb.Key5) {
		selectedSlotIndex = 4
	}
	if justPressed(eb.Key6) {
		selectedSlotIndex = 5
	}
	if justPressed(eb.Key7) {
		selectedSlotIndex = 6
	}
	if justPressed(eb.Key8) {
		selectedSlotIndex = 7
	}
	if justPressed(eb.Key9) {
		selectedSlotIndex = 8
	}
}
func GoToNextSlot() {
	if selectedSlotIndex+1 < len(playerInv.Slots) {
		selectedSlotIndex++
	} else {
		selectedSlotIndex = 0
	}
}
func RandomFillInventory() {
	resetInventory(playerInv)
	for i := range playerInv.Slots {
		addItem(
			playerInv,
			uint16(mathutil.RandRangeInt(1, len(items.Property))),
		)
		playerInv.Slots[i].Quantity = uint8(mathutil.RandRangeInt(1, 64))
	}
}
func DropSlotItem() {
	id := playerInv.Slots[selectedSlotIndex].ID
	if playerInv.Slots[selectedSlotIndex].Quantity > 0 {
		playerInv.Slots[selectedSlotIndex].Quantity--
		e := arche.SpawnDropItem(Space, ecsWorld, playerPos, id)
		b := comp.Body.Get(e)
		if IsFacingLeft {
			b.ApplyImpulseAtLocalPoint(
				inputAxisLast.Scale(200).Rotate(mathutil.Radians(45)), vec2{})
		}
		if IsFacingRight {
			b.ApplyImpulseAtLocalPoint(
				inputAxisLast.Scale(200).Rotate(mathutil.Radians(-45)), vec2{})
		}

	}
}

func playerFlyVelocityFunc(b *cm.Body, _ vec.Vec2, _, _ float64) {
	velocity := inputAxis.Unit().Scale(300)
	b.SetVelocityVector(velocity)
}

func ResetHitBlockHealth() {
	if hitShape != nil {
		if checkEntry(hitShape.Body) {
			e := getEntry(hitShape.Body)
			if e.HasComponent(comp.TagBlock) && e.HasComponent(comp.Health) {
				resetHealthComponent(e)
			}
		}
	}
}
func GiveDamageToBlock() {
	if hitShape != nil {
		if checkEntry(hitShape.Body) {
			e := getEntry(hitShape.Body)
			if e.HasComponent(comp.TagBreakable) && e.HasComponent(comp.Health) {
				isAttacking = true
				h := comp.Health.Get(e)
				h.Health -= 0.2
			}
		}
	} else {
		isAttacking = false
	}
}
func PlaceBlock() {
	if hitShape != nil {
		id := playerInv.Slots[selectedSlotIndex].ID
		if items.IsBlock(id) {
			if id != items.Air {
				placeBB := cm.NewBBForExtents(placeBlockPos, kar.BlockSize/2, kar.BlockSize/2)
				if !playerBody.ShapeAtIndex(0).BB.Intersects(placeBB) {
					if removeItem(playerInv, id) {
						arche.SpawnBlock(
							Space,
							ecsWorld,
							placeBlockPos,
							playerInv.Slots[selectedSlotIndex].ID,
						)
						gameWorld.Image.SetGray16(
							placeBlockPixelCoord.X,
							placeBlockPixelCoord.Y,
							color.Gray16{id},
						)
					}
				}
			}
		}
	}
}
func TakeInHand() {
	if slotIndex, ok := hasEmptySlot(playerInv); ok {
		temp := playerInv.HandSlot
		playerInv.HandSlot = playerInv.Slots[selectedSlotIndex]
		deleteSlot(playerInv, selectedSlotIndex)
		playerInv.Slots[slotIndex] = temp
	}
}

func toggleFlyMode() {
	switch playerFlyModeDisabled {
	case false:
		playerBody.Shapes[0].SetSensor(false)
		playerBody.SetVelocityUpdateFunc(VelocityFunc)
	case true:
		playerBody.Shapes[0].SetSensor(true)
		playerBody.SetVelocityUpdateFunc(playerFlyVelocityFunc)
	}
	playerFlyModeDisabled = !playerFlyModeDisabled
}
