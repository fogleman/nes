package nes

import (
	"encoding/gob"
	"log"
)

// https://github.com/asfdfdfd/fceux/blob/master/src/boards/225.cpp
// https://wiki.nesdev.com/w/index.php/INES_Mapper_225

type Mapper225 struct {
	*Cartridge
	chrBank  int
	prgBank1 int
	prgBank2 int
}

func NewMapper225(cartridge *Cartridge) Mapper {
	prgBanks := len(cartridge.PRG) / 0x4000
	return &Mapper225{cartridge, 0, 0, prgBanks - 1}
}

func (m *Mapper225) Save(encoder *gob.Encoder) error {
	encoder.Encode(m.chrBank)
	encoder.Encode(m.prgBank1)
	encoder.Encode(m.prgBank2)
	return nil
}

func (m *Mapper225) Load(decoder *gob.Decoder) error {
	decoder.Decode(&m.chrBank)
	decoder.Decode(&m.prgBank1)
	decoder.Decode(&m.prgBank2)
	return nil
}

func (m *Mapper225) Step() {
}

func (m *Mapper225) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		index := m.chrBank*0x2000 + int(address)
		return m.CHR[index]
	case address >= 0xC000:
		index := m.prgBank2*0x4000 + int(address-0xC000)
		return m.PRG[index]
	case address >= 0x8000:
		index := m.prgBank1*0x4000 + int(address-0x8000)
		return m.PRG[index]
	case address >= 0x6000:
		index := int(address) - 0x6000
		return m.SRAM[index]
	default:
		log.Fatalf("unhandled Mapper225 read at address: 0x%04X", address)
	}
	return 0
}

func (m *Mapper225) Write(address uint16, value byte) {
	if (address < 0x8000) {
		return
	}

	A := int(address)
	bank := (A >> 14) & 1
	m.chrBank = (A & 0x3f) | (bank << 6)
	prg := ((A >> 6) & 0x3f) | (bank << 6)
	mode := (A >> 12) & 1;
	if (mode == 1) {
		m.prgBank1 = prg
		m.prgBank2 = prg
	} else {
		m.prgBank1 = prg
		m.prgBank2 = prg + 1
	}
	mirr := (A >> 13) & 1
	if (mirr == 1) {
		m.Cartridge.Mirror = MirrorHorizontal
	} else {
		m.Cartridge.Mirror = MirrorVertical
	}

	// fmt.Println(address, mirr, mode, prg)
}
