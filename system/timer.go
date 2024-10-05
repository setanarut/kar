package system

import (
	"kar/comp"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type TimersSystem struct {
}

func NewTimersSystem() *TimersSystem {
	return &TimersSystem{}
}

func (s *TimersSystem) Init() {}

func (s *TimersSystem) Update() {
	comp.Index.Each(res.ECSWorld, func(e *donburi.Entry) {
		compIndex := comp.Index.Get(e)
		if compIndex.Index < ItemAnimFrameCount {
			compIndex.Index++
		} else {
			compIndex.Index = 0
		}
	})
	// comp.Timer.Each(res.World, timerComponentUpdateFunc)
}

func (s *TimersSystem) Draw(screen *ebiten.Image) {}

/* func TimerRemaining(t *types.DataTimer) time.Duration {
	return t.TimerDuration - t.Elapsed
}

func TimerRemainingSecondsString(t *types.DataTimer) string {
	return fmt.Sprintf("%.1fs", TimerRemaining(t).Abs().Seconds())
}

func timerComponentUpdateFunc(e *donburi.Entry) {
	timer := comp.Timer.Get(e)
	if timer.Elapsed < timer.TimerDuration {
		timer.Elapsed += res.TimerTick
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
*/
