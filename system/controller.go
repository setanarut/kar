package system

import (
	"fmt"
	"image"
	"kar"
	"kar/arc"
	"kar/engine/mathutil"
	"kar/items"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/anim"
	"github.com/setanarut/tilecollider"
)

type Controller struct {
	AnimPlayer *anim.AnimationPlayer
	Health     *arc.Health
	Inventory  *items.Inventory
	Rect       *arc.Rect

	fallingDamageTempPosY float64

	CurrentState string
	Collider     *tilecollider.Collider[uint16]

	VelX                                float64
	VelY                                float64
	JumpPower                           float64
	Gravity                             float64
	MaxFallSpeed                        float64
	MaxRunSpeed                         float64
	MaxWalkSpeed                        float64
	Acceleration                        float64
	Deceleration                        float64
	JumpHoldTime                        float64
	JumpBoost                           float64
	JumpTimer                           float64
	MinSpeedThresForJumpBoostMultiplier float64 // Yüksek zıplama için gereken minimum hız
	JumpBoostMultiplier                 float64 // Yüksek zıplamada kullanılacak çarpan
	SpeedJumpFactor                     float64 // Yatay hızın zıplama yüksekliğine etkisini kontrol eden çarpan
	ShortJumpVelocity                   float64 // Kısa zıplama için hız
	JumpReleaseTimer                    float64 // Zıplama tuşu bırakıldığında geçen süre

	IsOnFloor   bool
	IsSkidding  bool
	IsFalling   bool
	FlipXFactor float64

	SkiddingJumpEnabled bool

	// Input durumları
	IsBreakKeyPressed      bool
	IsAttackKeyJustPressed bool
	IsJumpKeyPressed       bool
	IsJumpKeyJustPressed   bool
	IsRunKeyPressed        bool
	InputAxis              image.Point
	AxisLast               image.Point

	WalkAcceleration float64
	WalkDeceleration float64
	RunAcceleration  float64
	RunDeceleration  float64

	HorizontalVelocity float64
	// Durum değişikliği için yeni alan
	previousState string
}

func NewController(velX, velY float64, tc *tilecollider.Collider[uint16]) *Controller {
	return &Controller{
		CurrentState:                        "falling",
		Collider:                            tc,
		VelX:                                velX,
		VelY:                                velY,
		JumpPower:                           -3.7,
		Gravity:                             0.19,
		MaxFallSpeed:                        100.0,
		Acceleration:                        0.08,
		Deceleration:                        0.1,
		JumpHoldTime:                        20.0,
		JumpBoost:                           -0.1,
		MinSpeedThresForJumpBoostMultiplier: 0.1,
		JumpBoostMultiplier:                 1.01,
		SpeedJumpFactor:                     0.3,
		ShortJumpVelocity:                   -2.0,
		JumpReleaseTimer:                    5,
		MaxWalkSpeed:                        1.6,
		MaxRunSpeed:                         3.0,
		WalkAcceleration:                    0.04,
		WalkDeceleration:                    0.04,
		RunAcceleration:                     0.04,
		RunDeceleration:                     0.04,

		FlipXFactor: 1,
	}

}

func (c *Controller) UpdateInput() {
	c.IsBreakKeyPressed = ebiten.IsKeyPressed(ebiten.KeyRight)
	c.IsRunKeyPressed = ebiten.IsKeyPressed(ebiten.KeyShift)
	c.IsJumpKeyPressed = ebiten.IsKeyPressed(ebiten.KeySpace)
	c.IsAttackKeyJustPressed = inpututil.IsKeyJustPressed(ebiten.KeyLeft)
	c.IsJumpKeyJustPressed = inpututil.IsKeyJustPressed(ebiten.KeySpace)
	c.InputAxis = image.Point{}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		c.InputAxis.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		c.InputAxis.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		c.InputAxis.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		c.InputAxis.X += 1
	}
	if !c.InputAxis.Eq(image.Point{}) {
		c.AxisLast = c.InputAxis
	}
}

