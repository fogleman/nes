package nes

import (
	"encoding/gob"
	"image"
)

type PPU struct {
	Memory           // memory interface
	console *Console // reference to parent object

	Cycle    int    // 0-340
	ScanLine int    // 0-261, 0-239=visible, 240=post, 241-260=vblank, 261=pre
	Frame    uint64 // frame counter

	// storage variables
	paletteData   [32]byte
	nameTableData [2048]byte
	oamData       [256]byte
	front         *image.RGBA
	back          *image.RGBA

	// PPU registers
	v uint16 // current vram address (15 bit)
	t uint16 // temporary vram address (15 bit)
	x byte   // fine x scroll (3 bit)
	w byte   // write toggle (1 bit)
	f byte   // even/odd frame flag (1 bit)

	register byte

	// NMI flags
	nmiOccurred bool
	nmiOutput   bool
	nmiPrevious bool
	nmiDelay    byte

	// background temporary variables
	nameTableByte      byte
	attributeTableByte byte
	lowTileByte        byte
	highTileByte       byte
	tileData           uint64

	// sprite temporary variables
	spriteCount      int
	spritePatterns   [8]uint32
	spritePositions  [8]byte
	spritePriorities [8]byte
	spriteIndexes    [8]byte

	// $2000 PPUCTRL
	flagNameTable       byte // 0: $2000; 1: $2400; 2: $2800; 3: $2C00
	flagIncrement       byte // 0: add 1; 1: add 32
	flagSpriteTable     byte // 0: $0000; 1: $1000; ignored in 8x16 mode
	flagBackgroundTable byte // 0: $0000; 1: $1000
	flagSpriteSize      byte // 0: 8x8; 1: 8x16
	flagMasterSlave     byte // 0: read EXT; 1: write EXT

	// $2001 PPUMASK
	flagGrayscale          byte // 0: color; 1: grayscale
	flagShowLeftBackground byte // 0: hide; 1: show
	flagShowLeftSprites    byte // 0: hide; 1: show
	flagShowBackground     byte // 0: hide; 1: show
	flagShowSprites        byte // 0: hide; 1: show
	flagRedTint            byte // 0: normal; 1: emphasized
	flagGreenTint          byte // 0: normal; 1: emphasized
	flagBlueTint           byte // 0: normal; 1: emphasized

	// $2002 PPUSTATUS
	flagSpriteZeroHit  byte
	flagSpriteOverflow byte

	// $2003 OAMADDR
	oamAddress byte

	// $2007 PPUDATA
	bufferedData byte // for buffered reads
}

func NewPPU(console *Console) *PPU {
	ppu := PPU{Memory: NewPPUMemory(console), console: console}
	ppu.front = image.NewRGBA(image.Rect(0, 0, 256, 240))
	ppu.back = image.NewRGBA(image.Rect(0, 0, 256, 240))
	ppu.Reset()
	return &ppu
}

