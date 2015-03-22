package nes

type Cartridge struct {
	PRG, CHR, SRAM []byte // PRG-ROM banks, CHR-ROM banks, // Save RAM
	Mapper, Mirror, Battery byte   // mapper type, mirroring type, battery present
}

func NewCartridge(prg, chr []byte, mapper, mirror, battery byte) *Cartridge {
	sram := make([]byte, 0x2000)
	return &Cartridge{prg, chr, sram, mapper, mirror, battery}
}
