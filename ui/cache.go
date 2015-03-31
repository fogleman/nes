package ui

import (
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
)

type Cache struct {
	homeDir string
	ch      chan string
}

func NewCache() *Cache {
	u, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	cache := Cache{}
	cache.homeDir = u.HomeDir
	cache.ch = make(chan string, 1024)
	return &cache
}

func (c *Cache) LoadThumbnail(path string) image.Image {
	im := CreateGenericThumbnail(path)
	hash, err := hashFile(path)
	if err != nil {
		return im
	}
	thumbnailPath := c.homeDir + "/.nes/thumbnail/" + hash + ".png"
	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		go c.downloadThumbnail(path, hash)
		return im
	} else {
		thumbnail, err := loadPNG(thumbnailPath)
		if err != nil {
			return im
		}
		return thumbnail
	}
}

func (c *Cache) downloadThumbnail(path, hash string) error {
	dir := c.homeDir + "/.nes/thumbnail/"
	thumbnailPath := dir + hash + ".png"
	url := "http://www.michaelfogleman.com/static/nes/" + hash + ".png"

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(thumbnailPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	c.ch <- path

	return nil
}
