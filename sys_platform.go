package kar

import (
	"kar/items"
	"kar/res"
)

type Platform struct {
	hit *HitInfo
}

func (p *Platform) Init() {
	p.hit = &HitInfo{}
}

func (p *Platform) Update() {
	q := filterPlatform.Query()
	for q.Next() {
		aabb, vel := q.Get()
		delta := tileCollider.Collide(*aabb, *(*Vec)(vel), nil)
		aabb.Pos = aabb.Pos.Add(delta)
		if vel.X != delta.X {
			vel.X *= -1
		}
	}
}
func (p *Platform) Draw() {
	q := filterPlatform.Query()
	for q.Next() {
		aabb, _ := q.Get()
		topLeftPos := aabb.TopLeft()
		colorMDIO.GeoM.Reset()
		colorMDIO.GeoM.Translate(topLeftPos.X, topLeftPos.Y)
		cameraRes.DrawWithColorM(res.BlockCrackFrames[items.Stone][0], colorM, colorMDIO, Screen)
		// DrawAABB(aabb)
	}
}
