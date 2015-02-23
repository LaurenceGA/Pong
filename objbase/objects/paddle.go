package objects

import (
	"math"
	"math/rand"

	"github.com/LaurenceGA/lib/colours"
	"github.com/LaurenceGA/lib/drawUtil"

	"github.com/LaurenceGA/Pong/gameinput"
	"github.com/LaurenceGA/Pong/objbase"
	"github.com/LaurenceGA/Pong/screen"
)

//Paddle is the object the ball is hit with
type Paddle struct {
	objbase.Base
	objbase.Collider
	Width, Height float64
	moveSpeed     float64
	targetOffset  float64
}

//GetBase returns a reference to it's base
func (p *Paddle) GetBase() *objbase.Base {
	return &p.Base
}

//GetNewOffset defines where the paddle will try to be
func (p *Paddle) GetNewOffset() {
	p.targetOffset = (rand.Float64() * p.Height) - (p.Height / 2)
	// fmt.Println(p.targetOffset)
}

//New initialises a paddle
func (p *Paddle) New(pos [2]float64, w, h float64, t string) *Paddle {
	p.Base.New(pos)
	p.Collider.New(pos[0], pos[1], w, h)
	p.Width, p.Height = w, h
	p.Tag = t
	p.moveSpeed = 400
	return p
}

//Draw will make a coloured paddle
func (p *Paddle) Draw() {
	drawUtil.DrawRect(p.Width, p.Height, colours.White)
}

//Step runs object logic
func (p *Paddle) Step(deltaTime float64) {
	p.Move(deltaTime)

	p.Collider.X = p.Pos[0]
	p.Collider.Y = p.Pos[1]

	movekeys := gameinput.Movekeys

	if p.Tag == "player" {
		p.Velocity[1] = 0
		if movekeys[0] {
			p.Velocity[1] = p.moveSpeed
		}
		if movekeys[1] {
			p.Velocity[1] = -p.moveSpeed
		}
	} else {
		for _, obj := range objbase.Instances {
			if o, ok := obj.(*Ball); ok {
				if math.Abs(p.Pos[1]-(o.Pos[1]+p.targetOffset)) > (p.moveSpeed/2)*deltaTime {
					if o.Pos[1]+p.targetOffset > p.Pos[1] {
						p.Velocity[1] = p.moveSpeed / 2
					} else if o.Pos[1]+p.targetOffset < p.Pos[1] {
						p.Velocity[1] = -p.moveSpeed / 2
					} else {
						p.Velocity[1] = 0
					}
				} else {
					p.Pos[1] = o.Pos[1] + p.targetOffset
				}
			}
		}
	}

	p.Pos[1] = clamp(p.Pos[1], p.Height/2, float64(screen.WindowDim[1])-p.Height/2)
}

//A simple mathematical clamp function
func clamp(p, p0, p1 float64) float64 {

	if p0 > p1 {
		panic("Improper use of clamp")
	} else if p < p0 {
		return p0
	} else if p > p1 {
		return p1
	}
	return p
}
