package kar

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Enemy struct {
	enemyRect *AABB
	hit       *HitInfo
}

func (e *Enemy) Init() {
	e.hit = &HitInfo{}
	e.enemyRect = &AABB{Half: enemyWormHalfSize}
}

func (e *Enemy) Update() {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		x, y := cameraRes.ScreenToWorld(ebiten.CursorPosition())
		SpawnEnemy(Vec{x, y}, Vec{0.5, 0})
	}

	// if world.Alive(currentPlayer) {

	// 	enemyQuery := filterEnemy.Query()

	// 	for enemyQuery.Next() {
	// 		epos, evel, _ := enemyQuery.Get()
	// 		enemyPos := (*Vec)(epos)
	// 		enemyVel := (*Vec)(evel)
	// 		e.enemyRect.Pos = *enemyPos

	// 		// switch enemyAI.Name {
	// 		// Player enemy collision
	// 		pBox := mapAABB.GetUnchecked(currentPlayer)
	// 		pVel := (*Vec)(mapVel.GetUnchecked(currentPlayer))

	// 		e.hit.Reset()
	// 		if OverlapSweep2(e.enemyRect, pBox, *enemyVel, *pVel, e.hit) {
	// 			*pVel = pVel.Add(e.hit.Delta)
	// 			// apply enemy collision delta to player
	// 			// pBox.Pos = pBox.Pos.Add(*pVel)

	// 			// if e.hit.Normal == v.Up {
	// 			// 	toRemove = append(toRemove, enemyQuery.Entity())
	// 			// 	*pVel = e.hit.Normal.Scale(2)
	// 			// }
	// 			// if e.hit.Normal == v.Down {
	// 			// 	pVel.Y = e.hit.Normal.Y * math.Abs(pVel.Y)
	// 			// 	mapHealth.GetUnchecked(currentPlayer).Current -= 10
	// 			// }
	// 			// if e.hit.Normal == v.Right || e.hit.Normal == v.Left {
	// 			// 	mapHealth.GetUnchecked(currentPlayer).Current -= 10
	// 			// 	pVel.X = e.hit.Normal.X * math.Abs(pVel.X)
	// 			// }

	// 		}

	// 		// delta := tileCollider.Collide(*e.enemyRect, *enemyVel, nil)
	// 		// *enemyPos = enemyPos.Add(delta)

	// 		// if enemyVel.X != delta.X {
	// 		// 	enemyVel.X *= -1
	// 		// }
	// 		// }
	// 	}
	// }

}
func (e *Enemy) Draw() {

	q := filterEnemy.Query()
	for q.Next() {
		pos, _, AI := q.Get()
		x, y := cameraRes.ApplyCameraTransformToPoint(pos.X, pos.Y)
		switch AI.Name {
		case "worm":
			vector.DrawFilledRect(
				Screen,
				float32(x-enemyWormHalfSize.X),
				float32(y-enemyWormHalfSize.Y),
				float32(enemyWormHalfSize.X*2),
				float32(enemyWormHalfSize.Y*2),
				color.RGBA{128, 0, 0, 10},
				false,
			)
		}
	}
}
