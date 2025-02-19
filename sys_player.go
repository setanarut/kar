package kar

import (
	"image"
	"image/color"
	"kar/items"
	"kar/v"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	bounceVelocity = -math.Sqrt(2 * SnowballGravity * SnowballBounceHeight)
	blockHealth    float64
	placeTile      image.Point
	isRayHit       bool
	dyingCountdown float64
)

type Player struct {
	itemHit      *HitInfo
	itemBox      AABB
	snowBallBox  AABB
	absXVelocity float64
	isOnFloor    bool
	isSkidding   bool
	playerTile   image.Point
}

func (c *Player) Init() {
	c.itemBox = AABB{Half: DropItemHalfSize}
	c.snowBallBox = AABB{Half: Vec{4, 4}}
	c.itemHit = &HitInfo{}
}

func (c *Player) Update() error {

	if !GameDataRes.CraftingState {
		if ECWorld.Alive(CurrentPlayer) {

			pBox, pVelocity, pHealth, ctrl, pFacing := MapPlayer.Get(CurrentPlayer)
			c.absXVelocity = math.Abs(pVelocity.X)
			// update animation player
			PlayerAnimPlayer.Update()

			if inpututil.IsKeyJustPressed(ebiten.Key9) {
				pHealth.Current = 0
			}

			// Death animation
			if pHealth.Current <= 0 {
				dyingCountdown += 0.1
				PlayerAnimPlayer.Data.Paused = true
				pBox.Pos.Y += dyingCountdown

				if pBox.Pos.Y > CameraRes.Y+CameraRes.Height {
					PlayerAnimPlayer.Data.Paused = false
					PreviousGameState = "playing"
					CurrentGameState = "menu"
					ColorM.ChangeHSV(1, 0, 0.5) // BW
					TextDO.ColorScale.Scale(0.5, 0.5, 0.5, 1)
					ECWorld.RemoveEntity(CurrentPlayer)
					dyingCountdown = 0
				}
				return nil
			}

			// Update input
			isBreakKeyPressed := ebiten.IsKeyPressed(ebiten.KeyRight)
			isRunKeyPressed := ebiten.IsKeyPressed(ebiten.KeyShift)
			isJumpKeyPressed := ebiten.IsKeyPressed(ebiten.KeySpace)
			isAttackKeyJustPressed := inpututil.IsKeyJustPressed(ebiten.KeyLeft)
			isJumpKeyJustPressed := inpututil.IsKeyJustPressed(ebiten.KeySpace)
			inputAxis := Vec{}

			if ebiten.IsKeyPressed(ebiten.KeyW) {
				inputAxis.Y -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				inputAxis.Y += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				inputAxis.X -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyD) {
				inputAxis.X += 1
			}
			if !inputAxis.IsZero() {
				// restrict facing direction to 4 directions (no diagonal)
				switch inputAxis {
				case v.Up, v.Down, v.Left, v.Right:
					*pFacing = Facing(inputAxis)
				default:
					*pFacing = Facing{}
				}
			}

			if math.Abs(pVelocity.X) > 0.01 {
				*pFacing = Facing{math.Copysign(1, pVelocity.X), 0}
			}

			// Update states
			switch ctrl.CurrentState {
			case "idle":
				// enter idle
				if ctrl.PreviousState != "idle" {
					ctrl.PreviousState = ctrl.CurrentState
					if pFacing.Y == 0 {
						PlayerAnimPlayer.SetAnim("idleRight")
					}
					if pFacing.X == 0 {
						PlayerAnimPlayer.SetAnim("idleUp")
					}
				}

				// while idle
				if pFacing.Y == -1 {
					PlayerAnimPlayer.SetAnim("idleUp")
				} else if pFacing.Y == 1 {
					PlayerAnimPlayer.SetAnim("idleDown")
				} else if pFacing.X == 1 {
					PlayerAnimPlayer.SetAnim("idleRight")
				} else if pFacing.X == -1 {
					PlayerAnimPlayer.SetAnim("idleRight")
				}

				// Handle specific transitions
				if isJumpKeyJustPressed {
					if c.absXVelocity > ctrl.MinSpeedThresForJumpBoostMultiplier {
						pVelocity.Y = ctrl.JumpPower * ctrl.JumpBoostMultiplier
					} else {
						pVelocity.Y = ctrl.JumpPower
					}
					ctrl.JumpTimer = 0
					ctrl.CurrentState = "jumping"
				} else if c.isOnFloor && c.absXVelocity > 0.01 {
					if c.absXVelocity > ctrl.MaxWalkSpeed {
						ctrl.CurrentState = "running"
					} else {
						ctrl.CurrentState = "walking"
					}
				} else if !c.isOnFloor && pVelocity.Y > 0.01 {
					ctrl.CurrentState = "falling"
				} else if isBreakKeyPressed && isRayHit {
					ctrl.CurrentState = "breaking"
				} else if pVelocity.Y != 0 && pVelocity.Y < -0.1 {
					ctrl.CurrentState = "jumping"
				}
				// exit idle
				if ctrl.PreviousState != ctrl.CurrentState {
					// fmt.Println("exit idle")
				}
			case "walking":
				// enter walking
				if ctrl.PreviousState != "walking" {
					ctrl.PreviousState = ctrl.CurrentState
					PlayerAnimPlayer.SetAnim("walkRight")
				}
				PlayerAnimPlayer.SetAnimFPS("walkRight", MapRange(c.absXVelocity, 0, ctrl.MaxRunSpeed, 4, 23))

				// Handle specific transitions
				if c.isSkidding {
					ctrl.CurrentState = "skidding"
				} else if pVelocity.Y > 0 && !c.isOnFloor {
					ctrl.CurrentState = "falling"
				} else if isJumpKeyJustPressed {
					ctrl.CurrentState = "jumping"
					if c.absXVelocity > ctrl.MinSpeedThresForJumpBoostMultiplier {
						pVelocity.Y = ctrl.JumpPower * ctrl.JumpBoostMultiplier
					} else {
						pVelocity.Y = ctrl.JumpPower
					}
					ctrl.JumpTimer = 0
				} else if c.absXVelocity == 0 {
					ctrl.CurrentState = "idle"
				} else if c.absXVelocity > ctrl.MaxWalkSpeed {
					ctrl.CurrentState = "running"
				}

				// exit walking
				if ctrl.PreviousState != ctrl.CurrentState {
					// fmt.Println("exit walking")
				}
			case "running":
				// enter running
				if ctrl.PreviousState != "running" {
					ctrl.PreviousState = ctrl.CurrentState
					PlayerAnimPlayer.SetAnim("walkRight")
				}

				// while running
				PlayerAnimPlayer.SetAnimFPS("walkRight", MapRange(c.absXVelocity, 0, ctrl.MaxRunSpeed, 4, 23))

				// Handle specific transitions
				if c.isSkidding {
					ctrl.CurrentState = "skidding"
				} else if pVelocity.Y > 0 && !c.isOnFloor {
					ctrl.CurrentState = "falling"
				} else if isJumpKeyJustPressed {
					if c.absXVelocity > ctrl.MinSpeedThresForJumpBoostMultiplier {
						pVelocity.Y = ctrl.JumpPower * ctrl.JumpBoostMultiplier
					} else {
						pVelocity.Y = ctrl.JumpPower
					}
					ctrl.JumpTimer = 0
					ctrl.CurrentState = "jumping"
				} else if c.absXVelocity < 0.01 {
					ctrl.CurrentState = "idle"
				} else if c.absXVelocity <= ctrl.MaxWalkSpeed {
					ctrl.CurrentState = "walking"
				}
				// exit running
				if ctrl.PreviousState != ctrl.CurrentState {
					// fmt.Println("exit running")
				}
			case "jumping":
				// enter running
				if ctrl.PreviousState != "jumping" {
					ctrl.PreviousState = ctrl.CurrentState
					PlayerAnimPlayer.SetAnim("jump")
				}

				// skidding jumpg physics
				if ctrl.PreviousState == "skidding" {
					if !isJumpKeyPressed && ctrl.JumpTimer < ctrl.JumpReleaseTimer {
						pVelocity.Y = ctrl.ShortJumpVelocity * 0.7 // Kısa zıplama gücünü azalt
						ctrl.JumpTimer = ctrl.JumpHoldTime
					} else if isJumpKeyPressed && ctrl.JumpTimer < ctrl.JumpHoldTime {
						pVelocity.Y += ctrl.JumpBoost * 0.7 // Boost gücünü azalt
						ctrl.JumpTimer++
					} else if pVelocity.Y >= 0.01 {
						ctrl.CurrentState = "falling"
					}
				} else {
					// normal skidding
					if !isJumpKeyPressed && ctrl.JumpTimer < ctrl.JumpReleaseTimer {
						pVelocity.Y = ctrl.ShortJumpVelocity
						ctrl.JumpTimer = ctrl.JumpHoldTime
					} else if isJumpKeyPressed && ctrl.JumpTimer < ctrl.JumpHoldTime {
						speedFactor := (c.absXVelocity / ctrl.MaxRunSpeed) * ctrl.SpeedJumpFactor
						pVelocity.Y += ctrl.JumpBoost * (1 + speedFactor)
						ctrl.JumpTimer++
					} else if pVelocity.Y >= 0 {
						ctrl.CurrentState = "falling"
					}
				}

				// apply air skidding Decel
				if inputAxis.X*pVelocity.X < 0 {
					pVelocity.X += math.Copysign(ctrl.AirSkiddingDecel, -pVelocity.X)
				}

				// exit jumping
				if ctrl.PreviousState != ctrl.CurrentState {
					// fmt.Println("exit jumping")
				}
			case "falling":
				// enter falling
				if ctrl.PreviousState != "falling" {
					ctrl.PreviousState = ctrl.CurrentState
					ctrl.FallingDamageTempPosY = pBox.Pos.Y + pBox.Half.Y
					PlayerAnimPlayer.SetAnim("jump")
				}

				// transitions
				if c.isOnFloor {
					if c.absXVelocity <= 0 {
						ctrl.CurrentState = "idle"
					} else if isRunKeyPressed {
						ctrl.CurrentState = "running"
					} else {
						ctrl.CurrentState = "walking"
					}

				}
				// exit falling
				if ctrl.PreviousState != ctrl.CurrentState {
					// fmt.Println("exit falling")
					d := int(((pBox.Pos.Y + pBox.Half.Y) - ctrl.FallingDamageTempPosY) / 30)
					if d > 3 {
						pHealth.Current -= d - 3
					}
				}
			case "breaking":
				// enter breaking
				if ctrl.PreviousState != "breaking" {
					// fmt.Println("enter breaking")
					ctrl.PreviousState = ctrl.CurrentState
				}
				// update animation states
				if pFacing.X == 1 {
					if c.absXVelocity > 0.01 {
						PlayerAnimPlayer.SetAnim("attackWalk")
					} else {
						PlayerAnimPlayer.SetAnim("attackRight")
					}
				} else if pFacing.X == -1 {
					if c.absXVelocity > 0.01 {
						PlayerAnimPlayer.SetAnim("attackWalk")
					} else {
						PlayerAnimPlayer.SetAnim("attackRight")
					}
				} else if pFacing.Y == 1 {
					PlayerAnimPlayer.SetAnim("attackDown")
				} else if pFacing.Y == -1 {
					PlayerAnimPlayer.SetAnim("attackUp")
				}

				// break block
				if isRayHit {
					blockID := TileMapRes.Get(GameDataRes.TargetBlockCoord.X, GameDataRes.TargetBlockCoord.Y)
					if !items.HasTag(blockID, items.UnbreakableBlock) {
						if items.IsBestTool(blockID, InventoryRes.CurrentSlotID()) {
							blockHealth += PlayerBestToolDamage
						} else {
							blockHealth += PlayerDefaultDamage
						}
					}
					// Destroy block
					if blockHealth >= 180 {

						// set air
						TileMapRes.Set(GameDataRes.TargetBlockCoord.X, GameDataRes.TargetBlockCoord.Y, items.Air)
						blockHealth = 0

						if items.HasTag(InventoryRes.CurrentSlotID(), items.Tool) {
							// damage the tool
							InventoryRes.CurrentSlot().Durability--
							// If durability is 0, destroy the tool.
							if InventoryRes.CurrentSlot().Durability <= 0 {
								InventoryRes.ClearCurrentSlot()
							}
						}
						// Spawn drop item
						x, y := TileMapRes.TileToWorldCenter(GameDataRes.TargetBlockCoord.X, GameDataRes.TargetBlockCoord.Y)
						dropid := items.Property[blockID].DropID
						if blockID == items.OakLeaves {
							if rand.N(2) == 0 {
								dropid = items.OakLeaves
							}
						}
						AppendToSpawnList(x, y, dropid, 0)
					}
				}
				// transitions
				if !isRayHit || (!isBreakKeyPressed && c.isOnFloor) {
					ctrl.CurrentState = "idle"
				} else if !c.isOnFloor && pVelocity.Y > 0.01 {
					ctrl.CurrentState = "falling"
				} else if !isBreakKeyPressed && isJumpKeyJustPressed {
					pVelocity.Y = ctrl.JumpPower
					ctrl.JumpTimer = 0
					ctrl.CurrentState = "jumping"
				}
				// exit breaking
				if ctrl.PreviousState != ctrl.CurrentState {
					// fmt.Println("exit breaking")
					blockHealth = 0
				}
			case "skidding":
				// enter skidding
				if ctrl.PreviousState != "skidding" {
					// fmt.Println("enter skidding")
					ctrl.PreviousState = ctrl.CurrentState
				}
				// Apply Skidding decel
				if c.absXVelocity > ctrl.SkiddingFriction {
					pVelocity.X += math.Copysign(ctrl.SkiddingFriction, -pVelocity.X)
				}
				if c.absXVelocity > 0.5 {
					PlayerAnimPlayer.SetAnim("skidding")
				}

				// Handle specific transitions
				if ctrl.SkiddingJumpEnabled && isJumpKeyJustPressed {
					// Yeni yöne doğru çok küçük sabit değerle başla
					pVelocity.X = math.Copysign(0.3, float64(inputAxis.X))
					pVelocity.Y = ctrl.JumpPower * 0.7 // Zıplama gücünü azalt
					ctrl.JumpTimer = 0
					ctrl.CurrentState = "jumping"
				} else if c.absXVelocity < 0.01 {
					ctrl.CurrentState = "idle"
				} else if !c.isSkidding {
					if c.absXVelocity > ctrl.MaxWalkSpeed {
						ctrl.CurrentState = "running"
					} else {
						ctrl.CurrentState = "walking"
					}
				}
				if ctrl.PreviousState != ctrl.CurrentState {
					// fmt.Println("exit skidding")
				}
			}

			// ########### UPDATE PHYSICS ################

			currentAccel, currentDecel, maxSpeed := ctrl.WalkAcceleration, ctrl.WalkDeceleration, ctrl.MaxWalkSpeed

			pVelocity.Y += ctrl.Gravity
			pVelocity.Y = min(ctrl.MaxFallSpeed, pVelocity.Y)

			if !c.isSkidding {
				if isRunKeyPressed {
					maxSpeed = ctrl.MaxRunSpeed
					currentAccel = ctrl.RunAcceleration
					currentDecel = ctrl.RunDeceleration
				} else if c.absXVelocity > ctrl.MaxWalkSpeed {
					currentDecel = ctrl.RunDeceleration
				}
			}

			if inputAxis.X > 0 {
				if pVelocity.X > maxSpeed {
					pVelocity.X = max(maxSpeed, pVelocity.X-currentDecel)
				} else {
					pVelocity.X = min(maxSpeed, pVelocity.X+currentAccel)
				}
			} else if inputAxis.X < 0 {
				if pVelocity.X < -maxSpeed {
					pVelocity.X = min(-maxSpeed, pVelocity.X+currentDecel)
				} else {
					pVelocity.X = max(-maxSpeed, pVelocity.X-currentAccel)
				}
			} else {
				if pVelocity.X > 0 {
					pVelocity.X = max(0, pVelocity.X-currentDecel)
				} else if pVelocity.X < 0 {
					pVelocity.X = min(0, pVelocity.X+currentDecel)
				}
			}

			c.isSkidding = math.Signbit(pVelocity.X*inputAxis.X) && inputAxis.X != 0

			// Player and tilemap collision
			TileCollider.Collide(*pBox, Vec(*pVelocity), func(collisionInfos []HitTileInfo, delta Vec) {
				c.isOnFloor = false
				pBox.Pos = pBox.Pos.Add(delta)
				// Reset velocity when collide
				for _, ci := range collisionInfos {

					tileID := TileMapRes.GetUnchecked(ci.TileCoords)

					if ci.Normal.Y == -1 {
						// Ground collision
						pVelocity.Y = 0
						c.isOnFloor = true
					}
					// Ceil collision
					if ci.Normal.Y == 1 { // TODO aynı anda olan çarpışmaları teke indir

						pVelocity.Y = 0

						switch tileID {
						case items.StoneBricks:
							// Destroy block when ceil hit
							TileMapRes.Set(ci.TileCoords.X, ci.TileCoords.Y, items.Air)
							wx, wy := TileMapRes.TileToWorldCenter(ci.TileCoords.X, ci.TileCoords.Y)
							SpawnEffect(tileID, wx, wy)
						case items.Random:
							if TileMapRes.Get(ci.TileCoords.X, ci.TileCoords.Y-1) == items.Air {
								TileMapRes.Set(ci.TileCoords.X, ci.TileCoords.Y, items.Bedrock)
								CeilBlockCoord = ci.TileCoords
								CeilBlockTick = 3
							}
						}

					}
					// Right or Left wall collision
					if ci.Normal.X == -1 || ci.Normal.X == 1 {
						// While running at maximum speed, hold down the right arrow key and hit the block to destroy it.
						if c.absXVelocity == ctrl.MaxRunSpeed && isBreakKeyPressed {
							TileMapRes.Set(ci.TileCoords.X, ci.TileCoords.Y, items.Air)
							wx, wy := TileMapRes.TileToWorldCenter(ci.TileCoords.X, ci.TileCoords.Y)
							SpawnEffect(tileID, wx, wy)
						}
						pVelocity.X = 0
						c.absXVelocity = 0
					}
				}
			},
			)

			// player facing raycast for target block
			c.playerTile = TileMapRes.WorldToTile(pBox.Pos.X, pBox.Pos.Y)
			targetBlockTemp := GameDataRes.TargetBlockCoord
			GameDataRes.TargetBlockCoord, isRayHit = TileMapRes.Raycast(
				c.playerTile,
				int(pFacing.X), int(pFacing.Y),
				RaycastDist,
			)

			// reset attack if block focus changed
			if !GameDataRes.TargetBlockCoord.Eq(targetBlockTemp) || !isRayHit {
				blockHealth = 0
			}

			// place block if IsAttackKeyJustPressed
			if isAttackKeyJustPressed {
				anyItemOverlapsWithPlaceRect := false
				// if slot item is block
				if isRayHit && items.HasTag(InventoryRes.CurrentSlot().ID, items.Block) {
					// Get tile rect
					placeTile = GameDataRes.TargetBlockCoord.Sub(image.Point{int(pFacing.X), int(pFacing.Y)})
					placeTileBox := AABB{
						Pos:  Vec{float64(placeTile.X*20) + 10, float64(placeTile.Y*20) + 10},
						Half: Vec{10, 10},
					}
					// check overlaps
					queryItem := FilterDroppedItem.Query(&ECWorld)
					for queryItem.Next() {
						_, itemPos, _, _, _ := queryItem.Get()
						c.itemBox.Pos = Vec(*itemPos)
						anyItemOverlapsWithPlaceRect = placeTileBox.Overlap(c.itemBox, nil)
						if anyItemOverlapsWithPlaceRect {
							queryItem.Close()
							break
						}
					}
					if !anyItemOverlapsWithPlaceRect {
						if !pBox.Overlap(placeTileBox, nil) {
							// place block
							TileMapRes.Set(placeTile.X, placeTile.Y, InventoryRes.CurrentSlotID())
							// remove item
							InventoryRes.RemoveItemFromSelectedSlot()
						}
					}
					// if slot item snowball, throw snowball
				} else if InventoryRes.CurrentSlot().ID == items.Snowball {
					if ctrl.CurrentState != "skidding" {
						switch *pFacing {
						case Facing{1, 0}:
							SpawnProjectile(items.Snowball, pBox.Pos.X, pBox.Pos.Y-4, SnowballSpeedX, SnowballMaxFallVelocity)
						case Facing{-1, 0}:
							SpawnProjectile(items.Snowball, pBox.Pos.X, pBox.Pos.Y-4, -SnowballSpeedX, SnowballMaxFallVelocity)
						}
						InventoryRes.RemoveItemFromSelectedSlot()
					}
				}

			}

			// Drop Item
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
				currentSlot := InventoryRes.CurrentSlot()
				if currentSlot.ID != items.Air {
					AppendToSpawnList(
						pBox.Pos.X,
						pBox.Pos.Y,
						currentSlot.ID,
						currentSlot.Durability,
					)
					InventoryRes.RemoveItemFromSelectedSlot()
					onInventorySlotChanged()
				}
			}
			// projectile physics
			q := FilterProjectile.Query(&ECWorld)
			for q.Next() {
				itemID, projectilePos, projectileVel := q.Get()
				// snowball physics
				if itemID.ID == items.Snowball {
					projectileVel.Y += SnowballGravity
					projectileVel.Y = min(projectileVel.Y, SnowballMaxFallVelocity)
					c.snowBallBox.Pos = Vec(*projectilePos)
					TileCollider.Collide(c.snowBallBox, Vec(*projectileVel), func(ci []HitTileInfo, delta Vec) {
						projectilePos.X += delta.X
						projectilePos.Y += delta.Y
						isHorizontalCollision := false
						for _, c := range ci {
							if c.Normal.Y == -1 {
								projectileVel.Y = bounceVelocity
							}
							if c.Normal.X == -1 && projectileVel.X > 0 && projectileVel.Y > 0 {
								isHorizontalCollision = true
							}
							if c.Normal.X == 1 && projectileVel.X < 0 && projectileVel.Y > 0 {
								isHorizontalCollision = true
							}
						}
						if isHorizontalCollision {
							if ECWorld.Alive(q.Entity()) {
								toRemove = append(toRemove, q.Entity())
							}
						}
					},
					)
				}
			}
		}
	}
	return nil
}
func (c *Player) Draw() {
	if DrawDebugHitboxesEnabled {
		// Draw player tile for debug
		x, y, w, h := TileMapRes.GetTileRect(c.playerTile.X, c.playerTile.Y)
		x, y = CameraRes.ApplyCameraTransformToPoint(x, y)
		vector.DrawFilledRect(
			Screen,
			float32(x),
			float32(y),
			float32(w),
			float32(h),
			color.RGBA{0, 0, 128, 10},
			false,
		)
	}
}
