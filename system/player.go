package system

import (
	"image"
	"kar"
	"kar/arc"
	"kar/items"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/tilecollider"
)

const (
	ballGravity         = 0.5
	ballXspeed          = 3.5
	ballMaxFallVelocity = 2.5
	ballBounceHeight    = 9
)

var (
	blockHealth                  float64
	targetBlockPos               image.Point
	placeBlock                   image.Point
	playerTile                   image.Point
	playerCenterX, playerCenterY float64
	isRayHit                     bool
	bounceVelocity               = -math.Sqrt(2 * ballGravity * ballBounceHeight)
)

func (c *Player) Init() {
	anim, hlt, rect := arc.MapPlayer.Get(player)
	ctrl.AnimPlayer = anim
	ctrl.Rect = rect
	ctrl.Health = hlt
	ctrl.Inventory = items.NewInventory()
	ctrl.Inventory.SetSlot(0, items.Snowball, 64, 0)
	ctrl.Inventory.SetSlot(5, items.DiamondPickaxe, 1, 0)
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
			targetBlockTemp := targetBlockPos
			targetBlockPos, isRayHit = tileMap.Raycast(
				playerTile,
				ctrl.AxisLast,
				kar.RaycastDist,
			)
			// reset attack if block focus changed
			if !targetBlockPos.Eq(targetBlockTemp) || !isRayHit {
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
			}

			// Place block
			if ctrl.IsAttackKeyJustPressed {
				anyItemOverlapsWithPlaceCoords := false
				if isRayHit && items.HasTag(ctrl.Inventory.CurrentSlot().ID, items.Block) {
					placeBlock = targetBlockPos.Sub(ctrl.AxisLast)
					queryItem := arc.FilterItem.Query(&kar.WorldECS)
					for queryItem.Next() {
						_, itemRect, _, _ := queryItem.Get()
						anyItemOverlapsWithPlaceCoords = itemRect.Overlaps(
							tileMap.GetTileRect(placeBlock),
						)
						if anyItemOverlapsWithPlaceCoords {
							queryItem.Close()
							break
						}
					}
					if !anyItemOverlapsWithPlaceCoords {
						if !ctrl.Rect.Overlaps(tileMap.GetTileRect(placeBlock)) {
							tileMap.Set(placeBlock.X, placeBlock.Y, ctrl.Inventory.CurrentSlotID())
							ctrl.Inventory.RemoveItemFromSelectedSlot()
						}
					}
				} else if ctrl.Inventory.CurrentSlot().ID == items.Snowball {
					switch ctrl.AxisLast {
					case image.Point{1, 0}:
						arc.SpawnSnowBall(playerCenterX, playerCenterY-4, ballXspeed, ballMaxFallVelocity)
					case image.Point{-1, 0}:
						arc.SpawnSnowBall(playerCenterX, playerCenterY-4, -ballXspeed, ballMaxFallVelocity)
					}
				}

			}

			// snowball physics
			q := arc.FilterMapSnowBall.Query(&kar.WorldECS)
			for q.Next() {
				_, rect, v := q.Get()
				v.VelY += ballGravity
				v.VelY = min(v.VelY, ballMaxFallVelocity)
				collider.Collide(rect.X, rect.Y, rect.W, rect.H, v.VelX, v.VelY,
					func(ci []tilecollider.CollisionInfo[uint16], dx, dy float64) {
						rect.X += dx
						rect.Y += dy
						isHorizontalCollision := false
						for _, c := range ci {
							if c.Normal[1] == -1 {
								v.VelY = bounceVelocity
							}
							if c.Normal[0] == -1 && v.VelX > 0 && v.VelY > 0 {
								isHorizontalCollision = true
							}
							if c.Normal[0] == 1 && v.VelX < 0 && v.VelY > 0 {
								isHorizontalCollision = true
							}
						}
						if isHorizontalCollision {
							if kar.WorldECS.Alive(q.Entity()) {
								toRemove = append(toRemove, q.Entity())
							}
						}
					})
			}

			// Remove dead player entity
			if ctrl.Health.Health <= 0 {
				toRemove = append(toRemove, player)
			}

		}

	}

}

func (c *Player) Draw() {

}
