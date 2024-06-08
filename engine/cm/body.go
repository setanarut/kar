package cm

import (
	"fmt"
	"log"
	"math"
)

// Body types
const (
	BODY_DYNAMIC = iota
	BODY_KINEMATIC
	BODY_STATIC
)

var bodyCur int = 0

// BodyVelocityFunc is rigid body velocity update function type.
type BodyVelocityFunc func(body *Body, gravity Vec2, damping float64, dt float64)

// BodyPositionFunc is rigid body position update function type.
type BodyPositionFunc func(body *Body, dt float64)

type Body struct {
	// UserData is an object that this constraint is associated with.
	//
	// You can use this get a reference to your game object or controller object from within callbacks.
	UserData interface{}

	id            int              // Body id
	velocity_func BodyVelocityFunc // Integration function
	position_func BodyPositionFunc // Integration function

	mass  float64 // Mass
	m_inv float64 // Mass inverse

	moi     float64 // Moment of inertia
	moi_inv float64 // Inverse of moment of inertia i

	cog Vec2 // Center of gravity

	position Vec2 // Position
	vel      Vec2 // Velocity
	force    Vec2 // Force

	angle  float64 // Angle (radians)
	w      float64 // Angular velocity,
	torque float64 // Torque (radians)

	transform Transform

	// "pseudo-velocities" used for eliminating overlap.
	// Erin Catto has some papers that talk about what these are.
	v_bias Vec2
	w_bias float64

	space *Space

	shapeList      []*Shape
	arbiterList    *Arbiter
	constraintList *Constraint

	sleepingRoot     *Body
	sleepingNext     *Body
	sleepingIdleTime float64
}

// String returns body id as string
func (b Body) String() string {
	return fmt.Sprint("Body ", b.id, ", Shapes ", b.shapeList)
}

// Shapes returns shapes attached to this body
func (b *Body) Shapes() []*Shape {
	return b.shapeList
}

// FirstShape returns first shape attached to this body
func (b *Body) FirstShape() *Shape {
	return b.shapeList[0]
}

// NewBody Initializes a rigid body with the given mass and moment of inertia.
//
// Guessing the moment of inertia is usually a bad idea. Use the moment estimation functions MomentFor*().
func NewBody(mass, moment float64) *Body {
	body := &Body{
		id:            bodyCur,
		cog:           Vec2{},
		position:      Vec2{},
		vel:           Vec2{},
		force:         Vec2{},
		v_bias:        Vec2{},
		transform:     NewTransformIdentity(),
		velocity_func: BodyUpdateVelocity,
		position_func: BodyUpdatePosition,
	}
	bodyCur++

	body.SetMass(mass)
	body.SetMoment(moment)
	body.SetAngle(0)

	return body
}

// NewStaticBody allocates and initializes a Body, and set it as a static body.
func NewStaticBody() *Body {
	body := NewBody(0, 0)
	body.SetType(BODY_STATIC)
	return body
}

// NewKinematicBody allocates and initializes a Body, and set it as a kinematic body.
func NewKinematicBody() *Body {
	body := NewBody(0, 0)
	body.SetType(BODY_KINEMATIC)
	return body
}

// SetAngle sets the angle of body.
func (body *Body) SetAngle(angle float64) {
	body.Activate()
	body.angle = angle
	body.SetTransform(body.position, angle)
}

// Moment returns moment of inertia of the body.
func (body Body) Moment() float64 {
	return body.moi
}

// SetMoment sets moment of inertia of the body.
func (body *Body) SetMoment(moment float64) {
	body.Activate()
	body.moi = moment
	body.moi_inv = 1 / moment
}

// Mass returns mass of the body
func (body *Body) Mass() float64 {
	return body.mass
}

// SetMass sets mass of the body
func (body *Body) SetMass(mass float64) {
	body.Activate()
	body.mass = mass
	body.m_inv = 1 / mass
}

