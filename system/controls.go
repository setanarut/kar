package system

import (
	"image"
	"image/color"
	"kar/arche"
	"kar/comp"
	"kar/engine"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/items"
	"kar/res"
	"kar/world"

	eb "github.com/hajimehoshi/ebiten/v2"
	iu "github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
	"github.com/yohamta/donburi"
)

var (
	attackSegQuery                                     cm.SegmentQueryInfo
	hitShape                                           *cm.Shape
	playerPos, placeBlockPos                           vec.Vec2
	currentBlockPos, attackSegEnd                      vec.Vec2
	playerPosMap, placeBlockPosMap, currentBlockPosMap image.Point
	hitItemID                                          uint16
)
var (
	justPressed  = iu.IsKeyJustPressed
	justReleased = iu.IsKeyJustReleased
	keyPressed   = eb.IsKeyPressed
	inputManager = &engine.InputManager{}
)
var (
	right = vec.Vec2{1, 0}
	left  = vec.Vec2{-1, 0}
	down  = vec.Vec2{0, 1}
	up    = vec.Vec2{0, -1}
	noDir = vec.Vec2{0, 0}
)
var (
	attacking   bool
	digDown     bool
	digUp       bool
	facingDown  bool
	facingLeft  bool
	facingRight bool
	facingUp    bool
	idle        bool
	isGround    bool
	noWASD      bool
	walking     bool
	walkLeft    bool
	walkRight   bool
)

type PlayerControlSystem struct {
}

