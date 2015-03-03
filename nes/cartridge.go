package nes

const (
	MirrorHorizontal = 0
	MirrorVertical   = 1
	MirrorQuad       = 2
)

type Cartridge struct {
	PRG     []byte
	CHR     []byte
	Mapper  int
	Mirror  int
	Battery bool
}
