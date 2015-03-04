package nes

const (
	MirrorHorizontal = 0
	MirrorVertical   = 1
	MirrorQuad       = 2
)

type Cartridge struct {
	PRG     []byte // PRG-ROM banks
	CHR     []byte // CHR-ROM banks
	Mapper  int    // mapper type
	Mirror  int    // mirroring mode
	Battery bool   // battery present
}

func (c *Cartridge) Read(address uint16) byte {
	index := (int(address) - 0x8000) % len(c.PRG)
	return c.PRG[index]
}
