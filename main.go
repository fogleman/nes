package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/fogleman/nes/ui"
	"github.com/fogleman/nes/util"
)

func main() {
	log.SetFlags(0)
	paths := getPaths()
	if len(paths) == 0 {
		log.Fatalln("no rom files specified or found")
	}
	ui.Run(paths)
}

func getPaths() []string {
	var arg string
	args := os.Args[1:]
	if len(args) == 1 {
		arg = args[0]
	} else {
		arg, _ = os.Getwd()
	}
	info, err := os.Stat(arg)
	if err != nil {
		return nil
	}
	if info.IsDir() {
		infos, err := ioutil.ReadDir(arg)
		if err != nil {
			return nil
		}
		var result []string
		for _, info := range infos {
			name := info.Name()
			name, _ = util.HandleZip(name)
			if !strings.HasSuffix(name, ".nes") {
				continue
			}
			result = append(result, path.Join(arg, name))
		}
		return result
	} else {
		arg, _ = util.HandleZip(arg)
		return []string{arg}
	}
}