func (ppu *PPU) Save(encoder *gob.Encoder) error {
	encoder.Encode(ppu.Cycle)
	encoder.Encode(ppu.ScanLine)
	encoder.Encode(ppu.Frame)
	encoder.Encode(ppu.paletteData)
	encoder.Encode(ppu.nameTableData)
	encoder.Encode(ppu.oamData)
	encoder.Encode(ppu.v)
	encoder.Encode(ppu.t)
	encoder.Encode(ppu.x)
	encoder.Encode(ppu.w)
	encoder.Encode(ppu.f)
	encoder.Encode(ppu.register)
	encoder.Encode(ppu.nmiOccurred)
	encoder.Encode(ppu.nmiOutput)
	encoder.Encode(ppu.nmiPrevious)
	encoder.Encode(ppu.nmiDelay)
	encoder.Encode(ppu.nameTableByte)
	encoder.Encode(ppu.attributeTableByte)
	encoder.Encode(ppu.lowTileByte)
	encoder.Encode(ppu.highTileByte)
	encoder.Encode(ppu.tileData)
	encoder.Encode(ppu.spriteCount)
	encoder.Encode(ppu.spritePatterns)
	encoder.Encode(ppu.spritePositions)
	encoder.Encode(ppu.spritePriorities)
	encoder.Encode(ppu.spriteIndexes)
	encoder.Encode(ppu.flagNameTable)
	encoder.Encode(ppu.flagIncrement)
	encoder.Encode(ppu.flagSpriteTable)
	encoder.Encode(ppu.flagBackgroundTable)
	encoder.Encode(ppu.flagSpriteSize)
	encoder.Encode(ppu.flagMasterSlave)
	encoder.Encode(ppu.flagGrayscale)
	encoder.Encode(ppu.flagShowLeftBackground)
	encoder.Encode(ppu.flagShowLeftSprites)
	encoder.Encode(ppu.flagShowBackground)
	encoder.Encode(ppu.flagShowSprites)
	encoder.Encode(ppu.flagRedTint)
	encoder.Encode(ppu.flagGreenTint)
	encoder.Encode(ppu.flagBlueTint)
	encoder.Encode(ppu.flagSpriteZeroHit)
	encoder.Encode(ppu.flagSpriteOverflow)
	encoder.Encode(ppu.oamAddress)
	encoder.Encode(ppu.bufferedData)
	return nil
}

func (ppu *PPU) Load(decoder *gob.Decoder) error {
	decoder.Decode(&ppu.Cycle)
	decoder.Decode(&ppu.ScanLine)
	decoder.Decode(&ppu.Frame)
	decoder.Decode(&ppu.paletteData)
	decoder.Decode(&ppu.nameTableData)
	decoder.Decode(&ppu.oamData)
	decoder.Decode(&ppu.v)
	decoder.Decode(&ppu.t)
	decoder.Decode(&ppu.x)
	decoder.Decode(&ppu.w)
	decoder.Decode(&ppu.f)
	decoder.Decode(&ppu.register)
	decoder.Decode(&ppu.nmiOccurred)
	decoder.Decode(&ppu.nmiOutput)
	decoder.Decode(&ppu.nmiPrevious)
	decoder.Decode(&ppu.nmiDelay)
	decoder.Decode(&ppu.nameTableByte)
	decoder.Decode(&ppu.attributeTableByte)
	decoder.Decode(&ppu.lowTileByte)
	decoder.Decode(&ppu.highTileByte)
	decoder.Decode(&ppu.tileData)
	decoder.Decode(&ppu.spriteCount)
	decoder.Decode(&ppu.spritePatterns)
	decoder.Decode(&ppu.spritePositions)
	decoder.Decode(&ppu.spritePriorities)
	decoder.Decode(&ppu.spriteIndexes)
	decoder.Decode(&ppu.flagNameTable)
	decoder.Decode(&ppu.flagIncrement)
	decoder.Decode(&ppu.flagSpriteTable)
	decoder.Decode(&ppu.flagBackgroundTable)
	decoder.Decode(&ppu.flagSpriteSize)
	decoder.Decode(&ppu.flagMasterSlave)
	decoder.Decode(&ppu.flagGrayscale)
	decoder.Decode(&ppu.flagShowLeftBackground)
	decoder.Decode(&ppu.flagShowLeftSprites)
	decoder.Decode(&ppu.flagShowBackground)
	decoder.Decode(&ppu.flagShowSprites)
	decoder.Decode(&ppu.flagRedTint)
	decoder.Decode(&ppu.flagGreenTint)
	decoder.Decode(&ppu.flagBlueTint)
	decoder.Decode(&ppu.flagSpriteZeroHit)
	decoder.Decode(&ppu.flagSpriteOverflow)
	decoder.Decode(&ppu.oamAddress)
	decoder.Decode(&ppu.bufferedData)
	return nil
}

