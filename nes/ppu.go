package nes

import "log"

type PPU struct {
	Cycle    int    // 0-340
	ScanLine int    // 0-261, 0-239=visible, 240=post, 241-260=vblank, 261=pre
	Frame    uint64 // frame counter
	VB       byte   // vertical blank flag
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

// tick updates Cycle, ScanLine and Frame counters
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
}

func (ppu *PPU) Step() {
	ppu.tick()
	if ppu.ScanLine == 241 && ppu.Cycle == 1 {
		ppu.VB = 1
	}
	if ppu.ScanLine == 261 && ppu.Cycle == 1 {
		ppu.VB = 0
	}
}

func (ppu *PPU) readStatus() byte {
	var result byte
	result |= ppu.VB << 7
	ppu.VB = 0
	return result
}
