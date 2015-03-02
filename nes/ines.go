package nes

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

type iNESHeader struct {
	Magic    [4]byte
	NumPRG   byte
	NumCHR   byte
	Control1 byte
	Control2 byte
	NumRAM   byte
	_        [7]byte
}

func LoadNESFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	header := iNESHeader{}
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return err
	}
	if header.Magic != [4]byte{78, 69, 83, 26} {
		return errors.New("invalid .nes file header")
	}
	fmt.Println(header)
	return nil
}
