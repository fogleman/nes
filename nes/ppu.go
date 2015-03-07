package nes

import (
	"image"
	"image/color"
	"log"
)

var palette [64]color.RGBA = [64]color.RGBA{
	color.RGBA{0x66, 0x66, 0x66, 0xFF}, color.RGBA{0x00, 0x2A, 0x88, 0xFF},
	color.RGBA{0x14, 0x12, 0xA7, 0xFF}, color.RGBA{0x3B, 0x00, 0xA4, 0xFF},
	color.RGBA{0x5C, 0x00, 0x7E, 0xFF}, color.RGBA{0x6E, 0x00, 0x40, 0xFF},
	color.RGBA{0x6C, 0x06, 0x00, 0xFF}, color.RGBA{0x56, 0x1D, 0x00, 0xFF},
	color.RGBA{0x33, 0x35, 0x00, 0xFF}, color.RGBA{0x0B, 0x48, 0x00, 0xFF},
	color.RGBA{0x00, 0x52, 0x00, 0xFF}, color.RGBA{0x00, 0x4F, 0x08, 0xFF},
	color.RGBA{0x00, 0x40, 0x4D, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF},
	color.RGBA{0x00, 0x00, 0x00, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF},
	color.RGBA{0xAD, 0xAD, 0xAD, 0xFF}, color.RGBA{0x15, 0x5F, 0xD9, 0xFF},
	color.RGBA{0x42, 0x40, 0xFF, 0xFF}, color.RGBA{0x75, 0x27, 0xFE, 0xFF},
	color.RGBA{0xA0, 0x1A, 0xCC, 0xFF}, color.RGBA{0xB7, 0x1E, 0x7B, 0xFF},
	color.RGBA{0xB5, 0x31, 0x20, 0xFF}, color.RGBA{0x99, 0x4E, 0x00, 0xFF},
	color.RGBA{0x6B, 0x6D, 0x00, 0xFF}, color.RGBA{0x38, 0x87, 0x00, 0xFF},
	color.RGBA{0x0C, 0x93, 0x00, 0xFF}, color.RGBA{0x00, 0x8F, 0x32, 0xFF},
	color.RGBA{0x00, 0x7C, 0x8D, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF},
	color.RGBA{0x00, 0x00, 0x00, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF},
	color.RGBA{0xFF, 0xFE, 0xFF, 0xFF}, color.RGBA{0x64, 0xB0, 0xFF, 0xFF},
	color.RGBA{0x92, 0x90, 0xFF, 0xFF}, color.RGBA{0xC6, 0x76, 0xFF, 0xFF},
	color.RGBA{0xF3, 0x6A, 0xFF, 0xFF}, color.RGBA{0xFE, 0x6E, 0xCC, 0xFF},
	color.RGBA{0xFE, 0x81, 0x70, 0xFF}, color.RGBA{0xEA, 0x9E, 0x22, 0xFF},
	color.RGBA{0xBC, 0xBE, 0x00, 0xFF}, color.RGBA{0x88, 0xD8, 0x00, 0xFF},
	color.RGBA{0x5C, 0xE4, 0x30, 0xFF}, color.RGBA{0x45, 0xE0, 0x82, 0xFF},
	color.RGBA{0x48, 0xCD, 0xDE, 0xFF}, color.RGBA{0x4F, 0x4F, 0x4F, 0xFF},
	color.RGBA{0x00, 0x00, 0x00, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF},
	color.RGBA{0xFF, 0xFE, 0xFF, 0xFF}, color.RGBA{0xC0, 0xDF, 0xFF, 0xFF},
	color.RGBA{0xD3, 0xD2, 0xFF, 0xFF}, color.RGBA{0xE8, 0xC8, 0xFF, 0xFF},
	color.RGBA{0xFB, 0xC2, 0xFF, 0xFF}, color.RGBA{0xFE, 0xC4, 0xEA, 0xFF},
	color.RGBA{0xFE, 0xCC, 0xC5, 0xFF}, color.RGBA{0xF7, 0xD8, 0xA5, 0xFF},
	color.RGBA{0xE4, 0xE5, 0x94, 0xFF}, color.RGBA{0xCF, 0xEF, 0x96, 0xFF},
	color.RGBA{0xBD, 0xF4, 0xAB, 0xFF}, color.RGBA{0xB3, 0xF3, 0xCC, 0xFF},
	color.RGBA{0xB5, 0xEB, 0xF2, 0xFF}, color.RGBA{0xB8, 0xB8, 0xB8, 0xFF},
	color.RGBA{0x00, 0x00, 0x00, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF},
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

	// $2007 PPUDATA
	data byte // for buffered reads

	paletteData   [32]byte
	nameTableData [2048]byte

	buffer *image.RGBA // color buffer
}

