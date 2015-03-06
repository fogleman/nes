package nes

type NES struct {
	CPU       *CPU
	RAM       []byte
	Cartridge *Cartridge
}

func NewNES(path string) (*NES, error) {
	cartridge, err := LoadNESFile(path)
	if err != nil {
		return nil, err
	}
	ram := make([]byte, 2048)
	nes := NES{nil, ram, cartridge}
	nes.CPU = NewCPU(NewCPUMemory(&nes))
	return &nes, nil
}
