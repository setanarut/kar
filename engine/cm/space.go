package cm

import (
	"kar/engine/vec"
	"log"
	"math"
	"sync"
	"unsafe"
)

const MAX_CONTACTS_PER_ARBITER = 2
const CONTACTS_BUFFER_SIZE = 1024

type Space struct {

	// Iterations is number of iterations to use in the impulse solver to solve contacts and other constrain.
	// Must be non-zero.
	Iterations uint

	// IdleSpeedThreshold is speed threshold for a body to be considered idle.
	// The default value of 0 means to let the space guess a good threshold based on gravity.
	IdleSpeedThreshold float64

	// SleepTimeThreshold is time a group of bodies must remain idle in order to fall asleep.
	// Enabling sleeping also implicitly enables the the contact graph.
	// The default value of INFINITY disables the sleeping algorithm.
	SleepTimeThreshold float64

	// StaticBody is the Space provided static body for a given Space.
	// This is merely provided for convenience and you are not required to use it.
	StaticBody *Body

	// Gravity to pass to rigid bodies when integrating velocity.
	Gravity vec.Vec2

	// Damping rate expressed as the fraction of velocity bodies retain each second.
	//
	// A value of 0.9 would mean that each body's velocity will drop 10% per second.
	// The default value is 1.0, meaning no Damping is applied.
	// @note This Damping value is different than those of DampedSpring and DampedRotarySpring.
	Damping float64

	// CollisionSlop is amount of encouraged penetration between colliding shapes.
	//
	// Used to reduce oscillating contacts and keep the collision cache warm.
	// Defaults to 0.1. If you have poor simulation quality,
	// increase this number as much as possible without allowing visible amounts of overlap.
	CollisionSlop float64

	// CollisionBias determines how fast overlapping shapes are pushed apart.
	//
	// Expressed as a fraction of the error remaining after each second.
	// Defaults to math.Pow(0.9, 60) meaning that Chipmunk fixes 10% of overlap each frame at 60Hz.
	CollisionBias float64

	// Number of frames that contact information should persist.
	// Defaults to 3. There is probably never a reason to change this value.
	CollisionPersistence uint

	Arbiters          []*Arbiter
	DynamicBodies     []*Body
	StaticBodies      []*Body
	PostStepCallbacks []*PostStepCallback
	UserData          interface{}
	staticShapes      *SpatialIndex
	dynamicShapes     *SpatialIndex

	stamp              uint
	curr_dt            float64
	rousedBodies       []*Body
	sleepingComponents []*Body
	shapeIDCounter     uint
	constraints        []*Constraint
	contactBuffersHead *ContactBuffer
	cachedArbiters     *HashSet[ShapePair, *Arbiter]
	pooledArbiters     sync.Pool
	locked             int
	usesWildcards      bool
	collisionHandlers  *HashSet[*CollisionHandler, *CollisionHandler]
	defaultHandler     *CollisionHandler
	skipPostStep       bool
}

// NewSpace allocates and initializes a Space
func NewSpace() *Space {
	space := &Space{
		Iterations:           10,
		IdleSpeedThreshold:   0.0,
		SleepTimeThreshold:   math.MaxFloat64,
		StaticBody:           NewBody(0, 0),
		Gravity:              vec.Vec2{},
		Damping:              1.0,
		CollisionSlop:        0.1,
		CollisionBias:        math.Pow(0.9, 60),
		CollisionPersistence: 3,
		DynamicBodies:        []*Body{},
		StaticBodies:         []*Body{},
		Arbiters:             []*Arbiter{},
		locked:               0,
		stamp:                0,
		shapeIDCounter:       1,
		staticShapes:         NewBBTree(ShapeGetBB, nil),
		sleepingComponents:   []*Body{},
		rousedBodies:         []*Body{},
		cachedArbiters:       NewHashSet[ShapePair, *Arbiter](arbiterSetEql),
		pooledArbiters:       sync.Pool{New: func() interface{} { return &Arbiter{} }},
		constraints:          []*Constraint{},
		collisionHandlers: NewHashSet[*CollisionHandler, *CollisionHandler](func(a, b *CollisionHandler) bool {
			if a.TypeA == b.TypeA && a.TypeB == b.TypeB {
				return true
			}
			if a.TypeB == b.TypeA && a.TypeA == b.TypeB {
				return true
			}
			return false
		}),
		PostStepCallbacks: []*PostStepCallback{},
		defaultHandler:    &CollisionHandlerDoNothing,
	}
	for i := 0; i < PooledBufferSize; i++ {
		space.pooledArbiters.Put(&Arbiter{})
	}
	space.dynamicShapes = NewBBTree(ShapeGetBB, space.staticShapes)
	space.dynamicShapes.class.(*BBTree).velocityFunc = BBTreeVelocityFunc(ShapeVelocityFunc)
	space.StaticBody.SetType(BODY_STATIC)
	return space
}

// func (space *Space) Arbiters() []*Arbiter {
// 	return space.arbiters
// }

// // DynamicBodies returns dynamic body list of the space
// func (space *Space) DynamicBodies() []*Body {
// 	return space.dynamicBodies
// }

// // StaticBodies returns static body list of the space
// func (space *Space) StaticBodies() []*Body {
// 	return space.staticBodies
// }

