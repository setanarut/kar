package cm

import (
	"fmt"
	"math"
)

const (
	Infinity     = math.MaxFloat64
	MagicEpsilon = 1e-5

	RadianConst = math.Pi / 180
	DegreeConst = 180 / math.Pi

	PooledBufferSize = 1024
)

// Arbiter states
const (
	// Arbiter is active and its the first collision.
	ArbiterStateFirstCollision = iota
	// Arbiter is active and its not the first collision.
	ArbiterStateNormal
	// Collision has been explicitly ignored.
	// Either by returning false from a begin collision handler or calling cmArbiterIgnore().
	ArbiterStateIgnore
	// Collison is no longer active. A space will cache an arbiter for up to cmSpace.collisionPersistence more steps.
	ArbiterStateCached
	// Collison arbiter is invalid because one of the shapes was removed.
	ArbiterStateInvalidated
)

// Collision Bitmask Category
const (
	NoGroup       uint = 0        // Value for group signifying that a shape is in no group.
	AllCategories uint = ^uint(0) // Value for Shape layers signifying that a shape is in every layer.
)

// ShapeFilterAll is s collision filter value for a shape that will collide with anything except SHAPE_FILTER_NONE.
var ShapeFilterAll = ShapeFilter{NoGroup, AllCategories, AllCategories}

// ShapeFilterNone is a collision filter value for a shape that does not collide with anything.
var ShapeFilterNone = ShapeFilter{NoGroup, ^AllCategories, ^AllCategories}

// CollisionBeginFunc is collision begin event function callback type.
//
// Returning false from a begin callback causes the collision to be ignored until the the separate callback is called when the objects stop colliding.
type CollisionBeginFunc func(arb *Arbiter, space *Space, userData interface{}) bool

// CollisionPreSolveFunc is collision pre-solve event function callback type.
//
// Returning false from a pre-step callback causes the collision to be ignored until the next step.
type CollisionPreSolveFunc func(arb *Arbiter, space *Space, userData interface{}) bool

// CollisionPostSolveFunc is collision post-solve event function callback type.
type CollisionPostSolveFunc func(arb *Arbiter, space *Space, userData interface{})

// CollisionSeparateFunc is collision separate event function callback type.
type CollisionSeparateFunc func(arb *Arbiter, space *Space, userData interface{})

type CollisionType uintptr

// CollisionHandler is struct that holds function callback pointers to configure custom collision handling.
// Collision handlers have a pair of types; when a collision occurs between two shapes that have these types, the collision handler functions are triggered.
type CollisionHandler struct {
	// Collision type identifier of the first shape that this handler recognizes.
	// In the collision handler callback, the shape with this type will be the first argument. Read only.
	TypeA CollisionType
	// Collision type identifier of the second shape that this handler recognizes.
	// In the collision handler callback, the shape with this type will be the second argument. Read only.
	TypeB CollisionType
	// This function is called when two shapes with types that match this collision handler begin colliding.
	BeginFunc CollisionBeginFunc
	// This function is called each step when two shapes with types that match this collision handler are colliding.
	// It's called before the collision solver runs so that you can affect a collision's outcome.
	PreSolveFunc CollisionPreSolveFunc
	// This function is called each step when two shapes with types that match this collision handler are colliding.
	// It's called after the collision solver runs so that you can read back information about the collision to trigger events in your game.
	PostSolveFunc CollisionPostSolveFunc
	// This function is called when two shapes with types that match this collision handler stop colliding.
	SeparateFunc CollisionSeparateFunc
	// This is a user definable context pointer that is passed to all of the collision handler functions.
	UserData interface{}
}

type Contact struct {
	r1, r2 Vec2

	nMass, tMass float64
	bounce       float64 // TODO: look for an alternate bounce solution

	jnAcc, jtAcc, jBias float64
	bias                float64

	hash HashValue
}

func (c *Contact) Clone() Contact {
	return Contact{
		r1:     c.r1,
		r2:     c.r2,
		nMass:  c.nMass,
		tMass:  c.tMass,
		bounce: c.bounce,
		jnAcc:  c.jnAcc,
		jtAcc:  c.jtAcc,
		jBias:  c.jBias,
		bias:   c.bias,
		hash:   c.hash,
	}
}

// CollisionInfo collision info struct
type CollisionInfo struct {
	a, b        *Shape
	collisionId uint32

	n     Vec2
	count int
	arr   []Contact
}

func (info *CollisionInfo) PushContact(p1, p2 Vec2, hash HashValue) {
	// if info.count > MAX_CONTACTS_PER_ARBITER {
	// 	log.Fatalln("Internal error: Tried to push too many contacts.")
	// }

	con := &info.arr[info.count]
	con.r1 = p1
	con.r2 = p2
	con.hash = hash

	info.count++
}

