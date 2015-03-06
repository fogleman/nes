package nes

import "log"

var palette [64]uint32 = [64]uint32{
	0x666666, 0x002A88, 0x1412A7, 0x3B00A4,
	0x5C007E, 0x6E0040, 0x6C0600, 0x561D00,
	0x333500, 0x0B4800, 0x005200, 0x004F08,
	0x00404D, 0x000000, 0x000000, 0x000000,
	0xADADAD, 0x155FD9, 0x4240FF, 0x7527FE,
	0xA01ACC, 0xB71E7B, 0xB53120, 0x994E00,
	0x6B6D00, 0x388700, 0x0C9300, 0x008F32,
	0x007C8D, 0x000000, 0x000000, 0x000000,
	0xFFFEFF, 0x64B0FF, 0x9290FF, 0xC676FF,
	0xF36AFF, 0xFE6ECC, 0xFE8170, 0xEA9E22,
	0xBCBE00, 0x88D800, 0x5CE430, 0x45E082,
	0x48CDDE, 0x4F4F4F, 0x000000, 0x000000,
	0xFFFEFF, 0xC0DFFF, 0xD3D2FF, 0xE8C8FF,
	0xFBC2FF, 0xFEC4EA, 0xFECCC5, 0xF7D8A5,
	0xE4E594, 0xCFEF96, 0xBDF4AB, 0xB3F3CC,
	0xB5EBF2, 0xB8B8B8, 0x000000, 0x000000,
}

type PPU struct {
	Memory      // memory interface
	nes    *NES // reference to parent object

	Cycle         int    // 0-340
	ScanLine      int    // 0-261, 0-239=visible, 240=post, 241-260=vblank, 261=pre
	Frame         uint64 // frame counter
	VerticalBlank byte   // vertical blank status

	// $2000 PPUCTRL
	flagNameTable       byte // 0: $2000; 1: $2400; 2: $2800; 3: $2C00
	flagIncrement       byte // 0: add 1; 1: add 32
	flagSpriteTable     byte // 0: $0000; 1: $1000; ignored in 8x16 mode
	flagBackgroundTable byte // 0: $0000; 1: $1000
	flagSpriteSize      byte // 0: 8x8; 1: 8x16
	flagMasterSlave     byte // 0: read EXT; 1: write EXT
	flagGenerateNMI     byte // 0: off; 1: on

	// $2001 PPUMASK
	flagGrayscale          byte // 0: color; 1: grayscale
	flagShowLeftBackground byte // 0: hide; 1: show
	flagShowLeftSprites    byte // 0: hide; 1: show
	flagShowBackground     byte // 0: hide; 1: show
	flagShowSprites        byte // 0: hide; 1: show
	flagRedTint            byte // 0: normal; 1: emphasized
	flagGreenTint          byte // 0: normal; 1: emphasized
	flagBlueTint           byte // 0: normal; 1: emphasized

	// $2003 OAMADDR
	oamAddress byte

	// $2004 OAMDATA
	oamData [256]byte

	// $2005 PPUSCROLL
	scroll uint16 // x & y scrolling coordinates

	// $2006 PPUADDR
	address uint16 // address used by $2007 PPUDATA

	paletteData   [32]byte
	nameTableData [2048]byte
}

func NewPPU(nes *NES) *PPU {
	ppu := PPU{Memory: nes.PPUMemory, nes: nes}
	ppu.Reset()
	return &ppu
}

func (ppu *PPU) Reset() {
	ppu.Cycle = 340
	ppu.ScanLine = 240
	ppu.Frame = 0
	ppu.VerticalBlank = 0
	ppu.writeControl(0)
	ppu.writeMask(0)
	ppu.writeOAMAddress(0)
}

func (ppu *PPU) ReadRegister(address uint16) byte {
	switch address {
	case 0x2002:
		return ppu.readStatus()
	case 0x2004:
		return ppu.readOAMData()
	case 0x2007:
		return ppu.readData()
	default:
		log.Fatalf("unhandled ppu register read at address: 0x%04X", address)
	}
	return 0
}

func (ppu *PPU) WriteRegister(address uint16, value byte) {
	switch address {
	case 0x2000:
		ppu.writeControl(value)
	case 0x2001:
		ppu.writeMask(value)
	case 0x2003:
		ppu.writeOAMAddress(value)
	case 0x2004:
		ppu.writeOAMData(value)
	case 0x2005:
		ppu.writeScroll(value)
	case 0x2006:
		ppu.writeAddress(value)
	case 0x2007:
		ppu.writeData(value)
	case 0x4014:
		ppu.writeDMA(value)
	default:
		log.Fatalf("unhandled ppu register write at address: 0x%04X", address)
	}
}