// DynamicBodyCount returns the total number of dynamic bodies in space
func (space *Space) DynamicBodyCount() int {
	return len(space.DynamicBodies)
}

// StaticBodyCount returns the total number of static bodies in space
func (space *Space) StaticBodyCount() int {
	return len(space.StaticBodies)
}

// SetGravity sets gravity and wake up all of the sleeping bodies since the gravity changed.
func (space *Space) SetGravity(gravity vec.Vec2) {
	space.Gravity = gravity

	// Wake up all of the bodies since the gravity changed.
	for _, component := range space.sleepingComponents {
		component.Activate()
	}
}

func (space *Space) SetStaticBody(body *Body) {
	if space.StaticBody != nil {
		space.StaticBody.space = nil
		panic("Internal Error: Changing the designated static body while the old one still had shapes attached.")
	}
	space.StaticBody = body
	body.space = space
}

func (space *Space) Activate(body *Body) {
	// if body.GetType() != BODY_DYNAMIC {
	// 	log.Fatalln("Attempting to activate a non-dynamic body")
	// }

	if space.locked != 0 {
		if !Contains(space.rousedBodies, body) {
			space.rousedBodies = append(space.rousedBodies, body)
		}
		return
	}

	// if body.sleepingRoot != nil && body.sleepingNext != nil {
	// 	log.Fatalln("Activating body non-NULL node pointers.")
	// }

	space.DynamicBodies = append(space.DynamicBodies, body)

	for _, shape := range body.shapeList {
		space.staticShapes.class.Remove(shape, shape.hashid)
		space.dynamicShapes.class.Insert(shape, shape.hashid)
	}

	for arbiter := body.arbiterList; arbiter != nil; arbiter = arbiter.Next(body) {
		bodyA := arbiter.body_a

		// Arbiters are shared between two bodies that are always woken up together.
		// You only want to restore the arbiter once, so bodyA is arbitrarily chosen to own the arbiter.
		// The edge case is when static bodies are involved as the static bodies never actually sleep.
		// If the static body is bodyB then all is good. If the static body is bodyA, that can easily be checked.
		if body == bodyA || bodyA.GetType() == BODY_STATIC {
			numContacts := arbiter.count
			contacts := arbiter.contacts

			// Restore contact values back to the space's contact buffer memory
			arbiter.contacts = space.ContactBufferGetArray()[:numContacts]
			copy(arbiter.contacts, contacts)
			space.PushContacts(numContacts)

			// reinsert the arbiter into the arbiter cache
			a := arbiter.a
			b := arbiter.b
			shapePair := ShapePair{a, b}
			arbHashId := HashPair(HashValue(unsafe.Pointer(a)), HashValue(unsafe.Pointer(b)))
			space.cachedArbiters.Insert(arbHashId, shapePair, func(_ ShapePair) *Arbiter {
				return arbiter
			})

			// update arbiters state
			arbiter.stamp = space.stamp
			space.Arbiters = append(space.Arbiters, arbiter)
		}
	}

	for constraint := body.constraintList; constraint != nil; constraint = constraint.Next(body) {
		if body == constraint.bodyA || constraint.bodyA.GetType() == BODY_STATIC {
			space.constraints = append(space.constraints, constraint)
		}
	}
}

func (space *Space) Deactivate(body *Body) {
	// if body.GetType() != BODY_DYNAMIC {
	// 	log.Fatalln("Attempting to deactivate non-dynamic body.")
	// }
	for i, v := range space.DynamicBodies {
		if v == body {
			space.DynamicBodies = append(space.DynamicBodies[:i], space.DynamicBodies[i+1:]...)
			break
		}
	}

	for _, shape := range body.shapeList {
		space.dynamicShapes.class.Remove(shape, shape.hashid)
		space.staticShapes.class.Insert(shape, shape.hashid)
	}

	for arb := body.arbiterList; arb != nil; arb = ArbiterNext(arb, body) {
		bodyA := arb.body_a
		if body == bodyA || bodyA.GetType() == BODY_STATIC {
			space.UncacheArbiter(arb)
			// Save contact values to a new block of memory so they won't time out
			contacts := make([]Contact, arb.count)
			copy(contacts, arb.contacts[:arb.count])
			arb.contacts = contacts

		}
	}

	for constraint := body.constraintList; constraint != nil; constraint = constraint.Next(body) {
		bodyA := constraint.bodyA
		if body == bodyA || bodyA.GetType() == BODY_STATIC {
			for i, c := range space.constraints {
				if c == constraint {
					space.constraints = append(space.constraints[0:i], space.constraints[i+1:]...)
				}
			}
		}
	}
}

// AddShape adds a collision shape to the simulation.
//
// If the shape is attached to a static body, it will be added as a static shape
func (space *Space) AddShape(shape *Shape) *Shape {
	var body *Body = shape.Body()

	// if shape.space == space {
	// 	log.Fatalln("You have already added this shape to this space. You must not add it a second time.")
	// }
	// if shape.space != nil {
	// 	log.Fatalln("You have already added this shape to another space. You cannot add it to a second.")
	// }
	// if space.locked != 0 {
	// 	log.Fatalln("This operation cannot be done safely during a call to SpaceStep() or during a query. Put these calls into a post-step callback.")
	// }

	isStatic := body.GetType() == BODY_STATIC
	if !isStatic {
		body.Activate()
	}
	body.AddShape(shape)

	shape.SetHashId(HashValue(space.shapeIDCounter))
	space.shapeIDCounter += 1
	shape.Update(body.transform)

	if isStatic {
		space.staticShapes.class.Insert(shape, shape.HashId())
	} else {
		space.dynamicShapes.class.Insert(shape, shape.HashId())
	}
	shape.SetSpace(space)

	return shape
}

