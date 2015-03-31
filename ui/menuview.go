package ui

import "github.com/go-gl/gl/v2.1/gl"

const (
	border = 10
	margin = 10
)

type MenuView struct {
	director *Director
}

func NewMenuView(director *Director) View {
	return &MenuView{director}
}

func (view *MenuView) Enter() {
	view.director.SetTitle("Select Game")
}

func (view *MenuView) Exit() {
}

func (view *MenuView) Update(dt float64) {
	window := view.director.window
	w, h := window.GetFramebufferSize()
	sx := 256 + margin*2
	sy := 240 + margin*2
	nx := (w - border*2) / sx
	ny := (h - border*2) / sy
	ox := (w - nx*sx) / 2
	oy := (h - ny*sy) / 2
	gl.PushMatrix()
	gl.Ortho(0, float64(w), float64(h), 0, -1, 1)
	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			x := ox + i*sx
			y := oy + j*sy
			drawTile(float32(x), float32(y))
		}
	}
	gl.PopMatrix()
}

func drawTile(x, y float32) {
	gl.Color3f(1, 1, 1)
	gl.Begin(gl.QUADS)
	gl.Vertex2f(x, y)
	gl.Vertex2f(x+256, y)
	gl.Vertex2f(x+256, y+240)
	gl.Vertex2f(x, y+240)
	gl.End()
}
