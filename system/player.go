package system

import (
	"image"
	"kar"
	"kar/arc"
	"kar/items"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/anim"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/tilecollider"
)

const (
	ballGravity         = 0.5
	ballSpeedX          = 3.5
	ballMaxFallVelocity = 2.5
	ballBounceHeight    = 9
)

var (
	bounceVelocity = -math.Sqrt(2 * ballGravity * ballBounceHeight)
	blockHealth    float64
	targetTile     image.Point
	placeTile      image.Point
	playerTile     image.Point
	// floorTile                    image.Point
	playerCenterX, playerCenterY float64
	isRayHit                     bool

	AnimationPlayerHand *anim.AnimationPlayer
)

func (c *Player) Init() {
	ap, hlt, rect := arc.MapPlayer.Get(player)
	ctrl.AnimPlayer = ap
	AnimationPlayerHand = ap
	ctrl.Rect = rect
	ctrl.Health = hlt
	ctrl.Inventory = items.NewInventory()
	ctrl.Inventory.SetSlot(0, items.Snowball, 64, 0)
	// ctrl.Inventory.SetSlot(1, items.IronPickaxe, 1, items.GetDefaultDurability(items.DiamondPickaxe))
	// ctrl.Inventory.SetSlot(2, items.IronAxe, 1, items.GetDefaultDurability(items.DiamondAxe))
	// ctrl.Inventory.SetSlot(3, items.IronShovel, 1, items.GetDefaultDurability(items.DiamondShovel))
	ctrl.EnterFalling()
}

type Player struct {
}

func (c *Player) Update() {
	if kar.WorldECS.Alive(player) {

		if !craftingState {
			ctrl.UpdateInput()
			ctrl.UpdateState()
			ctrl.UpdatePhysics()
			playerCenterX = ctrl.Rect.X + ctrl.Rect.W/2
			playerCenterY = ctrl.Rect.Y + ctrl.Rect.H/2
			playerTile = tileMap.WorldToTile(playerCenterX, playerCenterY)
			targetBlockTemp := targetTile
			targetTile, isRayHit = tileMap.Raycast(
				playerTile,
				ctrl.AxisLast,
				kar.RaycastDist,
			)
			// reset attack if block focus changed
			if !targetTile.Eq(targetBlockTemp) || !isRayHit {
				blockHealth = 0
			}

			// Drop Item
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
				currentSlot := ctrl.Inventory.CurrentSlot()
				if currentSlot.ID != items.Air {
					AppendToSpawnList(
						playerCenterX,
						playerCenterY,
						currentSlot.ID,
						currentSlot.Durability,
					)
					ctrl.Inventory.RemoveItemFromSelectedSlot()
				}
				onInventorySlotChanged()
			}

			// PLACE BLOCK
			if ctrl.IsAttackKeyJustPressed {
				anyItemOverlapsWithPlaceCoords := false
				// if slot item is block
				if isRayHit && items.HasTag(ctrl.Inventory.CurrentSlot().ID, items.Block) {
					placeTile = targetTile.Sub(ctrl.AxisLast)
					queryItem := arc.FilterItem.Query(&kar.WorldECS)
					// check overlaps
					for queryItem.Next() {
						_, itemRect, _, _ := queryItem.Get()
						anyItemOverlapsWithPlaceCoords = itemRect.Overlaps2(
							tileMap.GetTileRect(placeTile.X, placeTile.Y),
						)
						if anyItemOverlapsWithPlaceCoords {
							queryItem.Close()
							break
						}
					}
					if !anyItemOverlapsWithPlaceCoords {
						if !ctrl.Rect.Overlaps2(tileMap.GetTileRect(placeTile.X, placeTile.Y)) {
							// place block
							tileMap.Set(placeTile.X, placeTile.Y, ctrl.Inventory.CurrentSlotID())
							// remove item
							ctrl.Inventory.RemoveItemFromSelectedSlot()
						}
					}
					// if slot item snowball, spawn snowball
				} else if ctrl.Inventory.CurrentSlot().ID == items.Snowball {
					ctrl.Inventory.RemoveItemFromSelectedSlot()
					switch ctrl.AxisLast {
					case image.Point{1, 0}:
						arc.SpawnSnowBall(playerCenterX, playerCenterY-4, ballSpeedX, ballMaxFallVelocity)
					case image.Point{-1, 0}:
						arc.SpawnSnowBall(playerCenterX, playerCenterY-4, -ballSpeedX, ballMaxFallVelocity)
					}
				}

			}

			// snowball physics
			q := arc.FilterMapSnowBall.Query(&kar.WorldECS)
			for q.Next() {
				_, rect, v := q.Get()
				v.Y += ballGravity
				v.Y = min(v.Y, ballMaxFallVelocity)
				collider.Collide(
					rect.X,
					rect.Y,
					rect.W,
					rect.H,
					v.X,
					v.Y,
					func(ci []tilecollider.CollisionInfo[uint16], dx, dy float64) {
						rect.X += dx
						rect.Y += dy
						isHorizontalCollision := false
						for _, c := range ci {
							if c.Normal[1] == -1 {
								v.Y = bounceVelocity
							}
							if c.Normal[0] == -1 && v.X > 0 && v.Y > 0 {
								isHorizontalCollision = true
							}
							if c.Normal[0] == 1 && v.X < 0 && v.Y > 0 {
								isHorizontalCollision = true
							}
						}
						if isHorizontalCollision {
							if kar.WorldECS.Alive(q.Entity()) {
								toRemove = append(toRemove, q.Entity())
							}
						}
					},
				)
			}
			// Remove dead player entity
			if ctrl.Health.Current <= 0 {
				kar.Camera.ShakeEnabled = true
				kar.Camera.SmoothType = kamera.Lerp
				kar.Camera.AddTrauma(1)
				toRemove = append(toRemove, player)
			}
		}

	}
}

func (c *Player) Draw() {
	if kar.WorldECS.Alive(player) {
		// Draw player
		kar.ColorMDIO.GeoM.Reset()
		kar.ColorMDIO.GeoM.Scale(ctrl.FlipXFactor, 1)
		if ctrl.FlipXFactor == -1 {
			kar.ColorMDIO.GeoM.Translate(ctrl.Rect.X+ctrl.Rect.W, ctrl.Rect.Y)
		} else {
			kar.ColorMDIO.GeoM.Translate(ctrl.Rect.X, ctrl.Rect.Y)
		}

		kar.Camera.DrawWithColorM(ctrl.AnimPlayer.CurrentFrame, kar.ColorM, kar.ColorMDIO, kar.Screen)
	}
}
