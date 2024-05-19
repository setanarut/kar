package engine

import (
	"fmt"
	"time"
)

type Timer struct {
	Target  time.Duration
	Elapsed time.Duration
	tick    time.Duration
}

func NewTimer(duration time.Duration) Timer {

	return Timer{
		Target:  duration,
		Elapsed: 0,
		tick:    time.Second / 60,
	}

}

func (t *Timer) Update() {
	if t.Elapsed < t.Target {
		t.Elapsed += t.tick

	}
}

func (t *Timer) IsReady() bool {
	return t.Elapsed > t.Target
}
func (t *Timer) IsStart() bool {
	return t.Elapsed == 0
}

func (t *Timer) Reset() {
	t.Elapsed = 0
}

func (t *Timer) Remaining() time.Duration {
	return t.Target - t.Elapsed
}
func (t *Timer) RemainingSecondsString() string {
	return fmt.Sprintf("%.1fs", t.Remaining().Abs().Seconds())
}
