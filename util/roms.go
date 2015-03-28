package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/fogleman/nes/nes"
)

func testRom(path string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	console, err := nes.NewConsole(path)
	if err != nil {
		return err
	}
	console.StepSeconds(3)
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalln("Usage: go run util/roms.go roms_directory")
	}
	dir := args[0]
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, info := range infos {
		name := info.Name()
		if !strings.HasSuffix(name, ".nes") {
			continue
		}
		name = path.Join(dir, name)
		err := testRom(name)
		if err == nil {
			fmt.Println("OK  ", name)
		} else {
			fmt.Println("FAIL", name)
			fmt.Println(err)
		}
	}
}
