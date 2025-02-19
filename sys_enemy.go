package kar

import (
	"fmt"
	"image/color"
	"kar/v"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Enemy struct {
	enemyRect AABB
}

func (e *Enemy) Init() {
	e.enemyRect = AABB{Half: EnemyWormHalfSize}
}

func (e *Enemy) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		x, y := CameraRes.ScreenToWorld(ebiten.CursorPosition())
		SpawnEnemy(x, y, 0.5, 0)
	}

	if ECWorld.Alive(CurrentPlayer) {
		playerBox, playerVelocity, _, _, _ := MapPlayer.Get(CurrentPlayer)
		enemyQuery := FilterEnemy.Query(&ECWorld)

		for enemyQuery.Next() {

			enemyPos, enemyVel, enemyAI := enemyQuery.Get()

			switch enemyAI.Name {
			case "worm":
				// enemyPos.X += enemyVel.X
				// enemyPos.Y += enemyVel.Y

				TileCollider.Collide(
					e.enemyRect,
					Vec(*enemyVel),
					func(infos []HitTileInfo, delta Vec) {
						enemyPos.X += delta.X
						enemyPos.Y += delta.Y
						for _, info := range infos {
							if info.Normal == v.Left {
								enemyVel.X *= -1
							}
							if info.Normal == v.Right {
								enemyVel.X *= -1
							}
							// TileMapRes.Set(info.TileCoords[0], info.TileCoords[1], items.Air)

						}
					},
				)

				// player-enemy collision
				hit := &HitInfo{}
				wormBox := AABB{
					Pos:  Vec(*enemyPos),
					Half: EnemyWormHalfSize,
				}

				collided := wormBox.OverlapSweep(*playerBox, Vec(*playerVelocity), hit)

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
	return nil
}
func (e *Enemy) Draw() {
	q := FilterEnemy.Query(&ECWorld)
	for q.Next() {
		pos, _, AI := q.Get()
		x, y := CameraRes.ApplyCameraTransformToPoint(pos.X, pos.Y)
		switch AI.Name {
		case "worm":
			vector.DrawFilledRect(
				Screen,
				float32(x-EnemyWormHalfSize.X),
				float32(y-EnemyWormHalfSize.Y),
				float32(EnemyWormHalfSize.X*2),
				float32(EnemyWormHalfSize.Y*2),
				color.RGBA{128, 0, 0, 10},
				false,
			)
		}
	}
}