// IdleTime returns sleeping idle time of the body
func (body *Body) IdleTime() float64 {
	return body.sleepingIdleTime
}

// SetType sets the type of the body.
func (body *Body) SetType(newType int) {
	oldType := body.GetType()
	if oldType == newType {
		return
	}

	if newType == BODY_STATIC {
		body.sleepingIdleTime = Infinity
	} else {
		body.sleepingIdleTime = 0
	}

	if newType == BODY_DYNAMIC {
		body.mass = 0
		body.moi = 0
		body.m_inv = Infinity
		body.moi_inv = Infinity

		body.AccumulateMassFromShapes()
	} else {
		body.mass = Infinity
		body.moi = Infinity
		body.m_inv = 0
		body.moi_inv = 0

		body.vel = Vec2{}
		body.w = 0
	}

	// If the body is added to a space already, we'll need to update some space data structures.
	if body.space == nil {
		return
	}
	if body.space.locked != 0 {
		log.Fatalln("Space is locked")
	}

	if oldType != BODY_STATIC {
		body.Activate()
	}

	if oldType == BODY_STATIC {
		for i, b := range body.space.StaticBodies {
			if b == body {
				body.space.StaticBodies = append(body.space.StaticBodies[:i], body.space.StaticBodies[i+1:]...)
				break
			}
		}
		body.space.DynamicBodies = append(body.space.DynamicBodies, body)
	} else if newType == BODY_STATIC {
		for i, b := range body.space.DynamicBodies {
			if b == body {
				body.space.DynamicBodies = append(body.space.DynamicBodies[:i], body.space.DynamicBodies[i+1:]...)
				break
			}
		}
		body.space.StaticBodies = append(body.space.StaticBodies, body)
	}

	var fromIndex, toIndex *SpatialIndex
	if oldType == BODY_STATIC {
		fromIndex = body.space.staticShapes
	} else {
		fromIndex = body.space.dynamicShapes
	}

	if newType == BODY_STATIC {
		toIndex = body.space.staticShapes
	} else {
		toIndex = body.space.dynamicShapes
	}

	if oldType != newType {
		for _, shape := range body.shapeList {
			fromIndex.class.Remove(shape, shape.hashid)
			toIndex.class.Insert(shape, shape.hashid)
		}
	}
}

// GetType returns the type of the body.
func (body *Body) GetType() int {
	if body.sleepingIdleTime == Infinity {
		return BODY_STATIC
	}
	if body.mass == Infinity {
		return BODY_KINEMATIC
	}
	return BODY_DYNAMIC
}

// AccumulateMassFromShapes should *only* be called when shapes with mass info are modified, added or removed.
func (body *Body) AccumulateMassFromShapes() {
	if body == nil || body.GetType() != BODY_DYNAMIC {
		return
	}

	body.mass = 0
	body.moi = 0
	body.cog = Vec2{}

	// cache position, realign at the end
	pos := body.Position()

	for _, shape := range body.shapeList {
		info := shape.MassInfo()
		m := info.m

		if info.m > 0 {
			msum := body.mass + m
			body.moi += m*info.i + body.cog.DistanceSq(info.cog)*(m*body.mass)/msum
			body.cog = body.cog.Lerp(info.cog, m/msum)
			body.mass = msum
		}
	}

	body.m_inv = 1.0 / body.mass
	body.moi_inv = 1.0 / body.moi

	body.SetPosition(pos)
}

// CenterOfGravity returns the offset of the center of gravity in body local coordinates.
func (body Body) CenterOfGravity() Vec2 {
	return body.cog
}

// Angle returns the angle of the body.
func (body *Body) Angle() float64 {
	return body.angle
}

// Rotation returns the rotation vector of the body.
//
// (The x basis vector of it's transform.)
func (body *Body) Rotation() Vec2 {
	return Vec2{body.transform.a, body.transform.b}
}

