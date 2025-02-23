package kar

import (
	"image"
	"image/color"
	"kar/items"
	"kar/tilemap"
	"kar/v"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	bounceVelocity float64
	placeTile      image.Point
	dyingCountdown float64

	itemHit     *HitInfo
	itemBox     *AABB
	snowBallBox AABB
	isOnFloor   bool
	playerTile  image.Point
}

func (p *Player) Init() {
	p.bounceVelocity = -math.Sqrt(2 * SnowballGravity * SnowballBounceHeight)
	p.itemBox = &AABB{Half: dropItemHalfSize}
	p.snowBallBox = AABB{Half: Vec{4, 4}}
	p.itemHit = &HitInfo{}
}

func (p *Player) Update() {

	if gameDataRes.GameplayState == Playing {
		if world.Alive(currentPlayer) {
			animPlayer.Update()

			pBox, pVelocity, pHealth, ctrl, pFacing := mapPlayer.Get(currentPlayer)
			absXVelocity := math.Abs(pVelocity.X)
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
					pFacing.Vec = inputAxis
				default:
					pFacing.Vec = Vec{}
				}
			}

			if absXVelocity > 0.01 {
				pFacing.Vec = Vec{math.Copysign(1, pVelocity.X), 0}
			}

			isSkidding := inputAxis.X != 0 && (pVelocity.X*inputAxis.X < 0)

			// Death animation
			if pHealth.Current <= 0 {
				p.dyingCountdown += 0.1
				animPlayer.Data.Paused = true
				pBox.Pos.Y += p.dyingCountdown

				if pBox.Pos.Y > cameraRes.Y+cameraRes.Height {
					animPlayer.Data.Paused = false
					previousGameState = "playing"
					currentGameState = "menu"
					ColorM.ChangeHSV(1, 0, 0.5) // BW
					TextDO.ColorScale.Scale(0.5, 0.5, 0.5, 1)
					world.RemoveEntity(currentPlayer)
					p.dyingCountdown = 0
				}
			}

			// Update states
			switch ctrl.CurrentState {
			case "idle":
				// enter idle
				if ctrl.PreviousState != "idle" {
					ctrl.PreviousState = ctrl.CurrentState
					if pFacing.Y == 0 {
						animPlayer.SetAnim("idleRight")
					}
					if pFacing.X == 0 {
						animPlayer.SetAnim("idleUp")
					}
				}

				// while idle
				if pFacing.Y == -1 {
					animPlayer.SetAnim("idleUp")
				} else if pFacing.Y == 1 {
					animPlayer.SetAnim("idleDown")
				} else if pFacing.X == 1 {
					animPlayer.SetAnim("idleRight")
				} else if pFacing.X == -1 {
					animPlayer.SetAnim("idleRight")
				}

				// Handle specific transitions
				if isJumpKeyJustPressed {
					if absXVelocity > ctrl.MinSpeedThresForJumpBoostMultiplier {
						pVelocity.Y = ctrl.JumpPower * ctrl.JumpBoostMultiplier
					} else {
						pVelocity.Y = ctrl.JumpPower
					}
					ctrl.JumpTimer = 0
					ctrl.CurrentState = "jumping"
				} else if p.isOnFloor && absXVelocity > 0.01 {
					if absXVelocity > ctrl.MaxWalkSpeed {
						ctrl.CurrentState = "running"
					} else {
						ctrl.CurrentState = "walking"
					}
				} else if !p.isOnFloor && pVelocity.Y > 0.01 {
					ctrl.CurrentState = "falling"
				} else if isBreakKeyPressed && gameDataRes.IsRayHit {
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
					animPlayer.SetAnim("walkRight")
				}
				animPlayer.SetAnimFPS("walkRight", MapRange(absXVelocity, 0, ctrl.MaxRunSpeed, 4, 23))

				// Handle specific transitions
				if isSkidding {
					ctrl.CurrentState = "skidding"
				} else if pVelocity.Y > 0 && !p.isOnFloor {
					ctrl.CurrentState = "falling"
				} else if isJumpKeyJustPressed {
					ctrl.CurrentState = "jumping"
					if absXVelocity > ctrl.MinSpeedThresForJumpBoostMultiplier {
						pVelocity.Y = ctrl.JumpPower * ctrl.JumpBoostMultiplier
					} else {
						pVelocity.Y = ctrl.JumpPower
					}
					ctrl.JumpTimer = 0
				} else if absXVelocity == 0 {
					ctrl.CurrentState = "idle"
				} else if absXVelocity > ctrl.MaxWalkSpeed {
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
					animPlayer.SetAnim("walkRight")
				}

				// while running
				animPlayer.SetAnimFPS("walkRight", MapRange(absXVelocity, 0, ctrl.MaxRunSpeed, 4, 23))

				// Handle specific transitions
				if isSkidding {
					ctrl.CurrentState = "skidding"
				} else if pVelocity.Y > 0 && !p.isOnFloor {
					ctrl.CurrentState = "falling"
				} else if isJumpKeyJustPressed {
					if absXVelocity > ctrl.MinSpeedThresForJumpBoostMultiplier {
						pVelocity.Y = ctrl.JumpPower * ctrl.JumpBoostMultiplier
					} else {
						pVelocity.Y = ctrl.JumpPower
					}
					ctrl.JumpTimer = 0
					ctrl.CurrentState = "jumping"
				} else if absXVelocity < 0.01 {
					ctrl.CurrentState = "idle"
				} else if absXVelocity <= ctrl.MaxWalkSpeed {
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
					animPlayer.SetAnim("jump")
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
						speedFactor := (absXVelocity / ctrl.MaxRunSpeed) * ctrl.SpeedJumpFactor
						pVelocity.Y += ctrl.JumpBoost * (1 + speedFactor)
						ctrl.JumpTimer++
					} else if pVelocity.Y >= 0 {
						ctrl.CurrentState = "falling"
					}
				}

				// apply air skidding Decel
				if isSkidding {
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
					animPlayer.SetAnim("jump")
				}

				// TODO erken zıplama toleransı için mantık yaz.

				// transitions
				if p.isOnFloor {
					if absXVelocity <= 0 {
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
					if absXVelocity > 0.01 {
						animPlayer.SetAnim("attackWalk")
					} else {
						animPlayer.SetAnim("attackRight")
					}
				} else if pFacing.X == -1 {
					if absXVelocity > 0.01 {
						animPlayer.SetAnim("attackWalk")
					} else {
						animPlayer.SetAnim("attackRight")
					}
				} else if pFacing.Y == 1 {
					animPlayer.SetAnim("attackDown")
				} else if pFacing.Y == -1 {
					animPlayer.SetAnim("attackUp")
				}

				// break block
				if gameDataRes.IsRayHit {
					blockID := tileMapRes.Get(gameDataRes.TargetBlockCoord.X, gameDataRes.TargetBlockCoord.Y)
					if !items.HasTag(blockID, items.UnbreakableBlock) {
						if items.IsBestTool(blockID, inventoryRes.CurrentSlotID()) {
							gameDataRes.BlockHealth += PlayerBestToolDamage
						} else {
							gameDataRes.BlockHealth += PlayerDefaultDamage
						}
					}
					// Destroy block
					if gameDataRes.BlockHealth >= 180 {

						// set air
						tileMapRes.Set(gameDataRes.TargetBlockCoord.X, gameDataRes.TargetBlockCoord.Y, items.Air)
						gameDataRes.BlockHealth = 0

						if items.HasTag(inventoryRes.CurrentSlotID(), items.Tool) {
							// damage the tool
							inventoryRes.CurrentSlot().Durability--
							// If durability is 0, destroy the tool.
							if inventoryRes.CurrentSlot().Durability <= 0 {
								inventoryRes.ClearCurrentSlot()
							}
						}
						// Spawn drop item
						pos := tileMapRes.TileToWorldCenter(gameDataRes.TargetBlockCoord.X, gameDataRes.TargetBlockCoord.Y)
						dropid := items.Property[blockID].DropID
						if blockID == items.OakLeaves {
							if rand.N(2) == 0 {
								dropid = items.OakLeaves
							}
						}
						toSpawn = append(toSpawn, spawnData{pos, dropid, 0})
					}
				}
				// transitions
				if !gameDataRes.IsRayHit || (!isBreakKeyPressed && p.isOnFloor) {
					ctrl.CurrentState = "idle"
				} else if !p.isOnFloor && pVelocity.Y > 0.01 {
					ctrl.CurrentState = "falling"
				} else if !isBreakKeyPressed && isJumpKeyJustPressed {
					pVelocity.Y = ctrl.JumpPower
					ctrl.JumpTimer = 0
					ctrl.CurrentState = "jumping"
				}
				// exit breaking
				if ctrl.PreviousState != ctrl.CurrentState {
					// fmt.Println("exit breaking")
					gameDataRes.BlockHealth = 0
				}
			case "skidding":
				// enter skidding
				if ctrl.PreviousState != "skidding" {
					// fmt.Println("enter skidding")
					ctrl.PreviousState = ctrl.CurrentState
				}
				// Apply Skidding decel
				if absXVelocity > ctrl.SkiddingFriction {
					pVelocity.X += math.Copysign(ctrl.SkiddingFriction, -pVelocity.X)
				}
				if absXVelocity > 0.5 {
					animPlayer.SetAnim("skidding")
				}

				// Handle specific transitions
				if ctrl.SkiddingJumpEnabled && isJumpKeyJustPressed {
					// Yeni yöne doğru çok küçük sabit değerle başla
					pVelocity.X = math.Copysign(0.3, float64(inputAxis.X))
					pVelocity.Y = ctrl.JumpPower * 0.7 // Zıplama gücünü azalt
					ctrl.JumpTimer = 0
					ctrl.CurrentState = "jumping"
				} else if absXVelocity < 0.01 {
					ctrl.CurrentState = "idle"
				} else if !isSkidding {
					if absXVelocity > ctrl.MaxWalkSpeed {
						ctrl.CurrentState = "running"
					} else {
						ctrl.CurrentState = "walking"
					}
				}
				if ctrl.PreviousState != ctrl.CurrentState {
					// fmt.Println("exit skidding")
				}
			}

			/// Update Physics
			currentAccel := ctrl.WalkAcceleration
			currentDecel := ctrl.WalkDeceleration
			maxSpeed := ctrl.MaxWalkSpeed
			pVelocity.Y = min(ctrl.MaxFallSpeed, pVelocity.Y+ctrl.Gravity)

			if !isSkidding {
				if isRunKeyPressed {
					maxSpeed, currentAccel, currentDecel = ctrl.MaxRunSpeed, ctrl.RunAcceleration, ctrl.RunDeceleration
				} else if math.Abs(pVelocity.X) > ctrl.MaxWalkSpeed {
					currentDecel = ctrl.RunDeceleration
				}
			}

			targetSpeed := maxSpeed
			if inputAxis.X == 0 {
				targetSpeed = 0
			} else if inputAxis.X < 0 {
				targetSpeed = -maxSpeed
			}

			if pVelocity.X < targetSpeed {
				pVelocity.X = min(targetSpeed, pVelocity.X+currentAccel)
			} else {
				pVelocity.X = max(targetSpeed, pVelocity.X-currentDecel)
			}

			// Player and tilemap collision
			TileCollider.Collide(*pBox, pVelocity.Vec, func(collisionInfos []HitTileInfo, delta Vec) {
				p.isOnFloor = false
				pBox.Pos = pBox.Pos.Add(delta)
				// Reset velocity when collide
				for _, ci := range collisionInfos {

					tileID := tileMapRes.GetUnchecked(ci.TileCoords)

					if ci.Normal.Y == -1 {
						// Ground collision
						pVelocity.Y = 0
						p.isOnFloor = true
					}
					// Ceil collision
					if ci.Normal.Y == 1 { // TODO aynı anda olan çarpışmaları teke indir
						pVelocity.Y = 0

						switch tileID {
						case items.StoneBricks:
							// Destroy block when ceil hit
							tileMapRes.Set(ci.TileCoords.X, ci.TileCoords.Y, items.Air)
							effectPos := tileMapRes.TileToWorldCenter(ci.TileCoords.X, ci.TileCoords.Y)
							SpawnEffect(tileID, effectPos)
						case items.Random:
							if tileMapRes.GetUnchecked(ci.TileCoords.Add(tilemap.Up)) == items.Air {
								tileMapRes.Set(ci.TileCoords.X, ci.TileCoords.Y, items.Bedrock)
								ceilBlockCoord = ci.TileCoords
								ceilBlockTick = 3
							}
						}
					}
					// Right or Left wall collision
					if ci.Normal.X == -1 || ci.Normal.X == 1 {
						// While running at maximum speed, hold down the right arrow key and hit the block to destroy it.
						if absXVelocity == ctrl.MaxRunSpeed && isBreakKeyPressed {
							tileMapRes.Set(ci.TileCoords.X, ci.TileCoords.Y, items.Air)
							effectPos := tileMapRes.TileToWorldCenter(ci.TileCoords.X, ci.TileCoords.Y)
							SpawnEffect(tileID, effectPos)
						}
						pVelocity.X = 0
						// absXVelocity = 0
					}
				}
			},
			)

			// player facing raycast for target block
			p.playerTile = tileMapRes.WorldToTile(pBox.Pos.X, pBox.Pos.Y)
			targetBlockTemp := gameDataRes.TargetBlockCoord
			gameDataRes.TargetBlockCoord, gameDataRes.IsRayHit = tileMapRes.Raycast(
				p.playerTile,
				int(pFacing.X), int(pFacing.Y),
				RaycastDist,
			)

			// reset attack if block focus changed
			if !gameDataRes.TargetBlockCoord.Eq(targetBlockTemp) || !gameDataRes.IsRayHit {
				gameDataRes.BlockHealth = 0
			}

			// place block if IsAttackKeyJustPressed
			if isAttackKeyJustPressed {
				anyItemOverlapsWithPlaceRect := false
				// if slot item is block
				if gameDataRes.IsRayHit && items.HasTag(inventoryRes.CurrentSlot().ID, items.Block) {
					// Get tile rect
					p.placeTile = gameDataRes.TargetBlockCoord.Sub(image.Point{int(pFacing.X), int(pFacing.Y)})
					placeTileBox := &AABB{
						Pos:  Vec{float64(p.placeTile.X*20) + 10, float64(p.placeTile.Y*20) + 10},
						Half: Vec{10, 10},
					}
					// check overlaps
					queryItem := filterDroppedItem.Query()
					for queryItem.Next() {
						_, itemPos, _ := queryItem.Get()
						p.itemBox.Pos = itemPos.Vec
						anyItemOverlapsWithPlaceRect = placeTileBox.Overlap(p.itemBox, nil)
						if anyItemOverlapsWithPlaceRect {
							queryItem.Close()
							break
						}
					}
					if !anyItemOverlapsWithPlaceRect {
						if !pBox.Overlap(placeTileBox, nil) {
							// place block
							tileMapRes.Set(p.placeTile.X, p.placeTile.Y, inventoryRes.CurrentSlotID())
							// remove item
							inventoryRes.RemoveItemFromSelectedSlot()
						}
					}
					// if slot item snowball, throw snowball
				} else if inventoryRes.CurrentSlot().ID == items.Snowball {
					if ctrl.CurrentState != "skidding" {
						spawnPos := Vec{pBox.Pos.X, pBox.Pos.Y - 4}
						spawnVel := Vec{SnowballSpeedX, SnowballMaxFallVelocity}
						switch pFacing.Vec {
						case v.Right:
							SpawnProjectile(items.Snowball, spawnPos, spawnVel)
						case v.Left:
							SpawnProjectile(items.Snowball, spawnPos, spawnVel.NegX())
						}
						inventoryRes.RemoveItemFromSelectedSlot()
					}
				}
			}

			// drop Item
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
				currentSlot := inventoryRes.CurrentSlot()
				if currentSlot.ID != items.Air {
					toSpawn = append(toSpawn, spawnData{
						pBox.Pos,
						currentSlot.ID,
						currentSlot.Durability,
					})
					inventoryRes.RemoveItemFromSelectedSlot()
					onInventorySlotChanged()
				}
			}
			// projectile physics
			q := filterProjectile.Query()
			for q.Next() {
				itemID, projectilePos, projectileVel := q.Get()
				// snowball physics
				if itemID.ID == items.Snowball {
					projectileVel.Y += SnowballGravity
					projectileVel.Y = min(projectileVel.Y, SnowballMaxFallVelocity)
					p.snowBallBox.Pos = projectilePos.Vec
					TileCollider.Collide(p.snowBallBox, projectileVel.Vec, func(ci []HitTileInfo, delta Vec) {
						projectilePos.X += delta.X
						projectilePos.Y += delta.Y
						isHorizontalCollision := false
						for _, cinfo := range ci {
							if cinfo.Normal.Y == -1 {

								projectileVel.Y = p.bounceVelocity
							}
							if cinfo.Normal.X == -1 && projectileVel.X > 0 && projectileVel.Y > 0 {
								isHorizontalCollision = true
							}
							if cinfo.Normal.X == 1 && projectileVel.X < 0 && projectileVel.Y > 0 {
								isHorizontalCollision = true
							}
						}
						if isHorizontalCollision {
							if world.Alive(q.Entity()) {
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
	if DrawPlayerTileHitboxEnabled {
		// Draw player tile for debug
		x, y, w, h := tileMapRes.GetTileRect(c.playerTile.X, c.playerTile.Y)
		x, y = cameraRes.ApplyCameraTransformToPoint(x, y)
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
