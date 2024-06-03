package cm

import "math"

type DampedRotarySpring struct {
	*Constraint

	RestAngle, Stiffness, Damping float64
	SpringTorqueFunc              func(spring *DampedRotarySpring, relativeAngle float64) float64

	targetWrn, wCoef float64
	iSum, jAcc       float64
}

func defaultSpringTorque(spring *DampedRotarySpring, relativeAngle float64) float64 {
	return (relativeAngle - spring.RestAngle) * spring.Stiffness
}

func NewDampedRotarySpring(a, b *Body, restAngle, stiffness, damping float64) *Constraint {
	joint := &DampedRotarySpring{
		RestAngle:        restAngle,
		Stiffness:        stiffness,
		Damping:          damping,
		SpringTorqueFunc: defaultSpringTorque,
	}
	joint.Constraint = NewConstraint(joint, a, b)
	return joint.Constraint
}

func (spring *DampedRotarySpring) PreStep(dt float64) {
	a := spring.bodyA
	b := spring.bodyB

	moment := a.moi_inv + b.moi_inv
	spring.iSum = 1.0 / moment

	spring.wCoef = 1.0 - math.Exp(-spring.Damping*dt*moment)
	spring.targetWrn = 0

	jSpring := spring.SpringTorqueFunc(spring, a.angle-b.angle) * dt
	spring.jAcc = jSpring

	a.w -= jSpring * a.moi_inv
	b.w += jSpring * b.moi_inv
}

func (joint *DampedRotarySpring) ApplyCachedImpulse(dt_coef float64) {
	// nothing to do here
}

func (spring *DampedRotarySpring) ApplyImpulse(dt float64) {
	a := spring.bodyA
	b := spring.bodyB

	wrn := a.w - b.w

	wDamp := (spring.targetWrn - wrn) * spring.wCoef
	spring.targetWrn = wrn + wDamp

	jDamp := wDamp * spring.iSum
	spring.jAcc += jDamp

	a.w += jDamp * a.moi_inv
	b.w -= jDamp * b.moi_inv
}

func (joint *DampedRotarySpring) GetImpulse() float64 {
	return joint.jAcc
}
