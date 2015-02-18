package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/LaurenceGA/lib/colours"
	"github.com/LaurenceGA/lib/drawUtil"
	"github.com/LaurenceGA/lib/vector"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/gltext"
)

var (
	//Main
	windowWidth  = 640
	windowHeight = 480
	title        = "Pong"
	prevTime     time.Time
	deltaTime    float64
	padding      = 20.0
	rscore       = 0
	lscore       = 0
	//FPS
	maxFrames = time.Duration(60)
	accumTime float64
	frames    int
	fps       = 0.0
	//Paddle
	paddleWidth  = float64(15)
	paddleHeight = float64(80)
	movekeys     [2]bool
	moveSpeed    = float64(600)
	//Misc
	objects []object
	//Window
	cursorPos vector.Vector2
	screenDim vector.Vector2
	font      *gltext.Font
	smlFnt    *gltext.Font
	ballVerts [][2]float64
)

//Objects and what methods it has
type object interface {
	draw()
	step()
	xPos() float64
	yPos() float64
	getCol() bCollider
}

//Basic object properties
type obj struct {
	col colours.Colour
	pos vector.Vector2
}

//Collider fields
type bCollider struct {
	x, y, width, height float64
	tag                 string
}

//Defines ball's feilds
type ball struct {
	obj
	bCollider
	radius   float64
	velocity vector.Vector2
}

//X getter
func (b *ball) xPos() float64 {
	return b.pos.X()
}

//Y getter
func (b *ball) yPos() float64 {
	return b.pos.Y()
}

//Draw a pre defined list of verticies of a cricle
func (b *ball) draw() {
	//drawUtil.DrawSquare(b.radius, b.col)
	//drawUtil.DrawCircle(8, colours.White, 12)
	drawUtil.DrawVertexes(ballVerts, colours.White)
}

//Collider getter
func (b *ball) getCol() bCollider {
	return b.bCollider
}

//Ball logic: collision, movement
func (b *ball) step() {
	//Move
	b.pos = b.pos.Add(b.velocity.Mul(deltaTime))

	//Bounce of ceiling
	if (b.pos[1] - b.radius/2) <= 0 {
		b.velocity[1] *= -1
	} else if (b.pos[1] + b.radius/2) >= float64(windowHeight) {
		b.velocity[1] *= -1
	}

	b.bCollider.x = b.pos[0]
	b.bCollider.y = b.pos[1]

	for _, object := range objects {
		col := object.getCol()
		if col.tag == "paddle" {
			if ((b.pos[0]-b.radius/2) <= (col.x+col.width/2) &&
				(b.pos[0]-b.radius/2) >= (col.x-col.width/2)) ||
				((b.pos[0]+b.radius/2) <= (col.x+col.width/2) &&
					(b.pos[0]+b.radius/2 >= (col.x - col.width/2))) {
				if b.pos[1]-b.radius/2 <= col.y+col.height/2 &&
					b.pos[1]+b.radius/2 >= col.y-col.height/2 {
					b.velocity[0] *= -1
					b.velocity[1] = 400 * ((b.pos[1] - col.y) / (paddleHeight / 2))
				}
			}
		}
	}

	if b.pos[0] >= float64(windowWidth)+padding {
		lscore++
		//fmt.Printf("Left side:%v, Right side:%v\n", lscore, rscore)
		restartGame()
	} else if b.pos[0] <= -padding {
		rscore++
		//fmt.Printf("Left side:%v, Right side:%v\n", lscore, rscore)
		restartGame()
	}
}

//Define paddle fields
type paddle struct {
	obj
	bCollider
	width, height float64
	vspeed        float64
	id            int
}

//Draw the paddle rectangle
func (p *paddle) draw() {
	drawUtil.DrawRect(p.width, p.height, p.col)
}

//X getter
func (p *paddle) xPos() float64 {
	return p.pos.X()
}

//Y getter
func (p *paddle) yPos() float64 {
	return p.pos.Y()
}

