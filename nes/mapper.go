package nes

import "log"

type Mapper interface {
	Read(address uint16) byte
	Write(address uint16, value byte)
	Step()
}

func NewMapper(console *Console) Mapper {
	cartridge := console.Cartridge
	switch cartridge.Mapper {
	case 0:
		return NewMapper2(cartridge)
	case 1:
		return NewMapper1(cartridge)
	case 2:
		return NewMapper2(cartridge)
	case 3:
		return NewMapper3(cartridge)
	case 4:
		return NewMapper4(console, cartridge)
	case 7:
		return NewMapper7(cartridge)
	default:
		log.Fatalf("unsupported mapper: %d", cartridge.Mapper)
	}
	return nil
}
