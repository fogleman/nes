package nes

import "image"

type NES struct {
	CPU       *CPU
	PPU       *PPU
	CPUMemory Memory
	PPUMemory Memory
	RAM       []byte
	Cartridge *Cartridge
}

func NewNES(path string) (*NES, error) {
	cartridge, err := LoadNESFile(path)
	if err != nil {
		return nil, err
	}
	ram := make([]byte, 2048)
	nes := NES{nil, nil, nil, nil, ram, cartridge}
	nes.CPUMemory = NewCPUMemory(&nes)
	nes.PPUMemory = NewPPUMemory(&nes)
	nes.CPU = NewCPU(&nes)
	nes.PPU = NewPPU(&nes)
	return &nes, nil
}

func (nes *NES) Step() int {
	nes.CPU.PrintInstruction()
	cpuCycles := nes.CPU.Step()
	ppuCycles := cpuCycles * 3
	for i := 0; i < ppuCycles; i++ {
		nes.PPU.Step()
	}
	return cpuCycles
}

func (nes *NES) StepFrame() {
	frame := nes.PPU.Frame
	for frame == nes.PPU.Frame {
		nes.Step()
	}
}

func (nes *NES) Buffer() *image.RGBA {
	return nes.PPU.buffer
}
