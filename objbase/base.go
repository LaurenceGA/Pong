/*
Package objbase is meant to define basic objects types
*/
package objbase

import "github.com/LaurenceGA/lib/vector"

//Object defines an object
type Object interface {
	Move(float64)
	//New()
	Step(float64)
	Draw()
	GetBase() *Base
}

//Instances is a list of current objects
var Instances []Object

//Base defines the fields that all objects should have
type Base struct {
	Pos      vector.Vector2
	Velocity vector.Vector2
	Tag      string
}

//Move moves the object around
func (b *Base) Move(deltaTime float64) {
	b.Pos = b.Pos.Add(b.Velocity.Mul(deltaTime))
}

//New initialises a base type
func (b *Base) New(pos vector.Vector2) {
	b.Pos = pos
}

//Collider fields
type Collider struct {
	X, Y, Width, Height float64
}

//New initialises the collider
func (c *Collider) New(x, y, w, h float64) {
	c.X, c.Y, c.Width, c.Height = x, y, w, h
}
