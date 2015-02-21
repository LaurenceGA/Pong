package objects

import (
	"github.com/LaurenceGA/lib/colours"
	"github.com/LaurenceGA/lib/drawUtil"

	"github.com/LaurenceGA/Pong/game"
	"github.com/LaurenceGA/Pong/objbase"
	"github.com/LaurenceGA/Pong/screen"

	"github.com/LaurenceGA/lib/vector"
)

//Ball is a pong ball
type Ball struct {
	objbase.Base
	objbase.Collider
	Radius    float64
	BallVerts [][2]float64
}

//New initialises ball
func (b *Ball) New(pos [2]float64, r float64, v vector.Vector2) *Ball {
	b.Base.New(pos)
	b.Radius = r
	b.Collider.New(pos[0], pos[1], r, r)
	b.Velocity = v
	return b
}

//GetBase returns a reference to it's base
func (b *Ball) GetBase() *objbase.Base {
	return &b.Base
}

//Draw creates a circle that is the ball
func (b *Ball) Draw() {
	drawUtil.DrawVertexes(b.BallVerts, colours.White)
}

//Step defines logic: collision, movement
func (b *Ball) Step(deltaTime float64) {
	//Move
	b.Move(deltaTime)

	b.Collider.X = b.Pos[0]
	b.Collider.Y = b.Pos[1]

	//Bounce of ceiling
	if (b.Pos[1] - b.Radius/2) <= 0 {
		b.Velocity[1] *= -1
	} else if (b.Pos[1] + b.Radius/2) >= float64(screen.WindowDim[1]) {
		b.Velocity[1] *= -1
	}

	for _, object := range objbase.Instances {
		if o, ok := object.(*Paddle); ok {
			col := o.Collider
			//Collision should really be better implemented
			if ((b.Pos[0]-b.Radius/2) <= (col.X+col.Width/2) &&
				(b.Pos[0]-b.Radius/2) >= (col.X-col.Width/2)) ||
				((b.Pos[0]+b.Radius/2) <= (col.X+col.Width/2) &&
					(b.Pos[0]+b.Radius/2 >= (col.X - col.Width/2))) {
				if b.Pos[1]-b.Radius/2 <= col.Y+col.Height/2 &&
					b.Pos[1]+b.Radius/2 >= col.Y-col.Height/2 {
					b.Velocity[0] *= -1
					b.Velocity[1] = 400 * ((b.Pos[1] - col.Y) / (game.PaddleHeight / 2))
				}
			}
		}
	}

	padding := 20.0

	if b.Pos[0] >= float64(screen.WindowDim[0])+padding {
		game.Lscore++
		//fmt.Printf("Left side:%v, Right side:%v\n", lscore, rscore)
		game.Restart = true
	} else if b.Pos[0] <= -padding {
		game.Rscore++
		game.Restart = true
	}
}
