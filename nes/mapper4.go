package nes

import (
	"fmt"
	"log"
)

type Mapper4 struct {
	*Cartridge
	mode           byte
	r0, r1, r2, r3 byte
	r4, r5, r6, r7 byte
	prgOffset0     int
	prgOffset1     int
	prgOffset2     int
	prgOffset3     int
}

func NewMapper4(cartridge *Cartridge) Mapper {
	m := Mapper4{Cartridge: cartridge}
	m.prgOffset0 = m.prgBankOffset(0)
	m.prgOffset1 = m.prgBankOffset(1)
	m.prgOffset2 = m.prgBankOffset(-2)
	m.prgOffset3 = m.prgBankOffset(-1)
	return &m
}

func (m *Mapper4) prgBankOffset(index int) int {
	offset := index * 0x2000
	if offset < 0 {
		offset += len(m.PRG)
	}
	return offset
}

func (m *Mapper4) Read(address uint16) byte {
	switch {
	// case address < 0x1000:
	// 	return m.CHR[m.chrOffset0+int(address)-0x0000]
	// case address < 0x2000:
	// 	return m.CHR[m.chrOffset1+int(address)-0x1000]
	case address >= 0xE000:
		return m.PRG[m.prgOffset3+int(address)-0xE000]
	case address >= 0xC000:
		return m.PRG[m.prgOffset2+int(address)-0xC000]
	case address >= 0xA000:
		return m.PRG[m.prgOffset1+int(address)-0xA000]
	case address >= 0x8000:
		return m.PRG[m.prgOffset0+int(address)-0x8000]
	case address >= 0x6000:
		return m.SRAM[int(address)-0x6000]
	default:
		log.Fatalf("unhandled mapper4 read at address: 0x%04X", address)
	}
	return 0
}

func (m *Mapper4) Write(address uint16, value byte) {
	switch {
	// case address < 0x1000:
	// 	m.CHR[m.chrOffset0+int(address)-0x0000] = value
	// case address < 0x2000:
	// 	m.CHR[m.chrOffset1+int(address)-0x1000] = value
	case address >= 0x8000:
		m.writeRegister(address, value)
	case address >= 0x6000:
		m.SRAM[int(address)-0x6000] = value
	default:
		log.Fatalf("unhandled mapper4 write at address: 0x%04X", address)
	}
}

func (m *Mapper4) writeRegister(address uint16, value byte) {
	fmt.Println(address, value)
	switch {
	case address <= 0x9FFF && address%2 == 0:
		m.writeBankSelect(value)
	case address <= 0x9FFF && address%2 == 1:
		m.writeBankData(value)
	case address <= 0xBFFF && address%2 == 0:
		m.writeMirror(value)
	case address <= 0xBFFF && address%2 == 1:
		m.writeProtect(value)
	case address <= 0xDFFF && address%2 == 0:
		m.writeIRQLatch(value)
	case address <= 0xDFFF && address%2 == 1:
		m.writeIRQReload(value)
	case address <= 0xFFFF && address%2 == 0:
		m.writeIRQDisable(value)
	case address <= 0xFFFF && address%2 == 1:
		m.writeIRQEnable(value)
	}
}

func (m *Mapper4) writeBankSelect(value byte) {
}

func (m *Mapper4) writeBankData(value byte) {
}

func (m *Mapper4) writeMirror(value byte) {
}

func (m *Mapper4) writeProtect(value byte) {
}

func (m *Mapper4) writeIRQLatch(value byte) {
}

func (m *Mapper4) writeIRQReload(value byte) {
}

func (m *Mapper4) writeIRQDisable(value byte) {
}

func (m *Mapper4) writeIRQEnable(value byte) {
}

func (m *Mapper4) updateOffsets() {
}
