package kar

type Enemy struct {
	hit *HitInfo
}

func (e *Enemy) Init() {
	e.hit = &HitInfo{}
}

func (e *Enemy) Update() {
	enemyQuery := filterEnemy.Query()
	for enemyQuery.Next() {
		aabb, vel, _ := enemyQuery.Get()
		delta := tileCollider.Collide(*aabb, *(*Vec)(vel), nil)
		aabb.Pos = aabb.Pos.Add(delta)
		if vel.X != delta.X {
			vel.X *= -1
		}
	}
}
func (e *Enemy) Draw() {
	q := filterEnemy.Query()
	for q.Next() {
		aabb, _, _ := q.Get()
		DrawAABB(aabb)
	}
}
