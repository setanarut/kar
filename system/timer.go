package system

import (
	"kar/comp"

	"github.com/yohamta/donburi"
)

type Timers struct {
}

func NewTimersSystem() *Timers {
	return &Timers{}
}

func (s *Timers) Init() {}

func (s *Timers) Update() {

	// drop item animation frames
	comp.Index.Each(ecsWorld, itemAnimationIndexUpdateFunc)

	// spawn timer update
	comp.CollisionTimer.Each(ecsWorld, spawnTimerUpdateFunc)
}

func (s *Timers) Draw() {}

func spawnTimerUpdateFunc(e *donburi.Entry) {
	timerUpdate(comp.CollisionTimer.Get(e))
}

func itemAnimationIndexUpdateFunc(e *donburi.Entry) {
	compIndex := comp.Index.Get(e)
	if compIndex.Index < itemAnimFrameCount {
		compIndex.Index++
	} else {
		compIndex.Index = 0
	}
}