// RemoveShape removes a collision shape from the simulation.
func (space *Space) RemoveShape(shape *Shape) {
	body := shape.body

	// if !space.ContainsShape(shape) {
	// 	log.Fatalln("Shape is not in space")
	// }

	// if space.locked != 0 {
	// 	log.Fatalln("space.locked is not zero")
	// }

	isStatic := body.GetType() == BODY_STATIC
	if isStatic {
		body.ActivateStatic(shape)
	} else {
		body.Activate()
	}

	body.RemoveShape(shape)
	space.FilterArbiters(body, shape)
	if isStatic {
		space.staticShapes.class.Remove(shape, shape.hashid)
	} else {
		space.dynamicShapes.class.Remove(shape, shape.hashid)
	}
	shape.space = nil
	shape.hashid = 0
}

// AddBody adds a body to the space if not in space.
func (space *Space) AddBody(body *Body) *Body {
	// if space.ContainsBody(body) {
	// 	log.Fatalln("Body already added to space")
	// }

	// if body.space != nil {
	// 	log.Fatalln("Body already added to another space")
	// }
	if body.GetType() == BODY_STATIC {
		space.StaticBodies = append(space.StaticBodies, body)
	} else {
		space.DynamicBodies = append(space.DynamicBodies, body)
	}
	body.space = space
	return body
}

// AddBodyWidthShapes adds a body to the space with body's Shapes.
func (space *Space) AddBodyWidthShapes(body *Body) *Body {
	for _, v := range body.Shapes() {
		space.AddShape(v)
	}
	return space.AddBody(body)

}

// RemoveBody removes a body from the simulation
func (space *Space) RemoveBody(body *Body) {
	// if body == space.StaticBody {
	// 	log.Fatalln("Body is space.Staticbody")
	// }
	// if !space.ContainsBody(body) {
	// 	log.Fatalln("Body is not in space")
	// }
	// if space.locked != 0 {
	// 	log.Fatalln("Space is locked")
	// }
	body.Activate()
	if body.GetType() == BODY_STATIC {
		for i, b := range space.StaticBodies {
			if b == body {
				space.StaticBodies = append(space.StaticBodies[:i], space.StaticBodies[i+1:]...)
				break
			}
		}
	} else {
		for i, b := range space.DynamicBodies {
			if b == body {
				space.DynamicBodies = append(space.DynamicBodies[:i], space.DynamicBodies[i+1:]...)
				break
			}
		}
	}
	body.space = nil
}

// RemoveBodyWithShapes removes a body and body's shapes from the simulation
func (space *Space) RemoveBodyWithShapes(body *Body) {
	body.EachShape(func(s *Shape) {
		space.RemoveShape(s)
	})
	space.RemoveBody(body)
}

func (space *Space) AddConstraint(constraint *Constraint) *Constraint {
	// if constraint.space == space {
	// 	log.Fatalln("Already added to this space")
	// }
	// if constraint.space != nil {
	// 	log.Fatalln("Already added to another space")
	// }
	// if space.locked != 0 {
	// 	log.Fatalln("Space is locked")
	// }

	a := constraint.bodyA
	b := constraint.bodyB

	// if a == nil && b == nil {
	// 	log.Fatalln("Constraint is attached to a null body")
	// }

	a.Activate()
	b.Activate()
	space.constraints = append(space.constraints, constraint)

	// Push onto the heads of the bodies' constraint lists
	constraint.nextA = a.constraintList
	// possible nil pointer dereference (SA5011)
	a.constraintList = constraint
	constraint.nextB = b.constraintList
	b.constraintList = constraint
	constraint.space = space

	return constraint
}

func (space *Space) RemoveConstraint(constraint *Constraint) {

	// if !space.ContainsConstraint(constraint) {
	// 	log.Fatalln("Constraint not found")
	// }

	// if space.locked != 0 {
	// 	log.Fatalln("Space is locked")
	// }

	constraint.bodyA.Activate()
	constraint.bodyB.Activate()
	for i, c := range space.constraints {
		if c == constraint {
			space.constraints = append(space.constraints[:i], space.constraints[i+1:]...)
			break
		}
	}

	constraint.bodyA.RemoveConstraint(constraint)
	constraint.bodyB.RemoveConstraint(constraint)
	constraint.space = nil
}

func (space *Space) FilterArbiters(body *Body, filter *Shape) {
	space.Lock()

	space.cachedArbiters.Filter(func(arb *Arbiter) bool {
		return CachedArbitersFilter(arb, space, filter, body)
	})

	space.Unlock(true)
}

func (space *Space) ContainsConstraint(constraint *Constraint) bool {
	return constraint.space == space
}

func (space *Space) ContainsShape(shape *Shape) bool {
	return shape.space == space
}

