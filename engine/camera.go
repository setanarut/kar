package engine

import (
	"fmt"
	"kar/engine/mathutil"
	"kar/engine/vec"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ojrac/opensimplex-go"
)

// Camera object
type Camera struct {
	// Use the `Camera.LookAt()` function to align the center of the camera to the target.
	W, H, Rotation, ZoomFactor        float64
	tempTarget, centerOffset, topLeft vec.Vec2
	DrawOptions                       *ebiten.DrawImageOptions
	Lerp, TraumaEnabled               bool

	// camera shake options
	traumaOffset, ShakeSize                              vec.Vec2
	Trauma, TimeScale, MaxShakeAngle, Decay, delta, tick float64
	noise                                                opensimplex.Noise
}

// NewCamera returns new Camera
func NewCamera(lookAt vec.Vec2, w, h int) *Camera {
	c := &Camera{
		W:             float64(w),
		H:             float64(h),
		Rotation:      0,
		ZoomFactor:    0,
		DrawOptions:   &ebiten.DrawImageOptions{},
		Lerp:          false,
		TraumaEnabled: false,

		// Shake options
		Trauma:        0,
		ShakeSize:     vec.Vec2{150, 150},
		MaxShakeAngle: 30,
		TimeScale:     16,
		Decay:         0.333,

		// private
		traumaOffset: vec.Vec2{},
		topLeft:      vec.Vec2{},
		centerOffset: vec.Vec2{-(float64(w) * 0.5), -(float64(h) * 0.5)},
		tempTarget:   vec.Vec2{},
		noise:        opensimplex.New(1),
		delta:        1.0 / 60.0,
		tick:         0,
	}
	c.LookAt(lookAt)
	c.tempTarget = lookAt
	return c
}

// LookAt aligns the midpoint of the camera viewport to the target.
func (cam *Camera) LookAt(target vec.Vec2) {

	if cam.Lerp {
		cam.tempTarget = cam.tempTarget.Lerp(target, 0.1)
		cam.topLeft = cam.tempTarget
	} else {
		cam.topLeft = target
	}

	if cam.TraumaEnabled {
		if cam.Trauma > 0 {
			var shake = math.Pow(cam.Trauma, 2)
			cam.traumaOffset.X = cam.noise.Eval3(cam.tick*cam.TimeScale, 0, 0) * cam.ShakeSize.X * shake
			cam.traumaOffset.Y = cam.noise.Eval3(0, cam.tick*cam.TimeScale, 0) * cam.ShakeSize.Y * shake
			cam.Rotation = cam.noise.Eval3(0, 0, cam.tick*cam.TimeScale) * cam.MaxShakeAngle * shake
			cam.Trauma = mathutil.Clamp(cam.Trauma-(cam.delta*cam.Decay), 0, 1)
		}
		// offset
		cam.topLeft = cam.topLeft.Add(cam.traumaOffset)
	}

	cam.topLeft = cam.topLeft.Add(cam.centerOffset)
	cam.tick += cam.delta
	if cam.tick > 60000 {
		cam.tick = 0
	}
}
func (cam *Camera) AddTrauma(trauma_in float64) {
	cam.Trauma = mathutil.Clamp(cam.Trauma+trauma_in, 0, 1)
}

// SetSize returns center point of the camera
func (cam *Camera) SetSize(w, h float64) {
	cam.W, cam.H = w, h
	cam.centerOffset = vec.Vec2{-(w * 0.5), -(h * 0.5)}
}

// Center returns center point of the camera
func (cam *Camera) Center() vec.Vec2 {
	return cam.topLeft.Sub(cam.centerOffset)
}

// Reset resets all camera values to zero
func (cam *Camera) Reset() {
	cam.topLeft.X, cam.topLeft.Y = 0.0, 0.0
	cam.Rotation, cam.ZoomFactor = 0.0, 0
}

// String returns camera values as string
func (cam *Camera) String() string {
	return fmt.Sprintf(
		"CamX: %.1f\nCamY: %.1f\nCam Rotation: %.1f\nZoom factor: %.2f",
		cam.topLeft.X, cam.topLeft.Y, cam.Rotation, cam.ZoomFactor,
	)
}

// ScreenToWorld converts screen-space coordinates to world-space
func (cam *Camera) ScreenToWorld(screenX, screenY int) vec.Vec2 {
	g := ebiten.GeoM{}
	cam.ApplyCameraTransform(&g)
	if g.IsInvertible() {
		g.Invert()
		worldX, worldY := g.Apply(float64(screenX), float64(screenY))
		return vec.Vec2{worldX, worldY}
	} else {
		// When scaling it can happened that matrix is not invertable
		return vec.Vec2{math.NaN(), math.NaN()}
	}
}

// ApplyCameraTransform applies geometric transformation to given geoM
func (cam *Camera) ApplyCameraTransform(geoM *ebiten.GeoM) {
	geoM.Translate(-cam.topLeft.X, -cam.topLeft.Y)                                               // camera movement
	geoM.Translate(cam.centerOffset.X, cam.centerOffset.Y)                                       // rotate and scale from center.
	geoM.Rotate(cam.Rotation * 2 * math.Pi / 360)                                                // rotate
	geoM.Scale(math.Pow(1.01, float64(cam.ZoomFactor)), math.Pow(1.01, float64(cam.ZoomFactor))) // apply zoom factor
	geoM.Translate(math.Abs(cam.centerOffset.X), math.Abs(cam.centerOffset.Y))                   // restore center translation
}

// Draw applies the Camera's geometric transformation then draws the object on the screen with drawing options.
func (cam *Camera) Draw(worldObject *ebiten.Image, worldObjectOps *ebiten.DrawImageOptions, screen *ebiten.Image) {
	cam.DrawOptions = worldObjectOps
	cam.ApplyCameraTransform(&cam.DrawOptions.GeoM)
	screen.DrawImage(worldObject, cam.DrawOptions)
	cam.DrawOptions.GeoM.Reset()
}
