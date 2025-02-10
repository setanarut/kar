package kar

import (
	"fmt"
	"image"
	"image/color"
	"kar/items"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/tilecollider"
)

var (
	bounceVelocity = -math.Sqrt(2 * SnowballGravity * SnowballBounceHeight)
	blockHealth    float64
	placeTile      image.Point
	IsRayHit       bool
)

type Player struct {
	playerTile image.Point
}

func (c *Player) Init() {
}

func (c *Player) Update() {

	if !GameDataRes.CraftingState {
		// update animation player
		if !GameDataRes.CraftingState {
			PlayerAnimPlayer.Update()
		}

		if ECWorld.Alive(CurrentPlayer) {
			playerPos, playerSize, playerVelocity, playerHealth, ctrl, pFacing := MapPlayer.Get(CurrentPlayer)

			playerCenterX, playerCenterY := playerPos.X+playerSize.W/2, playerPos.Y+playerSize.H/2

			// Update input
			ctrl.IsBreakKeyPressed = ebiten.IsKeyPressed(ebiten.KeyRight)
			ctrl.IsRunKeyPressed = ebiten.IsKeyPressed(ebiten.KeyShift)
			ctrl.IsJumpKeyPressed = ebiten.IsKeyPressed(ebiten.KeySpace)
			ctrl.IsAttackKeyJustPressed = inpututil.IsKeyJustPressed(ebiten.KeyLeft)
			ctrl.IsJumpKeyJustPressed = inpututil.IsKeyJustPressed(ebiten.KeySpace)
			ctrl.InputAxis = image.Point{}

			if ebiten.IsKeyPressed(ebiten.KeyW) {
				ctrl.InputAxis.Y -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				ctrl.InputAxis.Y += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				ctrl.InputAxis.X -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyD) {
				ctrl.InputAxis.X += 1
			}
			if !ctrl.InputAxis.Eq(image.Point{}) {
				// restrict facing direction to 4 directions (no diagonal)
				switch ctrl.InputAxis {
				case image.Point{0, -1}, image.Point{0, 1}, image.Point{-1, 0}, image.Point{1, 0}:
					pFacing.Dir = ctrl.InputAxis
				default:
					pFacing.Dir = image.Point{0, 0}
				}
			}
			if playerVelocity.X > 0.01 {
				pFacing.Dir = image.Point{1, 0}
			} else if playerVelocity.X < -0.01 {
				pFacing.Dir = image.Point{-1, 0}
			}

			// Update states
			switch ctrl.CurrentState {
			case "idle":
				// enter idle
				if ctrl.PreviousState != "idle" {
					ctrl.PreviousState = ctrl.CurrentState
					if pFacing.Dir.Y == 0 {
						PlayerAnimPlayer.SetState("idleRight")
					}
					if pFacing.Dir.X == 0 {
						PlayerAnimPlayer.SetState("idleUp")
					}
				}

				// while idle
				if pFacing.Dir.Y == -1 {
					PlayerAnimPlayer.SetState("idleUp")
				} else if pFacing.Dir.Y == 1 {
					PlayerAnimPlayer.SetState("idleDown")
				} else if pFacing.Dir.X == 1 {
					PlayerAnimPlayer.SetState("idleRight")
				} else if pFacing.Dir.X == -1 {
					PlayerAnimPlayer.SetState("idleRight")
				}

				// Handle specific transitions
				if ctrl.IsJumpKeyJustPressed {
					if ctrl.HorizontalVelocity > ctrl.MinSpeedThresForJumpBoostMultiplier {
						playerVelocity.Y = ctrl.JumpPower * ctrl.JumpBoostMultiplier
					} else {
						playerVelocity.Y = ctrl.JumpPower
					}
					ctrl.JumpTimer = 0
					ctrl.CurrentState = "jumping"
				} else if ctrl.IsOnFloor && ctrl.HorizontalVelocity > 0.01 {
					if ctrl.HorizontalVelocity > ctrl.MaxWalkSpeed {
						ctrl.CurrentState = "running"
					} else {
						ctrl.CurrentState = "walking"
					}
				} else if !ctrl.IsOnFloor && playerVelocity.Y > 0.01 {
					ctrl.CurrentState = "falling"
				} else if ctrl.IsBreakKeyPressed && IsRayHit {
					ctrl.CurrentState = "breaking"
				} else if playerVelocity.Y != 0 && playerVelocity.Y < -0.1 {
					ctrl.CurrentState = "jumping"
				}
				// exit idle
				if ctrl.PreviousState != ctrl.CurrentState {
					fmt.Println("exit idle")
				}
			case "walking":
				// enter walking
				if ctrl.PreviousState != "walking" {
					ctrl.PreviousState = ctrl.CurrentState
					PlayerAnimPlayer.SetState("walkRight")
				}
				PlayerAnimPlayer.SetStateFPS("walkRight", MapRange(ctrl.HorizontalVelocity, 0, ctrl.MaxRunSpeed, 4, 23))

				// Handle specific transitions
				if ctrl.IsSkidding {
					ctrl.CurrentState = "skidding"
				} else if playerVelocity.Y > 0 && !ctrl.IsOnFloor {
					ctrl.CurrentState = "falling"
				} else if ctrl.IsJumpKeyJustPressed {
					ctrl.CurrentState = "jumping"
					if ctrl.HorizontalVelocity > ctrl.MinSpeedThresForJumpBoostMultiplier {
						playerVelocity.Y = ctrl.JumpPower * ctrl.JumpBoostMultiplier
					} else {
						playerVelocity.Y = ctrl.JumpPower
					}
					ctrl.JumpTimer = 0
				} else if ctrl.HorizontalVelocity <= 0 {
					ctrl.CurrentState = "idle"
				} else if ctrl.HorizontalVelocity > ctrl.MaxWalkSpeed {
					ctrl.CurrentState = "running"
				}

				// exit walking
				if ctrl.PreviousState != ctrl.CurrentState {
					fmt.Println("exit walking")
				}
			case "running":
				// enter running
				if ctrl.PreviousState != "running" {
					ctrl.PreviousState = ctrl.CurrentState
					PlayerAnimPlayer.SetState("walkRight")
				}

				// while running
				PlayerAnimPlayer.SetStateFPS("walkRight", MapRange(ctrl.HorizontalVelocity, 0, ctrl.MaxRunSpeed, 4, 23))

				// Handle specific transitions
				if ctrl.IsSkidding {
					ctrl.CurrentState = "skidding"
				} else if playerVelocity.Y > 0 && !ctrl.IsOnFloor {
					ctrl.CurrentState = "falling"
				} else if ctrl.IsJumpKeyJustPressed {
					if ctrl.HorizontalVelocity > ctrl.MinSpeedThresForJumpBoostMultiplier {
						playerVelocity.Y = ctrl.JumpPower * ctrl.JumpBoostMultiplier
					} else {
						playerVelocity.Y = ctrl.JumpPower
					}
					ctrl.JumpTimer = 0
					ctrl.CurrentState = "jumping"
				} else if ctrl.HorizontalVelocity < 0.01 {
					ctrl.CurrentState = "idle"
				} else if ctrl.HorizontalVelocity <= ctrl.MaxWalkSpeed {
					ctrl.CurrentState = "walking"
				}
				// exit running
				if ctrl.PreviousState != ctrl.CurrentState {
					fmt.Println("exit running")
				}
			case "jumping":
				// enter running
				if ctrl.PreviousState != "jumping" {
					ctrl.PreviousState = ctrl.CurrentState
					PlayerAnimPlayer.SetState("jump")
				}

				// skidding jumpg physics
				if ctrl.PreviousState == "skidding" {
					if !ctrl.IsJumpKeyPressed && ctrl.JumpTimer < ctrl.JumpReleaseTimer {
						playerVelocity.Y = ctrl.ShortJumpVelocity * 0.7 // Kısa zıplama gücünü azalt
						ctrl.JumpTimer = ctrl.JumpHoldTime
					} else if ctrl.IsJumpKeyPressed && ctrl.JumpTimer < ctrl.JumpHoldTime {
						playerVelocity.Y += ctrl.JumpBoost * 0.7 // Boost gücünü azalt
						ctrl.JumpTimer++
					} else if playerVelocity.Y >= 0.01 {
						ctrl.CurrentState = "falling"
					}
				} else {
					// normal skidding
					if !ctrl.IsJumpKeyPressed && ctrl.JumpTimer < ctrl.JumpReleaseTimer {
						playerVelocity.Y = ctrl.ShortJumpVelocity
						ctrl.JumpTimer = ctrl.JumpHoldTime
					} else if ctrl.IsJumpKeyPressed && ctrl.JumpTimer < ctrl.JumpHoldTime {
						speedFactor := (ctrl.HorizontalVelocity / ctrl.MaxRunSpeed) * ctrl.SpeedJumpFactor
						playerVelocity.Y += ctrl.JumpBoost * (1 + speedFactor)
						ctrl.JumpTimer++
					} else if playerVelocity.Y >= 0 {
						ctrl.CurrentState = "falling"
					}
				}

				// horizontal movement
				if ctrl.InputAxis.X < 0 && playerVelocity.X > 0 {
					playerVelocity.X -= ctrl.Deceleration
				} else if ctrl.InputAxis.X > 0 && playerVelocity.X < 0 {
					playerVelocity.X += ctrl.Deceleration
				}

				// exit jumping
				if ctrl.PreviousState != ctrl.CurrentState {
					fmt.Println("exit jumping")
				}
			case "falling":
				// enter falling
				if ctrl.PreviousState != "falling" {
					ctrl.PreviousState = ctrl.CurrentState
					ctrl.FallingDamageTempPosY = playerPos.Y
					PlayerAnimPlayer.SetState("jump")
				}

				// transitions
				if ctrl.IsOnFloor {
					if ctrl.HorizontalVelocity <= 0 {
						ctrl.CurrentState = "idle"
					} else if ctrl.IsRunKeyPressed {
						ctrl.CurrentState = "running"
					} else {
						ctrl.CurrentState = "walking"
					}

				}
				// exit falling
				if ctrl.PreviousState != ctrl.CurrentState {
					fmt.Println("exit falling")
					d := int((playerPos.Y - ctrl.FallingDamageTempPosY) / 30)
					if d > 3 {
						playerHealth.Current -= d - 3
					}
				}
			case "breaking":
				// enter running
				if ctrl.PreviousState != "breaking" {
					fmt.Println("enter breaking")
					ctrl.PreviousState = ctrl.CurrentState
				}
				// update animation states
				if pFacing.Dir.X == 1 {
					if ctrl.HorizontalVelocity > 0.01 {
						PlayerAnimPlayer.SetState("attackWalk")
					} else {
						PlayerAnimPlayer.SetState("attackRight")
					}
				} else if pFacing.Dir.X == -1 {
					if ctrl.HorizontalVelocity > 0.01 {
						PlayerAnimPlayer.SetState("attackWalk")
					} else {
						PlayerAnimPlayer.SetState("attackRight")
					}
				} else if pFacing.Dir.Y == 1 {
					PlayerAnimPlayer.SetState("attackDown")
				} else if pFacing.Dir.Y == -1 {
					PlayerAnimPlayer.SetState("attackUp")
				}

				// break block
				if IsRayHit {
					blockID := TileMapRes.Get(GameDataRes.TargetBlockCoord.X, GameDataRes.TargetBlockCoord.Y)
					if !items.HasTag(blockID, items.Unbreakable) {
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
				if !IsRayHit || (!ctrl.IsBreakKeyPressed && ctrl.IsOnFloor) {
					ctrl.CurrentState = "idle"
				} else if !ctrl.IsOnFloor && playerVelocity.Y > 0.01 {
					ctrl.CurrentState = "falling"
				} else if !ctrl.IsBreakKeyPressed && ctrl.IsJumpKeyJustPressed {
					playerVelocity.Y = ctrl.JumpPower
					ctrl.JumpTimer = 0
					ctrl.CurrentState = "jumping"
				}
				// exit breaking
				if ctrl.PreviousState != ctrl.CurrentState {
					fmt.Println("exit breaking")
					blockHealth = 0
				}
			case "skidding":
				// enter skidding
				if ctrl.PreviousState != "skidding" {
					fmt.Println("enter skidding")
					ctrl.PreviousState = ctrl.CurrentState
					PlayerAnimPlayer.SetState("skidding")
				}
				// Handle specific transitions
				if ctrl.SkiddingJumpEnabled && ctrl.IsJumpKeyJustPressed {
					playerVelocity.X = 0
					ctrl.HorizontalVelocity = 0
					// Yeni yöne doğru çok küçük sabit değerle başla
					if ctrl.InputAxis.X > 0 {
						playerVelocity.X = 0.3
					} else if ctrl.InputAxis.X < 0 {
						playerVelocity.X = -0.3
					}
					playerVelocity.Y = ctrl.JumpPower * 0.7 // Zıplama gücünü azalt
					ctrl.JumpTimer = 0
					ctrl.CurrentState = "jumping"
				} else if ctrl.HorizontalVelocity < 0.01 {
					ctrl.CurrentState = "idle"
				} else if !ctrl.IsSkidding {
					if ctrl.HorizontalVelocity > ctrl.MaxWalkSpeed {
						ctrl.CurrentState = "running"
					} else {
						ctrl.CurrentState = "walking"
					}
				}
				if ctrl.PreviousState != ctrl.CurrentState {
					fmt.Println("exit skidding")
				}
			}

			// ########### UPDATE PHYSICS ################

			maxSpeed := ctrl.MaxWalkSpeed
			currentAccel := ctrl.WalkAcceleration
			currentDecel := ctrl.WalkDeceleration
			ctrl.HorizontalVelocity = math.Abs(playerVelocity.X)

			playerVelocity.Y += ctrl.Gravity
			playerVelocity.Y = min(ctrl.MaxFallSpeed, playerVelocity.Y)

			if !ctrl.IsSkidding {
				if ctrl.IsRunKeyPressed {
					maxSpeed = ctrl.MaxRunSpeed
					currentAccel = ctrl.RunAcceleration
					currentDecel = ctrl.RunDeceleration
				} else if ctrl.HorizontalVelocity > ctrl.MaxWalkSpeed {
					currentDecel = ctrl.RunDeceleration
				}
			}

			if ctrl.InputAxis.X > 0 {
				if playerVelocity.X > maxSpeed {
					playerVelocity.X = max(maxSpeed, playerVelocity.X-currentDecel)
				} else {
					playerVelocity.X = min(maxSpeed, playerVelocity.X+currentAccel)
				}
			} else if ctrl.InputAxis.X < 0 {
				if playerVelocity.X < -maxSpeed {
					playerVelocity.X = min(-maxSpeed, playerVelocity.X+currentDecel)
				} else {
					playerVelocity.X = max(-maxSpeed, playerVelocity.X-currentAccel)
				}
			} else {
				if playerVelocity.X > 0 {
					playerVelocity.X = max(0, playerVelocity.X-currentDecel)
				} else if playerVelocity.X < 0 {
					playerVelocity.X = min(0, playerVelocity.X+currentDecel)
				}
			}

			ctrl.IsSkidding = (playerVelocity.X > 0 && ctrl.InputAxis.X == -1) || (playerVelocity.X < 0 && ctrl.InputAxis.X == 1)

			// Player and tilemap collision
			Collider.Collide(
				math.Round(playerPos.X),
				playerPos.Y,
				playerSize.W,
				playerSize.H,
				playerVelocity.X,
				playerVelocity.Y,
				func(collisionInfos []tilecollider.CollisionInfo[uint8], dx, dy float64) {
					ctrl.IsOnFloor = false

					playerPos.X += dx
					playerPos.Y += dy

					// Reset velocity when collide
					for _, collisionInfo := range collisionInfos {
						if collisionInfo.Normal[1] == -1 {
							// Ground collision
							playerVelocity.Y = 0
							ctrl.IsOnFloor = true // on floor collision
						}
						if collisionInfo.Normal[1] == 1 {
							// Ceil collision
							playerVelocity.Y = 0
						}
						if collisionInfo.Normal[0] == -1 || collisionInfo.Normal[0] == 1 {

							if ctrl.HorizontalVelocity == ctrl.MaxRunSpeed && ctrl.IsBreakKeyPressed {
								// Right of Left wall collision
								x := collisionInfo.TileCoords[0]
								y := collisionInfo.TileCoords[1]
								TileMapRes.Set(x, y, items.Air)
								wx, wy := TileMapRes.TileToWorldCenter(x, y)
								SpawnEffect(collisionInfo.TileID, wx, wy)
							}

							playerVelocity.X = 0
							ctrl.HorizontalVelocity = 0
						}
					}

					if inpututil.IsKeyJustPressed(ebiten.KeyS) {
						ids := make([]uint8, 0)
						for _, collisionInfo := range collisionInfos {
							if collisionInfo.Normal[1] == -1 {
								ids = append(ids, collisionInfo.TileID)
							}
						}
						if len(ids) == 2 {
							if ids[0] == items.Sand && ids[1] == items.GrassBlock {
								// Pipe logic is here
							}
						}
					}
				},
			)

			// player facing raycast for target block
			c.playerTile = TileMapRes.WorldToTile(playerCenterX, playerCenterY)
			targetBlockTemp := GameDataRes.TargetBlockCoord
			GameDataRes.TargetBlockCoord, IsRayHit = TileMapRes.Raycast(
				c.playerTile,
				pFacing.Dir,
				RaycastDist,
			)

			// reset attack if block focus changed
			if !GameDataRes.TargetBlockCoord.Eq(targetBlockTemp) || !IsRayHit {
				blockHealth = 0
			}

			// place block if IsAttackKeyJustPressed
			if ctrl.IsAttackKeyJustPressed {
				anyItemOverlapsWithPlaceRect := false
				// if slot item is block
				if IsRayHit && items.HasTag(InventoryRes.CurrentSlot().ID, items.Block) {
					// Get tile rect
					placeTile = GameDataRes.TargetBlockCoord.Sub(pFacing.Dir)
					placeTilePos := &Position{float64(placeTile.X * 20), float64(placeTile.Y * 20)}
					placeTileSize := &Size{20, 20}
					// check overlaps
					queryItem := FilterDroppedItem.Query(&ECWorld)
					for queryItem.Next() {
						_, itemPos, _, _, _ := queryItem.Get()
						anyItemOverlapsWithPlaceRect = Overlaps(itemPos, &DropItemSize, placeTilePos, placeTileSize)
						if anyItemOverlapsWithPlaceRect {
							queryItem.Close()
							break

						}
					}
					if !anyItemOverlapsWithPlaceRect {
						// oyuncu place tile ile çarpışıyormu
						if !Overlaps(playerPos, playerSize, placeTilePos, placeTileSize) {
							// place block
							TileMapRes.Set(placeTile.X, placeTile.Y, InventoryRes.CurrentSlotID())
							// remove item
							InventoryRes.RemoveItemFromSelectedSlot()
						}
					}
					// if slot item snowball, throw snowball
				} else if InventoryRes.CurrentSlot().ID == items.Snowball {
					if ctrl.CurrentState != "skidding" {
						switch pFacing.Dir {
						case image.Point{1, 0}:
							SpawnProjectile(items.Snowball, playerCenterX, playerCenterY-4, SnowballSpeedX, SnowballMaxFallVelocity)
						case image.Point{-1, 0}:
							SpawnProjectile(items.Snowball, playerCenterX, playerCenterY-4, -SnowballSpeedX, SnowballMaxFallVelocity)
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
						playerCenterX,
						playerCenterY,
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
					Collider.Collide(
						projectilePos.X,
						projectilePos.Y,
						DropItemSize.W,
						DropItemSize.H,
						projectileVel.X,
						projectileVel.Y,
						func(ci []tilecollider.CollisionInfo[uint8], dx, dy float64) {
							projectilePos.X += dx
							projectilePos.Y += dy
							isHorizontalCollision := false
							for _, c := range ci {
								if c.Normal[1] == -1 {
									projectileVel.Y = bounceVelocity
								}
								if c.Normal[0] == -1 && projectileVel.X > 0 && projectileVel.Y > 0 {
									isHorizontalCollision = true
								}
								if c.Normal[0] == 1 && projectileVel.X < 0 && projectileVel.Y > 0 {
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
