package cm

import (
	"kar/engine/vec"
	"math"
)

type Circle struct {
	*Shape
	c, transformC vec.Vec2
	radius        float64
}

func NewCircle(body *Body, radius float64, offset vec.Vec2) *Shape {
	circle := &Circle{
		c:      offset,
		radius: radius,
	}
	circle.Shape = NewShape(circle, body, CircleShapeMassInfo(0, radius, offset))
	return circle.Shape
}

func CircleShapeMassInfo(mass, radius float64, center vec.Vec2) *ShapeMassInfo {
	return &ShapeMassInfo{
		m:    mass,
		i:    MomentForCircle(1, 0, radius, vec.Vec2{}),
		cog:  center,
		area: AreaForCircle(0, radius),
	}
}

func (circle *Circle) CacheData(transform Transform) BB {
	circle.transformC = transform.Point(circle.c)
	return NewBBForCircle(circle.transformC, circle.radius)
}

func (circle *Circle) Radius() float64 {
	return circle.radius
}

func (circle *Circle) SetRadius(r float64) {
	circle.radius = r

	mass := circle.massInfo.m
	circle.massInfo = CircleShapeMassInfo(mass, circle.radius, circle.c)
	if mass > 0 {
		circle.body.AccumulateMassFromShapes()
	}
}

func (circle *Circle) TransformC() vec.Vec2 {
	return circle.transformC
}

func (circle *Circle) PointQuery(p vec.Vec2, info *PointQueryInfo) {
	delta := p.Sub(circle.transformC)
	d := delta.Length()
	r := circle.radius

	info.Shape = circle.Shape
	info.Point = circle.transformC.Add(delta.Scale(r / d))
	info.Distance = d - r

	if d > MagicEpsilon {
		info.Gradient = delta.Scale(1 / d)
	} else {
		info.Gradient = vec.Vec2{0, 1}
	}
}

func (circle *Circle) SegmentQuery(a, b vec.Vec2, radius float64, info *SegmentQueryInfo) {
	CircleSegmentQuery(circle.Shape, circle.transformC, circle.radius, a, b, radius, info)
}

func CircleSegmentQuery(shape *Shape, center vec.Vec2, r1 float64, a, b vec.Vec2, r2 float64, info *SegmentQueryInfo) {
	da := a.Sub(center)
	db := b.Sub(center)
	rsum := r1 + r2

	qa := da.Dot(da) - 2*da.Dot(db) + db.Dot(db)
	qb := da.Dot(db) - da.Dot(da)
	det := qb*qb - qa*(da.Dot(da)-rsum*rsum)

	if det >= 0 {
		t := (-qb - math.Sqrt(det)) / qa
		if 0 <= t && t <= 1 {
			n := da.Lerp(db, t).Normalize()

			info.Shape = shape
			info.Point = a.Lerp(b, t).Sub(n.Scale(r2))
			info.Normal = n
			info.Alpha = t
		}
	}
}
