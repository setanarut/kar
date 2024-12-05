package arc

import (
	"kar/engine/mathutil"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/anim"
)

type Controller struct {
	CurrentState                        string
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

	IsOnFloor bool
	// Input durumları
	IsLeftKeyPressed     bool
	IsRightKeyPressed    bool
	IsJumpKeyPressed     bool
	IsJumpKeyJustPressed bool
	IsRunKeyPressed      bool
}

func NewController(velX, velY float64) *Controller {
	return &Controller{
		CurrentState: "falling",

		VelX:                                velX,
		VelY:                                velY,
		JumpPower:                           -3.7,
		Gravity:                             0.19,
		MaxFallSpeed:                        6.0,
		MaxRunSpeed:                         2.5,
		MaxWalkSpeed:                        2.0,
		Acceleration:                        0.08,
		Deceleration:                        0.1,
		JumpHoldTime:                        20.0,
		JumpBoost:                           -0.1,
		MinSpeedThresForJumpBoostMultiplier: 0.1,
		JumpBoostMultiplier:                 1.01,
		SpeedJumpFactor:                     0.3,
		ShortJumpVelocity:                   -2.0,
		JumpReleaseTimer:                    5,
	}
}

func (c *Controller) UpdateInput() {
	c.IsRunKeyPressed = ebiten.IsKeyPressed(ebiten.KeyShift)
	c.IsLeftKeyPressed = ebiten.IsKeyPressed(ebiten.KeyA)
	c.IsRightKeyPressed = ebiten.IsKeyPressed(ebiten.KeyD)
	c.IsJumpKeyPressed = ebiten.IsKeyPressed(ebiten.KeySpace)
	c.IsJumpKeyJustPressed = inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

func (c *Controller) UpdatePhysics() {
	maxSpeed := c.MaxWalkSpeed
	if c.IsRunKeyPressed {
		maxSpeed = c.MaxRunSpeed
	}

	if c.IsRightKeyPressed {
		if c.VelX < maxSpeed {
			c.VelX += c.Acceleration
		}
	} else if c.IsLeftKeyPressed {
		if c.VelX > -maxSpeed {
			c.VelX -= c.Acceleration
		}
	} else {
		if c.VelX > 0 {
			c.VelX = max(0, c.VelX-c.Deceleration)
		} else if c.VelX < 0 {
			c.VelX = min(0, c.VelX+c.Deceleration)
		}
	}

	c.VelY += c.Gravity
	if c.VelY > c.MaxFallSpeed {
		c.VelY = c.MaxFallSpeed
	}
}

func (c *Controller) UpdateState(anim *anim.AnimationPlayer) {
	switch c.CurrentState {
	case "idle":
		anim.SetState("idleRight")
		if c.IsJumpKeyJustPressed {
			c.CurrentState = "jumping"
			c.VelY = c.JumpPower
			c.JumpTimer = 0
		} else if c.VelX != 0 {
			if c.IsRunKeyPressed {
				c.CurrentState = "running"
			} else {
				c.CurrentState = "walking"
			}
		}

	case "walking":
		anim.SetState("walkRight")
		fps := mathutil.MapRange(math.Abs(c.VelX), 0, c.MaxRunSpeed, 0, 10)
		anim.SetStateFPS("walkRight", fps)
		if c.VelY > 0 && !c.IsOnFloor {
			c.CurrentState = "falling"
		}

		// Koşma hızından yürüme hızına kademeli geçiş
		if math.Abs(c.VelX) > c.MaxWalkSpeed {
			if c.VelX > 0 {
				c.VelX = math.Max(c.MaxWalkSpeed, c.VelX-c.Deceleration)
			} else {
				c.VelX = math.Min(-c.MaxWalkSpeed, c.VelX+c.Deceleration)
			}
		}

		if c.IsJumpKeyJustPressed {
			c.CurrentState = "jumping"
			if math.Abs(c.VelX) > c.MinSpeedThresForJumpBoostMultiplier {
				c.VelY = c.JumpPower * c.JumpBoostMultiplier
			} else {
				c.VelY = c.JumpPower
			}
			c.JumpTimer = 0
		} else if math.Abs(c.VelX) <= 0 {
			c.CurrentState = "idle"
		}

		if c.IsRunKeyPressed {
			c.CurrentState = "running"
		}

	case "running":
		anim.SetState("walkRight")
		fps := mathutil.MapRange(math.Abs(c.VelX), 0, c.MaxRunSpeed, 0, 10)
		anim.SetStateFPS("walkRight", fps)
		if c.IsJumpKeyJustPressed {
			c.CurrentState = "jumping"
			if math.Abs(c.VelX) > c.MinSpeedThresForJumpBoostMultiplier {
				c.VelY = c.JumpPower * c.JumpBoostMultiplier
			} else {
				c.VelY = c.JumpPower
			}
			c.JumpTimer = 0
		} else if c.VelX == 0 {
			c.CurrentState = "idle"
		}

		if !c.IsRunKeyPressed {
			c.CurrentState = "walking"
		}

	case "jumping":
		anim.SetState("jump")
		if !c.IsJumpKeyPressed && c.JumpTimer < c.JumpReleaseTimer {
			c.VelY = c.ShortJumpVelocity
			c.JumpTimer = c.JumpHoldTime // Zıplama süresini bitir
		} else if c.IsJumpKeyPressed && c.JumpTimer < c.JumpHoldTime {
			speedFactor := (math.Abs(c.VelX) / c.MaxRunSpeed) * c.SpeedJumpFactor
			c.VelY += c.JumpBoost * (1 + speedFactor)
			c.JumpTimer++
		} else if c.VelY >= 0 {
			c.CurrentState = "falling"
		}

		if c.IsLeftKeyPressed && c.VelX > 0 {
			c.VelX -= c.Deceleration
		} else if c.IsRightKeyPressed && c.VelX < 0 {
			c.VelX += c.Deceleration
		}

	case "falling":
		anim.SetState("jump")
		if c.IsOnFloor {
			if c.VelX == 0 {
				c.CurrentState = "idle"
			} else if c.IsRunKeyPressed {
				c.CurrentState = "running"
			} else {
				c.CurrentState = "walking"
			}
		}
	}

}
