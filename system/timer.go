package system

import (
	"fmt"
	"kar/comp"
	"kar/constants"
	"kar/res"
	"time"

	"github.com/yohamta/donburi"
)

type Timers struct {
}

func NewTimersSystem() *Timers {
	return &Timers{}
}

func (s *Timers) Init() {}

func (s *Timers) Update() {

	comp.AttackTimer.Each(res.World, func(e *donburi.Entry) {
		t := comp.AttackTimer.Get(e)
		if t.Elapsed < t.TimerDuration {
			t.Elapsed += constants.TimerTick

		}
	})

}

func (s *Timers) Draw() {
}

func Remaining(t *comp.DataTimer) time.Duration {
	return t.TimerDuration - t.Elapsed
}

func RemainingSecondsString(t *comp.DataTimer) string {
	return fmt.Sprintf("%.1fs", Remaining(t).Abs().Seconds())
}

func ResetTimer(t *comp.DataTimer) {
	t.Elapsed = 0
}

func IsTimerReady(t *comp.DataTimer) bool {
	return t.Elapsed > t.TimerDuration
}
func IsTimerStart(t *comp.DataTimer) bool {
	return t.Elapsed == 0
}
