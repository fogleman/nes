package nes

type APU struct {
}

func NewAPU() *APU {
	return &APU{}
}

func (apu *APU) readRegister(address uint16) byte {
	switch address {
	case 0x4015:
		return apu.readStatus()
		// default:
		// 	log.Fatalf("unhandled apu register read at address: 0x%04X", address)
	}
	return 0
}

func (apu *APU) writeRegister(address uint16, value byte) {
	switch address {
	case 0x4015:
		apu.writeControl(value)
		// default:
		// 	log.Fatalf("unhandled apu register write at address: 0x%04X", address)
	}
}

// $4015: Status
func (apu *APU) readStatus() byte {
	return 0
}

// $4015: Control (---D NT21)
func (apu *APU) writeControl(value byte) {
}