// $2000: PPUCTRL
func (ppu *PPU) writeControl(value byte) {
	ppu.flagNameTable = (value >> 0) & 3
	ppu.flagIncrement = (value >> 2) & 1
	ppu.flagSpriteTable = (value >> 3) & 1
	ppu.flagBackgroundTable = (value >> 4) & 1
	ppu.flagSpriteSize = (value >> 5) & 1
	ppu.flagMasterSlave = (value >> 6) & 1
	ppu.flagGenerateNMI = (value >> 7) & 1
}

// $2001: PPUMASK
func (ppu *PPU) writeMask(value byte) {
	ppu.flagGrayscale = (value >> 0) & 1
	ppu.flagShowLeftBackground = (value >> 1) & 1
	ppu.flagShowLeftSprites = (value >> 2) & 1
	ppu.flagShowBackground = (value >> 3) & 1
	ppu.flagShowSprites = (value >> 4) & 1
	ppu.flagRedTint = (value >> 5) & 1
	ppu.flagGreenTint = (value >> 6) & 1
	ppu.flagBlueTint = (value >> 7) & 1
}

// $2002: PPUSTATUS
func (ppu *PPU) readStatus() byte {
	var result byte
	result |= ppu.VerticalBlank << 7
	ppu.VerticalBlank = 0
	return result
}

// $2003: OAMADDR
func (ppu *PPU) writeOAMAddress(value byte) {
	ppu.oamAddress = value
}

// $2004: OAMDATA (read)
func (ppu *PPU) readOAMData() byte {
	return ppu.oamData[ppu.oamAddress]
}

// $2004: OAMDATA (write)
func (ppu *PPU) writeOAMData(value byte) {
	ppu.oamData[ppu.oamAddress] = value
	ppu.oamAddress++
}

// $2005: PPUSCROLL
func (ppu *PPU) writeScroll(value byte) {
	ppu.scroll <<= 8
	ppu.scroll |= uint16(value)
}

// $2006: PPUADDR
func (ppu *PPU) writeAddress(value byte) {
	ppu.address <<= 8
	ppu.address |= uint16(value)
}

// $2007: PPUDATA (read)
func (ppu *PPU) readData() byte {
	value := ppu.Read(ppu.address)
	if ppu.flagIncrement == 0 {
		ppu.address += 1
	} else {
		ppu.address += 32
	}
	return value
}

// $2007: PPUDATA (write)
func (ppu *PPU) writeData(value byte) {
	ppu.Write(ppu.address, value)
	if ppu.flagIncrement == 0 {
		ppu.address += 1
	} else {
		ppu.address += 32
	}
}

// $4014: OAMDMA
func (ppu *PPU) writeDMA(value byte) {
	// TODO: stall CPU for 513 or 514 cycles
	cpu := ppu.nes.CPU
	address := uint16(value) << 8
	for i := 0; i < 256; i++ {
		ppu.oamData[ppu.oamAddress] = cpu.Read(address)
		ppu.oamAddress++
		address++
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

// Step executes a single PPU cycle
func (ppu *PPU) Step() {
	ppu.tick()
	if ppu.ScanLine == 241 && ppu.Cycle == 1 {
		ppu.VerticalBlank = 1
	}
	if ppu.ScanLine == 261 && ppu.Cycle == 1 {
		ppu.VerticalBlank = 0
	}
}

func (ppu *PPU) tileAttribute(nameTable, x, y int) byte {
	gx := x / 4
	gy := y / 4
	sx := (x % 4) / 2
	sy := (y % 4) / 2
	nameAddress := uint16(0x23c0 + 0x1000*nameTable + gy*8 + gx)
	attribute := ppu.Read(nameAddress)
	shift := byte((sy*2 + sx) * 2)
	return (attribute >> shift) & 3
}

func (ppu *PPU) tilePattern(nameTable, patternTable, x, y, row int) uint16 {
	index := y*32 + x
	nameAddress := uint16(0x2000 + 0x1000*nameTable + index)
	pattern := int(ppu.Read(nameAddress))
	patternAddress1 := uint16(0x1000*patternTable + pattern*16 + row)
	patternAddress2 := patternAddress1 + 8
	pattern1 := uint16(ppu.Read(patternAddress1))
	pattern2 := uint16(ppu.Read(patternAddress2))
	return pattern1<<8 | pattern2
}