func (sys *PlayerControlSystem) Init() {
}
func (sys *PlayerControlSystem) Update() {

	inputManager.UpdateWASDDirection()
	wasdLast := inputManager.LastPressedWASDDir
	wasd := inputManager.WASDDir

	if justPressed(eb.KeyF2) {
		ChipmunkDebugSpaceDrawing = !ChipmunkDebugSpaceDrawing
	}
	if justPressed(eb.KeyLeft) {
		wasdLast = left
	}
	if justPressed(eb.KeyRight) {
		wasdLast = right
	}

	facingRight = wasdLast.Equal(right) || wasd.Equal(right)
	facingLeft = wasdLast.Equal(left) || wasd.Equal(left)
	facingDown = wasdLast.Equal(down) || wasd.Equal(down)
	facingUp = wasdLast.Equal(up) || wasd.Equal(up)
	noWASD = wasd.Equal(noDir)
	walkRight = wasd.Equal(right)
	walkLeft = wasd.Equal(left)
	attacking = keyPressed(eb.KeyShiftRight)
	walking = walkLeft || walkRight
	idle = noWASD && !attacking && isGround
	digDown = facingDown && attacking
	digUp = facingUp && attacking

	comp.TagWASD.Each(res.ECSWorld, WASDPlatformerForce)
	comp.TagWASDFly.Each(res.ECSWorld, WASDFly)

	if player, ok := comp.TagPlayer.First(res.ECSWorld); ok {
		playerInventory := comp.Inventory.Get(player)
		playerBody := comp.Body.Get(player)
		playerPos = playerBody.Position()
		playerPosMap = world.WorldSpaceToPixelSpace(
			playerPos.Add(
				vec.Vec2{(res.BlockSize / 2), (res.BlockSize / 2)},
			),
		)

		playerAnimation := comp.AnimPlayer.Get(player)
		playerDrawOptions := comp.DrawOptions.Get(player)
		attackSegEnd = playerPos.Add(wasdLast.Scale(res.BlockSize * 3.5))
		hitShape = attackSegQuery.Shape

		if hitShape != nil {
			if checkEntry(hitShape.Body()) {
				e := getEntry(hitShape.Body())
				if e.HasComponent(comp.Item) {
					hitItemID = comp.Item.Get(e).ID
				}
			}
			currentBlockPos = hitShape.Body().Position()
			currentBlockPosMap = world.WorldSpaceToPixelSpace(currentBlockPos)
			placeBlockPos = currentBlockPos.Add(
				attackSegQuery.Normal.Scale(res.BlockSize),
			)
			placeBlockPosMap = world.WorldSpaceToPixelSpace(placeBlockPos)
		}

		attackSegQuery = res.Space.SegmentQueryFirst(
			playerPos,
			attackSegEnd,
			0,
			res.FilterPlayerRaycast)

		// Fly Mode
		if justPressed(eb.KeyG) {
			if player.HasComponent(comp.TagWASD) {
				playerBody.SetVelocity(0, 0)
				player.RemoveComponent(comp.TagWASD)
				player.AddComponent(comp.TagWASDFly)
				playerBody.FirstShape().SetSensor(true)
			} else {
				playerBody.SetVelocity(0, 0)
				player.RemoveComponent(comp.TagWASDFly)
				player.AddComponent(comp.TagWASD)
				playerBody.FirstShape().SetSensor(false)
			}
		}

		// Reset block health
		if justReleased(eb.KeyShiftRight) {
			if hitShape != nil {
				if checkEntry(hitShape.Body()) {
					e := getEntry(hitShape.Body())
					if e.HasComponent(comp.Item) && e.HasComponent(comp.Health) {
						resourcesetHealthComponent(e)
					}
				}
			}
		}

		// Reset block health
		if attackSegQuery.Shape == nil || attackSegQuery.Shape != hitShape {
			if hitShape != nil {
				if checkEntry(hitShape.Body()) {
					e := getEntry(hitShape.Body())
					if e.HasComponent(comp.Item) && e.HasComponent(comp.Health) {
						resourcesetHealthComponent(e)
					}
				}
			}
		}

		// Attack
		if keyPressed(eb.KeyShiftRight) {

			if attackSegQuery.Shape != nil &&
				attackSegQuery.Shape == hitShape {
				if hitShape != nil {
					if checkEntry(hitShape.Body()) {
						e := getEntry(hitShape.Body())
						if e.HasComponent(comp.Item) && e.HasComponent(comp.Health) {
							id := comp.Item.Get(e).ID
							if items.IsBreakable(id) && items.IsBlock(id) {
								h := comp.Health.Get(e)
								h.Health -= 0.2
							}
						}
					}
				}
			}

		}

		// Place block
		if justPressed(eb.KeySlash) {
			selectedSlotItemId := playerInventory.Slots[res.SelectedSlotIndex].ID
			if items.Property[selectedSlotItemId].Category&items.Block != 0 {
				if hitShape != nil {
					if selectedSlotItemId != items.Air {
						if playerPosMap != placeBlockPosMap {
							if inventoryManager.removeItem(
								playerInventory,
								playerInventory.Slots[res.SelectedSlotIndex].ID,
							) {
								arche.SpawnBlock(
									placeBlockPos,
									playerInventory.Slots[res.SelectedSlotIndex].ID,
								)
								MainWorld.Image.SetGray16(
									placeBlockPosMap.X,
									placeBlockPosMap.Y,
									color.Gray16{res.SelectedSlotItemID},
								)
							}
						}
					}
				}
			}

		}

		if justPressed(eb.KeyE) {
			i, ok := inventoryManager.hasEmptySlot(playerInventory)

			if ok {
				temp := playerInventory.HandSlot
				playerInventory.HandSlot = playerInventory.Slots[res.SelectedSlotIndex]
				inventoryManager.resetSlot(playerInventory, res.SelectedSlotIndex)
				playerInventory.Slots[i] = temp
			}

		}
		// Drop Item
		if justPressed(eb.KeyQ) {
			id := playerInventory.Slots[res.SelectedSlotIndex].ID
			if playerInventory.Slots[res.SelectedSlotIndex].Quantity > 0 {
				playerInventory.Slots[res.SelectedSlotIndex].Quantity--
				e := arche.SpawnDropItem(playerPos, id)
				comp.Body.Get(e).ApplyImpulseAtLocalPoint(
					wasdLast.Scale(150).Rotate(0.166),
					right)

			}
		}

		// Add random item to inventory
		if justPressed(eb.KeyR) {
			inventoryManager.reset(playerInventory)
			for i := range playerInventory.Slots {
				inventoryManager.addItemIfEmpty(
					playerInventory,
					uint16(mathutil.RandRangeInt(1, len(items.Property))),
				)
				playerInventory.Slots[i].Quantity = uint8(mathutil.RandRangeInt(1, 64))
			}

		}
		if justPressed(eb.KeyTab) {

			if res.SelectedSlotIndex+1 < len(playerInventory.Slots) {
				res.SelectedSlotIndex++
			} else {
				res.SelectedSlotIndex = 0
			}

		}

		if justPressed(eb.Key0) {
			inventoryManager.resetSlot(playerInventory, res.SelectedSlotIndex)
		}

		if justPressed(eb.Key1) {
			res.SelectedSlotIndex = 0
		}
		if justPressed(eb.Key2) {
			res.SelectedSlotIndex = 1
		}
		if justPressed(eb.Key3) {
			res.SelectedSlotIndex = 2
		}
		if justPressed(eb.Key4) {
			res.SelectedSlotIndex = 3
		}
		if justPressed(eb.Key5) {
			res.SelectedSlotIndex = 4
		}
		if justPressed(eb.Key6) {
			res.SelectedSlotIndex = 5
		}
		if justPressed(eb.Key7) {
			res.SelectedSlotIndex = 6
		}
		if justPressed(eb.Key8) {
			res.SelectedSlotIndex = 7
		}
		if justPressed(eb.Key9) {
			res.SelectedSlotIndex = 8
		}

		if idle && facingLeft && !walking {
			playerAnimation.SetState("idle_left")
			playerDrawOptions.FlipX = false
		} else if idle && facingRight {
			playerAnimation.SetState("idle_right")
			playerDrawOptions.FlipX = false
		} else {
			playerAnimation.SetState("idle_front")
			playerDrawOptions.FlipX = false
		}

		if !isGround && !idle {
			playerAnimation.SetState("jump")
		}

		if digDown {
			playerAnimation.SetState("dig_down")
		}
		if digUp {
			playerAnimation.SetState("dig_right")
		}

		if attacking && facingRight && !idle {
			playerAnimation.SetState("dig_right")
			playerDrawOptions.FlipX = false
		}
		if attacking && facingLeft && !idle {
			playerAnimation.SetState("dig_right")
			playerDrawOptions.FlipX = true
		}

		if walkRight && !attacking && isGround && !idle {
			playerAnimation.SetState("walk_right")
			playerDrawOptions.FlipX = false
		}
		if walkLeft && !attacking && isGround && !idle {
			playerAnimation.SetState("walk_right")
			playerDrawOptions.FlipX = true
		}

	}

	if justPressed(eb.KeyF4) {
		go util.WritePNG(
			res.Frames[items.Dirt][0],
			res.DesktopDir+"map.png",
		)
	}
	if justPressed(eb.KeyF6) {
		go util.WritePNG(
			world.ApplyColorMap(MainWorld.Image, items.ItemColorMap),
			res.DesktopDir+"map.png",
		)
	}
	if justPressed(eb.KeyF5) {
		go util.WritePNG(
			world.ApplyColorMap(
				MainWorld.ChunkImage(playerChunk),
				items.ItemColorMap,
			),
			res.DesktopDir+"playerChunk.png",
		)
	}

}

