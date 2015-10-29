package nes

import (
	"encoding/gob"
	"log"
)

type Mapper2 struct {
	*Cartridge
	prgBanks int
	prgBank1 int
	prgBank2 int
}

func NewMapper2(cartridge *Cartridge) Mapper {
	prgBanks := len(cartridge.PRG) / 0x4000
	prgBank1 := 0
	prgBank2 := prgBanks - 1
	return &Mapper2{cartridge, prgBanks, prgBank1, prgBank2}
}

func (m *Mapper2) Save(encoder *gob.Encoder) error {
	encoder.Encode(m.prgBanks)
	encoder.Encode(m.prgBank1)
	encoder.Encode(m.prgBank2)
	return nil
}

func (m *Mapper2) Load(decoder *gob.Decoder) error {
	decoder.Decode(&m.prgBanks)
	decoder.Decode(&m.prgBank1)
	decoder.Decode(&m.prgBank2)
	return nil
}

func (m *Mapper2) Step() {
}

func (m *Mapper2) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		return m.CHR[address]
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
		log.Fatalf("unhandled mapper2 read at address: 0x%04X", address)
	}
	return 0
}

func (m *Mapper2) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		m.CHR[address] = value
	case address >= 0x8000:
		m.prgBank1 = int(value) % m.prgBanks
	case address >= 0x6000:
		index := int(address) - 0x6000
		m.SRAM[index] = value
	default:
		log.Fatalf("unhandled mapper2 write at address: 0x%04X", address)
	}
}
