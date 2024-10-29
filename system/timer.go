package system

import (
	"kar"
	"kar/comp"
	"kar/types"

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

func timerIsReady(t *types.Timer) bool {
	return t.Elapsed > t.Duration
}

func timerUpdate(timer *types.Timer) {
	if timer.Elapsed < timer.Duration {
		timer.Elapsed += kar.TimerTick
	}
}

// func timerRemaining(t *types.Timer) time.Duration {
// 	return t.Duration - t.Elapsed
// }

// func timerRemainingSecondsString(t *types.Timer) string {
// 	return fmt.Sprintf("%.1fs", timerRemaining(t).Abs().Seconds())
// }

// func timerReset(t *types.Timer) {
// 	t.Elapsed = 0
// }

// func timerIsStart(t *types.Timer) bool {
// 	return t.Elapsed == 0
// }
