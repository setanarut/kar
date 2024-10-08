package system

import (
	"image"
	"image/color"
	"kar/arche"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/itm"
	"kar/res"
	"kar/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
	"github.com/yohamta/donburi"
)

var (
	attackSegmentQuery                                 cm.SegmentQueryInfo
	hitShape                                           *cm.Shape
	attackSegmentEnd                                   vec.Vec2
	playerPos, placeBlockPos, currentBlockPos          vec.Vec2
	playerPosMap, placeBlockPosMap, currentBlockPosMap image.Point
	HitItemID                                          uint16
)

var (
	FacingLeft  bool
	FacingRight bool
	FacingDown  bool
	FacingUp    bool
	DigUp       bool
	IsGround    bool
	Idle        bool
	Walking     bool
	WalkRight   bool
	WalkLeft    bool
	Attacking   bool
	IdleAttack  bool
	NoWASD      bool
	DigDown     bool
)

type PlayerControlSystem struct {
}

func NewPlayerControlSystem() *PlayerControlSystem {
	return &PlayerControlSystem{}
}

func (sys *PlayerControlSystem) Init() {
}

func (sys *PlayerControlSystem) Update() {

	res.Input.UpdateWASDDirection()

	if inpututil.IsKeyJustPressed(ebiten.KeyY) {
		ChipmunkDebugSpaceDrawing = !ChipmunkDebugSpaceDrawing
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		res.Input.LastPressedWASDDirection = res.Left
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		res.Input.LastPressedWASDDirection = res.Right
	}

	FacingRight = res.Input.LastPressedWASDDirection.Equal(res.Right) || res.Input.WASDDirection.Equal(res.Right)
	FacingLeft = res.Input.LastPressedWASDDirection.Equal(res.Left) || res.Input.WASDDirection.Equal(res.Left)
	FacingDown = res.Input.LastPressedWASDDirection.Equal(res.Down) || res.Input.WASDDirection.Equal(res.Down)
	FacingUp = res.Input.LastPressedWASDDirection.Equal(res.Up) || res.Input.WASDDirection.Equal(res.Up)
	NoWASD = res.Input.WASDDirection.Equal(res.Zero)
	WalkRight = res.Input.WASDDirection.Equal(res.Right)
	WalkLeft = res.Input.WASDDirection.Equal(res.Left)
	Attacking = ebiten.IsKeyPressed(ebiten.KeyShiftRight)
	Walking = WalkLeft || WalkRight
	Idle = NoWASD && !Attacking && IsGround
	DigDown = FacingDown && Attacking
	DigUp = FacingUp && Attacking
	IdleAttack = NoWASD && Attacking && IsGround

	comp.WASDTag.Each(res.ECSWorld, WASDPlatformerForce)
	comp.WASDFlyTag.Each(res.ECSWorld, WASDFly)

	if player, ok := comp.PlayerTag.First(res.ECSWorld); ok {
		playerInventory := comp.Inventory.Get(player)
		playerBody := comp.Body.Get(player)
		playerPos = playerBody.Position()
		playerPosMap = MainWorld.WorldSpaceToPixelSpace(playerPos.Add(vec.Vec2{(res.BlockSize / 2), (res.BlockSize / 2)}))

		playerAnimation := comp.AnimationPlayer.Get(player)
		playerDrawOptions := comp.DrawOptions.Get(player)
		attackSegmentEnd = playerPos.Add(res.Input.LastPressedWASDDirection.Scale(res.BlockSize * 3.5))
		hitShape = attackSegmentQuery.Shape

		if hitShape != nil {
			if checkEntry(hitShape.Body()) {
				e := getEntry(hitShape.Body())
				if e.HasComponent(comp.Item) {
					HitItemID = comp.Item.Get(e).ID
				}
			}
			currentBlockPos = hitShape.Body().Position()
			currentBlockPosMap = MainWorld.WorldSpaceToPixelSpace(currentBlockPos)
			placeBlockPos = currentBlockPos.Add(attackSegmentQuery.Normal.Scale(res.BlockSize))
			placeBlockPosMap = MainWorld.WorldSpaceToPixelSpace(placeBlockPos)
		}

		attackSegmentQuery = res.Space.SegmentQueryFirst(
			playerPos,
			attackSegmentEnd,
			0,
			res.FilterPlayerRaycast)

		// Fly Mode
		if inpututil.IsKeyJustPressed(ebiten.KeyG) {
			if player.HasComponent(comp.WASDTag) {
				player.RemoveComponent(comp.WASDTag)
				player.AddComponent(comp.WASDFlyTag)
				playerBody.SetVelocity(0, 0)
				playerBody.FirstShape().SetSensor(true)
			} else {
				playerBody.SetVelocity(0, 0)
				player.RemoveComponent(comp.WASDFlyTag)
				player.AddComponent(comp.WASDTag)
				playerBody.FirstShape().SetSensor(false)
			}
		}

		// Reset block health
		if inpututil.IsKeyJustReleased(ebiten.KeyShiftRight) {
			if hitShape != nil {
				if checkEntry(hitShape.Body()) {
					e := getEntry(hitShape.Body())
					if e.HasComponent(comp.Item) && e.HasComponent(comp.Health) {
						resetHealthComponent(e)
					}
				}
			}
		}

		// reset block health
		if attackSegmentQuery.Shape == nil || attackSegmentQuery.Shape != hitShape {
			if hitShape != nil {
				if checkEntry(hitShape.Body()) {
					e := getEntry(hitShape.Body())
					if e.HasComponent(comp.Item) && e.HasComponent(comp.Health) {
						resetHealthComponent(e)
					}
				}
			}
		}

		// Attack
		if ebiten.IsKeyPressed(ebiten.KeyShiftRight) {

			if attackSegmentQuery.Shape != nil && attackSegmentQuery.Shape == hitShape {
				if hitShape != nil {

					if checkEntry(hitShape.Body()) {
						e := getEntry(hitShape.Body())
						if e.HasComponent(comp.Item) && e.HasComponent(comp.BlockTag) && e.HasComponent(comp.Health) {
							it := comp.Item.Get(e)
							if itm.Items[it.ID].Breakable {
								h := comp.Health.Get(e)
								h.Health -= 0.2
							}
						}
					}
				}
			}

		}

		// Place block
		if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
			selectedSlotItemId := playerInventory.Slots[res.SelectedSlot].ID
			if itm.Items[selectedSlotItemId].Category&itm.CatBlock != 0 {
				if hitShape != nil {
					if selectedSlotItemId != itm.Air {
						if playerPosMap != placeBlockPosMap {
							if inventoryManager.removeItem(playerInventory, playerInventory.Slots[res.SelectedSlot].ID) {
								arche.SpawnBlock(placeBlockPos, placeBlockPosMap, playerInventory.Slots[res.SelectedSlot].ID)
								MainWorld.Image.SetGray16(placeBlockPosMap.X, placeBlockPosMap.Y, color.Gray16{res.SelectedItemID})
							}
						}
					}
				}
			}

		}

		// Drop Item
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			if hitShape != nil {
				id := playerInventory.Slots[res.SelectedSlot].ID
				if playerInventory.Slots[res.SelectedSlot].Quantity > 0 {
					playerInventory.Slots[res.SelectedSlot].Quantity--
					arche.SpawnDropItem(placeBlockPos, id, MainWorld.WorldPosToChunkCoord(playerPos))
				}
			}
		}

		// Add random item to inventory
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			if ebiten.IsKeyPressed(ebiten.KeyMetaLeft) {
				inventoryManager.addItem(playerInventory, uint16(mathutil.RandRangeInt(1, 74)))
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyTab) {

			if res.SelectedSlot+1 < len(playerInventory.Slots) {
				res.SelectedSlot++
			} else {
				res.SelectedSlot = 0
			}

		}
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			res.SelectedSlot = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			res.SelectedSlot = 1
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			res.SelectedSlot = 2
		}
		if inpututil.IsKeyJustPressed(ebiten.Key4) {
			res.SelectedSlot = 3
		}
		if inpututil.IsKeyJustPressed(ebiten.Key5) {
			res.SelectedSlot = 4
		}
		if inpututil.IsKeyJustPressed(ebiten.Key6) {
			res.SelectedSlot = 5
		}
		if inpututil.IsKeyJustPressed(ebiten.Key7) {
			res.SelectedSlot = 6
		}
		if inpututil.IsKeyJustPressed(ebiten.Key8) {
			res.SelectedSlot = 7
		}
		if inpututil.IsKeyJustPressed(ebiten.Key9) {
			res.SelectedSlot = 8
		}

		if Idle {
			playerAnimation.SetState("idle")
		}
		if !IsGround {
			playerAnimation.SetState("jump")
		}

		if DigDown {
			playerAnimation.SetState("dig_down")
		}
		if DigUp {
			playerAnimation.SetState("dig_right")
		}

		if Attacking && FacingRight {
			playerAnimation.SetState("dig_right")
			playerDrawOptions.FlipX = false
		}
		if Attacking && FacingLeft {
			playerAnimation.SetState("dig_right")
			playerDrawOptions.FlipX = true
		}

		if WalkRight && !Attacking && IsGround {
			playerAnimation.SetState("walk_right")
			playerDrawOptions.FlipX = false
		}
		if WalkLeft && !Attacking && IsGround {
			playerAnimation.SetState("walk_right")
			playerDrawOptions.FlipX = true
		}

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		go util.WritePNG(res.SpriteFrames[itm.Dirt][0], res.DesktopDir+"map.png")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF6) {
		go util.WritePNG(world.ApplyColorMap(MainWorld.Image, itm.ItemColorMap), res.DesktopDir+"map.png")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		go util.WritePNG(
			world.ApplyColorMap(MainWorld.ChunkImage(PlayerChunk), itm.ItemColorMap),
			res.DesktopDir+"playerChunk.png")
	}

}

