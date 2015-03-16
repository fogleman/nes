package nes

type Cartridge struct {
	PRG     []byte // PRG-ROM banks
	CHR     []byte // CHR-ROM banks
	SRAM    []byte // Save RAM
	Mapper  byte   // mapper type
	Mirror  byte   // mirroring mode
	Battery byte   // battery present
}

func NewCartridge(prg, chr []byte, mapper, mirror, battery byte) *Cartridge {
	sram := make([]byte, 0x2000)
	return &Cartridge{prg, chr, sram, mapper, mirror, battery}
}
