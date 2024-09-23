package system

import (
	"image/color"
	"kar/arche"
	"kar/comp"
	"kar/engine/mathutil"
	"kar/items"
	"kar/res"
	"kar/types"
	"time"

	"github.com/setanarut/cm"

	"github.com/setanarut/vec"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

var (
	attackSegmentQuery       cm.SegmentQueryInfo
	blockSpawnSegmentQuery   cm.SegmentQueryInfo
	hitShape                 *cm.Shape
	blockPlaceTimerData      *types.DataTimer
	attackSegmentEnd         vec.Vec2
	blockSpawnSegmentEnd     vec.Vec2
	placeBlockPos            vec.Vec2
	blockPlaceTimerDataReady bool
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
	blockPlaceTimerDataReady = false
	blockPlaceTimerData = &types.DataTimer{
		TimerDuration: time.Second / 6,
	}
}

func (sys *PlayerControlSystem) Update() {

	TimerUpdate(blockPlaceTimerData)
	res.Input.UpdateWASDDirection()
	res.Input.UpdateArrowDirection()

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

	comp.WASDTag.Each(res.World, WASDPlatformerForce)
	comp.WASDFlyTag.Each(res.World, WASDFly)

	if player, ok := comp.PlayerTag.First(res.World); ok {

		playerBody := comp.Body.Get(player)
		playerAnimation := comp.AnimationPlayer.Get(player)
		playerDrawOptions := comp.DrawOptions.Get(player)
		playerPosition := playerBody.Position()
		attackSegmentEnd = playerPosition.Add(res.Input.LastPressedWASDDirection.Scale(res.BlockSize * 3.5))
		attackSegmentQuery = res.Space.SegmentQueryFirst(
			playerPosition,
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

		// Block place
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			blockPlaceTimerDataReady = true
			blockSpawnSegmentEnd = playerPosition.Add(res.Down.Scale(res.BlockSize * 3.5))

			blockSpawnSegmentQuery = res.Space.SegmentQueryFirst(
				playerPosition,
				blockSpawnSegmentEnd,
				0,
				res.FilterPlayerRaycast)

			if blockSpawnSegmentQuery.Shape != nil {
				r := playerBody.FirstShape().Class.(*cm.Circle).Radius()
				dist := blockSpawnSegmentQuery.Point.Distance(playerPosition) - r

				if dist > res.BlockSize {
					centerDistance := blockSpawnSegmentQuery.Normal.Unit().Scale(res.BlockSize / 2)
					placeBlockPos = blockSpawnSegmentQuery.Point.Add(centerDistance)
					mapPos := mathutil.FromPoint(Terr.WorldSpaceToMapSpace(placeBlockPos))
					blockPosCenter := mapPos.Scale(res.BlockSize)

					air := color.Gray{uint8(items.Air)}

					if res.Terrain.GrayAt(int(mapPos.X), int(mapPos.Y)) == air {
						arche.SpawnBlock(blockPosCenter, Terr.WorldPosToChunkCoord(blockPosCenter), items.Dirt)
						res.Terrain.SetGray(int(mapPos.X), int(mapPos.Y), color.Gray{uint8(items.Dirt)})
					}

				}

			}

		} else {
			blockPlaceTimerDataReady = false
		}

		// Reset block health
		if inpututil.IsKeyJustReleased(ebiten.KeyShiftRight) {
			if hitShape != nil {
				if CheckEntry(hitShape.Body()) {
					e := GetEntry(hitShape.Body())
					if e.HasComponent(comp.Item) && e.HasComponent(comp.Health) {
						ResetHealthComponent(e)
					}
				}
			}
		}

		// reset block health
		if attackSegmentQuery.Shape == nil || attackSegmentQuery.Shape != hitShape {
			if hitShape != nil {
				if CheckEntry(hitShape.Body()) {
					e := GetEntry(hitShape.Body())
					if e.HasComponent(comp.Item) && e.HasComponent(comp.Health) {
						ResetHealthComponent(e)
					}
				}
			}
		}

		// Attack
		if ebiten.IsKeyPressed(ebiten.KeyShiftRight) {
			if attackSegmentQuery.Shape != nil && attackSegmentQuery.Shape == hitShape {
				if hitShape != nil {
					if CheckEntry(hitShape.Body()) {
						e := GetEntry(hitShape.Body())
						if e.HasComponent(comp.Item) {
							h := comp.Health.Get(e)
							h.Health -= 0.2
						}
					}
				}
			}
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

		hitShape = attackSegmentQuery.Shape
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
	velocity := res.Input.WASDDirection.Unit().Scale(mobileData.Speed)
	body.SetVelocityVector(body.Velocity().LerpDistance(velocity, mobileData.Accel))
}
