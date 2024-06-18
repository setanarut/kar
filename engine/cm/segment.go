package cm

import "kar/engine/vec"

// Segment is a segment Shape
type Segment struct {
	*Shape
	a, b, n                            vec.Vec2
	transformA, transformB, transformN vec.Vec2
	radius                             float64
	aTangent, bTangent                 vec.Vec2
}

func (seg *Segment) CacheData(transform Transform) BB {
	seg.transformA = transform.Point(seg.a)
	seg.transformB = transform.Point(seg.b)
	seg.transformN = transform.Vect(seg.n)

	var l, r, b, t float64

	if seg.transformA.X < seg.transformB.X {
		l = seg.transformA.X
		r = seg.transformB.X
	} else {
		l = seg.transformB.X
		r = seg.transformA.X
	}

	if seg.transformA.Y < seg.transformB.Y {
		b = seg.transformA.Y
		t = seg.transformB.Y
	} else {
		b = seg.transformB.Y
		t = seg.transformA.Y
	}

	rad := seg.radius
	return BB{l - rad, b - rad, r + rad, t + rad}
}

func (seg *Segment) SetRadius(r float64) {
	seg.radius = r

	mass := seg.massInfo.m
	seg.massInfo = NewSegmentMassInfo(seg.massInfo.m, seg.a, seg.b, seg.radius)
	if mass > 0 {
		seg.body.AccumulateMassFromShapes()
	}
}

func (seg *Segment) Radius() float64 {
	return seg.radius
}

func (seg *Segment) TransformA() vec.Vec2 {
	return seg.transformA
}

func (seg *Segment) TransformB() vec.Vec2 {
	return seg.transformB
}

func (seg *Segment) SetEndpoints(a, b vec.Vec2) {
	seg.a = a
	seg.b = b
	seg.n = b.Sub(a).Normalize().Perp()

	mass := seg.massInfo.m
	seg.massInfo = NewSegmentMassInfo(seg.massInfo.m, seg.a, seg.b, seg.radius)
	if mass > 0 {
		seg.body.AccumulateMassFromShapes()
	}
}

func (seg *Segment) Normal() vec.Vec2 {
	return seg.n
}

func (seg *Segment) A() vec.Vec2 {
	return seg.a
}

func (seg *Segment) B() vec.Vec2 {
	return seg.b
}

func (seg *Segment) PointQuery(p vec.Vec2, info *PointQueryInfo) {
	closest := p.ClosestPointOnSegment(seg.transformA, seg.transformB)

	delta := p.Sub(closest)
	d := delta.Length()
	r := seg.radius
	g := delta.Scale(1 / d)

	info.Shape = seg.Shape
	if d != 0 {
		info.Point = closest.Add(g.Scale(r))
	} else {
		info.Point = closest
	}
	info.Distance = d - r

	// Use the segment's normal if the distance is very small.
	if d > MagicEpsilon {
		info.Gradient = g
	} else {
		info.Gradient = seg.n
	}
}

func (seg *Segment) SegmentQuery(a, b vec.Vec2, r2 float64, info *SegmentQueryInfo) {
	n := seg.transformN
	d := seg.transformA.Sub(a).Dot(n)
	r := seg.radius + r2

	var flippedN vec.Vec2
	if d > 0 {
		flippedN = n.Neg()
	} else {
		flippedN = n
	}
	segOffset := flippedN.Scale(r).Sub(a)

	// Make the endpoints relative to 'a' and move them by the thickness of the segment.
	segA := seg.transformA.Add(segOffset)
	segB := seg.transformB.Add(segOffset)
	delta := b.Sub(a)

	if delta.Cross(segA)*delta.Cross(segB) <= 0 {
		dOffset := d
		if d > 0 {
			dOffset -= r
		} else {
			dOffset += r
		}
		ad := -dOffset
		bd := delta.Dot(n) - dOffset

		if ad*bd < 0 {
			t := ad / (ad - bd)

			info.Shape = seg.Shape
			info.Point = a.Lerp(b, t).Sub(flippedN.Scale(r2))
			info.Normal = flippedN
			info.Alpha = t
		}
	} else if r != 0 {
		info1 := SegmentQueryInfo{nil, b, vec.Vec2{}, 1}
		info2 := SegmentQueryInfo{nil, b, vec.Vec2{}, 1}
		CircleSegmentQuery(seg.Shape, seg.transformA, seg.radius, a, b, r2, &info1)
		CircleSegmentQuery(seg.Shape, seg.transformB, seg.radius, a, b, r2, &info2)

		if info1.Alpha < info2.Alpha {
			*info = info1
		} else {
			*info = info2
		}
	}
}

func NewSegment(body *Body, a, b vec.Vec2, r float64) *Shape {
	segment := &Segment{
		a: a,
		b: b,
		n: b.Sub(a).Normalize().ReversePerp(),

		radius:   r,
		aTangent: vec.Vec2{},
		bTangent: vec.Vec2{},
	}
	segment.Shape = NewShape(segment, body, NewSegmentMassInfo(0, a, b, r))
	return segment.Shape
}

func NewSegmentMassInfo(mass float64, a, b vec.Vec2, r float64) *ShapeMassInfo {
	return &ShapeMassInfo{
		m:    mass,
		i:    MomentForBox(1, a.Distance(b)+2*r, 2*r),
		cog:  a.Lerp(b, 0.5),
		area: AreaForSegment(a, b, r),
	}
}
