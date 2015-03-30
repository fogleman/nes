package ui

import (
	"log"
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
)

func init() {
	// we need a parallel OS thread to avoid audio stuttering
	runtime.GOMAXPROCS(2)

	// we need to keep OpenGL calls on a single thread
	runtime.LockOSThread()
}

func Run(path string) {
	console, err := nes.NewConsole(path)
	if err != nil {
		log.Fatalln(err)
	}

	portaudio.Initialize()
	defer portaudio.Terminate()

	if err := glfw.Init(); err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(width*scale, height*scale, path, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	gl.Enable(gl.TEXTURE_2D)

	timestamp := glfw.GetTime()

	audio := NewAudio()
	console.SetAudioChannel(audio.channel)
	if err := audio.Start(); err != nil {
		log.Fatalln(err)
	}
	defer audio.Stop()

	director := &Director{window}
	view := NewGameView(director, console)

	for !window.ShouldClose() {
		now := glfw.GetTime()
		elapsed := now - timestamp
		timestamp = now
		view.Update(timestamp, elapsed)
		view.Draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
