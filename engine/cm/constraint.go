package cm

import (
	"math"
)

type Constrainer interface {
	PreStep(dt float64)
	ApplyCachedImpulse(dt_coef float64)
	ApplyImpulse(dt float64)
	GetImpulse() float64
}

type ConstraintPreSolveFunc func(*Constraint, *Space)
type ConstraintPostSolveFunc func(*Constraint, *Space)

type Constraint struct {
	Class Constrainer
	space *Space

	bodyA, bodyB *Body
	nextA, nextB *Constraint

	maxForce, errorBias, maxBias float64

	collideBodies bool
	PreSolve      ConstraintPreSolveFunc
	PostSolve     ConstraintPostSolveFunc

	UserData interface{}
}

func NewConstraint(class Constrainer, a, b *Body) *Constraint {
	return &Constraint{
		Class: class,
		bodyA: a,
		bodyB: b,
		space: nil,

		maxForce:  Infinity,
		errorBias: math.Pow(1.0-0.1, 60.0),
		maxBias:   Infinity,

		collideBodies: true,
		PreSolve:      nil,
		PostSolve:     nil,
	}
}

func (c *Constraint) ActivateBodies() {
	c.bodyA.Activate()
	c.bodyB.Activate()
}
func (c *Constraint) BodyA() *Body {
	return c.bodyA
}
func (c *Constraint) BodyB() *Body {
	return c.bodyB
}

func (c Constraint) MaxForce() float64 {
	return c.maxForce
}

func (c *Constraint) SetMaxForce(max float64) {
	// if max < 0.0 {
	// 	log.Fatalln("Must be positive")
	// }
	c.ActivateBodies()
	c.maxForce = max
}

func (c Constraint) MaxBias() float64 {
	return c.maxBias
}

func (c *Constraint) SetMaxBias(max float64) {
	// if max < 0 {
	// 	log.Fatalln("Must be positive")
	// }
	c.ActivateBodies()
	c.maxBias = max
}

func (c Constraint) ErrorBias() float64 {
	return c.errorBias
}

func (c *Constraint) SetErrorBias(errorBias float64) {
	// if errorBias < 0 {
	// 	log.Fatalln("Must be positive")
	// }
	c.ActivateBodies()
	c.errorBias = errorBias
}

func (c *Constraint) Next(body *Body) *Constraint {
	if c.bodyA == body {
		return c.nextA
	} else {
		return c.nextB
	}
}

func (c *Constraint) SetCollideBodies(collideBodies bool) {
	c.ActivateBodies()
	c.collideBodies = collideBodies
}
