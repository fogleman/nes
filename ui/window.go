package ui

import (
	"image"
	"runtime"

	"github.com/fogleman/nes/nes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	width  = 256
	height = 240
	scale  = 3
	title  = "NES"
)

func init() {
	runtime.LockOSThread()
}

func createTexture() uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	return texture
}

func setTexture(texture uint32, im *image.RGBA) {
	size := im.Rect.Size()
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		int32(size.X), int32(size.Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(im.Pix))
}

func drawQuad(window *glfw.Window) {
	w, h := window.GetFramebufferSize()
	aspect := float32(w) / float32(h)
	var x, y, size float32
	size = 0.95
	if aspect >= 1 {
		x = size / aspect
		y = size
	} else {
		x = size
		y = size * aspect
	}
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-x, -y, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x, -y, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x, y, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-x, y, 1)
	gl.End()
}

func readKeys(window *glfw.Window, machine *nes.NES) {
	machine.SetPressed(1, nes.ButtonA, window.GetKey(glfw.KeyZ) == glfw.Press)
	machine.SetPressed(1, nes.ButtonB, window.GetKey(glfw.KeyX) == glfw.Press)
	machine.SetPressed(1, nes.ButtonSelect, window.GetKey(glfw.KeyRightShift) == glfw.Press)
	machine.SetPressed(1, nes.ButtonStart, window.GetKey(glfw.KeyEnter) == glfw.Press)
	machine.SetPressed(1, nes.ButtonUp, window.GetKey(glfw.KeyUp) == glfw.Press)
	machine.SetPressed(1, nes.ButtonDown, window.GetKey(glfw.KeyDown) == glfw.Press)
	machine.SetPressed(1, nes.ButtonLeft, window.GetKey(glfw.KeyLeft) == glfw.Press)
	machine.SetPressed(1, nes.ButtonRight, window.GetKey(glfw.KeyRight) == glfw.Press)
}

func Run(machine *nes.NES) {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(width*scale, height*scale, title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.TEXTURE_2D)
	texture := createTexture()

	for !window.ShouldClose() {
		// step emulator
		readKeys(window, machine)
		machine.StepFrame()
		setTexture(texture, machine.Buffer())
		// render frame
		gl.Clear(gl.COLOR_BUFFER_BIT)
		drawQuad(window)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