func (sys *PlayerControlSystem) Draw(screen *eb.Image) {}

func WASDPlatformerForce(e *donburi.Entry) {
	body := comp.Body.Get(e)
	p := body.Position()
	queryInfo := res.Space.SegmentQueryFirst(
		p,
		p.Add(vec.Vec2{0, res.BlockSize / 2}),
		0,
		res.FilterPlayerRaycast,
	)
	contactShape := queryInfo.Shape
	speed := res.BlockSize * 30
	bv := body.Velocity()
	body.SetVelocity(bv.X*0.9, bv.Y)
	// yerde
	if contactShape != nil {
		isGround = true
		// Zıpla
		if justPressed(eb.KeySpace) {
			body.ApplyImpulseAtLocalPoint(
				vec.Vec2{0, -(speed * 0.30)},
				body.CenterOfGravity(),
			)
		}
		if keyPressed(eb.KeyA) {
			body.ApplyForceAtLocalPoint(vec.Vec2{-speed, 0}, body.CenterOfGravity())
		}
		if keyPressed(eb.KeyD) {
			body.ApplyForceAtLocalPoint(vec.Vec2{speed, 0}, body.CenterOfGravity())
		}
	} else {
		isGround = false
		if keyPressed(eb.KeyA) {
			body.ApplyForceAtLocalPoint(vec.Vec2{-(speed), 0}, body.CenterOfGravity())
		}
		if keyPressed(eb.KeyD) {
			body.ApplyForceAtLocalPoint(vec.Vec2{speed, 0}, body.CenterOfGravity())
		}
	}
}

func WASDFly(e *donburi.Entry) {
	body := comp.Body.Get(e)
	mobileData := comp.Mobile.Get(e)
	velocity := inputManager.WASDDir.Unit().Scale(mobileData.Speed * 4)
	body.SetVelocityVector(
		body.Velocity().LerpDistance(velocity, mobileData.Accel),
	)
}
