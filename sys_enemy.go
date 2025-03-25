package kar

import (
	"kar/items"
	"kar/res"
	"math"
)

type Enemy struct {
	hit *HitInfo
}

func (p *Enemy) Init() {
	p.hit = &HitInfo{}
}

func (p *Enemy) Update() {
	q := filterEnemy.Query()
	for q.Next() {
		aabb, vel, mobileID, tick := q.Get()
		*tick += 0.09
		if *tick >= 2 {
			*tick = 0
		}
		switch *mobileID {
		case CrabID:
			enemyVel := *(*Vec)(vel)
			tileCollider.Collide(*aabb, enemyVel, func(hitInfos []HitTileInfo, delta Vec) {
				aabb.Pos = aabb.Pos.Add(delta)
				for _, hit := range hitInfos {
					spawnData := spawnEffectData{
						Pos: tileMapRes.TileToWorld(hit.TileCoords),
						Id:  tileMapRes.GetIDUnchecked(hit.TileCoords),
					}
					toSpawnEffect = append(toSpawnEffect, spawnData)
					tileMapRes.SetUnchecked(hit.TileCoords, items.Air)
				}
			})
		}
	}
}
func (p *Enemy) Draw() {
	q := filterEnemy.Query()
	for q.Next() {
		aabb, _, mobileID, idx := q.Get()
		switch *mobileID {
		case CrabID:
			// TODO draw worm enemy
			tl := aabb.TopLeft()
			colorMDIO.GeoM.Reset()
			colorMDIO.GeoM.Translate(tl.X, tl.Y)
			cameraRes.DrawWithColorM(res.Crab[int(math.Floor(float64(*idx)))], colorM, colorMDIO, Screen)
			// drawAABB(aabb)
		}

	}
}