func (space *Space) ContainsBody(body *Body) bool {
	return body.space == space
}

func (space *Space) PushFreshContactBuffer() {
	stamp := space.stamp
	head := space.contactBuffersHead

	if head == nil {
		space.contactBuffersHead = NewContactBuffer(stamp, nil)
	} else if stamp-head.next.stamp > space.CollisionPersistence {
		tail := head.next
		space.contactBuffersHead = tail.InitHeader(stamp, tail)
	} else {
		// Allocate a new buffer and push it into the ring
		buffer := NewContactBuffer(stamp, head)
		head.next = buffer
		space.contactBuffersHead = buffer
	}
}

func (space *Space) ContactBufferGetArray() []Contact {
	if space.contactBuffersHead.numContacts+MAX_CONTACTS_PER_ARBITER > CONTACTS_BUFFER_SIZE {
		space.PushFreshContactBuffer()
	}

	head := space.contactBuffersHead
	return head.contacts[head.numContacts : head.numContacts+MAX_CONTACTS_PER_ARBITER]
}

func (space *Space) ProcessComponents(dt float64) {
	sleep := space.SleepTimeThreshold != Infinity

	// calculate the kinetic energy of all the bodies
	if sleep {
		dv := space.IdleSpeedThreshold
		var dvsq float64
		if dv != 0 {
			dvsq = dv * dv
		} else {
			dvsq = space.Gravity.LengthSq() * dt * dt
		}

		// update idling and reset component nodes
		for _, body := range space.DynamicBodies {
			if body.GetType() != BODY_DYNAMIC {
				continue
			}

			// Need to deal with infinite mass objects
			var keThreshold float64
			if dvsq != 0 {
				keThreshold = body.mass * dvsq
			}
			if body.KineticEnergy() > keThreshold {
				body.sleepingIdleTime = 0
			} else {
				body.sleepingIdleTime += dt
			}
		}
	}

	// Awaken any sleeping bodies found and then push arbiters to the bodies' lists.
	for _, arb := range space.Arbiters {
		a := arb.body_a
		b := arb.body_b

		if sleep {
			if b.GetType() == BODY_KINEMATIC || a.IsSleeping() {
				a.Activate()
			}
			if a.GetType() == BODY_KINEMATIC || b.IsSleeping() {
				b.Activate()
			}
		}

		a.PushArbiter(arb)
		b.PushArbiter(arb)
	}

	if sleep {
		// Bodies should be held active if connected by a joint to a kinematic.
		for _, constraint := range space.constraints {
			if constraint.bodyB.GetType() == BODY_KINEMATIC {
				constraint.bodyA.Activate()
			}
			if constraint.bodyA.GetType() == BODY_KINEMATIC {
				constraint.bodyB.Activate()
			}
		}

		// Generate components and deactivate sleeping ones
		for i := 0; i < len(space.DynamicBodies); {
			body := space.DynamicBodies[i]

			if body.ComponentRoot() == nil {
				// Body not in a component yet. Perform a DFS to flood fill mark
				// the component in the contact graph using this body as the root.
				FloodFillComponent(body, body)

				// Check if the component should be put to sleep.
				if !ComponentActive(body, space.SleepTimeThreshold) {
					space.sleepingComponents = append(space.sleepingComponents, body)
					for item := body; item != nil; item = item.sleepingNext {
						space.Deactivate(item)
					}

					// Deactivate() removed the current body from the list.
					// Skip incrementing the index counter.
					continue
				}
			}

			i++

			// Only sleeping bodies retain their component node pointers.
			body.sleepingRoot = nil
			body.sleepingNext = nil
		}
	}
}

func (space *Space) Step(dt float64) {
	if dt == 0 {
		return
	}

	space.stamp++

	prev_dt := space.curr_dt
	space.curr_dt = dt

	// reset and empty the arbiter lists
	for _, arb := range space.Arbiters {
		arb.state = ArbiterStateNormal

		// If both bodies are awake, unthread the arbiter from the contact graph.
		if !arb.body_a.IsSleeping() && !arb.body_b.IsSleeping() {
			arb.Unthread()
		}
	}
	space.Arbiters = space.Arbiters[:0]

	space.Lock()
	{
		// Integrate positions
		for _, body := range space.DynamicBodies {
			body.position_func(body, dt)
		}

		// Find colliding pairs.
		space.PushFreshContactBuffer()
		space.dynamicShapes.class.Each(ShapeUpdateFunc)
		space.dynamicShapes.class.ReindexQuery(SpaceCollideShapesFunc, space)
	}
	space.Unlock(false)

	// Rebuild the contact graph (and detect sleeping components if sleeping is enabled)
	space.ProcessComponents(dt)

	space.Lock()
	{
		// Clear out old cached arbiters and call separate callbacks
		space.cachedArbiters.Filter(func(arb *Arbiter) bool {
			return SpaceArbiterSetFilter(arb, space)
		})

		// Prestep the arbiters and constraints.
		slop := space.CollisionSlop
		biasCoef := 1 - math.Pow(space.CollisionBias, dt)
		for _, arbiter := range space.Arbiters {
			arbiter.PreStep(dt, slop, biasCoef)
		}

		for _, constraint := range space.constraints {
			if constraint.PreSolve != nil {
				constraint.PreSolve(constraint, space)
			}

			constraint.Class.PreStep(dt)
		}

		// Integrate velocities.
		damping := math.Pow(space.Damping, dt)
		gravity := space.Gravity
		for _, body := range space.DynamicBodies {
			body.velocity_func(body, gravity, damping, dt)
		}

		// Apply cached impulses
		var dt_coef float64
		if prev_dt != 0 {
			dt_coef = dt / prev_dt
		}

		for _, arbiter := range space.Arbiters {
			arbiter.ApplyCachedImpulse(dt_coef)
		}

		for _, constraint := range space.constraints {
			constraint.Class.ApplyCachedImpulse(dt_coef)
		}

		// Run the impulse solver.
		var i uint
		for i = 0; i < space.Iterations; i++ {
			for _, arbiter := range space.Arbiters {
				arbiter.ApplyImpulse()
			}

			for _, constraint := range space.constraints {
				constraint.Class.ApplyImpulse(dt)
			}
		}

		// Run the constraint post-solve callbacks
		for _, constraint := range space.constraints {
			if constraint.PostSolve != nil {
				constraint.PostSolve(constraint, space)
			}
		}

		// run the post-solve callbacks
		for _, arb := range space.Arbiters {
			arb.handler.PostSolveFunc(arb, space, arb.handler)
		}
	}
	space.Unlock(true)
}

