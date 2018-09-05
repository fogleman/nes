package util

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

func HandleCompressedFile(fileName string) (string, error) {
	tempFolder := CreateTempFolder()
	var rom string
	r, err := zip.OpenReader(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return rom, err
		}
		defer rc.Close()
		if strings.HasSuffix(f.Name, ".nes") {
			rom = path.Join(tempFolder, f.Name)

			outFile, err := os.OpenFile(rom, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return rom, err
			}
			_, err = io.Copy(outFile, rc)

			outFile.Close()
			if err != nil {
				return rom, err
			}
		}
	}
	return rom, nil
}

func CreateTempFolder() string {
	tempFolder := "tmp"

	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		os.Mkdir(tempFolder, os.ModePerm)
	}
	return tempFolder
}