func (ppu *PPU) Reset() {
	ppu.Cycle = 340
	ppu.ScanLine = 240
	ppu.Frame = 0
	ppu.writeControl(0)
	ppu.writeMask(0)
	ppu.writeOAMAddress(0)
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

func (ppu *PPU) readRegister(address uint16) byte {
	switch address {
	case 0x2002:
		return ppu.readStatus()
	case 0x2004:
		return ppu.readOAMData()
	case 0x2007:
		return ppu.readData()
	}
	return 0
}

func (ppu *PPU) writeRegister(address uint16, value byte) {
	ppu.register = value
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
	ppu.nmiOutput = (value>>7)&1 == 1
	ppu.nmiChange()
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
	result := ppu.register & 0x1F
	result |= ppu.flagSpriteOverflow << 5
	result |= ppu.flagSpriteZeroHit << 6
	if ppu.nmiOccurred {
		result |= 1 << 7
	}
	ppu.nmiOccurred = false
	ppu.nmiChange()
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
	data := ppu.oamData[ppu.oamAddress]
	if (ppu.oamAddress & 0x03) == 0x02 {
		data = data & 0xE3
	}
	return data
}

// $2004: OAMDATA (write)
func (ppu *PPU) writeOAMData(value byte) {
	ppu.oamData[ppu.oamAddress] = value
	ppu.oamAddress++
}

// $2005: PPUSCROLL
func (ppu *PPU) writeScroll(value byte) {
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
		ppu.t = (ppu.t & 0x8FFF) | ((uint16(value) & 0x07) << 12)
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
		buffered := ppu.bufferedData
		ppu.bufferedData = value
		value = buffered
	} else {
		ppu.bufferedData = ppu.Read(ppu.v - 0x1000)
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
	cpu := ppu.console.CPU
	address := uint16(value) << 8
	for i := 0; i < 256; i++ {
		ppu.oamData[ppu.oamAddress] = cpu.Read(address)
		ppu.oamAddress++
		address++
	}
	cpu.stall += 513
	if cpu.Cycles%2 == 1 {
		cpu.stall++
	}
}

// NTSC Timing Helper Functions

func (ppu *PPU) incrementX() {
	// increment hori(v)
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
	// increment vert(v)
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
			y++
		}
		// put coarse Y back into v
		ppu.v = (ppu.v & 0xFC1F) | (y << 5)
	}
}

func (ppu *PPU) copyX() {
	// hori(v) = hori(t)
	// v: .....F.. ...EDCBA = t: .....F.. ...EDCBA
	ppu.v = (ppu.v & 0xFBE0) | (ppu.t & 0x041F)
}

func (ppu *PPU) copyY() {
	// vert(v) = vert(t)
	// v: .IHGF.ED CBA..... = t: .IHGF.ED CBA.....
	ppu.v = (ppu.v & 0x841F) | (ppu.t & 0x7BE0)
}

func (ppu *PPU) nmiChange() {
	nmi := ppu.nmiOutput && ppu.nmiOccurred
	if nmi && !ppu.nmiPrevious {
		// TODO: this fixes some games but the delay shouldn't have to be so
		// long, so the timings are off somewhere
		ppu.nmiDelay = 15
	}
	ppu.nmiPrevious = nmi
}

func (ppu *PPU) setVerticalBlank() {
	ppu.front, ppu.back = ppu.back, ppu.front
	ppu.nmiOccurred = true
	ppu.nmiChange()
}

func (ppu *PPU) clearVerticalBlank() {
	ppu.nmiOccurred = false
	ppu.nmiChange()
}

func (ppu *PPU) fetchNameTableByte() {
	v := ppu.v
	address := 0x2000 | (v & 0x0FFF)
	ppu.nameTableByte = ppu.Read(address)
}

func (ppu *PPU) fetchAttributeTableByte() {
	v := ppu.v
	address := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
	shift := ((v >> 4) & 4) | (v & 2)
	ppu.attributeTableByte = ((ppu.Read(address) >> shift) & 3) << 2
}