// Position returns the position of the body.
func (body *Body) Position() Vec2 {
	return body.transform.Point(Vec2{})
}

// SetPosition sets the position of the body.
func (body *Body) SetPosition(position Vec2) {
	body.Activate()
	body.position = body.transform.Vect(body.cog).Add(position)
	body.SetTransform(body.position, body.angle)
}

// Velocity returns the velocity of the body.
func (body *Body) Velocity() Vec2 {
	return body.vel
}

// SetVelocity sets the velocity of the body.
//
// Shorthand for Body.SetVelocityVector()
func (body *Body) SetVelocity(x, y float64) {
	body.Activate()
	body.vel = Vec2{x, y}
}

// SetVelocityVector sets the velocity of the body
func (body *Body) SetVelocityVector(v Vec2) {
	body.Activate()
	body.vel = v
}

// UpdateVelocity is the default velocity integration function.
func (body *Body) UpdateVelocity(gravity Vec2, damping, dt float64) {
	if body.GetType() == BODY_KINEMATIC {
		return
	}
	// if body.mass < 0 && body.moi < 0 {
	// 	log.Fatalln("Body's mass and moment must be positive")
	// }

	body.vel = body.vel.Scale(damping).Add(gravity.Add(body.force.Scale(body.m_inv)).Scale(dt))
	body.w = body.w*damping + body.torque*body.moi_inv*dt

	body.force = Vec2{}
	body.torque = 0
}

// Force returns the force applied to the body for the next time step.
func (body *Body) Force() Vec2 {
	return body.force
}

// SetForce sets the force applied to the body for the next time step.
func (body *Body) SetForce(force Vec2) {
	body.Activate()
	body.force = force
}

// Torque returns the torque applied to the body for the next time step.
func (body *Body) Torque() float64 {
	return body.torque
}

// SetTorque sets the torque applied to the body for the next time step.
func (body *Body) SetTorque(torque float64) {
	body.Activate()
	body.torque = torque
}

// AngularVelocity returns the angular velocity of the body.
func (body *Body) AngularVelocity() float64 {
	return body.w
}

// SetAngularVelocity sets the angular velocity of the body.
func (body *Body) SetAngularVelocity(angularVelocity float64) {
	body.Activate()
	body.w = angularVelocity
}

// SetTransform sets transform
func (body *Body) SetTransform(p Vec2, a float64) {
	rot := Vec2{math.Cos(a), math.Sin(a)}
	c := body.cog

	body.transform = NewTransformTranspose(
		rot.X, -rot.Y, p.X-(c.X*rot.X-c.Y*rot.Y),
		rot.Y, rot.X, p.Y-(c.X*rot.Y+c.Y*rot.X),
	)
}

// Transform returns body's transform
func (body *Body) Transform() Transform {
	return body.transform
}

// Activate wakes up a sleeping or idle body.
func (body *Body) Activate() {
	if !(body != nil && body.GetType() == BODY_DYNAMIC) {
		return
	}

	body.sleepingIdleTime = 0

	root := body.ComponentRoot()
	if root != nil && root.IsSleeping() {

		// if root.GetType() != BODY_DYNAMIC {
		// 	log.Fatalln("Non-dynamic root")
		// }

		space := root.space
		// in the chipmunk code they shadow body, so here I am not
		bodyToo := root
		for bodyToo != nil {
			next := bodyToo.sleepingNext
			bodyToo.sleepingIdleTime = 0
			bodyToo.sleepingRoot = nil
			bodyToo.sleepingNext = nil
			space.Activate(bodyToo)

			bodyToo = next
		}

		for i := 0; i < len(space.sleepingComponents); i++ {
			if space.sleepingComponents[i] == root {
				space.sleepingComponents = append(space.sleepingComponents[:i], space.sleepingComponents[i+1:]...)
				break
			}
		}
	}

	for arbiter := body.arbiterList; arbiter != nil; arbiter = arbiter.Next(body) {
		// Reset the idle timer of things the body is touching as well.
		// That way things don't get left hanging in the air.
		var other *Body
		if arbiter.body_a == body {
			other = arbiter.body_b
		} else {
			other = arbiter.body_a
		}
		if other.GetType() != BODY_STATIC {
			other.sleepingIdleTime = 0
		}
	}
}

