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
	texture      *Texture
	nx, ny, i, j int
	scroll       int
	t            float64
}

func NewMenuView(director *Director, paths []string) View {
	texture := NewTexture()
	return &MenuView{director, paths, texture, 0, 0, 0, 0, 0, 0}
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
		default:
			return
		}
		view.t = glfw.GetTime()
	}
	if action == glfw.Release {
		switch key {
		case glfw.KeyEnter:
			view.OnSelect()
		}
	}
}

func (view *MenuView) OnSelect() {
	index := view.nx*(view.j+view.scroll) + view.i
	if index >= len(view.paths) {
		return
	}
	view.director.PlayGame(view.paths[index])
}

func (view *MenuView) Enter() {
	gl.ClearColor(0.333, 0.333, 0.333, 1)
	view.director.SetTitle("Select Game")
	view.director.window.SetKeyCallback(view.OnKey)
}

func (view *MenuView) Exit() {
	view.director.window.SetKeyCallback(nil)
}

func (view *MenuView) Update(t, dt float64) {
	view.texture.Purge()
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
	view.texture.Bind()
	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			x := float32(ox + i*sx)
			y := float32(oy + j*sy)
			index := nx*(j+view.scroll) + i
			if index >= len(view.paths) {
				continue
			}
			path := view.paths[index]
			tx, ty, tw, th := view.texture.Lookup(path)
			drawThumbnail(x, y, tx, ty, tw, th)
		}
	}
	view.texture.Unbind()
	if int((t-view.t)*4)%2 == 0 {
		x := float32(ox + view.i*sx)
		y := float32(oy + view.j*sy)
		drawSelection(x, y, 8, 4)
	}
	gl.PopMatrix()
	view.nx = nx
	view.ny = ny
}

func (view *MenuView) clampSelection(nx, ny int) {
	n := len(view.paths)
	rows := n / nx
	if n%nx > 0 {
		rows++
	}
	maxScroll := rows - ny
	if view.i < 0 {
		view.i = nx - 1
	}
	if view.i >= nx {
		view.i = 0
	}
	if view.j < 0 {
		view.j = 0
		view.scroll--
	}
	if view.j >= ny {
		view.j = ny - 1
		view.scroll++
	}
	if view.scroll < 0 {
		view.scroll = maxScroll
		view.j = ny - 1
	}
	if view.scroll > maxScroll {
		view.scroll = 0
		view.j = 0
	}
}

func drawThumbnail(x, y, tx, ty, tw, th float32) {
	sx := x + 4
	sy := y + 4
	gl.Disable(gl.TEXTURE_2D)
	gl.Color3f(0.2, 0.2, 0.2)
	gl.Begin(gl.QUADS)
	gl.Vertex2f(sx, sy)
	gl.Vertex2f(sx+256, sy)
	gl.Vertex2f(sx+256, sy+240)
	gl.Vertex2f(sx, sy+240)
	gl.End()
	gl.Enable(gl.TEXTURE_2D)
	gl.Color3f(1, 1, 1)
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(tx, ty)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(tx+tw, ty)
	gl.Vertex2f(x+256, y)
	gl.TexCoord2f(tx+tw, ty+th)
	gl.Vertex2f(x+256, y+240)
	gl.TexCoord2f(tx, ty+th)
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
