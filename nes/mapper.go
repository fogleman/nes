package nes

import "log"

func NewMapper(cartridge *Cartridge) Memory {
	switch cartridge.Mapper {
	case 0:
		return NewMapper2(cartridge)
	case 2:
		return NewMapper2(cartridge)
	default:
		log.Fatalf("unsupported mapper: %d", cartridge.Mapper)
	}
	return nil
}
