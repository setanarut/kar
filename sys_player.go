package kar

import (
	"image"
	"kar/items"
	"kar/tilemap"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/v"
)

var elapsed = 0
var elapsed2 = 5
var blinking bool

type Player struct {
	placeTile  image.Point
	itemHit    *HitInfo
	isOnFloor  bool
	playerTile image.Point
}

func (p *Player) Init() {

	p.itemHit = &HitInfo{}
}

func (p *Player) Update() {

	if world.Alive(currentPlayer) {
		animPlayer.Update()

		if elapsed > 0 {
			elapsed--
			if elapsed2 > 0 {
				elapsed2--
			}
			if elapsed2 == 0 {
				elapsed2 = 5
				blinking = !blinking
			}
		}

		playerAABB, vl, pHealth, ctrl, fc := mapPlayer.GetUnchecked(currentPlayer)

		pFacing := (*Vec)(fc)
		playerVel := (*Vec)(vl)

		absXVelocity := math.Abs(playerVel.X)
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
				*pFacing = inputAxis
			default:
				*pFacing = Vec{}
			}
		}

		if absXVelocity > 0.01 {
			*pFacing = Vec{math.Copysign(1, playerVel.X), 0}
		}

		isSkidding := inputAxis.X != 0 && (playerVel.X*inputAxis.X < 0)

		if pHealth.Current <= 0 {
			// TODO ayrı olarak ölme animasyonu sistemi yaz
			world.RemoveEntity(currentPlayer)
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
					playerVel.Y = ctrl.JumpPower * ctrl.JumpBoostMultiplier
				} else {
					playerVel.Y = ctrl.JumpPower
				}
				ctrl.JumpTimer = 0
				ctrl.CurrentState = "jumping"
			} else if p.isOnFloor && absXVelocity > 0.01 {
				if absXVelocity > ctrl.MaxWalkSpeed {
					ctrl.CurrentState = "running"
				} else {
					ctrl.CurrentState = "walking"
				}
			} else if !p.isOnFloor && playerVel.Y > 0.01 {
				ctrl.CurrentState = "falling"
			} else if isBreakKeyPressed && gameDataRes.IsRayHit {
				ctrl.CurrentState = "breaking"
			} else if playerVel.Y != 0 && playerVel.Y < -0.1 {
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
			animPlayer.SetAnimFPS("walkRight", mapRange(absXVelocity, 0, ctrl.MaxRunSpeed, 4, 23))

			// Handle specific transitions
			if isSkidding {
				ctrl.CurrentState = "skidding"
			} else if playerVel.Y > 0 && !p.isOnFloor {
				ctrl.CurrentState = "falling"
			} else if isJumpKeyJustPressed {
				ctrl.CurrentState = "jumping"
				if absXVelocity > ctrl.MinSpeedThresForJumpBoostMultiplier {
					playerVel.Y = ctrl.JumpPower * ctrl.JumpBoostMultiplier
				} else {
					playerVel.Y = ctrl.JumpPower
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
			animPlayer.SetAnimFPS("walkRight", mapRange(absXVelocity, 0, ctrl.MaxRunSpeed, 4, 23))

			// Handle specific transitions
			if isSkidding {
				ctrl.CurrentState = "skidding"
			} else if playerVel.Y > 0 && !p.isOnFloor {
				ctrl.CurrentState = "falling"
			} else if isJumpKeyJustPressed {
				if absXVelocity > ctrl.MinSpeedThresForJumpBoostMultiplier {
					playerVel.Y = ctrl.JumpPower * ctrl.JumpBoostMultiplier
				} else {
					playerVel.Y = ctrl.JumpPower
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
					playerVel.Y = ctrl.ShortJumpVelocity * 0.7 // Kısa zıplama gücünü azalt
					ctrl.JumpTimer = ctrl.JumpHoldTime
				} else if isJumpKeyPressed && ctrl.JumpTimer < ctrl.JumpHoldTime {
					playerVel.Y += ctrl.JumpBoost * 0.7 // Boost gücünü azalt
					ctrl.JumpTimer++
				} else if playerVel.Y >= 0.01 {
					ctrl.CurrentState = "falling"
				}
			} else {
				// normal skidding
				if !isJumpKeyPressed && ctrl.JumpTimer < ctrl.JumpReleaseTimer {
					playerVel.Y = ctrl.ShortJumpVelocity
					ctrl.JumpTimer = ctrl.JumpHoldTime
				} else if isJumpKeyPressed && ctrl.JumpTimer < ctrl.JumpHoldTime {
					speedFactor := (absXVelocity / ctrl.MaxRunSpeed) * ctrl.SpeedJumpFactor
					playerVel.Y += ctrl.JumpBoost * (1 + speedFactor)
					ctrl.JumpTimer++
				} else if playerVel.Y >= 0 {
					ctrl.CurrentState = "falling"
				}
			}

			// apply air skidding Decel
			if isSkidding {
				playerVel.X += math.Copysign(ctrl.AirSkiddingDecel, -playerVel.X)
			}

			// exit jumping
			if ctrl.PreviousState != ctrl.CurrentState {
				// fmt.Println("exit jumping")
			}
		case "falling":
			// enter falling
			if ctrl.PreviousState != "falling" {
				ctrl.PreviousState = ctrl.CurrentState
				ctrl.FallingDamageTempPosY = playerAABB.Pos.Y + playerAABB.Half.Y
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
				d := int(((playerAABB.Pos.Y + playerAABB.Half.Y) - ctrl.FallingDamageTempPosY) / 30)
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
				blockID := tileMapRes.GetID(gameDataRes.TargetBlockCoord.X, gameDataRes.TargetBlockCoord.Y)
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
					pos := tileMapRes.TileToWorld(gameDataRes.TargetBlockCoord)
					dropid := items.Property[blockID].DropID
					if blockID == items.OakLeaves {
						if rand.N(2) == 0 {
							dropid = items.OakLeaves
						}
					}
					toSpawnItem = append(toSpawnItem, spawnItemData{pos, dropid, 0})
				}
			}
			// transitions
			if !gameDataRes.IsRayHit || (!isBreakKeyPressed && p.isOnFloor) {
				ctrl.CurrentState = "idle"
			} else if !p.isOnFloor && playerVel.Y > 0.01 {
				ctrl.CurrentState = "falling"
			} else if !isBreakKeyPressed && isJumpKeyJustPressed {
				playerVel.Y = ctrl.JumpPower
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
				playerVel.X += math.Copysign(ctrl.SkiddingFriction, -playerVel.X)
			}
			if absXVelocity > 0.5 {
				animPlayer.SetAnim("skidding")
			}

			// Handle specific transitions
			if ctrl.SkiddingJumpEnabled && isJumpKeyJustPressed {
				// Yeni yöne doğru çok küçük sabit değerle başla
				playerVel.X = math.Copysign(0.3, float64(inputAxis.X))
				playerVel.Y = ctrl.JumpPower * 0.7 // Zıplama gücünü azalt
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
		playerVel.Y = min(ctrl.MaxFallSpeed, playerVel.Y+ctrl.Gravity)

		if !isSkidding {
			if isRunKeyPressed {
				maxSpeed, currentAccel, currentDecel = ctrl.MaxRunSpeed, ctrl.RunAcceleration, ctrl.RunDeceleration
			} else if math.Abs(playerVel.X) > ctrl.MaxWalkSpeed {
				currentDecel = ctrl.RunDeceleration
			}
		}

		targetSpeed := maxSpeed
		if inputAxis.X == 0 {
			targetSpeed = 0
		} else if inputAxis.X < 0 {
			targetSpeed = -maxSpeed
		}

		if playerVel.X < targetSpeed {
			playerVel.X = min(targetSpeed, playerVel.X+currentAccel)
		} else {
			playerVel.X = max(targetSpeed, playerVel.X-currentDecel)
		}

		// Player and tilemap collision
		*playerVel = tileCollider.Collide(*playerAABB, *playerVel, func(hitInfos []HitTileInfo, delta Vec) {
			p.isOnFloor = false
			for _, hit := range hitInfos {
				tileID := tileMapRes.GetIDUnchecked(hit.TileCoords)
				if hit.Normal.Y == -1 {
					// Ground collision
					playerVel.Y = 0
					p.isOnFloor = true
				}
				// Ceil collision
				if hit.Normal.Y == 1 {
					playerVel.Y = 0

					switch tileID {
					case items.StoneBricks:
						// Destroy block when ceil hit
						tileMapRes.Set(hit.TileCoords.X, hit.TileCoords.Y, items.Air)
						SpawnEffect(tileMapRes.TileToWorld(hit.TileCoords), tileID)
					case items.Random:
						if tileMapRes.GetIDUnchecked(hit.TileCoords.Add(tilemap.Up)) == items.Air {
							tileMapRes.Set(hit.TileCoords.X, hit.TileCoords.Y, items.Bedrock)
							ceilBlockCoord = hit.TileCoords
							ceilBlockTick = 3
						}
					}
				}
				// Right or Left wall collision
				if hit.Normal.X == -1 || hit.Normal.X == 1 {
					// While running at maximum speed, hold down the right arrow key and hit the block to destroy it.
					if absXVelocity == ctrl.MaxRunSpeed && isBreakKeyPressed {
						tileMapRes.Set(hit.TileCoords.X, hit.TileCoords.Y, items.Air)
						SpawnEffect(tileMapRes.TileToWorld(hit.TileCoords), tileID)
					}
					playerVel.X = 0
				}
			}
		},
		)

		// Platform collision
		playerAABB.Pos = playerAABB.Pos.Add(*playerVel)
		hit := &HitInfo2{}
		pq := filterPlatform.Query()
		for pq.Next() {

			platformAABB, platformVel, platformType := pq.Get()

			if AABBPlatform(playerAABB, platformAABB, playerVel, (*Vec)(platformVel), hit) {
				if hit.Top {
					if *platformType == "solid" {
						playerAABB.Pos = playerAABB.Pos.Add(hit.Delta)
						playerVel.Y = 0
					}
				}
				if hit.Bottom {
					playerAABB.Pos = playerAABB.Pos.Add(hit.Delta)
					playerVel.Y = 0
					pvel := tileCollider.Collide(*playerAABB, Vec(*platformVel), nil)
					playerAABB.Pos = playerAABB.Pos.Add(pvel)
					p.isOnFloor = true
				}
				if hit.Right {
					playerAABB.Pos = playerAABB.Pos.Add(hit.Delta)
					playerVel.X = -1.01
				}
				if hit.Left {
					playerAABB.Pos = playerAABB.Pos.Add(hit.Delta)
					playerVel.X = +1.01
				}
			}
		}

		// Enemy collision
		hit.Reset()
		eq := filterEnemy.Query()
		for eq.Next() {
			enemyAABB, enemyVel, mobileID := eq.Get()
			if AABBPlatform(playerAABB, enemyAABB, playerVel, (*Vec)(enemyVel), hit) {
				if hit.Top {
					switch *mobileID {
					case WormID:
						playerAABB.Pos = playerAABB.Pos.Add(hit.Delta)
						playerVel.Y = 0
						if elapsed == 0 {
							pHealth.Current -= 5
							elapsed = 2 * 60
						}
					}
				}
				if hit.Bottom {
					playerAABB.Pos = playerAABB.Pos.Add(hit.Delta)
					playerVel.Y = 0
					pvel := tileCollider.Collide(*playerAABB, Vec(*enemyVel), nil)
					playerAABB.Pos = playerAABB.Pos.Add(pvel)
					p.isOnFloor = true
					toRemove = append(toRemove, eq.Entity())
				}
				if hit.Right {
					playerAABB.Pos = playerAABB.Pos.Add(hit.Delta)
					playerVel.X = -1.01
					if elapsed == 0 {
						pHealth.Current -= 5
						elapsed = 2 * 60
					}
				}
				if hit.Left {
					playerAABB.Pos = playerAABB.Pos.Add(hit.Delta)
					playerVel.X = +1.01
					if elapsed == 0 {
						pHealth.Current -= 5
						elapsed = 2 * 60
					}
				}
			}
		}

		// player facing raycast for target block
		p.playerTile = tileMapRes.WorldToTile(playerAABB.Pos.X, playerAABB.Pos.Y)
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
				// check place block overlaps
				queryItem := filterDroppedItem.Query()
				for queryItem.Next() {
					dropItemAABB.Pos = *(*Vec)(mapPos.GetUnchecked(queryItem.Entity()))
					anyItemOverlapsWithPlaceRect = Overlap(placeTileBox, dropItemAABB, nil)
					if anyItemOverlapsWithPlaceRect {
						queryItem.Close()
						break
					}
				}
				if !anyItemOverlapsWithPlaceRect {
					if !Overlap(playerAABB, placeTileBox, nil) {
						// place block
						tileMapRes.Set(p.placeTile.X, p.placeTile.Y, inventoryRes.CurrentSlotID())
						// remove item
						inventoryRes.RemoveItemFromSelectedSlot()
					}
				}
				// if slot item snowball, throw snowball
			} else if inventoryRes.CurrentSlot().ID == items.Snowball {
				if ctrl.CurrentState != "skidding" {
					spawnPos := Vec{playerAABB.Pos.X, playerAABB.Pos.Y - 4}
					spawnVel := Vec{SnowballSpeedX, SnowballMaxFallVelocity}
					switch *pFacing {
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
				toSpawnItem = append(toSpawnItem, spawnItemData{
					playerAABB.Pos,
					currentSlot.ID,
					currentSlot.Durability,
				})
				inventoryRes.RemoveItemFromSelectedSlot()
				onInventorySlotChanged()
			}
		}

	}
}
func (c *Player) Draw() {
	// Draw player
	if world.Alive(currentPlayer) {
		playerBox := mapAABB.GetUnchecked(currentPlayer)

		colorMDIO.GeoM.Reset()
		x := playerBox.Pos.X - playerBox.Half.X
		y := playerBox.Pos.Y - playerBox.Half.Y
		if mapFacing.GetUnchecked(currentPlayer).X == -1 {
			colorMDIO.GeoM.Scale(-1, 1)
			colorMDIO.GeoM.Translate(playerBox.Pos.X+playerBox.Half.X, y)
		} else {
			colorMDIO.GeoM.Translate(x, y)
		}
		if blinking {
			colorMDIO.ColorScale.SetA(0.2)
		} else {
			colorMDIO.ColorScale.SetA(1)
		}
		cameraRes.DrawWithColorM(animPlayer.CurrentFrame, colorM, colorMDIO, Screen)

	}
}
