package nes

import "encoding/gob"

type Cartridge struct {
	PRG, CHR, SRAM []byte // PRG-ROM banks, CHR-ROM banks, // Save RAM
	Mapper, Mirror, Battery byte   // mapper type, mirroring type, battery present
}

func NewCartridge(prg, chr []byte, mapper, mirror, battery byte) *Cartridge {
	sram := make([]byte, 0x2000)
	return &Cartridge{prg, chr, sram, mapper, mirror, battery}
}

func (cartridge *Cartridge) Save(encoder *gob.Encoder) error {
	encoder.Encode(cartridge.SRAM)
	encoder.Encode(cartridge.Mirror)
	return nil
}

func (cartridge *Cartridge) Load(decoder *gob.Decoder) error {
	decoder.Decode(&cartridge.SRAM)
	decoder.Decode(&cartridge.Mirror)
	return nil
}
