package system

import (
	"image"
	"image/color"
	"kar/arche"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/engine/util"
	"kar/items"
	"kar/resources"
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

	resources.Input.UpdateWASDDirection()

	if inpututil.IsKeyJustPressed(ebiten.KeyY) {
		ChipmunkDebugSpaceDrawing = !ChipmunkDebugSpaceDrawing
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		resources.Input.LastPressedWASDDirection = resources.Left
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		resources.Input.LastPressedWASDDirection = resources.Right
	}

	FacingRight = resources.Input.LastPressedWASDDirection.Equal(resources.Right) || resources.Input.WASDDirection.Equal(resources.Right)
	FacingLeft = resources.Input.LastPressedWASDDirection.Equal(resources.Left) || resources.Input.WASDDirection.Equal(resources.Left)
	FacingDown = resources.Input.LastPressedWASDDirection.Equal(resources.Down) || resources.Input.WASDDirection.Equal(resources.Down)
	FacingUp = resources.Input.LastPressedWASDDirection.Equal(resources.Up) || resources.Input.WASDDirection.Equal(resources.Up)
	NoWASD = resources.Input.WASDDirection.Equal(resources.Zero)
	WalkRight = resources.Input.WASDDirection.Equal(resources.Right)
	WalkLeft = resources.Input.WASDDirection.Equal(resources.Left)
	Attacking = ebiten.IsKeyPressed(ebiten.KeyShiftRight)
	Walking = WalkLeft || WalkRight
	Idle = NoWASD && !Attacking && IsGround
	DigDown = FacingDown && Attacking
	DigUp = FacingUp && Attacking
	IdleAttack = NoWASD && Attacking && IsGround

	comp.TagWASD.Each(resources.ECSWorld, WASDPlatformerForce)
	comp.TagWASDFly.Each(resources.ECSWorld, WASDFly)

	if player, ok := comp.TagPlayer.First(resources.ECSWorld); ok {
		playerInventory := comp.Inventory.Get(player)
		playerBody := comp.Body.Get(player)
		playerPos = playerBody.Position()
		playerPosMap = MainWorld.WorldSpaceToPixelSpace(playerPos.Add(vec.Vec2{(resources.BlockSize / 2), (resources.BlockSize / 2)}))

		playerAnimation := comp.AnimPlayer.Get(player)
		playerDrawOptions := comp.DrawOptions.Get(player)
		attackSegmentEnd = playerPos.Add(resources.Input.LastPressedWASDDirection.Scale(resources.BlockSize * 3.5))
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
			placeBlockPos = currentBlockPos.Add(attackSegmentQuery.Normal.Scale(resources.BlockSize))
			placeBlockPosMap = MainWorld.WorldSpaceToPixelSpace(placeBlockPos)
		}

		attackSegmentQuery = resources.Space.SegmentQueryFirst(
			playerPos,
			attackSegmentEnd,
			0,
			resources.FilterPlayerRaycast)

		// Fly Mode
		if inpututil.IsKeyJustPressed(ebiten.KeyG) {
			if player.HasComponent(comp.TagWASD) {
				player.RemoveComponent(comp.TagWASD)
				player.AddComponent(comp.TagWASDFly)
				playerBody.SetVelocity(0, 0)
				playerBody.FirstShape().SetSensor(true)
			} else {
				playerBody.SetVelocity(0, 0)
				player.RemoveComponent(comp.TagWASDFly)
				player.AddComponent(comp.TagWASD)
				playerBody.FirstShape().SetSensor(false)
			}
		}

		// resourceset block health
		if inpututil.IsKeyJustReleased(ebiten.KeyShiftRight) {
			if hitShape != nil {
				if checkEntry(hitShape.Body()) {
					e := getEntry(hitShape.Body())
					if e.HasComponent(comp.Item) && e.HasComponent(comp.Health) {
						resourcesetHealthComponent(e)
					}
				}
			}
		}

		// resourceset block health
		if attackSegmentQuery.Shape == nil || attackSegmentQuery.Shape != hitShape {
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
		if ebiten.IsKeyPressed(ebiten.KeyShiftRight) {

			if attackSegmentQuery.Shape != nil && attackSegmentQuery.Shape == hitShape {
				if hitShape != nil {

					if checkEntry(hitShape.Body()) {
						e := getEntry(hitShape.Body())
						if e.HasComponent(comp.Item) && e.HasComponent(comp.TagBlock) && e.HasComponent(comp.Health) {
							it := comp.Item.Get(e)
							if items.Property[it.ID].Breakable {
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
			selectedSlotItemId := playerInventory.Slots[resources.SelectedSlot].ID
			if items.Property[selectedSlotItemId].Category&items.CategoryBlock != 0 {
				if hitShape != nil {
					if selectedSlotItemId != items.Air {
						if playerPosMap != placeBlockPosMap {
							if inventoryManager.removeItem(playerInventory, playerInventory.Slots[resources.SelectedSlot].ID) {
								arche.SpawnBlock(placeBlockPos, placeBlockPosMap, playerInventory.Slots[resources.SelectedSlot].ID)
								MainWorld.Image.SetGray16(placeBlockPosMap.X, placeBlockPosMap.Y, color.Gray16{resources.SelectedItemID})
							}
						}
					}
				}
			}

		}

		// Drop Item
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			if hitShape != nil {
				id := playerInventory.Slots[resources.SelectedSlot].ID
				if playerInventory.Slots[resources.SelectedSlot].Quantity > 0 {
					playerInventory.Slots[resources.SelectedSlot].Quantity--
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

			if resources.SelectedSlot+1 < len(playerInventory.Slots) {
				resources.SelectedSlot++
			} else {
				resources.SelectedSlot = 0
			}

		}
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			resources.SelectedSlot = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			resources.SelectedSlot = 1
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			resources.SelectedSlot = 2
		}
		if inpututil.IsKeyJustPressed(ebiten.Key4) {
			resources.SelectedSlot = 3
		}
		if inpututil.IsKeyJustPressed(ebiten.Key5) {
			resources.SelectedSlot = 4
		}
		if inpututil.IsKeyJustPressed(ebiten.Key6) {
			resources.SelectedSlot = 5
		}
		if inpututil.IsKeyJustPressed(ebiten.Key7) {
			resources.SelectedSlot = 6
		}
		if inpututil.IsKeyJustPressed(ebiten.Key8) {
			resources.SelectedSlot = 7
		}
		if inpututil.IsKeyJustPressed(ebiten.Key9) {
			resources.SelectedSlot = 8
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
		go util.WritePNG(resources.SpriteStages[items.Dirt][0], resources.DesktopDir+"map.png")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF6) {
		go util.WritePNG(world.ApplyColorMap(MainWorld.Image, items.ItemColorMap), resources.DesktopDir+"map.png")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		go util.WritePNG(
			world.ApplyColorMap(MainWorld.ChunkImage(PlayerChunk), items.ItemColorMap),
			resources.DesktopDir+"playerChunk.png")
	}

}

func (sys *PlayerControlSystem) Draw(screen *ebiten.Image) {}

func WASDPlatformerForce(e *donburi.Entry) {
	body := comp.Body.Get(e)
	p := body.Position()
	queryInfo := resources.Space.SegmentQueryFirst(p, p.Add(vec.Vec2{0, resources.BlockSize / 2}), 0, resources.FilterPlayerRaycast)
	contactShape := queryInfo.Shape
	speed := resources.BlockSize * 30
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
	velocity := resources.Input.WASDDirection.Unit().Scale(mobileData.Speed * 4)
	body.SetVelocityVector(body.Velocity().LerpDistance(velocity, mobileData.Accel))
}
