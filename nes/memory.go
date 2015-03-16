package nes

import "log"

type Memory interface {
	Read(address uint16) byte
	Write(address uint16, value byte)
}

// CPU Memory Map

type cpuMemory struct {
	nes *NES
}

func NewCPUMemory(nes *NES) Memory {
	return &cpuMemory{nes}
}

func (mem *cpuMemory) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		return mem.nes.RAM[address%0x0800]
	case address < 0x4000:
		return mem.nes.PPU.readRegister(0x2000 + address%8)
	case address == 0x4014:
		return mem.nes.PPU.readRegister(address)
	case address == 0x4016:
		return mem.nes.Controller1.Read()
	case address == 0x4017:
		return mem.nes.Controller2.Read()
	case address >= 0x6000:
		return mem.nes.Mapper.Read(address)
	default:
		log.Fatalf("unhandled cpu memory read at address: 0x%04X", address)
	}
	return 0
}

func (mem *cpuMemory) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		mem.nes.RAM[address%0x0800] = value
	case address < 0x4000:
		mem.nes.PPU.writeRegister(0x2000+address%8, value)
	case address == 0x4014:
		mem.nes.PPU.writeRegister(address, value)
	case address == 0x4016:
		mem.nes.Controller1.Write(value)
	case address == 0x4017:
		mem.nes.Controller2.Write(value)
	case address < 0x4020:
		// TODO: I/O registers
	case address >= 0x6000:
		mem.nes.Mapper.Write(address, value)
	default:
		log.Fatalf("unhandled cpu memory write at address: 0x%04X", address)
	}
}

// PPU Memory Map

var MirrorLookup = [2][4]uint16{
	{0, 0, 1, 1},
	{0, 1, 0, 1},
}

type ppuMemory struct {
	nes *NES
}

func NewPPUMemory(nes *NES) Memory {
	return &ppuMemory{nes}
}

func (mem *ppuMemory) Read(address uint16) byte {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		return mem.nes.Mapper.Read(address)
	case address < 0x3F00:
		return mem.nes.PPU.nameTableData[mem.MirrorAddress(address)%2048]
	case address < 0x4000:
		return mem.nes.PPU.readPalette(address % 32)
	default:
		log.Fatalf("unhandled ppu memory read at address: 0x%04X", address)
	}
	return 0
}

func (mem *ppuMemory) Write(address uint16, value byte) {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		mem.nes.Mapper.Write(address, value)
	case address < 0x3F00:
		mem.nes.PPU.nameTableData[mem.MirrorAddress(address)%2048] = value
	case address < 0x4000:
		mem.nes.PPU.writePalette(address%32, value)
	default:
		log.Fatalf("unhandled ppu memory write at address: 0x%04X", address)
	}
}

func (mem *ppuMemory) MirrorAddress(address uint16) uint16 {
	cartridge := mem.nes.Cartridge
	address = (address - 0x2000) % 0x1000
	table := address / 0x0400
	offset := address % 0x0400
	return 0x2000 + MirrorLookup[cartridge.Mirror][table]*0x0400 + offset
}
