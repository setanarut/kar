package arc

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Controller struct {
	CurrentState string

	VelocityX float64
	VelocityY float64

	OnFloor                bool
	EarthGravity           float64
	JumpPower              float64
	MinJumpPower           float64
	MaxWalkSpeed           float64
	MaxRunSpeed            float64
	RunSpeedMultiplier     float64
	XFriction              float64
	Acceleration           float64
	calculatedJumpVelocity float64
	WalkJumpMultiplier     float64
	RunJumpMultiplier      float64

	isKeyLeftPressed      bool
	isKeyRightPressed     bool
	isKeyJumpJustPressed  bool
	isKeyJumpJustReleased bool
	isKeyRunPressed       bool
}

func NewController() *Controller {
	p := &Controller{
		CurrentState: "idle",

		EarthGravity:       0.50,
		JumpPower:          10.0,
		MinJumpPower:       3.00,
		MaxWalkSpeed:       10.0,
		MaxRunSpeed:        14.0,
		RunSpeedMultiplier: 1.50,
		WalkJumpMultiplier: 1.01,
		RunJumpMultiplier:  1.01,
		XFriction:          0.90,
		Acceleration:       0.30,
	}

	// İstenen yükseklik için zıplama hızını hesapla
	desiredHeight := 63.0
	p.calculatedJumpVelocity = p.calculateJumpPower(desiredHeight)

	return p
}

func (c *Controller) Update() {

	c.isKeyLeftPressed = ebiten.IsKeyPressed(ebiten.KeyA)
	c.isKeyRightPressed = ebiten.IsKeyPressed(ebiten.KeyD)
	c.isKeyJumpJustPressed = inpututil.IsKeyJustPressed(ebiten.KeySpace)
	c.isKeyJumpJustReleased = inpututil.IsKeyJustReleased(ebiten.KeySpace)
	c.isKeyRunPressed = ebiten.IsKeyPressed(ebiten.KeyShiftRight)

	// Durum güncellemesi
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
	}
}

func (c *Controller) calculateJumpPower(desiredHeight float64) float64 {
	jumpVelocity := math.Sqrt(2 * c.EarthGravity * desiredHeight)
	return -jumpVelocity
}

func (c *Controller) Idle() {
	// Sürtünme uygulaması - daha yumuşak bir kayma için değiştirildi
	c.VelocityX *= c.XFriction

	// Çok küçük hızları sıfırla (titreşimi önlemek için)
	if math.Abs(c.VelocityX) < 0.1 {
		c.VelocityX = 0
	}

	// Yerçekimi ve dikey hareket
	c.VelocityY += c.EarthGravity

	// Durum geçişleri
	if c.isKeyLeftPressed || c.isKeyRightPressed {
		if c.isKeyRunPressed {
			c.CurrentState = "running"
		} else {
			c.CurrentState = "walking"
		}
	}
	if c.OnFloor && c.isKeyJumpJustPressed {
		c.CurrentState = "jumping"
		c.VelocityY = c.calculatedJumpVelocity
		c.isKeyJumpJustReleased = false
	}
}

func (c *Controller) Walking() {
	// Yatay hareket
	if c.isKeyLeftPressed {
		c.VelocityX -= c.Acceleration
	} else if c.isKeyRightPressed {
		c.VelocityX += c.Acceleration
	}

	// Sürtünme uygulaması (basitleştirilmiş)
	c.VelocityX *= c.XFriction

	// Hız sınırlaması
	if c.VelocityX > c.MaxWalkSpeed {
		c.VelocityX = c.MaxWalkSpeed
	} else if c.VelocityX < -c.MaxWalkSpeed {
		c.VelocityX = -c.MaxWalkSpeed
	}

	// Yerçekimi ve dikey hareket
	c.VelocityY += c.EarthGravity

	// Durum geçişleri
	if c.isKeyRunPressed {
		c.CurrentState = "running"
	}
	if !c.isKeyLeftPressed && !c.isKeyRightPressed {
		c.CurrentState = "idle"
	}
	if c.OnFloor && c.isKeyJumpJustPressed {
		c.CurrentState = "jumping"
		c.VelocityY = c.calculatedJumpVelocity * c.WalkJumpMultiplier
		c.isKeyJumpJustReleased = false
	}
}