func (space *Space) Lock() {
	space.locked++
}

// IsLocked returns true from inside a callback when objects cannot be added/removed.
func (space *Space) IsLocked() bool {
	return space.locked > 0
}

func (space *Space) Unlock(runPostStep bool) {
	space.locked--

	// if space.locked < 0 {
	// 	log.Fatalln("Space lock underflow")
	// }

	if space.locked != 0 {
		return
	}

	for i := 0; i < len(space.rousedBodies); i++ {
		space.Activate(space.rousedBodies[i])
		space.rousedBodies[i] = nil
	}
	space.rousedBodies = space.rousedBodies[:0]

	if runPostStep && !space.skipPostStep {
		space.skipPostStep = true

		for _, callback := range space.PostStepCallbacks {
			f := callback.callback

			// Mark the func as NULL in case calling it calls SpaceRunPostStepCallbacks() again.
			// TODO: need more tests around this case I think.
			callback.callback = nil

			if f != nil {
				f(space, callback.key, callback.data)
			}
		}

		space.PostStepCallbacks = space.PostStepCallbacks[:0]
		space.skipPostStep = false
	}
}

func (space *Space) UncacheArbiter(arb *Arbiter) {
	a := arb.a
	b := arb.b
	shapePair := ShapePair{a, b}
	arbHashId := HashPair(HashValue(unsafe.Pointer(a)), HashValue(unsafe.Pointer(b)))
	space.cachedArbiters.Remove(arbHashId, shapePair)
	for i, a := range space.Arbiters {
		if a == arb {
			// leak-free delete from slice
			last := len(space.Arbiters) - 1
			space.Arbiters[i] = space.Arbiters[last]
			space.Arbiters[last] = nil
			space.Arbiters = space.Arbiters[:last]
			return
		}
	}
	panic("Arbiter not found")
}

func (space *Space) PushContacts(count int) {
	// if count > MAX_CONTACTS_PER_ARBITER {
	// 	log.Fatalln("Contact buffer overflow")
	// }
	space.contactBuffersHead.numContacts += count
}

func (space *Space) PopContacts(count int) {
	space.contactBuffersHead.numContacts -= count
}

func (space *Space) LookupHandler(a, b CollisionType, defaultHandler *CollisionHandler) *CollisionHandler {
	types := &CollisionHandler{TypeA: a, TypeB: b}
	handler := space.collisionHandlers.Find(HashPair(HashValue(a), HashValue(b)), types)
	if handler != nil {
		return handler
	}
	return defaultHandler
}

// NewCollisionHandler sets a collision handler to handle specific collision types.
//
// The methods are called only when shapes with the specified CollisionTypeA and CollisionTypeB collide.
//
// Use Shape.SetCollisionType() to set type.
func (space *Space) NewCollisionHandler(collisionTypeA, collisionTypeB CollisionType) *CollisionHandler {
	hash := HashPair(HashValue(collisionTypeA), HashValue(collisionTypeB))
	handler := &CollisionHandler{collisionTypeA, collisionTypeB, DefaultBegin, DefaultPreSolve, DefaultPostSolve, DefaultSeparate, nil}
	return space.collisionHandlers.Insert(hash, handler, func(a *CollisionHandler) *CollisionHandler { return a })
}

func (space *Space) NewWildcardCollisionHandler(collisionType CollisionType) *CollisionHandler {
	space.UseWildcardDefaultHandler()

	hash := HashPair(HashValue(collisionType), HashValue(WILDCARD_COLLISION_TYPE))
	handler := &CollisionHandler{collisionType, WILDCARD_COLLISION_TYPE, AlwaysCollide, AlwaysCollide, DoNothing, DoNothing, nil}
	return space.collisionHandlers.Insert(hash, handler, func(a *CollisionHandler) *CollisionHandler { return a })
}

func (space *Space) UseWildcardDefaultHandler() {
	if !space.usesWildcards {
		space.usesWildcards = true
		space.defaultHandler = &CollisionHandlerDefault
	}
}

