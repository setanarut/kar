package system

import (
	"kar/comp"
	"kar/res"

	"github.com/yohamta/donburi"
)

type TimersSystem struct {
}

func NewTimersSystem() *TimersSystem {
	return &TimersSystem{}
}

func (s *TimersSystem) Init() {}

func (s *TimersSystem) Update() {

	comp.AttackTimer.Each(res.World, func(e *donburi.Entry) {
		t := comp.AttackTimer.Get(e)
		if t.Elapsed < t.TimerDuration {
			t.Elapsed += res.TimerTick

		}
	})

}

func (s *TimersSystem) Draw() {}

// func timerRemaining(t *comp.DataTimer) time.Duration {
// 	return t.TimerDuration - t.Elapsed
// }

// func timerRemainingSecondsString(t *comp.DataTimer) string {
// 	return fmt.Sprintf("%.1fs", timerRemaining(t).Abs().Seconds())
// }

func timerReset(t *comp.DataTimer) {
	t.Elapsed = 0
}

func timerIsReady(t *comp.DataTimer) bool {
	return t.Elapsed > t.TimerDuration
}
func timerIsStart(t *comp.DataTimer) bool {
	return t.Elapsed == 0
}
