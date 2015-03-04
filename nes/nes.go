package nes

type NES struct {
	CPU       *CPU
	Cartridge *Cartridge
}

func NewNES(path string) (*NES, error) {
	cartridge, err := LoadNESFile(path)
	if err != nil {
		return nil, err
	}
	nes := NES{nil, cartridge}
	nes.CPU = NewCPU(NewCPUMemory(&nes))
	return &nes, nil
}
