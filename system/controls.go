package system

import (
	"image/color"
	"kar/arche"
	"kar/comp"
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

var HitShape *cm.Shape
var AttackSegmentQuery cm.SegmentQueryInfo
var BlockPlaceTimerData *types.DataTimer

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
	BlockPlaceTimerData = &types.DataTimer{
		TimerDuration: time.Second / 8,
	}
}

func (sys *PlayerControlSystem) Update() {
	TimerUpdate(BlockPlaceTimerData)
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

		AttackSegmentQuery = res.Space.SegmentQueryFirst(
			playerPosition,
			playerPosition.Add(res.Input.LastPressedWASDDirection.Scale(res.BlockSize*3.5)),
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

		if res.Input.IsPressedAndNotABC(ebiten.KeyArrowUp, ebiten.KeyArrowDown, ebiten.KeyArrowLeft, ebiten.KeyArrowRight) {

			if TimerIsReady(BlockPlaceTimerData) {
				TimerReset(BlockPlaceTimerData)
			}

			if TimerIsStart(BlockPlaceTimerData) {
				end := playerPosition.Add(res.Up.Scale(res.BlockSize * 3.5))
				// end.X = end.X + res.BlockSize

				seg := res.Space.SegmentQueryFirst(
					playerPosition,
					end,
					0,
					res.FilterPlayerRaycast)

				if seg.Shape != nil {
					r := playerBody.FirstShape().Class.(*cm.Circle).Radius()
					dist := seg.Point.Distance(playerPosition) - r
					if dist > res.BlockSize {
						centerDistance := seg.Normal.Normalize().Scale(res.BlockSize / 2)
						blockPos := seg.Point.Add(centerDistance)
						mapPos := vec.FromPoint(Terr.WorldSpaceToMapSpace(blockPos))
						blockPosCenter := mapPos.Scale(res.BlockSize)

						air := color.Gray{uint8(items.Air)}
						if res.Terrain.GrayAt(int(mapPos.X), int(mapPos.Y)) == air {
							arche.SpawnBlock(blockPosCenter, Terr.WorldPosToChunkCoord(blockPosCenter), items.Dirt)
							res.Terrain.SetGray(int(mapPos.X), int(mapPos.Y), color.Gray{uint8(items.Dirt)})
						}

					}

				}
			}

		}

		if res.Input.IsPressedAndNotABC(ebiten.KeyArrowDown, ebiten.KeyArrowUp, ebiten.KeyArrowLeft, ebiten.KeyArrowRight) {

			if TimerIsReady(BlockPlaceTimerData) {
				TimerReset(BlockPlaceTimerData)
			}

			if TimerIsStart(BlockPlaceTimerData) {
				seg := res.Space.SegmentQueryFirst(
					playerPosition,
					playerPosition.Add(res.Down.Scale(res.BlockSize*3.5)),
					0,
					res.FilterPlayerRaycast)

				if seg.Shape != nil {
					r := playerBody.FirstShape().Class.(*cm.Circle).Radius()
					dist := seg.Point.Distance(playerPosition) - r
					if dist > res.BlockSize {
						centerDistance := seg.Normal.Normalize().Scale(res.BlockSize / 2)
						blockPos := seg.Point.Add(centerDistance)
						mapPos := vec.FromPoint(Terr.WorldSpaceToMapSpace(blockPos))
						blockPosCenter := mapPos.Scale(res.BlockSize)

						air := color.Gray{uint8(items.Air)}
						if res.Terrain.GrayAt(int(mapPos.X), int(mapPos.Y)) == air {
							arche.SpawnBlock(blockPosCenter, Terr.WorldPosToChunkCoord(blockPosCenter), items.Dirt)
							res.Terrain.SetGray(int(mapPos.X), int(mapPos.Y), color.Gray{uint8(items.Dirt)})
						}
					}

				}
			}

		}
		if res.Input.IsPressedAndNotABC(ebiten.KeyArrowLeft, ebiten.KeyArrowRight, ebiten.KeyArrowUp, ebiten.KeyArrowDown) {

			if TimerIsReady(BlockPlaceTimerData) {
				TimerReset(BlockPlaceTimerData)
			}

			if TimerIsStart(BlockPlaceTimerData) {
				seg := res.Space.SegmentQueryFirst(
					playerPosition,
					playerPosition.Add(res.Left.Scale(res.BlockSize*3.5)),
					0,
					res.FilterPlayerRaycast)

				if seg.Shape != nil {
					r := playerBody.FirstShape().Class.(*cm.Circle).Radius()
					dist := seg.Point.Distance(playerPosition) - r
					if dist > res.BlockSize {
						centerDistance := seg.Normal.Normalize().Scale(res.BlockSize / 2)
						blockPos := seg.Point.Add(centerDistance)
						mapPos := vec.FromPoint(Terr.WorldSpaceToMapSpace(blockPos))
						blockPosCenter := mapPos.Scale(res.BlockSize)

						air := color.Gray{uint8(items.Air)}
						if res.Terrain.GrayAt(int(mapPos.X), int(mapPos.Y)) == air {
							arche.SpawnBlock(blockPosCenter, Terr.WorldPosToChunkCoord(blockPosCenter), items.Dirt)
							res.Terrain.SetGray(int(mapPos.X), int(mapPos.Y), color.Gray{uint8(items.Dirt)})
						}
					}

				}
			}

		}
		if res.Input.IsPressedAndNotABC(ebiten.KeyArrowRight, ebiten.KeyArrowDown, ebiten.KeyArrowLeft, ebiten.KeyArrowUp) {

			if TimerIsReady(BlockPlaceTimerData) {
				TimerReset(BlockPlaceTimerData)
			}

			if TimerIsStart(BlockPlaceTimerData) {
				seg := res.Space.SegmentQueryFirst(
					playerPosition,
					playerPosition.Add(res.Right.Scale(res.BlockSize*3.5)),
					0,
					res.FilterPlayerRaycast)

				if seg.Shape != nil {
					r := playerBody.FirstShape().Class.(*cm.Circle).Radius()
					dist := seg.Point.Distance(playerPosition) - r
					if dist > res.BlockSize {
						centerDistance := seg.Normal.Normalize().Scale(res.BlockSize / 2)
						blockPos := seg.Point.Add(centerDistance)
						mapPos := vec.FromPoint(Terr.WorldSpaceToMapSpace(blockPos))
						blockPosCenter := mapPos.Scale(res.BlockSize)

						air := color.Gray{uint8(items.Air)}
						if res.Terrain.GrayAt(int(mapPos.X), int(mapPos.Y)) == air {
							arche.SpawnBlock(blockPosCenter, Terr.WorldPosToChunkCoord(blockPosCenter), items.Dirt)
							res.Terrain.SetGray(int(mapPos.X), int(mapPos.Y), color.Gray{uint8(items.Dirt)})
						}
					}

				}
			}

		}

		if res.Input.ArrowDirection.Equal(vec.Vec2{}) {
			TimerReset(BlockPlaceTimerData)
		}

		// Reset block health
		if inpututil.IsKeyJustReleased(ebiten.KeyShiftRight) {
			if HitShape != nil {
				if CheckEntry(HitShape.Body()) {
					e := GetEntry(HitShape.Body())
					if e.HasComponent(comp.Item) && e.HasComponent(comp.Health) {
						ResetHealthComponent(e)
					}
				}
			}
		}

		// reset block health
		if AttackSegmentQuery.Shape == nil || AttackSegmentQuery.Shape != HitShape {
			if HitShape != nil {
				if CheckEntry(HitShape.Body()) {
					e := GetEntry(HitShape.Body())
					if e.HasComponent(comp.Item) && e.HasComponent(comp.Health) {
						ResetHealthComponent(e)
					}
				}
			}
		}

		// Attack
		if ebiten.IsKeyPressed(ebiten.KeyShiftRight) {
			if AttackSegmentQuery.Shape != nil && AttackSegmentQuery.Shape == HitShape {
				if HitShape != nil {
					if CheckEntry(HitShape.Body()) {
						e := GetEntry(HitShape.Body())
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

		if WalkRight && !Attacking {
			playerAnimation.SetState("walk_right")
			playerDrawOptions.FlipX = false
		}
		if WalkLeft && !Attacking {
			playerAnimation.SetState("walk_right")
			playerDrawOptions.FlipX = true
		}

		HitShape = AttackSegmentQuery.Shape
	}

}

func (sys *PlayerControlSystem) Draw(screen *ebiten.Image) {}

func WASDPlatformerForce(e *donburi.Entry) {

	body := comp.Body.Get(e)
	p := body.Position()
	// p.Add(vec.Vec2{0, -25}
	queryInfo := res.Space.SegmentQueryFirst(p, p.Add(vec.Vec2{0, res.BlockSize / 2}), 0, res.FilterPlayerRaycast)
	// queryInfoRight := res.Space.SegmentQueryFirst(p, p.Add(vec.Vec2{0, -25}), 0, res.FilterPlayerRaycast)
	contactShape := queryInfo.Shape
	speed := res.BlockSize * 30

	bv := body.Velocity()
	body.SetVelocity(bv.X*0.9, bv.Y)
	// body.SetVelocityVector(body.Velocity().ClampLenght(500))
	// if bv.X > res.BlockSize*5 {
	// 	body.SetVelocity(500, bv.Y)
	// }
	// if bv.X < -(res.BlockSize * 5) {
	// 	body.SetVelocity(-500, bv.Y)
	// }
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
func WASDPlatformer(e *donburi.Entry) {
	vel := vec.Vec2{}
	body := comp.Body.Get(e)
	p := body.Position()
	mobileData := comp.Mobile.Get(e)
	vel.X = res.Input.WASDDirection.X
	vel = vel.Scale(mobileData.Speed)
	vel = vel.Add(vec.Vec2{0, -500})
	body.SetVelocityVector(body.Velocity().LerpDistance(vel, mobileData.Accel))

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {

		queryInfo := res.Space.SegmentQueryFirst(p, p.Add(vec.Vec2{0, -50}), 0, res.FilterPlayerRaycast)
		contactShape := queryInfo.Shape

		if contactShape != nil {
			body.ApplyImpulseAtLocalPoint(vec.Vec2{0, 900}, body.CenterOfGravity())
		}

	}

}

func WASD4Directional(e *donburi.Entry) {
	body := comp.Body.Get(e)
	mobileData := comp.Mobile.Get(e)
	velocity := res.Input.WASDDirection.Normalize().Scale(mobileData.Speed)
	body.SetVelocityVector(body.Velocity().LerpDistance(velocity, mobileData.Accel))
}
func WASDFly(e *donburi.Entry) {
	body := comp.Body.Get(e)
	mobileData := comp.Mobile.Get(e)
	velocity := res.Input.WASDDirection.Normalize().Scale(mobileData.Speed)
	body.SetVelocityVector(body.Velocity().LerpDistance(velocity, mobileData.Accel))
}
