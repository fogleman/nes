package ui

import "github.com/go-gl/gl/v2.1/gl"

const textureSize = 4096
const textureDim = textureSize / 256
const textureCount = textureDim * textureDim

type Texture struct {
	cache   *Cache
	texture uint32
	lookup  map[string]int
	reverse [textureCount]string
	access  [textureCount]int
	counter int
}

func NewTexture(cache *Cache) *Texture {
	texture := createTexture()
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		textureSize, textureSize,
		0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	t := Texture{}
	t.cache = cache
	t.texture = texture
	t.lookup = make(map[string]int)
	return &t
}

func (t *Texture) Purge() {
	for {
		select {
		case path := <-t.cache.ch:
			delete(t.lookup, path)
		default:
			return
		}
	}
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.texture)
}

func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
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
	delete(t.lookup, t.reverse[index])
	t.mark(index)
	t.lookup[path] = index
	t.reverse[index] = path
	x := int32((index % textureDim) * 256)
	y := int32((index / textureDim) * 256)
	im := copyImage(t.cache.LoadThumbnail(path))
	size := im.Rect.Size()
	gl.TexSubImage2D(
		gl.TEXTURE_2D, 0, x, y, int32(size.X), int32(size.Y),
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(im.Pix))
	return index
}
