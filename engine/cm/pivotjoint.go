package cm

type PivotJoint struct {
	*Constraint
	AnchorA, AnchorB Vec2

	r1, r2 Vec2
	k      Mat2x2

	jAcc, bias Vec2
}

func NewPivotJoint(a, b *Body, pivot Vec2) *Constraint {
	var anchorA Vec2
	var anchorB Vec2

	if a != nil {
		anchorA = a.WorldToLocal(pivot)
	} else {
		anchorA = pivot
	}

	if b != nil {
		anchorB = b.WorldToLocal(pivot)
	} else {
		anchorB = pivot
	}

	return NewPivotJoint2(a, b, anchorA, anchorB)
}

func NewPivotJoint2(a, b *Body, anchorA, anchorB Vec2) *Constraint {
	joint := &PivotJoint{
		AnchorA: anchorA,
		AnchorB: anchorB,
		jAcc:    Vec2{},
	}
	constraint := NewConstraint(joint, a, b)
	joint.Constraint = constraint
	return constraint
}

func (joint *PivotJoint) PreStep(dt float64) {
	a := joint.Constraint.bodyA
	b := joint.Constraint.bodyB

	joint.r1 = a.transform.Vect(joint.AnchorA.Sub(a.cog))
	joint.r2 = b.transform.Vect(joint.AnchorB.Sub(b.cog))

	// Calculate mass tensor
	joint.k = k_tensor(a, b, joint.r1, joint.r2)

	// calculate bias velocity
	delta := b.position.Add(joint.r2).Sub(a.position.Add(joint.r1))
	joint.bias = delta.Mult(-bias_coef(joint.Constraint.errorBias, dt) / dt).ClampLenght(joint.Constraint.maxBias)
}

func (joint *PivotJoint) ApplyCachedImpulse(dt_coef float64) {
	apply_impulses(joint.bodyA, joint.bodyB, joint.r1, joint.r2, joint.jAcc.Mult(dt_coef))
}

func (joint *PivotJoint) ApplyImpulse(dt float64) {
	a := joint.Constraint.bodyA
	b := joint.Constraint.bodyB

	r1 := joint.r1
	r2 := joint.r2

	// compute relative velocity
	vr := relative_velocity(a, b, r1, r2)

	// compute normal impulse
	j := joint.k.Transform(joint.bias.Sub(vr))
	jOld := joint.jAcc
	joint.jAcc = joint.jAcc.Add(j).ClampLenght(joint.Constraint.maxForce * dt)
	j = joint.jAcc.Sub(jOld)

	apply_impulses(a, b, joint.r1, joint.r2, j)
}

func (joint *PivotJoint) GetImpulse() float64 {
	return joint.jAcc.Length()
}