func (space *Space) UseSpatialHash(dim float64, count int) {
	staticShapes := NewSpaceHash(dim, count, ShapeGetBB, nil)
	dynamicShapes := NewSpaceHash(dim, count, ShapeGetBB, staticShapes)

	space.staticShapes.class.Each(func(shape *Shape) {
		staticShapes.class.Insert(shape, shape.hashid)
	})
	space.dynamicShapes.class.Each(func(shape *Shape) {
		dynamicShapes.class.Insert(shape, shape.hashid)
	})

	space.staticShapes = staticShapes
	space.dynamicShapes = dynamicShapes
}

// EachBody calls func f for each body in the space
//
// Example:
//
//	space.EachBody(func(body *cm.Body) {
//		fmt.Println(body.Position())
//	})
func (space *Space) EachBody(f func(b *Body)) {
	space.Lock()
	defer space.Unlock(true)

	for _, b := range space.DynamicBodies {
		f(b)
	}

	for _, b := range space.StaticBodies {
		f(b)
	}

	for _, root := range space.sleepingComponents {
		b := root

		for b != nil {
			next := b.sleepingNext
			f(b)
			b = next
		}
	}
}

// EachStaticBody calls func f for each static body in the space
func (space *Space) EachStaticBody(f func(b *Body)) {
	space.Lock()
	defer space.Unlock(true)

	for _, b := range space.StaticBodies {
		f(b)
	}

}

// EachDynamicBody calls func f for each dynamic body in the space
func (space *Space) EachDynamicBody(f func(b *Body)) {
	space.Lock()
	defer space.Unlock(true)

	for _, b := range space.DynamicBodies {
		f(b)
	}

	for _, root := range space.sleepingComponents {
		b := root

		for b != nil {
			next := b.sleepingNext
			f(b)
			b = next
		}
	}
}

// EachStaticShape calls func f for each static shape in the space
func (space *Space) EachStaticShape(f func(*Shape)) {
	space.Lock()
	space.staticShapes.class.Each(func(shape *Shape) {
		f(shape)
	})
	space.Unlock(true)
}

// EachDynamicShape calls func f for each dynamic shape in the space
func (space *Space) EachDynamicShape(f func(*Shape)) {
	space.Lock()
	space.dynamicShapes.class.Each(func(shape *Shape) {
		f(shape)
	})
	space.Unlock(true)
}

// EachShape calls func f for each shape in the space
func (space *Space) EachShape(f func(*Shape)) {
	space.Lock()

	space.dynamicShapes.class.Each(func(shape *Shape) {
		f(shape)
	})
	space.staticShapes.class.Each(func(shape *Shape) {
		f(shape)
	})

	space.Unlock(true)
}

func (space *Space) EachConstraint(f func(*Constraint)) {
	space.Lock()

	for i := 0; i < len(space.constraints); i++ {
		f(space.constraints[i])
	}

	space.Unlock(true)
}

// Query the space at a point and return the nearest shape found. Returns NULL if no shapes were found.
func (space *Space) PointQueryNearest(point vec.Vec2, maxDistance float64, filter ShapeFilter) *PointQueryInfo {
	info := &PointQueryInfo{nil, vec.Vec2{}, maxDistance, vec.Vec2{}}
	context := &PointQueryContext{point, maxDistance, filter, nil}

	bb := NewBBForCircle(point, math.Max(maxDistance, 0))
	space.dynamicShapes.class.Query(context, bb, NearestPointQueryNearest, info)
	space.staticShapes.class.Query(context, bb, NearestPointQueryNearest, info)

	return info
}

func (space *Space) BBQuery(bb BB, filter ShapeFilter, f SpaceBBQueryFunc, data interface{}) {
	context := BBQueryContext{bb, filter, f}
	space.staticShapes.class.Query(&context, bb, space.bbQuery, data)

	space.Lock()
	space.dynamicShapes.class.Query(&context, bb, space.bbQuery, data)
	space.Unlock(true)
}

func (space *Space) ArrayForBodyType(bodyType int) *[]*Body {
	if bodyType == BODY_STATIC {
		return &space.StaticBodies
	}
	return &space.DynamicBodies
}

// SegmentQuery Perform a directed line segment query (like a raycast) against the space and yield each shape intersected.
// The filter is applied to the query and follows the same rules as the collision detection. Sensor shapes are included
func (space *Space) SegmentQuery(start, end vec.Vec2, radius float64, filter ShapeFilter, f SpaceSegmentQueryFunc, data interface{}) {
	context := SegmentQueryContext{start, end, radius, filter, f}
	space.Lock()

	space.staticShapes.class.SegmentQuery(&context, start, end, 1, segmentQuery, data)
	space.dynamicShapes.class.SegmentQuery(&context, start, end, 1, segmentQuery, data)

	space.Unlock(true)
}

// SegmentQueryFirst Perform a directed line segment query (like a raycast) against the space and return the first shape hit.
// Returns nil if no shapes were hit.
func (space *Space) SegmentQueryFirst(start, end vec.Vec2, radius float64, filter ShapeFilter) SegmentQueryInfo {
	info := SegmentQueryInfo{nil, end, vec.Vec2{}, 1}
	context := &SegmentQueryContext{start, end, radius, filter, nil}
	space.staticShapes.class.SegmentQuery(context, start, end, 1, queryFirst, &info)
	space.dynamicShapes.class.SegmentQuery(context, start, end, info.Alpha, queryFirst, &info)
	return info
}