func (ppu *PPU) fetchLowTileByte() {
	fineY := (ppu.v >> 12) & 7
	table := ppu.flagBackgroundTable
	tile := ppu.nameTableByte
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	ppu.lowTileByte = ppu.Read(address)
}

func (ppu *PPU) fetchHighTileByte() {
	fineY := (ppu.v >> 12) & 7
	table := ppu.flagBackgroundTable
	tile := ppu.nameTableByte
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	ppu.highTileByte = ppu.Read(address + 8)
}

func (ppu *PPU) storeTileData() {
	var data uint32
	for i := 0; i < 8; i++ {
		a := ppu.attributeTableByte
		p1 := (ppu.lowTileByte & 0x80) >> 7
		p2 := (ppu.highTileByte & 0x80) >> 6
		ppu.lowTileByte <<= 1
		ppu.highTileByte <<= 1
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	ppu.tileData |= uint64(data)
}

func (ppu *PPU) fetchTileData() uint32 {
	return uint32(ppu.tileData >> 32)
}

func (ppu *PPU) backgroundPixel() byte {
	if ppu.flagShowBackground == 0 {
		return 0
	}
	data := ppu.fetchTileData() >> ((7 - ppu.x) * 4)
	return byte(data & 0x0F)
}

func (ppu *PPU) spritePixel() (byte, byte) {
	if ppu.flagShowSprites == 0 {
		return 0, 0
	}
	for i := 0; i < ppu.spriteCount; i++ {
		offset := (ppu.Cycle - 1) - int(ppu.spritePositions[i])
		if offset < 0 || offset > 7 {
			continue
		}
		offset = 7 - offset
		color := byte((ppu.spritePatterns[i] >> byte(offset*4)) & 0x0F)
		if color%4 == 0 {
			continue
		}
		return byte(i), color
	}
	return 0, 0
}

func (ppu *PPU) renderPixel() {
	x := ppu.Cycle - 1
	y := ppu.ScanLine
	background := ppu.backgroundPixel()
	i, sprite := ppu.spritePixel()
	if x < 8 && ppu.flagShowLeftBackground == 0 {
		background = 0
	}
	if x < 8 && ppu.flagShowLeftSprites == 0 {
		sprite = 0
	}
	b := background%4 != 0
	s := sprite%4 != 0
	var color byte
	if !b && !s {
		color = 0
	} else if !b && s {
		color = sprite | 0x10
	} else if b && !s {
		color = background
	} else {
		if ppu.spriteIndexes[i] == 0 && x < 255 {
			ppu.flagSpriteZeroHit = 1
		}
		if ppu.spritePriorities[i] == 0 {
			color = sprite | 0x10
		} else {
			color = background
		}
	}
	c := Palette[ppu.readPalette(uint16(color))%64]
	ppu.back.SetRGBA(x, y, c)
}

func (ppu *PPU) fetchSpritePattern(i, row int) uint32 {
	tile := ppu.oamData[i*4+1]
	attributes := ppu.oamData[i*4+2]
	var address uint16
	if ppu.flagSpriteSize == 0 {
		if attributes&0x80 == 0x80 {
			row = 7 - row
		}
		table := ppu.flagSpriteTable
		address = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	} else {
		if attributes&0x80 == 0x80 {
			row = 15 - row
		}
		table := tile & 1
		tile &= 0xFE
		if row > 7 {
			tile++
			row -= 8
		}
		address = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	}
	a := (attributes & 3) << 2
	lowTileByte := ppu.Read(address)
	highTileByte := ppu.Read(address + 8)
	var data uint32
	for i := 0; i < 8; i++ {
		var p1, p2 byte
		if attributes&0x40 == 0x40 {
			p1 = (lowTileByte & 1) << 0
			p2 = (highTileByte & 1) << 1
			lowTileByte >>= 1
			highTileByte >>= 1
		} else {
			p1 = (lowTileByte & 0x80) >> 7
			p2 = (highTileByte & 0x80) >> 6
			lowTileByte <<= 1
			highTileByte <<= 1
		}
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	return data
}

func (ppu *PPU) evaluateSprites() {
	var h int
	if ppu.flagSpriteSize == 0 {
		h = 8
	} else {
		h = 16
	}
	count := 0
	for i := 0; i < 64; i++ {
		y := ppu.oamData[i*4+0]
		a := ppu.oamData[i*4+2]
		x := ppu.oamData[i*4+3]
		row := ppu.ScanLine - int(y)
		if row < 0 || row >= h {
			continue
		}
		if count < 8 {
			ppu.spritePatterns[count] = ppu.fetchSpritePattern(i, row)
			ppu.spritePositions[count] = x
			ppu.spritePriorities[count] = (a >> 5) & 1
			ppu.spriteIndexes[count] = byte(i)
		}
		count++
	}
	if count > 8 {
		count = 8
		ppu.flagSpriteOverflow = 1
	}
	ppu.spriteCount = count
}

// tick updates Cycle, ScanLine and Frame counters
func (ppu *PPU) tick() {
	if ppu.nmiDelay > 0 {
		ppu.nmiDelay--
		if ppu.nmiDelay == 0 && ppu.nmiOutput && ppu.nmiOccurred {
			ppu.console.CPU.triggerNMI()
		}
	}

	if ppu.flagShowBackground != 0 || ppu.flagShowSprites != 0 {
		if ppu.f == 1 && ppu.ScanLine == 261 && ppu.Cycle == 339 {
			ppu.Cycle = 0
			ppu.ScanLine = 0
			ppu.Frame++
			ppu.f ^= 1
			return
		}
	}
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
}

// Step executes a single PPU cycle
func (ppu *PPU) Step() {
	ppu.tick()

	renderingEnabled := ppu.flagShowBackground != 0 || ppu.flagShowSprites != 0
	preLine := ppu.ScanLine == 261
	visibleLine := ppu.ScanLine < 240
	// postLine := ppu.ScanLine == 240
	renderLine := preLine || visibleLine
	preFetchCycle := ppu.Cycle >= 321 && ppu.Cycle <= 336
	visibleCycle := ppu.Cycle >= 1 && ppu.Cycle <= 256
	fetchCycle := preFetchCycle || visibleCycle

	// background logic
	if renderingEnabled {
		if visibleLine && visibleCycle {
			ppu.renderPixel()
		}
		if renderLine && fetchCycle {
			ppu.tileData <<= 4
			switch ppu.Cycle % 8 {
			case 1:
				ppu.fetchNameTableByte()
			case 3:
				ppu.fetchAttributeTableByte()
			case 5:
				ppu.fetchLowTileByte()
			case 7:
				ppu.fetchHighTileByte()
			case 0:
				ppu.storeTileData()
			}
		}
		if preLine && ppu.Cycle >= 280 && ppu.Cycle <= 304 {
			ppu.copyY()
		}
		if renderLine {
			if fetchCycle && ppu.Cycle%8 == 0 {
				ppu.incrementX()
			}
			if ppu.Cycle == 256 {
				ppu.incrementY()
			}
			if ppu.Cycle == 257 {
				ppu.copyX()
			}
		}
	}

	// sprite logic
	if renderingEnabled {
		if ppu.Cycle == 257 {
			if visibleLine {
				ppu.evaluateSprites()
			} else {
				ppu.spriteCount = 0
			}
		}
	}

	// vblank logic
	if ppu.ScanLine == 241 && ppu.Cycle == 1 {
		ppu.setVerticalBlank()
	}
	if preLine && ppu.Cycle == 1 {
		ppu.clearVerticalBlank()
		ppu.flagSpriteZeroHit = 0
		ppu.flagSpriteOverflow = 0
	}
}