// ShapeMassInfo is mass info struct
type ShapeMassInfo struct {
	m, i, area float64
	cog        Vec2
}

// PointQueryInfo is point query info struct.
type PointQueryInfo struct {
	// The nearest shape, NULL if no shape was within range.
	Shape *Shape
	// The closest point on the shape's surface. (in world space coordinates)
	Point Vec2
	// The distance to the point. The distance is negative if the point is inside the shape.
	Distance float64
	// The gradient of the signed distance function.
	// The value should be similar to info.p/info.d, but accurate even for very small values of info.d.
	Gradient Vec2
}

// SegmentQueryInfo is segment query info struct.
type SegmentQueryInfo struct {
	// The shape that was hit, or NULL if no collision occurred.
	Shape *Shape
	// The point of impact.
	Point Vec2
	// The normal of the surface hit.
	Normal Vec2
	// The normalized distance along the query segment in the range [0, 1].
	Alpha float64
}

type SplittingPlane struct {
	v0, n Vec2
}

// ShapeFilter is fast collision filtering type that is used to determine if two objects collide before calling collision or query callbacks.
type ShapeFilter struct {
	// Two objects with the same non-zero group value do not collide.
	// This is generally used to group objects in a composite object together to disable self collisions.
	Group uint
	// A bitmask of user definable categories that this object belongs to.
	// The category/mask combinations of both objects in a collision must agree for a collision to occur.
	Categories uint
	// A bitmask of user definable category types that this object object collides with.
	// The category/mask combinations of both objects in a collision must agree for a collision to occur.
	Mask uint
}

// NewShapeFilter creates a new collision filter.
func NewShapeFilter(group, categories, mask uint) ShapeFilter {
	return ShapeFilter{group, categories, mask}
}

func (a ShapeFilter) Reject(b ShapeFilter) bool {
	// Reject the collision if:
	return (a.Group != 0 && a.Group == b.Group) ||
		// One of the category/mask combinations fails.
		(a.Categories&b.Mask) == 0 ||
		(b.Categories&a.Mask) == 0
}

// Mat2x2 is a 2x2 matrix type used for tensors and such.
type Mat2x2 struct {
	a, b, c, d float64
}

// Transform transforms Vector v
func (m *Mat2x2) Transform(v Vec2) Vec2 {
	return Vec2{v.X*m.a + v.Y*m.b, v.X*m.c + v.Y*m.d}
}

// MomentForBox calculates the moment of inertia for a solid box.
func MomentForBox(mass, width, height float64) float64 {
	return mass * (width*width + height*height) / 12.0
}

// MomentForBox2 calculates the moment of inertia for a solid box.
func MomentForBox2(mass float64, box BB) float64 {
	width := box.R - box.L
	height := box.T - box.B
	offset := Vec2{box.L + box.R, box.B + box.T}.Scale(0.5)

	// TODO: NaN when offset is 0 and m is INFINITY
	return MomentForBox(mass, width, height) + mass*offset.LengthSq()
}

// MomentForCircle calculates the moment of inertia for a circle.
//
// d1 and d2 are the inner and outer diameters. A solid circle has an inner diameter (d1) of 0.
func MomentForCircle(mass, d1, d2 float64, offset Vec2) float64 {
	return mass * (0.5*(d1*d1+d2*d2) + offset.LengthSq())
}

// MomentForSegment calculates the moment of inertia for a line segment.
//
// Beveling radius is not supported.
func MomentForSegment(mass float64, a, b Vec2, r float64) float64 {
	offset := a.Lerp(b, 0.5)
	length := b.Distance(a) + 2.0*r
	return mass * ((length*length+4.0*r*r)/12.0 + offset.LengthSq())
}

// MomentForPoly calculates the moment of inertia for a solid polygon shape assuming it's center of gravity is at it's centroid.
// The offset is added to each vertex.
func MomentForPoly(mass float64, count int, verts []Vec2, offset Vec2, r float64) float64 {
	if count == 2 {
		return MomentForSegment(mass, verts[0], verts[1], 0)
	}

	var sum1 float64
	var sum2 float64
	for i := 0; i < count; i++ {
		v1 := verts[i].Add(offset)
		v2 := verts[(i+1)%count].Add(offset)

		a := v2.Cross(v1)
		b := v1.Dot(v1) + v1.Dot(v2) + v2.Dot(v2)

		sum1 += a * b
		sum2 += a
	}

	return (mass * sum1) / (6.0 * sum2)
}