func (c *Controller) Running() {
	// Yatay hareket
	if c.isKeyLeftPressed {
		c.VelocityX -= c.Acceleration * c.RunSpeedMultiplier
	} else if c.isKeyRightPressed {
		c.VelocityX += c.Acceleration * c.RunSpeedMultiplier
	}

	// Sürtünme uygulaması (Walking ile aynı basitleştirilmiş mantık kullanılmalı)
	c.VelocityX *= c.XFriction

	// Hız sınırlaması
	if c.VelocityX > c.MaxRunSpeed {
		c.VelocityX = c.MaxRunSpeed
	} else if c.VelocityX < -c.MaxRunSpeed {
		c.VelocityX = -c.MaxRunSpeed
	}

	// Yerçekimi ve dikey hareket
	c.VelocityY += c.EarthGravity

	// Durum geçişleri
	if !c.isKeyRunPressed {
		c.CurrentState = "walking"
	}
	if !c.isKeyLeftPressed && !c.isKeyRightPressed {
		c.CurrentState = "idle"
	}
	if c.OnFloor && c.isKeyJumpJustPressed {
		c.CurrentState = "jumping"
		c.VelocityY = c.calculatedJumpVelocity * c.RunJumpMultiplier
		c.isKeyJumpJustReleased = false
	}
}

func (c *Controller) Jumping() {
	if c.isKeyJumpJustPressed {
		c.VelocityY = c.calculatedJumpVelocity
		c.isKeyJumpJustReleased = false
	}

	// Yatay hareket kontrolü
	acceleration := c.Acceleration
	if c.isKeyRunPressed {
		acceleration *= c.RunSpeedMultiplier
	}

	if c.isKeyLeftPressed {
		c.VelocityX -= acceleration
	} else if c.isKeyRightPressed {
		c.VelocityX += acceleration
	}

	// Sürtünme uygulaması
	c.VelocityX *= c.XFriction

	// Hız sınırlaması
	maxSpeed := c.MaxWalkSpeed
	if c.isKeyRunPressed {
		maxSpeed = c.MaxRunSpeed
	}
	if c.VelocityX > maxSpeed {
		c.VelocityX = maxSpeed
	} else if c.VelocityX < -maxSpeed {
		c.VelocityX = -maxSpeed
	}

	// Yerçekimi ve dikey hareket
	c.VelocityY += c.EarthGravity

	// Zıplama yüksekliği kontrolü
	if c.isKeyJumpJustReleased {
		if c.VelocityY < -c.MinJumpPower {
			c.VelocityY = -c.MinJumpPower
		}
		c.isKeyJumpJustReleased = true
	}

	// Durum geçişi
	if c.VelocityY > 0 {
		c.CurrentState = "falling"
	}
}

func (c *Controller) Falling() {
	// Yatay hareket kontrolü
	acceleration := c.Acceleration
	if c.isKeyRunPressed {
		acceleration *= c.RunSpeedMultiplier
	}

	if c.isKeyLeftPressed {
		c.VelocityX -= acceleration
	} else if c.isKeyRightPressed {
		c.VelocityX += acceleration
	}

	// Sürtünme uygulaması
	c.VelocityX *= c.XFriction

	// Hız sınırlaması
	maxSpeed := c.MaxWalkSpeed
	if c.isKeyRunPressed {
		maxSpeed = c.MaxRunSpeed
	}
	if c.VelocityX > maxSpeed {
		c.VelocityX = maxSpeed
	} else if c.VelocityX < -maxSpeed {
		c.VelocityX = -maxSpeed
	}

	// Yerçekimi ve dikey hareket
	c.VelocityY += c.EarthGravity

	// Durum geçişi
	if c.OnFloor {
		if c.isKeyLeftPressed || c.isKeyRightPressed {
			if c.isKeyRunPressed {
				c.CurrentState = "running"
			} else {
				c.CurrentState = "walking"
			}
		} else {
			c.CurrentState = "idle"
		}
	}
}
