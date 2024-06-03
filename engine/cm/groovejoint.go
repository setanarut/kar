package cm

type GrooveJoint struct {
	*Constraint

	GrooveN, GrooveA, GrooveB Vec2
	AnchorB                   Vec2

	grooveTn Vec2
	clamp    float64
	r1, r2   Vec2
	k        Mat2x2

	jAcc, bias Vec2
}

func NewGrooveJoint(a, b *Body, grooveA, grooveB, anchorB Vec2) *Constraint {
	joint := &GrooveJoint{
		GrooveA: grooveA,
		GrooveB: grooveB,
		GrooveN: grooveB.Sub(grooveA).Normalize().Perp(),
		AnchorB: anchorB,
		jAcc:    Vec2{},
	}
	joint.Constraint = NewConstraint(joint, a, b)
	return joint.Constraint
}

func (joint *GrooveJoint) PreStep(dt float64) {
	a := joint.bodyA
	b := joint.bodyB

	ta := a.transform.Point(joint.GrooveA)
	tb := a.transform.Point(joint.GrooveB)

	n := a.transform.Vect(joint.GrooveN)
	d := ta.Dot(n)

	joint.grooveTn = n
	joint.r2 = b.transform.Vect(joint.AnchorB.Sub(b.cog))

	td := b.position.Add(joint.r2).Cross(n)

	if td <= ta.Cross(n) {
		joint.clamp = 1
		joint.r1 = ta.Sub(a.position)
	} else if td >= tb.Cross(n) {
		joint.clamp = -1
		joint.r1 = tb.Sub(a.position)
	} else {
		joint.clamp = 0
		joint.r1 = n.Perp().Mult(-td).Add(n.Mult(d)).Sub(a.position)
	}

	joint.k = k_tensor(a, b, joint.r1, joint.r2)

	delta := b.position.Add(joint.r2).Sub(a.position.Add(joint.r1))
	joint.bias = delta.Mult(-bias_coef(joint.errorBias, dt) / dt).ClampLenght(joint.maxBias)
}

func (joint *GrooveJoint) ApplyCachedImpulse(dt_coef float64) {
	a := joint.bodyA
	b := joint.bodyB

	apply_impulses(a, b, joint.r1, joint.r2, joint.jAcc.Mult(dt_coef))
}

func (joint *GrooveJoint) grooveConstrain(j Vec2, dt float64) Vec2 {
	n := joint.grooveTn
	var jClamp Vec2
	if joint.clamp*j.Cross(n) > 0 {
		jClamp = j
	} else {
		jClamp = j.Project(n)
	}
	return jClamp.ClampLenght(joint.maxForce * dt)
}

func (joint *GrooveJoint) ApplyImpulse(dt float64) {
	a := joint.bodyA
	b := joint.bodyB

	r1 := joint.r1
	r2 := joint.r2

	vr := relative_velocity(a, b, r1, r2)

	j := joint.k.Transform(joint.bias.Sub(vr))
	jOld := joint.jAcc
	joint.jAcc = joint.grooveConstrain(jOld.Add(j), dt)
	j = joint.jAcc.Sub(jOld)

	apply_impulses(a, b, joint.r1, joint.r2, j)
}

func (joint *GrooveJoint) GetImpulse() float64 {
	return joint.jAcc.Length()
}
