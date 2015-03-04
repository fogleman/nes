package nes

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

const NESFileMagic = 0x1a53454e

type NESFileHeader struct {
	Magic    uint32  // iNES magic number
	NumPRG   byte    // number of PRG-ROM banks (16KB each)
	NumCHR   byte    // number of CHR-ROM banks (8KB each)
	Control1 byte    // control bits
	Control2 byte    // control bits
	NumRAM   byte    // PRG-RAM size (x 8KB)
	_        [7]byte // unused padding
}

// LoadNESFile reads the iNES file format and returns a Cartridge on success.
// http://wiki.nesdev.com/w/index.php/INES
// http://nesdev.com/NESDoc.pdf (page 28)
func LoadNESFile(path string) (*Cartridge, error) {
	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// read file header
	header := NESFileHeader{}
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	// verify header magic number
	if header.Magic != NESFileMagic {
		return nil, errors.New("invalid .nes file")
	}

	// mapper type
	mapper1 := int(header.Control1) >> 4
	mapper2 := int(header.Control2) >> 4
	mapper := mapper1 | mapper2<<4

	// mirroring type
	mirror1 := int(header.Control1) & 1
	mirror2 := int(header.Control1>>3) & 1
	mirror := mirror1 | mirror2<<1

	// battery-backed RAM
	battery := header.Control1&2 == 2

	// read trainer if present (unused)
	if header.Control1&4 == 4 {
		trainer := make([]byte, 512)
		if _, err := io.ReadFull(file, trainer); err != nil {
			return nil, err
		}
	}

	// read prg-rom bank(s)
	prg := make([]byte, int(header.NumPRG)*16384)
	if _, err := io.ReadFull(file, prg); err != nil {
		return nil, err
	}

	// read chr-rom bank(s)
	chr := make([]byte, int(header.NumCHR)*8192)
	if _, err := io.ReadFull(file, chr); err != nil {
		return nil, err
	}

	// success
	cartridge := Cartridge{prg, chr, mapper, mirror, battery}
	return &cartridge, nil
}
