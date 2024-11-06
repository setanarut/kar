package system

import (
	"fmt"
	"image"
	"image/color"
	"kar"
	"kar/arc"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/items"
	"kar/res"
	"kar/world"

	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

type vec2 = vec.Vec2

var (
	attackSegQuery                                             cm.SegmentQueryInfo
	hitShape                                                   *cm.Shape
	placeBlockPos, hitBlockPos, attackSegEnd                   vec2
	playerPixelCoord, placeBlockPixelCoord, hitBlockPixelCoord image.Point
	hitItemID                                                  uint16
)

var (
	filterPlayerRaycast = cm.ShapeFilter{
		0,
		arc.PlayerRayBit,
		cm.AllCategories &^ arc.PlayerBit &^ arc.DropItemBit,
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

	if kar.WorldECS.Alive(playerEntity) {

		playerPixelCoord = world.WorldToPixel(playerPos)
		attackSegEnd = playerPos.Add(inputAxisLast.Scale(kar.BlockSize * 3.5))
		hitShape = attackSegQuery.Shape
		if hitShape != nil {
			e := hitShape.Body.UserData.(ecs.Entity)
			if kar.WorldECS.Alive(e) {
				if hitShape.CollisionType == arc.Block {
					hitItemID = arc.ItemMapper.Get(e).ID
				} else {
					hitItemID = items.Air
				}
				hitBlockPos = hitShape.Body.Position()
				hitBlockPixelCoord = world.WorldToPixel(hitBlockPos)
				placeBlockPos = hitBlockPos.Add(attackSegQuery.Normal.Scale(kar.BlockSize))
				placeBlockPixelCoord = world.WorldToPixel(placeBlockPos)
			}

		}

		attackSegQuery = kar.Space.SegmentQueryFirst(
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

	if justPressed(eb.KeyK) {
		// playerBody.SetVelocity(1, 0)
		// fmt.Println(playerBody.Velocity(), playerBody.Position())
		playerBody.SetPosition(playerBody.Position().Add(vec.Vec2{10, 0}))
		kar.Space.ReindexShape(playerBody.ShapeAtIndex(0))
	}
	if justPressed(eb.KeyJ) {
		// playerBody.SetVelocity(1, 0)
		// fmt.Println(playerBody.Velocity(), playerBody.Position())
		playerBody.SetPosition(playerBody.Position().Add(vec.Vec2{0, 10}))
		kar.Space.ReindexShape(playerBody.ShapeAtIndex(0))
	}

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
		dropItemEntity := arc.SpawnDropItem(playerPos, id)
		bd := arc.BodyMapper.Get(dropItemEntity)
		if IsFacingLeft {
			bd.Body.ApplyImpulseAtLocalPoint(
				inputAxisLast.Scale(200).Rotate(mathutil.Radians(45)), vec2{})
		}
		if IsFacingRight {
			bd.Body.ApplyImpulseAtLocalPoint(
				inputAxisLast.Scale(200).Rotate(mathutil.Radians(-45)), vec2{})
		}

	}
}

func PlayerFlyVelocityFunc(b *cm.Body, _ vec.Vec2, _, _ float64) {
	v := inputAxis.Unit().Scale(300)
	b.SetVelocityVector(v)
}

func ResetHitBlockHealth() {
	if hitShape != nil {
		e := hitShape.Body.UserData.(ecs.Entity)
		if kar.WorldECS.Alive(e) {
			if hitShape.CollisionType == arc.Block {
				h := arc.HealthMapper.Get(e)
				h.Health = h.MaxHealth
			}
		}
	}
}

// BlockQuery := arc.BlockFilter.Query(&kar.WorldECS)
// for BlockQuery.Next() {
// 	healthComponent, _, body, _ := BlockQuery.Get()
// 	if healthComponent.Health <= 0 {
// 		pos := body.Position()
// 		body.Shapes[0].SetSensor(true)
// 		kar.Space.AddPostStepCallback(removeBodyPostStep, body, nil)
// 		s.toRemove = append(s.toRemove, BlockQuery.Entity())
// 		blockPos := world.WorldToPixel(pos)
// 		gameWorld.Image.SetGray16(blockPos.X, blockPos.Y, color.Gray16{items.Air})

// 		// dropID := items.Property[item.ID].Drops
// 		// arc.SpawnDropItem(pos, dropID)
// 	}
// }

func GiveDamageToBlock() {
	if hitShape != nil {
		e := hitShape.Body.UserData.(ecs.Entity)
		if kar.WorldECS.Alive(e) {
			itm := arc.ItemMapper.Get(e)
			if items.IsBreakable(itm.ID) && arc.HealthMapper.Has(e) {
				h := arc.HealthMapper.Get(e)
				if h.Health <= 0 {
					fmt.Println("Blok Yok Edildi!")
					removeBodyPostStep(kar.Space, hitShape.Body, nil)
					kar.WorldECS.RemoveEntity(e)

					gameWorld.Image.SetGray16(
						hitBlockPixelCoord.X,
						hitBlockPixelCoord.Y, color.Gray16{items.Air})

					dropID := items.Property[itm.ID].Drops
					arc.SpawnDropItem(hitBlockPos, dropID)

				}
				h.Health -= 0.2
				isAttacking = true
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
						arc.SpawnBlock(
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
		playerBody.SetVelocityUpdateFunc(PlayerFlyVelocityFunc)
	}
	playerFlyModeDisabled = !playerFlyModeDisabled
}