// ActivateStatic wakes up any sleeping or idle bodies touching this static body.
func (body *Body) ActivateStatic(filter *Shape) {
	// if body.GetType() != BODY_STATIC {
	// 	log.Fatalln("Body is not static")
	// }

	for arb := body.arbiterList; arb != nil; arb = arb.Next(body) {
		if filter == nil || filter == arb.a || filter == arb.b {
			if arb.body_a == body {
				arb.body_b.Activate()
			} else {
				arb.body_a.Activate()
			}
		}
	}
}

// IsSleeping returns true if the body is sleeping.
func (body *Body) IsSleeping() bool {
	return body.sleepingRoot != nil
}

// AddShape adds shape to the body and returns added shape
func (body *Body) AddShape(shape *Shape) *Shape {
	body.shapeList = append(body.shapeList, shape)
	if shape.MassInfo().m > 0 {
		body.AccumulateMassFromShapes()
	}
	return shape
}

// KineticEnergy returns the kinetic energy of this body.
func (body *Body) KineticEnergy() float64 {
	// Need to do some fudging to avoid NaNs
	vsq := body.vel.Dot(body.vel)
	wsq := body.w * body.w
	var a, b float64
	if vsq != 0 {
		a = vsq * body.mass
	}
	if wsq != 0 {
		b = wsq * body.moi
	}
	return a + b
}

func (body *Body) PushArbiter(arb *Arbiter) {
	next := body.arbiterList
	arb.ThreadForBody(body).next = next
	if next != nil {
		next.ThreadForBody(body).prev = arb
	}
	body.arbiterList = arb
}

func (root *Body) ComponentAdd(body *Body) {
	body.sleepingRoot = root

	if body != root {
		body.sleepingNext = root.sleepingNext
		root.sleepingNext = body
	}
}

func (body *Body) ComponentRoot() *Body {
	if body != nil {
		return body.sleepingRoot
	}
	return nil
}

// WorldToLocal converts from world to body local Coordinates.
//
// Convert a point in body local coordinates to world (absolute) coordinates.
func (body *Body) WorldToLocal(point Vec2) Vec2 {
	return NewTransformRigidInverse(body.transform).Point(point)
}

// LocalToWorld converts from body local to world coordinates.
//
// Convert a point in world (absolute) coordinates to body local coordinates affected by the position and rotation of the rigid body.
func (body *Body) LocalToWorld(point Vec2) Vec2 {
	return body.transform.Point(point)
}

// ApplyForceAtWorldPoint applies a force at world point.
func (body *Body) ApplyForceAtWorldPoint(force, point Vec2) {
	body.Activate()
	body.force = body.force.Add(force)

	r := point.Sub(body.transform.Point(body.cog))
	body.torque += r.Cross(force)
}

// ApplyForceAtLocalPoint applies a force at local point.
func (body *Body) ApplyForceAtLocalPoint(force, point Vec2) {
	body.ApplyForceAtWorldPoint(body.transform.Vect(force), body.transform.Point(point))
}

// ApplyImpulseAtWorldPoint applies impulse at world point
func (body *Body) ApplyImpulseAtWorldPoint(impulse, point Vec2) {
	body.Activate()

	r := point.Sub(body.transform.Point(body.cog))
	apply_impulse(body, impulse, r)
}

// ApplyImpulseAtLocalPoint applies impulse at local point
func (body *Body) ApplyImpulseAtLocalPoint(impulse, point Vec2) {
	body.ApplyImpulseAtWorldPoint(body.transform.Vect(impulse), body.transform.Point(point))
}

