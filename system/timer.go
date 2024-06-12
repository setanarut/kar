package system

import (
	"fmt"
	"kar/comp"
	"kar/res"
	"kar/types"
	"time"

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
		TimerUpdate(comp.AttackTimer.Get(e))
	})

}

func (s *TimersSystem) Draw() {}

func TimerRemaining(t *types.DataTimer) time.Duration {
	return t.TimerDuration - t.Elapsed
}

func TimerRemainingSecondsString(t *types.DataTimer) string {
	return fmt.Sprintf("%.1fs", TimerRemaining(t).Abs().Seconds())
}
func TimerUpdate(t *types.DataTimer) {
	if t.Elapsed < t.TimerDuration {
		t.Elapsed += res.TimerTick
	}
}
func TimerReset(t *types.DataTimer) {
	t.Elapsed = 0
}

func TimerIsReady(t *types.DataTimer) bool {
	return t.Elapsed > t.TimerDuration
}
func TimerIsStart(t *types.DataTimer) bool {
	return t.Elapsed == 0
}
