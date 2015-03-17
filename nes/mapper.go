package nes

import "log"

type Mapper interface {
	Read(address uint16) byte
	Write(address uint16, value byte)
	Step()
}

func NewMapper(nes *NES, cartridge *Cartridge) Mapper {
	switch cartridge.Mapper {
	case 0:
		return NewMapper2(cartridge)
	case 1:
		return NewMapper1(cartridge)
	case 2:
		return NewMapper2(cartridge)
	case 4:
		return NewMapper4(nes, cartridge)
	default:
		log.Fatalf("unsupported mapper: %d", cartridge.Mapper)
	}
	return nil
}
