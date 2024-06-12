package cm

import (
	"kar/engine/vec"
	"math"
)

type SlideJoint struct {
	*Constraint

	AnchorA, AnchorB vec.Vec2
	Min, Max         float64

	r1, r2, n vec.Vec2
	nMass     float64

	jnAcc, bias float64
}

func NewSlideJoint(a, b *Body, anchorA, anchorB vec.Vec2, min, max float64) *Constraint {
	joint := &SlideJoint{
		AnchorA: anchorA,
		AnchorB: anchorB,
		Min:     min,
		Max:     max,
		jnAcc:   0,
	}
	joint.Constraint = NewConstraint(joint, a, b)
	return joint.Constraint
}

func (joint *SlideJoint) PreStep(dt float64) {
	a := joint.bodyA
	b := joint.bodyB

	joint.r1 = a.transform.Vect(joint.AnchorA.Sub(a.cog))
	joint.r2 = b.transform.Vect(joint.AnchorB.Sub(b.cog))

	delta := b.position.Add(joint.r2).Sub(a.position.Add(joint.r1))
	dist := delta.Length()
	pdist := 0.0
	if dist > joint.Max {
		pdist = dist - joint.Max
		joint.n = delta.Normalize()
	} else if dist < joint.Min {
		pdist = joint.Min - dist
		joint.n = delta.Normalize().Neg()
	} else {
		joint.n = vec.Vec2{}
		joint.jnAcc = 0
	}

	// calculate the mass normal
	joint.nMass = 1.0 / k_scalar(a, b, joint.r1, joint.r2, joint.n)

	// calculate bias velocity
	maxBias := joint.maxBias
	joint.bias = clamp(-bias_coef(joint.errorBias, dt)*pdist/dt, -maxBias, maxBias)
}

func (joint *SlideJoint) ApplyCachedImpulse(dt_coef float64) {
	a := joint.bodyA
	b := joint.bodyB

	j := joint.n.Scale(joint.jnAcc * dt_coef)
	apply_impulses(a, b, joint.r1, joint.r2, j)
}

func (joint *SlideJoint) ApplyImpulse(dt float64) {
	if joint.n.Equal(vec.Vec2{}) {
		return
	}

	a := joint.bodyA
	b := joint.bodyB
	n := joint.n
	r1 := joint.r1
	r2 := joint.r2

	vr := relative_velocity(a, b, r1, r2)
	vrn := vr.Dot(n)

	jn := (joint.bias - vrn) * joint.nMass
	jnOld := joint.jnAcc
	joint.jnAcc = clamp(jnOld+jn, -joint.maxForce*dt, 0)
	jn = joint.jnAcc - jnOld

	apply_impulses(a, b, joint.r1, joint.r2, n.Scale(jn))
}

func (joint *SlideJoint) GetImpulse() float64 {
	return math.Abs(joint.jnAcc)
}
