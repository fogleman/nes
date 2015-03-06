package nes

type NES struct {
	CPU       *CPU
	PPU       *PPU
	RAM       []byte
	Cartridge *Cartridge
}

func NewNES(path string) (*NES, error) {
	cartridge, err := LoadNESFile(path)
	if err != nil {
		return nil, err
	}
	ppu := NewPPU()
	ram := make([]byte, 2048)
	nes := NES{nil, ppu, ram, cartridge}
	nes.CPU = NewCPU(NewCPUMemory(&nes))
	return &nes, nil
}

func (nes *NES) Step() int {
	cpuCycles := nes.CPU.Step()
	for i := 0; i < cpuCycles*3; i++ {
		nes.PPU.Step()
	}
	return cpuCycles
}