//Collider getter
func (p *paddle) getCol() bCollider {
	return p.bCollider
}

//Paddle logic: movement, collider.
func (p *paddle) step() {
	p.bCollider.x = p.pos[0]
	p.bCollider.y = p.pos[1]

	if p.id == 0 {
		p.vspeed = 0
		if movekeys[0] {
			p.vspeed = moveSpeed
		}
		if movekeys[1] {
			p.vspeed = -moveSpeed
		}
	} else {
		for _, obj := range objects {
			if o, ok := obj.(*ball); ok {
				p.pos[1] = o.pos[1]
			}
		}
	}
	p.pos[1] += p.vspeed * deltaTime
	p.pos[1] = clamp(p.pos[1], paddleHeight/2, float64(windowHeight)-paddleHeight/2)
}

//Displays and glfw errors
func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

//Initialises callbacks, drawing perspective, fonts etc.
func initOpenGl(window *glfw.Window) {
	monitor, _ := glfw.GetPrimaryMonitor()
	vidMode, _ := monitor.GetVideoMode()
	screenDim[0] = float64(vidMode.Width)
	screenDim[1] = float64(vidMode.Height)
	sw := &screenDim[0]
	sh := &screenDim[1]
	w, h := window.GetSize() // query window to get screen pixels
	width, height := window.GetFramebufferSize()
	window.SetPosition(int(*sw/2)-(w/2), int(*sh/2)-(h/2))

	window.SetSizeCallback(onResize)
	window.SetKeyCallback(onKey)
	glfw.SetErrorCallback(errorCallback)

	gl.Viewport(0, 0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.ClearColor(0, 0, 0, 1)

	font = drawUtil.InitGltext(8)
	smlFnt = drawUtil.InitGltext(2)
	ballVerts = drawUtil.MakeCircle(8, 12)
}

//fixed drawing perspective and object locations when the
//Window is resized
func onResize(window *glfw.Window, w, h int) {
	if w < 1 {
		w = 1
	}

	if h < 1 {
		h = 1
	}

	gl.Viewport(0, 0, w, h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	windowWidth = w
	windowHeight = h

	for _, o := range objects {
		p, ok := o.(*paddle)
		if ok {
			if p.id == 1 {
				p.pos[0] = float64(w - 20)
			}
		}
	}
}

//Render things to the buffer and cycle through all objects'
//Draw events. Plus draw things such as the background and score
func render() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.LoadIdentity()

	//Background
	drawUtil.DrawLine(float64(windowWidth/2), 0, float64(windowWidth/2),
		float64(windowHeight), colours.White, 2.0)
	// drawUtil.DrawDotLine(float64(windowWidth/2), 0, float64(windowWidth/2),
	// float64(windowHeight), colours.White, 2.0, 16)

	//FPS
	ferr := drawUtil.DrawString(10, 10, strconv.FormatFloat(fps, 'f', 2, 64), colours.White, smlFnt)
	if ferr != nil {
		fmt.Println(ferr)
	}

	//Score
	// drawUtil.DrawString(0, 0, "TEST", colours.White, font)
	var w, _ = drawUtil.GetBounds(font)
	lstr := strconv.Itoa(lscore)
	rstr := strconv.Itoa(rscore)

	if lscore < 10 {
		lstr = "0" + lstr
	}
	if rscore < 10 {
		rstr = "0" + rstr
	}

	err := drawUtil.DrawString(float32(windowWidth/2-5-len(lstr)*w), 10, lstr, colours.White, font)
	if err != nil {
		fmt.Println(err)
	}
	err1 := drawUtil.DrawString(float32(windowWidth/2+10), 10, rstr, colours.White, font)
	if err1 != nil {
		fmt.Println(err1)
	}

	for _, o := range objects {
		drawObject(o)
	}
}

//handles calling a given objects draw function
//while also making sure it is in the right place and drawn correctly
func drawObject(o object) {
	gl.PushMatrix()
	//position
	gl.Translatef(float32(o.xPos()), float32(o.yPos()), 0.0)
	o.draw()
	gl.PopMatrix()
}

//Create an object and store it's reference
func createObject(o object) {
	objects = append(objects, o)
}

//Creates room's objects
func startupObjects() {
	//Make ball
	createObject(&ball{obj: obj{colours.White, vector.Vector2{100, 100}},
		radius: 10, velocity: vector.Vector2{350, 300}})
	//Make paddles
	createObject(&paddle{id: 0, obj: obj{colours.White,
		vector.Vector2{20, float64(windowHeight) / 2}},
		width: paddleWidth, height: paddleHeight,
		bCollider: bCollider{20, float64(windowHeight) / 2, paddleWidth, paddleHeight, "paddle"}})
	createObject(&paddle{id: 1, obj: obj{colours.White,
		vector.Vector2{float64(windowWidth) - 20, float64(windowHeight) / 2}},
		width: paddleWidth, height: paddleHeight,
		bCollider: bCollider{float64(windowWidth) - 20, float64(windowHeight) / 2, paddleWidth, paddleHeight, "paddle"}})
}

//Runs the step event logic for every object
func runSteps() {
	//Processes
	for _, o := range objects {
		o.step()
	}
}

//Specifically finds whether there is intent to move
func getMoveInp(w *glfw.Window) (bool, bool) {
	//u, d bool
	var u, d bool = false, false
	if w.GetKey(glfw.KeyUp) == glfw.Press {
		u = true
	}

	if w.GetKey(glfw.KeyDown) == glfw.Press {
		d = true
	}

	return u, d
}

//Sets all necesary inputs not handles by the key callback
func getInp(w *glfw.Window) {
	movekeys[0], movekeys[1] = getMoveInp(w)
	cursorPos[0], cursorPos[1] = w.GetCursorPosition()
	cursorPos[1] = float64(windowHeight) - cursorPos[1]

	//fmt.Printf("Mousex:%v Mousey:%v\n", cursorPos[0], cursorPos[1])
}

//Registers key events
func onKey(w *glfw.Window, key glfw.Key, sancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
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

//Resets the court, not the whole game
func restartGame() {
	objects = objects[:0]
	startupObjects()
}

//Window close code
func onExit() {
	//Things to do when exiting
}

func main() {
	//Initialise glfw3
	if !glfw.Init() {
		panic("Can't initialise glfw")
	}

	//Ensure termination at end of main func
	defer glfw.Terminate()

	//Create window
	window, err := glfw.CreateWindow(windowWidth, windowHeight, title, nil, nil)
	//Panic if we can't do it
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if gl.Init() != 0 {
		panic("GLEW init failed")
	}
	gl.GetError()

	initOpenGl(window)
	defer font.Release()
	defer smlFnt.Release()

	//Vsync 0=off, 1=on
	glfw.SwapInterval(0)

	//runtime.LockOSThread()

	startupObjects() // Instantiate objects

	prevTime = time.Now()

	for !window.ShouldClose() {
		deltaTime = float64(time.Now().Sub(prevTime).Seconds())
		prevTime = time.Now()
		//fmt.Println(deltaTime)

		//Do things
		getInp(window) // Get input
		runSteps()     // Run the step function for every object
		// rt := time.Now()
		render() //Draw all things that need to be drawn
		//renderTime := time.Now().Sub(rt).Seconds()
		//fmt.Println(renderTime)

		window.SwapBuffers() // Display new buffer
		glfw.PollEvents()

		//Wait out the rest of one second / max frames
		time.Sleep((time.Second / maxFrames) - (time.Now().Sub(prevTime)))

		accumTime += deltaTime
		frames++
		fps = float64(frames) / accumTime
		//fps = 1 / deltaTime
	}

	if window.ShouldClose() {
		onExit()
	}
}
