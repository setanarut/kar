package system

import (
	"image"
	"image/color"
	"kar"
	"kar/arc"
	"kar/engine/mathutil"
	"kar/items"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/tilecollider"
)

var (
	bounceVelocity = -math.Sqrt(2 * kar.SnowballGravity * kar.SnowballBounceHeight)
	blockHealth    float64
	// targetTile     image.Point
	placeTile image.Point

	isRayHit bool
)

type Player struct {
	playerTile image.Point
}

func (c *Player) Init() {
}

func (c *Player) Update() {
	// update animation player
	if !kar.GameDataRes.CraftingState {
		kar.PlayerAnimPlayer.Update()
	}

	if kar.ECWorld.Alive(kar.CurrentPlayer) {
		playerPos, playerSize, playerVelocity, playerHealth, playerController, pFacing := arc.MapPlayer.Get(kar.CurrentPlayer)

		playerCenterX, playerCenterY := playerPos.X+playerSize.W/2, playerPos.Y+playerSize.H/2

		if !kar.GameDataRes.CraftingState {

			// Update input
			playerController.IsBreakKeyPressed = ebiten.IsKeyPressed(ebiten.KeyRight)
			playerController.IsRunKeyPressed = ebiten.IsKeyPressed(ebiten.KeyShift)
			playerController.IsJumpKeyPressed = ebiten.IsKeyPressed(ebiten.KeySpace)
			playerController.IsAttackKeyJustPressed = inpututil.IsKeyJustPressed(ebiten.KeyLeft)
			playerController.IsJumpKeyJustPressed = inpututil.IsKeyJustPressed(ebiten.KeySpace)
			playerController.InputAxis = image.Point{}

			if ebiten.IsKeyPressed(ebiten.KeyW) {
				playerController.InputAxis.Y -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				playerController.InputAxis.Y += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				playerController.InputAxis.X -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyD) {
				playerController.InputAxis.X += 1
			}
			if !playerController.InputAxis.Eq(image.Point{}) {
				// restrict facing direction to 4 directions (no diagonal)
				switch playerController.InputAxis {
				case image.Point{0, -1}, image.Point{0, 1}, image.Point{-1, 0}, image.Point{1, 0}:
					pFacing.Dir = playerController.InputAxis
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
			switch playerController.CurrentState {
			case "idle":
				// enter idle
				if playerController.PreviousState != "idle" {
					if pFacing.Dir.Y == 0 {
						kar.PlayerAnimPlayer.SetState("idleRight")
					}
					if pFacing.Dir.X == 0 {
						kar.PlayerAnimPlayer.SetState("idleUp")
					}
				}

				// while idle
				if pFacing.Dir.Y == -1 {
					kar.PlayerAnimPlayer.SetState("idleUp")
				} else if pFacing.Dir.Y == 1 {
					kar.PlayerAnimPlayer.SetState("idleDown")
				} else if pFacing.Dir.X == 1 {
					kar.PlayerAnimPlayer.SetState("idleRight")
				} else if pFacing.Dir.X == -1 {
					kar.PlayerAnimPlayer.SetState("idleRight")
				}

				// All transitions with common exit code
				if playerController.IsJumpKeyJustPressed ||
					(playerController.IsOnFloor && playerController.HorizontalVelocity > 0.01) ||
					(!playerController.IsOnFloor && playerVelocity.Y > 0.01) ||
					(playerController.IsBreakKeyPressed && isRayHit) ||
					(playerVelocity.Y != 0 && playerVelocity.Y < -0.1) {
					playerController.PreviousState = playerController.CurrentState

					// ------- Common exit code  here -----------

					// Handle specific transitions
					if playerController.IsJumpKeyJustPressed {
						playerController.CurrentState = "jumping"
						if playerController.HorizontalVelocity > playerController.MinSpeedThresForJumpBoostMultiplier {
							playerVelocity.Y = playerController.JumpPower * playerController.JumpBoostMultiplier
						} else {
							playerVelocity.Y = playerController.JumpPower
						}
						playerController.JumpTimer = 0
					} else if playerController.IsOnFloor && playerController.HorizontalVelocity > 0.01 {
						if playerController.HorizontalVelocity > playerController.MaxWalkSpeed {
							playerController.CurrentState = "running"
						} else {
							playerController.CurrentState = "walking"
						}
					} else if !playerController.IsOnFloor && playerVelocity.Y > 0.01 {
						playerController.CurrentState = "falling"
					} else if playerController.IsBreakKeyPressed && isRayHit {
						playerController.CurrentState = "breaking"
					} else if playerVelocity.Y != 0 && playerVelocity.Y < -0.1 {
						playerController.CurrentState = "jumping"
					}
					break
				}

				// Update previous state if no transition occurred
				playerController.PreviousState = playerController.CurrentState
			case "walking":
				// enter walking
				if playerController.PreviousState != "walking" {
					kar.PlayerAnimPlayer.SetState("walkRight")
				}

				kar.PlayerAnimPlayer.SetStateFPS("walkRight", mathutil.MapRange(playerController.HorizontalVelocity, 0, playerController.MaxRunSpeed, 4, 23))

				// All transitions with common exit code
				if playerController.IsSkidding || (playerVelocity.Y > 0 && !playerController.IsOnFloor) ||
					playerController.IsJumpKeyJustPressed || playerController.HorizontalVelocity <= 0 ||
					playerController.HorizontalVelocity > playerController.MaxWalkSpeed {

					// Common exit code

					// Handle specific transitions
					if playerController.IsSkidding {
						playerController.PreviousState = playerController.CurrentState
						playerController.CurrentState = "skidding"
					} else if playerVelocity.Y > 0 && !playerController.IsOnFloor {
						playerController.PreviousState = playerController.CurrentState
						playerController.CurrentState = "falling"
						playerController.FallingDamageTempPosY = playerPos.Y
					} else if playerController.IsJumpKeyJustPressed {
						playerController.PreviousState = playerController.CurrentState
						playerController.CurrentState = "jumping"
						if playerController.HorizontalVelocity > playerController.MinSpeedThresForJumpBoostMultiplier {
							playerVelocity.Y = playerController.JumpPower * playerController.JumpBoostMultiplier
						} else {
							playerVelocity.Y = playerController.JumpPower
						}
						playerController.JumpTimer = 0
					} else if playerController.HorizontalVelocity <= 0 {
						playerController.PreviousState = playerController.CurrentState
						playerController.CurrentState = "idle"
					} else if playerController.HorizontalVelocity > playerController.MaxWalkSpeed {
						playerController.PreviousState = playerController.CurrentState
						playerController.CurrentState = "running"
					}
					break
				}

				// Update previous state if no transition occurred
				playerController.PreviousState = playerController.CurrentState

			case "running":
				// enter running
				if playerController.PreviousState != "running" {
					kar.PlayerAnimPlayer.SetState("walkRight")
				}

				// while running
				kar.PlayerAnimPlayer.SetStateFPS("walkRight", mathutil.MapRange(playerController.HorizontalVelocity, 0, playerController.MaxRunSpeed, 4, 23))

				// All transitions with common exit code
				if playerController.IsSkidding || (playerVelocity.Y > 0 && !playerController.IsOnFloor) ||
					playerController.IsJumpKeyJustPressed || playerController.HorizontalVelocity < 0.01 ||
					playerController.HorizontalVelocity <= playerController.MaxWalkSpeed {

					playerController.PreviousState = playerController.CurrentState

					// Common exit code

					// Handle specific transitions
					if playerController.IsSkidding {
						playerController.CurrentState = "skidding"
					} else if playerVelocity.Y > 0 && !playerController.IsOnFloor {
						playerController.CurrentState = "falling"
					} else if playerController.IsJumpKeyJustPressed {
						if playerController.HorizontalVelocity > playerController.MinSpeedThresForJumpBoostMultiplier {
							playerVelocity.Y = playerController.JumpPower * playerController.JumpBoostMultiplier
						} else {
							playerVelocity.Y = playerController.JumpPower
						}
						playerController.JumpTimer = 0
						playerController.CurrentState = "skidding"
					} else if playerController.HorizontalVelocity < 0.01 {
						playerController.CurrentState = "idle"
					} else if playerController.HorizontalVelocity <= playerController.MaxWalkSpeed {
						playerController.CurrentState = "walking"
					}
					break
				}

				// Update previous state if no transition occurred
				playerController.PreviousState = playerController.CurrentState

			case "jumping":
				if playerVelocity.Y != 0 && playerVelocity.Y > playerController.JumpPower+0.1 {
					kar.PlayerAnimPlayer.SetState("jump")
				}
				// skidding jumpg physics
				if playerController.PreviousState == "skidding" {
					if !playerController.IsJumpKeyPressed && playerController.JumpTimer < playerController.JumpReleaseTimer {
						playerVelocity.Y = playerController.ShortJumpVelocity * 0.7 // Kısa zıplama gücünü azalt
						playerController.JumpTimer = playerController.JumpHoldTime
					} else if playerController.IsJumpKeyPressed && playerController.JumpTimer < playerController.JumpHoldTime {
						playerVelocity.Y += playerController.JumpBoost * 0.7 // Boost gücünü azalt
						playerController.JumpTimer++
					} else if playerVelocity.Y >= 0.01 {
						playerController.PreviousState = playerController.CurrentState
						playerController.CurrentState = "falling"
						break
					}
				} else {
					// normal skidding
					if !playerController.IsJumpKeyPressed && playerController.JumpTimer < playerController.JumpReleaseTimer {
						playerVelocity.Y = playerController.ShortJumpVelocity
						playerController.JumpTimer = playerController.JumpHoldTime
					} else if playerController.IsJumpKeyPressed && playerController.JumpTimer < playerController.JumpHoldTime {
						speedFactor := (playerController.HorizontalVelocity / playerController.MaxRunSpeed) * playerController.SpeedJumpFactor
						playerVelocity.Y += playerController.JumpBoost * (1 + speedFactor)
						playerController.JumpTimer++
					} else if playerVelocity.Y >= 0 {
						playerController.PreviousState = playerController.CurrentState
						playerController.CurrentState = "falling"
						break
					}
				}

				// horizontal movement
				if playerController.InputAxis.X < 0 && playerVelocity.X > 0 {
					playerVelocity.X -= playerController.Deceleration
				} else if playerController.InputAxis.X > 0 && playerVelocity.X < 0 {
					playerVelocity.X += playerController.Deceleration
				}
				// Update previous state if no transition occurred
				playerController.PreviousState = playerController.CurrentState

			case "falling":
				// enter falling
				if playerController.PreviousState != "falling" {
					playerController.FallingDamageTempPosY = playerPos.Y
				}

				// while falling
				if playerVelocity.Y > 0.1 {
					kar.PlayerAnimPlayer.SetState("jump")
				}

				// transitions
				if playerController.IsOnFloor {

					d := int((playerPos.Y - playerController.FallingDamageTempPosY) / 30)
					if d > 3 {
						playerHealth.Current -= d - 3
					}

					playerController.PreviousState = playerController.CurrentState
					if playerController.HorizontalVelocity <= 0 {
						playerController.CurrentState = "idle"
					} else if playerController.IsRunKeyPressed {
						playerController.CurrentState = "running"
					} else {
						playerController.CurrentState = "walking"
					}

					break
				}
				// Update previous state if no transition occurred
				playerController.PreviousState = playerController.CurrentState
			case "breaking":
				// set animation states
				if pFacing.Dir.X == 1 {
					if playerController.HorizontalVelocity > 0.01 {
						kar.PlayerAnimPlayer.SetState("attackWalk")
					} else {
						kar.PlayerAnimPlayer.SetState("attackRight")
					}
				} else if pFacing.Dir.X == -1 {
					if playerController.HorizontalVelocity > 0.01 {
						kar.PlayerAnimPlayer.SetState("attackWalk")
					} else {
						kar.PlayerAnimPlayer.SetState("attackRight")
					}
				} else if pFacing.Dir.Y == 1 {
					kar.PlayerAnimPlayer.SetState("attackDown")
				} else if pFacing.Dir.Y == -1 {
					kar.PlayerAnimPlayer.SetState("attackUp")
				}

				if isRayHit {
					blockID := kar.TileMapRes.Get(kar.GameDataRes.TargetBlockCoord.X, kar.GameDataRes.TargetBlockCoord.Y)
					if !items.HasTag(blockID, items.Unbreakable) {
						if items.IsBestTool(blockID, kar.InventoryRes.CurrentSlotID()) {
							blockHealth += kar.PlayerBestToolDamage
						} else {
							blockHealth += kar.PlayerDefaultDamage
						}
					}
					// Destroy block
					if blockHealth >= 180 {
						blockHealth = 0
						kar.TileMapRes.Set(kar.GameDataRes.TargetBlockCoord.X, kar.GameDataRes.TargetBlockCoord.Y, items.Air)
						if items.HasTag(kar.InventoryRes.CurrentSlotID(), items.Tool) {
							kar.InventoryRes.CurrentSlot().Durability--
							if kar.InventoryRes.CurrentSlot().Durability <= 0 {
								kar.InventoryRes.ClearCurrentSlot()
							}
						}
						// spawn drop item
						x, y := kar.TileMapRes.TileToWorldCenter(kar.GameDataRes.TargetBlockCoord.X, kar.GameDataRes.TargetBlockCoord.Y)
						dropid := items.Property[blockID].DropID
						if blockID == items.OakLeaves {
							if rand.N(2) == 0 {
								dropid = items.OakLeaves
							}
						}
						AppendToSpawnList(x, y, dropid, 0)
					}
				}
				if !isRayHit {
					playerController.PreviousState = playerController.CurrentState
					playerController.CurrentState = "idle"
					// break
				}
				if !playerController.IsOnFloor && playerVelocity.Y > 0.01 {
					playerController.PreviousState = playerController.CurrentState
					playerController.CurrentState = "falling"
					// break
				} else if !playerController.IsBreakKeyPressed && playerController.IsOnFloor {
					playerController.PreviousState = playerController.CurrentState
					playerController.CurrentState = "idle"
					// break
				} else if !playerController.IsBreakKeyPressed && playerController.IsJumpKeyJustPressed {
					playerVelocity.Y = playerController.JumpPower
					playerController.JumpTimer = 0
					playerController.PreviousState = playerController.CurrentState
					playerController.CurrentState = "jumping"
					// break
				}
				// exit breaking
				if playerController.CurrentState != "breaking" {
					blockHealth = 0
				}
				// Update previous state if no transition occurred
				playerController.PreviousState = playerController.CurrentState

			case "skidding":
				// enter skidding
				if playerController.PreviousState != "skidding" {
					kar.PlayerAnimPlayer.SetState("skidding")
				}

				// <- while skidding scope ->

				// All transitions with common exit code
				if playerController.SkiddingJumpEnabled && playerController.IsJumpKeyJustPressed ||
					playerController.HorizontalVelocity < 0.01 || !playerController.IsSkidding {
					playerController.PreviousState = playerController.CurrentState

					// Common exit code

					// Handle specific transitions
					if playerController.SkiddingJumpEnabled && playerController.IsJumpKeyJustPressed {
						playerVelocity.X = 0
						playerController.HorizontalVelocity = 0

						// Yeni yöne doğru çok küçük sabit değerle başla
						if playerController.InputAxis.X > 0 {
							playerVelocity.X = 0.3
						} else if playerController.InputAxis.X < 0 {
							playerVelocity.X = -0.3
						}

						playerVelocity.Y = playerController.JumpPower * 0.7 // Zıplama gücünü azalt
						playerController.JumpTimer = 0
						playerController.CurrentState = "jumping"
					} else if playerController.HorizontalVelocity < 0.01 {
						playerController.CurrentState = "idle"
					} else if !playerController.IsSkidding {
						if playerController.HorizontalVelocity > playerController.MaxWalkSpeed {
							playerController.CurrentState = "running"
						} else {
							playerController.CurrentState = "walking"
						}
					}
					break // exit from skidding state
				}

				// Update previous state if no transition occurred
				playerController.PreviousState = playerController.CurrentState

			}

			// ########### UPDATE PHYSICS ################

			maxSpeed := playerController.MaxWalkSpeed
			currentAccel := playerController.WalkAcceleration
			currentDecel := playerController.WalkDeceleration
			playerController.HorizontalVelocity = math.Abs(playerVelocity.X)

			playerVelocity.Y += playerController.Gravity
			playerVelocity.Y = min(playerController.MaxFallSpeed, playerVelocity.Y)

			// Enemy collisions
			enemyQuery := arc.FilterEnemy.Query(&kar.ECWorld)
			for enemyQuery.Next() {
				enemyPos, enemySize, _, _ := enemyQuery.Get()
				collInfo := CheckCollision(playerPos, playerSize, playerVelocity, enemyPos, enemySize)

				playerVelocity.X += collInfo.DeltaX
				playerVelocity.Y += collInfo.DeltaY

				if collInfo.Collided {
					switch collInfo.Normal[0] {
					case 1:
						playerVelocity.X += 3
						// c.Health.Health -= 8
					case -1:
						playerVelocity.X -= 3
						// c.Health.Health -= 8
					}
					switch collInfo.Normal[1] {
					case 1:
						playerVelocity.Y = 0
					case -1:
						playerVelocity.Y = -5
					}
				}
			}

			if !playerController.IsSkidding {
				if playerController.IsRunKeyPressed {
					maxSpeed = playerController.MaxRunSpeed
					currentAccel = playerController.RunAcceleration
					currentDecel = playerController.RunDeceleration
				} else if playerController.HorizontalVelocity > playerController.MaxWalkSpeed {
					currentDecel = playerController.RunDeceleration
				}
			}

			if playerController.InputAxis.X > 0 {
				if playerVelocity.X > maxSpeed {
					playerVelocity.X = max(maxSpeed, playerVelocity.X-currentDecel)
				} else {
					playerVelocity.X = min(maxSpeed, playerVelocity.X+currentAccel)
				}
			} else if playerController.InputAxis.X < 0 {
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

			playerController.IsSkidding = (playerVelocity.X > 0 && playerController.InputAxis.X == -1) || (playerVelocity.X < 0 && playerController.InputAxis.X == 1)

			// Player and tilemap collision
			kar.Collider.Collide(
				math.Round(playerPos.X),
				playerPos.Y,
				playerSize.W,
				playerSize.H,
				playerVelocity.X,
				playerVelocity.Y,
				func(collisionInfos []tilecollider.CollisionInfo[uint8], dx, dy float64) {
					playerController.IsOnFloor = false

					playerPos.X += dx
					playerPos.Y += dy

					// Reset velocity when collide
					for _, collisionInfo := range collisionInfos {
						if collisionInfo.Normal[1] == -1 {
							// Ground collision
							playerVelocity.Y = 0
							playerController.IsOnFloor = true // on floor collision
						}
						if collisionInfo.Normal[1] == 1 {
							// Ceil collision
							playerVelocity.Y = 0
						}
						if collisionInfo.Normal[0] == -1 {
							// Right wall collision
							playerVelocity.X = 0
							playerController.HorizontalVelocity = 0
						}
						if collisionInfo.Normal[0] == 1 {
							// Left wall collision
							playerVelocity.X = 0
							playerController.HorizontalVelocity = 0
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
			c.playerTile = kar.TileMapRes.WorldToTile(playerCenterX, playerCenterY)
			targetBlockTemp := kar.GameDataRes.TargetBlockCoord
			kar.GameDataRes.TargetBlockCoord, isRayHit = kar.TileMapRes.Raycast(
				c.playerTile,
				pFacing.Dir,
				kar.RaycastDist,
			)

			// reset attack if block focus changed
			if !kar.GameDataRes.TargetBlockCoord.Eq(targetBlockTemp) || !isRayHit {
				blockHealth = 0
			}

			// place block if IsAttackKeyJustPressed
			if playerController.IsAttackKeyJustPressed {
				anyItemOverlapsWithPlaceRect := false
				itemSize := &arc.Size{kar.GameDataRes.DropItemW, kar.GameDataRes.DropItemH}
				// if slot item is block
				if isRayHit && items.HasTag(kar.InventoryRes.CurrentSlot().ID, items.Block) {
					// Get tile rect
					placeTile = kar.GameDataRes.TargetBlockCoord.Sub(pFacing.Dir)
					placeTilePos := &arc.Position{float64(placeTile.X * 20), float64(placeTile.Y * 20)}
					placeTileSize := &arc.Size{20, 20}
					// check overlaps
					queryItem := arc.FilterDroppedItem.Query(&kar.ECWorld)
					for queryItem.Next() {
						_, itemPos, _, _, _ := queryItem.Get()
						anyItemOverlapsWithPlaceRect = Overlaps(itemPos, itemSize, placeTilePos, placeTileSize)
						if anyItemOverlapsWithPlaceRect {
							queryItem.Close()
							break

						}
					}
					if !anyItemOverlapsWithPlaceRect {
						// oyuncu place tile ile çarpışıyormu
						if !Overlaps(playerPos, playerSize, placeTilePos, placeTileSize) {
							// place block
							kar.TileMapRes.Set(placeTile.X, placeTile.Y, kar.InventoryRes.CurrentSlotID())
							// remove item
							kar.InventoryRes.RemoveItemFromSelectedSlot()
						}
					}
					// if slot item snowball, throw snowball
				} else if kar.InventoryRes.CurrentSlot().ID == items.Snowball {
					if playerController.CurrentState != "skidding" {
						switch pFacing.Dir {
						case image.Point{1, 0}:
							arc.SpawnProjectile(items.Snowball, playerCenterX, playerCenterY-4, kar.SnowballSpeedX, kar.SnowballMaxFallVelocity)
						case image.Point{-1, 0}:
							arc.SpawnProjectile(items.Snowball, playerCenterX, playerCenterY-4, -kar.SnowballSpeedX, kar.SnowballMaxFallVelocity)
						}
						kar.InventoryRes.RemoveItemFromSelectedSlot()
					}
				}

			}

			// Drop Item
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
				currentSlot := kar.InventoryRes.CurrentSlot()
				if currentSlot.ID != items.Air {
					AppendToSpawnList(
						playerCenterX,
						playerCenterY,
						currentSlot.ID,
						currentSlot.Durability,
					)
					kar.InventoryRes.RemoveItemFromSelectedSlot()
					onInventorySlotChanged()
				}
			}
			// projectile physics
			q := arc.FilterProjectile.Query(&kar.ECWorld)
			for q.Next() {
				itemID, projectilePos, projectileVel := q.Get()
				// snowball physics
				if itemID.ID == items.Snowball {
					projectileVel.Y += kar.SnowballGravity
					projectileVel.Y = min(projectileVel.Y, kar.SnowballMaxFallVelocity)
					kar.Collider.Collide(
						projectilePos.X,
						projectilePos.Y,
						kar.GameDataRes.DropItemW,
						kar.GameDataRes.DropItemH,
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
								if kar.ECWorld.Alive(q.Entity()) {
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

	if kar.DrawDebugHitboxesEnabled {
		// Draw player tile for debug
		x, y, w, h := kar.TileMapRes.GetTileRect(c.playerTile.X, c.playerTile.Y)
		x, y = kar.CameraRes.ApplyCameraTransformToPoint(x, y)
		vector.DrawFilledRect(
			kar.Screen,
			float32(x),
			float32(y),
			float32(w),
			float32(h),
			color.RGBA{0, 0, 128, 10},
			false,
		)
	}
}