func (space *Space) TimeStep() float64 {
	return space.curr_dt
}

func (space *Space) PostStepCallback(key interface{}) *PostStepCallback {
	for i := 0; i < len(space.PostStepCallbacks); i++ {
		callback := space.PostStepCallbacks[i]
		if callback != nil && callback.key == key {
			return callback
		}
	}
	return nil
}

// AddPostStepCallback defines a callback to be run just before Space.Step() finishes.
//
// The main reason you want to define post-step callbacks is to get around
// the restriction that you cannot call the add/remove methods from a collision handler callback.
// Post-step callbacks run right before the next (or current) call to Space.Step() returns when it is safe to add and remove objects.
// You can only schedule one post-step callback per key value, this prevents you from accidentally removing an object twice.
// Registering a second callback for the same key is a no-op.
//
// example:
// type PostStepCallbackFunc func(space *Space, key interface{}, data interface{})
func (space *Space) AddPostStepCallback(f PostStepCallbackFunc, key, data interface{}) bool {
	if key == nil || space.PostStepCallback(key) == nil {
		callback := &PostStepCallback{
			key:  key,
			data: data,
		}
		if f != nil {
			callback.callback = f
		} else {
			callback.callback = PostStepDoNothing
		}
		space.PostStepCallbacks = append(space.PostStepCallbacks, callback)
		return true
	}
	return false
}

// ShapeQuery queries a space for any shapes overlapping the this shape and call the callback for each shape found.
func (space *Space) ShapeQuery(shape *Shape, callback func(shape *Shape, points *ContactPointSet)) bool {
	body := shape.body
	var bb BB
	if body != nil {
		bb = shape.Update(body.transform)
	} else {
		bb = shape.bb
	}

	var anyCollision bool

	shapeQuery := func(obj interface{}, b *Shape, collisionId uint32, _ interface{}) uint32 {
		a := obj.(*Shape)
		if a.Filter.Reject(b.Filter) || a == b {
			return collisionId
		}

		contactPointSet := ShapesCollide(a, b)
		if contactPointSet.Count > 0 {
			if callback != nil {
				callback(b, &contactPointSet)
			}
			anyCollision = !(a.Sensor || b.Sensor)
		}

		return collisionId
	}

	space.Lock()
	{
		space.dynamicShapes.class.Query(shape, bb, shapeQuery, nil)
		space.staticShapes.class.Query(shape, bb, shapeQuery, nil)
	}
	space.Unlock(true)

	return anyCollision
}

func (space *Space) bbQuery(obj interface{}, shape *Shape, collisionId uint32, data interface{}) uint32 {
	context := obj.(*BBQueryContext)
	if !shape.Filter.Reject(context.filter) && shape.BB().Intersects(context.bb) {
		context.f(shape, data)
	}
	return collisionId
}

func PostStepDoNothing(space *Space, key, data interface{}) {}

func SpaceCollideShapesFunc(obj interface{}, b *Shape, collisionId uint32, vspace interface{}) uint32 {
	a := obj.(*Shape)
	space := vspace.(*Space)

	// Reject any of the simple cases
	if QueryReject(a, b) {
		return collisionId
	}

	// Narrow-phase collision detection.
	info := Collide(a, b, collisionId, space.ContactBufferGetArray())

	if info.count == 0 {
		// shapes are not colliding
		return info.collisionId
	}

	//  Push contacts
	space.PushContacts(info.count)

	// Get an arbiter from space->arbiterSet for the two shapes.
	// This is where the persistent contact magic comes from.
	shapePair := ShapePair{info.a, info.b}
	arbHashId := HashPair(HashValue(unsafe.Pointer(info.a)), HashValue(unsafe.Pointer(info.b)))
	arb := space.cachedArbiters.Insert(arbHashId, shapePair, func(shapes ShapePair) *Arbiter {
		arb := space.pooledArbiters.Get().(*Arbiter)
		arb.Init(shapes.a, shapes.b)
		return arb
	})
	arb.Update(&info, space)

	if arb.state == ArbiterStateFirstCollision && !arb.handler.BeginFunc(arb, space, arb.handler.UserData) {
		arb.Ignore()
	}

	// Ignore the arbiter if it has been flagged
	if arb.state != ArbiterStateIgnore &&
		// Call PreSolve
		arb.handler.PreSolveFunc(arb, space, arb.handler.UserData) &&
		// Check (again) in case the pre-solve() callback called ArbiterIgnored().
		arb.state != ArbiterStateIgnore &&
		// Process, but don't add collisions for sensors.
		!(a.Sensor || b.Sensor) &&
		// Don't process collisions between two infinite mass bodies.
		// This includes collisions between two kinematic bodies, or a kinematic body and a static body.
		!(a.body.mass == Infinity && b.body.mass == Infinity) {
		space.Arbiters = append(space.Arbiters, arb)
	} else {
		space.PopContacts(info.count)
		arb.contacts = nil
		arb.count = 0

		// Normally arbiters are set as used after calling the post-solve callback.
		// However, post-solve() callbacks are not called for sensors or arbiters rejected from pre-solve.
		if arb.state != ArbiterStateIgnore {
			arb.state = ArbiterStateNormal
		}
	}

	// Time stamp the arbiter so we know it was used recently.
	arb.stamp = space.stamp
	return info.collisionId
}

