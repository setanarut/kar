package system

import (
	"fmt"
	"kar/comp"
	"kar/constants"
	"kar/model"
	"kar/res"
	"time"
)

type Timers struct {
}

func NewTimersSystem() *Timers {
	return &Timers{}
}

func (s *Timers) Init() {}

func (s *Timers) Update() {

	if p, ok := comp.PlayerTag.First(res.World); ok {
		c := comp.Char.Get(p)
		UpdateTimer(c.ShootCooldown)
	}
}

func (s *Timers) Draw() {
}

func Remaining(t *model.TimerData) time.Duration {
	return t.Target - t.Elapsed
}

func RemainingSecondsString(t *model.TimerData) string {
	return fmt.Sprintf("%.1fs", Remaining(t).Abs().Seconds())
}

func ResetTimer(t *model.TimerData) {
	t.Elapsed = 0
}

func IsTimerReady(t *model.TimerData) bool {
	return t.Elapsed > t.Target
}
func IsTimerStart(t *model.TimerData) bool {
	return t.Elapsed == 0
}
func UpdateTimer(t *model.TimerData) {
	if t.Elapsed < t.Target {
		t.Elapsed += constants.TimerTick

	}
}