func (c *Controller) UpdatePhysics() {
	maxSpeed := c.MaxWalkSpeed
	currentAccel := c.WalkAcceleration
	currentDecel := c.WalkDeceleration
	c.HorizontalVelocity = math.Abs(c.VelX)

	c.VelY += c.Gravity
	c.VelY = min(c.MaxFallSpeed, c.VelY)

	// Enemy collisions
	enemyQuery := arc.FilterEnemy.Query(&kar.WorldECS)
	for enemyQuery.Next() {
		rect, _, _ := enemyQuery.Get()
		collInfo := c.Rect.CheckCollision(rect, ctrl.VelX, ctrl.VelY)

		c.VelX += collInfo.DeltaX
		c.VelY += collInfo.DeltaY

		if collInfo.Collided {
			switch collInfo.Normal[0] {
			case 1:
				fmt.Println("left collide")
				c.VelX += 3
				c.Health.Health -= 8
			case -1:
				fmt.Println("right collide")
				c.VelX -= 3
				c.Health.Health -= 8
			}
			switch collInfo.Normal[1] {
			case 1:
				c.VelY = 0
				fmt.Println("floor collide")
			case -1:
				c.VelY = -5
				fmt.Println("ceil collide")
			}
		}
	}

	if !c.IsSkidding {
		if c.IsRunKeyPressed {
			maxSpeed = c.MaxRunSpeed
			currentAccel = c.RunAcceleration
			currentDecel = c.RunDeceleration
		} else if c.HorizontalVelocity > c.MaxWalkSpeed {
			currentDecel = c.RunDeceleration
		}
	}

	if c.InputAxis.X > 0 {
		if c.VelX > maxSpeed {
			c.VelX = max(maxSpeed, c.VelX-currentDecel)
		} else {
			c.VelX = min(maxSpeed, c.VelX+currentAccel)
		}
	} else if c.InputAxis.X < 0 {
		if c.VelX < -maxSpeed {
			c.VelX = min(-maxSpeed, c.VelX+currentDecel)
		} else {
			c.VelX = max(-maxSpeed, c.VelX-currentAccel)
		}
	} else {
		if c.VelX > 0 {
			c.VelX = max(0, c.VelX-currentDecel)
		} else if c.VelX < 0 {
			c.VelX = min(0, c.VelX+currentDecel)
		}
	}

	c.IsSkidding = (c.VelX > 0 && c.InputAxis.X == -1) || (c.VelX < 0 && c.InputAxis.X == 1)

	if c.VelX > 0.01 {
		c.FlipXFactor = 1 // sağa gidiyor
		c.AxisLast.X = 1
	} else if c.VelX < -0.01 {
		c.FlipXFactor = -1 // sola gidiyor
		c.AxisLast.X = -1
	}

	// Player and tilemap collision
	c.Collider.Collide(math.Round(c.Rect.X), c.Rect.Y, c.Rect.W, c.Rect.H, c.VelX, c.VelY, c.HandleCollision)
}

