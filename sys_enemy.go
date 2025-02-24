package kar

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Enemy struct {
	enemyRect AABB
}

func (e *Enemy) Init() {
	e.enemyRect = AABB{Half: enemyWormHalfSize}
}

func (e *Enemy) Update() {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		x, y := cameraRes.ScreenToWorld(ebiten.CursorPosition())
		SpawnEnemy(Vec{x, y}, Vec{0.5, 0})
	}

	if world.Alive(currentPlayer) {
		playerBox := mapAABB.GetUnchecked(currentPlayer)
		playerVelocity := (*Vec)(mapVel.GetUnchecked(currentPlayer))

		enemyQuery := filterEnemy.Query()
		for enemyQuery.Next() {
			p, v, enemyAI := enemyQuery.Get()
			enemyPos := (*Vec)(p)
			enemyVel := (*Vec)(v)
			switch enemyAI.Name {
			case "worm":
				*enemyPos = enemyPos.Add(*enemyVel)
				e.enemyRect.Pos = *enemyPos
				TileCollider.Collide(
					e.enemyRect,
					*enemyVel,
					func(infos []HitTileInfo, delta Vec) {
						for _, info := range infos {
							if info.Normal.X == -1 {
								enemyVel.X *= -1
							}
							if info.Normal.X == 1 {
								enemyVel.X *= -1
							}
							// TileMapRes.Set(info.TileCoords[0], info.TileCoords[1], items.Air)

						}
					},
				)

				// player-enemy collision
				hit := &HitInfo{}

				collided := e.enemyRect.OverlapSweep(playerBox, *playerVelocity, hit)

				if collided {
					playerVelocity.X += hit.Delta.X
					playerVelocity.Y += hit.Delta.Y
					playerBox.Pos.X += playerVelocity.X
					playerBox.Pos.Y += playerVelocity.Y

					if hit.Normal.Y < 0 {
						toRemove = append(toRemove, enemyQuery.Entity())
						playerVelocity.Y = -2
						// playerVelocity.X += enemyVel.X
						fmt.Println("TOP")
					}
					if hit.Normal.Y > 0 {
						playerVelocity.Y = 2
						fmt.Println("ALT")
					}
					if hit.Normal.X < 0 {
						playerVelocity.X = -2
						// playerHealth.Current -= 6
						// TODO oyuncu çarpınca çarpışma devre dışı kalsın
						// oyuncuya yanıp sönme componenti ekle
						fmt.Println("SAĞ")
					}
					if hit.Normal.X > 0 {
						playerVelocity.X = 2
						fmt.Println("SOL")
						// playerHealth.Current -= 6
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
