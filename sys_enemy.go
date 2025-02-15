package kar

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/tilecollider"
)

type Enemy struct {
}

func (e *Enemy) Init() {
}

func (e *Enemy) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		x, y := CameraRes.ScreenToWorld(ebiten.CursorPosition())
		SpawnEnemy(x, y, 0.5, 0)
	}

	if ECWorld.Alive(CurrentPlayer) {
		playerPos, playerSize, playerVelocity, _, _, _ := MapPlayer.Get(CurrentPlayer)
		enemyQuery := FilterEnemy.Query(&ECWorld)

		for enemyQuery.Next() {

			enemyPos, enemyVel, enemyAI := enemyQuery.Get()

			switch enemyAI.Name {
			case "worm":
				// enemyPos.X += enemyVel.X
				// enemyPos.Y += enemyVel.Y

				Collider.Collide(
					enemyPos.X,
					enemyPos.Y,
					EnemyWormSize.W,
					EnemyWormSize.H,
					enemyVel.X,
					enemyVel.Y,
					func(infos []tilecollider.CollisionInfo[uint8], dx, dy float64) {
						enemyPos.X += dx
						enemyPos.Y += dy
						for _, info := range infos {
							if info.Normal == [2]int{-1, 0} {
								enemyVel.X *= -1
							}
							if info.Normal == [2]int{1, 0} {
								enemyVel.X *= -1
							}
							// TileMapRes.Set(info.TileCoords[0], info.TileCoords[1], items.Air)

						}
					},
				)

				// player-enemy collision
				collInfo := CheckCollision(playerPos, playerSize, playerVelocity, enemyPos, &EnemyWormSize)

				if collInfo.Collided {
					playerPos.X += collInfo.DeltaX
					playerPos.Y += collInfo.DeltaY
					if collInfo.Normal[1] == -1 {
						toRemove = append(toRemove, enemyQuery.Entity())
						playerVelocity.Y = -1
						// playerVelocity.X += enemyVel.X
						fmt.Println("TOP")
					}
					if collInfo.Normal[1] == 1 {
						toRemove = append(toRemove, enemyQuery.Entity())
						// playerVelocity.Y = -5
						playerVelocity.Y = 1
						fmt.Println("TOP")
					}
					if collInfo.Normal[0] == -1 {
						// playerVelocity.X = -20
						// playerHealth.Current -= 6
						fmt.Println("SAĞ")
					}
					if collInfo.Normal[0] == 1 {
						// playerVelocity.X = 20
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
		pos, _, _ := q.Get()
		x, y := CameraRes.ApplyCameraTransformToPoint(pos.X, pos.Y)
		vector.DrawFilledRect(
			Screen,
			float32(x),
			float32(y),
			float32(EnemyWormSize.W),
			float32(EnemyWormSize.H),
			color.RGBA{128, 0, 0, 10},
			false,
		)
	}
}
