package nes

import "log"

type PPU struct {
	Cycle    int
	ScanLine int
	Frame    uint64
	VB       byte
}

func NewPPU() *PPU {
	ppu := PPU{}
	ppu.Reset()
	return &ppu
}

func (ppu *PPU) Reset() {
	ppu.Cycle = 340
	ppu.ScanLine = 240
	ppu.Frame = 0
}

func (ppu *PPU) Read(address uint16) byte {
	switch address {
	case 0x2002:
		return ppu.readStatus()
	default:
		log.Fatalf("unhandled ppu read at address: 0x%04X", address)
	}
	return 0
}

func (ppu *PPU) Write(address uint16, value byte) {
	switch address {
	default:
		log.Fatalf("unhandled ppu write at address: 0x%04X", address)
	}
}

func (ppu *PPU) tick() {
	ppu.Cycle++
	if ppu.Cycle > 340 {
		ppu.Cycle = 0
		ppu.ScanLine++
		if ppu.ScanLine > 261 {
			ppu.ScanLine = 0
			ppu.Frame++
		}
	}
	if ppu.ScanLine == 241 && ppu.Cycle == 1 {
		ppu.VB = 1
	}
	if ppu.ScanLine == 261 && ppu.Cycle == 1 {
		ppu.VB = 0
	}
}

func (ppu *PPU) Step() {
	ppu.tick()
}

func (ppu *PPU) readStatus() byte {
	var result byte
	result |= ppu.VB << 7
	ppu.VB = 0
	return result
}