// VelocityAtLocalPoint returns the velocity of a point on a body.
//
// Get the world (absolute) velocity of a point on a rigid body specified in body local coordinates.
func (body *Body) VelocityAtLocalPoint(point Vec2) Vec2 {
	r := body.transform.Vect(point.Sub(body.cog))
	return body.vel.Add(r.Perp().Scale(body.w))
}

// VelocityAtWorldPoint returns the velocity of a point on a body.
//
// Get the world (absolute) velocity of a point on a rigid body specified in world coordinates.
func (body *Body) VelocityAtWorldPoint(point Vec2) Vec2 {
	r := point.Sub(body.transform.Point(body.cog))
	return body.vel.Add(r.Perp().Scale(body.w))
}

// RemoveConstraint removes constraint from the body.
func (body *Body) RemoveConstraint(constraint *Constraint) {
	body.constraintList = filterConstraints(body.constraintList, body, constraint)
}

// RemoveShape removes collision shape from the body.
func (body *Body) RemoveShape(shape *Shape) {
	for i, s := range body.shapeList {
		if s == shape {
			// leak-free delete from slice
			last := len(body.shapeList) - 1
			body.shapeList[i] = body.shapeList[last]
			body.shapeList[last] = nil
			body.shapeList = body.shapeList[:last]
			break
		}
	}
	if body.GetType() == BODY_DYNAMIC && shape.massInfo.m > 0 {
		body.AccumulateMassFromShapes()
	}
}

// SetVelocityUpdateFunc sets the callback used to update a body's velocity.
func (body *Body) SetVelocityUpdateFunc(f BodyVelocityFunc) {
	body.velocity_func = f
}

// SetPositionUpdateFunc sets the callback used to update a body's position.
func (body *Body) SetPositionUpdateFunc(f BodyPositionFunc) {
	body.position_func = f
}

// EachArbiter calls f once for each arbiter that is currently active on the body.
func (body *Body) EachArbiter(f func(*Arbiter)) {
	arb := body.arbiterList
	for arb != nil {
		next := arb.Next(body)
		swapped := arb.swapped

		arb.swapped = body == arb.body_b
		f(arb)

		arb.swapped = swapped
		arb = next
	}
}

// EachShape calls f once for each shape attached to this body
func (body *Body) EachShape(f func(*Shape)) {
	for i := 0; i < len(body.shapeList); i++ {
		f(body.shapeList[i])
	}
}

// EachConstraint calls f once for each constraint attached to this body
func (body *Body) EachConstraint(f func(*Constraint)) {
	constraint := body.constraintList
	for constraint != nil {
		next := constraint.Next(body)
		f(constraint)
		constraint = next
	}
}

func filterConstraints(node *Constraint, body *Body, filter *Constraint) *Constraint {
	if node == filter {
		return node.Next(body)
	} else if node.bodyA == body {
		node.nextA = filterConstraints(node.nextA, body, filter)
	} else {
		node.nextB = filterConstraints(node.nextB, body, filter)
	}
	return node
}

// BodyUpdateVelocity is default velocity integration function.
func BodyUpdateVelocity(body *Body, gravity Vec2, damping, dt float64) {
	if body.GetType() == BODY_KINEMATIC {
		return
	}

	body.vel = body.vel.Scale(damping).Add(gravity.Add(body.force.Scale(body.m_inv)).Scale(dt))
	body.w = body.w*damping + body.torque*body.moi_inv*dt

	body.force = Vec2{}
	body.torque = 0
}

// BodyUpdatePosition is default position integration function.
func BodyUpdatePosition(body *Body, dt float64) {
	body.position = body.position.Add(body.vel.Add(body.v_bias).Scale(dt))
	body.angle = body.angle + (body.w+body.w_bias)*dt
	body.SetTransform(body.position, body.angle)

	body.v_bias = Vec2{}
	body.w_bias = 0
}
