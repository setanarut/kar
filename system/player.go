package system

import (
	"image/color"
	"kar"
	"kar/arche"
	"kar/comp"
	"kar/items"
	"kar/res"
	"kar/types"
	"kar/util"
	"kar/world"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/anim"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

const (
	CooldownTimeSec  = 3.0
	MaxFallSpeed     = 270.0
	MaxFallSpeedCap  = 240.0
	MaxSpeed         = 153.75
	MaxWalkSpeed     = 93.75
	MinSlowDownSpeed = 33.75
	MinSpeed         = 4.453125
	RunAcceleration  = 200.390625
	SkidFriction     = 365.625
	StompSpeed       = 240.0
	StompSpeedCap    = -60.0
	WalkAcceleration = 133.59375
	WalkFriction     = 182.8125
)

var (
	jumpSpeeds        = [3]float64{-240.0, -240.0, -300.0}
	longJumpGravities = [3]float64{450.0, 421.875, 562.5}
	gravities         = [3]float64{1575.0, 1350.0, 2025.0}
	speedThresholds   = [2]float64{60, 138.75}
)

// States
var (
	isAttacking        bool
	isCrouching        bool
	isFacingLeft       bool
	isFacingRight      bool
	isFacingUp         bool
	isFacingLeftLast   bool
	isFacingRightLast  bool
	isFacingDown       bool
	isFalling          bool
	isIdle             bool
	isOnFloor          bool
	isRunning          bool
	isSkiding          bool
	isDigDown, isDigUp bool
	// isJumping    bool
)

// var speedScale = 0.0
var (
	minSpeedTemp = MinSpeed
	maxSpeedTemp = MaxWalkSpeed
	acceleration = WalkAcceleration
	delta        = 1 / 60.0

	speedThreshold int = 0
)

type Player struct{}

func (sys *Player) Init() {
	playerBody.SetVelocityUpdateFunc(playerDefaultVelocityFunc)
}
func (sys *Player) Draw() {}
func (sys *Player) Update() {

	ProcessInput()

	if playerEntry.Valid() {

		playerPixelCoord = world.WorldPosToPixelCoord(playerBody.Position())
		playerAnimation := comp.AnimPlayer.Get(playerEntry)
		playerDrawOptions := comp.DrawOptions.Get(playerEntry)
		attackSegEnd = playerBody.Position().Add(inputAxisLast.Scale(kar.BlockSize * 3.5))
		hitShape = attackSegQuery.Shape

		if hitShape != nil {
			if checkEntry(hitShape.Body) {
				e := getEntry(hitShape.Body)
				if e.HasComponent(comp.Block) {
					hitItemID = comp.Item.Get(e).ID
				} else {
					hitItemID = items.Air
				}
			}
			hitBlockPos = hitShape.Body.Position()
			hitBlockPixelCoord = world.WorldPosToPixelCoord(hitBlockPos)
			placeBlockPos = hitBlockPos.Add(attackSegQuery.Normal.Scale(kar.BlockSize))
			placeBlockPixelCoord = world.WorldPosToPixelCoord(placeBlockPos)
		}

		attackSegQuery = cmSpace.SegmentQueryFirst(
			playerPos,
			attackSegEnd,
			0,
			filterPlayerRaycast)

		// Fly Mode
		if justPressed(ebiten.KeyG) {
			toggleFlyMode()
		}

		if justReleased(ebiten.KeyShiftRight) {
			ResetHitBlockHealth()
		}
		if attackSegQuery.Shape == nil || attackSegQuery.Shape != hitShape {
			ResetHitBlockHealth()
		}

		// Give damage to block
		if pressed(ebiten.KeyShiftRight) {
			GiveDamageToBlock()
		}

		// Place block
		if justPressed(ebiten.KeySlash) {
			PlaceBlock()
		}

		// Eğer boş slot varsa eline al
		if justPressed(ebiten.KeyE) {
			TakeInHand()
		}
		// Drop Item
		if justPressed(ebiten.KeyQ) {
			DropSlotItem()
		}

		// Adds random items to inventory
		if justPressed(ebiten.KeyR) {
			RandomFillInventory()

		}
		if justPressed(ebiten.KeyTab) {
			GoToNextSlot()
		}
		if justPressed(ebiten.Key0) {
			deleteSlot(inventory, selectedSlotIndex)
		}

		UpdateSlotInput()
		UpdateAnimationStates(playerAnimation, playerDrawOptions)

	}

	UpdateFunctionKeys()
}

func onFloor() bool {
	groundNormal := vec.Vec2{}
	playerBody.EachArbiter(func(arb *cm.Arbiter) {
		n := arb.Normal().Neg()
		if n.Y < groundNormal.Y {
			groundNormal = n
		}
	})

	isOnFloor = groundNormal.Y < 0
	return isOnFloor
}

func ProcessInput() {
	isOnFloor = onFloor()
	isAttacking = pressed(ebiten.KeyShiftRight)
	isIdle = inputAxis.Equal(zero) && !isAttacking && isOnFloor
	isFacingDown = inputAxisLast.Equal(down) || inputAxis.Equal(down)
	isFacingUp = inputAxisLast.Equal(up) || inputAxis.Equal(up)
	isFacingRightLast = inputAxisLast.Equal(right)
	isFacingLeftLast = inputAxisLast.Equal(left)
	isFacingRight = inputAxisLast.Equal(right) || inputAxis.Equal(right)
	isDigDown = isFacingDown && isAttacking
	isDigUp = isFacingUp && isAttacking

	if isOnFloor {
		isRunning = ebiten.IsKeyPressed(ebiten.KeyShiftLeft)
		isCrouching = ebiten.IsKeyPressed(ebiten.KeyDown)
		if isCrouching && inputAxis.X != 0 {
			isCrouching = false
			inputAxis.X = 0.0
		}
	}
}

func UpdateAnimationStates(anim *anim.AnimationPlayer, opt *types.DrawOptions) {

	if isIdle && isFacingLeftLast {
		anim.SetState("idle_left")
		opt.FlipX = false
	} else if isIdle && isFacingRightLast {
		anim.SetState("idle_right")
		opt.FlipX = false
	} else {
		anim.SetState("idle_front")
		opt.FlipX = false
	}

	if !isOnFloor && !isIdle {
		anim.SetState("jump")
	}

	if isDigDown {
		anim.SetState("dig_down")
	}
	if isDigUp {
		anim.SetState("dig_right")
	}

	if isAttacking && isFacingRight && !isIdle {
		anim.SetState("dig_right")
		opt.FlipX = false
	}
	if isAttacking && isFacingLeft && !isIdle {
		anim.SetState("dig_right")
		opt.FlipX = true
	}

	if inputAxis.Equal(right) && !isAttacking && isOnFloor && !isIdle {
		anim.SetState("walk_right")
		opt.FlipX = false
	}
	if inputAxis.Equal(left) && !isAttacking && isOnFloor && !isIdle {
		anim.SetState("walk_right")
		opt.FlipX = true
	}
}

func toggleFlyMode() {
	playerFlyModeDisabled = !playerFlyModeDisabled
	switch playerFlyModeDisabled {
	case true:
		playerBody.Shapes[0].SetSensor(false)
		playerBody.SetVelocityUpdateFunc(playerDefaultVelocityFunc)
	case false:
		playerBody.Shapes[0].SetSensor(true)
		playerBody.SetVelocityUpdateFunc(playerFlyVelocityFunc)
	}
}

func playerFlyVelocityFunc(b *cm.Body, _ vec.Vec2, _, _ float64) {
	velocity := inputAxis.Unit().Scale(300)
	b.SetVelocityVector(velocity)
}

func playerDefaultVelocityFunc(body *cm.Body, grav vec.Vec2, damping, dt float64) {
	velocity := playerBody.Velocity()

	// process_jump()
	if isOnFloor {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			// isJumping = true

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
		if ebiten.IsKeyPressed(ebiten.KeySpace) && !isFalling {
			gravity = longJumpGravities[speedThreshold]
		}
		velocity.Y = velocity.Y + gravity*delta
		if velocity.Y > MaxFallSpeed {
			velocity.Y = MaxFallSpeedCap
		}
	}

	if velocity.Y > 0 {
		// isJumping = false
		isFalling = true
	} else if isOnFloor {
		isFalling = false
	}

	// process_walk()
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

	// speedScale = math.Abs(velocity.X) / MaxSpeed

	playerBody.SetVelocityVector(velocity)
}

func UpdateFunctionKeys() {
	if justPressed(ebiten.KeyF2) {
		debugDrawingEnabled = !debugDrawingEnabled
	}
	if justPressed(ebiten.KeyF4) {
		go util.WritePNG(
			res.Frames[items.Dirt][0],
			desktopDir+"map.png",
		)
	}
	if justPressed(ebiten.KeyF6) {
		go util.WritePNG(
			world.ApplyColorMap(gameWorld.Image, items.ItemColorMap),
			desktopDir+"map.png",
		)
	}
	if justPressed(ebiten.KeyF5) {
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
	if justPressed(ebiten.Key1) {
		selectedSlotIndex = 0
	}
	if justPressed(ebiten.Key2) {
		selectedSlotIndex = 1
	}
	if justPressed(ebiten.Key3) {
		selectedSlotIndex = 2
	}
	if justPressed(ebiten.Key4) {
		selectedSlotIndex = 3
	}
	if justPressed(ebiten.Key5) {
		selectedSlotIndex = 4
	}
	if justPressed(ebiten.Key6) {
		selectedSlotIndex = 5
	}
	if justPressed(ebiten.Key7) {
		selectedSlotIndex = 6
	}
	if justPressed(ebiten.Key8) {
		selectedSlotIndex = 7
	}
	if justPressed(ebiten.Key9) {
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
			uint16(util.RandRangeInt(1, len(items.Property))),
		)
		inventory.Slots[i].Quantity = uint8(util.RandRangeInt(1, 64))
	}
}
func DropSlotItem() {
	id := inventory.Slots[selectedSlotIndex].ID
	if inventory.Slots[selectedSlotIndex].Quantity > 0 {
		inventory.Slots[selectedSlotIndex].Quantity--
		e := arche.SpawnDropItem(cmSpace, ecsWorld, playerPos, id)
		b := comp.Body.Get(e)
		if isFacingLeft {
			b.ApplyImpulseAtLocalPoint(
				inputAxisLast.Scale(200).Rotate(util.Radians(45)), vec.Vec2{})
		}
		if isFacingRight {
			b.ApplyImpulseAtLocalPoint(
				inputAxisLast.Scale(200).Rotate(util.Radians(-45)), vec.Vec2{})
		}

	}
}

func ResetHitBlockHealth() {
	if hitShape != nil {
		if checkEntry(hitShape.Body) {
			e := getEntry(hitShape.Body)
			if e.HasComponent(comp.Block) && e.HasComponent(comp.Health) {
				resetHealthComponent(e)
			}
		}
	}
}
func GiveDamageToBlock() {
	if hitShape != nil {
		if checkEntry(hitShape.Body) {
			e := getEntry(hitShape.Body)
			if e.HasComponent(comp.Breakable) && e.HasComponent(comp.Health) {
				h := comp.Health.Get(e)
				h.Health -= 0.2
			}
		}
	}
}

func PlaceBlock() {
	if hitShape != nil {
		if items.IsBlock(inventory.Slots[selectedSlotIndex].ID) {
			if inventory.Slots[selectedSlotIndex].ID != items.Air {
				if playerPixelCoord != placeBlockPixelCoord {
					if removeItem(inventory, inventory.Slots[selectedSlotIndex].ID) {
						arche.SpawnBlock(
							cmSpace,
							ecsWorld,
							placeBlockPos,
							inventory.Slots[selectedSlotIndex].ID,
						)
						gameWorld.Image.SetGray16(
							placeBlockPixelCoord.X,
							placeBlockPixelCoord.Y,
							color.Gray16{selectedSlotItemID},
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

func MoveToward(from, to, delta float64) float64 {
	if math.Abs(to-from) <= delta {
		return to
	}
	if to > from {
		return from + delta
	}
	return from - delta
}
