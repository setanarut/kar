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
	"math"

	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/anim"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
	"github.com/yohamta/donburi"
)

const velScale = 2.0

const (
	CooldownTimeSec  = 3.0 * velScale
	MaxFallSpeed     = 270.0 * velScale
	MaxFallSpeedCap  = 240.0 * velScale
	MaxSpeed         = 153.75 * velScale
	MaxWalkSpeed     = 93.75 * velScale
	MinSlowDownSpeed = 33.75 * velScale
	MinSpeed         = 4.453125 * velScale
	RunAcceleration  = 200.390625 * velScale
	SkidFriction     = 365.625 * velScale
	StompSpeed       = 240.0 * velScale
	StompSpeedCap    = -60.0 * velScale
	WalkAcceleration = 133.59375 * velScale
	WalkFriction     = 182.8125 * velScale
)

var (
	jumpSpeeds        = [3]float64{-240.0 * velScale, -240.0 * velScale, -300.0 * velScale}
	longJumpGravities = [3]float64{450.0 * velScale, 421.875 * velScale, 562.5 * velScale}
	gravities         = [3]float64{1575.0 * velScale, 1350.0 * velScale, 2025.0 * velScale}
	speedThresholds   = [2]float64{60 * velScale, 138.75 * velScale}
)

// States
var (
	isAttacking                   bool
	isCrouching                   bool
	isFacingLeft, isFacingRight   bool
	isFacingUp, isFacingDown      bool
	isFalling                     bool
	isIdle                        bool
	isOnFloor                     bool
	isRunning                     bool
	isSkiding                     bool
	isDigDown, isDigUp            bool
	isWalkingLeft, isWalkingRight bool
	// isJumping    bool
)
var playerFlyModeDisabled = true

// var speedScale = 0.0
var (
	minSpeedTemp = MinSpeed
	maxSpeedTemp = MaxWalkSpeed
	acceleration = WalkAcceleration
	delta        = 1 / 60.0

	speedThreshold int = 0
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
	playerEntry         *donburi.Entry
	playerVel           vec2
	playerSpawnPos      vec2
	playerBody          *cm.Body
	playerInv           *types.Inventory
	filterPlayerRaycast = cm.ShapeFilter{
		0,
		arche.PlayerRayBit,
		cm.AllCategories &^ arche.PlayerBit &^ arche.DropItemBit,
	}
)

type Player struct {
}

func (plr *Player) Init() {
	// if playerEntry.Valid() {
	// 	playerBody.SetVelocityUpdateFunc(playerDefaultVelocityFunc)
	// }
}
func (plr *Player) Draw() {

}
func (plr *Player) Update() {
	ProcessInput()

	if playerEntry.Valid() {
		playerPixelCoord = world.WorldToPixel(playerPos)
		playerAnimation := comp.AnimPlayer.Get(playerEntry)
		playerDrawOptions := comp.DrawOptions.Get(playerEntry)
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
		if isFacingLeft {
			b.ApplyImpulseAtLocalPoint(
				inputAxisLast.Scale(200).Rotate(mathutil.Radians(45)), vec2{})
		}
		if isFacingRight {
			b.ApplyImpulseAtLocalPoint(
				inputAxisLast.Scale(200).Rotate(mathutil.Radians(-45)), vec2{})
		}

	}
}

