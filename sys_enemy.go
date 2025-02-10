package kar

import (
	"kar/items"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/tilecollider"
)

type Enemy struct {
	normal [2]int
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
			pos, vel, ai := enemyQuery.Get()
			switch ai.Name {
			case "worm":
				pos.X += vel.X
				pos.Y += vel.Y
				e.normal = [2]int{0, 0}

				Collider.Collide(
					pos.X,
					pos.Y,
					EnemyWormSize.W,
					EnemyWormSize.H,
					vel.X,
					vel.Y,
					func(infos []tilecollider.CollisionInfo[uint8], dx, dy float64) {
						// Apply tilemap collision response
						// pos.X += dx
						// pos.Y += dy
						for _, info := range infos {
							TileMapRes.Set(info.TileCoords[0], info.TileCoords[1], items.Air)
						}
					},
				)
				// player-enemy collision
				collInfo := CheckCollision(playerPos, playerSize, playerVelocity, pos, &EnemyWormSize)
				// playerVelocity.X += collInfo.DeltaX
				playerVelocity.Y += collInfo.DeltaY

				if collInfo.Collided {
					if collInfo.Normal[1] == -1 {
						playerVelocity.Y = -5
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
		ColorMDIO.GeoM.Reset()
		ColorMDIO.GeoM.Translate(pos.X, pos.Y)
		CameraRes.DrawWithColorM(res.Icon8[items.Sand], ColorM, ColorMDIO, Screen)
	}
}