// AreaForCircle returns area of a hollow circle.
//
// r1 and r2 are the inner and outer diameters. A solid circle has an inner diameter of 0.
func AreaForCircle(r1, r2 float64) float64 {
	return math.Pi * math.Abs(r1*r1-r2*r2)
}

// AreaForSegment calculates the area of a fattened (capsule shaped) line segment.
func AreaForSegment(a, b Vec2, r float64) float64 {
	return r * (math.Pi*r + 2.0*a.Distance(b))
}

// AreaForPoly calculates the signed area of a polygon.
//
// A Clockwise winding gives positive area. This is probably backwards from what you expect, but matches Chipmunk's the winding for poly shapes.
func AreaForPoly(count int, verts []Vec2, r float64) float64 {
	var area float64
	var perimeter float64
	for i := 0; i < count; i++ {
		v1 := verts[i]
		v2 := verts[(i+1)%count]

		area += v1.Cross(v2)
		perimeter += v1.Distance(v2)
	}

	return r*(math.Pi*math.Abs(r)+perimeter) + area/2.0
}

// CentroidForPoly calculates the natural centroid of a polygon.
func CentroidForPoly(count int, verts []Vec2) Vec2 {
	var sum float64
	vsum := Vec2{}

	for i := 0; i < count; i++ {
		v1 := verts[i]
		v2 := verts[(i+1)%count]
		cross := v1.Cross(v2)

		sum += cross
		vsum = vsum.Add(v1.Add(v2).Scale(cross))
	}

	return vsum.Scale(1.0 / (3.0 * sum))
}

// DebugInfo returns info of space
func DebugInfo(space *Space) string {
	arbiters := len(space.Arbiters)
	points := 0

	for i := 0; i < arbiters; i++ {
		points += int(space.Arbiters[i].count)
	}

	constraints := len(space.constraints) + points*int(space.Iterations)
	if arbiters > maxArbiters {
		maxArbiters = arbiters
	}
	if points > maxPoints {
		maxPoints = points
	}
	if constraints > maxConstraints {
		maxConstraints = constraints
	}

	var ke float64
	for _, body := range space.DynamicBodies {
		if body.mass == Infinity || body.moi == Infinity {
			continue
		}
		ke += body.mass*body.vel.Dot(body.vel) + body.moi*body.w*body.w
	}

	return fmt.Sprintf(`Arbiters: %d (%d) - Contact Points: %d (%d)
Other Constraints: %d, Iterations: %d
Constraints x Iterations: %d (%d)
KE: %e`, arbiters, maxArbiters,
		points, maxPoints, len(space.constraints), space.Iterations, constraints, maxConstraints, ke)
}

func k_scalar_body(body *Body, r, n Vec2) float64 {
	rcn := r.Cross(n)
	return body.m_inv + body.moi_inv*rcn*rcn
}

func k_scalar(a, b *Body, r1, r2, n Vec2) float64 {
	return k_scalar_body(a, r1, n) + k_scalar_body(b, r2, n)
}

func normal_relative_velocity(a, b *Body, r1, r2, n Vec2) float64 {
	return relative_velocity(a, b, r1, r2).Dot(n)
}

func k_tensor(a, b *Body, r1, r2 Vec2) Mat2x2 {
	m_sum := a.m_inv + b.m_inv

	// start with Identity*m_sum
	k11 := m_sum
	k12 := 0.0
	k21 := 0.0
	k22 := m_sum

	// add the influence from r1
	a_i_inv := a.moi_inv
	r1xsq := r1.X * r1.X * a_i_inv
	r1ysq := r1.Y * r1.Y * a_i_inv
	r1nxy := -r1.X * r1.Y * a_i_inv
	k11 += r1ysq
	k12 += r1nxy
	k21 += r1nxy
	k22 += r1xsq

	// add the influence from r2
	b_i_inv := b.moi_inv
	r2xsq := r2.X * r2.X * b_i_inv
	r2ysq := r2.Y * r2.Y * b_i_inv
	r2nxy := -r2.X * r2.Y * b_i_inv
	k11 += r2ysq
	k12 += r2nxy
	k21 += r2nxy
	k22 += r2xsq

	// invert
	det := k11*k22 - k12*k21
	// if det == 0.0 {
	// 	log.Fatalln("Unsolvable constraint")
	// }

	det_inv := 1.0 / det
	return Mat2x2{
		k22 * det_inv, -k12 * det_inv,
		-k21 * det_inv, k11 * det_inv,
	}
}

func bias_coef(errorBias, dt float64) float64 {
	return 1.0 - math.Pow(errorBias, dt)
}

var maxArbiters, maxPoints, maxConstraints int
