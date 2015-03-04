package nes

type Memory interface {
	Read(address uint16) byte
	Read16(address uint16) uint16
	Write(address uint16, value byte)
}

type CPUMemory struct {
	NES *NES
	RAM []byte
}

func NewCPUMemory(nes *NES) Memory {
	ram := make([]byte, 2048)
	return &CPUMemory{nes, ram}
}

func (mem *CPUMemory) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		return mem.RAM[address%0x0800]
	case address >= 0x8000:
		return mem.NES.Cartridge.Read(address)
	}
	return 0
}

func (mem *CPUMemory) Read16(address uint16) uint16 {
	lo := uint16(mem.Read(address))
	hi := uint16(mem.Read(address + 1))
	return hi<<8 | lo
}

func (mem *CPUMemory) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		mem.RAM[address%0x0800] = value
	}
}
