package kar

import (
	"kar/items"
	"math"
)

type Projectile struct {
	snowBallBox    AABB
	bounceVelocity float64
}

func (p *Projectile) Init() {
	p.snowBallBox = AABB{Half: Vec{4, 4}}
	p.bounceVelocity = -math.Sqrt(2 * SnowballGravity * SnowballBounceHeight)
}
func (p *Projectile) Update() {
	// projectile physics
	q := filterProjectile.Query()
	for q.Next() {
		itemID, projectilePos, projectileVel := q.Get()
		// snowball physics
		if itemID.ID == items.Snowball {
			projectileVel.Y += SnowballGravity
			projectileVel.Y = min(projectileVel.Y, SnowballMaxFallVelocity)
			p.snowBallBox.Pos.X = projectilePos.X
			p.snowBallBox.Pos.Y = projectilePos.Y
			tileCollider.Collide(p.snowBallBox, *(*Vec)(projectileVel), func(ci []HitTileInfo, delta Vec) {
				projectilePos.X += delta.X
				projectilePos.Y += delta.Y
				isHorizontalCollision := false
				for _, cinfo := range ci {
					if cinfo.Normal.Y == -1 {

						projectileVel.Y = p.bounceVelocity
					}
					if cinfo.Normal.X == -1 && projectileVel.X > 0 && projectileVel.Y > 0 {
						isHorizontalCollision = true
					}
					if cinfo.Normal.X == 1 && projectileVel.X < 0 && projectileVel.Y > 0 {
						isHorizontalCollision = true
					}
				}
				if isHorizontalCollision {
					if world.Alive(q.Entity()) {
						toRemove = append(toRemove, q.Entity())
					}
				}
			},
			)
		}
	}
}
func (p *Projectile) Draw() {}
