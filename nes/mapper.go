package nes

import "log"

type Mapper interface {
	Read(address uint16) byte
	Write(address uint16, value byte)
}

func NewMapper(cartridge *Cartridge) Mapper {
	switch cartridge.Mapper {
	case 0:
		return NewMapper2(cartridge)
	case 1:
		return NewMapper1(cartridge)
	case 2:
		return NewMapper2(cartridge)
	case 4:
		return NewMapper4(cartridge)
	default:
		log.Fatalf("unsupported mapper: %d", cartridge.Mapper)
	}
	return nil
}
