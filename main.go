package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"fmt"

	"github.com/fogleman/nes/ui"
)

func main() {
	log.SetFlags(0)
	args, fullscreen := extractFlags()
	paths := getPaths(args)
	if len(paths) == 0 {
		log.Fatalln("no rom files specified or found")
	}
	ui.Run(paths, fullscreen)
}

func extractFlags() ([]string, bool) {
	args := os.Args

	for i, x := range args {
		//Probably should add a keyboard shortcut for this also
		if x == "--fullscreen" {
			a := append(args[:i], args[i+1:]...)
			fmt.Printf("fullscreen - args - %v -- a %b", args, a)
			return a, true
		}
	}

	return args, false
}

func getPaths(argsv []string ) []string {
	var arg string
	args := argsv[1:]
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
			if !strings.HasSuffix(name, ".nes") {
				continue
			}
			result = append(result, path.Join(arg, name))
		}
		return result
	} else {
		return []string{arg}
	}
}