func (c *Controller) HandleCollision(collisionInfos []tilecollider.CollisionInfo[uint16], dx, dy float64) {
	c.IsOnFloor = false

	// Apply tilemap collision response
	c.Rect.X += dx
	c.Rect.Y += dy

	// Reset velocity when collide
	for _, collisionInfo := range collisionInfos {
		if collisionInfo.Normal[1] == -1 {
			c.VelY = 0
			c.IsOnFloor = true // on floor collision
		}
		if collisionInfo.Normal[1] == 1 {
			c.ResetVelocityX()
		}
		if collisionInfo.Normal[0] == -1 {
			c.ResetVelocityX()
		}
		if collisionInfo.Normal[0] == 1 {
			c.ResetVelocityX()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		ids := make([]uint16, 0)
		for _, collisionInfo := range collisionInfos {
			if collisionInfo.Normal[1] == -1 {
				ids = append(ids, collisionInfo.TileID)
			}
		}
		if len(ids) == 2 {
			if ids[0] == items.Sand && ids[1] == items.GrassBlock {
				fmt.Println(ids[0], items.DisplayName(ids[0]), ids[1], items.DisplayName(ids[1]))
				// Pipe logic is here
			}
		}
	}

}

func (c *Controller) ResetVelocityX() {
	c.VelX = 0
	c.HorizontalVelocity = 0
}

func (c *Controller) Skidding() {
	if c.SkiddingJumpEnabled && c.IsJumpKeyJustPressed {
		c.ResetVelocityX()

		// Yeni yöne doğru çok küçük sabit değerle başla
		if c.InputAxis.X > 0 {
			c.VelX = 0.3
		} else if c.InputAxis.X < 0 {
			c.VelX = -0.3
		}

		c.VelY = c.JumpPower * 0.7 // Zıplama gücünü azalt
		c.JumpTimer = 0
		c.ChangeState("jumping")
		return
	}

	// Mevcut mantık devam eder...
	if c.HorizontalVelocity < 0.01 {
		c.ChangeState("idle")
	} else if !c.IsSkidding {
		if c.HorizontalVelocity > c.MaxWalkSpeed {
			c.ChangeState("running")
		} else {
			c.ChangeState("walking")
		}
	}
}

func (c *Controller) Falling() {
	if c.VelY > 0.1 {
		c.AnimPlayer.SetStateAndReset("jump")
	}
	if c.IsOnFloor {
		if c.HorizontalVelocity <= 0 {
			c.ChangeState("idle")
		} else if c.IsRunKeyPressed {
			c.ChangeState("running")
		} else {
			c.ChangeState("walking")
		}
	}
}

func (c *Controller) Breaking() {

	// set animation states
	if c.AxisLast.X == 1 {
		if c.HorizontalVelocity > 0.01 {
			c.AnimPlayer.SetStateAndReset("attackWalk")
		} else {
			c.AnimPlayer.SetStateAndReset("attackRight")
		}
	} else if c.AxisLast.X == -1 {
		if c.HorizontalVelocity > 0.01 {
			c.AnimPlayer.SetStateAndReset("attackWalk")
		} else {
			c.AnimPlayer.SetStateAndReset("attackRight")
		}
		c.FlipXFactor = -1
	} else if c.AxisLast.Y == 1 {
		c.AnimPlayer.SetStateAndReset("attackDown")
	} else if c.AxisLast.Y == -1 {
		c.AnimPlayer.SetStateAndReset("attackUp")
	}

	if isRayHit {
		blockID := tileMap.Get(targetTile.X, targetTile.Y)
		if !items.HasTag(blockID, items.Unbreakable) {
			if items.IsBestTool(blockID, c.Inventory.CurrentSlotID()) {
				blockHealth += kar.PlayerBestToolDamage
			} else {
				blockHealth += kar.PlayerDefaultDamage
			}
		}
		// Destroy block
		if blockHealth >= 180 {
			blockHealth = 0
			tileMap.Set(targetTile.X, targetTile.Y, items.Air)

			if items.HasTag(c.Inventory.CurrentSlotID(), items.Tool) {
				c.Inventory.CurrentSlot().Durability--
				if c.Inventory.CurrentSlot().Durability <= 0 {
					c.Inventory.ClearCurrentSlot()
				}
			}

			// spawn drop item
			x, y := tileMap.TileToWorldCenter(targetTile.X, targetTile.Y)
			AppendToSpawnList(x, y, items.Property[blockID].DropID, 0)
		}
	}

	if !isRayHit {
		c.ChangeState("idle")
	}

	if !c.IsOnFloor && c.VelY > 0.01 {
		c.ChangeState("falling")
	} else if !c.IsBreakKeyPressed && c.IsOnFloor {
		c.ChangeState("idle")
	} else if !c.IsBreakKeyPressed && c.IsJumpKeyJustPressed {
		c.ChangeState("jumping")
		c.VelY = c.JumpPower
		c.JumpTimer = 0
	}
}

func (c *Controller) Jumping() {
	if c.VelY != 0 && c.VelY > c.JumpPower+0.1 {
		c.AnimPlayer.SetStateAndReset("jump")
	}
	// Skidding'den geldiyse özel durum
	if c.previousState == "skidding" {
		if !c.IsJumpKeyPressed && c.JumpTimer < c.JumpReleaseTimer {
			c.VelY = c.ShortJumpVelocity * 0.7 // Kısa zıplama gücünü azalt
			c.JumpTimer = c.JumpHoldTime
		} else if c.IsJumpKeyPressed && c.JumpTimer < c.JumpHoldTime {
			c.VelY += c.JumpBoost * 0.7 // Boost gücünü azalt
			c.JumpTimer++
		} else if c.VelY >= 0.01 {
			c.ChangeState("falling")
		}
	} else {
		// Normal jumping mantığı aynen devam eder
		if !c.IsJumpKeyPressed && c.JumpTimer < c.JumpReleaseTimer {
			c.VelY = c.ShortJumpVelocity
			c.JumpTimer = c.JumpHoldTime
		} else if c.IsJumpKeyPressed && c.JumpTimer < c.JumpHoldTime {
			speedFactor := (c.HorizontalVelocity / c.MaxRunSpeed) * c.SpeedJumpFactor
			c.VelY += c.JumpBoost * (1 + speedFactor)
			c.JumpTimer++
		} else if c.VelY >= 0 {
			c.ChangeState("falling")
		}
	}

	// Yatay hareket kontrolü
	if c.InputAxis.X < 0 && c.VelX > 0 {
		c.VelX -= c.Deceleration
	} else if c.InputAxis.X > 0 && c.VelX < 0 {
		c.VelX += c.Deceleration
	}
}

func (c *Controller) Running() {
	c.AnimPlayer.Animations["walkRight"].FPS = mathutil.MapRange(c.HorizontalVelocity, 0, c.MaxRunSpeed, 4, 23)

	// Kayma durumu kontrolü
	if c.IsSkidding {
		c.ChangeState("skidding")
		return
	}

	if c.VelY > 0 && !c.IsOnFloor {
		c.ChangeState("falling")
	}

	if c.IsJumpKeyJustPressed {
		c.ChangeState("jumping")
		if c.HorizontalVelocity > c.MinSpeedThresForJumpBoostMultiplier {
			c.VelY = c.JumpPower * c.JumpBoostMultiplier
		} else {
			c.VelY = c.JumpPower
		}
		c.JumpTimer = 0
	} else if c.HorizontalVelocity < 0.01 {
		c.ChangeState("idle")
	} else if c.HorizontalVelocity <= c.MaxWalkSpeed {
		c.ChangeState("walking")
	}
}

func (c *Controller) Walking() {
	c.AnimPlayer.Animations["walkRight"].FPS = mathutil.MapRange(c.HorizontalVelocity, 0, c.MaxRunSpeed, 4, 23)
	// Kayma durumu kontrolü
	if c.IsSkidding {
		c.ChangeState("skidding")
		return
	}

	if c.VelY > 0 && !c.IsOnFloor {
		c.ChangeState("falling")
	}

	if c.IsJumpKeyJustPressed {
		c.ChangeState("jumping")
		if c.HorizontalVelocity > c.MinSpeedThresForJumpBoostMultiplier {
			c.VelY = c.JumpPower * c.JumpBoostMultiplier
		} else {
			c.VelY = c.JumpPower
		}
		c.JumpTimer = 0
	} else if c.HorizontalVelocity <= 0 {
		c.ChangeState("idle")
	} else if c.HorizontalVelocity > c.MaxWalkSpeed {
		c.ChangeState("running")
	}
}

func (c *Controller) Idle() {

	if c.AxisLast.Y == -1 {
		c.AnimPlayer.SetStateAndReset("idleUp")
	} else if c.AxisLast.Y == 1 {
		c.AnimPlayer.SetStateAndReset("idleDown")
	}
	if c.IsJumpKeyJustPressed {
		c.VelY = c.JumpPower
		c.JumpTimer = 0
		if c.HorizontalVelocity > c.MinSpeedThresForJumpBoostMultiplier {
			c.VelY = c.JumpPower * c.JumpBoostMultiplier
		} else {
			c.VelY = c.JumpPower
		}
		c.JumpTimer = 0
		// c.changeState("jumping")
	} else if c.IsOnFloor && c.HorizontalVelocity > 0.01 {
		if c.HorizontalVelocity > c.MaxWalkSpeed {
			c.ChangeState("running")
		} else {
			c.ChangeState("walking")
		}
	} else if !c.IsOnFloor && c.VelY > 0.01 {
		c.ChangeState("falling")
	} else if c.IsBreakKeyPressed && isRayHit {
		c.ChangeState("breaking")
	}

	if c.VelY != 0 && c.VelY < -0.1 {
		c.ChangeState("jumping")
	}
}

func (c *Controller) UpdateState() {
	switch c.CurrentState {
	case "idle":
		c.Idle()
	case "walking":
		c.Walking()
	case "running":
		c.Running()
	case "jumping":
		c.Jumping()
	case "falling":
		c.Falling()
	case "breaking":
		c.Breaking()
	case "skidding":
		c.Skidding()
	}
}

// func (c *Controller) exitRunning()  {}
// func (c *Controller) exitJumping()  {}
// func (c *Controller) exitFalling()  {}

func (c *Controller) EnterWalking() {
	c.AnimPlayer.SetStateAndReset("walkRight")
}
func (c *Controller) EnterRunning() {
	c.AnimPlayer.SetStateAndReset("walkRight")
}

func (c *Controller) EnterIdle() {
	if c.AxisLast.Y == 0 {
		c.AnimPlayer.SetStateAndReset("idleRight")
	}
	if c.AxisLast.X == 0 {
		c.AnimPlayer.SetStateAndReset("idleUp")
	}
}

// func (c *Controller) enterAttacking() {
// }

func (c *Controller) ExitBreaking() {
	blockHealth = 0
}

// func (c *Controller) enterJumping() {

// }

func (c *Controller) EnterFalling() {
	c.fallingDamageTempPosY = c.Rect.Y
}

func (c *Controller) ExitFalling() {
	d := int((c.Rect.Y - c.fallingDamageTempPosY) / 30)
	if d > 3 {
		c.Health.Health -= d - 3
	}
}

func (c *Controller) EnterSkidding() {
	c.AnimPlayer.SetStateAndReset("skidding")
}

func (c *Controller) ChangeState(newState string) {
	if c.CurrentState == newState {
		return
	}

	// Mevcut durumdan çık
	switch c.CurrentState {
	case "breaking":
		c.ExitBreaking()
	// case "idle":
	// c.exitIdle()
	// case "walking":
	// 	c.exitWalking()
	// case "running":
	// 	c.exitRunning()
	// case "jumping":
	// 	c.exitJumping()
	case "falling":
		c.ExitFalling()
	}

	c.previousState = c.CurrentState
	c.CurrentState = newState

	// Yeni duruma gir
	switch newState {
	case "idle":
		c.EnterIdle()
	case "walking":
		c.EnterWalking()
	case "running":
		c.EnterRunning()
	case "falling":
		c.EnterFalling()
	case "skidding":
		c.EnterSkidding()
		// case "breaking":
		// 	c.enterBreaking()
		// case "jumping":
		// 	c.enterJumping()
	}
}