func onFloor() bool {
	groundNormal := vec2{}
	playerBody.EachArbiter(func(arb *cm.Arbiter) {
		n := arb.Normal().Neg()
		if n.Y < groundNormal.Y {
			groundNormal = n
		}
	})
	return groundNormal.Y < 0
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
							Space,
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

func playerDefaultVelocityFunc(body *cm.Body, grav vec.Vec2, damping, dt float64) {
	velocity := playerBody.Velocity()
	isOnFloor = onFloor()

	if isOnFloor {
		if justPressed(eb.KeySpace) {
			var speed = math.Abs(velocity.X)
			speedThreshold = len(speedThresholds)

			for i := 0; i < len(speedThresholds); i++ {
				if speed < speedThresholds[i] {
					speedThreshold = i
					break
				}
			}
			velocity.Y = jumpSpeeds[speedThreshold]

		}
	} else {
		var gravity = gravities[speedThreshold]
		if pressed(eb.KeySpace) && !isFalling {
			gravity = longJumpGravities[speedThreshold]
		}
		velocity.Y = velocity.Y + gravity*delta
		if velocity.Y > MaxFallSpeed {
			velocity.Y = MaxFallSpeedCap
		}
	}

	if velocity.Y > 0 {
		isFalling = true
	} else if isOnFloor {
		isFalling = false
	}

	if inputAxis.X != 0 {
		if isOnFloor {
			if velocity.X != 0 {
				isFacingLeft = inputAxis.X < 0.0
				isSkiding = velocity.X < 0.0 != isFacingLeft
			}
			if isSkiding {
				minSpeedTemp = MinSlowDownSpeed
				maxSpeedTemp = MaxWalkSpeed
				acceleration = SkidFriction
			} else if isRunning {
				minSpeedTemp = MinSpeed
				maxSpeedTemp = MaxSpeed
				acceleration = RunAcceleration
			} else {
				minSpeedTemp = MinSpeed
				maxSpeedTemp = MaxWalkSpeed
				acceleration = WalkAcceleration
			}
		} else if isRunning && math.Abs(velocity.X) > MaxWalkSpeed {
			maxSpeedTemp = MaxSpeed
		} else {
			maxSpeedTemp = MaxWalkSpeed
		}
		var target_speed = inputAxis.X * maxSpeedTemp
		velocity.X = MoveToward(velocity.X, target_speed, acceleration*delta)
	} else if isOnFloor && velocity.X != 0 {
		if !isSkiding {
			acceleration = WalkFriction
		}
		if inputAxis.Y != 0 {
			minSpeedTemp = MinSlowDownSpeed
		} else {
			minSpeedTemp = MinSpeed
		}
		if math.Abs(velocity.X) < minSpeedTemp {
			velocity.X = 0.0
		} else {
			velocity.X = MoveToward(velocity.X, 0.0, acceleration*delta)
		}
	}
	if math.Abs(velocity.X) < MinSlowDownSpeed {
		isSkiding = false
	}
	playerBody.SetVelocityVector(velocity)
}

func MoveToward(from, to, delta float64) float64 {
	if math.Abs(to-from) <= delta {
		return to
	}
	if to > from {
		return from + delta
	}
	return from - delta
}

func toggleFlyMode() {
	switch playerFlyModeDisabled {
	case false:
		playerBody.Shapes[0].SetSensor(false)
		playerBody.SetVelocityUpdateFunc(playerDefaultVelocityFunc)
	case true:
		playerBody.Shapes[0].SetSensor(true)
		playerBody.SetVelocityUpdateFunc(playerFlyVelocityFunc)
	}
	playerFlyModeDisabled = !playerFlyModeDisabled
}

func ProcessInput() {
	isOnFloor = onFloor()
	isAttacking = pressed(eb.KeyShiftRight)
	isIdle = inputAxis.Equal(zero) && !isAttacking && isOnFloor

	isFacingDown = inputAxisLast.Equal(down) || inputAxis.Equal(down)
	isFacingUp = inputAxisLast.Equal(up) || inputAxis.Equal(up)
	isFacingRight = inputAxisLast.Equal(right) || inputAxis.Equal(right)
	isFacingLeft = inputAxisLast.Equal(left) || inputAxis.Equal(left)

	isWalkingLeft = !isIdle && isOnFloor && playerBody.Velocity().X < 0.0
	isWalkingRight = !isIdle && isOnFloor && playerBody.Velocity().X > 0.0

	isDigDown = isFacingDown && isAttacking
	isDigUp = isFacingUp && isAttacking

	if isOnFloor {
		isRunning = pressed(eb.KeyAltRight)
		isCrouching = pressed(eb.KeyDown)
		if isCrouching && inputAxis.X != 0 {
			isCrouching = false
			inputAxis.X = 0.0
		}
	}
}

func UpdateAnimationStates(anim *anim.AnimationPlayer, opt *types.DrawOptions) {
	if isRunning {
		anim.SetStateFPS("walk_right", 25)
	} else {
		anim.SetStateFPS("walk_right", 15)
	}
	if isIdle && inputAxisLast.Equal(left) {
		anim.SetState("idle_left")
		opt.FlipX = false
	} else if isIdle && inputAxisLast.Equal(right) {
		anim.SetState("idle_right")
		opt.FlipX = false
	} else {
		anim.SetState("idle_front")
		opt.FlipX = false
	}
	if !isOnFloor && !isIdle {
		anim.SetState("jump")
		if isFacingLeft {
			opt.FlipX = true
		} else {
			opt.FlipX = false
		}
	}
	if isDigDown && !isFacingLeft {
		anim.SetState("dig_down")
	}
	if isDigUp {
		anim.SetState("dig_right")
	}
	// dig right
	if isAttacking && isFacingRight {
		anim.SetState("dig_right")
		opt.FlipX = false
	}
	// dig left
	if isAttacking && isFacingLeft {
		anim.SetState("dig_right")
		opt.FlipX = true
	}
	if isWalkingRight {
		anim.SetState("walk_right")
		opt.FlipX = false
	}
	if isWalkingLeft {
		anim.SetState("walk_right")
		opt.FlipX = true
	}
}