func QueryReject(a, b *Shape) bool {
	if a.body == b.body {
		return true
	}
	if a.Filter.Reject(b.Filter) {
		return true
	}
	if !a.bb.Intersects(b.bb) {
		return true
	}
	if QueryRejectConstraints(a.body, b.body) {
		return true
	}
	return false
}

func QueryRejectConstraints(a, b *Body) bool {
	for constraint := a.constraintList; constraint != nil; constraint = constraint.Next(a) {
		if !constraint.collideBodies && ((constraint.bodyA == a && constraint.bodyB == b) ||
			(constraint.bodyA == b && constraint.bodyB == a)) {
			return true
		}
	}

	return false
}

func ComponentActive(root *Body, threshold float64) bool {
	for item := root; item != nil; item = item.sleepingNext {
		if item.sleepingIdleTime < threshold {
			return true
		}
	}
	return false
}

func FloodFillComponent(root *Body, body *Body) {
	// Kinematic bodies cannot be put to sleep and prevent bodies they are touching from sleeping.
	// Static bodies are effectively sleeping all the time.
	if body.GetType() != BODY_DYNAMIC {
		return
	}

	// body.sleeping.root
	other_root := body.ComponentRoot()
	if other_root == nil {
		root.ComponentAdd(body)

		for arb := body.arbiterList; arb != nil; arb = ArbiterNext(arb, body) {
			if body == arb.body_a {
				FloodFillComponent(root, arb.body_b)
			} else {
				FloodFillComponent(root, arb.body_a)
			}
		}

		for constraint := body.constraintList; constraint != nil; constraint = constraint.Next(body) {
			if body == constraint.bodyA {
				FloodFillComponent(root, constraint.bodyB)
			} else {
				FloodFillComponent(root, constraint.bodyA)
			}
		}
	} else {
		if other_root != root {
			log.Fatalln("Inconsistency detected in the contact graph (FFC)")
		}
	}
}

func ArbiterNext(arb *Arbiter, body *Body) *Arbiter {
	if arb.body_a == body {
		return arb.thread_a.next
	}
	return arb.thread_b.next
}

func Contains(bodies []*Body, body *Body) bool {
	for i := 0; i < len(bodies); i++ {
		if bodies[i] == body {
			return true
		}
	}
	return false
}

func NearestPointQueryNearest(obj interface{}, shape *Shape, collisionId uint32, out interface{}) uint32 {
	context := obj.(*PointQueryContext)
	if !shape.Filter.Reject(context.filter) && !shape.Sensor {
		info := shape.PointQuery(context.point)
		if info.Distance < out.(*PointQueryInfo).Distance {
			outp := out.(*PointQueryInfo)
			*outp = info
		}
	}

	return collisionId
}

func segmentQuery(obj interface{}, shape *Shape, data interface{}) float64 {
	context := obj.(*SegmentQueryContext)
	var info SegmentQueryInfo

	if !shape.Filter.Reject(context.filter) && shape.SegmentQuery(context.start, context.end, context.radius, &info) {
		context.f(shape, info.Point, info.Normal, info.Alpha, data)
	}

	return 1
}

func queryFirst(obj interface{}, shape *Shape, data interface{}) float64 {
	context := obj.(*SegmentQueryContext)
	out := data.(*SegmentQueryInfo)
	var info SegmentQueryInfo

	if !shape.Filter.Reject(context.filter) &&
		!shape.Sensor &&
		shape.SegmentQuery(context.start, context.end, context.radius, &info) &&
		info.Alpha < out.Alpha {
		*out = info
	}

	return out.Alpha
}

func arbiterSetEql(shapes ShapePair, arb *Arbiter) bool {
	a := shapes.a
	b := shapes.b

	return (a == arb.a && b == arb.b) || (b == arb.a && a == arb.b)
}

var ShapeVelocityFunc = func(obj interface{}) vec.Vec2 {
	return obj.(*Shape).body.vel
}

var ShapeUpdateFunc = func(shape *Shape) {
	shape.CacheBB()
}

type PostStepCallback struct {
	callback PostStepCallbackFunc
	key      interface{}
	data     interface{}
}

type PointQueryContext struct {
	point       vec.Vec2
	maxDistance float64
	filter      ShapeFilter
	f           SpacePointQueryFunc
}

type ShapePair struct {
	a, b *Shape
}

type BBQueryContext struct {
	bb     BB
	filter ShapeFilter
	f      SpaceBBQueryFunc
}

type SegmentQueryContext struct {
	start, end vec.Vec2
	radius     float64
	filter     ShapeFilter
	f          SpaceSegmentQueryFunc
}

type SpacePointQueryFunc func(*Shape, vec.Vec2, float64, vec.Vec2, interface{})
type SpaceBBQueryFunc func(shape *Shape, data interface{})
type SpaceSegmentQueryFunc func(shape *Shape, point, normal vec.Vec2, alpha float64, data interface{})
type PostStepCallbackFunc func(space *Space, key interface{}, data interface{})
