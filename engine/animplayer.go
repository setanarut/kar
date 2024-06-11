package engine

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// AnimationPlayer plays and manages animations.
type AnimationPlayer struct {
	SpriteSheet, CurrentFrame *ebiten.Image
	Animations                map[string]*Animation
	Paused                    bool
	CurrentFrameIndex         int
	CurrentState              string
	Tick                      float64
}

// NewAnimationPlayer returns new AnimationPlayer with spriteSheet
func NewAnimationPlayer(spriteSheet *ebiten.Image) *AnimationPlayer {
	return &AnimationPlayer{
		SpriteSheet:       spriteSheet,
		Paused:            false,
		Animations:        make(map[string]*Animation),
		CurrentFrameIndex: 0,
	}

}

// AddStateAnimation adds new Animation to this AnimationPlayer and returns the added animation.
//
// x and y are top-left coordinates of the first frame's rectangle.
//
// w and h are the width and height of the first frame's rectangle.
//
// The sub-rectangles repeat to the right by the frameCount amount.
func (ap *AnimationPlayer) AddStateAnimation(stateName string, x, y, w, h, frameCount int, pingpong bool, vertical bool) *Animation {

	subImages := []*ebiten.Image{}
	frameRect := image.Rect(x, y, x+w, y+h)
	for i := 0; i < frameCount; i++ {
		subImages = append(subImages, ap.SpriteSheet.SubImage(frameRect).(*ebiten.Image))
		if vertical {
			frameRect.Min.Y += h
			frameRect.Max.Y += h
		} else {
			frameRect.Min.X += w
			frameRect.Max.X += w
		}
	}

	if pingpong {
		for i := frameCount - 2; i > 0; i-- {
			subImages = append(subImages, subImages[i])
		}
	}

	anim := &Animation{
		FPS:    15.0,
		Frames: subImages,
		Name:   stateName,
	}

	ap.CurrentState = stateName
	ap.Animations[stateName] = anim
	ap.CurrentFrame = ap.Animations[ap.CurrentState].Frames[ap.CurrentFrameIndex]

	return anim
}

// SetFPS overwrites the FPS of all animations.
func (ap *AnimationPlayer) SetFPS(fps float64) {
	for _, anim := range ap.Animations {
		anim.FPS = fps
	}
}

// AddAnimation adds the given animation to this player.
// Adds the name of the animation as a map key.
func (ap *AnimationPlayer) AddAnimation(a *Animation) {
	ap.Animations[a.Name] = a
}

// State returns current active animation state
func (ap *AnimationPlayer) State() string {
	return ap.CurrentState
}

// CurrentStateFPS returns FPS of the current animation state
func (ap *AnimationPlayer) CurrentStateFPS() float64 {
	return ap.Animations[ap.State()].FPS
}

// SetStateReset sets the current animation state. Each time the state changes, the animation resets to the first frame.
func (ap *AnimationPlayer) SetStateReset(state string) {
	if ap.CurrentState != state {
		ap.CurrentState = state
		ap.Tick = 0
		ap.CurrentFrameIndex = 0
	}
}

// SetState sets the current animation state.
func (ap *AnimationPlayer) SetState(state string) {
	if ap.CurrentState != state {
		ap.CurrentState = state
	}
}

// PauseAtFrame pauses the current animation at the frame. If index is out of range it does nothing.
func (ap *AnimationPlayer) PauseAtFrame(frameIndex int) {
	if frameIndex < len(ap.Animations[ap.State()].Frames) && frameIndex >= 0 {
		ap.Paused = true
		ap.CurrentFrameIndex = frameIndex
	}
}

// Update updates AnimationPlayer
func (ap *AnimationPlayer) Update() {
	if !ap.Paused {
		ap.Tick += ap.Animations[ap.CurrentState].FPS / 60.0
		ap.CurrentFrameIndex = int(math.Floor(ap.Tick))
		if ap.CurrentFrameIndex >= len(ap.Animations[ap.CurrentState].Frames) {
			ap.Tick = 0
			ap.CurrentFrameIndex = 0
		}
	}

	// update image
	ap.CurrentFrame = ap.Animations[ap.CurrentState].Frames[ap.CurrentFrameIndex]
}

// Animation for AnimationPlayer
type Animation struct {
	Frames []*ebiten.Image
	FPS    float64
	Name   string
}

// func oscillate(input, min, max int) int {
// 	rang := max - min
// 	return min + int(math.Abs(float64(((input+rang)%(rang*2))-rang)))
// }
