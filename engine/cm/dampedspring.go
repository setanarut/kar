package cm

import (
	"math"
)

type DampedSpringForceFunc func(spring *DampedSpring, dist float64) float64

type DampedSpring struct {
	*Constraint

	AnchorA, AnchorB               Vec2
	RestLength, Stiffness, Damping float64
	SpringForceFunc                DampedSpringForceFunc

	targetVrn, vCoef float64

	r1, r2 Vec2
	nMass  float64
	n      Vec2

	jAcc float64
}

func NewDampedSpring(a, b *Body, anchorA, anchorB Vec2, restLength, stiffness, damping float64) *Constraint {
	spring := &DampedSpring{
		AnchorA:         anchorA,
		AnchorB:         anchorB,
		RestLength:      restLength,
		Stiffness:       stiffness,
		Damping:         damping,
		SpringForceFunc: DefaultSpringForce,
		jAcc:            0,
	}
	spring.Constraint = NewConstraint(spring, a, b)
	return spring.Constraint
}

func (spring *DampedSpring) PreStep(dt float64) {
	a := spring.bodyA
	b := spring.bodyB

	spring.r1 = a.transform.Vect(spring.AnchorA.Sub(a.cog))
	spring.r2 = b.transform.Vect(spring.AnchorB.Sub(b.cog))

	delta := b.position.Add(spring.r2).Sub(a.position.Add(spring.r1))
	dist := delta.Length()
	if dist != 0 {
		spring.n = delta.Scale(1.0 / dist)
	} else {
		spring.n = delta.Scale(1.0 / Infinity)
	}

	k := k_scalar(a, b, spring.r1, spring.r2, spring.n)

	// if k == 0 {
	// 	log.Fatalln("Unsolvable spring")
	// }

	spring.nMass = 1.0 / k

	spring.targetVrn = 0
	spring.vCoef = 1.0 - math.Exp(-spring.Damping*dt*k)

	fSpring := spring.SpringForceFunc(spring, dist)
	spring.jAcc = fSpring * dt
	apply_impulses(a, b, spring.r1, spring.r2, spring.n.Scale(spring.jAcc))
}

func (spring *DampedSpring) ApplyCachedImpulse(dt_coef float64) {
	// nothing to do here
}

func (spring *DampedSpring) ApplyImpulse(dt float64) {
	a := spring.bodyA
	b := spring.bodyB

	n := spring.n
	r1 := spring.r1
	r2 := spring.r2

	vrn := normal_relative_velocity(a, b, r1, r2, n)

	vDamp := (spring.targetVrn - vrn) * spring.vCoef
	spring.targetVrn = vrn + vDamp

	jDamp := vDamp * spring.nMass
	spring.jAcc += jDamp
	apply_impulses(a, b, spring.r1, spring.r2, spring.n.Scale(jDamp))
}

func (spring *DampedSpring) GetImpulse() float64 {
	return spring.jAcc
}

func DefaultSpringForce(spring *DampedSpring, dist float64) float64 {
	return (spring.RestLength - dist) * spring.Stiffness
}
