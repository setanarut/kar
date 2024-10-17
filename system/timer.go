package system

import (
	"fmt"
	"kar/comp"
	"kar/res"
	"kar/types"
	"time"

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

	// drop item animation frames
	comp.Index.Each(res.ECSWorld, ItemAnimationIndexUpdateCallbackFunc)

	// spawn timer update
	comp.SpawnTimer.Each(res.ECSWorld, spawnTimerComponentUpdateCallbackFunc)
}

func (s *TimersSystem) Draw(screen *ebiten.Image) {}

func spawnTimerComponentUpdateCallbackFunc(e *donburi.Entry) {
	TimerUpdate(comp.SpawnTimer.Get(e))
}

func ItemAnimationIndexUpdateCallbackFunc(e *donburi.Entry) {
	compIndex := comp.Index.Get(e)
	if compIndex.Index < ItemAnimFrameCount {
		compIndex.Index++
	} else {
		compIndex.Index = 0
	}
}
func TimerRemaining(t *types.Timer) time.Duration {
	return t.Duration - t.Elapsed
}

func TimerRemainingSecondsString(t *types.Timer) string {
	return fmt.Sprintf("%.1fs", TimerRemaining(t).Abs().Seconds())
}

func Timerreset(t *types.Timer) {
	t.Elapsed = 0
}

func TimerIsReady(t *types.Timer) bool {
	return t.Elapsed > t.Duration
}
func TimerIsStart(t *types.Timer) bool {
	return t.Elapsed == 0
}
func TimerUpdate(timer *types.Timer) {
	if timer.Elapsed < timer.Duration {
		timer.Elapsed += res.TimerTick
	}
}
