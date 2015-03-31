package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	border = 10
	margin = 10
)

type MenuView struct {
	director     *Director
	paths        []string
	nx, ny, i, j int
	t            float64
}

func NewMenuView(director *Director, paths []string) View {
	return &MenuView{director, paths, 0, 0, 0, 0, 0}
}

func (view *MenuView) OnKey(
	window *glfw.Window, key glfw.Key, scancode int, action glfw.Action,
	mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyUp:
			view.j--
		case glfw.KeyDown:
			view.j++
		case glfw.KeyLeft:
			view.i--
		case glfw.KeyRight:
			view.i++
		case glfw.KeyEnter:
			view.OnSelect()
		}
	}
	view.t = glfw.GetTime()
}

func (view *MenuView) OnSelect() {
	index := view.nx*view.j + view.i
	view.director.PlayGame(view.paths[index])
}

func (view *MenuView) Enter() {
	view.director.SetTitle("Select Game")
	view.director.window.SetKeyCallback(view.OnKey)
}

func (view *MenuView) Exit() {
	view.director.window.SetKeyCallback(nil)
}

func (view *MenuView) Update(t, dt float64) {
	window := view.director.window
	w, h := window.GetFramebufferSize()
	sx := 256 + margin*2
	sy := 240 + margin*2
	nx := (w - border*2) / sx
	ny := (h - border*2) / sy
	ox := (w - nx*sx) / 2
	oy := (h - ny*sy) / 2
	view.clampSelection(nx, ny)
	gl.PushMatrix()
	gl.Ortho(0, float64(w), float64(h), 0, -1, 1)
	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			x := ox + i*sx
			y := oy + j*sy
			drawThumbnail(float32(x), float32(y))
			if int((t-view.t)*4)%2 == 0 && i == view.i && j == view.j {
				drawSelection(float32(x), float32(y), 8, 4)
			}
		}
	}
	gl.PopMatrix()
	view.nx = nx
	view.ny = ny
}

func (view *MenuView) clampSelection(nx, ny int) {
	if view.i < 0 {
		view.i = 0
	}
	if view.i >= nx {
		view.i = nx - 1
	}
	if view.j < 0 {
		view.j = 0
	}
	if view.j >= ny {
		view.j = ny - 1
	}
}

func drawThumbnail(x, y float32) {
	gl.Color3f(1, 1, 1)
	gl.Begin(gl.QUADS)
	gl.Vertex2f(x, y)
	gl.Vertex2f(x+256, y)
	gl.Vertex2f(x+256, y+240)
	gl.Vertex2f(x, y+240)
	gl.End()
}

func drawSelection(x, y, p, w float32) {
	gl.LineWidth(w)
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2f(x-p, y-p)
	gl.Vertex2f(x+256+p, y-p)
	gl.Vertex2f(x+256+p, y+240+p)
	gl.Vertex2f(x-p, y+240+p)
	gl.Vertex2f(x-p, y-p)
	gl.End()
}
