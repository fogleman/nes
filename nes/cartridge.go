package nes

import "encoding/gob"

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

func (cartridge *Cartridge) Save(encoder *gob.Encoder) error {
	encoder.Encode(cartridge.PRG)
	encoder.Encode(cartridge.CHR)
	encoder.Encode(cartridge.SRAM)
	encoder.Encode(cartridge.Mirror)
	return nil
}

func (cartridge *Cartridge) Load(decoder *gob.Decoder) error {
	decoder.Decode(&cartridge.PRG)
	decoder.Decode(&cartridge.CHR)
	decoder.Decode(&cartridge.SRAM)
	decoder.Decode(&cartridge.Mirror)
	return nil
}
