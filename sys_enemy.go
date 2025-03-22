package kar

import "kar/items"

type Enemy struct {
	hit *HitInfo
}

func (p *Enemy) Init() {
	p.hit = &HitInfo{}
}

func (p *Enemy) Update() {
	q := filterEnemy.Query()
	for q.Next() {
		aabb, vel, mobileID := q.Get()
		switch *mobileID {
		case WormID:
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
		aabb, _, mobileID := q.Get()
		switch *mobileID {
		case WormID:
			// TODO draw worm enemy
			drawAABB(aabb)
		}

	}
}
