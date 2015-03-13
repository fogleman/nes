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

	// PPU registers
	v uint16 // current vram address (15 bit)
	t uint16 // temporary vram address (15 bit)
	x byte   // fine x scroll (3 bit)
	w byte   // write toggle (1 bit)
	f byte   // even/odd frame flag (1 bit)

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
	oamPrimary   [256]byte
	oamSecondary [32]byte

	// $2005 PPUSCROLL
	scroll uint16 // x & y scrolling coordinates

	// $2007 PPUDATA
	data byte // for buffered reads

	paletteData   [32]byte
	nameTableData [2048]byte
	tileData      [128]byte
	tileIndex     byte

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
	// t: ....BA.. ........ = d: ......BA
	ppu.t = (ppu.t & 0xF3FF) | ((uint16(value) & 0x03) << 10)
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
	// w:                   = 0
	ppu.w = 0
	return result
}

// $2003: OAMADDR
func (ppu *PPU) writeOAMAddress(value byte) {
	ppu.oamAddress = value
}

// $2004: OAMDATA (read)
func (ppu *PPU) readOAMData() byte {
	return ppu.oamPrimary[ppu.oamAddress]
}

// $2004: OAMDATA (write)
func (ppu *PPU) writeOAMData(value byte) {
	ppu.oamPrimary[ppu.oamAddress] = value
	ppu.oamAddress++
}

// $2005: PPUSCROLL
func (ppu *PPU) writeScroll(value byte) {
	ppu.scroll <<= 8
	ppu.scroll |= uint16(value)
	if ppu.w == 0 {
		// t: ........ ...HGFED = d: HGFED...
		// x:               CBA = d: .....CBA
		// w:                   = 1
		ppu.t = (ppu.t & 0xFFE0) | (uint16(value) >> 3)
		ppu.x = value & 0x07
		ppu.w = 1
	} else {
		// t: .CBA..HG FED..... = d: HGFEDCBA
		// w:                   = 0
		ppu.t = (ppu.t & 0x8FFF) | ((uint16(value) & 0x03) << 12)
		ppu.t = (ppu.t & 0xFC1F) | ((uint16(value) & 0xF8) << 2)
		ppu.w = 0
	}
}

// $2006: PPUADDR
func (ppu *PPU) writeAddress(value byte) {
	if ppu.w == 0 {
		// t: ..FEDCBA ........ = d: ..FEDCBA
		// t: .X...... ........ = 0
		// w:                   = 1
		ppu.t = (ppu.t & 0x80FF) | ((uint16(value) & 0x3F) << 8)
		ppu.w = 1
	} else {
		// t: ........ HGFEDCBA = d: HGFEDCBA
		// v                    = t
		// w:                   = 0
		ppu.t = (ppu.t & 0xFF00) | uint16(value)
		ppu.v = ppu.t
		ppu.w = 0
	}
}

// $2007: PPUDATA (read)
func (ppu *PPU) readData() byte {
	value := ppu.Read(ppu.v)
	// emulate buffered reads
	if ppu.v%0x4000 < 0x3F00 {
		buffered := ppu.data
		ppu.data = value
		value = buffered
	} else {
		ppu.data = ppu.Read(ppu.v - 0x1000)
	}
	// increment address
	if ppu.flagIncrement == 0 {
		ppu.v += 1
	} else {
		ppu.v += 32
	}
	return value
}

// $2007: PPUDATA (write)
func (ppu *PPU) writeData(value byte) {
	ppu.Write(ppu.v, value)
	if ppu.flagIncrement == 0 {
		ppu.v += 1
	} else {
		ppu.v += 32
	}
}

// $4014: OAMDMA
func (ppu *PPU) writeDMA(value byte) {
	// TODO: stall CPU for 513 or 514 cycles
	cpu := ppu.nes.CPU
	address := uint16(value) << 8
	for i := 0; i < 256; i++ {
		ppu.oamPrimary[ppu.oamAddress] = cpu.Read(address)
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

func (ppu *PPU) incrementX() {
	// if coarse X == 31
	if ppu.v&0x001F == 31 {
		// coarse X = 0
		ppu.v &= 0xFFE0
		// switch horizontal nametable
		ppu.v ^= 0x0400
	} else {
		// increment coarse X
		ppu.v++
	}
}

func (ppu *PPU) incrementY() {
	// if fine Y < 7
	if ppu.v&0x7000 != 0x7000 {
		// increment fine Y
		ppu.v += 0x1000
	} else {
		// fine Y = 0
		ppu.v &= 0x8FFF
		// let y = coarse Y
		y := (ppu.v & 0x03E0) >> 5
		if y == 29 {
			// coarse Y = 0
			y = 0
			// switch vertical nametable
			ppu.v ^= 0x0800
		} else if y == 31 {
			// coarse Y = 0, nametable not switched
			y = 0
		} else {
			// increment coarse Y
			y += 1
			// put coarse Y back into v
			ppu.v = (ppu.v & 0xFC1F) | (y << 5)
		}
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
			ppu.f ^= 1
		}
	}
	if ppu.f == 1 && ppu.ScanLine == 0 && ppu.Cycle == 0 {
		ppu.Cycle++
	}
}

// Step executes a single PPU cycle
func (ppu *PPU) Step() {
	ppu.tick()

	// if rendering is enabled
	if ppu.flagShowBackground != 0 || ppu.flagShowSprites != 0 {
		// fetch tile data
		if ppu.ScanLine < 240 || ppu.ScanLine == 261 {
			if ppu.Cycle < 249 || (ppu.Cycle >= 321 && ppu.Cycle < 337) {
				if ppu.Cycle%2 == 1 {
				}
			}
		}

		// pre-render line
		if ppu.ScanLine == 261 {
			if ppu.Cycle >= 280 && ppu.Cycle <= 304 {
				// vert(v) = vert(t)
				// v: .IHGF.ED CBA..... = t: .IHGF.ED CBA.....
				ppu.v = (ppu.v & 0x841F) | (ppu.t & 0x7BE0)
			}
		}

		// pre-render and render lines
		if ppu.ScanLine < 240 || ppu.ScanLine == 261 {
			if (ppu.Cycle <= 256 || ppu.Cycle >= 328) && ppu.Cycle%8 == 0 {
				// increment hori(v)
				ppu.incrementX()
			}
			if ppu.Cycle == 256 {
				// increment vert(v)
				ppu.incrementY()
			}
			if ppu.Cycle == 257 {
				// hori(v) = hori(t)
				// v: .....F.. ...EDCBA = t: .....F.. ...EDCBA
				ppu.v = (ppu.v & 0xFBE0) | (ppu.t & 0x041F)
			}
		}
	}

	// tile address      = 0x2000 | (v & 0x0FFF)
	// attribute address = 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)

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
		y := ppu.oamPrimary[index+0]
		t := ppu.oamPrimary[index+1]
		f := ppu.oamPrimary[index+2]
		x := ppu.oamPrimary[index+3]
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
