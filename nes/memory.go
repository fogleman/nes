package nes

import "log"

type Memory interface {
	Read(address uint16) byte
	Write(address uint16, value byte)
}

// CPU Memory Map

type cpuMemory struct {
	console *Console
}

func NewCPUMemory(console *Console) Memory {
	return &cpuMemory{console}
}

func (mem *cpuMemory) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		return mem.console.RAM[address%0x0800]
	case address < 0x4000:
		return mem.console.PPU.readRegister(0x2000 + address%8)
	case address == 0x4014:
		return mem.console.PPU.readRegister(address)
	case address == 0x4015:
		return mem.console.APU.readRegister(address)
	case address == 0x4016:
		return mem.console.Controller1.Read()
	case address == 0x4017:
		return mem.console.Controller2.Read()
	case address < 0x6000:
		// TODO: I/O registers
	case address >= 0x6000:
		return mem.console.Mapper.Read(address)
	default:
		log.Fatalf("unhandled cpu memory read at address: 0x%04X", address)
	}
	return 0
}

func (mem *cpuMemory) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		mem.console.RAM[address%0x0800] = value
	case address < 0x4000:
		mem.console.PPU.writeRegister(0x2000+address%8, value)
	case address < 0x4014:
		mem.console.APU.writeRegister(address, value)
	case address == 0x4014:
		mem.console.PPU.writeRegister(address, value)
	case address == 0x4015:
		mem.console.APU.writeRegister(address, value)
	case address == 0x4016:
		mem.console.Controller1.Write(value)
		mem.console.Controller2.Write(value)
	case address == 0x4017:
		mem.console.APU.writeRegister(address, value)
	case address < 0x6000:
		// TODO: I/O registers
	case address >= 0x6000:
		mem.console.Mapper.Write(address, value)
	default:
		log.Fatalf("unhandled cpu memory write at address: 0x%04X", address)
	}
}

// PPU Memory Map

type ppuMemory struct {
	console *Console
}

func NewPPUMemory(console *Console) Memory {
	return &ppuMemory{console}
}

func (mem *ppuMemory) Read(address uint16) byte {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		return mem.console.Mapper.Read(address)
	case address < 0x3F00:
		mode := mem.console.Cartridge.Mirror
		return mem.console.PPU.nameTableData[MirrorAddress(mode, address)%2048]
	case address < 0x4000:
		return mem.console.PPU.readPalette(address % 32)
	default:
		log.Fatalf("unhandled ppu memory read at address: 0x%04X", address)
	}
	return 0
}

func (mem *ppuMemory) Write(address uint16, value byte) {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		mem.console.Mapper.Write(address, value)
	case address < 0x3F00:
		mode := mem.console.Cartridge.Mirror
		mem.console.PPU.nameTableData[MirrorAddress(mode, address)%2048] = value
	case address < 0x4000:
		mem.console.PPU.writePalette(address%32, value)
	default:
		log.Fatalf("unhandled ppu memory write at address: 0x%04X", address)
	}
}

// Mirroring Modes

const (
	MirrorHorizontal = 0
	MirrorVertical   = 1
	MirrorSingle0    = 2
	MirrorSingle1    = 3
	MirrorFour       = 4
)

var MirrorLookup = [...][4]uint16{
	{0, 0, 1, 1},
	{0, 1, 0, 1},
	{0, 0, 0, 0},
	{1, 1, 1, 1},
	{0, 1, 2, 3},
}

func MirrorAddress(mode byte, address uint16) uint16 {
	address = (address - 0x2000) % 0x1000
	table := address / 0x0400
	offset := address % 0x0400
	return 0x2000 + MirrorLookup[mode][table]*0x0400 + offset
}
