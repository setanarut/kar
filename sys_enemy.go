package kar

import (
	"image/color"
	"kar/v"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Enemy struct {
	enemyRect    AABB
	enemyHitInfo *HitInfo
}

func (e *Enemy) Init() {
	e.enemyHitInfo = &HitInfo{}
	e.enemyRect = AABB{Half: enemyWormHalfSize}
}

func (e *Enemy) Update() {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		x, y := cameraRes.ScreenToWorld(ebiten.CursorPosition())
		SpawnEnemy(Vec{x, y}, Vec{0.5, 0})
	}

	playerBox := mapAABB.GetUnchecked(currentPlayer)
	playerVelocity := (*Vec)(mapVel.GetUnchecked(currentPlayer))

	if world.Alive(currentPlayer) {
		enemyQuery := filterEnemy.Query()
		for enemyQuery.Next() {
			epos, evel, enemyAI := enemyQuery.Get()
			enemyPos := (*Vec)(epos)
			enemyVel := (*Vec)(evel)
			switch enemyAI.Name {
			case "worm":
				*enemyPos = enemyPos.Add(*enemyVel)
				e.enemyRect.Pos = *enemyPos

				// Enemy tilemap collision
				delta := tileCollider.Collide(e.enemyRect, *enemyVel, nil)
				if enemyVel.X != delta.X {
					enemyVel.X *= -1
				}

				// Player enemy collision
				if e.enemyRect.OverlapSweep(playerBox, *playerVelocity, e.enemyHitInfo) {

					*playerVelocity = playerVelocity.Add(e.enemyHitInfo.Delta)
					playerBox.Pos = playerBox.Pos.Add(*playerVelocity)

					if e.enemyHitInfo.Normal == v.Up {
						toRemove = append(toRemove, enemyQuery.Entity())
						*playerVelocity = e.enemyHitInfo.Normal.Scale(3)
					}
					if e.enemyHitInfo.Normal == v.Down {
						*playerVelocity = e.enemyHitInfo.Normal.Scale(3)
					}

					// Horizontal collision
					if e.enemyHitInfo.Normal == v.Right || e.enemyHitInfo.Normal == v.Left {
						enemyVel.X *= -1
						*playerVelocity = e.enemyHitInfo.Normal.Scale(3)
						// TODO oyuncu çarpınca çarpışma devre dışı kalsın
						// oyuncuya yanıp sönme componenti ekle
					}

				}

			case "other":
				// other
			}
		}
	}
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
