package system

import (
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
	"github.com/yohamta/donburi"
)

type PlayerControl struct {
}

func (sys *PlayerControl) Init() {
}
func (sys *PlayerControl) Update() {

	updateWASDDirection()

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

	comp.TagWASD.Each(ecsWorld, wasdFunc)
	comp.TagWASDFly.Each(ecsWorld, wasdFlyFunc)

	if playerEntry.Valid() {
		playerBody := comp.Body.Get(playerEntry)
		playerPosMap = world.WorldSpaceToPixelSpace(playerPos)
		playerAnimation := comp.AnimPlayer.Get(playerEntry)
		playerDrawOptions := comp.DrawOptions.Get(playerEntry)
		attackSegEnd = playerPos.Add(wasdLast.Scale(kar.BlockSize * 3.5))
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
			currentBlockPos = hitShape.Body.Position()
			currentBlockPosMap = world.WorldSpaceToPixelSpace(currentBlockPos)
			placeBlockPos = currentBlockPos.Add(
				attackSegQuery.Normal.Scale(kar.BlockSize),
			)
			placeBlockPosMap = world.WorldSpaceToPixelSpace(placeBlockPos)
		}

		attackSegQuery = space.SegmentQueryFirst(
			playerPos,
			attackSegEnd,
			0,
			filterPlayerRaycast)

		// Fly Mode
		if justPressed(eb.KeyG) {
			checkFlyMode(playerEntry, playerBody)
		}

		// Reset block health
		if justReleased(eb.KeyShiftRight) {
			if hitShape != nil {
				if checkEntry(hitShape.Body) {
					e := getEntry(hitShape.Body)
					if e.HasComponent(comp.Item) && e.HasComponent(comp.Health) {
						resourcesetHealthComponent(e)
					}
				}
			}
		}

		// Reset block health
		if attackSegQuery.Shape == nil || attackSegQuery.Shape != hitShape {
			if hitShape != nil {
				if checkEntry(hitShape.Body) {
					e := getEntry(hitShape.Body)
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
					if checkEntry(hitShape.Body) {
						e := getEntry(hitShape.Body)
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
			selectedSlotItemId := inventory.Slots[selectedSlotIndex].ID
			if items.Property[selectedSlotItemId].Category&items.Block != 0 {
				if hitShape != nil {
					if selectedSlotItemId != items.Air {
						if playerPosMap != placeBlockPosMap {
							if removeItem(
								inventory,
								inventory.Slots[selectedSlotIndex].ID,
							) {
								arche.SpawnBlock(space, ecsWorld,
									placeBlockPos,
									inventory.Slots[selectedSlotIndex].ID,
								)
								mainWorld.Image.SetGray16(
									placeBlockPosMap.X,
									placeBlockPosMap.Y,
									color.Gray16{selectedSlotItemID},
								)
							}
						}
					}
				}
			}

		}

		if justPressed(eb.KeyE) {
			i, ok := hasEmptySlot(inventory)

			if ok {
				temp := inventory.HandSlot
				inventory.HandSlot = inventory.Slots[selectedSlotIndex]
				resetSlot(inventory, selectedSlotIndex)
				inventory.Slots[i] = temp
			}

		}
		// Drop Item
		if justPressed(eb.KeyQ) {
			id := inventory.Slots[selectedSlotIndex].ID
			if inventory.Slots[selectedSlotIndex].Quantity > 0 {
				inventory.Slots[selectedSlotIndex].Quantity--
				e := arche.SpawnDropItem(space, ecsWorld, playerPos, id)
				b := comp.Body.Get(e)
				if facingLeft {
					b.ApplyImpulseAtLocalPoint(
						wasdLast.Scale(200).Rotate(mathutil.Radians(45)), vec.Vec2{})
				}
				if facingRight {
					b.ApplyImpulseAtLocalPoint(
						wasdLast.Scale(200).Rotate(mathutil.Radians(-45)), vec.Vec2{})
				}

			}
		}

		// Add random item to inventory
		if justPressed(eb.KeyR) {
			resetInventory(inventory)
			for i := range inventory.Slots {
				addItem(
					inventory,
					uint16(mathutil.RandRangeInt(1, len(items.Property))),
				)
				inventory.Slots[i].Quantity = uint8(mathutil.RandRangeInt(1, 64))
			}

		}
		if justPressed(eb.KeyTab) {

			if selectedSlotIndex+1 < len(inventory.Slots) {
				selectedSlotIndex++
			} else {
				selectedSlotIndex = 0
			}

		}
		if justPressed(eb.Key0) {
			resetSlot(inventory, selectedSlotIndex)
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

	if justPressed(eb.KeyF2) {
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
			world.ApplyColorMap(mainWorld.Image, items.ItemColorMap),
			desktopDir+"map.png",
		)
	}
	if justPressed(eb.KeyF5) {
		go util.WritePNG(
			world.ApplyColorMap(
				mainWorld.ChunkImage(playerChunk),
				items.ItemColorMap,
			),
			desktopDir+"playerChunk.png",
		)
	}

}

func (sys *PlayerControl) Draw() {}

func wasdFunc(e *donburi.Entry) {
	body := comp.Body.Get(e)
	p := body.Position()
	queryInfo := space.SegmentQueryFirst(
		p,
		p.Add(vec.Vec2{0, kar.BlockSize / 2}),
		0,
		filterPlayerRaycast,
	)
	contactShape := queryInfo.Shape
	speed := kar.BlockSize * 30
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
		if walkLeft {
			body.ApplyForceAtLocalPoint(vec.Vec2{-speed, 0}, body.CenterOfGravity())
		}
		if walkRight {
			body.ApplyForceAtLocalPoint(vec.Vec2{speed, 0}, body.CenterOfGravity())
		}
	} else {
		isGround = false
		if walkLeft {
			body.ApplyForceAtLocalPoint(vec.Vec2{-(speed), 0}, body.CenterOfGravity())
		}
		if walkRight {
			body.ApplyForceAtLocalPoint(vec.Vec2{speed, 0}, body.CenterOfGravity())
		}
	}
}

func wasdFlyFunc(e *donburi.Entry) {
	body := comp.Body.Get(e)
	mobileData := comp.Mobile.Get(e)
	velocity := wasd.Unit().Scale(mobileData.Speed * 4)
	body.SetVelocityVector(
		body.Velocity().LerpDistance(velocity, mobileData.Accel),
	)
}

func checkFlyMode(player *donburi.Entry, playerBody *cm.Body) {
	if player.HasComponent(comp.TagWASD) {
		playerBody.SetVelocity(0, 0)
		player.RemoveComponent(comp.TagWASD)
		player.AddComponent(comp.TagWASDFly)
		playerBody.Shapes[0].SetSensor(true)
	} else {
		playerBody.SetVelocity(0, 0)
		player.RemoveComponent(comp.TagWASDFly)
		player.AddComponent(comp.TagWASD)
		playerBody.Shapes[0].SetSensor(false)
	}
}

func updateWASDDirection() {
	wasd = vec.Vec2{}
	if keyPressed(eb.KeyW) {
		wasd.Y -= 1
	}
	if keyPressed(eb.KeyS) {
		wasd.Y += 1
	}
	if keyPressed(eb.KeyA) {
		wasd.X -= 1
	}
	if keyPressed(eb.KeyD) {
		wasd.X += 1
	}
	if !wasd.Equal(vec.Vec2{}) {
		wasdLast = wasd
	}
}
