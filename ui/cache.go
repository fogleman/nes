package ui

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
)

const baseThumbnailURL = "http://www.michaelfogleman.com/static/nes/"

var homeDir string

func init() {
	u, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	homeDir = u.HomeDir
}

func EnsureThumbnail(path string) {
	hash, err := hashFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	LoadThumbnail(hash)
}

func LoadThumbnail(hash string) {
	path := homeDir + "/.nes/thumbnail/" + hash + ".png"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		DownloadThumbnail(hash) // TODO: goroutine
	} else {
		fmt.Println("exists")
	}
}

func DownloadThumbnail(hash string) error {
	dir := homeDir + "/.nes/thumbnail/"
	path := dir + hash + ".png"
	url := baseThumbnailURL + hash + ".png"

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}
