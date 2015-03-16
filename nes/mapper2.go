package nes

import "log"

type Mapper2 struct {
	*Cartridge
	prgBanks int
	chrBanks int
	prgBank1 int
	prgBank2 int
	chrBank  int
}

func NewMapper2(cartridge *Cartridge) Memory {
	prgBanks := len(cartridge.PRG) / 0x4000
	chrBanks := len(cartridge.CHR) / 0x2000
	prgBank1 := 0
	prgBank2 := prgBanks - 1
	chrBank := 0
	return &Mapper2{cartridge, prgBanks, chrBanks, prgBank1, prgBank2, chrBank}
}

func (m *Mapper2) Read(address uint16) byte {
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
		log.Fatalf("unhandled mapper2 read at address: 0x%04X", address)
	}
	return 0
}

func (m *Mapper2) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		index := m.chrBank*0x2000 + int(address)
		m.CHR[index] = value
	case address >= 0x8000:
		m.prgBank1 = int(value)
	case address >= 0x6000:
		index := int(address) - 0x6000
		m.SRAM[index] = value
	default:
		log.Fatalf("unhandled mapper2 write at address: 0x%04X", address)
	}
}
