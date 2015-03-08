package nes

import (
	"image"
	"log"
)

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
		if ppu.flagGenerateNMI != 0 {
			ppu.nes.CPU.triggerNMI()
		}
	}
	if ppu.Cycle == 1 && ppu.ScanLine == 261 {
		ppu.VerticalBlank = 0
	}
}

func (ppu *PPU) pattern(attribute, patternTable, pattern, row int) [8]byte {
	patternAddress1 := uint16(0x1000*patternTable + pattern*16 + row)
	patternAddress2 := patternAddress1 + 8
	pattern1 := ppu.Read(patternAddress1)
	pattern2 := ppu.Read(patternAddress2)
	var result [8]byte
	for i := 0; i < 8; i++ {
		p1 := (pattern1 & 1)
		p2 := (pattern2 & 1) << 1
		index := byte(attribute) | p1 | p2
		result[7-i] = ppu.readPalette(uint16(index))
		pattern1 >>= 1
		pattern2 >>= 1
	}
	return result
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

func (ppu *PPU) tilePattern(attribute, nameTable, x, y, row int) [8]byte {
	index := y*32 + x
	address := uint16(0x2000 + 0x400*nameTable + index)
	pattern := int(ppu.Read(address))
	patternTable := int(ppu.flagBackgroundTable)
	return ppu.pattern(attribute, patternTable, pattern, row)
}

func (ppu *PPU) tileRow(nameTable, x, y, row int) [8]byte {
	attribute := int(ppu.tileAttribute(nameTable, x, y) << 2)
	return ppu.tilePattern(attribute, nameTable, x, y, row)
}

func (ppu *PPU) renderNameTableLine(nameTable, y int) []byte {
	result := make([]byte, 256)
	ty := y / 8
	row := y % 8
	for tx := 0; tx < 32; tx++ {
		tile := ppu.tileRow(nameTable, tx, ty, row)
		for i := 0; i < 8; i++ {
			x := tx*8 + i
			result[x] = tile[i]
		}
	}
	return result
}

func (ppu *PPU) renderSpriteLine() []byte {
	result := make([]byte, 256)
	for i := 0; i < 64; i++ {
		index := i * 4
		y := ppu.oamData[index+0]
		t := ppu.oamData[index+1]
		f := ppu.oamData[index+2]
		x := ppu.oamData[index+3]
		row := int(y) - ppu.ScanLine
		if row < 0 || row > 7 {
			continue
		}
		pattern := t
		bank := ppu.flagSpriteTable
		if ppu.flagSpriteSize == 1 {
			bank = t & 1
			pattern = t & 0xFE
		}
		attribute := (f&3)<<2 | 16
		tile := ppu.pattern(int(attribute), int(bank), int(pattern), row)
		for j := 0; j < 8; j++ {
			index := int(x) + j
			if index > 255 {
				continue
			}
			if result[index]%4 == 0 {
				result[index] = tile[j]
			}
		}
	}
	return result
}

func (ppu *PPU) renderScanLine() {
	sx := int(ppu.scroll >> 8)
	sy := int(ppu.scroll & 0xFF)

	y := ppu.ScanLine + sy
	nameTable := int(ppu.flagNameTable)
	if y >= 240 {
		y -= 240
		nameTable += 2
	}
	nameTable1 := (nameTable + 0) % 4
	nameTable2 := (nameTable + 1) % 4

	line1 := ppu.renderNameTableLine(nameTable1, y)
	line2 := ppu.renderNameTableLine(nameTable2, y)
	line := make([]byte, 0, 512)
	line = append(line, line1...)
	line = append(line, line2...)
	sprites := ppu.renderSpriteLine()

	for i := 0; i < 256; i++ {
		background := line[sx+i]
		sprite := sprites[i]
		p := sprite
		if sprite%4 == 0 {
			p = background
		}
		c := palette[p]
		ppu.buffer.SetRGBA(i, ppu.ScanLine, c)
	}
}
