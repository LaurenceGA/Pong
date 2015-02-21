package gameinput

import glfw "github.com/go-gl/glfw3"

//Move keys are the ones that move the paddle
var Movekeys [2]bool

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

//GetInp sets all necesary inputs not handles by the key callback
func GetInp(w *glfw.Window) {
	Movekeys[0], Movekeys[1] = getMoveInp(w)
	//cursorPos[0], cursorPos[1] = w.GetCursorPosition()
	//cursorPos[1] = float64(windowHeight) - cursorPos[1]

	//fmt.Printf("Mousex:%v Mousey:%v\n", cursorPos[0], cursorPos[1])
}

//OnKey registers key events
func OnKey(w *glfw.Window, key glfw.Key, sancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
}
