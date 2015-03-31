package ui

import (
	"image"

	"github.com/fogleman/nes/nes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const padding = 0

type GameView struct {
	director *Director
	console  *nes.Console
	title    string
	texture  uint32
}

func NewGameView(director *Director, console *nes.Console, title string) View {
	texture := createTexture()
	return &GameView{director, console, title, texture}
}

func (view *GameView) Enter() {
	view.director.SetTitle(view.title)
	view.console.SetAudioChannel(view.director.audio.channel)
}

func (view *GameView) Exit() {
	view.console.SetAudioChannel(nil)
}

func (view *GameView) Update(t, dt float64) {
	window := view.director.window
	console := view.console
	if readKey(window, glfw.KeyEscape) {
		view.director.ShowMenu()
	}
	if readKey(window, glfw.KeyR) {
		console.Reset()
	}
	updateControllers(window, console)
	console.StepSeconds(dt)
	gl.BindTexture(gl.TEXTURE_2D, view.texture)
	setTexture(console.Buffer())
	drawQuad(view.director.window)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func createTexture() uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}

func setTexture(im *image.RGBA) {
	size := im.Rect.Size()
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		int32(size.X), int32(size.Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(im.Pix))
}

func drawQuad(window *glfw.Window) {
	w, h := window.GetFramebufferSize()
	s1 := float32(w) / float32(width)
	s2 := float32(h) / float32(height)
	f := float32(1 - padding)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}

func readKey(window *glfw.Window, key glfw.Key) bool {
	return window.GetKey(key) == glfw.Press
}

func readKeys(window *glfw.Window, turbo bool) [8]bool {
	var result [8]bool
	result[nes.ButtonA] = readKey(window, glfw.KeyZ) || (turbo && readKey(window, glfw.KeyA))
	result[nes.ButtonB] = readKey(window, glfw.KeyX) || (turbo && readKey(window, glfw.KeyS))
	result[nes.ButtonSelect] = readKey(window, glfw.KeyRightShift)
	result[nes.ButtonStart] = readKey(window, glfw.KeyEnter)
	result[nes.ButtonUp] = readKey(window, glfw.KeyUp)
	result[nes.ButtonDown] = readKey(window, glfw.KeyDown)
	result[nes.ButtonLeft] = readKey(window, glfw.KeyLeft)
	result[nes.ButtonRight] = readKey(window, glfw.KeyRight)
	return result
}

func readJoystick(joy glfw.Joystick, turbo bool) [8]bool {
	var result [8]bool
	if !glfw.JoystickPresent(joy) {
		return result
	}
	axes := glfw.GetJoystickAxes(joy)
	buttons := glfw.GetJoystickButtons(joy)
	result[nes.ButtonA] = buttons[0] == 1 || (turbo && buttons[2] == 1)
	result[nes.ButtonB] = buttons[1] == 1 || (turbo && buttons[3] == 1)
	result[nes.ButtonSelect] = buttons[6] == 1
	result[nes.ButtonStart] = buttons[7] == 1
	result[nes.ButtonUp] = axes[1] < -0.5
	result[nes.ButtonDown] = axes[1] > 0.5
	result[nes.ButtonLeft] = axes[0] < -0.5
	result[nes.ButtonRight] = axes[0] > 0.5
	return result
}

func combineButtons(a, b [8]bool) [8]bool {
	var result [8]bool
	for i := 0; i < 8; i++ {
		result[i] = a[i] || b[i]
	}
	return result
}

func updateControllers(window *glfw.Window, console *nes.Console) {
	turbo := console.PPU.Frame%6 < 3
	k1 := readKeys(window, turbo)
	j1 := readJoystick(glfw.Joystick1, turbo)
	j2 := readJoystick(glfw.Joystick2, turbo)
	console.SetButtons1(combineButtons(k1, j1))
	console.SetButtons2(j2)
}
