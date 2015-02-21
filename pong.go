package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/LaurenceGA/Pong/game"
	"github.com/LaurenceGA/Pong/gameinput"
	"github.com/LaurenceGA/Pong/objbase"
	"github.com/LaurenceGA/Pong/objbase/objects"
	"github.com/LaurenceGA/Pong/screen"

	"github.com/LaurenceGA/lib/colours"
	"github.com/LaurenceGA/lib/drawUtil"
	"github.com/LaurenceGA/lib/vector"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/gltext"
)

var (
	//Main
	prevTime  time.Time
	deltaTime float64
	padding   = 20.0
	// rscore       = 0
	// lscore       = 0
	//FPS
	maxFrames = time.Duration(60)
	accumTime float64
	frames    int
	fps       = 0.0
	moveSpeed = float64(400)
	won       bool
	winStr    string
	//Window
	cursorPos vector.Vector2
	screenDim vector.Vector2
	font      *gltext.Font
	smlFnt    *gltext.Font
	ballVerts [][2]float64
)

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
	screen.WindowDim[0], screen.WindowDim[1] = w, h
	screen.ScreenDim[0], screen.ScreenDim[1] = int(*sw), int(*sh)

	window.SetSizeCallback(onResize)
	window.SetKeyCallback(gameinput.OnKey)
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
	game.WindowWidth = w
	game.WindowHeight = h
	screen.WindowDim[0], screen.WindowDim[1] = w, h

	for _, o := range objbase.Instances {
		p, ok := o.(*objects.Paddle)
		if ok {
			if p.Tag == "computer" {
				p.Pos[0] = float64(w - 20)
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
	drawUtil.DrawLine(float64(game.WindowWidth/2), 0, float64(game.WindowWidth/2),
		float64(game.WindowHeight), colours.White, 2.0)
	// drawUtil.DrawDotLine(float64(windowWidth/2), 0, float64(windowWidth/2),
	// float64(windowHeight), colours.White, 2.0, 16)

	//FPS
	ferr := drawUtil.DrawString(10, 10, strconv.FormatFloat(fps, 'f', 2, 64), colours.White, smlFnt)
	if ferr != nil {
		fmt.Println(ferr)
	}

	var w, h = drawUtil.GetBounds(font)

	if won {
		err := drawUtil.DrawString(float32((game.WindowWidth/2)-w*len(winStr)/2),
			float32((game.WindowHeight/2)-h/2), winStr, colours.White, font)
		if err != nil {
			fmt.Println(err)
		}
	}

	//Score
	lstr := strconv.Itoa(game.Lscore)
	rstr := strconv.Itoa(game.Rscore)

	if game.Lscore < 10 {
		lstr = "0" + lstr
	}
	if game.Rscore < 10 {
		rstr = "0" + rstr
	}

	err := drawUtil.DrawString(float32(game.WindowWidth/2-5-len(lstr)*w), 10, lstr, colours.White, font)
	if err != nil {
		fmt.Println(err)
	}
	err1 := drawUtil.DrawString(float32(game.WindowWidth/2+10), 10, rstr, colours.White, font)
	if err1 != nil {
		fmt.Println(err1)
	}

	for _, o := range objbase.Instances {
		drawObject(o)
	}
}

//handles calling a given objects draw function
//while also making sure it is in the right place and drawn correctly
func drawObject(o objbase.Object) {
	gl.PushMatrix()
	//position
	gl.Translatef(float32(o.GetBase().Pos[0]), float32(o.GetBase().Pos[1]), 0.0)
	o.Draw()
	gl.PopMatrix()
}

//Runs the step event logic for every object
func runSteps() {
	//Processes
	for _, o := range objbase.Instances {
		o.Step(deltaTime)
	}
}

//Window close code
func onExit() {
	//Things to do when exiting
}

func win(winner string) {
	objbase.Instances = objbase.Instances[:0]
	winStr = winner + " wins!"
	won = true
}

//CreateObject creates and object and stores it's reference
func createObject(o objbase.Object) {
	objbase.Instances = append(objbase.Instances, o)
}

//RestartGame resets the court, not the whole game
func restartGame() {
	objbase.Instances = objbase.Instances[:0]
	startupObjects()
}

//StartupObjects creates room's objects
func startupObjects() {
	b := (&objects.Ball{}).New(vector.Vector2{50, float64(screen.WindowDim[1]) / 2}, 10,
		vector.Vector2{600, 60})
	b.BallVerts = ballVerts
	createObject(b)
	createObject((&objects.Paddle{}).New(vector.Vector2{20, float64(screen.WindowDim[1]) / 2},
		game.PaddleWidth, game.PaddleHeight, "player"))
	createObject((&objects.Paddle{}).New(vector.Vector2{float64(screen.WindowDim[0]) - 20, float64(screen.WindowDim[1]) / 2},
		game.PaddleWidth, game.PaddleHeight, "computer"))
}

func main() {
	//Initialise glfw3
	if !glfw.Init() {
		panic("Can't initialise glfw")
	}

	//Ensure termination at end of main func
	defer glfw.Terminate()

	//Create window
	window, err := glfw.CreateWindow(game.WindowWidth, game.WindowHeight, game.Title, nil, nil)
	//Panic if we can't do it
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	//OpenGl requires that it is executed in the main thread
	runtime.LockOSThread()

	if gl.Init() != 0 {
		panic("GLEW init failed")
	}
	gl.GetError()

	initOpenGl(window)
	defer font.Release()
	defer smlFnt.Release()

	//Vsync 0=off, 1=on
	glfw.SwapInterval(0)

	startupObjects() // Instantiate objects

	prevTime = time.Now()

	for !window.ShouldClose() {
		deltaTime = float64(time.Now().Sub(prevTime).Seconds())
		prevTime = time.Now()
		//fmt.Println(deltaTime)

		//Do things
		//getInp(window) // Get input
		gameinput.GetInp(window)
		runSteps() // Run the step function for every object
		// rt := time.Now()
		render() //Draw all things that need to be drawn
		//renderTime := time.Now().Sub(rt).Seconds()
		//fmt.Println(renderTime)

		if game.Lscore >= 10 {
			win("Player")
		} else if game.Rscore >= 10 {
			win("Computer")
		}

		if game.Restart {
			restartGame()
			game.Restart = false
		}

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
