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
		aabb, vel, _ := q.Get()
		tileCollider.Collide(*aabb, *(*Vec)(vel), func(hitInfos []HitTileInfo, delta Vec) {
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
func (p *Enemy) Draw() {
	q := filterEnemy.Query()
	for q.Next() {
		aabb, _, _ := q.Get()
		drawAABB(aabb)
	}
}
