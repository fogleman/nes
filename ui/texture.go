package ui

import (
	"image"
	"image/draw"

	"github.com/go-gl/gl/v2.1/gl"
)

const textureSize = 4096
const textureDim = textureSize / 256
const textureCount = textureDim * textureDim

type Texture struct {
	texture uint32
	im      *image.RGBA
	lookup  map[string]int
	access  [textureCount]int
	counter int
	dirty   bool
}

func NewTexture() *Texture {
	t := Texture{}
	t.texture = createTexture()
	t.im = image.NewRGBA(image.Rect(0, 0, textureSize, textureSize))
	t.lookup = make(map[string]int)
	return &t
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.texture)
}

func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (t *Texture) Sync() {
	if t.dirty {
		t.Bind()
		setTexture(t.im)
		t.Unbind()
		t.dirty = false
	}
}

func (t *Texture) Lookup(path string) (x, y, dx, dy float32) {
	if index, ok := t.lookup[path]; ok {
		return t.coord(index)
	} else {
		return t.coord(t.load(path))
	}
}

func (t *Texture) mark(index int) {
	t.counter++
	t.access[index] = t.counter
}

func (t *Texture) lru() int {
	minIndex := 0
	minValue := t.counter + 1
	for i, n := range t.access {
		if n < minValue {
			minIndex = i
			minValue = n
		}
	}
	return minIndex
}

func (t *Texture) coord(index int) (x, y, dx, dy float32) {
	x = float32(index%textureDim) / textureDim
	y = float32(index/textureDim) / textureDim
	dx = 1.0 / textureDim
	dy = 1.0 / textureDim
	return
}

func (t *Texture) load(path string) int {
	index := t.lru()
	t.mark(index)
	t.lookup[path] = index
	t.dirty = true
	x := (index % textureDim) * 256
	y := (index / textureDim) * 256
	r := image.Rect(x, y, x+256, y+240)
	im := CreateGenericThumbnail(path)
	draw.Draw(t.im, r, im, image.ZP, draw.Src)
	return index
}
