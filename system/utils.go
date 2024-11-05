package system

import (
	"kar"
	"kar/arc"
	"kar/items"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/cm"
)

func getSprite(id uint16) *ebiten.Image {
	im, ok := res.Images[id]
	if ok {
		return im
	} else {
		if len(res.Frames[id]) > 0 {
			return res.Frames[id][0]
		} else {
			return res.Images[items.Air]
		}
	}
}

func removeBodyPostStep(space *cm.Space, body, data interface{}) {
	space.RemoveBodyWithShapes(body.(*cm.Body))
}

func TimerIsReady(t *arc.Timer) bool {
	return t.Elapsed > t.Duration
}

func TimerUpdate(timer *arc.Timer) {
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
