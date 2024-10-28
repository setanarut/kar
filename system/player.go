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
	"kar/types"
	"kar/world"

	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/anim"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
	"github.com/yohamta/donburi"
)

type vec2 = vec.Vec2

var axisLast vec2
var axis vec2

var (
	attackSegQuery                                             cm.SegmentQueryInfo
	hitShape                                                   *cm.Shape
	playerPos, placeBlockPos, hitBlockPos, attackSegEnd        vec2
	playerPixelCoord, placeBlockPixelCoord, hitBlockPixelCoord image.Point
	hitItemID                                                  uint16
)
var (
	attacking, digDown, digUp, facingDown, facingLeft, facingRight bool
	facingUp, idle, isGround, noWASD, walking, walkLeft, walkRight bool
)
var (
	playerEntry         *donburi.Entry
	playerVel           vec2
	playerSpawnPos      vec2
	playerBody          *cm.Body
	playerInv           *types.Inventory
	jumptime            = 0.0
	damp                = vec2{0, -1000}
	filterPlayerRaycast = cm.ShapeFilter{
		0,
		arche.PlayerRayBit,
		cm.AllCategories &^ arche.PlayerBit &^ arche.DropItemBit,
	}
)

type Player struct {
}

func (plr *Player) Init() {
}
func (plr *Player) Draw() {

}
func (plr *Player) Update() {

	axis = GetAxis()
	if !axis.Equal(vec2{}) {
		axisLast = axis
	}
	if justPressed(eb.KeyLeft) {
		axisLast = left
	}
	if justPressed(eb.KeyRight) {
		axisLast = right
	}

	facingRight = axisLast.Equal(right) || axis.Equal(right)
	facingLeft = axisLast.Equal(left) || axis.Equal(left)
	facingDown = axisLast.Equal(down) || axis.Equal(down)
	facingUp = axisLast.Equal(up) || axis.Equal(up)
	noWASD = axis.Equal(zero)
	walkRight = axis.Equal(right)
	walkLeft = axis.Equal(left)
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
		attackSegEnd = playerPos.Add(axisLast.Scale(kar.BlockSize * 3.5))
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
			deleteSlot(playerInv, selectedSlotIndex)
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
		e := arche.SpawnDropItem(cmSpace, ecsWorld, playerPos, id)
		b := comp.Body.Get(e)
		if facingLeft {
			b.ApplyImpulseAtLocalPoint(
				axisLast.Scale(200).Rotate(mathutil.Radians(45)), vec2{})
		}
		if facingRight {
			b.ApplyImpulseAtLocalPoint(
				axisLast.Scale(200).Rotate(mathutil.Radians(-45)), vec2{})
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
	velocity := axis.Unit().Scale(mobileData.Speed * 4)
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
func GetAxis() vec2 {
	axis := vec2{}
	if pressed(eb.KeyW) {
		axis.Y -= 1
	}
	if pressed(eb.KeyS) {
		axis.Y += 1
	}
	if pressed(eb.KeyA) {
		axis.X -= 1
	}
	if pressed(eb.KeyD) {
		axis.X += 1
	}
	return axis
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
		id := playerInv.HandSlot.ID
		if items.IsBlock(id) {
			if id != items.Air {
				placeBB := cm.NewBBForExtents(placeBlockPos, kar.BlockSize/2, kar.BlockSize/2)
				if !playerBody.ShapeAtIndex(0).BB.Intersects(placeBB) {
					if removeHandItem(playerInv, id) {
						arche.SpawnBlock(
							cmSpace,
							ecsWorld,
							placeBlockPos,
							playerInv.HandSlot.ID,
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