func (sys *PlayerControlSystem) Draw(screen *ebiten.Image) {}

func WASDPlatformerForce(e *donburi.Entry) {
	body := comp.Body.Get(e)
	p := body.Position()
	queryInfo := res.Space.SegmentQueryFirst(p, p.Add(vec.Vec2{0, res.BlockSize / 2}), 0, res.FilterPlayerRaycast)
	contactShape := queryInfo.Shape
	speed := res.BlockSize * 30
	bv := body.Velocity()
	body.SetVelocity(bv.X*0.9, bv.Y)
	// yerde
	if contactShape != nil {
		IsGround = true
		// Zıpla
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			body.ApplyImpulseAtLocalPoint(vec.Vec2{0, -(speed * 0.30)}, body.CenterOfGravity())
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			body.ApplyForceAtLocalPoint(vec.Vec2{-speed, 0}, body.CenterOfGravity())
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			body.ApplyForceAtLocalPoint(vec.Vec2{speed, 0}, body.CenterOfGravity())
		}
	} else {
		IsGround = false
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			body.ApplyForceAtLocalPoint(vec.Vec2{-(speed), 0}, body.CenterOfGravity())
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			body.ApplyForceAtLocalPoint(vec.Vec2{speed, 0}, body.CenterOfGravity())
		}
	}
}

func WASDFly(e *donburi.Entry) {
	body := comp.Body.Get(e)
	mobileData := comp.Mobile.Get(e)
	velocity := res.Input.WASDDirection.Unit().Scale(mobileData.Speed * 4)
	body.SetVelocityVector(body.Velocity().LerpDistance(velocity, mobileData.Accel))
}
