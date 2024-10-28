package system

import (
	"fmt"
	"image/color"
	"kar"
	"kar/arche"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/items"
	"kar/res"
	"kar/types"
	"kar/world"

	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/anim"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
	"github.com/yohamta/donburi"
)

type vec2 = vec.Vec2

var jumptime = 0.0
var damp = vec2{0, -1000}

type Player struct {
}

func (plr *Player) Init() {
}
func (plr *Player) Draw() {

}
func (plr *Player) Update() {

	UpdateWASDInput()

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
	noWASD = wasd.Equal(zero)
	walkRight = wasd.Equal(right)
	walkLeft = wasd.Equal(left)
	attacking = pressed(eb.KeyShiftRight)
	walking = walkLeft || walkRight
	idle = noWASD && !attacking && isGround
	digDown = facingDown && attacking
	digUp = facingUp && attacking

	comp.TagWASD.Each(ecsWorld, MovementFunc)
	comp.TagWASDFly.Each(ecsWorld, MovementFlyFunc)

	if playerEntry.Valid() {
		playerPixelCoord = world.WorldToPixel(playerPos)
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
			hitBlockPos = hitShape.Body.Position()
			hitBlockPixelCoord = world.WorldToPixel(hitBlockPos)
			placeBlockPos = hitBlockPos.Add(attackSegQuery.Normal.Scale(kar.BlockSize))
			placeBlockPixelCoord = world.WorldToPixel(placeBlockPos)
		}

		attackSegQuery = cmSpace.SegmentQueryFirst(
			playerPos,
			attackSegEnd,
			0,
			filterPlayerRaycast)

		// Fly Mode
		if justPressed(eb.KeyG) {
			CheckFlyMode(playerEntry, playerBody)
		}

		if justReleased(eb.KeyShiftRight) {
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
			deleteSlot(inventory, selectedSlotIndex)
		}

		UpdateSlotInput()
		UpdateAnimationStates(playerAnimation, playerDrawOptions)

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
func UpdateAnimationStates(anim *anim.AnimationPlayer, opt *types.DrawOptions) {
	if idle && facingLeft && !walking {
		anim.SetState("idle_left")
		opt.FlipX = false
	} else if idle && facingRight {
		anim.SetState("idle_right")
		opt.FlipX = false
	} else {
		anim.SetState("idle_front")
		opt.FlipX = false
	}

	if !isGround && !idle {
		anim.SetState("jump")
	}

	if digDown {
		anim.SetState("dig_down")
	}
	if digUp {
		anim.SetState("dig_right")
	}

	if attacking && facingRight && !idle {
		anim.SetState("dig_right")
		opt.FlipX = false
	}
	if attacking && facingLeft && !idle {
		anim.SetState("dig_right")
		opt.FlipX = true
	}

	if walkRight && !attacking && isGround && !idle {
		anim.SetState("walk_right")
		opt.FlipX = false
	}
	if walkLeft && !attacking && isGround && !idle {
		anim.SetState("walk_right")
		opt.FlipX = true
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
	if selectedSlotIndex+1 < len(inventory.Slots) {
		selectedSlotIndex++
	} else {
		selectedSlotIndex = 0
	}
}
func RandomFillInventory() {
	resetInventory(inventory)
	for i := range inventory.Slots {
		addItem(
			inventory,
			uint16(mathutil.RandRangeInt(1, len(items.Property))),
		)
		inventory.Slots[i].Quantity = uint8(mathutil.RandRangeInt(1, 64))
	}
}
func DropSlotItem() {
	id := inventory.Slots[selectedSlotIndex].ID
	if inventory.Slots[selectedSlotIndex].Quantity > 0 {
		inventory.Slots[selectedSlotIndex].Quantity--
		e := arche.SpawnDropItem(cmSpace, ecsWorld, playerPos, id)
		b := comp.Body.Get(e)
		if facingLeft {
			b.ApplyImpulseAtLocalPoint(
				wasdLast.Scale(200).Rotate(mathutil.Radians(45)), vec2{})
		}
		if facingRight {
			b.ApplyImpulseAtLocalPoint(
				wasdLast.Scale(200).Rotate(mathutil.Radians(-45)), vec2{})
		}

	}
}

func isOnFloor() bool {
	groundNormal := vec2{}
	playerBody.EachArbiter(func(arb *cm.Arbiter) {
		n := arb.Normal().Neg()
		if n.Y < groundNormal.Y {
			groundNormal = n
		}
	})

	return groundNormal.Y < 0
}

func MovementFunc(e *donburi.Entry) {
	body := comp.Body.Get(e)
	// p := body.Position()
	// queryInfo := cmSpace.SegmentQueryFirst(
	// 	p,
	// 	p.Add(vect{0, kar.BlockSize / 2}),
	// 	0,
	// 	filterPlayerRaycast,
	// )
	// contactShape := queryInfo.Shape
	speed := kar.BlockSize * 30
	bv := body.Velocity()
	body.SetVelocity(bv.X*0.9, bv.Y)
	// yerde
	if isOnFloor() {
		isGround = true
		// Zıpla
		if justPressed(eb.KeySpace) {
			body.ApplyImpulseAtLocalPoint(
				vec2{0, -(speed * 0.30)},
				body.CenterOfGravity(),
			)
		}
		if walkLeft {
			body.ApplyForceAtLocalPoint(vec2{-speed, 0}, body.CenterOfGravity())
		}
		if walkRight {
			body.ApplyForceAtLocalPoint(vec2{speed, 0}, body.CenterOfGravity())
		}
	} else {
		isGround = false
		if walkLeft {
			body.ApplyForceAtLocalPoint(vec2{-(speed), 0}, body.CenterOfGravity())
		}
		if walkRight {
			body.ApplyForceAtLocalPoint(vec2{speed, 0}, body.CenterOfGravity())
		}
	}
}

func MovementFunc2(e *donburi.Entry) {
	body := comp.Body.Get(e)
	p := body.Position()
	queryInfo := cmSpace.SegmentQueryFirst(
		p,
		p.Add(vec2{0, kar.BlockSize / 2}),
		0,
		filterPlayerRaycast,
	)
	contactShape := queryInfo.Shape
	speed := kar.BlockSize * 30
	vel := body.Velocity()
	body.SetVelocity(vel.X*0.9, vel.Y)
	if justReleased(eb.KeySpace) {
		jumptime = 0
	}
	if pressed(eb.KeySpace) {
		jumptime += 0.01
		if body.Velocity().Y > -200 {
			if jumptime < 1.0 {
				body.ApplyForceAtLocalPoint(damp.Scale(-jumptime), body.CenterOfGravity())
			}
			if isGround {
				body.ApplyImpulseAtLocalPoint(vec2{0, -300}, body.CenterOfGravity())
			}
		}
	}

	// yerde
	if contactShape != nil {
		isGround = true
		// Zıpla
		if justPressed(eb.KeySpace) {
			body.ApplyImpulseAtLocalPoint(vec2{0, -(speed * 0.30)}, body.CenterOfGravity())
		}

		if walkLeft {
			body.ApplyForceAtLocalPoint(vec2{-speed, 0}, body.CenterOfGravity())
		}
		if walkRight {
			body.ApplyForceAtLocalPoint(vec2{speed, 0}, body.CenterOfGravity())
		}
	} else {
		isGround = false
		if walkLeft {
			body.ApplyForceAtLocalPoint(vec2{-(speed), 0}, body.CenterOfGravity())
		}
		if walkRight {
			body.ApplyForceAtLocalPoint(vec2{speed, 0}, body.CenterOfGravity())
		}
	}
}

func MovementFlyFunc(e *donburi.Entry) {
	body := comp.Body.Get(e)
	mobileData := comp.Mobile.Get(e)
	velocity := wasd.Unit().Scale(mobileData.Speed * 4)
	body.SetVelocityVector(
		body.Velocity().LerpDistance(velocity, mobileData.Accel),
	)
}
func CheckFlyMode(player *donburi.Entry, playerBody *cm.Body) {
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
func UpdateWASDInput() {
	wasd = vec2{}
	if pressed(eb.KeyW) {
		wasd.Y -= 1
	}
	if pressed(eb.KeyS) {
		wasd.Y += 1
	}
	if pressed(eb.KeyA) {
		wasd.X -= 1
	}
	if pressed(eb.KeyD) {
		wasd.X += 1
	}
	if !wasd.Equal(vec2{}) {
		wasdLast = wasd
	}
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
				h := comp.Health.Get(e)
				h.Health -= 0.2
			}
		}
	}
}
func PlaceBlock() {
	if hitShape != nil {
		id := inventory.Slots[selectedSlotIndex].ID
		if items.IsBlock(id) {
			if id != items.Air {
				if playerPixelCoord != placeBlockPixelCoord {
					if removeItem(inventory, inventory.Slots[selectedSlotIndex].ID) {
						arche.SpawnBlock(
							cmSpace,
							ecsWorld,
							placeBlockPos,
							inventory.Slots[selectedSlotIndex].ID,
						)
						fmt.Println(placeBlockPixelCoord, DisplayName(id))
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
	if slotIndex, ok := hasEmptySlot(inventory); ok {
		temp := inventory.HandSlot
		inventory.HandSlot = inventory.Slots[selectedSlotIndex]
		deleteSlot(inventory, selectedSlotIndex)
		inventory.Slots[slotIndex] = temp
	}
}
