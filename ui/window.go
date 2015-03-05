package ui

import (
	"image"
	"image/draw"
	"image/png"
	"os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	width  = 256
	height = 224
	scale  = 4
	title  = "NES"
)

func init() {
	runtime.LockOSThread()
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(file)
}

func convertImage(im image.Image) *image.RGBA {
	rgba := image.NewRGBA(im.Bounds())
	draw.Draw(rgba, rgba.Bounds(), im, image.Point{0, 0}, draw.Src)
	return rgba
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

func setTexture(texture uint32, im *image.RGBA, offset int) {
	size := im.Rect.Size()
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		int32(size.X), int32(size.Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(im.Pix[offset:]))
}

func Run() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
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

	im, err := loadImage("texture.png")
	if err != nil {
		panic(err)
	}
	rgba := convertImage(im)

	texture := createTexture()

	frame := 0
	for !window.ShouldClose() {
		setTexture(texture, rgba, (frame*4)%1024)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Begin(gl.QUADS)
		gl.TexCoord2f(0, 1)
		gl.Vertex3f(-1, -1, 1)
		gl.TexCoord2f(1, 1)
		gl.Vertex3f(1, -1, 1)
		gl.TexCoord2f(1, 0)
		gl.Vertex3f(1, 1, 1)
		gl.TexCoord2f(0, 0)
		gl.Vertex3f(-1, 1, 1)
		gl.End()
		window.SwapBuffers()
		glfw.PollEvents()
		frame++
	}
}
