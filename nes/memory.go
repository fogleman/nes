package nes

import "log"

type Memory interface {
	Read(address uint16) byte
	Read16(address uint16) uint16
	Write(address uint16, value byte)
}

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
	case address >= 0x6000:
		return mem.nes.Cartridge.Read(address)
	default:
		log.Fatalf("unhandled read at address: 0x%04X", address)
	}
	return 0
}

func (mem *cpuMemory) Read16(address uint16) uint16 {
	lo := uint16(mem.Read(address))
	hi := uint16(mem.Read(address + 1))
	return hi<<8 | lo
}

func (mem *cpuMemory) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		mem.nes.RAM[address%0x0800] = value
	case address >= 0x6000:
		mem.nes.Cartridge.Write(address, value)
	default:
		log.Fatalf("unhandled write at address: 0x%04X", address)
	}
}