func NewPPU(nes *NES) *PPU {
	ppu := PPU{Memory: nes.PPUMemory, nes: nes}
	ppu.buffer = image.NewRGBA(image.Rect(0, 0, 256, 240))
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

func (ppu *PPU) readRegister(address uint16) byte {
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

func (ppu *PPU) writeRegister(address uint16, value byte) {
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
	// emulate buffered reads
	if ppu.address%0x4000 < 0x3F00 {
		buffered := ppu.data
		ppu.data = value
		value = buffered
	} else {
		ppu.data = ppu.Read(ppu.address - 0x1000)
	}
	// increment address
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

func (ppu *PPU) readPalette(address uint16) byte {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	return ppu.paletteData[address]
}

func (ppu *PPU) writePalette(address uint16, value byte) {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	ppu.paletteData[address] = value
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
	if ppu.Cycle == 1 && ppu.ScanLine < 240 {
		ppu.renderScanLine()
	}
	if ppu.Cycle == 1 && ppu.ScanLine == 241 {
		ppu.VerticalBlank = 1
	}
	if ppu.Cycle == 1 && ppu.ScanLine == 261 {
		ppu.VerticalBlank = 0
	}
}

func (ppu *PPU) tileAttribute(nameTable, x, y int) byte {
	gx := x / 4
	gy := y / 4
	sx := (x % 4) / 2
	sy := (y % 4) / 2
	address := uint16(0x23c0 + 0x400*nameTable + gy*8 + gx)
	attribute := ppu.Read(address)
	shift := byte((sy*2 + sx) * 2)
	return (attribute >> shift) & 3
}

func (ppu *PPU) tilePattern(nameTable, x, y, row int) (byte, byte) {
	// fetch pattern index from name table
	index := y*32 + x
	address := uint16(0x2000 + 0x400*nameTable + index)
	pattern := int(ppu.Read(address))
	// fetch pattern data from pattern table
	patternTable := int(ppu.flagBackgroundTable)
	patternAddress1 := uint16(0x1000*patternTable + pattern*16 + row)
	patternAddress2 := patternAddress1 + 8
	pattern1 := ppu.Read(patternAddress1)
	pattern2 := ppu.Read(patternAddress2)
	return pattern1, pattern2
}

func (ppu *PPU) tileRow(nameTable, x, y, row int) [8]byte {
	attribute := ppu.tileAttribute(nameTable, x, y) << 2
	pattern1, pattern2 := ppu.tilePattern(nameTable, x, y, row)
	var result [8]byte
	for i := 0; i < 8; i++ {
		p1 := (pattern1 & 1)
		p2 := (pattern2 & 1) << 1
		index := attribute | p1 | p2
		result[7-i] = ppu.paletteData[index]
		pattern1 >>= 1
		pattern2 >>= 1
	}
	return result
}

func (ppu *PPU) renderScanLine() {
	nameTable := int(ppu.flagNameTable)
	y := ppu.ScanLine
	ty := y / 8
	row := y % 8
	for tx := 0; tx < 32; tx++ {
		tile := ppu.tileRow(nameTable, tx, ty, row)
		for i := 0; i < 8; i++ {
			x := tx*8 + i
			c := palette[tile[i]]
			ppu.buffer.SetRGBA(x, y, c)
		}
	}
}
