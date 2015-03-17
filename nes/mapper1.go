package nes

import "log"

type Mapper1 struct {
	*Cartridge
	shiftRegister byte
	control       byte
	prgMode       byte
	chrMode       byte
	prgBank       byte
	chrBank0      byte
	chrBank1      byte
	prgOffset0    int
	prgOffset1    int
	chrOffset0    int
	chrOffset1    int
}

func NewMapper1(cartridge *Cartridge) Mapper {
	prgOffset1 := len(cartridge.PRG) - 0x4000
	return &Mapper1{cartridge, 0x10, 0, 0, 0, 0, 0, 0, 0, prgOffset1, 0, 0}
}

func (m *Mapper1) Step() {
}

func (m *Mapper1) Read(address uint16) byte {
	switch {
	case address < 0x1000:
		return m.CHR[m.chrOffset0+int(address)-0x0000]
	case address < 0x2000:
		return m.CHR[m.chrOffset1+int(address)-0x1000]
	case address >= 0xC000:
		return m.PRG[m.prgOffset1+int(address)-0xC000]
	case address >= 0x8000:
		return m.PRG[m.prgOffset0+int(address)-0x8000]
	case address >= 0x6000:
		return m.SRAM[int(address)-0x6000]
	default:
		log.Fatalf("unhandled mapper1 read at address: 0x%04X", address)
	}
	return 0
}

func (m *Mapper1) Write(address uint16, value byte) {
	switch {
	case address < 0x1000:
		m.CHR[m.chrOffset0+int(address)-0x0000] = value
	case address < 0x2000:
		m.CHR[m.chrOffset1+int(address)-0x1000] = value
	case address >= 0x8000:
		m.loadRegister(address, value)
	case address >= 0x6000:
		m.SRAM[int(address)-0x6000] = value
	default:
		log.Fatalf("unhandled mapper1 write at address: 0x%04X", address)
	}
}

func (m *Mapper1) loadRegister(address uint16, value byte) {
	if value&0x80 == 0x80 {
		m.shiftRegister = 0x10
		m.writeControl(m.control | 0x0C)
	} else {
		complete := m.shiftRegister&1 == 1
		m.shiftRegister >>= 1
		m.shiftRegister |= (value & 1) << 4
		if complete {
			m.writeRegister(address, m.shiftRegister)
			m.shiftRegister = 0x10
		}
	}
}

func (m *Mapper1) writeRegister(address uint16, value byte) {
	switch {
	case address <= 0x9FFF:
		m.writeControl(value)
	case address <= 0xBFFF:
		m.writeCHRBank0(value)
	case address <= 0xDFFF:
		m.writeCHRBank1(value)
	case address <= 0xFFFF:
		m.writePRGBank(value)
	}
}

// Control (internal, $8000-$9FFF)
func (m *Mapper1) writeControl(value byte) {
	m.control = value
	m.chrMode = (value >> 4) & 1
	m.prgMode = (value >> 2) & 3
	mirror := value & 3
	switch mirror {
	case 0:
		m.Cartridge.Mirror = MirrorSingle0
	case 1:
		m.Cartridge.Mirror = MirrorSingle1
	case 2:
		m.Cartridge.Mirror = MirrorVertical
	case 3:
		m.Cartridge.Mirror = MirrorHorizontal
	}
	m.updateOffsets()
}

// CHR bank 0 (internal, $A000-$BFFF)
func (m *Mapper1) writeCHRBank0(value byte) {
	m.chrBank0 = value
	m.updateOffsets()
}

// CHR bank 1 (internal, $C000-$DFFF)
func (m *Mapper1) writeCHRBank1(value byte) {
	m.chrBank1 = value
	m.updateOffsets()
}

// PRG bank (internal, $E000-$FFFF)
func (m *Mapper1) writePRGBank(value byte) {
	m.prgBank = value & 0x0F
	m.updateOffsets()
}

// PRG ROM bank mode (0, 1: switch 32 KB at $8000, ignoring low bit of bank number;
//                    2: fix first bank at $8000 and switch 16 KB bank at $C000;
//                    3: fix last bank at $C000 and switch 16 KB bank at $8000)
// CHR ROM bank mode (0: switch 8 KB at a time; 1: switch two separate 4 KB banks)
func (m *Mapper1) updateOffsets() {
	switch m.prgMode {
	case 0, 1:
		m.prgOffset0 = int(m.prgBank&0xFE) * 0x4000
		m.prgOffset1 = m.prgOffset0 + 0x4000
	case 2:
		m.prgOffset0 = 0
		m.prgOffset1 = int(m.prgBank) * 0x4000
	case 3:
		m.prgOffset0 = int(m.prgBank) * 0x4000
		m.prgOffset1 = len(m.Cartridge.PRG) - 0x4000
	}
	switch m.chrMode {
	case 0:
		m.chrOffset0 = int(m.chrBank0&0xFE) * 0x1000
		m.chrOffset1 = m.chrOffset0 + 0x1000
	case 1:
		m.chrOffset0 = int(m.chrBank0) * 0x1000
		m.chrOffset1 = int(m.chrBank1) * 0x1000
	}
}
