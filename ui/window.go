package ui

import (
	"image"
	"runtime"

	"code.google.com/p/portaudio-go/portaudio"

	"github.com/fogleman/nes/nes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	width   = 256
	height  = 240
	scale   = 3
	padding = 0
	title   = "NES"
)

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
	gl.BindTexture(gl.TEXTURE_2D, texture)
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
	gl.Vertex3f(-x, -y, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x, -y, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x, y, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-x, y, 1)
	gl.End()
}

func readKey(window *glfw.Window, key glfw.Key) bool {
	return window.GetKey(key) == glfw.Press
}

func readKeys(window *glfw.Window) [8]bool {
	var result [8]bool
	result[nes.ButtonA] = readKey(window, glfw.KeyZ)
	result[nes.ButtonB] = readKey(window, glfw.KeyX)
	result[nes.ButtonSelect] = readKey(window, glfw.KeyRightShift)
	result[nes.ButtonStart] = readKey(window, glfw.KeyEnter)
	result[nes.ButtonUp] = readKey(window, glfw.KeyUp)
	result[nes.ButtonDown] = readKey(window, glfw.KeyDown)
	result[nes.ButtonLeft] = readKey(window, glfw.KeyLeft)
	result[nes.ButtonRight] = readKey(window, glfw.KeyRight)
	return result
}

func readJoystick(joy glfw.Joystick) [8]bool {
	var result [8]bool
	if !glfw.JoystickPresent(joy) {
		return result
	}
	axes := glfw.GetJoystickAxes(joy)
	buttons := glfw.GetJoystickButtons(joy)
	result[nes.ButtonA] = buttons[0] == 1
	result[nes.ButtonB] = buttons[1] == 1
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
	k1 := readKeys(window)
	j1 := readJoystick(glfw.Joystick1)
	j2 := readJoystick(glfw.Joystick2)
	console.SetButtons1(combineButtons(k1, j1))
	console.SetButtons2(j2)
}

func Run(console *nes.Console) {
	// we need a parallel OS thread to avoid audio stuttering
	runtime.GOMAXPROCS(2)

	// we need to keep OpenGL calls on a single thread
	runtime.LockOSThread()

	portaudio.Initialize()
	defer portaudio.Terminate()

	audio := NewAudio()
	if err := audio.Start(); err != nil {
		panic(err)
	}
	defer audio.Stop()

	console.SetAudioChannel(audio.channel)

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
		updateControllers(window, console)
		console.StepFrame()
		setTexture(texture, console.Buffer())
		gl.Clear(gl.COLOR_BUFFER_BIT)
		drawQuad(window)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
